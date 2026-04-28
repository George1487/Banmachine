package entity

import "time"

type SubmissionEntity struct {
	tableName struct{} `pg:"submissions"`

	ID             string    `pg:"submission_id,pk"`
	LabID          string    `pg:"lab_id,notnull"`
	StudentID      string    `pg:"student_id,notnull"`
	Status         string    `pg:"status,notnull"`
	SourceFileName string    `pg:"source_file_name,notnull"`
	MimeType       string    `pg:"mime_type,notnull"`
	StorageKey     string    `pg:"storage_key,notnull"`
	SubmittedAt    time.Time `pg:"submitted_at,notnull"`
	CreatedAt      time.Time `pg:"created_at,notnull"`
}
