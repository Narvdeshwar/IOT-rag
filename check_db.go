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
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	var count int
	err = pool.QueryRow(context.Background(), "SELECT count(*) FROM telemetry_chunks").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Telemetry Chunks: %d\n", count)

	err = pool.QueryRow(context.Background(), "SELECT count(*) FROM sensor_events").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sensor Events: %d\n", count)

	err = pool.QueryRow(context.Background(), "SELECT count(*) FROM telemetry_chunks WHERE embedding IS NOT NULL").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Chunks with Embeddings: %d\n", count)
}
