package endpoints

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/containrrr/shoutrrr"
)

func GetTitle(env envelope.Envelope, config Config) string {
	if config.TextDisable {
		return ""
	}

	return env.Message.Subject
}

func GetBody(env envelope.Envelope, config Config) (string, error) {
	if config.TextDisable {
		return "", nil
	}

	var buffer bytes.Buffer
	if err := config.BodyTemplate.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func GetAttachments(env envelope.Envelope, config Config) []*envelope.Attachment {
	if config.AttachmentsDisable {
		return []*envelope.Attachment{}
	}

	return env.Attachments
}

func FilterImagesOnly(atts []*envelope.Attachment) []*envelope.Attachment {
	imgAtts := []*envelope.Attachment{}
	for _, a := range atts {
		if a.IsImage() {
			imgAtts = append(imgAtts, a)
		}
	}

	return imgAtts
}

func SenderFactory(kind string, config map[string]string) (Sender, error) {
	switch kind {
	case "console":
		return NewConsole(), nil
	case "telegram":
		return NewTelegram(config["token"], config["chat_id"]), nil
	case "shoutrrr":
		router, err := shoutrrr.CreateSender(strings.Split(config["urls"], "\n")...)
		if err != nil {
			return nil, err
		}
		return NewShoutrrr(router), nil
	default:
		return nil, fmt.Errorf("invalid sender type: %s", kind)
	}
}
