package lib

import (
	"github.com/prometheus/pushgateway/storage"
	"time"
)

func ClearMetrics(ms storage.MetricStore) {
	ms.SubmitWriteRequest(storage.WriteRequest{
		Labels:    labels,
		Timestamp: time.Now(),
	})
}
