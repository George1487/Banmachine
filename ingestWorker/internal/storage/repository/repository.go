package repository

import (
	"context"
	"errors"
	"ingestWorker/internal/models"
)

var ErrNoPendingJobs = errors.New("no pending ingest jobs")

type Repository interface {
	ClaimNextPendingJob(ctx context.Context) (*models.IngestJob, *models.Submission, error)
	MarkJobDone(ctx context.Context, jobID string, parsed models.ParsedSubmission) error
	MarkJobFailed(ctx context.Context, jobID, submissionID, errorMessage string) error
}
