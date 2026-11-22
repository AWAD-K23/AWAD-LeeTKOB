package checker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"leetkob/backend/internal/sandbox"
)

// LocalScanner runs libclang inside the already-isolated checker image.
// It analyzes source only; it never executes the submitted program.
type LocalScanner struct{ Path string }

func (s LocalScanner) Run(parent context.Context, source []byte, _, mode string) (sandbox.Result, error) {
	if mode != "analyze" || len(source) > 1024*1024 {
		return sandbox.Result{}, errors.New("invalid analysis request")
	}
	dir, err := os.MkdirTemp("", "leetkob-scan-")
	if err != nil {
		return sandbox.Result{}, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "main.cpp")
	if err := os.WriteFile(path, source, 0600); err != nil {
		return sandbox.Result{}, err
	}

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, s.Path, path)
	// The compiler process must not inherit model credentials or local settings.
	command.Env = []string{"PATH=/usr/local/bin:/usr/bin:/bin", "LANG=C.UTF-8", "HOME=/tmp"}
	command.Stderr = io.Discard
	output := &boundedOutput{limit: 256 * 1024, cancel: cancel}
	command.Stdout = output
	if err := command.Run(); err != nil || output.exceeded {
		return sandbox.Result{}, errors.New("syntax analysis failed or exceeded limits")
	}

	return sandbox.Result{Verdict: "Executed", Output: output.String()}, nil
}

type boundedOutput struct {
	bytes.Buffer
	limit    int
	cancel   context.CancelFunc
	exceeded bool
}

func (b *boundedOutput) Write(data []byte) (int, error) {
	if len(data) > b.limit-b.Len() {
		b.exceeded = true
		b.cancel()
		return len(data), nil
	}
	return b.Buffer.Write(data)
}
