package rules

import (
	"testing"
)

func TestEvaluate_MatchingRule(t *testing.T) {
	rs := New([]Rule{
		{Port: 8080, Protocol: "tcp", Action: ActionIgnore, Comment: "dev server"},
		{Port: 22, Protocol: "tcp", Action: ActionAlert, Comment: "ssh"},
	})

	action, found := rs.Evaluate(8080, "tcp")
	if !found {
		t.Fatal("expected rule to be found for port 8080/tcp")
	}
	if action != ActionIgnore {
		t.Errorf("expected action %q, got %q", ActionIgnore, action)
	}
}

func TestEvaluate_NoMatchDefaultsToAlert(t *testing.T) {
	rs := New([]Rule{
		{Port: 8080, Protocol: "tcp", Action: ActionIgnore},
	})

	action, found := rs.Evaluate(9090, "tcp")
	if found {
		t.Fatal("expected no rule to be found for port 9090/tcp")
	}
	if action != ActionAlert {
		t.Errorf("expected default action %q, got %q", ActionAlert, action)
	}
}

func TestEvaluate_CaseInsensitiveProtocol(t *testing.T) {
	rs := New([]Rule{
		{Port: 53, Protocol: "UDP", Action: ActionIgnore},
	})

	action, found := rs.Evaluate(53, "udp")
	if !found {
		t.Fatal("expected rule to match case-insensitively")
	}
	if action != ActionIgnore {
		t.Errorf("expected action %q, got %q", ActionIgnore, action)
	}
}

func TestEvaluate_FirstMatchingRuleWins(t *testing.T) {
	rs := New([]Rule{
		{Port: 80, Protocol: "tcp", Action: ActionIgnore},
		{Port: 80, Protocol: "tcp", Action: ActionAlert},
	})

	action, found := rs.Evaluate(80, "tcp")
	if !found {
		t.Fatal("expected rule to be found for port 80/tcp")
	}
	if action != ActionIgnore {
		t.Errorf("expected first matching rule to win: action %q, got %q", ActionIgnore, action)
	}
}

func TestValidate_ValidRules(t *testing.T) {
	rs := New([]Rule{
		{Port: 443, Protocol: "tcp", Action: ActionAlert},
		{Port: 53, Protocol: "udp", Action: ActionIgnore},
	})
	if err := rs.Validate(); err != nil {
		t.Errorf("expected no validation error, got: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	rs := New([]Rule{
		{Port: 0, Protocol: "tcp", Action: ActionAlert},
	})
	if err := rs.Validate(); err == nil {
		t.Error("expected validation error for port 0")
	}
}

func TestValidate_InvalidProtocol(t *testing.T) {
	rs := New([]Rule{
		{Port: 80, Protocol: "icmp", Action: ActionAlert},
	})
	if err := rs.Validate(); err == nil {
		t.Error("expected validation error for protocol icmp")
	}
}

func TestValidate_InvalidAction(t *testing.T) {
	rs := New([]Rule{
		{Port: 80, Protocol: "tcp", Action: "notify"},
	})
	if err := rs.Validate(); err == nil {
		t.Error("expected validation error for unknown action")
	}
}
