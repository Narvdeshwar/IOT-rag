package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var _ types.LLM = (*GeminiLLM)(nil)

type GeminiLLM struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiLLM(ctx context.Context, apiKey string) (*GeminiLLM, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.L.Error("failed to create gemini client", "err", err)
		return nil, err
	}
	model := client.GenerativeModel("gemini-1.5-flash")
	return &GeminiLLM{
		client: client,
		model:  model,
	}, nil
}

func (l *GeminiLLM) StreamComplete(ctx context.Context, prompt string, onToken func(string)) (string, error) {
	iter := l.model.GenerateContentStream(ctx, genai.Text(prompt))
	var fullText strings.Builder
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.L.Error("gemini stream failed", "err", err)
			return "", err
		}
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					token := fmt.Sprintf("%v", part)
					fullText.WriteString(token)
					onToken(token)
				}
			}
		}
	}
	return fullText.String(), nil
}
