package models

import (
	"time"
)

type RetentionPolicy struct {
	EnvelopeCount  *int
	EnvelopeAge    *time.Duration
	MinEnvelopeAge time.Duration
	AttachmentSize *int64
}

func (rp RetentionPolicy) EnvelopeAgeDate() time.Time {
	date := time.Now()
	if rp.EnvelopeAge != nil && *rp.EnvelopeAge > rp.MinEnvelopeAge {
		date = date.Add(-*rp.EnvelopeAge)
	} else {
		date = date.Add(-rp.MinEnvelopeAge)
	}

	return date
}
