package parser

import "context"

type ResultParser struct {
	RawText        string
	StructuredData []byte
}

type DocumentParser interface {
	SupportMimeType(mimeType string) bool
	ParseDocument(ctx context.Context, filePath string) (*ResultParser, error)
}
