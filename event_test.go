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
var mockBatchEventsResponse = map[string]any{
	"events": []map[string]any{
		{
			"lago_id":                    "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"transaction_id":             "event_1234",
			"lago_customer_id":           nil,
			"code":                       "bm_code",
			"timestamp":                  "2025-07-03T15:35:00Z",
			"precise_total_amount_cents": "100",
			"properties": map[string]any{
				"value": "100",
			},
			"lago_subscription_id":     nil,
			"external_subscription_id": "sub_1234",
			"created_at":               "2025-07-03T15:35:22Z",
		},
	},
}

func batchHandlerFunc(c *qt.C, assertRequestFunc func(*qt.C, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestFunc(c, r)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockBatchEventsResponse)
	}))
}

func assertBatchEventListResponse(c *qt.C, result []Event) {
	c.Assert(result, qt.HasLen, 1)

	event := result[0]
	c.Assert(event.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(event.Code, qt.Equals, "bm_code")
	c.Assert(event.Timestamp.Format(time.RFC3339), qt.Equals, "2025-07-03T15:35:00Z")
	c.Assert(event.PreciseTotalAmountCents, qt.Equals, "100")
	c.Assert(event.Properties, qt.DeepEquals, map[string]any{
		"value": "100",
	})
	c.Assert(event.LagoSubscriptionID, qt.IsNil)
	c.Assert(event.ExternalSubscriptionID, qt.Equals, "sub_1234")
	c.Assert(event.CreatedAt.Format(time.RFC3339), qt.Equals, "2025-07-03T15:35:22Z")
}

func TestEventsBatch(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Event().Batch(context.Background(), []EventInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Post \"http://localhost:88888/api/v1/events/batch\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When events are provided", func(t *testing.T) {
		c := qt.New(t)

		server := batchHandlerFunc(c, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/events/batch")
		})

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.Event().Batch(context.Background(), []EventInput{
			{
				TransactionID:           "event_1234",
				ExternalSubscriptionID:  "sub_1234",
				Code:                    "bm_code",
				Timestamp:               "1751549700",
				PreciseTotalAmountCents: "100",
				Properties: map[string]any{
					"value": "100",
				},
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		assertBatchEventListResponse(c, result)
	})
}
