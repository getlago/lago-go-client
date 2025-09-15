package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"

	lt "github.com/getlago/lago-go-client/testing"

	. "github.com/getlago/lago-go-client"
)

// Mock JSON response structure
var mockPaymentRequest = map[string]interface{}{
	"lago_id":         "89b6b61e-4dbc-4307-ac96-4abcfa9e3e2d",
	"email":           "gavin@overdue.test",
	"amount_cents":    19955,
	"amount_currency": "EUR",
	"payment_status":  "pending",
	"created_at":      "2024-06-30T10:59:51Z",
	"invoices": []map[string]interface{}{
		{
			"lago_id":                                 "f8e194df-5d90-4382-b146-c881d2c67f28",
			"sequential_id":                           15,
			"number":                                  "LAG-1234-001-002",
			"issuing_date":                            "2022-06-02",
			"payment_dispute_lost_at":                 "2022-04-29T08:59:51Z",
			"payment_due_date":                        "2022-06-02",
			"payment_overdue":                         true,
			"invoice_type":                            "one_off",
			"version_number":                          2,
			"status":                                  "finalized",
			"payment_status":                          "pending",
			"currency":                                "EUR",
			"net_payment_term":                        0,
			"fees_amount_cents":                       10000,
			"taxes_amount_cents":                      2000,
			"coupons_amount_cents":                    0,
			"payment_requests_amount_cents":           0,
			"sub_total_excluding_taxes_amount_cents":  10000,
			"sub_total_including_taxes_amount_cents":  12000,
			"prepaid_credit_amount_cents":             0,
			"total_amount_cents":                      12000,
			"progressive_billing_credit_amount_cents": 0,
			"file_url":                                "https://lago-files/invoice_002.pdf",
		},
	},
}

var mockPaymentRequestResponse = map[string]interface{}{
	"payment_request": mockPaymentRequest,
}

func assertPaymentRequestResponse(c *qt.C, result *PaymentRequest) {
	c.Assert(result.LagoID.String(), qt.Equals, "89b6b61e-4dbc-4307-ac96-4abcfa9e3e2d")
	c.Assert(result.Email, qt.Equals, "gavin@overdue.test")
	c.Assert(result.Invoices, qt.HasLen, 1)
	invoice := result.Invoices[0]
	c.Assert(invoice.LagoID.String(), qt.Equals, "f8e194df-5d90-4382-b146-c881d2c67f28")
}

func TestPaymentRequestRequest_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		paymentRequestUUID, _ := uuid.Parse("89b6b61e-4dbc-4307-ac96-4abcfa9e3e2d")

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.PaymentRequest().Get(context.Background(), paymentRequestUUID)
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/payment_requests/89b6b61e-4dbc-4307-ac96-4abcfa9e3e2d\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When get is called", func(t *testing.T) {
		c := qt.New(t)

		paymentRequestUUID, _ := uuid.Parse("89b6b61e-4dbc-4307-ac96-4abcfa9e3e2d")

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/payment_requests/89b6b61e-4dbc-4307-ac96-4abcfa9e3e2d").
			MockResponse(mockPaymentRequestResponse)
		defer server.Close()

		result, err := server.Client().PaymentRequest().Get(context.Background(), paymentRequestUUID)
		c.Assert(err == nil, qt.IsTrue)
		assertPaymentRequestResponse(c, result)
	})
}
