package lago_test

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
)

func TestPaymentListFilters(t *testing.T) {
	for _, customerScoped := range []bool{false, true} {
		name := "payments"
		if customerScoped {
			name = "customer payments"
		}
		t.Run(name, func(t *testing.T) {
			c := qt.New(t)
			page, perPage := 2, 5
			var from int64 = 0
			var to int64 = math.MaxInt64
			want := url.Values{
				"page": {"2"}, "per_page": {"5"}, "invoice_id": {"1a901a90-1a90-1a90-1a90-1a901a901a90"},
				"payment_status[]": {"succeeded", "failed"}, "payment_statuses[]": {"pending"},
				"amount_from": {"0"}, "amount_to": {"9223372036854775807"}, "receipt_number": {"Rcpt & +/#1"},
				"created_at_from": {"2026-09-01"}, "created_at_to": {"2026-09-07"},
				"payment_provider_type[]": {"stripe", "gocardless"}, "currency": {"EUR"},
				"invoice_number": {"LAG & +/#2"}, "payment_type[]": {"manual", "provider"},
				"payable_type[]": {"Invoice", "PaymentRequest"}, "search_term": {"pi_3 & +/#"},
			}
			path := "/api/v1/payments"
			if customerScoped {
				path = "/api/v1/customers/cust_1/payments"
			} else {
				want.Set("external_customer_id", "cust_1")
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Check(r.Method, qt.Equals, http.MethodGet)
				c.Check(r.URL.Path, qt.Equals, path)
				c.Check(r.URL.Query(), qt.DeepEquals, want)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"payments":[],"meta":{"current_page":2,"total_count":7}}`))
			}))
			defer server.Close()
			client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
			var result *PaymentResult
			var err *Error
			if customerScoped {
				result, err = client.Customer().GetPaymentList(context.Background(), "cust_1", &CustomerPaymentListInput{
					Page: &page, PerPage: &perPage, InvoiceID: want.Get("invoice_id"), PaymentStatus: want["payment_status[]"], PaymentStatuses: want["payment_statuses[]"],
					AmountFrom: &from, AmountTo: &to, ReceiptNumber: want.Get("receipt_number"), CreatedAtFrom: want.Get("created_at_from"), CreatedAtTo: want.Get("created_at_to"),
					PaymentProviderType: want["payment_provider_type[]"], Currency: Currency("EUR"), InvoiceNumber: want.Get("invoice_number"),
					PaymentType: want["payment_type[]"], PayableType: want["payable_type[]"], SearchTerm: want.Get("search_term"),
				})
			} else {
				result, err = client.Payment().GetList(context.Background(), &PaymentListInput{
					Page: &page, PerPage: &perPage, ExternalCustomerID: "cust_1", InvoiceID: want.Get("invoice_id"), PaymentStatus: want["payment_status[]"], PaymentStatuses: want["payment_statuses[]"],
					AmountFrom: &from, AmountTo: &to, ReceiptNumber: want.Get("receipt_number"), CreatedAtFrom: want.Get("created_at_from"), CreatedAtTo: want.Get("created_at_to"),
					PaymentProviderType: want["payment_provider_type[]"], Currency: Currency("EUR"), InvoiceNumber: want.Get("invoice_number"),
					PaymentType: want["payment_type[]"], PayableType: want["payable_type[]"], SearchTerm: want.Get("search_term"),
				})
			}
			c.Assert(err == nil, qt.IsTrue)
			c.Check(result.Meta.TotalCount, qt.Equals, 7)
		})
	}
}

func TestPaymentListWithoutFilters(t *testing.T) {
	c := qt.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Check(r.URL.RawQuery, qt.Equals, "")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payments":[],"meta":{"total_count":0}}`))
	}))
	defer server.Close()
	client := New().SetBaseURL(server.URL)
	for _, input := range []*PaymentListInput{nil, {}} {
		_, err := client.Payment().GetList(context.Background(), input)
		c.Assert(err == nil, qt.IsTrue)
	}
}

func TestPaymentListExistingJSONTags(t *testing.T) {
	c := qt.New(t)
	page := 2
	data, err := json.Marshal(&PaymentListInput{Page: &page, ExternalCustomerID: "cust_1", InvoiceID: "inv"})
	c.Assert(err == nil, qt.IsTrue)
	c.Check(string(data), qt.JSONEquals, map[string]string{"page": "2", "external_customer_id": "cust_1", "invoice_id": "inv"})
}
