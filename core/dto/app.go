package dto

import (
	"context"
	"io/fs"
)

type App interface {
	AttachmentFS() fs.FS
	EventList(ctx context.Context, req *EventListRequest) (*EventListResponse, error)
	Info(ctx context.Context) (*InfoResponse, error)
	MessageDelete(ctx context.Context, req *MessageDeleteRequest) error
	MessageEventList(ctx context.Context, req *EventListRequest) (*EventListResponse, error)
	MessageGet(ctx context.Context, req *MessageGetRequest) (*MessageGetResponse, error)
	MessageHandle(ctx context.Context, req *MessageHandleRequest) error
	MessageList(ctx context.Context, req *MessageListRequest) (*MessageListResponse, error)
	SMTPLogin(ctx context.Context, req *SMTPLoginRequest) error
	Version() *VersionResponse
}
