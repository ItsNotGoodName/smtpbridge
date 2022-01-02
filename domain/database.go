package domain

type Database interface {
	Close() error
	AttachmentRepository() AttachmentRepositoryPort
	MessageRepository() MessageRepositoryPort
}
