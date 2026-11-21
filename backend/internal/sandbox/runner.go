// Package sandbox is the only part of the worker allowed to control Docker.
package sandbox

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Runner struct {
	Image string
	Wall  time.Duration
}
type Result struct {
	Verdict   string
	Output    string
	Container string
}

func containerArgs(name, image, mode string) []string {
	command := `clang++ -std=c++17 -O0 /work/main.cpp -o /work/program || exit 42; exec prlimit --cpu=2:2 -- /work/program < /work/input`
	if mode == "analyze" {
		command = `exec /usr/local/bin/scan /work/main.cpp`
	}
	// Source arrives via a tar stream, never a host bind mount. No environment
	// from the API or worker is forwarded to this container.
	command = "tar -x -C /work && " + command
	return []string{"create", "--interactive", "--name", name, "--network", "none", "--memory", "256m", "--memory-swap", "256m",
		"--pids-limit", "32", "--cpus", "1", "--ulimit", "cpu=15:15", "--user", "65534:65534",
		"--read-only", "--cap-drop", "ALL", "--security-opt", "no-new-privileges",
		"--tmpfs", "/work:rw,exec,nosuid,size=64m,uid=65534,gid=65534,mode=700",
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=32m,mode=1777", image, "sh", "-c", command}
}

// limitedBuffer prevents an untrusted program from filling worker memory.
type limitedBuffer struct {
	mu       sync.Mutex
	b        bytes.Buffer
	limit    int
	cancel   context.CancelFunc
	exceeded bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := len(p)
	remaining := b.limit - b.b.Len()
	if len(p) > remaining {
		b.b.Write(p[:remaining])
		b.exceeded = true
		b.cancel()
	} else {
		b.b.Write(p)
	}
	return n, nil
}

func (r Runner) Run(parent context.Context, source []byte, input, mode string) (Result, error) {
	image := r.Image
	if image == "" {
		image = "leetkob-sandbox:pa3"
	}
	wall := r.Wall
	if wall == 0 {
		wall = 20 * time.Second
	}
	name := "leetkob-" + uuid.NewString()
	ctx, cancel := context.WithTimeout(parent, wall)
	defer cancel()
	// Always use a fresh context for cleanup: a timed-out context cannot kill
	// the container. Removing it forcibly kills every descendant too.
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		_ = exec.CommandContext(cleanup, "docker", "rm", "-f", name).Run()
	}()
	if out, err := exec.CommandContext(ctx, "docker", containerArgs(name, image, mode)...).CombinedOutput(); err != nil {
		return Result{}, fmt.Errorf("create sandbox: %w: %s", err, out)
	}
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	for _, file := range []struct {
		name string
		data []byte
	}{{"main.cpp", source}, {"input", []byte(input)}} {
		if err := tw.WriteHeader(&tar.Header{Name: file.name, Mode: 0600, Uid: 65534, Gid: 65534, Size: int64(len(file.data))}); err != nil {
			return Result{}, err
		}
		if _, err := tw.Write(file.data); err != nil {
			return Result{}, err
		}
	}
	if err := tw.Close(); err != nil {
		return Result{}, err
	}
	output := &limitedBuffer{limit: 64 * 1024, cancel: cancel}
	cmd := exec.CommandContext(ctx, "docker", "start", "--attach", "--interactive", name)
	cmd.Stdin = &archive
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()
	result := Result{Verdict: "Executed", Output: output.b.String(), Container: name}
	if output.exceeded {
		result.Verdict = "Output Limit Exceeded"
		return result, nil
	}
	if ctx.Err() != nil {
		result.Verdict = "Time Limit Exceeded"
		return result, nil
	}
	state, err := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.State.ExitCode}} {{.State.OOMKilled}}", name).Output()
	if err != nil {
		return Result{}, err
	}
	fields := strings.Fields(string(state))
	if len(fields) != 2 {
		return Result{}, errors.New("invalid container state")
	}
	switch {
	case fields[1] == "true":
		result.Verdict = "Memory Limit Exceeded"
	case fields[0] == "42":
		result.Verdict = "Compilation Error"
	case fields[0] == "137" || fields[0] == "152":
		result.Verdict = "Time Limit Exceeded"
	case fields[0] != "0":
		result.Verdict = "Runtime Error"
	}
	return result, nil
}
