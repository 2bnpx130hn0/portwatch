package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var baseFingerprintEntries = func() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-5 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-4 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
}()

func TestFingerprint_GroupsByPortProtocol(t *testing.T) {
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestFingerprint_CountsActions(t *testing.T) {
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{})
	// port 80 should have allow:2 alert:1
	var r80 *FingerprintResult
	for i := range results {
		if results[i].Port == 80 {
			r80 = &results[i]
			break
		}
	}
	if r80 == nil {
		t.Fatal("port 80 not found")
	}
	if r80.Actions["allow"] != 2 {
		t.Errorf("expected allow:2, got %d", r80.Actions["allow"])
	}
	if r80.Actions["alert"] != 1 {
		t.Errorf("expected alert:1, got %d", r80.Actions["alert"])
	}
}

func TestFingerprint_FingerprintStringStable(t *testing.T) {
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{})
	for _, r := range results {
		if r.Fingerprint == "" {
			t.Errorf("port %d has empty fingerprint", r.Port)
		}
	}
	// running again should produce identical fingerprints
	results2 := Fingerprint(baseFingerprintEntries, FingerprintOptions{})
	for i := range results {
		if results[i].Fingerprint != results2[i].Fingerprint {
			t.Errorf("fingerprint not stable for port %d", results[i].Port)
		}
	}
}

func TestFingerprint_SinceFilter(t *testing.T) {
	now := time.Now()
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{
		Since: now.Add(-90 * time.Minute),
	})
	// only port 22 (last two entries) should remain
	if len(results) != 1 || results[0].Port != 22 {
		t.Errorf("expected only port 22, got %+v", results)
	}
}

func TestFingerprint_ActionFilter(t *testing.T) {
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{Action: "alert"})
	for _, r := range results {
		if _, ok := r.Actions["alert"]; !ok {
			t.Errorf("port %d has no alert actions but was included", r.Port)
		}
	}
}

func TestRenderFingerprint_Text(t *testing.T) {
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{})
	var buf bytes.Buffer
	RenderFingerprint(results, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Error("expected header in text output")
	}
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
}

func TestRenderFingerprint_JSON(t *testing.T) {
	results := Fingerprint(baseFingerprintEntries, FingerprintOptions{})
	var buf bytes.Buffer
	RenderFingerprint(results, "json", &buf)
	if !strings.Contains(buf.String(), "Fingerprint") {
		t.Error("expected JSON output with Fingerprint field")
	}
}

func TestRenderFingerprint_Empty(t *testing.T) {
	var buf bytes.Buffer
	RenderFingerprint([]FingerprintResult{}, "text", &buf)
	if !strings.Contains(buf.String(), "no fingerprint") {
		t.Error("expected empty message")
	}
}
