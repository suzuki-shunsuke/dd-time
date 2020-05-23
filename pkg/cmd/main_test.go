package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suzuki-shunsuke/go-ptr"
	"github.com/suzuki-shunsuke/gomic/gomic"
	"github.com/zorkian/go-datadog-api"

	"github.com/suzuki-shunsuke/dd-time/pkg/mock"
)

func Test_validateParams(t *testing.T) {
	data := []struct {
		title  string
		isErr  bool
		params Params
	}{
		{
			title: "DATADOG_API_KEY is required",
			isErr: true,
		},
		{
			title: "executed command isn't passed to dd-time",
			isErr: true,
			params: Params{
				DataDogAPIKey: "xxx",
			},
		},
		{
			title: "success",
			params: Params{
				DataDogAPIKey: "xxx",
				Args:          []string{"sleep", "5"},
			},
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := validateParams(d.params)
			if d.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_send(t *testing.T) {
	data := []struct {
		title    string
		isErr    bool
		ddClient metricsPoster
		metrics  []datadog.Metric
	}{
		{
			title:    "success",
			ddClient: mock.NewMetricsPoster(t, gomic.DoNothing),
			metrics:  []datadog.Metric{},
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := send(d.metrics, d.ddClient)
			if d.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_getMetrics(t *testing.T) {
	data := []struct {
		title    string
		exp      []datadog.Metric
		duration float64
		now      time.Time
		params   Params
	}{
		{
			title:    "success",
			duration: 5,
			now:      time.Date(2019, 2, 10, 12, 0, 0, 0, time.UTC),
			params: Params{
				MetricName: "command-execution-time",
			},
			exp: []datadog.Metric{
				{
					Metric: ptr.PStr("command-execution-time"),
					Points: []datadog.DataPoint{{ptr.PFloat64(1.5498e+09), ptr.PFloat64(5)}},
				},
			},
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			assert.Equal(t, d.exp, getMetrics(d.duration, d.now, d.params))
		})
	}
}
