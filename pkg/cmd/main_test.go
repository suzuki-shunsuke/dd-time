package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suzuki-shunsuke/gomic/gomic"

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
		duration float64
		now      time.Time
		ddClient metricsPoster
		params   Params
	}{
		{
			title:    "success",
			duration: 5,
			now:      time.Date(2019, 2, 10, 12, 0, 0, 0, time.UTC),
			ddClient: mock.NewMetricsPoster(t, gomic.DoNothing),
			params: Params{
				MetricName: "command-execution-time",
			},
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := send(d.duration, d.now, d.params, d.ddClient)
			if d.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
