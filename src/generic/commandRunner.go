package generic

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Docker içinde olup olmadığımızı /.dockerenv dosyasına bakarak anlıyoruz
func isDocker() bool {
	if _, err := os.Stat("/.dockerenv"); !os.IsNotExist(err) {
		return true
	}
	// Alternatif kontrol: /proc/1/cgroup (Docker/Containerd için)
	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil && strings.Contains(string(data), "docker") {
		return true
	}
	return false
}

// NewCmd — dizin belirtmeden komut çalıştırır
// Docker'da: nsenter --target 1 --mount --uts --ipc --net --pid -- <name> <args...>
// Dev'de  : doğrudan exec.CommandContext
func NewCmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	if isDocker() {
		// bash -l (login shell) kullanarak hostun tüm profil dosyalarını (.profile, /etc/profile vb.) yüklemesini sağlıyoruz.
		// Bu sayede PATH dahil tüm ortam değişkenleri hosttaki gibi olur.
		fullCmd := fmt.Sprintf("%s %s", shellEscape(name), joinEscaped(args))
		nsArgs := buildNsenterArgs("bash", "-l", "-c", fullCmd)
		log.Printf("[DOCKER] Running nsenter %v", nsArgs)
		return exec.CommandContext(ctx, "nsenter", nsArgs...)
	}
	log.Printf("[LOCAL] Running %s %v", name, args)
	return exec.CommandContext(ctx, name, args...)
}

// NewCmdInDir — belirli bir dizinde komut çalıştırır
// Docker'da nsenter içinden sh -c "cd <dir> && <name> <args...>" olarak çalışır
// Dev'de   : cmd.Dir set edilerek normal çalışır
func NewCmdInDir(ctx context.Context, dir string, name string, args ...string) *exec.Cmd {
	if isDocker() {
		// nsenter ile login shell üzerinden dizine girip komutu çalıştırıyoruz
		fullCmd := fmt.Sprintf("cd %s && %s %s",
			shellEscape(dir),
			shellEscape(name),
			joinEscaped(args),
		)
		nsArgs := buildNsenterArgs("bash", "-l", "-c", fullCmd)
		log.Printf("[DOCKER] Running nsenter %v", nsArgs)
		return exec.CommandContext(ctx, "nsenter", nsArgs...)
	}
	log.Printf("[LOCAL] Running (Dir: %s) %s %v", dir, name, args)
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
