package models

import (
	"io"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
)

type DTOMessageCreate struct {
	Date    time.Time
	Subject string
	From    string
	To      []string
	Text    string
	HTML    string
}

type DTOAttachmentCreate struct {
	Name string
	Data io.Reader
}

type DTOEnvelopeListRequest struct {
	Search        string
	SearchSubject bool
	SearchText    bool
	Ascending     bool
	Order         dtoEnvelopeField
}

type dtoEnvelopeField string

func NewDTOEnvelopeField(s string) dtoEnvelopeField {
	switch s {
	case DTOEnvelopeFieldFrom:
		return DTOEnvelopeFieldFrom
	case DTOEnvelopeFieldSubject:
		return DTOEnvelopeFieldSubject
	default:
		return DTOEnvelopeFieldCreatedAt
	}
}

const (
	DTOEnvelopeFieldCreatedAt = "created_at"
	DTOEnvelopeFieldFrom      = "from"
	DTOEnvelopeFieldSubject   = "subject"
)

type DTOEnvelopeListResult struct {
	PageResult pagination.PageResult
	Envelopes  []Envelope
}

type DTOAttachmentListRequest struct {
	Ascending bool
}

type DTOAttachmentListResult struct {
	PageResult  pagination.PageResult
	Attachments []Attachment
}

type DTOTraceListRequest struct {
	Ascending bool
}

type DTOTraceListResult struct {
	PageResult pagination.PageResult
	Traces     [][]Trace
}

type DTORuleCreate struct {
	Name       string
	Expression string
	Endpoints  []int64
}

type DTORuleUpdate struct {
	ID         int64
	Name       *string
	Expression *string
	Enable     *bool
	Endpoints  *[]int64
}
