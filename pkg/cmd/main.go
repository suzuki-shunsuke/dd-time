package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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

type withExitCodeError struct {
	err  error
	code int
}

func wrapWithExitCode(err error, code int) error {
	return &withExitCodeError{
		err:  err,
		code: code,
	}
}

func (err *withExitCodeError) ExitCode() int {
	return err.code
}

func (err *withExitCodeError) Error() string {
	return err.err.Error()
}

func (err *withExitCodeError) Unwrap() error {
	return err.err
}

func GetExitCode(err error) int {
	var ecerr *withExitCodeError
	if errors.As(err, &ecerr) {
		return ecerr.ExitCode()
	}
	return 1
}

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

	var ddOutput io.Writer = os.Stderr
	if params.Output != "" {
		switch params.Output {
		case "/dev/stderr":
		case "/dev/null":
			f, err := os.Open("/dev/null")
			if err == nil {
				defer f.Close()
				ddOutput = f
			} else {
				ddOutput = bytes.NewBufferString("")
			}
		case "/dev/stdout":
			ddOutput = os.Stdout
		default:
			var (
				f   io.WriteCloser
				err error
			)
			if params.Append {
				f, err = os.OpenFile(params.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			} else {
				f, err = os.Open(params.Output)
			}
			if err == nil {
				defer f.Close()
				ddOutput = f
			} else {
				ddOutput = bytes.NewBufferString("")
			}
		}
	}

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
				return wrapWithExitCode(err, cmd.ProcessState.ExitCode())
			}
			return send(getMetrics(duration, time.Now(), params), ddClient, ddOutput)
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

func send(metrics []datadog.Metric, ddClient metricsPoster, ddOutput io.Writer) error {
	if err := ddClient.PostMetrics(metrics); err != nil {
		fmt.Fprintln(ddOutput, "send a time series metrics to DataDog:", err)
		return nil
	}
	return nil
}
