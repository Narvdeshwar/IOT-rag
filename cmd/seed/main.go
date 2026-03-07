package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/narvdeshwar/IOT-rag/internal/db"
	"github.com/narvdeshwar/IOT-rag/internal/telemetry"
)

func main() {
	cfg := config.Load()
	db.Init(context.Background(), cfg.PostgresURL)
	defer db.Pool.Close()

	ctx := context.Background()

	// STEP 1: seed devices
	for i := 1; i <= 50; i++ {
		deviceID := fmt.Sprintf("dev-%03d", i)

		_, err := db.Pool.Exec(ctx, `
		INSERT INTO devices (id, name, zone, type)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (id) DO NOTHING
		`,
			deviceID,
			fmt.Sprintf("Device %03d", i),
			"Zone 3",
			"sensor",
		)

		if err != nil {
			panic(err)
		}
	}

	// STEP 2: seed sensor events
	for i := 0; i < 1000; i++ {

		event := telemetry.SensorEvent{
			DeviceID:  fmt.Sprintf("dev-%03d", rand.Intn(50)+1),
			EventTime: time.Now().Add(-time.Duration(rand.Intn(7*24)) * time.Hour),
			Metric:    []string{"voltage", "current", "temp", "faults"}[rand.Intn(4)],
			Value:     rand.Float64() * 100,
			Unit:      []string{"V", "A", "°C", "count"}[rand.Intn(4)],
			Metadata:  json.RawMessage(`{"zone":"Zone 3","status":"active"}`),
		}

		_, err := db.Pool.Exec(ctx, `
		INSERT INTO sensor_events (device_id, event_time, metric, value, unit, metadata)
		VALUES ($1,$2,$3,$4,$5,$6)
		`,
			event.DeviceID,
			event.EventTime,
			event.Metric,
			event.Value,
			event.Unit,
			event.Metadata,
		)

		if err != nil {
			panic(err)
		}
	}

	fmt.Println("✅ Seeded 50 devices + 1000 events")
}
