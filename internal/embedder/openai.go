package embedder

import (
	"context"

	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/types"
	"github.com/sashabaranov/go-openai"
)

var _ types.Embedder = (*OpenAIEmbedder)(nil)

type OpenAIEmbedder struct {
	client *openai.Client
}

func NewEmbedder(apiKey string) *OpenAIEmbedder {
	return &OpenAIEmbedder{
		client: openai.NewClient(apiKey),
	}
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	res, err := e.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel("text-embedding-3-small"),
	})
	if err != nil {
		logger.L.Error("embed failed", "err", err)
		return nil, err
	}
	return res.Data[0].Embedding, nil
}
