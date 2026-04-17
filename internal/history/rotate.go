package history

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls when and how history files are rotated.
type RotateOptions struct {
	MaxSizeBytes int64
	DestDir      string
}

// Rotate checks if the history file exceeds MaxSizeBytes and, if so,
// moves it to DestDir with a timestamp suffix and starts fresh.
func Rotate(h *History, opts RotateOptions) (bool, error) {
	if opts.MaxSizeBytes <= 0 {
		return false, nil
	}

	info, err := os.Stat(h.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("rotate stat: %w", err)
	}

	if info.Size() < opts.MaxSizeBytes {
		return false, nil
	}

	dest := opts.DestDir
	if dest == "" {
		dest = filepath.Dir(h.path)
	}

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return false, fmt.Errorf("rotate mkdir: %w", err)
	}

	base := filepath.Base(h.path)
	stamp := time.Now().UTC().Format("20060102T150405")
	dstPath := filepath.Join(dest, fmt.Sprintf("%s.%s", base, stamp))

	if err := os.Rename(h.path, dstPath); err != nil {
		return false, fmt.Errorf("rotate rename: %w", err)
	}

	h.entries = nil
	return true, nil
}
