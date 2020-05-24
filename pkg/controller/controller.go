package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	ddog "github.com/suzuki-shunsuke/dd-time/pkg/datadog"
	"github.com/suzuki-shunsuke/dd-time/pkg/execute"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
)

type Controller struct {
	Exec    Executor
	DataDog DataDogClient
	Now     func() time.Time
	Since   func(time.Time) time.Duration
}

type Executor interface {
	Run(ctx context.Context, arg string, args ...string) error
}

type DataDogClient interface {
	Send(*ddog.Params) error
}

func New(apiKey string) *Controller {
	return &Controller{
		Exec:    execute.New(),
		Now:     time.Now,
		Since:   time.Since,
		DataDog: ddog.New(apiKey),
	}
}

func (ctrl *Controller) Main(ctx context.Context, params Params) int {
	ddOutput, closeOutput := ctrl.getDDOutput(params.Output, params.Append)
	if closeOutput != nil {
		defer closeOutput()
	}
	msg, code := ctrl.core(ctx, params)
	if msg != "" {
		fmt.Fprintln(ddOutput, msg)
	}
	return code
}

func (ctrl *Controller) core(ctx context.Context, params Params) (string, int) {
	if err := ctrl.validateParams(params); err != nil {
		return err.Error(), 1
	}

	startT := ctrl.Now()
	if err := ctrl.Exec.Run(ctx, params.Args[0], params.Args[1:]...); err != nil {
		return err.Error(), ecerror.GetExitCode(err)
	}
	duration := ctrl.Since(startT).Seconds()

	if err := ctrl.DataDog.Send(&ddog.Params{
		MetricName: params.MetricName,
		MetricHost: params.MetricHost,
		Tags:       append(params.Tags, ddog.GetTags()...),
		Duration:   duration,
		Now:        float64(startT.Unix()),
	}); err != nil {
		return "send a time series metrics to DataDog: " + err.Error(), 0
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
)

func (ctrl *Controller) validateParams(params Params) error {
	if len(params.Args) == 0 {
		return errors.New("executed command isn't passed to dd-time")
	}
	return nil
}

func (ctrl *Controller) getDDOutput(output string, appended bool) (io.Writer, func() error) {
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
