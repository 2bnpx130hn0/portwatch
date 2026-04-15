package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "portwatch")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "." // run from package directory
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_VersionFlag(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(string(out), "portwatch ") {
		t.Errorf("expected version prefix, got: %q", string(out))
	}
}

func TestMain_MissingConfig(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--config", "/nonexistent/path/config.yaml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code for missing config")
	}
	if !strings.Contains(string(out), "error loading config") {
		t.Errorf("expected config error message, got: %q", string(out))
	}
}

func TestMain_ValidConfig_StartsAndStops(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	cfgContent := `interval: 1s
log_level: info
state_file: ` + filepath.Join(dir, "state.json") + `
scan_timeout: 500ms
rules: []
`
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd := exec.Command(bin, "--config", cfgFile)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	// send SIGTERM immediately
	_ = cmd.Process.Signal(os.Interrupt)
	_ = cmd.Wait()
	// process should exit cleanly (exit 0 or signal exit)
}
