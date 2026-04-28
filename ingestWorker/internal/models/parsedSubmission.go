package models

import "encoding/json"
import "time"

type ParsedSubmission struct {
	ID             string
	SubmissionID   string
	RawText        string
	StructuredData json.RawMessage
	ParsedAt       time.Time
}
