package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ddog "github.com/suzuki-shunsuke/dd-time/pkg/datadog"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
)

func TestController_validateParams(t *testing.T) {
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
	ctrl := &Controller{}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := ctrl.validateParams(d.params)
			if d.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestController_getDDOutput(t *testing.T) {
	data := []struct {
		title    string
		output   string
		appended bool
	}{
		{
			title: "default is os.Stderr",
		},
		{
			title:  "/dev/stderr",
			output: "/dev/stderr",
		},
		{
			title:  "/dev/null",
			output: "/dev/null",
		},
		{
			title:  "/dev/stdout",
			output: "/dev/stdout",
		},
	}
	ctrl := &Controller{}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			_, closeFn := ctrl.getDDOutput(d.output, d.appended)
			if closeFn != nil {
				closeFn()
			}
		})
	}
}

type mockExecutor struct {
	code int
}

func (exc *mockExecutor) Run(ctx context.Context, arg string, args ...string) error {
	if exc.code == 0 {
		return nil
	}
	return ecerror.Wrap(errors.New("command is failure"), exc.code)
}

type mockDataDog struct {
	err error
}

func (dd *mockDataDog) Send(*ddog.Params) error {
	return dd.err
}

func TestController_core(t *testing.T) {
	data := []struct {
		title  string
		params Params
		code   int
		ctrl   *Controller
	}{
		{
			title: "invalid parameters",
			code:  1,
			ctrl: &Controller{
				Exec:    &mockExecutor{},
				DataDog: &mockDataDog{},
				Now:     time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
		{
			title: "succeed",
			params: Params{
				Args: []string{"echo", "hello"},
			},
			ctrl: &Controller{
				Exec:    &mockExecutor{},
				DataDog: &mockDataDog{},
				Now:     time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
		{
			title: "succeed even if it is failed to send a metrics to DataDog",
			params: Params{
				Args: []string{"echo", "hello"},
			},
			ctrl: &Controller{
				Exec: &mockExecutor{},
				DataDog: &mockDataDog{
					err: errors.New("internal server error"),
				},
				Now: time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
		{
			title: "command is failure",
			params: Params{
				Args: []string{"echo", "hello"},
			},
			code: 3,
			ctrl: &Controller{
				Exec: &mockExecutor{
					code: 3,
				},
				DataDog: &mockDataDog{},
				Now:     time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
	}
	ctx := context.Background()
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			_, code := d.ctrl.core(ctx, d.params)
			require.Equal(t, d.code, code)
		})
	}
}

func TestController_Main(t *testing.T) {
	data := []struct {
		title  string
		params Params
		code   int
		ctrl   *Controller
	}{
		{
			title: "invalid parameters",
			code:  1,
			ctrl: &Controller{
				Exec:    &mockExecutor{},
				DataDog: &mockDataDog{},
				Now:     time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
		{
			title: "succeed",
			params: Params{
				Args: []string{"echo", "hello"},
			},
			ctrl: &Controller{
				Exec:    &mockExecutor{},
				DataDog: &mockDataDog{},
				Now:     time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
		{
			title: "succeed even if it is failed to send a metrics to DataDog",
			params: Params{
				Args: []string{"echo", "hello"},
			},
			ctrl: &Controller{
				Exec: &mockExecutor{},
				DataDog: &mockDataDog{
					err: errors.New("internal server error"),
				},
				Now: time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
		{
			title: "command is failure",
			params: Params{
				Args: []string{"echo", "hello"},
			},
			code: 3,
			ctrl: &Controller{
				Exec: &mockExecutor{
					code: 3,
				},
				DataDog: &mockDataDog{},
				Now:     time.Now,
				Since: func(t time.Time) time.Duration {
					return 5 * time.Second
				},
			},
		},
	}
	ctx := context.Background()
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			code := d.ctrl.Main(ctx, d.params)
			require.Equal(t, d.code, code)
		})
	}
}

func TestNew(t *testing.T) {
	require.NotNil(t, New(""))
}
