package core

type Database interface {
	Close() error
	AttachmentRepository() AttachmentRepositoryPort
	MessageRepository() MessageRepositoryPort
}
