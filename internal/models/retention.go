package models

import (
	"time"
)

type RetentionPolicy struct {
	MinAge         time.Duration
	EnvelopeCount  *int
	EnvelopeAge    *time.Duration
	AttachmentSize *int64
}

func (rp RetentionPolicy) AgeDate() time.Time {
	date := time.Now()
	if rp.EnvelopeAge != nil && *rp.EnvelopeAge > rp.MinAge {
		date = date.Add(-*rp.EnvelopeAge)
	} else {
		date = date.Add(-rp.MinAge)
	}

	return date
}
