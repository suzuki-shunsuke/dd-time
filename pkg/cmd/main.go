package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
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

func Main(ctx context.Context, params Params) error {
	if err := validateParams(params); err != nil {
		return err
	}

	var ddClient metricsPoster
	if params.DataDogAPIKey != "" {
		ddClient = datadog.NewClient(params.DataDogAPIKey, "")
	}

	ddOutput, closeOutput := getDDOutput(params.Output, params.Append)
	if closeOutput != nil {
		defer closeOutput()
	}

	duration, err := execute(ctx, &params)
	if err != nil {
		return err
	}

	if ddClient == nil {
		return nil
	}
	return send(getMetrics(duration, time.Now(), params), ddClient, ddOutput)
}

func getCircleCITags() []string {
	// https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables
	envs := []string{
		"CIRCLECI",
		"CIRCLE_BRANCH",
		"CIRCLE_BUILD_NUM",
		"CIRCLE_BUILD_URL",
		"CIRCLE_JOB",
		"CIRCLE_NODE_INDEX",
		"CIRCLE_NODE_TOTAL",
		"CIRCLE_PR_NUMBER",
		"CIRCLE_PR_REPONAME",
		"CIRCLE_PR_USERNAME",
		"CIRCLE_PROJECT_REPONAME",
		"CIRCLE_PROJECT_USERNAME",
		"CIRCLE_REPOSITORY_URL",
		"CIRCLE_SHA1",
		"CIRCLE_TAG",
		"CIRCLE_USERNAME",
		"CIRCLE_WORKFLOW_ID",
	}
	arr := make([]string, len(envs))
	for i, e := range envs {
		arr[i] = strings.ToLower(e) + ":" + os.Getenv(e)
	}
	return arr
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
