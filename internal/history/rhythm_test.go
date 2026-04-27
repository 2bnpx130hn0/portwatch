package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseRhythmEntries() []Entry {
	now := time.Now()
	var entries []Entry
	// Port 80/tcp — very regular, every 60s
	for i := 0; i < 6; i++ {
		entries = append(entries, Entry{
			Port:      80,
			Protocol:  "tcp",
			Action:    "allow",
			Timestamp: now.Add(time.Duration(i) * 60 * time.Second),
		})
	}
	// Port 443/tcp — irregular
	offsets := []time.Duration{0, 5 * time.Second, 90 * time.Second, 95 * time.Second, 300 * time.Second}
	for _, off := range offsets {
		entries = append(entries, Entry{
			Port:      443,
			Protocol:  "tcp",
			Action:    "alert",
			Timestamp: now.Add(off),
		})
	}
	return entries
}

func TestRhythm_DetectsRegular(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 3})
	var found *RhythmResult
	for i := range results {
		if results[i].Port == 80 {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected port 80 in results")
	}
	if !found.Regular {
		t.Errorf("expected port 80 to be regular, got stddev=%v avg=%v", found.PeriodStddev, found.PeriodAvg)
	}
}

func TestRhythm_DetectsIrregular(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 3})
	var found *RhythmResult
	for i := range results {
		if results[i].Port == 443 {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected port 443 in results")
	}
	if found.Regular {
		t.Errorf("expected port 443 to be irregular")
	}
}

func TestRhythm_MinOccurrencesFilters(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 10})
	if len(results) != 0 {
		t.Errorf("expected no results with high MinOccurrences, got %d", len(results))
	}
}

func TestRhythm_RegularOnly(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 3, RegularOnly: true})
	for _, r := range results {
		if !r.Regular {
			t.Errorf("RegularOnly=true but got irregular port %d", r.Port)
		}
	}
}

func TestRhythm_ActionFilter(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 3, Action: "alert"})
	for _, r := range results {
		if r.Port == 80 {
			t.Errorf("port 80 (allow) should be excluded when filtering by alert")
		}
	}
}

func TestRenderRhythm_Text(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 3})
	var buf bytes.Buffer
	RenderRhythm(results, "text", &buf)
	if !strings.Contains(buf.String(), "PORT") {
		t.Errorf("expected header in text output, got: %s", buf.String())
	}
}

func TestRenderRhythm_JSON(t *testing.T) {
	entries := baseRhythmEntries()
	results := Rhythm(entries, RhythmOptions{MinOccurrences: 3})
	var buf bytes.Buffer
	RenderRhythm(results, "json", &buf)
	if !strings.Contains(buf.String(), "period_avg_ms") {
		t.Errorf("expected JSON keys, got: %s", buf.String())
	}
}

func TestRenderRhythm_Empty(t *testing.T) {
	var buf bytes.Buffer
	RenderRhythm(nil, "text", &buf)
	if !strings.Contains(buf.String(), "no rhythm") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
