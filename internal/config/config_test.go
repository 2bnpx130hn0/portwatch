package config

import (
	"os"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTemp(t, "rules: []\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != DefaultInterval {
		t.Errorf("interval: got %v, want %v", cfg.Interval, DefaultInterval)
	}
	if cfg.StateFile != DefaultStateFile {
		t.Errorf("state_file: got %q, want %q", cfg.StateFile, DefaultStateFile)
	}
}

func TestLoad_ExplicitValues(t *testing.T) {
	path := writeTemp(t, `
interval: 10s
state_file: /tmp/state.json
rules:
  - port: 22
    protocol: tcp
    action: allow
    comment: ssh
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("interval: got %v, want 10s", cfg.Interval)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Port != 22 {
		t.Errorf("rules not loaded correctly: %+v", cfg.Rules)
	}
}

func TestLoad_RuleDefaults(t *testing.T) {
	path := writeTemp(t, "rules:\n  - port: 80\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Rules[0].Protocol != "tcp" {
		t.Errorf("protocol default: got %q, want tcp", cfg.Rules[0].Protocol)
	}
	if cfg.Rules[0].Action != "alert" {
		t.Errorf("action default: got %q, want alert", cfg.Rules[0].Action)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	path := writeTemp(t, "rules:\n  - port: 99999\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid port, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
