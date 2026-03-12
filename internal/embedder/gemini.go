package embedder

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/types"
	"google.golang.org/api/option"
)

var _ types.Embedder = (*GeminiEmbedder)(nil)

type GeminiEmbedder struct {
	client *genai.Client
	model  *genai.EmbeddingModel
}

func NewGeminiEmbedder(ctx context.Context, apiKey string) (*GeminiEmbedder, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.L.Error("failed to create gemini client", "err", err)
		return nil, err
	}
	model := client.EmbeddingModel("embedding-001")
	return &GeminiEmbedder{
		client: client,
		model:  model,
	}, nil
}

func (e *GeminiEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	res, err := e.model.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		logger.L.Error("gemini embed failed", "err", err)
		return nil, err
	}
	return res.Embedding.Values, nil
}
