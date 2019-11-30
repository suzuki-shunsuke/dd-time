package mock

// Don't edit this file.
// This file is generated by gomic 0.6.0.
// https://github.com/suzuki-shunsuke/gomic

import (
	testing "testing"

	gomic "github.com/suzuki-shunsuke/gomic/gomic"
	"github.com/zorkian/go-datadog-api"
)

type (
	// MetricsPoster is a mock.
	MetricsPoster struct {
		t                      *testing.T
		name                   string
		callbackNotImplemented gomic.CallbackNotImplemented
		impl                   struct {
			PostMetrics func(series []datadog.Metric) (r0 error)
		}
	}
)

// NewMetricsPoster returns MetricsPoster .
func NewMetricsPoster(t *testing.T, cb gomic.CallbackNotImplemented) *MetricsPoster {
	return &MetricsPoster{
		t: t, name: "MetricsPoster", callbackNotImplemented: cb}
}

// PostMetrics is a mock method.
func (mock MetricsPoster) PostMetrics(series []datadog.Metric) (r0 error) {
	methodName := "PostMetrics" // nolint: goconst
	if mock.impl.PostMetrics != nil {
		return mock.impl.PostMetrics(series)
	}
	if mock.callbackNotImplemented != nil {
		mock.callbackNotImplemented(mock.t, mock.name, methodName)
	} else {
		gomic.DefaultCallbackNotImplemented(mock.t, mock.name, methodName)
	}
	return mock.fakeZeroPostMetrics(series)
}

// SetFuncPostMetrics sets a method and returns the mock.
func (mock *MetricsPoster) SetFuncPostMetrics(impl func(series []datadog.Metric) (r0 error)) *MetricsPoster {
	mock.impl.PostMetrics = impl
	return mock
}

// SetReturnPostMetrics sets a fake method.
func (mock *MetricsPoster) SetReturnPostMetrics(r0 error) *MetricsPoster {
	mock.impl.PostMetrics = func([]datadog.Metric) error {
		return r0
	}
	return mock
}

// fakeZeroPostMetrics is a fake method which returns zero values.
func (mock MetricsPoster) fakeZeroPostMetrics(series []datadog.Metric) (r0 error) {
	return r0
}
