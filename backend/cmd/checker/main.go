// checker writes prediction JSONL from any dataset using the published layout.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"leetkob/backend/internal/checker"
	"leetkob/backend/internal/hardrules"
	"leetkob/backend/internal/localenv"
	"leetkob/backend/internal/sandbox"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	data := flag.String("data", "", "dataset folder")
	predictions := flag.String("predictions", "-", "output JSONL path, or - for stdout")
	scannerPath := flag.String("scanner", "docker", "scanner executable, or docker")
	callLog := flag.String("call-log", "", "new audit JSONL file, required for enabled model")
	flag.Parse()

	if *data == "" {
		return fmt.Errorf("--data is required")
	}
	if err := localenv.Load(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	var scanner hardrules.Scanner = sandbox.Runner{}
	if *scannerPath != "docker" {
		scanner = checker.LocalScanner{Path: *scannerPath}
	}

	engine := checker.Checker{Scanner: scanner}
	closeAudit, err := configureModel(&engine, *callLog)
	if err != nil {
		return err
	}
	defer closeAudit()

	var output io.Writer = os.Stdout
	if *predictions != "-" {
		file, err := os.OpenFile(*predictions, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		output = file
	}

	return engine.Run(ctx, *data, output)
}
