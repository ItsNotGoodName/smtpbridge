//go:build dev

package api

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

const AttachmentURI = "/attachment/"

func convertAttachments(r *http.Request, attachments []dto.Attachment) {
	for i := range attachments {
		attachments[i].File = "http://" + r.Host + AttachmentURI + attachments[i].File
	}
}
