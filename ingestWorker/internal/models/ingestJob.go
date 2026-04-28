package models

import "time"

type IngestJobStatus string

const (
	IngestJobStatusPending    IngestJobStatus = "pending"
	IngestJobStatusProcessing IngestJobStatus = "processing"
	IngestJobStatusDone       IngestJobStatus = "done"
	IngestJobStatusFailed     IngestJobStatus = "failed"
)

type IngestJob struct {
	ID           string
	SubmissionID string
	Status       IngestJobStatus
	CreatedAt    time.Time
	StartedAt    *time.Time
	FinishedAt   *time.Time
	ErrorMessage *string
}
