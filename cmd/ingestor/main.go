package main

import (
	"context"
	"time"

	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/ingestor"
)

func main() {
	cfg := config.Load()
	db.Init(context.Background(), cfg.PostgresURL)
	defer db.Pool.Close()

	w := ingestor.NewWorker()
	for {
		w.ProcessBatch(context.Background())
		time.Sleep(30 * time.Second)
	}
}
