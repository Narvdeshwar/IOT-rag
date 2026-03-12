package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/narvdeshwar/IOT-rag/internal/types"
)

var _ types.Embedder = (*OllamaEmbedder)(nil)

type OllamaEmbedder struct {
	url   string
	model string
}

func NewOllamaEmbedder(url, model string) *OllamaEmbedder {
	return &OllamaEmbedder{
		url:   url,
		model: model,
	}
}

type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	reqBody, _ := json.Marshal(ollamaEmbedRequest{
		Model:  e.model,
		Prompt: text,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/embeddings", e.url), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama embed failed: %s", resp.Status)
	}

	var res ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Embedding, nil
}
