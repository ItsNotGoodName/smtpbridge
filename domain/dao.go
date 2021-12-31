package domain

type DAO struct {
	Attachment AttachmentRepositoryPort
	Message    MessageRepositoryPort
	Endpoint   EndpointRepositoryPort
}
