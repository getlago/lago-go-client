package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockPaymentMethodListResponse = map[string]any{
	"payment_methods": []map[string]any{
		{
			"lago_id":               "pm-lago-id-1",
			"is_default":            true,
			"payment_provider_code": "stripe-eu-1",
			"payment_provider_name": "Stripe EU",
			"payment_provider_type": "stripe",
			"provider_method_id":    "pm_stripe_123",
			"created_at":            "2024-01-15T10:00:00Z",
		},
		{
			"lago_id":               "pm-lago-id-2",
			"is_default":            false,
			"payment_provider_code": "adyen-us-1",
			"payment_provider_name": "Adyen US",
			"payment_provider_type": "adyen",
			"provider_method_id":    "pm_adyen_456",
			"created_at":            "2024-02-20T12:00:00Z",
		},
	},
	"meta": map[string]any{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  2,
	},
}

var mockPaymentMethodResponse = map[string]any{
	"payment_method": map[string]any{
		"lago_id":               "pm-lago-id-1",
		"is_default":            true,
		"payment_provider_code": "stripe-eu-1",
		"payment_provider_name": "Stripe EU",
		"payment_provider_type": "stripe",
		"provider_method_id":    "pm_stripe_123",
		"created_at":            "2024-01-15T10:00:00Z",
	},
}

func TestCustomerRequest_GetPaymentMethodList(t *testing.T) {
	t.Run("When the customer is not found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/INVALID_CUSTOMER/payment_methods").
			MockResponseWithCode(404, map[string]any{
				"status": 404,
				"error":  "Not Found",
				"code":   "customer_not_found",
			})
		defer server.Close()

		result, err := server.Client().Customer().GetPaymentMethodList(context.Background(), "INVALID_CUSTOMER", &CustomerPaymentMethodListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 404)
		c.Assert(err.Message, qt.Equals, "Not Found")
		c.Assert(err.ErrorCode, qt.Equals, "customer_not_found")
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_methods").
			MatchQuery("").
			MockResponse(mockPaymentMethodListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetPaymentMethodList(context.Background(), "CUSTOMER_1", &CustomerPaymentMethodListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentMethods, qt.HasLen, 2)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 2)

		pm := result.PaymentMethods[0]
		c.Assert(pm.LagoID, qt.Equals, "pm-lago-id-1")
		c.Assert(pm.IsDefault, qt.Equals, true)
		c.Assert(pm.PaymentProviderCode, qt.Equals, "stripe-eu-1")
		c.Assert(pm.PaymentProviderName, qt.Equals, "Stripe EU")
		c.Assert(pm.PaymentProviderType, qt.Equals, "stripe")
		c.Assert(pm.ProviderMethodID, qt.Equals, "pm_stripe_123")
		c.Assert(pm.CreatedAt, qt.Equals, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))

		pm2 := result.PaymentMethods[1]
		c.Assert(pm2.LagoID, qt.Equals, "pm-lago-id-2")
		c.Assert(pm2.IsDefault, qt.Equals, false)
		c.Assert(pm2.PaymentProviderCode, qt.Equals, "adyen-us-1")
		c.Assert(pm2.PaymentProviderName, qt.Equals, "Adyen US")
		c.Assert(pm2.PaymentProviderType, qt.Equals, "adyen")
		c.Assert(pm2.ProviderMethodID, qt.Equals, "pm_adyen_456")
		c.Assert(pm2.CreatedAt, qt.Equals, time.Date(2024, 2, 20, 12, 0, 0, 0, time.UTC))
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_methods").
			MatchQuery(map[string]string{
				"per_page": "10",
				"page":     "1",
			}).
			MockResponse(mockPaymentMethodListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetPaymentMethodList(context.Background(), "CUSTOMER_1", &CustomerPaymentMethodListInput{
			PerPage: Ptr(10),
			Page:    Ptr(1),
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentMethods, qt.HasLen, 2)
	})
}

func TestCustomerRequest_DestroyPaymentMethod(t *testing.T) {
	t.Run("When the payment method is not found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_methods/invalid-pm-id").
			MockResponseWithCode(404, map[string]any{
				"status": 404,
				"error":  "Not Found",
				"code":   "payment_method_not_found",
			})
		defer server.Close()

		result, err := server.Client().Customer().DestroyPaymentMethod(context.Background(), "CUSTOMER_1", "invalid-pm-id")
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 404)
		c.Assert(err.Message, qt.Equals, "Not Found")
		c.Assert(err.ErrorCode, qt.Equals, "payment_method_not_found")
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_methods/pm-lago-id-1").
			MockResponse(mockPaymentMethodResponse)
		defer server.Close()

		result, err := server.Client().Customer().DestroyPaymentMethod(context.Background(), "CUSTOMER_1", "pm-lago-id-1")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, "pm-lago-id-1")
		c.Assert(result.IsDefault, qt.Equals, true)
		c.Assert(result.PaymentProviderCode, qt.Equals, "stripe-eu-1")
		c.Assert(result.PaymentProviderName, qt.Equals, "Stripe EU")
		c.Assert(result.PaymentProviderType, qt.Equals, "stripe")
		c.Assert(result.ProviderMethodID, qt.Equals, "pm_stripe_123")
		c.Assert(result.CreatedAt, qt.Equals, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))
	})
}

func TestCustomerRequest_SetPaymentMethodAsDefault(t *testing.T) {
	t.Run("When the payment method is not found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_methods/invalid-pm-id/set_as_default").
			MockResponseWithCode(404, map[string]any{
				"status": 404,
				"error":  "Not Found",
				"code":   "payment_method_not_found",
			})
		defer server.Close()

		result, err := server.Client().Customer().SetPaymentMethodAsDefault(context.Background(), "CUSTOMER_1", "invalid-pm-id")
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 404)
		c.Assert(err.Message, qt.Equals, "Not Found")
		c.Assert(err.ErrorCode, qt.Equals, "payment_method_not_found")
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_methods/pm-lago-id-1/set_as_default").
			MockResponse(mockPaymentMethodResponse)
		defer server.Close()

		result, err := server.Client().Customer().SetPaymentMethodAsDefault(context.Background(), "CUSTOMER_1", "pm-lago-id-1")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, "pm-lago-id-1")
		c.Assert(result.IsDefault, qt.Equals, true)
		c.Assert(result.PaymentProviderCode, qt.Equals, "stripe-eu-1")
		c.Assert(result.PaymentProviderName, qt.Equals, "Stripe EU")
		c.Assert(result.PaymentProviderType, qt.Equals, "stripe")
		c.Assert(result.ProviderMethodID, qt.Equals, "pm_stripe_123")
		c.Assert(result.CreatedAt, qt.Equals, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))
	})
}
