package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
)

type Telegram struct {
	Token  string
	ChatID string
	Client *http.Client
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{
		Token:  token,
		ChatID: chatID,
		Client: &http.Client{},
	}
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (t *Telegram) Send(ctx context.Context, env endpoint.Envelope) error {
	// Send with 0 attachments
	if len(env.Attachments) == 0 {
		return t.sendMessage(ctx, env.Message.Text())
	}

	// TODO: use sendMediaGroup when more than 1 attachment

	// Send with 1 attachment
	if err := t.sendPhoto(ctx, env.Message.Text(), env.Attachments[0].Name, env.Attachments[0].Data); err != nil {
		return err
	}

	// Send rest of attachments
	length := len(env.Attachments)
	for i := 1; i < length; i++ {
		if err := t.sendPhoto(ctx, "", env.Attachments[i].Name, env.Attachments[i].Data); err != nil {
			return err
		}
	}

	return nil
}

func (t *Telegram) sendMessage(ctx context.Context, text string) error {
	// Create and send request
	if len(text) > 4096 {
		text = text[:4096]
	}

	// Create request
	values := url.Values{"chat_id": {t.ChatID}, "text": {text}}
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+t.Token+"/sendMessage", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := t.Client.Do(req)
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

func (t *Telegram) sendPhoto(ctx context.Context, caption, name string, photo []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Photo
	w, err := writer.CreateFormFile("photo", name)
	if err != nil {
		return err
	}
	w.Write(photo)

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
	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+t.Token+"/sendPhoto?chat_id="+t.ChatID, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := t.Client.Do(req)
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
