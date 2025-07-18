package lago

import (
	"context"
	"net/http"
	"testing"
	"time"

	lt "github.com/getlago/lago-go-client/testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"
)

var mockInvoice = map[string]any{
	"lago_id":                                 "1a901a90-1a90-1a90-1a90-1a901a901a90",
	"billing_entity_code":                     "acme_corp",
	"sequential_id":                           2,
	"number":                                  "LAG-1234-001-002",
	"issuing_date":                            "2022-04-30",
	"payment_dispute_lost_at":                 "2022-09-14T16:35:31Z",
	"payment_due_date":                        "2022-04-30",
	"payment_overdue":                         true,
	"net_payment_term":                        30,
	"invoice_type":                            "subscription",
	"status":                                  "finalized",
	"payment_status":                          "succeeded",
	"currency":                                "EUR",
	"fees_amount_cents":                       100,
	"coupons_amount_cents":                    10,
	"credit_notes_amount_cents":               10,
	"sub_total_excluding_taxes_amount_cents":  100,
	"taxes_amount_cents":                      20,
	"sub_total_including_taxes_amount_cents":  120,
	"prepaid_credit_amount_cents":             0,
	"progressive_billing_credit_amount_cents": 0,
	"total_amount_cents":                      100,
	"version_number":                          3,
	"self_billed":                             false,
	"file_url":                                "https://getlago.com/invoice/file",
	"created_at":                              "2022-04-29T08:59:51Z",
	"updated_at":                              "2022-04-29T08:59:51Z",
	"customer": map[string]any{
		"lago_id":                      "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"sequential_id":                1,
		"slug":                         "LAG-1234-001",
		"external_id":                  "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
		"billing_entity_code":          "acme_corp",
		"address_line1":                "5230 Penfield Ave",
		"address_line2":                nil,
		"applicable_timezone":          "America/Los_Angeles",
		"city":                         "Woodland Hills",
		"country":                      "US",
		"currency":                     "USD",
		"email":                        "dinesh@piedpiper.test",
		"legal_name":                   "Coleman-Blair",
		"legal_number":                 "49-008-2965",
		"logo_url":                     "http://hooli.com/logo.png",
		"name":                         "Gavin Belson",
		"firstname":                    "Gavin",
		"lastname":                     "Belson",
		"account_type":                 "customer",
		"customer_type":                "company",
		"phone":                        "1-171-883-3711 x245",
		"state":                        "CA",
		"tax_identification_number":    "EU123456789",
		"timezone":                     "America/Los_Angeles",
		"url":                          "http://hooli.com",
		"zipcode":                      "91364",
		"net_payment_term":             30,
		"created_at":                   "2022-04-29T08:59:51Z",
		"updated_at":                   "2022-04-29T08:59:51Z",
		"finalize_zero_amount_invoice": "inherit",
		"skip_invoice_custom_sections": false,
		"billing_configuration": map[string]any{
			"invoice_grace_period":  3,
			"payment_provider":      "stripe",
			"payment_provider_code": "stripe-eu-1",
			"provider_customer_id":  "cus_12345",
			"sync":                  true,
			"sync_with_provider":    true,
			"document_locale":       "fr",
			"provider_payment_methods": []string{
				"card",
				"sepa_debit",
				"us_bank_account",
				"bacs_debit",
				"link",
				"boleto",
				"crypto",
				"customer_balance",
			},
			"integration_customers": []map[string]any{},
		},
		"shipping_address": map[string]any{
			"address_line1": "5230 Penfield Ave",
			"address_line2": nil,
			"city":          "Woodland Hills",
			"country":       "US",
			"state":         "CA",
			"zipcode":       "91364",
		},
		"metadata": []map[string]any{
			{
				"lago_id":            "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"key":                "Purchase Order",
				"value":              "123456789",
				"display_in_invoice": true,
				"created_at":         "2022-04-29T08:59:51Z",
			},
		},
		"integration_customers": []map[string]any{
			{
				"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"type":                 "netsuite",
				"integration_code":     "netsuite-eu-1",
				"external_customer_id": "cus_12345",
				"sync_with_provider":   true,
				"subsidiary_id":        "2",
			},
		},
	},
	"metadata": []map[string]any{
		{
			"lago_id":    "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"key":        "digital_ref_id",
			"value":      "INV-0123456-98765",
			"created_at": "2022-04-29T08:59:51Z",
		},
	},
	"applied_taxes": []map[string]any{
		{
			"lago_id":           "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_tax_id":       "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"tax_name":          "TVA",
			"tax_code":          "french_standard_vat",
			"tax_rate":          20,
			"tax_description":   "French standard VAT",
			"amount_cents":      2000,
			"amount_currency":   "USD",
			"created_at":        "2022-09-14T16:35:31Z",
			"lago_invoice_id":   "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"fees_amount_cents": 20000,
		},
	},
}

