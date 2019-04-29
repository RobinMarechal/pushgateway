package lib

import (
	"github.com/prometheus/common/log"
	"github.com/prometheus/pushgateway/storage"
	"time"
)

func ClearMetrics(ms storage.MetricStore) {
	mfs := ms.GetMetricFamiliesMap()
	for _, labels := range mfs {
		ms.SubmitWriteRequest(storage.WriteRequest{
			Labels:    labels.Labels,
			Timestamp: time.Now(),
		})
	}

	log.Debug("Cleared metrics.")
}
