package senders

import (
	"context"
	"io"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

type Transformer interface {
	Title(ctx context.Context, env models.Envelope) (string, error)
	Body(ctx context.Context, env models.Envelope) (string, error)
	Files(ctx context.Context, env models.Envelope) ([]models.Attachment, error)
	Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error)
	Path(ctx context.Context, att models.Attachment) (string, error)
}
