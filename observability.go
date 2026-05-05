package lago

import (
	"log"
	"sort"
	"strconv"
)

// DefaultRateLimitThresholds is the default set of usage thresholds at which
// LoggingRateLimitObserver emits a warning.
var DefaultRateLimitThresholds = []float64{0.80, 0.90, 0.95}

// LoggingRateLimitObserver is a ready-to-use OnRateLimitInfo callback that
// logs a warning each time rate limit usage crosses one of the configured
// thresholds.
//
// Example:
//
//	client := lago.New().SetApiKey("...")
//	client.RetryPolicy.OnRateLimitInfo = lago.NewLoggingRateLimitObserver(nil, nil)
type LoggingRateLimitObserver struct {
	thresholds []float64 // sorted descending
	logger     *log.Logger
}

// NewLoggingRateLimitObserver returns a LoggingRateLimitObserver as a
// RateLimitInfoCallback.
//
// Pass nil for thresholds to use DefaultRateLimitThresholds (80/90/95%).
// Pass nil for logger to use the standard logger.
func NewLoggingRateLimitObserver(thresholds []float64, logger *log.Logger) RateLimitInfoCallback {
	if len(thresholds) == 0 {
		thresholds = DefaultRateLimitThresholds
	}
	if logger == nil {
		logger = log.Default()
	}

	sorted := make([]float64, len(thresholds))
	copy(sorted, thresholds)
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))

	o := &LoggingRateLimitObserver{thresholds: sorted, logger: logger}
	return o.observe
}

func (o *LoggingRateLimitObserver) observe(info *RateLimitInfo) {
	pct, ok := info.UsagePct()
	if !ok {
		return
	}
	for _, t := range o.thresholds {
		if pct >= t {
			o.logger.Printf(
				"lago: rate limit at %.0f%% (limit=%s, remaining=%s, reset=%ss, %s %s)",
				pct*100,
				intPtrString(info.Limit),
				intPtrString(info.Remaining),
				intPtrString(info.Reset),
				info.Method, info.URL,
			)
			return
		}
	}
}

func intPtrString(v *int) string {
	if v == nil {
		return "<nil>"
	}
	return strconv.Itoa(*v)
}
