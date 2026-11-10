package lago_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/getlago/lago-go-client"
)

func TestCustomerExternalIDPath(t *testing.T) {
	tests := []struct {
		name, suffix, method string
		call                 func(*CustomerRequest, string) *Error
	}{
		{"Get", "", "GET", func(cr *CustomerRequest, id string) *Error { _, err := cr.Get(context.Background(), id); return err }},
		{"Delete", "", "DELETE", func(cr *CustomerRequest, id string) *Error { _, err := cr.Delete(context.Background(), id); return err }},
		{"PortalUrl", "portal_url", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.PortalUrl(context.Background(), id)
			return err
		}},
		{"CheckoutUrl", "checkout_url", "POST", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.CheckoutUrl(context.Background(), id)
			return err
		}},
		{"CurrentUsage", "current_usage", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.CurrentUsage(context.Background(), id, &CustomerUsageInput{})
			return err
		}},
		{"PastUsage", "past_usage", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.PastUsage(context.Background(), id, &CustomerPastUsageInput{})
			return err
		}},
		{"ProjectedUsage", "projected_usage", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.ProjectedUsage(context.Background(), id, &CustomerUsageInput{})
			return err
		}},
		{"GetWalletList", "wallets", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetWalletList(context.Background(), id, &CustomerWalletListInput{})
			return err
		}},
		{"GetInvoiceList", "invoices", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetInvoiceList(context.Background(), id, &CustomerInvoiceListInput{})
			return err
		}},
		{"GetCreditNoteList", "credit_notes", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetCreditNoteList(context.Background(), id, &CustomerCreditNoteListInput{})
			return err
		}},
		{"GetPaymentList", "payments", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetPaymentList(context.Background(), id, &CustomerPaymentListInput{})
			return err
		}},
		{"GetPaymentRequestList", "payment_requests", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetPaymentRequestList(context.Background(), id, &CustomerPaymentRequestListInput{})
			return err
		}},
		{"GetAppliedCouponList", "applied_coupons", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetAppliedCouponList(context.Background(), id, &CustomerAppliedCouponListInput{})
			return err
		}},
		{"GetSubscriptionList", "subscriptions", "GET", func(cr *CustomerRequest, id string) *Error {
			_, err := cr.GetSubscriptionList(context.Background(), id, &CustomerSubscriptionListInput{})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, id := range []string{"tenant#customer", "tenant?customer", "tenant/customer", "tenant%customer"} {
				t.Run(id, func(t *testing.T) {
					called := false
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						called = true
						want := "/api/v1/customers/" + url.PathEscape(id)
						if tt.suffix != "" {
							want += "/" + tt.suffix
						}
						if r.URL.EscapedPath() != want {
							t.Errorf("path = %q, want %q", r.URL.EscapedPath(), want)
						}
						if r.Method != tt.method {
							t.Errorf("method = %q, want %q", r.Method, tt.method)
						}
						query := r.URL.Query()
						if tt.name == "CurrentUsage" || tt.name == "ProjectedUsage" {
							if query.Get("apply_taxes") != "false" {
								t.Errorf("apply_taxes = %q, want false", query.Get("apply_taxes"))
							}
							query.Del("apply_taxes")
							query.Del("external_subscription_id")
						}
						if tt.name == "PastUsage" {
							query.Del("external_subscription_id")
						}
						if len(query) != 0 {
							t.Errorf("unexpected query parameters: %v", query)
						}
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte("{}"))
					}))
					defer server.Close()
					if err := tt.call(New().SetBaseURL(server.URL).Customer(), id); err != nil {
						t.Fatal(err)
					}
					if !called {
						t.Fatal("server received no request")
					}
				})
			}
		})
	}
}
