package history

import (
	"sort"
	"time"
)

// ChainLink represents a single step in a port activity chain.
type ChainLink struct {
	Entry    Entry
	GapSince time.Duration // time since previous link in chain
}

// Chain holds an ordered sequence of related history entries.
type Chain struct {
	Key   string // e.g. "80/tcp"
	Links []ChainLink
}

// ChainOptions controls how chains are built.
type ChainOptions struct {
	Protocol  string
	Port      int
	Action    string
	Since     time.Time
	MaxGap    time.Duration // if > 0, break chain when gap exceeds this
}

// BuildChain constructs an ordered activity chain for a specific port/protocol.
func BuildChain(entries []Entry, opts ChainOptions) Chain {
	var filtered []Entry
	for _, e := range entries {
		if opts.Protocol != "" && !equalFold(e.Protocol, opts.Protocol) {
			continue
		}
		if opts.Port != 0 && e.Port != opts.Port {
			continue
		}
		if opts.Action != "" && !equalFold(e.Action, opts.Action) {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		filtered = append(filtered, e)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	key := ""
	if len(filtered) > 0 {
		key = itoa(filtered[0].Port) + "/" + filtered[0].Protocol
	}

	chain := Chain{Key: key}
	for i, e := range filtered {
		var gap time.Duration
		if i > 0 {
			gap = e.Timestamp.Sub(filtered[i-1].Timestamp)
			if opts.MaxGap > 0 && gap > opts.MaxGap {
				break
			}
		}
		chain.Links = append(chain.Links, ChainLink{Entry: e, GapSince: gap})
	}
	return chain
}
