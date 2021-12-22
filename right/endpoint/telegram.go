package endpoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Telegram struct {
	Token  string
	ChatID string
}

func NewTelegram(token string, chatID string) *Telegram {
	return &Telegram{
		Token:  token,
		ChatID: chatID,
	}
}

type TelegramResponse struct {
	OK bool `json:"ok"`
}

func (t *Telegram) Send(message *app.Message) error {
	response, err := http.PostForm(
		"https://api.telegram.org/bot"+t.Token+"/sendMessage",
		url.Values{
			"chat_id": {t.ChatID}, "text": {message.Text},
		})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	result := &TelegramResponse{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if !result.OK {
		return errors.New("Telegram response is not OK")
	}

	for _, attachment := range message.Attachments {
		t.SendPicture(attachment.Name, attachment.Data)
	}

	return nil
}

func (t *Telegram) SendPicture(name string, pic []byte) error {
	// just need bytes to send picture
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	w, err := writer.CreateFormFile("photo", name)
	if err != nil {
		return err
	}
	w.Write(pic)
	w, err = writer.CreateFormField("caption")
	if err != nil {
		return err
	}
	w.Write([]byte(name))
	writer.Close()

	// make a http post request
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+t.Token+"/sendPhoto?chat_id="+t.ChatID, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
