package lago_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

func TestRateLimitError_Error(t *testing.T) {
	c := qt.New(t)
	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Message:        "Too Many Requests",
		Limit:          100,
		Remaining:      0,
		Reset:          60,
	}

	errStr := rlErr.Error()
	c.Assert(len(errStr) > 0, qt.IsTrue)

	// Verify the error string contains the rate limit info
	c.Assert(strings.Contains(errStr, "429"), qt.IsTrue)
}

func TestParseRateLimitHeaders(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expectedErr *RateLimitError
	}{
		{
			name: "All headers present",
			headers: http.Header{
				"x-ratelimit-limit":     {"100"},
				"x-ratelimit-remaining": {"50"},
				"x-ratelimit-reset":     {"30"},
			},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
				Limit:          100,
				Remaining:      50,
				Reset:          30,
			},
		},
		{
			name: "Missing headers",
			headers: http.Header{
				"x-ratelimit-limit": {"100"},
			},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
				Limit:          100,
				Remaining:      0,
				Reset:          0,
			},
		},
		{
			name:    "Empty headers",
			headers: http.Header{},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
			},
		},
		{
			name: "Invalid numeric values",
			headers: http.Header{
				"x-ratelimit-limit":     {"not-a-number"},
				"x-ratelimit-remaining": {"also-not-a-number"},
				"x-ratelimit-reset":     {"still-not-a-number"},
			},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
				Limit:          0,
				Remaining:      0,
				Reset:          0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			baseErr := &Error{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
			}

			rlErr := ParseRateLimitError(baseErr, tt.headers)

			c.Assert(rlErr.HTTPStatusCode, qt.Equals, tt.expectedErr.HTTPStatusCode)
			c.Assert(rlErr.Limit, qt.Equals, tt.expectedErr.Limit)
			c.Assert(rlErr.Remaining, qt.Equals, tt.expectedErr.Remaining)
			c.Assert(rlErr.Reset, qt.Equals, tt.expectedErr.Reset)
		})
	}
}

func TestDefaultRetryPolicy(t *testing.T) {
	c := qt.New(t)
	policy := DefaultRetryPolicy()

	c.Assert(policy.MaxAttempts, qt.Equals, 3)
	c.Assert(policy.EnableRetry, qt.IsTrue)
	c.Assert(policy.InitialBackoff, qt.Equals, 1)
	c.Assert(policy.BackoffMultiplier, qt.Equals, 2.0)
}

func TestWaitForRateLimit(t *testing.T) {
	tests := []struct {
		name           string
		resetSeconds   int
		minWaitTime    time.Duration
		maxWaitTime    time.Duration
	}{
		{
			name:           "Wait with reset header",
			resetSeconds:   1,
			minWaitTime:    1 * time.Second,
			maxWaitTime:    2 * time.Second,
		},
		{
			name:           "Fallback exponential backoff when reset is 0",
			resetSeconds:   0,
			minWaitTime:    1 * time.Second, // InitialBackoff * 2^0 = 1 second
			maxWaitTime:    2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			ctx, cancel := context.WithTimeout(context.Background(), tt.maxWaitTime+1*time.Second)
			defer cancel()

			rlErr := &RateLimitError{
				HTTPStatusCode: 429,
				Reset:          tt.resetSeconds,
			}

			policy := DefaultRetryPolicy()

			startTime := time.Now()
			err := WaitForRateLimit(ctx, rlErr, policy)
			duration := time.Since(startTime)

			c.Assert(err, qt.IsNil)
			c.Assert(duration >= tt.minWaitTime, qt.IsTrue, qt.Commentf(
				"Wait duration too short: %v, want at least %v", duration, tt.minWaitTime))
		})
	}
}

func TestWaitForRateLimitContextCancellation(t *testing.T) {
	c := qt.New(t)
	ctx, cancel := context.WithCancel(context.Background())

	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Reset:          10, // Long wait
	}

	policy := DefaultRetryPolicy()

	// Cancel context immediately
	cancel()

	err := WaitForRateLimit(ctx, rlErr, policy)

	c.Assert(err, qt.Equals, context.Canceled)
}

func TestRetryPolicyDisabled(t *testing.T) {
	c := qt.New(t)
	policy := &RetryPolicy{
		EnableRetry: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Reset:          10,
	}

	startTime := time.Now()
	err := WaitForRateLimit(ctx, rlErr, policy)
	duration := time.Since(startTime)

	c.Assert(err, qt.IsNil)
	// Should return immediately without waiting (allow small overhead for execution)
	c.Assert(duration < 100*time.Millisecond, qt.IsTrue, qt.Commentf(
		"Wait duration too long: %v, should return immediately", duration))
}

func TestClientRateLimitRetry(t *testing.T) {
	c := qt.New(t)
	requestCount := 0

	// Create a mock server that returns 429 on first request, then 200
	mockServer := lt.NewMockServer(c)
	defer func() {
		// Verify the server received requests
		c.Assert(requestCount >= 2, qt.IsTrue, qt.Commentf(
			"Expected at least 2 requests (initial + 1 retry), got %d", requestCount))
		mockServer.Close()
	}()

	// Track requests
	originalHandler := mockServer.server.Config.Handler
	mockServer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			// First request returns 429 with rate limit headers
			w.Header().Set("x-ratelimit-limit", "100")
			w.Header().Set("x-ratelimit-remaining", "0")
			w.Header().Set("x-ratelimit-reset", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"status": 429, "error": "Too Many Requests"}`))
		} else {
			// Subsequent requests return success
			w.Header().Set("x-ratelimit-limit", "100")
			w.Header().Set("x-ratelimit-remaining", "99")
			w.Header().Set("x-ratelimit-reset", "60")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": "success"}`))
		}
	})

	client := mockServer.Client()

	ctx := context.Background()
	_, err := client.Get(ctx, &ClientRequest{
		Path: "/api/v1/test",
	})

	// The client should either succeed after retry or fail with a rate limit error if retries exhausted
	c.Assert(requestCount >= 1, qt.IsTrue)
}

func TestSetRetryPolicy(t *testing.T) {
	c := qt.New(t)
	client := New()

	// Verify default policy
	c.Assert(client.RetryPolicy, qt.Not(qt.IsNil))
	c.Assert(client.RetryPolicy.EnableRetry, qt.IsTrue)

	// Create custom policy
	customPolicy := &RetryPolicy{
		MaxAttempts:       5,
		EnableRetry:       true,
		InitialBackoff:    2,
		BackoffMultiplier: 1.5,
	}

	client.SetRetryPolicy(customPolicy)

	c.Assert(client.RetryPolicy.MaxAttempts, qt.Equals, 5)

	// Test disabling retries
	client.SetRetryPolicy(nil)
	c.Assert(client.RetryPolicy.EnableRetry, qt.IsFalse)
}
