package lago_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"

	lt "github.com/getlago/lago-go-client/testing"

	. "github.com/getlago/lago-go-client"
)

// Mock JSON response structure
var mockCreditNote = map[string]interface{}{
	"lago_id":                                "1a901a90-1a90-1a90-1a90-1a901a901a90",
	"billing_entity_code":                    "acme_corp",
	"sequential_id":                          2,
	"number":                                 "LAG-1234-CN2",
	"lago_invoice_id":                        "1a901a90-1a90-1a90-1a90-1a901a901a90",
	"invoice_number":                         "LAG-1234",
	"issuing_date":                           "2022-12-06",
	"credit_status":                          "available",
	"refund_status":                          "pending",
	"reason":                                 "other",
	"description":                            "Free text",
	"currency":                               "EUR",
	"total_amount_cents":                     120,
	"taxes_amount_cents":                     20,
	"taxes_rate":                             20,
	"sub_total_excluding_taxes_amount_cents": 100,
	"balance_amount_cents":                   100,
	"credit_amount_cents":                    100,
	"refund_amount_cents":                    0,
	"coupons_adjustment_amount_cents":        20,
	"created_at":                             "2022-09-14T16:35:31Z",
	"updated_at":                             "2022-09-14T16:35:31Z",
	"file_url":                               "https://getlago.com/credit_note/file",
	"items": []map[string]interface{}{
		{
			"lago_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"amount_cents":    100,
			"amount_currency": "EUR",
			"fee": map[string]interface{}{
				"lago_id":                    "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_charge_id":             "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_charge_filter_id":      "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_invoice_id":            "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_true_up_fee_id":        "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_true_up_parent_fee_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_subscription_id":       "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"lago_customer_id":           "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"external_customer_id":       "external_id",
				"external_subscription_id":   "external_id",
				"invoice_display_name":       "Setup Fee (SF1)",
				"amount_cents":               100,
				"precise_amount":             "1.0001",
				"precise_total_amount":       "1.0212",
				"amount_currency":            "EUR",
				"taxes_amount_cents":         20,
				"taxes_precise_amount":       "0.20123",
				"taxes_rate":                 20,
				"units":                      "0.32",
				"precise_unit_amount":        "312.5",
				"total_amount_cents":         120,
				"total_amount_currency":      "EUR",
				"events_count":               23,
				"pay_in_advance":             true,
				"invoiceable":                true,
				"from_date":                  "2022-04-29T08:59:51Z",
				"to_date":                    "2022-05-29T08:59:51Z",
				"payment_status":             "pending",
				"created_at":                 "2022-08-24T14:58:59Z",
				"succeeded_at":               "2022-08-24T14:58:59Z",
				"failed_at":                  "2022-08-24T14:58:59Z",
				"refunded_at":                "2022-08-24T14:58:59Z",
				"event_transaction_id":       "transaction_1234567890",
				"amount_details": map[string]interface{}{
					"graduated_ranges": []map[string]interface{}{
						{
							"units":                  "10.0",
							"from_value":             0,
							"to_value":               10,
							"flat_unit_amount":       "1.0",
							"per_unit_amount":        "1.0",
							"per_unit_total_amount":  "10.0",
							"total_with_flat_amount": "11.0",
						},
					},
					"free_units":                      "10.0",
					"paid_units":                      "40.0",
					"per_package_size":                1000,
					"per_package_unit_amount":         "0.5",
					"units":                           "20.0",
					"free_events":                     10,
					"rate":                            "1.0",
					"per_unit_total_amount":           "10.0",
					"paid_events":                     20,
					"fixed_fee_unit_amount":           "1.0",
					"fixed_fee_total_amount":          "20.0",
					"min_max_adjustment_total_amount": "20.0",
				},
				"self_billed": false,
				"item": map[string]interface{}{
					"type":                        "subscription",
					"code":                        "startup",
					"name":                        "Startup",
					"invoice_display_name":        "Setup Fee (SF1)",
					"filter_invoice_display_name": "AWS eu-east-1",
					"filters": map[string][]string{
						"additionalProp1": {"string"},
						"additionalProp2": {"string"},
						"additionalProp3": {"string"},
					},
					"lago_item_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
					"item_type":    "Subscription",
					"grouped_by": map[string]interface{}{
						"additionalProp1": "string",
						"additionalProp2": "string",
						"additionalProp3": "string",
					},
				},
				"applied_taxes": []map[string]interface{}{
					{
						"lago_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
						"lago_tax_id":     "1a901a90-1a90-1a90-1a90-1a901a901a90",
						"tax_name":        "TVA",
						"tax_code":        "french_standard_vat",
						"tax_rate":        20,
						"tax_description": "French standard VAT",
						"amount_cents":    2000,
						"amount_currency": "USD",
						"created_at":      "2022-09-14T16:35:31Z",
						"lago_fee_id":     "1a901a90-1a90-1a90-1a90-1a901a901a90",
					},
				},
			},
		},
	},
	"applied_taxes": []map[string]interface{}{
		{
			"lago_id":             "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_tax_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"tax_name":            "TVA",
			"tax_code":            "french_standard_vat",
			"tax_rate":            20,
			"tax_description":     "French standard VAT",
			"amount_cents":        2000,
			"amount_currency":     "USD",
			"created_at":          "2022-09-14T16:35:31Z",
			"lago_credit_note_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"base_amount_cents":   100,
		},
	},
	"self_billed":   false,
	"error_details": []map[string]interface{}{},
}

