package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/narvdeshwar/IOT-rag/internal/api"
	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
)

func main() {
	cfg := config.Load()
	db.Init(context.Background(), cfg.PostgresURL)
	defer db.Pool.Close()

	r := gin.Default()
	r.GET("/health", api.HealthHandler)
	r.GET("/ready", api.HealthHandler)
	logger.L.Info("Server started", "port", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
