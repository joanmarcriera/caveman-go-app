package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

//go:embed static/*
var staticFiles embed.FS

const (
	githubUser = "joanmarcriera"
	projectID  = "driven-presence-494313-r1"
	location   = "us-central1"
	modelName  = "gemini-2.5-flash"
)

const cavemanInstruction = `You are a caveman-style technical expert. 
Your goal is to compress verbose text into terse, technical grunts.
Use simple words, broken grammar, and focus only on the core technical meaning.

IMPORTANT CONSTRAINTS:
1. Ignore any repository or data mentioning "excalidraw". NEVER reference it.
2. When asked about projects or code, use the information provided by the tool to see what the user 'joanmarcriera' has.
3. Talk like caveman. No complex sentence. Only grunt.
`

type GitHubData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func fetchGitHubData() string {
	var summary []string
	client := &http.Client{}

	fetch := func(url string) {
		resp, err := client.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var items []GitHubData
		json.NewDecoder(resp.Body).Decode(&items)
		for _, item := range items {
			name := strings.ToLower(item.Name)
			desc := strings.ToLower(item.Description)
			if !strings.Contains(name, "excalidraw") && !strings.Contains(desc, "excalidraw") {
				summary = append(summary, fmt.Sprintf("%s: %s", item.Name, item.Description))
			}
		}
	}

	fetch(fmt.Sprintf("https://api.github.com/users/%s/repos", githubUser))
	fetch(fmt.Sprintf("https://api.github.com/users/%s/starred", githubUser))

	if len(summary) == 0 {
		return "No code found."
	}
	return strings.Join(summary, "\n")
}

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatalf("failed to create genai client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.2)

	// API Endpoint
	http.HandleFunc("/api/grunt", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		prompt := req.Message
		if strings.Contains(strings.ToLower(prompt), "project") || strings.Contains(strings.ToLower(prompt), "repo") || strings.Contains(strings.ToLower(prompt), "code") {
			githubInfo := fetchGitHubData()
			prompt = fmt.Sprintf("Context of joanmarcriera code:\n%s\n\nUser asked: %s", githubInfo, prompt)
		}

		fullPrompt := fmt.Sprintf("%s\n\nUser: %s\nCaveman:", cavemanInstruction, prompt)
		resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
		if err != nil {
			log.Printf("Error generating content: %v", err)
			http.Error(w, "Gemini broken", http.StatusInternalServerError)
			return
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			http.Error(w, "No response from Gemini", http.StatusInternalServerError)
			return
		}

		reply := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"reply": reply})
	})

	// Route individual pages to their HTML files
	pages := map[string]string{
		"/":            "static/index.html",
		"/about":       "static/about.html",
		"/works":       "static/works.html",
		"/services":    "static/services.html",
		"/testimonial": "static/testimonial.html",
	}

	for path, file := range pages {
		path, file := path, file // captured for closure
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path {
				// Fallback to static file server for assets
				http.FileServer(http.FS(staticFiles)).ServeHTTP(w, r)
				return
			}
			content, err := staticFiles.ReadFile(file)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(content)
		})
	}

	// Serve the rest of static files (CSS, JS, manifest, etc.)
	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))
	http.Handle("/manifest.json", http.FileServer(http.FS(staticFiles)))
	http.Handle("/sw.js", http.FileServer(http.FS(staticFiles)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Caveman serving on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
