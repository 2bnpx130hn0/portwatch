package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var baseHotspotEntries = func() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-3 * time.Hour)},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)},
		{Port: 22, Protocol: "tcp", Action: "alert", Timestamp: now.Add(-10 * time.Minute)},
		{Port: 53, Protocol: "udp", Action: "allow", Timestamp: now.Add(-5 * time.Minute)},
	}
}()

func TestHotspot_RankedByCount(t *testing.T) {
	result := Hotspot(baseHotspotEntries, HotspotOptions{})
	if len(result) == 0 {
		t.Fatal("expected results")
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80 as hottest, got %d", result[0].Port)
	}
	if result[0].Count != 3 {
		t.Errorf("expected count 3, got %d", result[0].Count)
	}
}

func TestHotspot_TopN(t *testing.T) {
	result := Hotspot(baseHotspotEntries, HotspotOptions{TopN: 2})
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestHotspot_FilterByAction(t *testing.T) {
	result := Hotspot(baseHotspotEntries, HotspotOptions{Action: "alert"})
	for _, h := range result {
		for action := range h.Actions {
			if action != "alert" {
				t.Errorf("unexpected action %q in filtered results", action)
			}
		}
	}
	if len(result) == 0 {
		t.Error("expected at least one alert hotspot")
	}
}

func TestHotspot_FilterByProtocol(t *testing.T) {
	result := Hotspot(baseHotspotEntries, HotspotOptions{Protocol: "udp"})
	if len(result) != 1 {
		t.Fatalf("expected 1 udp hotspot, got %d", len(result))
	}
	if result[0].Port != 53 {
		t.Errorf("expected port 53, got %d", result[0].Port)
	}
}

func TestHotspot_SinceFilter(t *testing.T) {
	now := time.Now()
	result := Hotspot(baseHotspotEntries, HotspotOptions{Since: now.Add(-20 * time.Minute)})
	for _, h := range result {
		if h.LastSeen.Before(now.Add(-20 * time.Minute)) {
			t.Errorf("entry %d/%s older than since filter", h.Port, h.Protocol)
		}
	}
}

func TestRenderHotspot_Text(t *testing.T) {
	var buf bytes.Buffer
	RenderHotspot([]HotspotEntry{
		{Port: 80, Protocol: "tcp", Count: 5, LastSeen: time.Now()},
	}, "text", &buf)
	if !strings.Contains(buf.String(), "80") {
		t.Error("expected port 80 in text output")
	}
}

func TestRenderHotspot_JSON(t *testing.T) {
	var buf bytes.Buffer
	RenderHotspot([]HotspotEntry{
		{Port: 443, Protocol: "tcp", Count: 2, LastSeen: time.Now()},
	}, "json", &buf)
	if !strings.Contains(buf.String(), "443") {
		t.Error("expected port 443 in json output")
	}
}

func TestRenderHotspot_Empty(t *testing.T) {
	var buf bytes.Buffer
	RenderHotspot([]HotspotEntry{}, "text", &buf)
	if !strings.Contains(buf.String(), "no hotspots") {
		t.Error("expected empty message")
	}
}
