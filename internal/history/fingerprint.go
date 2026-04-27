package history

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FingerprintResult holds a unique behavioral fingerprint for a port/protocol pair.
type FingerprintResult struct {
	Port        int
	Protocol    string
	Fingerprint string
	Actions     map[string]int
	FirstSeen   time.Time
	LastSeen    time.Time
	EventCount  int
}

// FingerprintOptions controls fingerprint generation.
type FingerprintOptions struct {
	Since    time.Time
	Protocol string
	Action   string
}

// Fingerprint generates a stable behavioral fingerprint for each port/protocol
// pair based on the sequence and distribution of observed actions.
func Fingerprint(entries []Entry, opts FingerprintOptions) []FingerprintResult {
	type key struct {
		port  int
		proto string
	}

	type bucket struct {
		actions   map[string]int
		first     time.Time
		last      time.Time
		count     int
	}

	buckets := map[key]*bucket{}

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		proto := strings.ToLower(e.Protocol)
		if opts.Protocol != "" && proto != strings.ToLower(opts.Protocol) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		k := key{port: e.Port, proto: proto}
		if buckets[k] == nil {
			buckets[k] = &bucket{
				actions: map[string]int{},
				first:   e.Timestamp,
				last:    e.Timestamp,
			}
		}
		b := buckets[k]
		b.actions[strings.ToLower(e.Action)]++
		b.count++
		if e.Timestamp.Before(b.first) {
			b.first = e.Timestamp
		}
		if e.Timestamp.After(b.last) {
			b.last = e.Timestamp
		}
	}

	results := make([]FingerprintResult, 0, len(buckets))
	for k, b := range buckets {
		results = append(results, FingerprintResult{
			Port:        k.port,
			Protocol:    k.proto,
			Fingerprint: buildFingerprint(k.port, k.proto, b.actions),
			Actions:     b.actions,
			FirstSeen:   b.first,
			LastSeen:    b.last,
			EventCount:  b.count,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Port != results[j].Port {
			return results[i].Port < results[j].Port
		}
		return results[i].Protocol < results[j].Protocol
	})
	return results
}

func buildFingerprint(port int, proto string, actions map[string]int) string {
	keys := make([]string, 0, len(actions))
	for k := range actions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s:%d", k, actions[k]))
	}
	return fmt.Sprintf("%d/%s[%s]", port, proto, strings.Join(parts, ","))
}
