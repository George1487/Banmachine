package postgresql

import (
	"context"
	"errors"
	"fmt"
	"ingestWorker/internal/models"
	"ingestWorker/internal/storage/entity"
	"ingestWorker/internal/storage/repository"
	"time"

	"github.com/go-pg/pg/v10"
)

const (
	setStatus          = "status = ?"
	setStartedAt       = "started_at = ?"
	setFinishedAt      = "finished_at = ?"
	setErrorMessage    = "error_message = ?"
	setErrorMessageNil = "error_message = NULL"
	whereSubmissionID  = "submission_id = ?"
	whereIngestJobID   = "ingest_job_id = ?"
)

var _ repository.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *pg.DB
}

func NewPostgresRepository(db *pg.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ClaimNextPendingJob(ctx context.Context) (*models.IngestJob, *models.Submission, error) {
	returnIngestJob := new(models.IngestJob)
	returnSubmission := new(models.Submission)

	err := r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		jobRow := new(entity.IngestJobEntity)

		_, err := tx.QueryOneContext(ctx, jobRow,
			`
		SELECT * FROM ingest_jobs
		WHERE status = ?
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT 1
		`,
			string(models.IngestJobStatusPending),
		)
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return repository.ErrNoPendingJobs
			}
			return fmt.Errorf("failed to query next pending job: %w", err)
		}

		now := time.Now()

		_, err = tx.Model((*entity.IngestJobEntity)(nil)).
			Context(ctx).
			Set(setStatus, string(models.IngestJobStatusProcessing)).
			Set(setStartedAt, now).
			Where(whereIngestJobID, jobRow.ID).
			Update()
		if err != nil {
			return fmt.Errorf("failed to update job status to processing: %w", err)
		}

		submissionRow := new(entity.SubmissionEntity)
		err = tx.Model(submissionRow).
			Context(ctx).
			Where(whereSubmissionID, jobRow.SubmissionID).
			Select()
		if err != nil {
			return fmt.Errorf("failed to query submission for job: %w", err)
		}

		_, err = tx.Model((*entity.SubmissionEntity)(nil)).
			Context(ctx).
			Set(setStatus, string(models.SubmissionStatusParsing)).
			Where(whereSubmissionID, jobRow.SubmissionID).
			Update()
		if err != nil {
			return fmt.Errorf("failed to update submission status to parsing: %w", err)
		}

		returnIngestJob = &models.IngestJob{
			ID:           jobRow.ID,
			SubmissionID: jobRow.SubmissionID,
			Status:       models.IngestJobStatusProcessing,
			CreatedAt:    jobRow.CreatedAt,
			StartedAt:    &now,
			FinishedAt:   jobRow.FinishedAt,
			ErrorMessage: jobRow.ErrorMessage,
		}

		returnSubmission = &models.Submission{
			ID:             submissionRow.ID,
			LabID:          submissionRow.LabID,
			StudentID:      submissionRow.StudentID,
			Status:         models.SubmissionStatusParsing,
			SourceFileName: submissionRow.SourceFileName,
			MimeType:       submissionRow.MimeType,
			StorageKey:     submissionRow.StorageKey,
			SubmittedAt:    submissionRow.SubmittedAt,
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return returnIngestJob, returnSubmission, nil
}

func (r *PostgresRepository) MarkJobDone(ctx context.Context, jobID string, parsed models.ParsedSubmission) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO parsed_submissions (
				parsed_submission_id,
				submission_id,
				raw_text,
				structured_data,
				parsed_at
			) VALUES (?, ?, ?, ?::jsonb, ?)
		`, parsed.ID, parsed.SubmissionID, parsed.RawText, string(parsed.StructuredData), parsed.ParsedAt)
		if err != nil {
			return fmt.Errorf("failed to insert parsed submission: %w", err)
		}

		_, err = tx.Model((*entity.SubmissionEntity)(nil)).
			Context(ctx).
			Set(setStatus, string(models.SubmissionStatusParsed)).
			Where(whereSubmissionID, parsed.SubmissionID).
			Update()
		if err != nil {
			return fmt.Errorf("failed to update submission status to parsed: %w", err)
		}

		_, err = tx.Model((*entity.IngestJobEntity)(nil)).
			Context(ctx).
			Set(setStatus, string(models.IngestJobStatusDone)).
			Set(setFinishedAt, time.Now()).
			Set(setErrorMessageNil).
			Where(whereIngestJobID, jobID).
			Update()
		if err != nil {
			return fmt.Errorf("failed to update ingest job status to done: %w", err)
		}

		return nil
	})
}

func (r *PostgresRepository) MarkJobFailed(ctx context.Context, jobID, submissionID, errorMessage string) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		_, err := tx.Model((*entity.SubmissionEntity)(nil)).
			Context(ctx).
			Set(setStatus, string(models.SubmissionStatusFailed)).
			Where(whereSubmissionID, submissionID).
			Update()
		if err != nil {
			return fmt.Errorf("failed to update submission status to failed: %w", err)
		}

		_, err = tx.Model((*entity.IngestJobEntity)(nil)).
			Context(ctx).
			Set(setStatus, string(models.IngestJobStatusFailed)).
			Set(setFinishedAt, time.Now()).
			Set(setErrorMessage, errorMessage).
			Where(whereIngestJobID, jobID).
			Update()
		if err != nil {
			return fmt.Errorf("failed to update ingest job status to failed: %w", err)
		}

		return nil
	})
}
