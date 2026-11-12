package lago_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

const v2WebhookEndpointID = "1a901a90-1a90-1a90-1a90-1a901a901a90"

func standardWebhooksEndpointInput() *WebhookEndpointInput {
	return &WebhookEndpointInput{
		WebhookURL:    "https://example.com/webhook",
		SignatureAlgo: StandardWebhooks,
	}
}

func TestWebhookEndpointV2Request_Get(t *testing.T) {
	t.Run("Sends expands[]=signing_keys and decodes the signing keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v2/webhook_endpoints/" + v2WebhookEndpointID).
			MatchQuery(map[string][]string{"expands[]": {"signing_keys"}}).
			MockResponse(mockStandardWebhooksEndpointResponse(mockSigningKeys))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Get(context.Background(), v2WebhookEndpointID, WebhookEndpointV2ExpandSigningKeys)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SignatureAlgo, qt.Equals, StandardWebhooks)
		assertSigningKeys(c, result.SigningKeys)
	})

	t.Run("Sends no query without an expansion and decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v2/webhook_endpoints/" + v2WebhookEndpointID).
			MatchQuery("").
			MockResponse(mockStandardWebhooksEndpointResponse(nil))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Get(context.Background(), v2WebhookEndpointID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookURL, qt.Equals, "https://example.com/webhook")
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})

	t.Run("Returns a v2 404 as an Error", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v2/webhook_endpoints/"+v2WebhookEndpointID).
			MockResponseWithCode(404, `{"code":"resource_not_found","details":[{"code":"not_found","field":"webhook_endpoint"}]}`)
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Get(context.Background(), v2WebhookEndpointID)
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.HTTPStatusCode, qt.Equals, 404)
		c.Assert(err.ErrorCode, qt.Equals, "resource_not_found")
		c.Assert(err.Details, qt.DeepEquals, []ErrorDetailV2{{Code: "not_found", Field: "webhook_endpoint"}})
	})
}

func TestWebhookEndpointV2Request_GetList(t *testing.T) {
	t.Run("Sends the cursor params set and decodes the cursor metadata", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v2/webhook_endpoints").
			MatchQuery(map[string][]string{
				"limit":               {"2"},
				"after":               {"cursor_a"},
				"include_total_count": {"true"},
			}).
			MockResponse(map[string]any{
				"webhook_endpoints": []map[string]any{
					mockStandardWebhooksEndpointResponse(nil)["webhook_endpoint"].(map[string]any),
				},
				"meta": map[string]any{
					"next_cursor": "cursor_b",
					"prev_cursor": nil,
					"total_count": 3,
				},
			})
		defer server.Close()

		limit := 2
		after := "cursor_a"
		includeTotalCount := true
		result, err := server.Client().WebhookEndpointV2().GetList(context.Background(), &WebhookEndpointV2ListInput{
			Limit:             &limit,
			After:             &after,
			IncludeTotalCount: &includeTotalCount,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookEndpoints, qt.HasLen, 1)
		c.Assert(result.WebhookEndpoints[0].SignatureAlgo, qt.Equals, StandardWebhooks)
		c.Assert(result.WebhookEndpoints[0].SigningKeys, qt.HasLen, 0)
		c.Assert(*result.Meta.NextCursor, qt.Equals, "cursor_b")
		c.Assert(result.Meta.PrevCursor, qt.IsNil)
		c.Assert(*result.Meta.TotalCount, qt.Equals, 3)
		c.Assert(result.Meta.TotalCountEstimated, qt.IsFalse)
	})

	t.Run("Sends the before cursor", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v2/webhook_endpoints").
			MatchQuery(map[string][]string{"before": {"cursor_b"}}).
			MockResponse(map[string]any{
				"webhook_endpoints": []map[string]any{},
				"meta":              map[string]any{"next_cursor": nil, "prev_cursor": "cursor_a", "total_count_estimated": true},
			})
		defer server.Close()

		before := "cursor_b"
		result, err := server.Client().WebhookEndpointV2().GetList(context.Background(), &WebhookEndpointV2ListInput{Before: &before})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookEndpoints, qt.HasLen, 0)
		c.Assert(result.Meta.NextCursor, qt.IsNil)
		c.Assert(*result.Meta.PrevCursor, qt.Equals, "cursor_a")
		c.Assert(result.Meta.TotalCount, qt.IsNil)
		c.Assert(result.Meta.TotalCountEstimated, qt.IsTrue)
	})

	t.Run("Sends no query with a nil input", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v2/webhook_endpoints").
			MatchQuery("").
			MockResponse(map[string]any{"webhook_endpoints": []map[string]any{}, "meta": map[string]any{}})
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().GetList(context.Background(), nil)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookEndpoints, qt.HasLen, 0)
	})
}

