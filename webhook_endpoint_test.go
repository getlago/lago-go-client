package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockWebhookEndpointResponse = map[string]any{
	"webhook_endpoint": map[string]any{
		"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"webhook_url":          "https://example.com/webhook",
		"signature_algo":       "jwt",
		"created_at":           "2025-01-01T10:00:00Z",
	},
}

func TestWebhookEndpointRequest_Create(t *testing.T) {
	t.Run("Wraps the body under the webhook_endpoint key", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/webhook_endpoints").
			MatchJSONBody(`{"webhook_endpoint":{"webhook_url":"https://example.com/webhook","signature_algo":"jwt"}}`).
			MockResponse(mockWebhookEndpointResponse)
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Create(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: JWT,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookURL, qt.Equals, "https://example.com/webhook")
	})
}

func TestWebhookEndpointRequest_Update(t *testing.T) {
	t.Run("Wraps the body under the webhook_endpoint key", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MatchJSONBody(`{"webhook_endpoint":{"webhook_url":"https://example.com/webhook","signature_algo":"jwt"}}`).
			MockResponse(mockWebhookEndpointResponse)
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Update(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: JWT,
		}, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookURL, qt.Equals, "https://example.com/webhook")
	})
}
