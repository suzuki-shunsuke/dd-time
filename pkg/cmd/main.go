package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/suzuki-shunsuke/go-timeout/timeout"
	"github.com/zorkian/go-datadog-api"
)

type (
	Params struct {
		MetricName    string
		MetricHost    string
		DataDogAPIKey string
		Args          []string
		Tags          []string
	}

	metricsPoster interface {
		PostMetrics(series []datadog.Metric) error
	}
)

func validateParams(params Params) error {
	if params.DataDogAPIKey == "" {
		return errors.New("The environment variable 'DATADOG_API_KEY' is required")
	}
	if len(params.Args) == 0 {
		return errors.New("executed command isn't passed to dd-time")
	}
	return nil
}

func Main(params Params) error {
	if err := validateParams(params); err != nil {
		return err
	}

	ddClient := datadog.NewClient(params.DataDogAPIKey, "")

	cmd := exec.Command(params.Args[0], params.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)

	runner := timeout.NewRunner(0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentSignals := map[os.Signal]struct{}{}
	exitChan := make(chan error, 1)

	var startT time.Time
	go func() {
		startT = time.Now()
		exitChan <- runner.Run(ctx, cmd)
	}()

	for {
		select {
		case err := <-exitChan:
			duration := time.Since(startT).Seconds()
			if err != nil {
				return err
			}
			return send(getMetrics(duration, time.Now(), params), ddClient)
		case sig := <-signalChan:
			if _, ok := sentSignals[sig]; ok {
				continue
			}
			sentSignals[sig] = struct{}{}
			runner.SendSignal(sig.(syscall.Signal))
		}
	}
}

func getMetrics(
	duration float64, now time.Time, params Params,
) []datadog.Metric {
	nowF := float64(now.Unix())
	metric := datadog.Metric{
		Metric: &params.MetricName,
		Tags:   params.Tags,
		Points: []datadog.DataPoint{{&nowF, &duration}},
	}
	if params.MetricHost != "" {
		metric.Host = &params.MetricHost
	}
	return []datadog.Metric{metric}
}

func send(metrics []datadog.Metric, ddClient metricsPoster) error {
	if err := ddClient.PostMetrics(metrics); err != nil {
		fmt.Fprintln(os.Stderr, "send a time series metrics to DataDog: %w", err)
		return nil
	}
	return nil
}
