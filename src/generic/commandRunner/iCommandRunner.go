package genericCommandrunner

import "context"

type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}
