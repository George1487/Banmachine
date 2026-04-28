package entity

import "time"

type IngestJobEntity struct {
	tableName struct{} `pg:"ingest_jobs"`

	ID           string     `pg:"ingest_job_id,pk"`
	SubmissionID string     `pg:"submission_id,notnull"`
	Status       string     `pg:"status,notnull"`
	CreatedAt    time.Time  `pg:"created_at,notnull"`
	StartedAt    *time.Time `pg:"started_at"`
	FinishedAt   *time.Time `pg:"finished_at"`
	ErrorMessage *string    `pg:"error_message"`
}