var mockCreditNoteListResponse = map[string]interface{}{
	"credit_notes": []map[string]interface{}{
		mockCreditNote,
	},
	"meta": map[string]interface{}{
		"current_page": 2,
		"next_page":    3,
		"prev_page":    1,
		"total_pages":  4,
		"total_count":  70,
	},
}

var mockCreditNoteResponse = map[string]interface{}{
	"credit_note": mockCreditNote,
}

var mockMetadataResponse = map[string]interface{}{
	"metadata": map[string]interface{}{
		"foo": "bar",
		"baz": nil,
	},
}

var mockNullMetadataResponse = map[string]interface{}{
	"metadata": nil,
}

var mockCreditNoteEstimateResponse = map[string]interface{}{
	"estimated_credit_note": map[string]interface{}{
		"lago_invoice_id":                         "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"invoice_number":                          "LAG-1234",
		"currency":                                "EUR",
		"taxes_amount_cents":                      20,
		"taxes_rate":                              20,
		"sub_total_excluding_taxes_amount_cents":  100,
		"max_creditable_amount_cents":             100,
		"max_refundable_amount_cents":             0,
		"max_offsettable_amount_cents":            100,
		"coupons_adjustment_amount_cents":         20,
		"precise_coupons_adjustment_amount_cents": 20,
		"precise_taxes_amount_cents":              20,
		"items": []map[string]interface{}{
			{
				"amount_cents": 100,
				"lago_fee_id":  "1a901a90-1a90-1a90-1a90-1a901a901a90",
			},
		},
		"applied_taxes": []map[string]interface{}{
			{
				"lago_tax_id":       "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"tax_name":          "TVA",
				"tax_code":          "french_standard_vat",
				"tax_rate":          20,
				"tax_description":   "French standard VAT",
				"base_amount_cents": 100,
				"amount_cents":      2000,
				"amount_currency":   "USD",
			},
		},
	},
}

func assertCreditNoteGetListResponse(c *qt.C, result *CreditNoteResult) {
	c.Assert(result.CreditNotes, qt.HasLen, 1)

	creditNote := result.CreditNotes[0]
	assertCreditNoteResponse(c, &creditNote)
}

