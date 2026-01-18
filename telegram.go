package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

// ---- UI ----

func sendMenu(chatID int64, text string) error {
	token := os.Getenv("TELEGRAM_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
		"reply_markup": map[string]any{
			"keyboard": [][]map[string]string{
				{
					{"text": "Add"}, {"text": "List"},
					{"text": "Delete"}, {"text": "ℹ️Help"},
				},
			},
			"resize_keyboard":   true,
			"one_time_keyboard": false,
		},
	}

	return postJSON(url, payload)
}

func sendText(chatID int64, text string) error {
	token := os.Getenv("TELEGRAM_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
		"link_preview_options": map[string]any{
			"is_disabled": true,
		},
	}

	return postJSON(url, payload)
}

func sendTextWithRemove(chatID int64, text string) error {
	token := os.Getenv("TELEGRAM_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
		"reply_markup": map[string]any{
			"remove_keyboard": true,
			"selective":       true,
		},
		"link_preview_options": map[string]any{
			"is_disabled": true,
		},
	}

	return postJSON(url, payload)
}

func postJSON(url string, data any) error {
	b, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ---- Telegram API ----

func getUpdates(offset int) ([]Update, error) {
	token := os.Getenv("TELEGRAM_TOKEN")
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%d",
		token,
		offset,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var updates UpdateResponse
	if err := json.Unmarshal(body, &updates); err != nil {
		return nil, err
	}

	return updates.Result, nil
}
