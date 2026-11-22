package main

import (
	"path/filepath"
	"testing"

	"leetkob/backend/internal/checker"
)

func TestOfflineNeverConfiguresProvider(t *testing.T) {
	t.Setenv("MODEL_NAME", "invalid-provider")
	t.Setenv("MODEL_API_KEY", "fake-secret")
	for _, mode := range []string{"", "1"} {
		t.Setenv("MODEL_OFF", mode)
		engine := checker.Checker{}
		closeLog, err := configureModel(&engine, filepath.Join(t.TempDir(), "unused.jsonl"))
		if err != nil || engine.ModelFor != nil {
			t.Fatal("offline mode configured a provider", err)
		}
		closeLog()
	}
}

func TestInvalidModeFailsBeforeProvider(t *testing.T) {
	t.Setenv("MODEL_OFF", "false")
	if _, err := configureModel(&checker.Checker{}, ""); err == nil {
		t.Fatal("ambiguous mode accepted")
	}
}
