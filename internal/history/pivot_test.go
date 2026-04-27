package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func basePivotEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
		{Port: 443, Protocol: "tcp", Action: "allow", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "allow", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "alert", Timestamp: now},
		{Port: 53, Protocol: "udp", Action: "alert", Timestamp: now},
	}
}

func TestPivot_RowsAndCols(t *testing.T) {
	entries := basePivotEntries()
	pt := Pivot(entries, PivotOptions{RowField: "protocol", ColField: "action"})

	if len(pt.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(pt.Rows))
	}
	if len(pt.Cols) != 2 {
		t.Fatalf("expected 2 cols, got %d", len(pt.Cols))
	}
}

func TestPivot_CellCounts(t *testing.T) {
	entries := basePivotEntries()
	pt := Pivot(entries, PivotOptions{RowField: "protocol", ColField: "action"})

	if pt.Cells["udp"]["alert"] != 2 {
		t.Errorf("expected udp/alert=2, got %d", pt.Cells["udp"]["alert"])
	}
	if pt.Cells["tcp"]["allow"] != 2 {
		t.Errorf("expected tcp/allow=2, got %d", pt.Cells["tcp"]["allow"])
	}
}

func TestPivot_FilterByAction(t *testing.T) {
	entries := basePivotEntries()
	pt := Pivot(entries, PivotOptions{RowField: "port", ColField: "action", Action: "alert"})

	// only alert entries should be counted
	for row, cols := range pt.Cells {
		for col, count := range cols {
			if col != "alert" {
				t.Errorf("row %s: unexpected col %s with count %d", row, col, count)
			}
		}
	}
}

func TestPivot_SinceFilter(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Port: 80, Protocol: "tcp", Action: "allow", Timestamp: now.Add(-2 * time.Hour)},
		{Port: 80, Protocol: "tcp", Action: "alert", Timestamp: now},
	}
	pt := Pivot(entries, PivotOptions{
		RowField: "port",
		ColField: "action",
		Since:    now.Add(-30 * time.Minute).Unix(),
	})

	if pt.Cells["80"]["allow"] != 0 {
		t.Errorf("expected old allow entry excluded")
	}
	if pt.Cells["80"]["alert"] != 1 {
		t.Errorf("expected recent alert entry included")
	}
}

func TestRenderPivot_TextContainsHeaders(t *testing.T) {
	entries := basePivotEntries()
	pt := Pivot(entries, PivotOptions{RowField: "protocol", ColField: "action"})
	var buf bytes.Buffer
	RenderPivot(pt, "text", &buf)
	out := buf.String()
	if !strings.Contains(out, "action") {
		t.Error("expected column header in text output")
	}
}

func TestRenderPivot_JSONContainsCells(t *testing.T) {
	entries := basePivotEntries()
	pt := Pivot(entries, PivotOptions{RowField: "protocol", ColField: "action"})
	var buf bytes.Buffer
	RenderPivot(pt, "json", &buf)
	out := buf.String()
	if !strings.Contains(out, "cells") {
		t.Error("expected 'cells' key in JSON output")
	}
}
