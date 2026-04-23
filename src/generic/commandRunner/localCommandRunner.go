package genericCommandrunner

import (
	"context"
	"os/exec"
)

type LocalCommandRunner struct{}

func (r *LocalCommandRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Output()
}
