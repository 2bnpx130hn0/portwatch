package history

import (
	"testing"
	"time"
)

func baseForecastEntries() []Entry {
	now := time.Now().Truncate(time.Hour)
	entries := []Entry{}
	// Generate increasing counts per hour to give a positive slope
	for h := 0; h < 5; h++ {
		count := h + 1
		for i := 0; i < count; i++ {
			entries = append(entries, Entry{
				Port:      8080,
				Protocol:  "tcp",
				Action:    "alert",
				Timestamp: now.Add(-time.Duration(5-h) * time.Hour),
			})
		}
	}
	return entries
}

func TestForecast_ReturnsSteps(t *testing.T) {
	entries := baseForecastEntries()
	results := Forecast(entries, ForecastOptions{Steps: 3, BucketSize: time.Hour})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestForecast_PositiveSlope(t *testing.T) {
	entries := baseForecastEntries()
	results := Forecast(entries, ForecastOptions{Steps: 2, BucketSize: time.Hour})
	if len(results) < 2 {
		t.Fatal("expected at least 2 results")
	}
	if results[1].Predicted < results[0].Predicted {
		t.Errorf("expected increasing prediction, got %v then %v", results[0].Predicted, results[1].Predicted)
	}
}

func TestForecast_FilterByProtocol(t *testing.T) {
	entries := baseForecastEntries()
	entries = append(entries, Entry{
		Port: 53, Protocol: "udp", Action: "allow",
		Timestamp: time.Now().Add(-time.Minute),
	})
	results := Forecast(entries, ForecastOptions{
		Protocol:   "udp",
		Steps:      2,
		BucketSize: time.Hour,
	})
	for _, r := range results {
		if r.Protocol != "udp" {
			t.Errorf("expected udp, got %s", r.Protocol)
		}
	}
}

func TestForecast_EmptyEntries(t *testing.T) {
	results := Forecast(nil, ForecastOptions{Steps: 3, BucketSize: time.Hour})
	if len(results) != 0 {
		t.Errorf("expected empty results for nil entries, got %d", len(results))
	}
}

func TestForecast_DefaultBucketSize(t *testing.T) {
	entries := baseForecastEntries()
	// BucketSize=0 should default to 1 hour
	results := Forecast(entries, ForecastOptions{Steps: 1})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Predicted < 0 {
		t.Error("predicted value should not be negative")
	}
}

func TestForecast_SinceFilter(t *testing.T) {
	entries := baseForecastEntries()
	// Add old entries that should be excluded
	for i := 0; i < 100; i++ {
		entries = append(entries, Entry{
			Port: 9999, Protocol: "tcp", Action: "alert",
			Timestamp: time.Now().Add(-72 * time.Hour),
		})
	}
	since := time.Now().Add(-6 * time.Hour)
	results := Forecast(entries, ForecastOptions{
		Steps:      2,
		BucketSize: time.Hour,
		Since:      since,
	})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