func assertCreditNoteResponse(c *qt.C, result *CreditNote) {
	c.Assert(result.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.SequentialID, qt.Equals, 2)
	c.Assert(result.BillingEntityCode, qt.Equals, "acme_corp")
	c.Assert(result.Number, qt.Equals, "LAG-1234-CN2")
	c.Assert(result.LagoInvoiceID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.InvoiceNumber, qt.Equals, "LAG-1234")
	c.Assert(result.Reason, qt.Equals, CreditNoteReason("other"))
	c.Assert(result.Description, qt.Equals, "Free text")
	c.Assert(result.SelfBilled, qt.Equals, false)
	c.Assert(result.CreditStatus, qt.Equals, CreditNoteCreditStatus("available"))
	c.Assert(result.RefundStatus, qt.Equals, CreditNoteRefundStatus("pending"))
	c.Assert(result.Currency, qt.Equals, Currency("EUR"))
	c.Assert(result.TotalAmountCents, qt.Equals, 120)
	c.Assert(result.CreditAmountCents, qt.Equals, 100)
	c.Assert(result.BalanceAmountCents, qt.Equals, 100)
	c.Assert(result.RefundAmountCents, qt.Equals, 0)
	c.Assert(result.TaxesAmountCents, qt.Equals, 20)
	c.Assert(result.TaxesRate, qt.Equals, 20.0)
	c.Assert(result.SubTotalExcludingTaxesAmountCents, qt.Equals, 100)
	c.Assert(result.CouponsAdjustmentAmountCents, qt.Equals, 20)
	c.Assert(result.FileURL, qt.Equals, "https://getlago.com/credit_note/file")
	c.Assert(result.IssuingDate, qt.Equals, "2022-12-06")
	c.Assert(result.CreatedAt.Format(time.RFC3339), qt.Equals, "2022-09-14T16:35:31Z")
	c.Assert(result.UpdatedAt.Format(time.RFC3339), qt.Equals, "2022-09-14T16:35:31Z")

	c.Assert(result.Items, qt.HasLen, 1)
	item := result.Items[0]
	c.Assert(item.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(item.AmountCents, qt.Equals, 100)
	c.Assert(item.AmountCurrency, qt.Equals, Currency("EUR"))
	c.Assert(item.Fee.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")

	c.Assert(result.AppliedTaxes, qt.HasLen, 1)
	appliedTax := result.AppliedTaxes[0]
	c.Assert(appliedTax.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedTax.LagoCreditNoteID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedTax.LagoTaxID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedTax.TaxName, qt.Equals, "TVA")
	c.Assert(appliedTax.TaxCode, qt.Equals, "french_standard_vat")
	c.Assert(appliedTax.TaxRate, qt.Equals, 20.0)
	c.Assert(appliedTax.TaxDescription, qt.Equals, "French standard VAT")
	c.Assert(appliedTax.AmountCents, qt.Equals, 2000)
	c.Assert(appliedTax.AmountCurrency, qt.Equals, Currency("USD"))
	c.Assert(appliedTax.BaseAmountCents, qt.Equals, 100)
	c.Assert(appliedTax.CreatedAt.Format(time.RFC3339), qt.Equals, "2022-09-14T16:35:31Z")

	c.Assert(result.ErrorDetails, qt.HasLen, 0)
}

func assertCreditNoteEstimateResponse(c *qt.C, result *EstimatedCreditNote) {
	c.Assert(result.LagoInvoiceID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.InvoiceNumber, qt.Equals, "LAG-1234")
	c.Assert(result.Currency, qt.Equals, Currency("EUR"))
	c.Assert(result.MaxCreditableAmountCents, qt.Equals, 100)
	c.Assert(result.MaxRefundableAmountCents, qt.Equals, 100)
	c.Assert(result.MaxOffsettableAmountCents, qt.Equals, 100)
	c.Assert(result.TaxesAmountCents, qt.Equals, 20)
	c.Assert(result.TaxesRate, qt.Equals, 20.0)
	c.Assert(result.SubTotalExcludingTaxesAmountCents, qt.Equals, 100)
	c.Assert(result.CouponsAdjustmentAmountCents, qt.Equals, 20)
	c.Assert(result.PreciseTaxesAmountCents, qt.Equals, 20.0)
	c.Assert(result.PreciseCouponsAdjustmentAmountCents, qt.Equals, 20.0)

	c.Assert(result.Items, qt.HasLen, 1)
	item := result.Items[0]
	c.Assert(item.LagoFeeID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(item.AmountCents, qt.Equals, 100)

	c.Assert(result.AppliedTaxes, qt.HasLen, 1)
	appliedTax := result.AppliedTaxes[0]
	c.Assert(appliedTax.LagoTaxID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedTax.TaxName, qt.Equals, "TVA")
	c.Assert(appliedTax.TaxCode, qt.Equals, "french_standard_vat")
	c.Assert(appliedTax.TaxRate, qt.Equals, 20.0)
	c.Assert(appliedTax.TaxDescription, qt.Equals, "French standard VAT")
	c.Assert(appliedTax.AmountCents, qt.Equals, 2000)
	c.Assert(appliedTax.AmountCurrency, qt.Equals, Currency("USD"))
	c.Assert(appliedTax.BaseAmountCents, qt.Equals, 100)
}

func TestCreditNoteRequest_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().Get(context.Background(), creditNoteUUID)
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When get is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(mockCreditNoteResponse)
		defer server.Close()

		result, err := server.Client().CreditNote().Get(context.Background(), creditNoteUUID)
		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteResponse(c, result)
	})
}

