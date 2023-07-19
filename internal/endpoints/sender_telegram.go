package endpoints

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
)

type Telegram struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (t *Telegram) Send(cc *core.Context, env envelope.Envelope, config Config) error {
	atts := FilterImagesOnly(GetAttachments(env, config))
	text, err := GetBody(env, config)
	if err != nil {
		return err
	}

	// Send with 0 attachments
	if len(atts) == 0 {
		return t.sendMessage(cc, text)
	}

	// TODO: use sendMediaGroup when more than 1 attachment

	// Send with 1 attachment
	file, err := files.GetFile(cc, atts[0])
	if err != nil {
		return err
	}
	if err := t.sendPhoto(cc, text, atts[0].Name, file); err != nil {
		return err
	}

	// Send rest of attachments
	length := len(atts)
	for i := 1; i < length; i++ {
		file, err := files.GetFile(cc, atts[i])
		if err != nil {
			return err
		}
		if err := t.sendPhoto(cc, "", atts[i].Name, file); err != nil {
			return err
		}
	}

	return nil
}

func (t *Telegram) sendMessage(cc *core.Context, text string) error {
	if text == "" {
		return nil
	}

	// Create and send request
	if len(text) > 4096 {
		text = text[:4096]
	}

	// Create request
	values := url.Values{"chat_id": {t.chatID}, "text": {text}}
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+t.token+"/sendMessage", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req = req.WithContext(cc.Context())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := t.client.Do(req)
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

func (t *Telegram) sendPhoto(cc *core.Context, caption, name string, file *os.File) error {
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
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+t.token+"/sendPhoto?chat_id="+t.chatID, body)
	if err != nil {
		return err
	}
	req = req.WithContext(cc.Context())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := t.client.Do(req)
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
