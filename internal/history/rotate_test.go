package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRotate_BelowThreshold(t *testing.T) {
	dir := t.TempDir()
	h := New(filepath.Join(dir, "history.json"))
	_ = h.Record(Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()})

	rotated, err := Rotate(h, RotateOptions{MaxSizeBytes: 1 << 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Fatal("expected no rotation")
	}
}

func TestRotate_ExceedsThreshold(t *testing.T) {
	dir := t.TempDir()
	h := New(filepath.Join(dir, "history.json"))
	_ = h.Record(Entry{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: time.Now()})

	rotated, err := Rotate(h, RotateOptions{MaxSizeBytes: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Fatal("expected rotation")
	}

	// original file should be gone
	if _, err := os.Stat(h.path); !os.IsNotExist(err) {
		t.Fatal("expected original file to be removed")
	}

	// a rotated file should exist in same dir
	matches, _ := filepath.Glob(filepath.Join(dir, "history.json.*"))
	if len(matches) == 0 {
		t.Fatal("expected rotated file to exist")
	}
}

func TestRotate_CustomDestDir(t *testing.T) {
	dir := t.TempDir()
	archiveDir := filepath.Join(dir, "archive")
	h := New(filepath.Join(dir, "history.json"))
	_ = h.Record(Entry{Port: 443, Protocol: "tcp", Action: "alert", Timestamp: time.Now()})

	_, err := Rotate(h, RotateOptions{MaxSizeBytes: 1, DestDir: archiveDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(archiveDir, "history.json.*"))
	if len(matches) == 0 {
		t.Fatal("expected rotated file in archive dir")
	}
}

func TestRotate_NoFile(t *testing.T) {
	dir := t.TempDir()
	h := New(filepath.Join(dir, "missing.json"))

	rotated, err := Rotate(h, RotateOptions{MaxSizeBytes: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Fatal("expected no rotation for missing file")
	}
}

func TestRotate_ZeroThreshold(t *testing.T) {
	dir := t.TempDir()
	h := New(filepath.Join(dir, "history.json"))
	_ = h.Record(Entry{Port: 22, Protocol: "tcp", Action: "warn", Timestamp: time.Now()})

	rotated, err := Rotate(h, RotateOptions{MaxSizeBytes: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Fatal("expected no rotation when threshold is zero")
	}
}
