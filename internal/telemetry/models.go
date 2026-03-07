package telemetry

import (
	"encoding/json"
	"time"

	"github.com/pgvector/pgvector-go"
)

type Device struct {
	ID          string    `json:"device_id"`
	Name        string    `json:"name"`
	Zone        string    `json:"zone"`
	Type        string    `json:"type"`
	LastSeen    time.Time `json:"last_seen"`
}

type SensorEvent struct {
	ID         int64           `json:"id"`
	DeviceID   string          `json:"device_id"`
	EventTime  time.Time       `json:"event_time"`
	Metric     string          `json:"metric"`
	Value      float64         `json:"value"`
	Unit       string          `json:"unit"`
	Metadata   json.RawMessage `json:"metadata"`
}

type EmbeddedChunk struct {
	ID         int64           `json:"id"`
	Content    string          `json:"content"`
	DeviceID   string          `json:"device_id"`
	EventTime  time.Time       `json:"event_time"`
	Metadata   json.RawMessage `json:"metadata"`
	Embedding  pgvector.Vector `json:"embedding"`
}
