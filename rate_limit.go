package lago

import (
	"context"
	"math"
	"time"

	"github.com/go-resty/resty/v2"
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

// retryCount holds retry state for a specific request.
type retryCount struct {
	count int
}

// setupRateLimitRetry configures resty's built-in retry mechanism with a custom
// condition to detect rate limit (429) responses and a backoff strategy that
// respects the x-ratelimit-reset header when available.
//
// The retry logic:
// 1. Intercepts 429 responses and checks if retry is enabled and attempts remain
// 2. If x-ratelimit-reset header is present, waits that many seconds
// 3. Otherwise uses exponential backoff starting at InitialBackoff seconds
// 4. Retries up to MaxAttempts (including initial request)
func (rp *RetryPolicy) setupRateLimitRetry(client *resty.Client) {
	if !rp.EnableRetry || rp.MaxAttempts < 1 {
		return
	}

	// Configure retry condition: only retry on 429 (Too Many Requests)
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		if err != nil {
			return false
		}
		return r.StatusCode() == 429
	})

	// Configure retry attempts
	client.SetRetryCount(rp.MaxAttempts - 1) // SetRetryCount is retries, not total attempts

	// Configure backoff strategy
	client.SetRetryWaitTime(100 * time.Millisecond).
		SetRetryMaxWaitTime(10 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) time.Duration {
			// Check for x-ratelimit-reset header (seconds until window resets)
			if resetHeader := resp.Header().Get("x-ratelimit-reset"); resetHeader != "" {
				// If the header is present, parse it and use it as the wait duration
				var resetSeconds int
				_, _ = time.Parse("2006-01-02 15:04:05 MST", resetHeader)
				// Try to parse as integer (seconds)
				if _, err := time.Parse(time.RFC3339, resetHeader); err == nil {
					// If it's a timestamp, calculate the difference
					resetTime, _ := time.Parse(time.RFC3339, resetHeader)
					waitDuration := time.Until(resetTime)
					if waitDuration > 0 {
						return waitDuration
					}
				}
				// If it's a simple integer, use it as seconds
				if n, err := time.ParseDuration(resetHeader + "s"); err == nil {
					return n
				}
			}

			// Fallback: exponential backoff if no reset header or parsing failed
			// Calculate backoff: InitialBackoff * (BackoffMultiplier ^ attemptNumber)
			retryCount := client.RetryCount - resp.Request.Attempt + 1
			if retryCount < 1 {
				retryCount = 1
			}
			backoffSeconds := int(float64(rp.InitialBackoff) * math.Pow(rp.BackoffMultiplier, float64(retryCount-1)))
			return time.Duration(backoffSeconds) * time.Second
		})
}

// WaitForRateLimit blocks the current goroutine until the rate limit window resets.
// It uses the x-ratelimit-reset value from a RateLimitError if available, otherwise
// uses exponential backoff.
//
// This is useful for cases where the client chooses not to use automatic retries
// but still wants to wait before making the next request.
func WaitForRateLimit(ctx context.Context, rlErr *RateLimitError, policy *RetryPolicy) error {
	if rlErr == nil || policy == nil || !policy.EnableRetry {
		return nil
	}

	var waitDuration time.Duration

	if rlErr.Reset > 0 {
		// Use the Reset value from the error (seconds until window resets)
		waitDuration = time.Duration(rlErr.Reset) * time.Second
	} else {
		// Fallback to exponential backoff
		backoffSeconds := int(float64(policy.InitialBackoff) * math.Pow(policy.BackoffMultiplier, 1))
		waitDuration = time.Duration(backoffSeconds) * time.Second
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
