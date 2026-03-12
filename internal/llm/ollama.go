package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/narvdeshwar/IOT-rag/internal/types"
)

var _ types.LLM = (*OllamaLLM)(nil)

type OllamaLLM struct {
	url   string
	model string
}

func NewOllamaLLM(url, model string) *OllamaLLM {
	return &OllamaLLM{
		url:   url,
		model: model,
	}
}

type ollamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []ollamaChatMessage `json:"messages"`
}

type ollamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatResponse struct {
	Message ollamaChatMessage `json:"message"`
	Done    bool              `json:"done"`
}

func (l *OllamaLLM) StreamComplete(ctx context.Context, prompt string, onToken func(string)) (string, error) {
	reqBody, _ := json.Marshal(ollamaChatRequest{
		Model: l.model,
		Messages: []ollamaChatMessage{
			{Role: "system", Content: SystemPrompt},
			{Role: "user", Content: prompt},
		},
	})

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/chat", l.url), bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama chat failed: %s", resp.Status)
	}

	var fullText strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk ollamaChatResponse
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			return "", err
		}
		token := chunk.Message.Content
		fullText.WriteString(token)
		onToken(token)
		if chunk.Done {
			break
		}
	}

	return fullText.String(), nil
}
