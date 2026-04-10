package lago

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"
)

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
}

// DefaultRetryPolicy returns a RetryPolicy with sensible defaults.
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts:       3,
		EnableRetry:       true,
		InitialBackoff:    1,
		BackoffMultiplier: 2.0,
	}
}

// waitDuration calculates how long to wait before retrying.
// It prefers the x-ratelimit-reset header value if available,
// otherwise falls back to exponential backoff.
func (rp *RetryPolicy) waitDuration(resp *http.Response, attempt int) time.Duration {
	// Try to use x-ratelimit-reset header (seconds until window resets)
	if resp != nil {
		if resetStr := resp.Header.Get("x-ratelimit-reset"); resetStr != "" {
			if resetSeconds, err := strconv.Atoi(resetStr); err == nil && resetSeconds > 0 {
				return time.Duration(resetSeconds) * time.Second
			}
		}
	}

	// Fallback: exponential backoff
	backoffSeconds := float64(rp.InitialBackoff) * math.Pow(rp.BackoffMultiplier, float64(attempt))
	return time.Duration(backoffSeconds * float64(time.Second))
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
