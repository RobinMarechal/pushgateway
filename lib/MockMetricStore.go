package lib

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/pushgateway/storage"
)

type MockMetricStore struct {
	LastWriteRequest storage.WriteRequest
}

func (m *MockMetricStore) SubmitWriteRequest(req storage.WriteRequest) {
	m.LastWriteRequest = req
}

func (m *MockMetricStore) GetMetricFamilies() []*dto.MetricFamily {
	panic("not implemented")
}

func (m *MockMetricStore) GetMetricFamiliesMap() storage.GroupingKeyToMetricGroup {
	panic("not implemented")
}

func (m *MockMetricStore) Shutdown() error {
	return nil
}

func (m *MockMetricStore) Healthy() error {
	return nil
}

func (m *MockMetricStore) Ready() error {
	return nil
}

