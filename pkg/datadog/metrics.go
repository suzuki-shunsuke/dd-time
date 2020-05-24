package datadog

import (
	"github.com/zorkian/go-datadog-api"
)

type Params struct {
	MetricName string
	MetricHost string
	Tags       []string
	Duration   float64
	Now        float64
}

type Client struct {
	Poster MetricsPoster
}

type MetricsPoster interface {
	PostMetrics(series []datadog.Metric) error
}

func New(apiKey string) *Client {
	var ddClient MetricsPoster
	if apiKey != "" {
		ddClient = datadog.NewClient(apiKey, "")
	}
	return &Client{
		Poster: ddClient,
	}
}

func (client *Client) Send(params *Params) error {
	if client.Poster == nil {
		return nil
	}
	return client.Poster.PostMetrics(client.getMetrics(params))
}

func (client *Client) getMetrics(params *Params) []datadog.Metric {
	metric := datadog.Metric{
		Metric: &params.MetricName,
		Tags:   params.Tags,
		Points: []datadog.DataPoint{{&params.Now, &params.Duration}},
	}
	if params.MetricHost != "" {
		metric.Host = &params.MetricHost
	}
	return []datadog.Metric{metric}
}
