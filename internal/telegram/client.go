package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const sendMessageUrl = "https://api.telegram.org/bot%s/sendMessage"

type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Client struct {
	BotToken string
	Logger   *slog.Logger
}

func NewClient(botToken string) *Client {
	return &Client{
		BotToken: botToken,
		Logger:   slog.Default(),
	}
}

func (c *Client) SendMessage(chatID int64, message string) error {
	url := fmt.Sprintf(sendMessageUrl, c.BotToken)

	msg := Message{
		ChatID: chatID,
		Text:   message,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // We send messages infrequently, so disable keep-alives
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Logger.Error("Failed to send message", slog.Int("status_code", resp.StatusCode))
		return nil
	}

	c.Logger.Info("Message sent successfully", slog.Int64("chat_id", chatID))
	return nil
}
