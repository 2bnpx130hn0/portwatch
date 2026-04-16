package watcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Config holds watcher configuration.
type Config struct {
	Interval  time.Duration
	Protocols []string
	PortRange [2]int // [min, max]
	StateFile string
}

// Watcher periodically scans ports and emits alerts on changes.
type Watcher struct {
	cfg     Config
	scanner *scanner.Scanner
	rules   *rules.Engine
	alerter *alert.Alerter
	store   *state.Store
	logger  *slog.Logger
}

// New creates a Watcher with the provided dependencies.
func New(cfg Config, sc *scanner.Scanner, re *rules.Engine, al *alert.Alerter, logger *slog.Logger) *Watcher {
	return &Watcher{
		cfg:     cfg,
		scanner: sc,
		rules:   re,
		alerter: al,
		store:   state.New(cfg.StateFile),
		logger:  logger,
	}
}

// Run starts the watch loop until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	if _, err := w.store.Load(); err != nil {
		w.logger.Warn("could not load previous state", "err", err)
	}

	// Run an immediate tick before waiting for the first interval.
	if err := w.tick(); err != nil {
		w.logger.Error("scan tick failed", "err", err)
	}

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(); err != nil {
				w.logger.Error("scan tick failed", "err", err)
			}
		}
	}
}

func (w *Watcher) tick() error {
	ports, err := w.scanner.OpenPorts(w.cfg.Protocols, w.cfg.PortRange[0], w.cfg.PortRange[1])
	if err != nil {
		return err
	}

	next := &state.Snapshot{Timestamp: time.Now().UTC(), Ports: ports}
	prev := w.store.Current()

	if prev != nil {
		added, removed := state.Diff(prev, next)
		w.processChanges(added, "opened")
		w.processChanges(removed, "closed")
	} else {
		w.logger.Info("initial snapshot captured", "ports", ports)
	}

	return w.store.Save(next)
}

func (w *Watcher) processChanges(changes map[string][]int, direction string) {
	for proto, ports := range changes {
		for _, port := range ports {
			event := rules.Event{Port: port, Protocol: proto, Direction: direction}
			action := w.rules.Evaluate(event)
			w.alerter.Notify(event, action)
		}
	}
}
