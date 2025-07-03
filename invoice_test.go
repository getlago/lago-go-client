package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	qt "github.com/frankban/quicktest"
)

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

func paymentUrlHandlerFunc(c *qt.C, assertRequestFunc func(*qt.C, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestFunc(c, r)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockInvoicePaymentUrlResponse)
	}))
}

func assertPaymentUrlResponse(c *qt.C, result *InvoicePaymentDetails) {
	c.Assert(result.LagoCustomerID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.LagoInvoiceID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(result.ExternalCustomerID, qt.Equals, "cust_1234")
	c.Assert(result.PaymentProvider, qt.Equals, "stripe")
	c.Assert(result.PaymentUrl, qt.Equals, "https://example.com/payment")
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

		server := paymentUrlHandlerFunc(c, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/invoices/1a901a90-1a90-1a90-1a90-1a901a901a90/payment_url")
		})

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		result, err := client.Invoice().PaymentUrl(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")

		fmt.Println(err)

		c.Assert(err == nil, qt.IsTrue)
		fmt.Println(result.ExternalCustomerID)
		assertPaymentUrlResponse(c, result)
	})
}
