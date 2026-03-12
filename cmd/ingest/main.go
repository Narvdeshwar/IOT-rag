package main

import (
	"context"
	"fmt"
	"log"

	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/embedder"
	"github.com/narvdeshwar/IOT-rag/internal/ingestor"
	"github.com/narvdeshwar/IOT-rag/internal/types"
)

func main() {
	cfg := config.Load()
	db.Init(context.Background(), cfg.PostgresURL)
	defer db.Pool.Close()

	var emb types.Embedder
	// Prioritize Ollama
	if cfg.OllamaURL != "" {
		emb = embedder.NewOllamaEmbedder(cfg.OllamaURL, "nomic-embed-text")
		fmt.Println("Using Ollama for embeddings...")
	} else if cfg.GeminiKey != "" {
		var err error
		emb, err = embedder.NewGeminiEmbedder(context.Background(), cfg.GeminiKey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Using Gemini for embeddings...")
	} else {
		emb = embedder.NewEmbedder(cfg.OpenAIKey)
		fmt.Println("Using OpenAI for embeddings...")
	}

	worker := ingestor.NewWorker(emb)
	ctx := context.Background()

	total := 0
	for {
		processed := worker.ProcessBatch(ctx)
		if processed == 0 {
			break
		}
		total += processed
		fmt.Printf("Processed %d events...\n", total)
	}

	fmt.Printf("\nIngestion complete! Total events processed and embedded: %d\n", total)
}
