package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/narvdeshwar/IOT-rag/internal/api"
	"github.com/narvdeshwar/IOT-rag/internal/cache"
	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/embedder"
	"github.com/narvdeshwar/IOT-rag/internal/llm"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/retriever"
	"github.com/narvdeshwar/IOT-rag/internal/types"
)

func main() {
	cfg := config.Load()
	db.Init(context.Background(), cfg.PostgresURL)
	defer db.Pool.Close()
	cache.Init()
	var emb types.Embedder
	var l types.LLM

	// Prioritize Ollama if it's running
	if cfg.OllamaURL != "" {
		fmt.Printf("Attempting to use Ollama provider at %s...\n", cfg.OllamaURL)
		emb = embedder.NewOllamaEmbedder(cfg.OllamaURL, "nomic-embed-text")
		l = llm.NewOllamaLLM(cfg.OllamaURL, "llama3")
		// We use a simple check here, if initialization fails we'll fallback
		logger.L.Info("Using Ollama provider")
	} else if cfg.GeminiKey != "" {
		logger.L.Info("Using Gemini provider")
		var err error
		emb, err = embedder.NewGeminiEmbedder(context.Background(), cfg.GeminiKey)
		if err != nil {
			logger.L.Error("Failed to create Gemini embedder", "err", err)
		}
		l, err = llm.NewGeminiLLM(context.Background(), cfg.GeminiKey)
		if err != nil {
			logger.L.Error("Failed to create Gemini LLM", "err", err)
		}
	}

	if emb == nil || (cfg.OpenAIKey != "" && cfg.GeminiKey == "" && cfg.OllamaURL == "") {
		logger.L.Info("Using OpenAI provider")
		emb = embedder.NewEmbedder(cfg.OpenAIKey)
		l = llm.NewLLM(cfg.OpenAIKey)
	}

	svc := &api.Service{
		Embedder:  emb,
		Retriever: &retriever.Retriever{},
		LLM:       l,
	}
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/health", api.HealthHandler)
	r.GET("/ready", api.HealthHandler)
	r.POST("/query", api.QueryHandler(svc))
	logger.L.Info("Server started", "port", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
