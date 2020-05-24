package execute

import (
	"context"
	"os/exec"

	"github.com/suzuki-shunsuke/dd-time/pkg/env"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/go-timeout/timeout"
)

type Executor struct {
	Env *env.Env
}

func New() *Executor {
	return &Executor{
		Env: env.New(),
	}
}

func (exc *Executor) Run(ctx context.Context, arg string, args ...string) error {
	cmd := exec.Command(arg, args...)
	cmd.Stdin = exc.Env.Stdin
	cmd.Stdout = exc.Env.Stdout
	cmd.Stderr = exc.Env.Stderr
	cmd.Env = exc.Env.Environ

	runner := timeout.NewRunner(0)

	if err := runner.Run(ctx, cmd); err != nil {
		return ecerror.Wrap(err, cmd.ProcessState.ExitCode())
	}
	return nil
}
