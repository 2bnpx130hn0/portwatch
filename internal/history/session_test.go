package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseSessionEntries() []Entry {
	now := time.Now().UTC().Truncate(time.Second)
	return []Entry{
		{Port: 80, Protocol: "TCP", Action: "allow", Timestamp: now},
		{Port: 80, Protocol: "TCP", Action: "allow", Timestamp: now.Add(1 * time.Minute)},
		{Port: 80, Protocol: "TCP", Action: "allow", Timestamp: now.Add(2 * time.Minute)},
		// gap > 5 min — new session
		{Port: 80, Protocol: "TCP", Action: "allow", Timestamp: now.Add(10 * time.Minute)},
		{Port: 443, Protocol: "TCP", Action: "alert", Timestamp: now.Add(30 * time.Second)},
	}
}

func TestBuildSessions_GroupsByGap(t *testing.T) {
	entries := baseSessionEntries()
	sessions := BuildSessions(entries, SessionOptions{})
	// port 80 should produce 2 sessions; port 443 one
	count80 := 0
	for _, s := range sessions {
		if s.Port == 80 {
			count80++
		}
	}
	if count80 != 2 {
		t.Errorf("expected 2 sessions for port 80, got %d", count80)
	}
}

func TestBuildSessions_CountIsCorrect(t *testing.T) {
	entries := baseSessionEntries()
	sessions := BuildSessions(entries, SessionOptions{})
	for _, s := range sessions {
		if s.Port == 80 && s.Count == 3 {
			return
		}
	}
	t.Error("expected a session for port 80 with count 3")
}

func TestBuildSessions_FilterByAction(t *testing.T) {
	entries := baseSessionEntries()
	sessions := BuildSessions(entries, SessionOptions{Action: "alert"})
	if len(sessions) != 1 || sessions[0].Port != 443 {
		t.Errorf("expected 1 alert session on port 443, got %+v", sessions)
	}
}

func TestBuildSessions_SinceFilter(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	entries := []Entry{
		{Port: 22, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-1 * time.Hour)},
		{Port: 22, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-30 * time.Minute)},
	}
	sessions := BuildSessions(entries, SessionOptions{Since: now.Add(-45 * time.Minute)})
	if len(sessions) != 1 {
		t.Errorf("expected 1 session after since filter, got %d", len(sessions))
	}
}

func TestRenderSessions_Text(t *testing.T) {
	now := time.Now().UTC()
	sessions := []Session{
		{Port: 8080, Protocol: "tcp", Action: "allow", Start: now, End: now.Add(2 * time.Minute), Duration: 2 * time.Minute, Count: 5},
	}
	var buf bytes.Buffer
	RenderSessions(sessions, "text", &buf)
	if !strings.Contains(buf.String(), "8080") {
		t.Error("expected port 8080 in text output")
	}
}

func TestRenderSessions_JSON(t *testing.T) {
	now := time.Now().UTC()
	sessions := []Session{
		{Port: 9090, Protocol: "udp", Action: "alert", Start: now, End: now, Duration: 0, Count: 1},
	}
	var buf bytes.Buffer
	RenderSessions(sessions, "json", &buf)
	if !strings.Contains(buf.String(), "9090") {
		t.Error("expected port 9090 in json output")
	}
}

func TestRenderSessions_Empty(t *testing.T) {
	var buf bytes.Buffer
	RenderSessions(nil, "text", &buf)
	if !strings.Contains(buf.String(), "no sessions") {
		t.Error("expected 'no sessions' message for empty input")
	}
}
