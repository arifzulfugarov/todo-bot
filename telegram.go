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
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	Message *Message `json:"message"`
	Data    string   `json: "data"`
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

func sendDeletionMenu(chatID int64) error {

	fileData, _ := os.ReadFile(filename)
	var thingstodo map[int64][]Task
	token := os.Getenv("TELEGRAM_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	json.Unmarshal(fileData, &thingstodo)

	tasks := thingstodo[chatID]

	if len(tasks) == 0 {
		return sendText(chatID, "Your list is empty! 🎉")
	}

	var rows [][]map[string]any
	for i, task := range tasks {
		displayText := fmt.Sprintf("%d. %s", i+1, task.WhatToDo)

		if len(displayText) > 35 {
			displayText = displayText[:32] + "..."
		}

		button := map[string]any{
			"text":          displayText,
			"callback_data": fmt.Sprintf("del_%d", i+1),
		}
		rows = append(rows, []map[string]any{button})
	}

	payload := map[string]any{
		"chat_id": chatID,
		"text":    "🗑️ Select a task to delete:",
		"reply_markup": map[string]any{
			"inline_keyboard": rows,
		},
	}

	return postJSON(url, payload)
}

func answerCallback(callbackQueryID string) error {
	token := os.Getenv("TELEGRAM_TOKEN")
	url := "https://api.telegram.org/bot" + token + "/answerCallbackQuery"

	payload := map[string]any{
		"callback_query_id": callbackQueryID, // The unique ID from the update
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