var mockInvoiceListResponse = map[string]any{
	"invoices": []map[string]interface{}{
		mockInvoice,
	},
	"meta": map[string]interface{}{
		"current_page": 2,
		"next_page":    3,
		"prev_page":    1,
		"total_pages":  4,
		"total_count":  70,
	},
}

// Mock JSON response structure
var mockInvoicePaymentUrlResponse = map[string]any{
	"invoice_payment_details": map[string]any{
		"lago_customer_id":     "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_invoice_id":      "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"external_customer_id": "cust_1234",
		"payment_provider":     "stripe",
		"payment_url":          "https://example.com/payment",
	},
}

func assertInvoiceGetListResponse(c *qt.C, result *InvoiceResult) {
	c.Assert(result.Invoices, qt.HasLen, 1)

	invoice := result.Invoices[0]
	assertInvoiceResponse(c, &invoice)
}

func assertInvoiceResponse(c *qt.C, result *Invoice) {
	c.Assert(result.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.SequentialID, qt.Equals, 2)
	c.Assert(result.BillingEntityCode, qt.Equals, "acme_corp")
	c.Assert(result.Number, qt.Equals, "LAG-1234-001-002")
	c.Assert(result.IssuingDate, qt.Equals, "2022-04-30")
	c.Assert(result.PaymentDisputeLostAt.Format(time.RFC3339), qt.Equals, "2022-09-14T16:35:31Z")
	c.Assert(result.PaymentDueDate, qt.Equals, "2022-04-30")
	c.Assert(result.PaymentOverdue, qt.Equals, true)
	c.Assert(result.InvoiceType, qt.Equals, InvoiceType("subscription"))
	c.Assert(result.Status, qt.Equals, InvoiceStatus("finalized"))
	c.Assert(result.PaymentStatus, qt.Equals, InvoicePaymentStatus("succeeded"))
	c.Assert(result.Currency, qt.Equals, Currency("EUR"))
	c.Assert(result.FeesAmountCents, qt.Equals, 100)
	c.Assert(result.TaxesAmountCents, qt.Equals, 20)
	c.Assert(result.CouponsAmountCents, qt.Equals, 10)
	c.Assert(result.CreditNotesAmountCents, qt.Equals, 10)
	c.Assert(result.SubTotalExcludingTaxesAmountCents, qt.Equals, 100)
	c.Assert(result.SubTotalIncludingTaxesAmountCents, qt.Equals, 120)
	c.Assert(result.TotalAmountCents, qt.Equals, 100)
	c.Assert(result.TotalDueAmountCents, qt.Equals, 0)
	c.Assert(result.PrepaidCreditAmountCents, qt.Equals, 0)
	c.Assert(result.ProgressiveBillingCreditAmountCents, qt.Equals, 0)
	c.Assert(result.NetPaymentTerm, qt.Equals, 30)
	c.Assert(result.FileURL, qt.Equals, "https://getlago.com/invoice/file")
	c.Assert(result.VersionNumber, qt.Equals, 3)

	c.Assert(result.Metadata, qt.HasLen, 1)
	metadata := result.Metadata[0]
	c.Assert(metadata.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(metadata.Key, qt.Equals, "digital_ref_id")
	c.Assert(metadata.Value, qt.Equals, "INV-0123456-98765")
	c.Assert(metadata.CreatedAt.Format(time.RFC3339), qt.Equals, "2022-04-29T08:59:51Z")

	customer := result.Customer
	c.Assert(customer.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(customer.ExternalID, qt.Equals, "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba")

	c.Assert(result.BillingPeriods, qt.HasLen, 0)
	c.Assert(result.Subscriptions, qt.HasLen, 0)
}

func assertPaymentUrlResponse(c *qt.C, result *InvoicePaymentDetails) {
	c.Assert(result.LagoCustomerID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.LagoInvoiceID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.ExternalCustomerID, qt.Equals, "cust_1234")
	c.Assert(result.PaymentProvider, qt.Equals, "stripe")
	c.Assert(result.PaymentUrl, qt.Equals, "https://example.com/payment")
}

func TestInvoiceGetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Invoice().GetList(context.Background(), &InvoiceListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/invoices\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.HandlerFunc(c, mockInvoiceListResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "GET")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/invoices")
			c.Assert(r.URL.Query().Encode(), qt.Equals, "")

			query := r.URL.Query()
			c.Assert(query.Get("per_page"), qt.Equals, "")
			c.Assert(query.Get("page"), qt.Equals, "")
			c.Assert(query.Get("payment_overdue"), qt.Equals, "")
			c.Assert(query.Get("partially_paid"), qt.Equals, "")
			c.Assert(query.Get("self_billed"), qt.Equals, "")
			c.Assert(query.Get("payment_dispute_lost"), qt.Equals, "")
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.Invoice().GetList(context.Background(), &InvoiceListInput{})
		// The method interface should return `error` and not `*Error` but that would break the API.
		// See https://go.dev/doc/faq#nil_error.

		c.Assert(err == nil, qt.IsTrue)
		assertInvoiceGetListResponse(c, result)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.HandlerFunc(c, mockInvoiceListResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "GET")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/invoices")

			query := r.URL.Query()
			c.Assert(query.Get("per_page"), qt.Equals, "10")
			c.Assert(query.Get("page"), qt.Equals, "1")
			c.Assert(query.Get("customer_external_id"), qt.Equals, "CUSTOMER_1")
			c.Assert(query.Get("invoice_type"), qt.Equals, "subscription")
			c.Assert(query.Get("status"), qt.Equals, "finalized")
			c.Assert(query.Get("payment_status"), qt.Equals, "succeeded")
			c.Assert(query.Get("issuing_date_from"), qt.Equals, "2022-09-14T00:00:00Z")
			c.Assert(query.Get("issuing_date_to"), qt.Equals, "2022-09-14T23:59:59Z")
			c.Assert(query.Get("amount_from"), qt.Equals, "10")
			c.Assert(query.Get("amount_to"), qt.Equals, "1000")
			c.Assert(query.Get("search_term"), qt.Equals, "credit")
			c.Assert(query["billing_entity_ids[]"], qt.DeepEquals, []string{"1a901a90-1a90-1a90-1a90-1a901a901a90", "1a901a90-1a90-1a90-1a90-1a901a901a91"})
			c.Assert(query.Get("currency"), qt.Equals, "EUR")
			c.Assert(query.Get("payment_overdue"), qt.Equals, "true")
			c.Assert(query.Get("partially_paid"), qt.Equals, "true")
			c.Assert(query.Get("self_billed"), qt.Equals, "false")
			c.Assert(query.Get("payment_dispute_lost"), qt.Equals, "true")
			c.Assert(query.Get("metadata[key1]"), qt.Equals, "10")
			c.Assert(query.Get("metadata[key2]"), qt.Equals, "value2")
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		// selfBilled := false
		entityUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")
		entityUUID2, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a91")

		result, err := client.Invoice().GetList(context.Background(), &InvoiceListInput{
			PerPage:            Ptr(10),
			Page:               Ptr(1),
			IssuingDateFrom:    "2022-09-14T00:00:00Z",
			IssuingDateTo:      "2022-09-14T23:59:59Z",
			ExternalCustomerID: "CUSTOMER_1",
			InvoiceType:        SubscriptionInvoiceType,
			Status:             InvoiceStatusFinalized,
			PaymentStatus:      InvoicePaymentStatusSucceeded,
			AmountFrom:         Ptr(10),
			AmountTo:           Ptr(1000),
			SearchTerm:         "credit",
			BillingEntityIDs:   []uuid.UUID{entityUUID, entityUUID2},
			Currency:           EUR,
			PaymentOverdue:     Ptr(true),
			PartiallyPaid:      Ptr(true),
			SelfBilled:         Ptr(false),
			PaymentDisputeLost: Ptr(true),
			Metadata:           &InvoiceListInputMetadata{"key1": 10, "key2": "value2"},
		})

		c.Assert(err == nil, qt.IsTrue)
		assertInvoiceGetListResponse(c, result)
	})
}

func TestPaymentUrl(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Invoice().PaymentUrl(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Post \"http://localhost:88888/api/v1/invoices/1a901a90-1a90-1a90-1a90-1a901a901a90/payment_url\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("With an invoiceID in the request", func(t *testing.T) {
		c := qt.New(t)

		server := lt.HandlerFunc(c, mockInvoicePaymentUrlResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/invoices/1a901a90-1a90-1a90-1a90-1a901a901a90/payment_url")
		})

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.Invoice().PaymentUrl(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")

		c.Assert(err == nil, qt.IsTrue)
		assertPaymentUrlResponse(c, result)
	})
}
