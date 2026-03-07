package ingestor

import (
	"context"

	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/embedder"
	"github.com/narvdeshwar/IOT-rag/internal/logger"
	"github.com/narvdeshwar/IOT-rag/internal/telemetry"
)

type Worker struct {
	embedder *embedder.OpenAIEmbedder
}

func NewWorker() *Worker {
	return &Worker{
		embedder: embedder.NewEmbedder(),
	}
}

func (w *Worker) ProcessBatch(ctx context.Context) {
	rows, _ := db.Pool.Query(ctx, "SELECT *from sensor_events WHERE id NOT IN (SELECT id from telemetry_chunks) LIMIT 100")
	defer rows.Close()
	var events []telemetry.SensorEvent
	for rows.Next() {
		var e telemetry.SensorEvent
		rows.Scan(&e.ID, &e.DeviceID, &e.EventTime, &e.Metric, &e.Value, &e.Unit, &e.Metadata)
		events = append(events, e)
	}

	for _, e := range events {
		text := ChunkEvent(e)
		vec, err := w.embedder.Embed(ctx, text)
		if err != nil {
			logger.L.Error("Skipping event due to embed error", "id", e.ID, "err", err)
			continue
		}
		_, _ = db.Pool.Exec(ctx, `
			INSERT INTO telemetry_chunks (content, device_id, event_time, metadata, embedding)
			VALUES ($1,$2,$3,$4,$5)`,
			text, e.DeviceID, e.EventTime, e.Metadata, vec)
	}
	logger.L.Info("batch embedded", "count", len(events))
}
