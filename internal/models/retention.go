package models

import "time"

type RetentionPolicy struct {
	EnvelopeCount  int
	EnvelopeAge    time.Duration
	AttachmentSize int64
}
