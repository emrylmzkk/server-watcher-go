package genericCommandrunner

import (
	"context"
	"os/exec"
)

type NsenterCommandRunner struct{}

func (r *NsenterCommandRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	// nsenter + gerçek komut
	allArgs := append([]string{
		"--target", "1",
		"--mount",
		"--uts",
		"--ipc",
		"--net",
		"--pid",
		"--",
		name,
	}, args...)

	cmd := exec.CommandContext(ctx, "nsenter", allArgs...)
	return cmd.Output()
}
