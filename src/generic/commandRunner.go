package generic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Docker içinde olup olmadığımızı /.dockerenv dosyasına bakarak anlıyoruz
func isDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return !os.IsNotExist(err)
}

// NewCmd — dizin belirtmeden komut çalıştırır
// Docker'da: nsenter --target 1 --mount --uts --ipc --net --pid -- <name> <args...>
// Dev'de  : doğrudan exec.CommandContext
func NewCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	if isDocker() {
		nsArgs := buildNsenterArgs(name, args...)
		return exec.CommandContext(ctx, "nsenter", nsArgs...)
	}
	return exec.CommandContext(ctx, name, args...)
}

// NewCmdInDir — belirli bir dizinde komut çalıştırır
// Docker'da nsenter içinden sh -c "cd <dir> && <name> <args...>" olarak çalışır
// Dev'de   : cmd.Dir set edilerek normal çalışır
func NewCmdInDir(ctx context.Context, dir string, name string, args ...string) *exec.Cmd {
	if isDocker() {
		// nsenter ile shell üzerinden dizine girip komutu çalıştırıyoruz
		shellCmd := fmt.Sprintf("cd %s && %s %s",
			shellEscape(dir),
			shellEscape(name),
			joinEscaped(args),
		)
		nsArgs := buildNsenterArgs("sh", "-c", shellCmd)
		return exec.CommandContext(ctx, "nsenter", nsArgs...)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	return cmd
}

// ── helpers ──────────────────────────────────────────────────────────────────

func buildNsenterArgs(name string, extra ...string) []string {
	base := []string{
		"--target", "1",
		"--mount", "--uts", "--ipc", "--net", "--pid",
		"--", name,
	}
	return append(base, extra...)
}

func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func joinEscaped(args []string) string {
	escaped := make([]string, len(args))
	for i, a := range args {
		escaped[i] = shellEscape(a)
	}
	return strings.Join(escaped, " ")
}
