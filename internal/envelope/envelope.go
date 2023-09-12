package envelope

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/jaytaylor/html2text"
	"github.com/samber/lo"
)

func New(msg models.Message, atts ...models.Attachment) models.Envelope {
	return models.Envelope{
		Message:     msg,
		Attachments: atts,
	}
}

func NewMessage(r models.DTOMessageCreate) models.Message {
	text := r.Text
	if isHTML(r.Text) {
		var err error
		text, err = html2text.FromString(r.Text, html2text.Options{
			TextOnly: true,
		})
		if err != nil {
			text = r.Text
		}
	}

	to := lo.Filter(lo.Uniq(r.To), func(to string, _ int) bool {
		return strings.Trim(to, " ") != ""
	})

	return models.Message{
		From:      r.From,
		To:        to,
		CreatedAt: models.NewTime(time.Now()),
		Subject:   r.Subject,
		Text:      text,
		HTML:      r.HTML,
		Date:      models.NewTime(r.Date),
	}
}

func NewDataAttachment(name string, data io.Reader) (models.DataAttachment, error) {
	rd := bufio.NewReaderSize(data, 512)
	b, err := rd.Peek(512)
	if err != nil {
		return models.DataAttachment{}, err
	}

	mimeT := http.DetectContentType(b)
	extension := fileExtension(name, mimeT)

	return models.DataAttachment{
		Attachment: models.Attachment{
			Name:      name,
			Mime:      mimeT,
			Extension: extension,
		},
		Reader: rd,
	}, nil
}
