package lago_test

import (
	"context"
	"testing"
	"time"

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
		"name":                 "My Webhook Endpoint",
		"event_types":          []string{"customer.created"},
		"created_at":           "2025-01-01T10:00:00Z",
	},
}

func TestWebhookEndpointRequest_Create(t *testing.T) {
	t.Run("Wraps the body under the webhook_endpoint key", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/webhook_endpoints").
			MatchJSONBody(`{"webhook_endpoint":{"webhook_url":"https://example.com/webhook","signature_algo":"jwt","name":"My Webhook Endpoint","event_types":["customer.created"]}}`).
			MockResponse(mockWebhookEndpointResponse)
		defer server.Close()

		name := "My Webhook Endpoint"
		result, err := server.Client().WebhookEndpoint().Create(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: JWT,
			Name:          &name,
			EventTypes:    []string{"customer.created"},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookURL, qt.Equals, "https://example.com/webhook")
		c.Assert(*result.Name, qt.Equals, "My Webhook Endpoint")
		c.Assert(result.EventTypes, qt.DeepEquals, []string{"customer.created"})
	})
}

func TestWebhookEndpointRequest_Update(t *testing.T) {
	t.Run("Wraps the body under the webhook_endpoint key", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MatchJSONBody(`{"webhook_endpoint":{"webhook_url":"https://example.com/webhook","signature_algo":"jwt","name":"My Webhook Endpoint","event_types":["customer.created"]}}`).
			MockResponse(mockWebhookEndpointResponse)
		defer server.Close()

		name := "My Webhook Endpoint"
		result, err := server.Client().WebhookEndpoint().Update(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: JWT,
			Name:          &name,
			EventTypes:    []string{"customer.created"},
		}, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookURL, qt.Equals, "https://example.com/webhook")
		c.Assert(*result.Name, qt.Equals, "My Webhook Endpoint")
		c.Assert(result.EventTypes, qt.DeepEquals, []string{"customer.created"})
	})
}

func mockStandardWebhooksEndpointResponse(signingKeys any) map[string]any {
	webhookEndpoint := map[string]any{
		"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"webhook_url":          "https://example.com/webhook",
		"signature_algo":       "standard_webhooks",
		"name":                 "My Webhook Endpoint",
		"event_types":          []string{"customer.created"},
		"created_at":           "2025-01-01T10:00:00Z",
	}
	if signingKeys != nil {
		webhookEndpoint["signing_keys"] = signingKeys
	}

	return map[string]any{"webhook_endpoint": webhookEndpoint}
}

var mockSigningKeys = []map[string]any{
	{
		"lago_id":    "3c903c90-3c90-3c90-3c90-3c903c903c90",
		"secret_key": "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw",
		"created_at": "2025-01-01T10:00:00Z",
	},
	{
		"lago_id":    "4d904d90-4d90-4d90-4d90-4d904d904d90",
		"secret_key": "whsec_C2FVsBQIhrscChlQIMV+b5sSYspob7oD",
		"created_at": "2025-02-01T12:30:00Z",
	},
}

const standardWebhooksEndpointBody = `{"webhook_endpoint":{"webhook_url":"https://example.com/webhook","signature_algo":"standard_webhooks"}}`

func assertSigningKeys(c *qt.C, signingKeys []WebhookEndpointSigningKey) {
	c.Assert(signingKeys, qt.HasLen, 2)
	c.Assert(signingKeys[0].LagoID.String(), qt.Equals, "3c903c90-3c90-3c90-3c90-3c903c903c90")
	c.Assert(signingKeys[0].SecretKey, qt.Equals, "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw")
	c.Assert(signingKeys[0].CreatedAt.Equal(time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)), qt.IsTrue)
	c.Assert(signingKeys[1].LagoID.String(), qt.Equals, "4d904d90-4d90-4d90-4d90-4d904d904d90")
	c.Assert(signingKeys[1].SecretKey, qt.Equals, "whsec_C2FVsBQIhrscChlQIMV+b5sSYspob7oD")
	c.Assert(signingKeys[1].CreatedAt.Equal(time.Date(2025, 2, 1, 12, 30, 0, 0, time.UTC)), qt.IsTrue)
}

func TestWebhookEndpointRequest_StandardWebhooks(t *testing.T) {
	t.Run("Create sends standard_webhooks and decodes the signing keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/webhook_endpoints").
			MatchJSONBody(standardWebhooksEndpointBody).
			MockResponse(mockStandardWebhooksEndpointResponse(mockSigningKeys))
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Create(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: StandardWebhooks,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SignatureAlgo, qt.Equals, StandardWebhooks)
		assertSigningKeys(c, result.SigningKeys)
	})

	t.Run("Update sends standard_webhooks and decodes the signing keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MatchJSONBody(standardWebhooksEndpointBody).
			MockResponse(mockStandardWebhooksEndpointResponse(mockSigningKeys))
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Update(context.Background(), &WebhookEndpointInput{
			WebhookURL:    "https://example.com/webhook",
			SignatureAlgo: StandardWebhooks,
		}, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SignatureAlgo, qt.Equals, StandardWebhooks)
		assertSigningKeys(c, result.SigningKeys)
	})

	t.Run("Get decodes the signing keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(mockStandardWebhooksEndpointResponse(mockSigningKeys))
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		assertSigningKeys(c, result.SigningKeys)
	})

	t.Run("Get without signing_keys decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(mockStandardWebhooksEndpointResponse(nil))
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})

	t.Run("Get with an empty signing_keys list decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(mockStandardWebhooksEndpointResponse([]map[string]any{}))
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})

	t.Run("GetList without signing_keys decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		endpoint := mockStandardWebhooksEndpointResponse(nil)["webhook_endpoint"]
		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/webhook_endpoints").
			MockResponse(map[string]any{
				"webhook_endpoints": []any{endpoint},
				"meta": map[string]any{
					"current_page": 1,
					"total_pages":  1,
					"total_count":  1,
				},
			})
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().GetList(context.Background(), &WebhookEndpointListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WebhookEndpoints, qt.HasLen, 1)
		c.Assert(result.WebhookEndpoints[0].SignatureAlgo, qt.Equals, StandardWebhooks)
		c.Assert(result.WebhookEndpoints[0].SigningKeys, qt.HasLen, 0)
	})

	t.Run("Delete without signing_keys decodes with no keys", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/webhook_endpoints/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(mockStandardWebhooksEndpointResponse(nil))
		defer server.Close()

		result, err := server.Client().WebhookEndpoint().Delete(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.SignatureAlgo, qt.Equals, StandardWebhooks)
		c.Assert(result.SigningKeys, qt.HasLen, 0)
	})
}
