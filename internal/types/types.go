package types

import (
	"context"
)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type LLM interface {
	StreamComplete(ctx context.Context, prompt string, onToken func(string)) (string, error)
}
