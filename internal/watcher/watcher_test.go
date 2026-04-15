package watcher_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"log/slog"
	"os"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watcher"
)

func newWatcher(t *testing.T, interval time.Duration) *watcher.Watcher {
	t.Helper()
	sc := scanner.New(50 * time.Millisecond)
	re := rules.New(nil) // no rules — defaults to alert
	al := alert.New(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := watcher.Config{
		Interval:  interval,
		Protocols: []string{"tcp"},
		PortRange: [2]int{1, 1024},
		StateFile: filepath.Join(t.TempDir(), "state.json"),
	}
	return watcher.New(cfg, sc, re, al, logger)
}

func TestRun_CancelImmediately(t *testing.T) {
	w := newWatcher(t, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Run

	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestRun_TickOnce(t *testing.T) {
	w := newWatcher(t, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Run should return with context.DeadlineExceeded after at least one tick.
	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected deadline error")
	}
}

func TestRun_StateFileCreated(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")

	sc := scanner.New(50 * time.Millisecond)
	re := rules.New(nil)
	al := alert.New(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := watcher.Config{
		Interval:  30 * time.Millisecond,
		Protocols: []string{"tcp"},
		PortRange: [2]int{1, 100},
		StateFile: statePath,
	}
	w := watcher.New(cfg, sc, re, al, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("state file not created: %v", err)
	}
}
