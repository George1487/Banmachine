package models

import "time"

type SubmissionStatus string

const (
	SubmissionStatusUploaded SubmissionStatus = "uploaded"
	SubmissionStatusParsing  SubmissionStatus = "parsing"
	SubmissionStatusParsed   SubmissionStatus = "parsed"
	SubmissionStatusFailed   SubmissionStatus = "failed"
)

type Submission struct {
	ID             string
	LabID          string
	StudentID      string
	Status         SubmissionStatus
	SourceFileName string
	MimeType       string
	StorageKey     string
	SubmittedAt    time.Time
}
