package app

type DAO struct {
	Attachment AttachmentRepositoryPort
	Message    MessageRepositoryPort
	Endpoint   EndpointRepositoryPort
}

func NewDAO(Attachment AttachmentRepositoryPort, Message MessageRepositoryPort, Endpoint EndpointRepositoryPort) DAO {
	return DAO{
		Attachment,
		Message,
		Endpoint,
	}
}
