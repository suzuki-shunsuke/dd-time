package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	ddog "github.com/suzuki-shunsuke/dd-time/pkg/datadog"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/go-timeout/timeout"
	"github.com/zorkian/go-datadog-api"
)

func Main(ctx context.Context, params Params) int {
	ddOutput, closeOutput := getDDOutput(params.Output, params.Append)
	if closeOutput != nil {
		defer closeOutput()
	}
	msg, code := core(ctx, params)
	if msg != "" {
		fmt.Fprintln(ddOutput, msg)
	}
	return code
}

func core(ctx context.Context, params Params) (string, int) {
	if err := validateParams(params); err != nil {
		return err.Error(), 1
	}

	var ddClient metricsPoster
	if params.DataDogAPIKey != "" {
		ddClient = datadog.NewClient(params.DataDogAPIKey, "")
	}

	duration, err := execute(ctx, &params)
	if err != nil {
		return err.Error(), ecerror.GetExitCode(err)
	}

	if ddClient == nil {
		return "", 0
	}
	if err := send(getMetrics(duration, time.Now(), params), ddClient); err != nil {
		return err.Error(), 0
	}
	return "", 0
}

type (
	Params struct {
		Append        bool
		MetricName    string
		MetricHost    string
		DataDogAPIKey string
		Output        string
		Args          []string
		Tags          []string
	}

	metricsPoster interface {
		PostMetrics(series []datadog.Metric) error
	}
)

func validateParams(params Params) error {
	if len(params.Args) == 0 {
		return errors.New("executed command isn't passed to dd-time")
	}
	return nil
}

func getDDOutput(output string, appended bool) (io.Writer, func() error) {
	switch output {
	case "", "/dev/stderr":
		return os.Stderr, nil
	case "/dev/null":
		return bytes.NewBufferString(""), nil
	case "/dev/stdout":
		return os.Stdout, nil
	default:
		var (
			f   io.WriteCloser
			err error
		)
		if appended {
			f, err = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			f, err = os.Open(output)
		}
		if err != nil {
			return bytes.NewBufferString(""), nil
		}
		return f, f.Close
	}
}

func execute(ctx context.Context, params *Params) (float64, error) {
	cmd := exec.Command(params.Args[0], params.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	runner := timeout.NewRunner(0)

	startT := time.Now()
	if err := runner.Run(ctx, cmd); err != nil {
		return 0, ecerror.Wrap(err, cmd.ProcessState.ExitCode())
	}
	return time.Since(startT).Seconds(), nil
}

func getMetrics(
	duration float64, now time.Time, params Params,
) []datadog.Metric {
	nowF := float64(now.Unix())
	metric := datadog.Metric{
		Metric: &params.MetricName,
		Tags:   append(params.Tags, ddog.GetTags()...),
		Points: []datadog.DataPoint{{&nowF, &duration}},
	}
	if params.MetricHost != "" {
		metric.Host = &params.MetricHost
	}
	return []datadog.Metric{metric}
}

func send(metrics []datadog.Metric, ddClient metricsPoster) error {
	if err := ddClient.PostMetrics(metrics); err != nil {
		return fmt.Errorf("send a time series metrics to DataDog: %w", err)
	}
	return nil
}
