package senders

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

type Telegram struct {
	token  string
	chatID string
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{
		token:  token,
		chatID: chatID,
	}
}

func (t Telegram) Send(ctx context.Context, env models.Envelope, tr Transformer) error {
	var images []models.Attachment
	{
		files, err := tr.Files(ctx, env)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsImage() {
				images = append(images, file)
			}
		}
	}

	body, err := tr.Body(ctx, env)
	if err != nil {
		return err
	}
	// Telegrams's body cannot be empty
	if body == "" {
		body, err = tr.Title(ctx, env)
		if err != nil {
			return err
		}
	}

	// Send with 0 attachments
	if len(images) == 0 {
		return t.sendMessage(ctx, body)
	}

	// TODO: use sendMediaGroup when more than 1 attachment

	// Send with 1 attachment
	rd, err := tr.Reader(ctx, images[0])
	if err != nil {
		return err
	}
	if err := t.sendPhoto(ctx, body, images[0].Name, rd); err != nil {
		return err
	}

	// Send rest of attachments
	length := len(images)
	if length > 10 {
		length = 10
	}
	for i := 1; i < length; i++ {
		rd, err := tr.Reader(ctx, images[i])
		if err != nil {
			return err
		}
		if err := t.sendPhoto(ctx, "", images[i].Name, rd); err != nil {
			return err
		}
	}

	return nil
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (t *Telegram) sendMessage(ctx context.Context, text string) error {
	if text == "" {
		return nil
	}

	// Create and send request
	if len(text) > 4096 {
		text = text[:4096]
	}

	// Create request
	values := url.Values{"chat_id": {t.chatID}, "text": {text}}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.telegram.org/bot"+t.token+"/sendMessage", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	res := &telegramResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Description)
	}

	return nil
}

func (t *Telegram) sendPhoto(ctx context.Context, caption, name string, file io.ReadCloser) error {
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Photo
	w, err := writer.CreateFormFile("photo", name)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	// Caption
	if caption != "" {
		w, err = writer.CreateFormField("caption")
		if err != nil {
			return err
		}
		if len(caption) > 1024 {
			caption = caption[:1024]
		}
		w.Write([]byte(caption))
	}

	// Close
	if err := writer.Close(); err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.telegram.org/bot"+t.token+"/sendPhoto?chat_id="+t.chatID, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	res := &telegramResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Description)
	}

	return nil
}
