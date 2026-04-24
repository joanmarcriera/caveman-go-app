# 🦴 Caveman Tech: Monolithic Agent App

A single-binary Go application that transforms verbose technical jargon into terse, technical "caveman grunts." It features a themed web UI, a Go-based agent backend using Gemini 2.5 Flash, and automatic GitHub integration for a specific user.

## 📸 Application Preview (Visual Mockup)

```text
+---------------------------------------+
|             CAVEMAN TECH              |
+---------------------------------------+
|                                       |
|  [GRUNT: Me Caveman. Talk tech. Me    |
|   grunt.]                             |
|                                       |
|          [I implemented a scalable    |
|           K8s cluster today]          |
|                                       |
|  [GRUNT: K8s big. Scalable. Good.]    |
|                                       |
|          [What repos does joan have?] |
|                                       |
|  [GRUNT: Joan have many code.         |
|   ai-factory, android-app, monkey...] |
|                                       |
+---------------------------------------+
| [ Type technical thing... ]  [GRUNT]  |
+---------------------------------------+
```

## 🚀 Features

- **Monolithic Build:** Go backend and HTML/JS frontend embedded into a single executable binary.
- **Caveman Agent:** Powered by **Gemini 2.5 Flash** with a strictly enforced "Caveman Grunt" persona.
- **GitHub Awareness:** Automatically fetches repositories and stars for user `joanmarcriera` using the GitHub API.
- **Strict Filtering:** Automatically ignores and refuses to mention any data related to `excalidraw`.
- **PWA Support:** Installable on iPhone/Android as a standalone app via "Add to Home Screen."

## 🧠 How It Works

### The Caveman Prompt
The heart of the agent is a specific system instruction that enforces the persona:
> "You are a caveman-style technical expert. Your goal is to compress verbose text into terse, technical grunts. Use simple words, broken grammar, and focus only on the core technical meaning."

### The Data Tooling
When the Go backend detects keywords like "repo" or "project," it triggers a manual "tool" function:
1. It queries `https://api.github.com/users/joanmarcriera/repos`.
2. It filters out any result where the name or description contains "excalidraw."
3. It injects this filtered technical context into the prompt before sending it to Gemini.

## 🛠️ Build & Run

### Prerequisites
- Go 1.22+
- Google Cloud Project with Vertex AI enabled.
- Authenticated `gcloud` environment.

### Installation
```bash
cd caveman-go-app
go build -o caveman-app
```

### Execution
```bash
./caveman-app
```
The server will start on `http://localhost:8080`.

## 📱 Mobile Installation (iPhone)
1. Ensure your iPhone is on the same Wi-Fi as your computer.
2. Access the app via your computer's local IP (e.g., `http://192.168.1.50:8080`).
3. Tap **Share** (the box with an arrow).
4. Select **"Add to Home Screen"**.
5. The app will now appear on your home screen with a standalone, borderless UI.

## 🛠️ Tech Stack
- **Backend:** Go (Golang)
- **AI Model:** Google Gemini 2.5 Flash (via Vertex AI Go SDK)
- **Frontend:** Vanilla HTML5, CSS3 (Stone-themed), and JavaScript
- **Deployment:** Single monolithic binary with `//go:embed`