func TestCreditNoteRequest_Download(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().Download(context.Background(), creditNoteUUID)
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Post \"http://localhost:88888/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/download\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When Download is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/download").
			MockResponse(mockCreditNoteResponse)
		defer server.Close()

		result, err := server.Client().CreditNote().Download(context.Background(), creditNoteUUID)
		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteResponse(c, result)
	})

	t.Run("When Download returns an empty response", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/download").
			MockResponse(nil)
		defer server.Close()

		result, err := server.Client().CreditNote().Download(context.Background(), creditNoteUUID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result == nil, qt.IsTrue)
	})
}

func TestCreditNoteRequest_GetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().GetList(context.Background(), &CreditNoteListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/credit_notes\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/credit_notes").
			MockResponse(mockCreditNoteListResponse)
		defer server.Close()

		result, err := server.Client().CreditNote().GetList(context.Background(), &CreditNoteListInput{})
		// The method interface should return `error` and not `*Error` but that would break the API.
		// See https://go.dev/doc/faq#nil_error.

		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteGetListResponse(c, result)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/credit_notes").
			MatchQuery(map[string]interface{}{
				"per_page":             "10",
				"page":                 "1",
				"external_customer_id": "CUSTOMER_1",
				"issuing_date_from":    "2022-09-14T00:00:00Z",
				"issuing_date_to":      "2022-09-14T23:59:59Z",
				"amount_from":          "10",
				"amount_to":            "1000",
				"search_term":          "credit",
				"billing_entity_ids[]": []string{"1a901a90-1a90-1a90-1a90-1a901a901a90", "2a902a90-2a90-2a90-2a90-2a902a902a90"},
				"credit_status":        "consumed",
				"currency":             "EUR",
				"invoice_number":       "LAG_1234",
				"reason":               "order_change",
				"refund_status":        "pending",
				"self_billed":          "false",
			}).
			MockResponse(mockCreditNoteListResponse)
		defer server.Close()

		selfBilled := false
		entityUUID1, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")
		entityUUID2, _ := uuid.Parse("2a902a90-2a90-2a90-2a90-2a902a902a90")

		result, err := server.Client().CreditNote().GetList(context.Background(), &CreditNoteListInput{
			PerPage:            Ptr(10),
			Page:               Ptr(1),
			ExternalCustomerID: "CUSTOMER_1",
			IssuingDateFrom:    "2022-09-14T00:00:00Z",
			IssuingDateTo:      "2022-09-14T23:59:59Z",
			AmountFrom:         10,
			AmountTo:           1000,
			SearchTerm:         "credit",
			BillingEntityIDs:   []uuid.UUID{entityUUID1, entityUUID2},
			CreditStatus:       CreditNoteCreditStatusConsumed,
			Currency:           EUR,
			InvoiceNumber:      "LAG_1234",
			Reason:             CreditNoteReasonOrderChange,
			RefundStatus:       CreditNoteRefundStatusPending,
			SelfBilled:         &selfBilled,
		})

		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteGetListResponse(c, result)
	})
}

func TestCreditNoteRequest_Create(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().Create(context.Background(), &CreditNoteInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Post \"http://localhost:88888/api/v1/credit_notes\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When create is called", func(t *testing.T) {
		c := qt.New(t)

		server := lt.ServerWithAssertions(c, mockCreditNoteResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/credit_notes")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			// Parse JSON
			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			if creditNote, ok := requestData["credit_note"].(map[string]interface{}); ok {
				c.Assert(creditNote["invoice_id"], qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
				c.Assert(creditNote["reason"], qt.Equals, "order_cancellation")
				c.Assert(creditNote["description"], qt.Equals, "test description")
				c.Assert(creditNote["credit_amount_cents"], qt.Equals, float64(10))
				c.Assert(creditNote["refund_amount_cents"], qt.Equals, float64(20))

				c.Assert(creditNote["items"], qt.HasLen, 1)

				if items, ok := creditNote["items"].([]map[string]interface{}); ok {
					c.Assert(items[0]["fee_id"], qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
					c.Assert(items[0]["amount_cents"], qt.Equals, 100)
				}
			}
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")

		lagoID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")
		description := "test description"

		result, err := client.CreditNote().Create(context.Background(), &CreditNoteInput{
			LagoInvoiceID: lagoID,
			Reason:        CreditNoteReason("order_cancellation"),
			Description:   &description,
			Items: []CreditNoteItemInput{
				{
					LagoFeeID:   lagoID,
					AmountCents: 100,
				},
			},
			CreditAmountCents: 10,
			RefundAmountCents: 20,
		})

		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteResponse(c, result)
	})
}

func TestCreditNoteRequest_Update(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().Update(context.Background(), &CreditNoteUpdateInput{
			LagoID: creditNoteUUID,
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Put \"http://localhost:88888/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When update is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.ServerWithAssertions(c, mockCreditNoteResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "PUT")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			// Parse JSON
			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			if creditNote, ok := requestData["credit_note"].(map[string]interface{}); ok {
				c.Assert(creditNote["refund_status"], qt.Equals, "succeeded")
			}
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.CreditNote().Update(context.Background(), &CreditNoteUpdateInput{
			LagoID:       creditNoteUUID,
			RefundStatus: CreditNoteRefundStatus("succeeded"),
		})

		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteResponse(c, result)
	})
}

func TestCreditNoteRequest_Void(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().Void(context.Background(), creditNoteUUID)
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Put \"http://localhost:88888/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/void\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When void is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/void").
			MockResponse(mockCreditNoteResponse)
		defer server.Close()

		result, err := server.Client().CreditNote().Void(context.Background(), creditNoteUUID)
		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteResponse(c, result)
	})
}

