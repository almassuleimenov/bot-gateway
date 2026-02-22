# 🚀 Architecture AI Bot: Gateway (Go)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Architecture](https://img.shields.io/badge/Architecture-Microservices-orange?style=flat)](#)
[![Status](https://img.shields.io/badge/Status-MVP_Ready-green?style=flat)](#)

> **The "Body" of our AI ecosystem.** > Этот микросервис на Go выступает в роли высокопроизводительного шлюза (Gateway), который связывает Telegram API с "мозгом" на Python.

---

## 🏗 Как это работает?

Gateway спроектирован по принципу **Non-blocking Reactive Flow**:
1. Принимает **Webhooks** от Telegram.
2. Моментально отвечает серверу Telegram `200 OK` (чтобы избежать повторных запросов).
3. В фоновом режиме пробрасывает запрос в **AI-Service (Python)** через REST API.
4. Доставляет интеллектуальный ответ клиенту.

## 🔥 Фишки
- **Lightweight:** Минимальное потребление памяти благодаря компилируемому Go.
- **Async Processing:** Обработка сообщений без задержек.
- **Safety first:** Строгая типизация и обработка ошибок Telegram API.
- **Clean Architecture:** Разделение на `models`, `services` и `handlers`.

---

## 🛠 Технологический стек
* **Language:** Go (Golang)
* **Router:** `go-chi/chi` (легкий и быстрый)
* **API:** Telegram Bot API
* **Communication:** HTTP/JSON (REST)

---
