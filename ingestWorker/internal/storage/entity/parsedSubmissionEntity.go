package entity

import "encoding/json"
import "time"

type ParsedSubmissionEntity struct {
	tableName struct{} `pg:"parsed_submissions"`

	ID             string          `pg:"parsed_submission_id,pk"`
	SubmissionID   string          `pg:"submission_id,notnull"`
	RawText        string          `pg:"raw_text,notnull"`
	StructuredData json.RawMessage `pg:"structured_data,notnull"`
	ParsedAt       time.Time       `pg:"parsed_at,notnull"`
	CreatedAt      time.Time       `pg:"created_at,notnull"`
}
