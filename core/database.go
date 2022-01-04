package core

import "errors"

var (
	ErrDatabaseCleanup = errors.New("database cleanup")
)

type Database interface {
	Close() error
	AttachmentRepository() AttachmentRepositoryPort
	MessageRepository() MessageRepositoryPort
}
