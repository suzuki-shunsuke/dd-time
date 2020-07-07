package datadog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suzuki-shunsuke/go-ptr"
	"github.com/zorkian/go-datadog-api"
)

func Test_getMetrics(t *testing.T) {
	data := []struct {
		title  string
		exp    []datadog.Metric
		params Params
	}{
		{
			title: "success",
			params: Params{
				MetricName: "command-execution-time",
				Duration:   5,
				Now:        float64(time.Date(2019, 2, 10, 12, 0, 0, 0, time.UTC).Unix()),
			},
			exp: []datadog.Metric{
				{
					Metric: ptr.PStr("command-execution-time"),
					Points: []datadog.DataPoint{{ptr.PFloat64(1.5498e+09), ptr.PFloat64(5)}},
				},
			},
		},
	}
	client := &Client{}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			assert.Equal(t, d.exp, client.getMetrics(d.params))
		})
	}
}
