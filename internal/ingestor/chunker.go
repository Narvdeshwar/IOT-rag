package ingestor

import (
	"fmt"

	"github.com/narvdeshwar/IOT-rag/internal/telemetry"
)

func ChunkEvent(e telemetry.SensorEvent) string {
	return fmt.Sprintf(`Device %s | %s: %.2f %s at %s | Metadata: %s`,
		e.DeviceID, e.Metric, e.Value, e.Unit, e.EventTime.Format("2006-01-02 15:04"), e.Metadata)
}
