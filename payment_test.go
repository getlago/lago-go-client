package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

func mockPayment(providerFields map[string]any) map[string]any {
	payment := map[string]any{
		"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_customer_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"lago_payable_id":      "3c903c90-3c90-3c90-3c90-3c903c903c90",
		"payable_type":         "Invoice",
		"amount_cents":         1200,
		"amount_currency":      "EUR",
		"status":               "processing",
		"payment_status":       "succeeded",
		"type":                 "provider",
		"reference":            "REF-123456",
		"external_payment_id":  "pi_123",
		"created_at":           "2022-04-29T08:59:51Z",
		"invoice_ids":          []string{"3c903c90-3c90-3c90-3c90-3c903c903c90"},
		"invoice_numbers":      []string{"LAG-1234-001-002"},
		"external_customer_id": "CUSTOMER_1",
	}
	for key, value := range providerFields {
		payment[key] = value
	}

	return payment
}

var providerPaymentFields = map[string]any{
	"payment_provider_code": "stripe_prod",
	"payment_provider_type": "stripe",
	"provider_payment_id":   "pi_123",
	"provider_customer_id":  "cus_123",
}

var manualPaymentFields = map[string]any{
	"payment_provider_code": nil,
	"payment_provider_type": nil,
	"provider_payment_id":   nil,
	"provider_customer_id":  nil,
}

func assertPaymentFields(c *qt.C, payment *Payment) {
	c.Assert(payment.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(payment.LagoCustomerID.String(), qt.Equals, "2b902b90-2b90-2b90-2b90-2b902b902b90")
	c.Assert(payment.LagoPayableID.String(), qt.Equals, "3c903c90-3c90-3c90-3c90-3c903c903c90")
	c.Assert(payment.PayableType, qt.Equals, "Invoice")
	c.Assert(payment.Status, qt.Equals, "processing")
	c.Assert(payment.PaymentStatus, qt.Equals, "succeeded")
}

func TestPaymentRequest_Get(t *testing.T) {
	t.Run("When the payment comes from a payment provider", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/payments/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(map[string]any{"payment": mockPayment(providerPaymentFields)})
		defer server.Close()

		result, err := server.Client().Payment().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		assertPaymentFields(c, result)
		c.Assert(*result.PaymentProviderCode, qt.Equals, "stripe_prod")
		c.Assert(*result.PaymentProviderType, qt.Equals, "stripe")
		c.Assert(*result.ProviderPaymentID, qt.Equals, "pi_123")
		c.Assert(*result.ProviderCustomerID, qt.Equals, "cus_123")
	})

	t.Run("When the payment is manual", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/payments/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(map[string]any{"payment": mockPayment(manualPaymentFields)})
		defer server.Close()

		result, err := server.Client().Payment().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		assertPaymentFields(c, result)
		c.Assert(result.PaymentProviderCode, qt.IsNil)
		c.Assert(result.PaymentProviderType, qt.IsNil)
		c.Assert(result.ProviderPaymentID, qt.IsNil)
		c.Assert(result.ProviderCustomerID, qt.IsNil)
	})
}

func TestPaymentRequest_GetList(t *testing.T) {
	t.Run("When payments are found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/payments").
			MockResponse(map[string]any{
				"payments": []map[string]any{mockPayment(providerPaymentFields)},
				"meta": map[string]any{
					"current_page": 1,
					"next_page":    0,
					"prev_page":    0,
					"total_pages":  1,
					"total_count":  1,
				},
			})
		defer server.Close()

		result, err := server.Client().Payment().GetList(context.Background(), &PaymentListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Payments, qt.HasLen, 1)
		assertPaymentFields(c, &result.Payments[0])
		c.Assert(*result.Payments[0].PaymentProviderCode, qt.Equals, "stripe_prod")
		c.Assert(*result.Payments[0].PaymentProviderType, qt.Equals, "stripe")
		c.Assert(*result.Payments[0].ProviderPaymentID, qt.Equals, "pi_123")
		c.Assert(*result.Payments[0].ProviderCustomerID, qt.Equals, "cus_123")
	})
}
