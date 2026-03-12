package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		log.Fatal("POSTGRES_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer pool.Close()

	stmts := []string{
		"DROP INDEX IF EXISTS telemetry_chunks_embedding_idx",
		"ALTER TABLE telemetry_chunks DROP COLUMN IF EXISTS embedding",
		"ALTER TABLE telemetry_chunks DROP COLUMN IF EXISTS event_id",
		"ALTER TABLE telemetry_chunks ADD COLUMN event_id BIGINT UNIQUE",
		"ALTER TABLE telemetry_chunks ADD COLUMN embedding VECTOR(768)",
		"CREATE INDEX telemetry_chunks_embedding_idx ON telemetry_chunks USING hnsw (embedding vector_cosine_ops)",
	}

	for _, stmt := range stmts {
		if _, err := pool.Exec(context.Background(), stmt); err != nil {
			log.Fatalf("failed: %s\n  error: %v", stmt, err)
		}
		fmt.Println("OK:", stmt)
	}

	fmt.Println("\nMigration complete. Column 'embedding' is now VECTOR(768) for Gemini embedding-001.")
}
