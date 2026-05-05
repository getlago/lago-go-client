package lago_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

func TestRateLimitInfo_UsagePct(t *testing.T) {
	c := qt.New(t)

	// Full headers, 80% used
	info := &RateLimitInfo{Limit: intPtr(100), Remaining: intPtr(20)}
	pct, ok := info.UsagePct()
	c.Assert(ok, qt.IsTrue)
	c.Assert(pct, qt.Equals, 0.80)

	// Saturated
	info = &RateLimitInfo{Limit: intPtr(100), Remaining: intPtr(0)}
	pct, ok = info.UsagePct()
	c.Assert(ok, qt.IsTrue)
	c.Assert(pct, qt.Equals, 1.0)

	// Missing limit
	info = &RateLimitInfo{Remaining: intPtr(20)}
	_, ok = info.UsagePct()
	c.Assert(ok, qt.IsFalse)

	// Missing remaining
	info = &RateLimitInfo{Limit: intPtr(100)}
	_, ok = info.UsagePct()
	c.Assert(ok, qt.IsFalse)

	// Zero limit
	info = &RateLimitInfo{Limit: intPtr(0), Remaining: intPtr(0)}
	_, ok = info.UsagePct()
	c.Assert(ok, qt.IsFalse)

	// Nil receiver
	var nilInfo *RateLimitInfo
	_, ok = nilInfo.UsagePct()
	c.Assert(ok, qt.IsFalse)
}

func TestOnRateLimitInfo_FiresOnSuccess(t *testing.T) {
	c := qt.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "20")
		w.Header().Set("x-ratelimit-reset", "5")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": "success"}`))
	}))
	defer server.Close()

	var (
		mu       sync.Mutex
		captured []*RateLimitInfo
	)

	client := New().SetBaseURL(server.URL).SetApiKey("k")
	client.RetryPolicy.OnRateLimitInfo = func(info *RateLimitInfo) {
		mu.Lock()
		defer mu.Unlock()
		captured = append(captured, info)
	}

	_, err := client.Get(context.Background(), &ClientRequest{Path: "test"})
	c.Assert(err == nil, qt.IsTrue, qt.Commentf("err: %v", err))

	mu.Lock()
	defer mu.Unlock()
	c.Assert(len(captured), qt.Equals, 1)
	info := captured[0]
	c.Assert(info.Limit, qt.Not(qt.IsNil))
	c.Assert(*info.Limit, qt.Equals, 100)
	c.Assert(*info.Remaining, qt.Equals, 20)
	c.Assert(*info.Reset, qt.Equals, 5)
	pct, ok := info.UsagePct()
	c.Assert(ok, qt.IsTrue)
	c.Assert(pct, qt.Equals, 0.80)
	c.Assert(info.Method, qt.Equals, "GET")
}

func TestOnRateLimitInfo_NotCalledWhenHeadersAbsent(t *testing.T) {
	c := qt.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": "success"}`))
	}))
	defer server.Close()

	called := 0
	client := New().SetBaseURL(server.URL).SetApiKey("k")
	client.RetryPolicy.OnRateLimitInfo = func(info *RateLimitInfo) {
		called++
	}

	_, err := client.Get(context.Background(), &ClientRequest{Path: "test"})
	c.Assert(err == nil, qt.IsTrue, qt.Commentf("err: %v", err))
	c.Assert(called, qt.Equals, 0)
}

func TestOnRateLimitInfo_PanicIsRecovered(t *testing.T) {
	c := qt.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "1")
		w.Header().Set("x-ratelimit-reset", "5")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": "success"}`))
	}))
	defer server.Close()

	client := New().SetBaseURL(server.URL).SetApiKey("k")
	client.RetryPolicy.OnRateLimitInfo = func(info *RateLimitInfo) {
		panic("intentional")
	}

	_, err := client.Get(context.Background(), &ClientRequest{Path: "test"})
	c.Assert(err == nil, qt.IsTrue, qt.Commentf("err: %v", err))
}

func TestOnRateLimitInfo_FiresOnceAfter429RetrySucceeds(t *testing.T) {
	c := qt.New(t)
	requests := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			w.Header().Set("x-ratelimit-reset", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"status":429,"error":"Too Many Requests"}`))
			return
		}
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "50")
		w.Header().Set("x-ratelimit-reset", "5")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":"ok"}`))
	}))
	defer server.Close()

	var captured []*RateLimitInfo
	client := New().SetBaseURL(server.URL).SetApiKey("k")
	client.RetryPolicy.OnRateLimitInfo = func(info *RateLimitInfo) {
		captured = append(captured, info)
	}

	_, err := client.Get(context.Background(), &ClientRequest{Path: "test"})
	c.Assert(err == nil, qt.IsTrue, qt.Commentf("err: %v", err))
	c.Assert(requests, qt.Equals, 2)
	c.Assert(len(captured), qt.Equals, 1)
	c.Assert(*captured[0].Remaining, qt.Equals, 50)
}

func TestLoggingRateLimitObserver_LogsAboveThreshold(t *testing.T) {
	c := qt.New(t)

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	observer := NewLoggingRateLimitObserver([]float64{0.80, 0.90, 0.95}, logger)
	observer(&RateLimitInfo{
		Limit:     intPtr(100),
		Remaining: intPtr(4), // 96% used
		Reset:     intPtr(10),
		Method:    "GET",
		URL:       "https://x",
	})

	out := buf.String()
	c.Assert(strings.Contains(out, "96%"), qt.IsTrue, qt.Commentf("logger output: %q", out))
}

func TestLoggingRateLimitObserver_SilentBelowThreshold(t *testing.T) {
	c := qt.New(t)

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	observer := NewLoggingRateLimitObserver([]float64{0.80}, logger)
	observer(&RateLimitInfo{
		Limit:     intPtr(100),
		Remaining: intPtr(50), // 50% used
	})

	c.Assert(buf.Len(), qt.Equals, 0)
}

func TestLoggingRateLimitObserver_DefaultThresholds(t *testing.T) {
	c := qt.New(t)

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	observer := NewLoggingRateLimitObserver(nil, logger)
	observer(&RateLimitInfo{Limit: intPtr(100), Remaining: intPtr(15)}) // 85%
	c.Assert(strings.Contains(buf.String(), "85%"), qt.IsTrue)
}
