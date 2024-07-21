package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Webhook struct {
	url  string
	name string
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func New(url string, name string) *Webhook {
	return &Webhook{url: url, name: name}
}

func (w *Webhook) Send(content string) error {
	msg := &Message{
		Username: w.name,
		Content:  content,
	}

	return w.SendMsg(msg)
}

func (w *Webhook) SendMsg(msg *Message) error {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(w.url, "application/json", payload)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errMsg []byte
		_, _ = resp.Body.Read(errMsg)

		return fmt.Errorf("error response: %d: %s", resp.StatusCode, string(errMsg))
	}

	return nil
}
