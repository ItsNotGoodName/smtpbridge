package attachment

import "context"

type AttachmentService struct {
	attachmentRepository Repository
}

func NewAttachmentService(attachmentRepository Repository) *AttachmentService {
	return &AttachmentService{
		attachmentRepository: attachmentRepository,
	}
}

func (as *AttachmentService) Create(ctx context.Context, param *Param) (*Attachment, error) {
	att, err := New(param)
	if err != nil {
		return nil, err
	}

	return att, as.attachmentRepository.Create(ctx, att)
}