func TestWebhookEndpointV2Request_Create(t *testing.T) {
	t.Run("Sends standard_webhooks with expands[]=signing_keys and decodes the signing keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v2/webhook_endpoints").
			MatchQuery(map[string][]string{"expands[]": {"signing_keys"}}).
			MatchJSONBody(standardWebhooksEndpointBody).
			MockResponse(mockStandardWebhooksEndpointResponse(mockSigningKeys))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Create(context.Background(), standardWebhooksEndpointInput(), WebhookEndpointV2ExpandSigningKeys)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SignatureAlgo, qt.Equals, StandardWebhooks)
		assertSigningKeys(c, result.SigningKeys)
	})

	t.Run("Sends no query without an expansion and decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v2/webhook_endpoints").
			MatchQuery("").
			MatchJSONBody(standardWebhooksEndpointBody).
			MockResponse(mockStandardWebhooksEndpointResponse(nil))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Create(context.Background(), standardWebhooksEndpointInput())
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})

	t.Run("Returns a v2 422 with the allowed values of the detail", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v2/webhook_endpoints").
			MockResponseWithCode(422, map[string]any{
				"code": "validation_errors",
				"details": []map[string]any{
					{
						"code":  "value_is_invalid",
						"field": "signature_algo",
						"meta":  map[string]any{"allowed_values": []string{"jwt", "hmac", "standard_webhooks"}},
					},
				},
			})
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Create(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: SignatureAlgo("unknown"),
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(err.HTTPStatusCode, qt.Equals, 422)
		c.Assert(err.ErrorCode, qt.Equals, "validation_errors")
		c.Assert(err.ErrorDetail, qt.IsNil)
		c.Assert(err.Details, qt.HasLen, 1)
		c.Assert(err.Details[0].Code, qt.Equals, "value_is_invalid")
		c.Assert(err.Details[0].Field, qt.Equals, "signature_algo")
		c.Assert(err.Details[0].Meta["allowed_values"], qt.DeepEquals, []any{"jwt", "hmac", "standard_webhooks"})
	})
}

func TestWebhookEndpointV2Request_Update(t *testing.T) {
	t.Run("Sends standard_webhooks with expands[]=signing_keys and decodes the signing keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v2/webhook_endpoints/" + v2WebhookEndpointID).
			MatchQuery(map[string][]string{"expands[]": {"signing_keys"}}).
			MatchJSONBody(standardWebhooksEndpointBody).
			MockResponse(mockStandardWebhooksEndpointResponse(mockSigningKeys))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Update(context.Background(), standardWebhooksEndpointInput(), v2WebhookEndpointID, WebhookEndpointV2ExpandSigningKeys)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SignatureAlgo, qt.Equals, StandardWebhooks)
		assertSigningKeys(c, result.SigningKeys)
	})

	t.Run("Sends no query without an expansion and decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v2/webhook_endpoints/" + v2WebhookEndpointID).
			MatchQuery("").
			MatchJSONBody(standardWebhooksEndpointBody).
			MockResponse(mockStandardWebhooksEndpointResponse(nil))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Update(context.Background(), standardWebhooksEndpointInput(), v2WebhookEndpointID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})
}

func TestWebhookEndpointV2Request_Delete(t *testing.T) {
	t.Run("Sends DELETE to the v2 path", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v2/webhook_endpoints/" + v2WebhookEndpointID).
			MockResponse(mockStandardWebhooksEndpointResponse(nil))
		defer server.Close()

		result, err := server.Client().WebhookEndpointV2().Delete(context.Background(), v2WebhookEndpointID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.LagoID.String(), qt.Equals, v2WebhookEndpointID)
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})
}

func TestWebhookEndpointV2Request_BaseURL(t *testing.T) {
	t.Run("Replaces the v1 path of a base URL set on the HTTP client", func(t *testing.T) {
		c := qt.New(t)

		var gotPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"webhook_endpoint":{"lago_id":"` + v2WebhookEndpointID + `"}}`))
		}))
		defer server.Close()

		client := New().SetApiKey("test_api_key")
		client.HttpClient.SetBaseURL(server.URL + "/api/v1/")

		result, err := client.WebhookEndpointV2().Get(context.Background(), v2WebhookEndpointID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.LagoID.String(), qt.Equals, v2WebhookEndpointID)
		c.Assert(gotPath, qt.Equals, "/api/v2/webhook_endpoints/"+v2WebhookEndpointID)
	})

	t.Run("Keeps v1 requests of the same client on the v1 path", func(t *testing.T) {
		c := qt.New(t)

		var gotPaths []string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPaths = append(gotPaths, r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"webhook_endpoint":{"lago_id":"` + v2WebhookEndpointID + `"}}`))
		}))
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")

		_, err := client.WebhookEndpointV2().Get(context.Background(), v2WebhookEndpointID)
		c.Assert(err == nil, qt.IsTrue)
		_, err = client.WebhookEndpoint().Get(context.Background(), v2WebhookEndpointID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(gotPaths, qt.DeepEquals, []string{
			"/api/v2/webhook_endpoints/" + v2WebhookEndpointID,
			"/api/v1/webhook_endpoints/" + v2WebhookEndpointID,
		})
	})
}
