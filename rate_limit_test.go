package lago_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
)

func intPtr(v int) *int { return &v }

func TestRateLimitError_Error(t *testing.T) {
	c := qt.New(t)
	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Message:        "Too Many Requests",
		Limit:          intPtr(100),
		Remaining:      intPtr(0),
		Reset:          intPtr(60),
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
				"X-Ratelimit-Limit":     {"100"},
				"X-Ratelimit-Remaining": {"50"},
				"X-Ratelimit-Reset":     {"30"},
			},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
				Limit:          intPtr(100),
				Remaining:      intPtr(50),
				Reset:          intPtr(30),
			},
		},
		{
			name: "Missing headers",
			headers: http.Header{
				"X-Ratelimit-Limit": {"100"},
			},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
				Limit:          intPtr(100),
				Remaining:      nil,
				Reset:          nil,
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
				"X-Ratelimit-Limit":     {"not-a-number"},
				"X-Ratelimit-Remaining": {"also-not-a-number"},
				"X-Ratelimit-Reset":     {"still-not-a-number"},
			},
			expectedErr: &RateLimitError{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
				Limit:          nil,
				Remaining:      nil,
				Reset:          nil,
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
			c.Assert(rlErr.Limit, qt.DeepEquals, tt.expectedErr.Limit)
			c.Assert(rlErr.Remaining, qt.DeepEquals, tt.expectedErr.Remaining)
			c.Assert(rlErr.Reset, qt.DeepEquals, tt.expectedErr.Reset)
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
		name        string
		reset       *int
		attempt     int
		minWaitTime time.Duration
		maxWaitTime time.Duration
	}{
		{
			name:        "Wait with reset header",
			reset:       intPtr(1),
			attempt:     0,
			minWaitTime: 1 * time.Second,
			maxWaitTime: 2 * time.Second,
		},
		{
			name:        "Fallback exponential backoff attempt 0",
			reset:       nil,
			attempt:     0,
			minWaitTime: 1 * time.Second, // InitialBackoff * 2^0 = 1 second
			maxWaitTime: 2 * time.Second,
		},
		{
			name:        "Fallback exponential backoff attempt 1",
			reset:       nil,
			attempt:     1,
			minWaitTime: 2 * time.Second, // InitialBackoff * 2^1 = 2 seconds
			maxWaitTime: 3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			ctx, cancel := context.WithTimeout(context.Background(), tt.maxWaitTime+1*time.Second)
			defer cancel()

			rlErr := &RateLimitError{
				HTTPStatusCode: 429,
				Reset:          tt.reset,
			}

			policy := DefaultRetryPolicy()

			startTime := time.Now()
			err := WaitForRateLimit(ctx, rlErr, policy, tt.attempt)
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
		Reset:          intPtr(10), // Long wait
	}

	policy := DefaultRetryPolicy()

	// Cancel context immediately
	cancel()

	err := WaitForRateLimit(ctx, rlErr, policy, 0)

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
		Reset:          intPtr(10),
	}

	startTime := time.Now()
	err := WaitForRateLimit(ctx, rlErr, policy, 0)
	duration := time.Since(startTime)

	c.Assert(err, qt.IsNil)
	// Should return immediately without waiting (allow small overhead for execution)
	c.Assert(duration < 100*time.Millisecond, qt.IsTrue, qt.Commentf(
		"Wait duration too long: %v, should return immediately", duration))
}

func TestClientRateLimitRetry(t *testing.T) {
	c := qt.New(t)
	requestCount := 0

	// Create a custom httptest server that returns 429 on first request, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-ratelimit-limit", "100")
			w.Header().Set("x-ratelimit-remaining", "0")
			w.Header().Set("x-ratelimit-reset", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"status": 429, "error": "Too Many Requests"}`))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data": "success"}`))
		}
	}))
	defer server.Close()

	client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")

	ctx := context.Background()
	_, err := client.Get(ctx, &ClientRequest{
		Path: "test",
	})

	// Should have retried: first request got 429, second got 200
	c.Assert(requestCount, qt.Equals, 2)
	c.Assert(err == nil, qt.IsTrue, qt.Commentf("expected nil error after successful retry, got: %v", err))
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
