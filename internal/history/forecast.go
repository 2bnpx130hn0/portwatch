package history

import (
	"math"
	"sort"
	"time"
)

// ForecastResult holds a predicted activity level for a future bucket.
type ForecastResult struct {
	Bucket    time.Time
	Predicted float64
	Protocol  string
	Action    string
}

// ForecastOptions configures the Forecast function.
type ForecastOptions struct {
	Protocol  string
	Action    string
	BucketSize time.Duration
	Steps      int   // number of future buckets to predict
	Since      time.Time
}

// Forecast uses a simple linear regression over historical bucketed counts
// to predict future activity levels.
func Forecast(entries []Entry, opts ForecastOptions) []ForecastResult {
	if opts.BucketSize <= 0 {
		opts.BucketSize = time.Hour
	}
	if opts.Steps <= 0 {
		opts.Steps = 3
	}

	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		filtered = append(filtered, e)
	}

	if len(filtered) == 0 {
		return nil
	}

	buckets := map[int64]float64{}
	var minT int64 = math.MaxInt64
	for _, e := range filtered {
		key := e.Timestamp.Truncate(opts.BucketSize).Unix()
		buckets[key]++
		if key < minT {
			minT = key
		}
	}

	keys := make([]int64, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	n := float64(len(keys))
	var sumX, sumY, sumXY, sumX2 float64
	for i, k := range keys {
		x := float64(i)
		y := buckets[k]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	var slope, intercept float64
	if denom != 0 {
		slope = (n*sumXY - sumX*sumY) / denom
		intercept = (sumY - slope*sumX) / n
	}

	lastKey := keys[len(keys)-1]
	lastTime := time.Unix(lastKey, 0)
	results := make([]ForecastResult, 0, opts.Steps)
	for s := 1; s <= opts.Steps; s++ {
		x := n + float64(s) - 1
		predicted := math.Max(0, intercept+slope*x)
		results = append(results, ForecastResult{
			Bucket:    lastTime.Add(time.Duration(s) * opts.BucketSize),
			Predicted: math.Round(predicted*100) / 100,
			Protocol:  opts.Protocol,
			Action:    opts.Action,
		})
	}
	return results
}
