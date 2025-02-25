package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Client struct {
	ApiBaseUrl string
	BotToken   string
	Logger     *slog.Logger
}

func NewClient(apiBaseUrl string, botToken string) *Client {
	return &Client{
		ApiBaseUrl: apiBaseUrl,
		BotToken:   botToken,
		Logger:     slog.Default(),
	}
}

func (c *Client) SendMessage(chatID int64, message string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", c.ApiBaseUrl, c.BotToken)

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger.Error("Failed to close response body", slog.Any("error", err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code received: %d", resp.StatusCode)
	}

	c.Logger.Info("Message sent successfully", slog.Int64("chat_id", chatID))
	return nil
}
