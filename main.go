package tlgbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Send(botToken string, chatID int, text string) (*http.Response, error) {
	t := NewTextMessage(botToken, chatID, text)
	return t.Send()
}

func SendWithContext(ctx context.Context, botToken string, chatID int, text string) (*http.Response, error) {
	t := NewTextMessage(botToken, chatID, text)
	return t.Send(ctx)
}

type TextMessage struct {
	BotToken string `json:"-"`
	ChatID   int    `json:"chat_id"`
	Text     string `json:"text"`
}

func NewTextMessage(botToken string, chatID int, text string) TextMessage {
	return TextMessage{
		botToken,
		chatID,
		text,
	}
}
func (t *TextMessage) Send(ctx ...context.Context) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(t); err != nil {
		return nil, err
	}

	var (
		c   *http.Client
		req *http.Request
		err error
	)

	URL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)
	if len(ctx) > 0 {
		req, err = http.NewRequestWithContext(ctx[0], "POST", URL, &buf)
	} else {
		c.Timeout = 10 * time.Second
		req, err = http.NewRequest("POST", URL, &buf)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}
