package endpoint

import (
	"encoding/json"
	"errors"
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

//func (t *Telegram) Name() string {
//	return "telegram"
//}

//func (t *Telegram) Capabilities() []app.Capability {
//	return []app.Capability{app.CapabilityText, app.CapabilityImage}
//}

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

	return nil
}
