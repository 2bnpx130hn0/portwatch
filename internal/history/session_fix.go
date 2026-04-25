package history

import "strings"

// strings import shim — BuildSessions uses strings.EqualFold and strings.ToLower
// which are referenced in session.go but the import must live in the same file
// or a sibling. This file exists solely to satisfy the build; the real logic
// is inlined in session.go via the import block below.
//
// NOTE: Go requires each file to declare its own imports. The actual import
// of "strings" is declared inside session.go's import block directly.
// This file is intentionally empty beyond the package declaration.

// sessionKey is a comparable key for grouping entries.
type sessionKey struct {
	port  int
	proto string
}

func normSessionProto(p string) string {
	return strings.ToLower(p)
}
