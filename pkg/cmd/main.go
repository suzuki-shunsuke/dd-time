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

func Main(params Params) error {
	if len(params.Args) == 0 {
		return errors.New("executed command isn't passed to dd-time")
	}

	ddClient := datadog.NewClient(params.DataDogAPIKey, "")

	cmd := exec.Command(params.Args[0], params.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
			duration := float64(time.Since(startT))
			if err != nil {
				return err
			}
			return send(duration, params, ddClient)
		case sig := <-signalChan:
			if _, ok := sentSignals[sig]; ok {
				continue
			}
			sentSignals[sig] = struct{}{}
			runner.SendSignal(sig.(syscall.Signal))
		}
	}
}

func send(duration float64, params Params, ddClient metricsPoster) error {
	now := float64(time.Now().Unix())
	metric := datadog.Metric{
		Metric: &params.MetricName,
		Tags:   params.Tags,
		Points: []datadog.DataPoint{{&now, &duration}},
	}
	if params.MetricHost != "" {
		metric.Host = &params.MetricHost
	}
	if err := ddClient.PostMetrics([]datadog.Metric{metric}); err != nil {
		fmt.Fprintln(os.Stderr, "send a time series metrics to DataDog: %w", err)
		return nil
	}

	return nil
}
