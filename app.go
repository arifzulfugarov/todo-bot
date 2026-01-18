package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	WhatToDo     string `json:"whattodo"`
	NumberOfTask int    `json:"numberoftask"`
}

const filename = "thingstodo.json"
const helpMessage = `ℹ️Help

➕ Add — create a new task  
📋 List — view your tasks  
❌ Delete — remove a task  

Just tap the buttons below 👇`

func addTask(chatID int64, description string) error {
	var thingstodo map[int64][]Task

	text := strings.TrimSpace(description)

	if len(text) > 300 {
		return fmt.Errorf("Task is too long")
	}

	if text == "" {
		return fmt.Errorf("Please provide text message only")
	}

	fileData, err := os.ReadFile(filename)
	if err == nil {
		if err := json.Unmarshal(fileData, &thingstodo); err != nil {
			return err
		}
	}

	if thingstodo == nil {
		thingstodo = make(map[int64][]Task)
	}

	tasks := thingstodo[chatID]

	newTask := Task{
		WhatToDo:     description,
		NumberOfTask: len(thingstodo) + 1,
	}

	tasks = append(tasks, newTask)
	thingstodo[chatID] = tasks

	updateData, err := json.MarshalIndent(thingstodo, "", "")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, updateData, 0644)
}

func list(chatID int64) (string, error) {

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return "No tasks found, add some tasks", nil
	}

	var thingstodo map[int64][]Task
	if err := json.Unmarshal(fileData, &thingstodo); err != nil {
		return "", err
	}

	tasks := thingstodo[chatID]
	if len(tasks) == 0 {
		return "No tasks found, add some tasks", nil
	}

	var b strings.Builder
	for i, task := range tasks {
		fmt.Fprintf(&b, "%d. %s\n", i+1, task.WhatToDo)
	}

	return b.String(), nil

}

func deleteTask(chatID int64, number int) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var thingstodo map[int64][]Task
	if err := json.Unmarshal(fileData, &thingstodo); err != nil {
		return err
	}

	tasks := thingstodo[chatID]

	if len(tasks) == 0 {
		return fmt.Errorf("The list is empty, please add something")
	}

	index := number - 1
	if index < 0 || index >= len(tasks) {
		return fmt.Errorf("task %d doesn't exist", number)
	}

	// delete task from this user's slice
	tasks = append(tasks[:index], tasks[index+1:]...)

	// renumber
	for i := range tasks {
		tasks[i].NumberOfTask = i + 1
	}

	thingstodo[chatID] = tasks

	updateData, err := json.MarshalIndent(thingstodo, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, updateData, 0644)

}

func main() {

	var chatState = make(map[int64]string)

	log.Println("Bot started")

	offset := 0

	for {
		updates, err := getUpdates(offset)
		if err != nil {
			log.Println("error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, u := range updates {

			// Inside your main loop
			chatID := u.Message.Chat.ID

			// Save the user's message ID (e.g., "Add" or the task text) to the trash bin
			h := messageHistory[chatID]
			h = append(h, u.Message.MessageID)
			messageHistory[chatID] = h

			if u.CallbackQuery != nil {
				data := u.CallbackQuery.Data
				chatID := u.CallbackQuery.Message.Chat.ID

				if strings.HasPrefix(data, "del_") {
					numStr := strings.TrimPrefix(data, "del_")
					number, _ := strconv.Atoi(numStr)

					err := deleteTask(chatID, number)
					if err != nil {
						sendText(chatID, "Error: "+err.Error())
					} else {
						clearAndRefresh(chatID)
						sendText(chatID, fmt.Sprintf("✅ Task #%d is done, good job!!!", number))
					}
				}

				answerCallback(u.CallbackQuery.ID)
				offset = u.UpdateID + 1
				continue
			}

			if u.Message == nil {
				offset = u.UpdateID + 1
				continue
			}

			text := u.Message.Text
			state := chatState[chatID]

			if u.Message.Text == "" {
				// UC-5: Provide clear feedback for invalid format
				sendText(chatID, "Invalid input. Please provide text only (no stickers, photos, or audio).")
				offset = u.UpdateID + 1
				continue
			}

			switch state {

			case "add":
				err := addTask(chatID, text)
				if err != nil {
					sendText(chatID, "❌ "+err.Error())
				} else {
					clearAndRefresh(chatID)

					currentList, _ := list(chatID)
					sendMenu(chatID, "✅ Task added! Here is your updated list:\n\n"+currentList)
				}

				chatState[chatID] = ""

			case "delete":
				number, err := strconv.Atoi(text)
				if err != nil {
					sendText(chatID, "Please send a valid number that normal people use")
					break
				}

				err = deleteTask(chatID, number)
				if err != nil {
					sendText(chatID, err.Error())
				} else {
					sendText(chatID, "Task deleted, great job!!!")
				}

				chatState[chatID] = ""

			default:
				switch text {
				case "/start":
					welcomeMMessage := "Welcome! \nI’m your personal *Todo Bot*.\nI help you manage tasks directly in Telegram.\n\nWhat can I help you with today?"
					sendMenu(chatID, welcomeMMessage)
				case "Add":
					chatState[chatID] = "add"
					sendTextWithRemove(chatID, " Please type your task below:")
				case "List":
					result, err := list(chatID)
					if err != nil {
						sendText(chatID, "Failed to list tasks")
					} else {
						sendText(chatID, result)
					}
				case "ℹ️Help":
					sendText(chatID, helpMessage)

				case "Delete":
					sendDeletionMenu(chatID)

				default:

					sendText(chatID, "❓ I didn't understand that. Please use the buttons below.")

				}

			}

			offset = u.UpdateID + 1
		}

		time.Sleep(2 * time.Second)
	}
}
