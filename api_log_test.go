package lago

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

// Mock JSON response structure
var GetListMockResponse = map[string]interface{}{
	"api_logs": []map[string]interface{}{
		{
			"request_id":     "8fae2f0e-fe8e-44d3-bbf7-1c552eba3a24",
			"client":         "LagoClient.0.0.0",
			"http_method":    "post",
			"http_status":    200,
			"request_origin": "https://request-origin.com",
			"request_path":   "/api/v1/billable_metrics",
			"api_version": "v1",
			"logged_at":   "2025-06-20T14:34:25Z",
			"created_at":  "2025-06-20T14:34:25Z",
			"request_body": map[string]interface{}{
				"billable_metric": map[string]interface{}{
					"name":               "Storage",
					"code":               "storage",
					"aggregation_type":   "sum_agg",
					"description":        "GB of storage used in my application",
					"recurring":          false,
					"rounding_function":  "round",
					"rounding_precision": 2,
					"field_name":         "gb",
					"weighted_interval":  "seconds",
				},
			},
			"request_response": map[string]interface{}{
				"billable_metric": map[string]interface{}{
					"lago_id":                  "4caa4455-07f2-4760-a697-f2644005eb43",
					"name":                     "Storage",
					"code":                     "storage",
					"description":              "GB of storage used in my application",
					"aggregation_type":         "sum_agg",
					"weighted_interval":        "seconds",
					"recurring":                false,
					"rounding_function":        "round",
					"rounding_precision":       2,
					"created_at":               "2025-06-20T14:34:25Z",
					"field_name":               "gb",
					"expression":               nil,
					"active_subscriptions_count": 0,
					"draft_invoices_count":       0,
					"plans_count":                0,
				},
			},
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

var GetMockResponse = map[string]interface{}{
	"api_log": map[string]interface{}{
		"request_id":     "8fae2f0e-fe8e-44d3-bbf7-1c552eba3a24",
		"client":         "LagoClient.0.0.0",
		"http_method":    "post",
		"http_status":    200,
		"request_origin": "https://request-origin.com",
		"request_path":   "/api/v1/billable_metrics",
		"api_version": "v1",
		"logged_at":   "2025-06-20T14:34:25Z",
		"created_at":  "2025-06-20T14:34:25Z",
		"request_body": map[string]interface{}{
			"billable_metric": map[string]interface{}{
				"name":               "Storage",
				"code":               "storage",
				"aggregation_type":   "sum_agg",
				"description":        "GB of storage used in my application",
				"recurring":          false,
				"rounding_function":  "round",
				"rounding_precision": 2,
				"field_name":         "gb",
				"weighted_interval":  "seconds",
			},
		},
		"request_response": map[string]interface{}{
			"billable_metric": map[string]interface{}{
				"lago_id":                  "4caa4455-07f2-4760-a697-f2644005eb43",
				"name":                     "Storage",
				"code":                     "storage",
				"description":              "GB of storage used in my application",
				"aggregation_type":         "sum_agg",
				"weighted_interval":        "seconds",
				"recurring":                false,
				"rounding_function":        "round",
				"rounding_precision":       2,
				"created_at":               "2025-06-20T14:34:25Z",
				"field_name":               "gb",
				"expression":               nil,
				"active_subscriptions_count": 0,
				"draft_invoices_count":       0,
				"plans_count":                0,
			},
		},
	},
}

func apiLogTestServer(c *qt.C, response any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func assertApiLog(c *qt.C, apiLog ApiLog) {
	c.Assert(apiLog.RequestId.String(), qt.Equals, "8fae2f0e-fe8e-44d3-bbf7-1c552eba3a24")
	c.Assert(apiLog.ApiVersion, qt.Equals, "v1")
	c.Assert(apiLog.HttpMethod, qt.Equals, PostMethod)
	c.Assert(apiLog.HttpStatus, qt.Equals, 200)
	c.Assert(apiLog.RequestBody, qt.IsNotNil)
	c.Assert(apiLog.RequestOrigin, qt.Equals, "https://request-origin.com")
	c.Assert(apiLog.RequestPath, qt.Equals, "/api/v1/billable_metrics")
	c.Assert(apiLog.RequestResponse, qt.IsNotNil)
	c.Assert(apiLog.LoggedAt.Format(time.RFC3339), qt.Equals, "2025-06-20T14:34:25Z")
	c.Assert(apiLog.CreatedAt.Format(time.RFC3339), qt.Equals, "2025-06-20T14:34:25Z")
}

func TestApiLogGetList(t *testing.T) {
	t.Run("When query for all api logs", func(t *testing.T) {
		c := qt.New(t)

		server := apiLogTestServer(c, GetListMockResponse)
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.ApiLog().GetList(context.Background(), &ApiLogListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.ApiLogs, qt.HasLen, 1)
		apiLog := result.ApiLogs[0]
		assertApiLog(c, apiLog)
	})
}

func TestApiLogGet(t *testing.T) {
	t.Run("When query for a specific api log", func(t *testing.T) {
		c := qt.New(t)

		server := apiLogTestServer(c, GetMockResponse)
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.ApiLog().Get(context.Background(), "8fae2f0e-fe8e-44d3-bbf7-1c552eba3a24")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		assertApiLog(c, *result)
	})
}