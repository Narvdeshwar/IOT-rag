package ingestor

import (
	"context"
	"sync"

	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/telemetry"
	"github.com/narvdeshwar/IOT-rag/internal/types"
	"github.com/pgvector/pgvector-go"
)

type Worker struct {
	embedder types.Embedder
}

func NewWorker(emb types.Embedder) *Worker {
	return &Worker{
		embedder: emb,
	}
}

func (w *Worker) ProcessBatch(ctx context.Context) int {
	rows, err := db.Pool.Query(ctx, "SELECT id, device_id, event_time, metric, value, unit, metadata from sensor_events WHERE id NOT IN (SELECT event_id from telemetry_chunks) LIMIT 100")
	if err != nil {
		logger.L.Error("failed to query sensor events", "err", err)
		return 0
	}
	defer rows.Close()

	var events []telemetry.SensorEvent
	for rows.Next() {
		var e telemetry.SensorEvent
		err := rows.Scan(&e.ID, &e.DeviceID, &e.EventTime, &e.Metric, &e.Value, &e.Unit, &e.Metadata)
		if err != nil {
			logger.L.Error("failed to scan event", "err", err)
			continue
		}
		events = append(events, e)
	}

	if len(events) == 0 {
		return 0
	}

	// Use concurrency for embeddings
	numWorkers := 10
	if len(events) < numWorkers {
		numWorkers = len(events)
	}

	type result struct {
		event telemetry.SensorEvent
		text  string
		vec   []float32
		err   error
	}

	jobs := make(chan telemetry.SensorEvent, len(events))
	results := make(chan result, len(events))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range jobs {
				text := ChunkEvent(e)
				vec, err := w.embedder.Embed(ctx, text)
				results <- result{event: e, text: text, vec: vec, err: err}
			}
		}()
	}

	// Send jobs
	for _, e := range events {
		jobs <- e
	}
	close(jobs)

	// Wait for workers in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	processedCount := 0
	for res := range results {
		if res.err != nil {
			logger.L.Error("Skipping event due to embed error", "id", res.event.ID, "err", res.err)
			continue
		}

		_, err = db.Pool.Exec(ctx, `
			INSERT INTO telemetry_chunks (content, device_id, event_time, metadata, event_id, embedding)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (event_id) DO NOTHING`,
			res.text, res.event.DeviceID, res.event.EventTime, res.event.Metadata, res.event.ID, pgvector.NewVector(res.vec))
		if err != nil {
			logger.L.Error("failed to insert chunk", "id", res.event.ID, "err", err)
		} else {
			processedCount++
		}
	}

	logger.L.Info("batch processed", "count", processedCount)
	return processedCount
}
