# 📝 Telegram Todo Bot (Go)

A simple multi-user **Telegram Todo Bot** written in **Go**, built to learn:
- Go fundamentals
- JSON-based persistence
- Telegram Bot API
- State handling in chatbots
- Basic backend design principles

Each Telegram chat has its **own isolated task list**.

---

## 🚀 Features

- ➕ Add tasks
- 📋 List tasks
- ❌ Delete tasks by number
- 👥 Multi-user support (tasks are stored per chat ID)
- 💾 Persistent storage using JSON
- 🔄 Long-polling with Telegram `getUpdates`
- ☁️ Deployable to Railway / cloud environments

---

## 🧠 What I Learned From This Project

This project was built as a **learning exercise**, focusing on:

- Go slices vs maps (`map[int64][]Task`)
- Passing data explicitly between layers
- Handling per-user state in chatbots
- JSON encoding / decoding and data modeling
- Error handling and debugging runtime issues
- Structuring a real, runnable backend service

---

## 🛠 Tech Stack

- **Language:** Go
- **API:** Telegram Bot API
- **Storage:** JSON file
- **Deployment:** Railway (optional)
- **Version Control:** Git + GitHub

---

## 🧩 Project Structure

```text
.
├── main.go          # Bot logic, state machine, polling loop
├── telegram.go      # Telegram API interaction
├── go.mod
├── README.md
└── thingstodo.json  # Created at runtime
