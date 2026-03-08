package retriever

import (
	"context"

	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/telemetry"
	"github.com/pgvector/pgvector-go"
)

type Retriever struct{}

func (r *Retriever) Search(ctx context.Context, vec []float32, topK int) ([]telemetry.EmbeddedChunk, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, content, device_id, event_time, metadata,
		       1 - (embedding <=> $1::vector) AS similarity
		FROM telemetry_chunks
		WHERE 1 - (embedding <=> $1::vector) > 0.75
		ORDER BY embedding <=> $1::vector
		LIMIT $2`, pgvector.NewVector(vec), topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []telemetry.EmbeddedChunk
	for rows.Next() {
		var c telemetry.EmbeddedChunk
		rows.Scan(&c.ID, &c.Content, &c.DeviceID, &c.EventTime, &c.Metadata)
		chunks = append(chunks, c)
	}
	return chunks, nil
}
