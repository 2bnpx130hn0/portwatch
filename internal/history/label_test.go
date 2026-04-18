package history

import (
	"testing"
)

func baseLabelEntries() []Entry {
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow"},
		{Port: 443, Protocol: "tcp", Action: "allow"},
		{Port: 22, Protocol: "tcp", Action: "alert"},
	}
}

func TestLabel_AddsLabel(t *testing.T) {
	entries := baseLabelEntries()
	out := Label(entries, 80, "tcp", "env", "prod")
	if out[0].Labels["env"] != "prod" {
		t.Fatalf("expected label env=prod on port 80")
	}
	if out[1].Labels != nil && out[1].Labels["env"] != "" {
		t.Fatalf("expected no label on port 443")
	}
}

func TestLabel_CaseInsensitiveProtocol(t *testing.T) {
	entries := baseLabelEntries()
	out := Label(entries, 22, "TCP", "tier", "infra")
	if out[2].Labels["tier"] != "infra" {
		t.Fatalf("expected label tier=infra on port 22")
	}
}

func TestLabel_PreservesExisting(t *testing.T) {
	entries := baseLabelEntries()
	out := Label(entries, 80, "tcp", "env", "prod")
	out = Label(out, 80, "tcp", "owner", "team-a")
	if out[0].Labels["env"] != "prod" {
		t.Fatalf("expected env label preserved")
	}
	if out[0].Labels["owner"] != "team-a" {
		t.Fatalf("expected owner label set")
	}
}

func TestRemoveLabel_RemovesKey(t *testing.T) {
	entries := baseLabelEntries()
	out := Label(entries, 80, "tcp", "env", "prod")
	out = RemoveLabel(out, 80, "tcp", "env")
	if v := out[0].Labels["env"]; v != "" {
		t.Fatalf("expected env label removed, got %q", v)
	}
}

func TestFilterByLabel_ReturnsMatching(t *testing.T) {
	entries := baseLabelEntries()
	out := Label(entries, 80, "tcp", "env", "prod")
	out = Label(out, 443, "tcp", "env", "prod")
	out = Label(out, 22, "tcp", "env", "staging")
	filtered := FilterByLabel(out, "env", "prod")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(filtered))
	}
}

func TestFilterByLabel_NoMatch(t *testing.T) {
	entries := baseLabelEntries()
	filtered := FilterByLabel(entries, "env", "prod")
	if len(filtered) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(filtered))
	}
}
