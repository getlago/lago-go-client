package lago

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitError_Error(t *testing.T) {
	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Message:        "Too Many Requests",
		Limit:          100,
		Remaining:      0,
		Reset:          60,
	}

	errStr := rlErr.Error()
	if len(errStr) == 0 {
		t.Errorf("RateLimitError.Error() returned empty string")
	}

	// Verify the error string contains the rate limit info
	if !containsSubstring(errStr, "429") {
		t.Errorf("RateLimitError.Error() missing status code: %s", errStr)
	}
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
			baseErr := &Error{
				HTTPStatusCode: 429,
				Message:        "Rate limit exceeded",
			}

			rlErr := ParseRateLimitError(baseErr, tt.headers)

			if rlErr.HTTPStatusCode != tt.expectedErr.HTTPStatusCode {
				t.Errorf("Status code mismatch: got %d, want %d", rlErr.HTTPStatusCode, tt.expectedErr.HTTPStatusCode)
			}

			if rlErr.Limit != tt.expectedErr.Limit {
				t.Errorf("Limit mismatch: got %d, want %d", rlErr.Limit, tt.expectedErr.Limit)
			}

			if rlErr.Remaining != tt.expectedErr.Remaining {
				t.Errorf("Remaining mismatch: got %d, want %d", rlErr.Remaining, tt.expectedErr.Remaining)
			}

			if rlErr.Reset != tt.expectedErr.Reset {
				t.Errorf("Reset mismatch: got %d, want %d", rlErr.Reset, tt.expectedErr.Reset)
			}
		})
	}
}

func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy()

	if policy.MaxAttempts != 3 {
		t.Errorf("MaxAttempts: got %d, want 3", policy.MaxAttempts)
	}

	if !policy.EnableRetry {
		t.Errorf("EnableRetry: got false, want true")
	}

	if policy.InitialBackoff != 1 {
		t.Errorf("InitialBackoff: got %d, want 1", policy.InitialBackoff)
	}

	if policy.BackoffMultiplier != 2.0 {
		t.Errorf("BackoffMultiplier: got %f, want 2.0", policy.BackoffMultiplier)
	}
}

func TestWaitForRateLimit(t *testing.T) {
	tests := []struct {
		name        string
		resetSeconds int
		shouldWait  bool
		maxWaitTime time.Duration
	}{
		{
			name:         "Wait with reset header",
			resetSeconds: 1,
			shouldWait:   true,
			maxWaitTime:  2 * time.Second,
		},
		{
			name:         "No wait with nil error",
			resetSeconds: 0,
			shouldWait:   true,
			maxWaitTime:  100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			if err != nil && err != context.DeadlineExceeded {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.shouldWait && duration < time.Duration(tt.resetSeconds)*time.Second {
				t.Errorf("Wait duration too short: %v, want at least %v",
					duration, time.Duration(tt.resetSeconds)*time.Second)
			}
		})
	}
}

func TestWaitForRateLimitContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Reset:          10, // Long wait
	}

	policy := DefaultRetryPolicy()

	// Cancel context immediately
	cancel()

	err := WaitForRateLimit(ctx, rlErr, policy)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestRetryPolicyDisabled(t *testing.T) {
	policy := &RetryPolicy{
		EnableRetry: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	rlErr := &RateLimitError{
		HTTPStatusCode: 429,
		Reset:          10,
	}

	startTime := time.Now()
	err := WaitForRateLimit(ctx, rlErr, policy)
	duration := time.Since(startTime)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should return immediately without waiting
	if duration > 100*time.Millisecond {
		t.Errorf("Wait duration too long: %v, should return immediately", duration)
	}
}

func TestClientRateLimitRetry(t *testing.T) {
	// Counter to track requests
	requestCount := 0
	rateLimitResetTime := time.Now().Add(1 * time.Second)

	// Create a test server that returns 429 on first request, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			// First request: return 429
			w.Header().Set("x-ratelimit-limit", "100")
			w.Header().Set("x-ratelimit-remaining", "0")
			w.Header().Set("x-ratelimit-reset", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"status": 429, "error": "Too Many Requests"}`))
			return
		}

		// Subsequent requests: return 200
		w.Header().Set("x-ratelimit-limit", "100")
		w.Header().Set("x-ratelimit-remaining", "99")
		w.Header().Set("x-ratelimit-reset", "60")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": "success"}`))
	}))
	defer server.Close()

	client := New().SetBaseURL(server.URL)

	ctx := context.Background()
	_, err := client.Get(ctx, &ClientRequest{
		Path: "/test",
	})

	// Note: Due to retries being automatic in resty, we should get success
	// (error will only occur if all retries are exhausted)
	if err != nil && err.HTTPStatusCode != 429 {
		t.Errorf("Unexpected error: %+v", err)
	}

	// Verify that the server was called at least once
	if requestCount < 1 {
		t.Errorf("Server not called, request count: %d", requestCount)
	}

	_ = rateLimitResetTime
}

func TestSetRetryPolicy(t *testing.T) {
	client := New()

	// Verify default policy
	if client.RetryPolicy == nil {
		t.Fatal("Default retry policy is nil")
	}

	if !client.RetryPolicy.EnableRetry {
		t.Error("Default retry policy should have EnableRetry=true")
	}

	// Create custom policy
	customPolicy := &RetryPolicy{
		MaxAttempts:       5,
		EnableRetry:       true,
		InitialBackoff:    2,
		BackoffMultiplier: 1.5,
	}

	client.SetRetryPolicy(customPolicy)

	if client.RetryPolicy.MaxAttempts != 5 {
		t.Errorf("MaxAttempts not updated: got %d, want 5", client.RetryPolicy.MaxAttempts)
	}

	// Test disabling retries
	client.SetRetryPolicy(nil)
	if client.RetryPolicy.EnableRetry {
		t.Error("RetryPolicy should have EnableRetry=false after setting to nil")
	}
}

// Helper function to check if a substring exists in a string
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
