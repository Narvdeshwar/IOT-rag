package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/narvdeshwar/IOT-rag/internal/api"
	"github.com/narvdeshwar/IOT-rag/internal/cache"
	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/embedder"
	"github.com/narvdeshwar/IOT-rag/internal/llm"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/retriever"
)

func main() {
	cfg := config.Load()
	db.Init(context.Background(), cfg.PostgresURL)
	defer db.Pool.Close()
	cache.Init()
	svc := &api.Service{
		Embedder:  embedder.NewEmbedder(),
		Retriever: &retriever.Retriever{},
		LLM:       llm.NewLLM(),
	}
	r := gin.Default()
	r.GET("/health", api.HealthHandler)
	r.GET("/ready", api.HealthHandler)
	r.POST("/query", api.QueryHandler(svc))
	logger.L.Info("Server started", "port", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
