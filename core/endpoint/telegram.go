package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
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

func (t *Telegram) Send(ctx context.Context, text string, atts []Attachment) error {
	// Send with 0 attachments
	if len(atts) == 0 {
		return t.sendMessage(ctx, text)
	}

	// TODO: use sendMediaGroup when more than 1 attachment

	// Send with 1 attachment
	if err := t.sendPhoto(ctx, text, atts[0].Name, atts[0].Data); err != nil {
		return err
	}

	// Send rest of attachments
	length := len(atts)
	for i := 1; i < length; i++ {
		if err := t.sendPhoto(ctx, "", atts[i].Name, atts[i].Data); err != nil {
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
