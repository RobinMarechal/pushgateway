package lib

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/pushgateway/storage"
	"math"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
)

var (
	mf1 = &dto.MetricFamily{
		Name: proto.String("mf2"),
		Help: proto.String("doc string 2"),
		Type: dto.MetricType_GAUGE.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("job"),
						Value: proto.String("job1"),
					},
					{
						Name:  proto.String("instance"),
						Value: proto.String("instance2"),
					},
				},
				Gauge: &dto.Gauge{
					Value: proto.Float64(math.Inf(+1)),
				},
				TimestampMs: proto.Int64(54321),
			},
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("job"),
						Value: proto.String("job1"),
					},
					{
						Name:  proto.String("instance"),
						Value: proto.String("instance2"),
					},
				},
				Gauge: &dto.Gauge{
					Value: proto.Float64(math.Inf(-1)),
				},
			},
		},
	}
	mf2 = &dto.MetricFamily{
		Name: proto.String("mf3"),
		Type: dto.MetricType_UNTYPED.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("job"),
						Value: proto.String("job1"),
					},
					{
						Name:  proto.String("instance"),
						Value: proto.String("instance1"),
					},
				},
				Untyped: &dto.Untyped{
					Value: proto.Float64(42),
				},
			},
		},
	}
	mf3 = &dto.MetricFamily{
		Name: proto.String("mf4"),
		Type: dto.MetricType_UNTYPED.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("job"),
						Value: proto.String("job3"),
					},
					{
						Name:  proto.String("instance"),
						Value: proto.String("instance2"),
					},
				},
				Untyped: &dto.Untyped{
					Value: proto.Float64(3.4345),
				},
			},
		},
	}
)

func TestClearMetrics(t *testing.T) {
	dms := storage.NewDiskMetricStore("", time.Minute, nil)

	labels1 := map[string]string{
		"job":      "job1",
		"instance": "instance1",
	}

	labels2 := map[string]string{
		"job":      "job2",
		"instance": "instance2",
	}

	// Submit a single simple metric family.
	ts1 := time.Now()
	dms.SubmitWriteRequest(storage.WriteRequest{
		Labels:         labels1,
		Timestamp:      ts1,
		MetricFamilies: map[string]*dto.MetricFamily{"mf3": mf3},
	})
	time.Sleep(20 * time.Millisecond) // Give loop() time to process.

	// Submit two metric families for a different instance.
	ts2 := ts1.Add(time.Second)
	dms.SubmitWriteRequest(storage.WriteRequest{
		Labels:         labels2,
		Timestamp:      ts2,
		MetricFamilies: map[string]*dto.MetricFamily{"mf1": mf1, "mf2": mf2},
	})
	time.Sleep(20 * time.Millisecond) // Give loop() time to process.

	ClearMetrics(dms)

	time.Sleep(20 * time.Millisecond) // Give loop() time to process.

	if expected, got := 0, len(dms.GetMetricFamiliesMap()); got != expected {
		t.Errorf("Wanted mfs length %v, got %v.", expected, got)
	}
}