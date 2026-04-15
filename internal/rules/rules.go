package rules

import (
	"fmt"
	"strings"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAlert  Action = "alert"
	ActionIgnore Action = "ignore"
)

// Rule represents a single port monitoring rule.
type Rule struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"` // tcp or udp
	Action   Action `yaml:"action"`
	Comment  string `yaml:"comment"`
}

// RuleSet holds a collection of rules and provides evaluation logic.
type RuleSet struct {
	Rules []Rule
}

// New creates a new RuleSet from the provided rules.
func New(rules []Rule) *RuleSet {
	return &RuleSet{Rules: rules}
}

// Evaluate checks a port/protocol pair against the rule set.
// It returns the matching Action and whether a rule was found.
// If no rule matches, it defaults to ActionAlert.
func (rs *RuleSet) Evaluate(port int, protocol string) (Action, bool) {
	protocol = strings.ToLower(protocol)
	for _, r := range rs.Rules {
		if r.Port == port && strings.ToLower(r.Protocol) == protocol {
			return r.Action, true
		}
	}
	return ActionAlert, false
}

// Validate checks that all rules in the set are well-formed.
func (rs *RuleSet) Validate() error {
	for i, r := range rs.Rules {
		if r.Port < 1 || r.Port > 65535 {
			return fmt.Errorf("rule %d: port %d is out of valid range (1-65535)", i, r.Port)
		}
		proto := strings.ToLower(r.Protocol)
		if proto != "tcp" && proto != "udp" {
			return fmt.Errorf("rule %d: protocol %q is invalid, must be tcp or udp", i, r.Protocol)
		}
		if r.Action != ActionAlert && r.Action != ActionIgnore {
			return fmt.Errorf("rule %d: action %q is invalid, must be alert or ignore", i, r.Action)
		}
	}
	return nil
}
