package lib

import (
	"github.com/prometheus/common/log"
	"github.com/prometheus/pushgateway/storage"
	"time"
)

func ScheduledClear(ms storage.MetricStore, interval *time.Duration) {
	log.Info("Starting clearing scheduler every " + interval.String() + "m.")
	for range time.Tick(*interval) {
		ClearMetrics(ms)
	}
}
