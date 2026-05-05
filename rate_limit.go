package lago

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RateLimitInfo holds parsed rate limit headers from a Lago API response.
//
// It is delivered to the OnRateLimitInfo callback after every successful
// request so callers can build observability around the rate limit (warn at
// thresholds, emit metrics, etc.).
type RateLimitInfo struct {
	// Limit is the parsed x-ratelimit-limit header (nil when absent or unparseable).
	Limit *int
	// Remaining is the parsed x-ratelimit-remaining header (nil when absent or unparseable).
	Remaining *int
	// Reset is the parsed x-ratelimit-reset header in seconds (nil when absent or unparseable).
	Reset *int
	// Method is the HTTP method of the call (GET, POST, ...).
	Method string
	// URL is the request URL.
	URL string
}

// UsagePct returns the fraction of the rate limit currently used in [0.0, 1.0].
//
// It returns (0, false) when the headers aren't usable (missing limit, zero
// limit, missing remaining).
func (i *RateLimitInfo) UsagePct() (float64, bool) {
	if i == nil || i.Limit == nil || i.Remaining == nil || *i.Limit <= 0 {
		return 0, false
	}
	return 1.0 - float64(*i.Remaining)/float64(*i.Limit), true
}

// RateLimitInfoCallback is invoked after every successful response with parsed
// rate limit headers.
type RateLimitInfoCallback func(info *RateLimitInfo)

// RetryPolicy defines the configuration for rate limit retry behavior.
type RetryPolicy struct {
	// MaxAttempts is the maximum number of attempts (including the initial request).
	// Default is 3 (initial + 2 retries).
	MaxAttempts int

	// EnableRetry controls whether automatic retries are enabled.
	// Default is true.
	EnableRetry bool

	// InitialBackoff is the initial backoff duration in seconds for exponential backoff.
	// Used when x-ratelimit-reset header is missing or invalid.
	// Default is 1 second.
	InitialBackoff int

	// BackoffMultiplier is the multiplier for exponential backoff.
	// Default is 2 (doubles the wait time on each retry).
	BackoffMultiplier float64

	// MaxRetryDelay is the maximum duration to wait before a retry, in seconds.
	// Applies to both header-based and exponential backoff delays.
	// Default is 20 seconds.
	MaxRetryDelay int

	// OnRateLimitInfo, when set, is invoked after every successful (non-429)
	// response with parsed rate limit headers. Use it to build observability
	// (warn at 80/90/95%, emit metrics, etc.). Panics from the callback are
	// recovered and logged so a buggy observer cannot break the request flow.
	OnRateLimitInfo RateLimitInfoCallback
}

// DefaultRetryPolicy returns a RetryPolicy with sensible defaults.
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts:       3,
		EnableRetry:       true,
		InitialBackoff:    1,
		BackoffMultiplier: 2.0,
		MaxRetryDelay:     20,
	}
}

// parseRateLimitInfo extracts x-ratelimit-* headers from a response into a
// RateLimitInfo. Returns nil when no rate limit headers are present (for
// example, on a self-hosted Lago instance that doesn't enforce limits).
func parseRateLimitInfo(resp *http.Response, method, url string) *RateLimitInfo {
	if resp == nil {
		return nil
	}

	info := &RateLimitInfo{Method: method, URL: url}
	hasAny := false

	if v := resp.Header.Get("x-ratelimit-limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			info.Limit = intPtr(n)
			hasAny = true
		}
	}
	if v := resp.Header.Get("x-ratelimit-remaining"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			info.Remaining = intPtr(n)
			hasAny = true
		}
	}
	if v := resp.Header.Get("x-ratelimit-reset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			info.Reset = intPtr(n)
			hasAny = true
		}
	}

	if !hasAny {
		return nil
	}
	return info
}

// emitRateLimitInfo invokes the configured OnRateLimitInfo callback if any.
// Panics from the callback are recovered and logged so the underlying request
// flow is never affected.
func (rp *RetryPolicy) emitRateLimitInfo(resp *http.Response, method, url string) {
	if rp == nil || rp.OnRateLimitInfo == nil {
		return
	}
	info := parseRateLimitInfo(resp, method, url)
	if info == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("lago: OnRateLimitInfo callback panicked: %v", r)
		}
	}()
	rp.OnRateLimitInfo(info)
}

// waitDuration calculates how long to wait before retrying.
// It prefers the x-ratelimit-reset header value if available,
// otherwise falls back to exponential backoff.
func (rp *RetryPolicy) waitDuration(resp *http.Response, attempt int) time.Duration {
	var duration time.Duration

	// Try to use x-ratelimit-reset header (seconds until window resets)
	if resp != nil {
		if resetStr := resp.Header.Get("x-ratelimit-reset"); resetStr != "" {
			if resetSeconds, err := strconv.Atoi(resetStr); err == nil && resetSeconds > 0 {
				duration = time.Duration(resetSeconds) * time.Second
			}
		}
	}

	// Fallback: exponential backoff
	if duration == 0 {
		backoffSeconds := float64(rp.InitialBackoff) * math.Pow(rp.BackoffMultiplier, float64(attempt))
		duration = time.Duration(backoffSeconds * float64(time.Second))
	}

	// Cap at MaxRetryDelay
	if rp.MaxRetryDelay > 0 {
		maxDelay := time.Duration(rp.MaxRetryDelay) * time.Second
		if duration > maxDelay {
			duration = maxDelay
		}
	}

	return duration
}

// WaitForRateLimit blocks the current goroutine until the rate limit window resets.
// It uses the x-ratelimit-reset value from a RateLimitError if available, otherwise
// uses exponential backoff based on the attempt number.
//
// The attempt parameter is 0-based and controls exponential backoff when the reset
// header is unavailable. This is useful for cases where the client chooses not to
// use automatic retries but still wants to wait before making the next request.
func WaitForRateLimit(ctx context.Context, rlErr *RateLimitError, policy *RetryPolicy, attempt int) error {
	if rlErr == nil || policy == nil || !policy.EnableRetry {
		return nil
	}

	var waitDuration time.Duration

	if rlErr.Reset != nil && *rlErr.Reset > 0 {
		// Use the Reset value from the error (seconds until window resets)
		waitDuration = time.Duration(*rlErr.Reset) * time.Second
	} else {
		// Fallback to exponential backoff
		backoffSeconds := float64(policy.InitialBackoff) * math.Pow(policy.BackoffMultiplier, float64(attempt))
		waitDuration = time.Duration(backoffSeconds * float64(time.Second))
	}

	// Cap at MaxRetryDelay
	if policy.MaxRetryDelay > 0 {
		maxDelay := time.Duration(policy.MaxRetryDelay) * time.Second
		if waitDuration > maxDelay {
			waitDuration = maxDelay
		}
	}

	// Use a timer with context to allow early cancellation
	timer := time.NewTimer(waitDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
