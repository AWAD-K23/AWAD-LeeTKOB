package main

import (
	"testing"

	"leetkob/backend/internal/semantic"
)

func TestModelConfiguration(t *testing.T) {
	t.Setenv("MODEL_NAME", semantic.SupportedModel)
	t.Setenv("MODEL_API_KEY", "fake-key")
	t.Setenv("MODEL_LAB_BUDGET_MICRO_USD", "1000000")

	for _, mode := range []string{"", "1"} {
		t.Setenv("MODEL_OFF", mode)
		config, err := readModelConfig()
		if err != nil || config.Enabled {
			t.Fatal("model enabled without explicit opt-in", err)
		}
	}

	t.Setenv("MODEL_OFF", "0")
	config, err := readModelConfig()
	if err != nil || !config.Enabled || config.BudgetMicroUSD != 1000000 {
		t.Fatal("valid settings rejected", err)
	}

	t.Setenv("MODEL_LAB_BUDGET_MICRO_USD", "-1")
	if _, err := readModelConfig(); err == nil {
		t.Fatal("negative budget accepted")
	}

	t.Setenv("MODEL_OFF", "false")
	if _, err := readModelConfig(); err == nil {
		t.Fatal("ambiguous mode accepted")
	}
}