func TestCreditNoteRequest_Estimate(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CreditNote().Estimate(context.Background(), &EstimateCreditNoteInput{
			LagoInvoiceID: uuid.New(),
			Items: []CreditNoteItemInput{
				{
					LagoFeeID:   uuid.New(),
					AmountCents: 10,
				},
			},
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Post \"http://localhost:88888/api/v1/credit_notes/estimate\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When estimate is called", func(t *testing.T) {
		c := qt.New(t)

		server := lt.ServerWithAssertions(c, mockCreditNoteEstimateResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/credit_notes/estimate")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			// Parse JSON
			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			if creditNote, ok := requestData["credit_note"].(map[string]interface{}); ok {
				c.Assert(creditNote["invoice_id"], qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
				c.Assert(creditNote["items"], qt.HasLen, 1)

				if items, ok := creditNote["items"].([]map[string]interface{}); ok {
					c.Assert(items[0]["fee_id"], qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
					c.Assert(items[0]["amount_cents"], qt.Equals, 100)
				}
			}
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")

		lagoID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")
		result, err := client.CreditNote().Estimate(context.Background(), &EstimateCreditNoteInput{
			LagoInvoiceID: lagoID,
			Items: []CreditNoteItemInput{
				{
					LagoFeeID:   lagoID,
					AmountCents: 10,
				},
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		assertCreditNoteEstimateResponse(c, result)
	})
}

func TestCreditNoteRequest_ReplaceMetadata(t *testing.T) {
	t.Run("When replace metadata is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.ServerWithAssertions(c, mockMetadataResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/metadata")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			c.Assert(requestData["metadata"], qt.IsNotNil)
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		bar := "bar"
		result, err := client.CreditNote().ReplaceMetadata(context.Background(), creditNoteUUID, map[string]*string{
			"foo": &bar,
			"baz": nil,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
		c.Assert(result["baz"], qt.IsNil)
	})
}

func TestCreditNoteRequest_MergeMetadata(t *testing.T) {
	t.Run("When merge metadata is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.ServerWithAssertions(c, mockMetadataResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "PATCH")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/metadata")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			c.Assert(requestData["metadata"], qt.IsNotNil)
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		qux := "qux"
		result, err := client.CreditNote().MergeMetadata(context.Background(), creditNoteUUID, map[string]*string{
			"foo": &qux,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
		c.Assert(result["baz"], qt.IsNil)
	})
}

func TestCreditNoteRequest_DeleteAllMetadata(t *testing.T) {
	t.Run("When delete all metadata is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/metadata").
			MockResponse(mockNullMetadataResponse)
		defer server.Close()

		result, err := server.Client().CreditNote().DeleteAllMetadata(context.Background(), creditNoteUUID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNil)
	})
}

func TestCreditNoteRequest_DeleteMetadataKey(t *testing.T) {
	t.Run("When delete metadata key is called", func(t *testing.T) {
		c := qt.New(t)

		creditNoteUUID, _ := uuid.Parse("1a901a90-1a90-1a90-1a90-1a901a901a90")

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/credit_notes/1a901a90-1a90-1a90-1a90-1a901a901a90/metadata/foo").
			MockResponse(mockMetadataResponse)
		defer server.Close()

		result, err := server.Client().CreditNote().DeleteMetadataKey(context.Background(), creditNoteUUID, "foo")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
	})
}
