package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"ingestWorker/internal/config"
	"ingestWorker/internal/models"
	"ingestWorker/internal/parser"
	"ingestWorker/internal/parser/docx"
	filestorage "ingestWorker/internal/storage/fileStorage"
	miniostorage "ingestWorker/internal/storage/fileStorage/minio"
	"ingestWorker/internal/storage/repository"
	postgresrepo "ingestWorker/internal/storage/repository/postgresql"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Service struct {
	repo         repository.Repository
	storage      filestorage.FileStorage
	parsers      []parser.DocumentParser
	pollInterval time.Duration
}

func Run(cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := openDB(cfg.DBDSN)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("db close failed: %v", closeErr)
		}
	}()

	storage, err := openMinIO(cfg)
	if err != nil {
		return err
	}

	service := &Service{
		repo:         postgresrepo.NewPostgresRepository(db),
		storage:      storage,
		parsers:      []parser.DocumentParser{&docx.DocxParser{}},
		pollInterval: cfg.PollInterval,
	}

	log.Printf("ingest worker started: concurrency=%d poll_interval=%s", cfg.WorkerConcurrency, cfg.PollInterval)
	service.RunPool(ctx, cfg.WorkerConcurrency)
	log.Printf("ingest worker stopped")

	return nil
}

func openDB(dsn string) (*pg.DB, error) {
	options, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DB_DSN: %w", err)
	}

	db := pg.Connect(options)
	if err := db.Ping(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	return db, nil
}

func openMinIO(cfg *config.Config) (filestorage.FileStorage, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("connect minio: %w", err)
	}

	return miniostorage.NewMiniOStorage(client, cfg.MinioBucket), nil
}

func (s *Service) RunPool(ctx context.Context, concurrency int) {
	var wg sync.WaitGroup

	for workerID := 1; workerID <= concurrency; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s.runWorker(ctx, id)
		}(workerID)
	}

	<-ctx.Done()
	log.Printf("shutdown signal received, waiting for workers")
	wg.Wait()
}

func (s *Service) runWorker(ctx context.Context, workerID int) {
	log.Printf("worker %d started", workerID)

	for {
		if err := ctx.Err(); err != nil {
			log.Printf("worker %d stopping: %v", workerID, err)
			return
		}

		job, submission, err := s.repo.ClaimNextPendingJob(ctx)
		if err != nil {
			if errors.Is(err, repository.ErrNoPendingJobs) {
				s.sleep(ctx)
				continue
			}

			log.Printf("worker %d: claim job failed: %v", workerID, err)
			s.sleep(ctx)
			continue
		}

		log.Printf("worker %d: job claimed job_id=%s submission_id=%s", workerID, job.ID, submission.ID)

		jobCtx := context.WithoutCancel(ctx)
		if err := s.processJob(jobCtx, workerID, job, submission); err != nil {
			log.Printf("worker %d: job failed job_id=%s submission_id=%s err=%v", workerID, job.ID, submission.ID, err)
			if markErr := s.repo.MarkJobFailed(jobCtx, job.ID, submission.ID, sanitizeErrorMessage(err)); markErr != nil {
				log.Printf("worker %d: mark failed also failed job_id=%s submission_id=%s err=%v", workerID, job.ID, submission.ID, markErr)
			}
			continue
		}

		log.Printf("worker %d: job done job_id=%s submission_id=%s", workerID, job.ID, submission.ID)
	}
}

func (s *Service) processJob(ctx context.Context, workerID int, job *models.IngestJob, submission *models.Submission) error {
	selectedParser := selectParser(submission.MimeType, s.parsers)
	if selectedParser == nil {
		return fmt.Errorf("no parser registered for mime type %s", submission.MimeType)
	}

	log.Printf("worker %d: parsing started job_id=%s submission_id=%s storage_key=%s", workerID, job.ID, submission.ID, submission.StorageKey)

	fileBytes, err := s.storage.GetFile(ctx, submission.StorageKey)
	if err != nil {
		return fmt.Errorf("download file from minio: %w", err)
	}

	tempPath, err := writeTempSubmissionFile(submission, fileBytes)
	if err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	defer func() {
		if removeErr := os.Remove(tempPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			log.Printf("remove temp file failed path=%s err=%v", tempPath, removeErr)
		}
	}()

	parseResult, err := selectedParser.ParseDocument(ctx, tempPath)
	if err != nil {
		return fmt.Errorf("parse document: %w", err)
	}

	parsed := models.ParsedSubmission{
		ID:             uuid.NewString(),
		SubmissionID:   submission.ID,
		RawText:        parseResult.RawText,
		StructuredData: json.RawMessage(parseResult.StructuredData),
		ParsedAt:       time.Now(),
	}

	if err := s.repo.MarkJobDone(ctx, job.ID, parsed); err != nil {
		return fmt.Errorf("mark job done: %w", err)
	}

	log.Printf("worker %d: parsing finished job_id=%s submission_id=%s", workerID, job.ID, submission.ID)

	if err := s.storage.DeleteFile(ctx, submission.StorageKey); err != nil {
		log.Printf("worker %d: minio cleanup failed job_id=%s submission_id=%s storage_key=%s err=%v", workerID, job.ID, submission.ID, submission.StorageKey, err)
	}

	return nil
}

func selectParser(mimeType string, parsers []parser.DocumentParser) parser.DocumentParser {
	for _, documentParser := range parsers {
		if documentParser.SupportMimeType(mimeType) {
			return documentParser
		}
	}

	return nil
}

func writeTempSubmissionFile(submission *models.Submission, data []byte) (string, error) {
	suffix := tempFileSuffix(submission)
	file, err := os.CreateTemp("", "ingestworker-*"+suffix)
	if err != nil {
		return "", err
	}

	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return "", err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

func tempFileSuffix(submission *models.Submission) string {
	if submission == nil {
		return ".tmp"
	}

	if ext := strings.TrimSpace(filepath.Ext(submission.SourceFileName)); ext != "" {
		return ext
	}

	switch strings.TrimSpace(strings.ToLower(submission.MimeType)) {
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/docx":
		return ".docx"
	default:
		return ".tmp"
	}
}

func sanitizeErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	message := strings.TrimSpace(err.Error())
	if len(message) <= 2048 {
		return message
	}

	return message[:2048]
}

func (s *Service) sleep(ctx context.Context) {
	timer := time.NewTimer(s.pollInterval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
