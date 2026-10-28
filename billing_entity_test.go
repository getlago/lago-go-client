package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockBillingEntityResponse = map[string]any{
	"billing_entity": map[string]any{
		"lago_id":    "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"name":       "Acme FR",
		"code":       "acme_fr",
		"country":    "FR",
		"einvoicing": true,
		"created_at": "2025-01-01T10:00:00Z",
	},
}

func TestBillingEntityRequest_Create(t *testing.T) {
	t.Run("Sends einvoicing when set to true", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/billing_entities").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","code":"acme_fr","billing_configuration":{},"einvoicing":true}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		einvoicing := true
		result, err := server.Client().BillingEntity().Create(context.Background(), &BillingEntityCreateInput{
			Name:       "Acme FR",
			Code:       "acme_fr",
			Einvoicing: &einvoicing,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Einvoicing, qt.IsTrue)
	})

	t.Run("Sends einvoicing when set to false", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/billing_entities").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","code":"acme_fr","billing_configuration":{},"einvoicing":false}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		einvoicing := false
		_, err := server.Client().BillingEntity().Create(context.Background(), &BillingEntityCreateInput{
			Name:       "Acme FR",
			Code:       "acme_fr",
			Einvoicing: &einvoicing,
		})
		c.Assert(err == nil, qt.IsTrue)
	})

	t.Run("Omits einvoicing when not set", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/billing_entities").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","code":"acme_fr","billing_configuration":{}}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		_, err := server.Client().BillingEntity().Create(context.Background(), &BillingEntityCreateInput{
			Name: "Acme FR",
			Code: "acme_fr",
		})
		c.Assert(err == nil, qt.IsTrue)
	})
}

func TestBillingEntityRequest_Update(t *testing.T) {
	t.Run("Sends einvoicing when set to true", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","billing_configuration":{},"einvoicing":true}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		einvoicing := true
		result, err := server.Client().BillingEntity().Update(context.Background(), "acme_fr", &BillingEntityUpdateInput{
			Name:       "Acme FR",
			Einvoicing: &einvoicing,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Einvoicing, qt.IsTrue)
	})

	t.Run("Sends einvoicing when set to false", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","billing_configuration":{},"einvoicing":false}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		einvoicing := false
		_, err := server.Client().BillingEntity().Update(context.Background(), "acme_fr", &BillingEntityUpdateInput{
			Name:       "Acme FR",
			Einvoicing: &einvoicing,
		})
		c.Assert(err == nil, qt.IsTrue)
	})

	t.Run("Omits einvoicing when not set", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","billing_configuration":{}}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		_, err := server.Client().BillingEntity().Update(context.Background(), "acme_fr", &BillingEntityUpdateInput{
			Name: "Acme FR",
		})
		c.Assert(err == nil, qt.IsTrue)
	})
}

func TestBillingEntityRequest_Get(t *testing.T) {
	t.Run("Decodes einvoicing", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		result, err := server.Client().BillingEntity().Get(context.Background(), "acme_fr")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "acme_fr")
		c.Assert(result.Einvoicing, qt.IsTrue)
	})
}
