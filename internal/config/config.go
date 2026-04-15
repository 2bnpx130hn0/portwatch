package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	Interval time.Duration `yaml:"interval"`
	StateFile string        `yaml:"state_file"`
	Rules     []RuleConfig  `yaml:"rules"`
}

// RuleConfig represents a single port rule entry in the config file.
type RuleConfig struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Action   string `yaml:"action"`
	Comment  string `yaml:"comment"`
}

// defaults applied when fields are omitted.
const (
	DefaultInterval  = 30 * time.Second
	DefaultStateFile = "/var/lib/portwatch/state.json"
)

// Load reads a YAML config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	applyDefaults(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultInterval
	}
	if cfg.StateFile == "" {
		cfg.StateFile = DefaultStateFile
	}
	for i := range cfg.Rules {
		if cfg.Rules[i].Protocol == "" {
			cfg.Rules[i].Protocol = "tcp"
		}
		if cfg.Rules[i].Action == "" {
			cfg.Rules[i].Action = "alert"
		}
	}
}

func validate(cfg *Config) error {
	for _, r := range cfg.Rules {
		if r.Port < 1 || r.Port > 65535 {
			return fmt.Errorf("config: invalid port %d in rules", r.Port)
		}
	}
	return nil
}
