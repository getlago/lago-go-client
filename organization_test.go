package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"

	"github.com/google/uuid"
)

var mockOrganizationResponse = `{
	"organization": {
		"name": "Lago",
		"email": "billing@getlago.com",
		"address_line1": "123 Main Street",
		"address_line2": "Suite 100",
		"city": "Paris",
		"zipcode": "75001",
		"state": "Ile-de-France",
		"country": "FR",
		"default_currency": "EUR",
		"legal_name": "Lago SAS",
		"legal_number": "FR123456789",
		"document_numbering": "per_organization",
		"document_number_prefix": "LAGO",
		"net_payment_term": 30,
		"tax_identification_number": "FR00123456789",
		"webhook_url": "https://example.com/webhook",
		"webhook_urls": ["https://example.com/webhook", "https://example.com/webhook2"],
		"timezone": "Europe/Paris",
		"email_settings": ["invoice.finalized", "credit_note.created"],
		"finalize_zero_amount_invoice": true,
		"events_store": "clickhouse",
		"billing_configuration": {
			"invoice_grace_period": 3,
			"invoice_footer": "Thanks for your business",
			"document_locale": "en"
		},
		"taxes": [{
			"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
			"name": "VAT",
			"code": "vat_20",
			"rate": 20.0,
			"description": "Standard VAT",
			"applied_to_organization": true,
			"created_at": "2022-04-29T08:59:51Z"
		}],
		"created_at": "2022-04-29T08:59:51Z"
	}
}`

func TestOrganization_Update(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Organization().Update(context.Background(), &OrganizationInput{
			Name: "Lago",
		})

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/organizations").
			MatchJSONBody(`{
				"organization": {
					"name": "Lago",
					"email": "billing@getlago.com",
					"address_line1": "123 Main Street",
					"address_line2": "Suite 100",
					"city": "Paris",
					"zipcode": "75001",
					"state": "Ile-de-France",
					"country": "FR",
					"default_currency": "EUR",
					"legal_name": "Lago SAS",
					"legal_number": "FR123456789",
					"document_numbering": "per_organization",
					"document_number_prefix": "LAGO",
					"net_payment_term": 30,
					"tax_identification_number": "FR00123456789",
					"webhook_url": "https://example.com/webhook",
					"timezone": "Europe/Paris",
					"email_settings": ["invoice.finalized", "credit_note.created"],
					"finalize_zero_amount_invoice": true,
					"billing_configuration": {
						"invoice_grace_period": 3,
						"invoice_footer": "Thanks for your business",
						"document_locale": "en"
					}
				}
			}`).
			MockResponse(mockOrganizationResponse)
		defer server.Close()

		result, err := server.Client().Organization().Update(context.Background(), &OrganizationInput{
			Name:                      "Lago",
			Email:                     "billing@getlago.com",
			AddressLine1:              "123 Main Street",
			AddressLine2:              "Suite 100",
			City:                      "Paris",
			Zipcode:                   "75001",
			State:                     "Ile-de-France",
			Country:                   "FR",
			DefaultCurrency:           EUR,
			LegalName:                 "Lago SAS",
			LegalNumber:               "FR123456789",
			DocumentNumbering:         DocumentNumberingPerOrganization,
			DocumentNumberPrefix:      "LAGO",
			NetPaymentTerm:            30,
			TaxIdentificationNumber:   "FR00123456789",
			WebhookURL:                "https://example.com/webhook",
			Timezone:                  "Europe/Paris",
			EmailSettings:             []string{"invoice.finalized", "credit_note.created"},
			FinalizeZeroAmountInvoice: true,
			BillingConfiguration: OrganizationBillingConfigurationInput{
				InvoiceGracePeriod: 3,
				InvoiceFooter:      "Thanks for your business",
				DocumentLocale:     "en",
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Name, qt.Equals, "Lago")
		c.Assert(result.Email, qt.Equals, "billing@getlago.com")
		c.Assert(result.AddressLine1, qt.Equals, "123 Main Street")
		c.Assert(result.AddressLine2, qt.Equals, "Suite 100")
		c.Assert(result.City, qt.Equals, "Paris")
		c.Assert(result.Zipcode, qt.Equals, "75001")
		c.Assert(result.State, qt.Equals, "Ile-de-France")
		c.Assert(result.Country, qt.Equals, "FR")
		c.Assert(result.DefaultCurrency, qt.Equals, EUR)
		c.Assert(result.LegalName, qt.Equals, "Lago SAS")
		c.Assert(result.LegalNumber, qt.Equals, "FR123456789")
		c.Assert(result.DocumentNumbering, qt.Equals, DocumentNumberingPerOrganization)
		c.Assert(result.DocumentNumberPrefix, qt.Equals, "LAGO")
		c.Assert(result.NetPaymentTerm, qt.Equals, 30)
		c.Assert(result.TaxIdentificationNumber, qt.Equals, "FR00123456789")
		c.Assert(result.WebhookURL, qt.Equals, "https://example.com/webhook")
		c.Assert(result.WebhookURLs, qt.DeepEquals, []string{"https://example.com/webhook", "https://example.com/webhook2"})
		c.Assert(result.Timezone, qt.Equals, "Europe/Paris")
		c.Assert(result.EmailSettings, qt.DeepEquals, []string{"invoice.finalized", "credit_note.created"})
		c.Assert(result.FinalizeZeroAmountInvoice, qt.Equals, true)
		c.Assert(result.EventsStore, qt.Equals, OrganizationEventsStoreClickhouse)
		c.Assert(result.BillingConfiguration.InvoiceGracePeriod, qt.Equals, 3)
		c.Assert(result.BillingConfiguration.InvoiceFooter, qt.Equals, "Thanks for your business")
		c.Assert(result.BillingConfiguration.DocumentLocale, qt.Equals, "en")
		c.Assert(result.Taxes, qt.HasLen, 1)
		c.Assert(result.Taxes[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Taxes[0].Name, qt.Equals, "VAT")
		c.Assert(result.Taxes[0].Code, qt.Equals, "vat_20")
		c.Assert(result.Taxes[0].Rate, qt.Equals, float32(20.0))
		c.Assert(result.Taxes[0].Description, qt.Equals, "Standard VAT")
		c.Assert(result.Taxes[0].AppliedToOrganization, qt.Equals, true)
		c.Assert(result.CreatedAt, qt.Equals, time.Date(2022, 4, 29, 8, 59, 51, 0, time.UTC))
	})

	t.Run("When the server returns an error response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/organizations").
			MatchJSONBody(`{
				"organization": {
					"email": "invalid-email",
					"billing_configuration": {}
				}
			}`).
			MockResponseWithCode(422, map[string]any{
				"status": 422,
				"error":  "Unprocessable Entity",
				"code":   "validation_errors",
				"error_details": map[string]any{
					"email": []string{"invalid_email_format"},
				},
			})
		defer server.Close()

		result, err := server.Client().Organization().Update(context.Background(), &OrganizationInput{
			Email: "invalid-email",
		})

		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 422)
		c.Assert(err.Message, qt.Equals, "Unprocessable Entity")
		c.Assert(err.ErrorCode, qt.Equals, "validation_errors")
		c.Assert(err.ErrorDetail, qt.IsNotNil)
		details, detailErr := err.ErrorDetail.Details()
		c.Assert(detailErr, qt.IsNil)
		c.Assert(details["email"], qt.DeepEquals, []string{"invalid_email_format"})
	})
}
