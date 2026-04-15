package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/watcher"
)

var version = "dev"

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	sc := scanner.New(cfg.ScanTimeout)
	rl := rules.New(cfg.Rules)
	al := alert.New(cfg.LogLevel)
	st, err := state.New(cfg.StateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising state: %v\n", err)
		os.Exit(1)
	}

	w := watcher.New(sc, rl, al, st, cfg.Interval)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Printf("portwatch %s starting (interval: %s, state: %s)\n", version, cfg.Interval, cfg.StateFile)

	if err := w.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "watcher exited with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("portwatch stopped")
}
