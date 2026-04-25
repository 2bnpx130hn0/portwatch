package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseChainEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)},
	}
}

func TestBuildChain_FiltersPort(t *testing.T) {
	entries := baseChainEntries()
	c := BuildChain(entries, ChainOptions{Port: 80, Protocol: "tcp"})
	if len(c.Links) != 3 {
		t.Fatalf("expected 3 links, got %d", len(c.Links))
	}
	if c.Key != "80/tcp" {
		t.Errorf("expected key 80/tcp, got %s", c.Key)
	}
}

func TestBuildChain_OrderedByTime(t *testing.T) {
	entries := baseChainEntries()
	c := BuildChain(entries, ChainOptions{Port: 80, Protocol: "tcp"})
	for i := 1; i < len(c.Links); i++ {
		if c.Links[i].Entry.Timestamp.Before(c.Links[i-1].Entry.Timestamp) {
			t.Errorf("links not ordered at index %d", i)
		}
	}
}

func TestBuildChain_GapCalculated(t *testing.T) {
	entries := baseChainEntries()
	c := BuildChain(entries, ChainOptions{Port: 80, Protocol: "tcp"})
	if c.Links[0].GapSince != 0 {
		t.Errorf("first link gap should be 0")
	}
	if c.Links[1].GapSince == 0 {
		t.Errorf("subsequent link should have non-zero gap")
	}
}

func TestBuildChain_MaxGapBreaks(t *testing.T) {
	entries := baseChainEntries()
	c := BuildChain(entries, ChainOptions{Port: 80, Protocol: "tcp", MaxGap: 30 * time.Minute})
	if len(c.Links) != 1 {
		t.Errorf("expected chain broken after 1 link, got %d", len(c.Links))
	}
}

func TestBuildChain_Empty(t *testing.T) {
	c := BuildChain([]Entry{}, ChainOptions{Port: 9999})
	if len(c.Links) != 0 {
		t.Errorf("expected empty chain")
	}
}

func TestRenderChain_Text(t *testing.T) {
	entries := baseChainEntries()
	c := BuildChain(entries, ChainOptions{Port: 80, Protocol: "tcp"})
	var buf bytes.Buffer
	RenderChain(c, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "80/tcp") {
		t.Errorf("expected 80/tcp in output, got: %s", out)
	}
}

func TestRenderChain_JSON(t *testing.T) {
	entries := baseChainEntries()
	c := BuildChain(entries, ChainOptions{Port: 80, Protocol: "tcp"})
	var buf bytes.Buffer
	RenderChain(c, "json", &buf)
	out := buf.String()
	if !strings.Contains(out, `"key"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestRenderChain_DefaultsToText(t *testing.T) {
	c := Chain{}
	var buf bytes.Buffer
	RenderChain(c, "", &buf)
	if !strings.Contains(buf.String(), "no chain entries") {
		t.Errorf("expected empty message")
	}
}
