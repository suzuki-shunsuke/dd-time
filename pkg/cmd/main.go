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

	ddOutput, closeOutput := getDDOutput(params.Output, params.Append)
	if closeOutput != nil {
		defer closeOutput()
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

func getCircleCITags() []string {
	// https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
	return []string{
		"circleci:true",
		"circle_branch:" + os.Getenv("CIRCLE_BRANCH"),
		"circle_build_num:" + os.Getenv("CIRCLE_BUILD_NUM"),
		"circle_build_url:" + os.Getenv("CIRCLE_BUILD_URL"),
		"circle_job:" + os.Getenv("CIRCLE_JOB"),
		"circle_node_index:" + os.Getenv("CIRCLE_NODE_INDEX"),
		"circle_node_total:" + os.Getenv("CIRCLE_NODE_TOTAL"),
		"circle_pr_number:" + os.Getenv("CIRCLE_PR_NUMBER"),
		"circle_pr_reponame:" + os.Getenv("CIRCLE_PR_REPONAME"),
		"circle_pr_username:" + os.Getenv("CIRCLE_PR_USERNAME"),
		"circle_project_reponame:" + os.Getenv("CIRCLE_PROJECT_REPONAME"),
		"circle_project_username:" + os.Getenv("CIRCLE_PROJECT_USERNAME"),
		"circle_repository_url:" + os.Getenv("CIRCLE_REPOSITORY_URL"),
		"circle_sha1:" + os.Getenv("CIRCLE_SHA1"),
		"circle_tag:" + os.Getenv("CIRCLE_TAG"),
		"circle_username:" + os.Getenv("CIRCLE_USERNAME"),
		"circle_workflow_id:" + os.Getenv("CIRCLE_WORKFLOW_ID"),
	}
}

func getTags() []string {
	if os.Getenv("CIRCLECI") == "true" {
		return getCircleCITags()
	}
	return nil
}

func getMetrics(
	duration float64, now time.Time, params Params,
) []datadog.Metric {
	nowF := float64(now.Unix())
	metric := datadog.Metric{
		Metric: &params.MetricName,
		Tags:   append(params.Tags, getTags()...),
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
