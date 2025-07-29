package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	lt "github.com/getlago/lago-go-client/testing"

	. "github.com/getlago/lago-go-client"
)

// Mock JSON response structure
var mockResponse = map[string]interface{}{
	"applied_coupons": []map[string]interface{}{
		{
			"lago_id":                      "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_coupon_id":               "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"coupon_code":                  "APPLIED_COUPON",
			"coupon_name":                  "Startup Deal",
			"lago_customer_id":             "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"external_customer_id":         "CUSTOMER_1",
			"status":                       "active",
			"amount_cents":                 2000,
			"amount_cents_remaining":       50,
			"amount_currency":              "EUR",
			"percentage_rate":              nil,
			"frequency":                    "recurring",
			"frequency_duration":           3,
			"frequency_duration_remaining": 1,
			"expiration_at":                "2022-04-29T08:59:51Z",
			"created_at":                   "2022-04-29T08:59:51Z",
			"terminated_at":                "2022-04-29T08:59:51Z",
			"credits": []map[string]interface{}{
				{
					"lago_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
					"amount_cents":    1200,
					"amount_currency": "EUR",
					"before_taxes":    false,
					"item": map[string]interface{}{
						"lago_item_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
						"type":         "coupon",
						"code":         "startup_deal",
						"name":         "Startup Deal",
					},
					"invoice": map[string]interface{}{
						"lago_id":        "1a901a90-1a90-1a90-1a90-1a901a901a90",
						"payment_status": "succeeded",
					},
				},
			},
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

func assertAppliedCoupontGetListResponse(c *qt.C, result *AppliedCouponResult) {
	c.Assert(result.AppliedCoupons, qt.HasLen, 1)
	appliedCoupon := result.AppliedCoupons[0]
	c.Assert(appliedCoupon.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedCoupon.LagoCouponID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedCoupon.CouponCode, qt.Equals, "APPLIED_COUPON")
	c.Assert(appliedCoupon.CouponName, qt.Equals, "Startup Deal")
	c.Assert(appliedCoupon.LagoCustomerID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(appliedCoupon.ExternalCustomerID, qt.Equals, "CUSTOMER_1")
	c.Assert(appliedCoupon.Status, qt.Equals, AppliedCouponStatusActive)
	c.Assert(appliedCoupon.AmountCents, qt.Equals, 2000)
	c.Assert(appliedCoupon.AmountCentsRemaining, qt.Equals, 50)
	c.Assert(appliedCoupon.AmountCurrency, qt.Equals, EUR)
	c.Assert(appliedCoupon.PercentageRate, qt.Equals, 0.0)
	c.Assert(appliedCoupon.Frequency, qt.Equals, CouponFrequencyRecurring)
	c.Assert(appliedCoupon.FrequencyDuration, qt.Equals, 3)
	c.Assert(appliedCoupon.FrequencyDurationRemaining, qt.Equals, 1)
	c.Assert(appliedCoupon.ExpirationAt.Format(time.RFC3339), qt.Equals, "2022-04-29T08:59:51Z")
	c.Assert(appliedCoupon.CreatedAt.Format(time.RFC3339), qt.Equals, "2022-04-29T08:59:51Z")
	c.Assert(appliedCoupon.TerminatedAt.Format(time.RFC3339), qt.Equals, "2022-04-29T08:59:51Z")
	c.Assert(appliedCoupon.Credits, qt.HasLen, 1)
	credit := appliedCoupon.Credits[0]
	c.Assert(credit.AmountCents, qt.Equals, 1200)
	c.Assert(credit.AmountCurrency, qt.Equals, EUR)
	c.Assert(credit.BeforeTaxes, qt.Equals, false)
	c.Assert(credit.Item.Type, qt.Equals, InvoiceCreditItemType("coupon"))
	c.Assert(credit.Item.Code, qt.Equals, "startup_deal")
	c.Assert(credit.Item.Name, qt.Equals, "Startup Deal")
	c.Assert(credit.Invoice.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
	c.Assert(credit.Invoice.PaymentStatus, qt.Equals, InvoicePaymentStatusSucceeded)
	c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
	c.Assert(result.Meta.NextPage, qt.Equals, 0)
	c.Assert(result.Meta.PrevPage, qt.Equals, 0)
	c.Assert(result.Meta.TotalPages, qt.Equals, 1)
	c.Assert(result.Meta.TotalCount, qt.Equals, 1)
}

func TestAppliedCouponGetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.AppliedCoupon().GetList(context.Background(), &AppliedCouponListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/applied_coupons\": dial tcp: address 88888: invalid port"}`)

	})

	t.Run("When no parameter is provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).MatchMethod("GET").MatchPath("/api/v1/applied_coupons").MatchQuery("").MockResponse(mockResponse)
		defer server.Close()

		result, err := server.Client().AppliedCoupon().GetList(context.Background(), &AppliedCouponListInput{})
		// The method interface should return `error` and not `*Error` but that would break the API.
		// See https://go.dev/doc/faq#nil_error.
		c.Assert(err == nil, qt.IsTrue)
		assertAppliedCoupontGetListResponse(c, result)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/applied_coupons").
			MatchQuery(map[string]interface{}{
				"external_customer_id": "CUSTOMER_1",
				"status":               "active",
				"per_page":             "10",
				"page":                 "1",
				"coupon_code[]":        []string{"BLACK_FRIDAY", "CYBER_MONDAY"},
			}).
			MockResponse(mockResponse)
		defer server.Close()

		result, err := server.Client().AppliedCoupon().GetList(context.Background(), &AppliedCouponListInput{
			ExternalCustomerID: "CUSTOMER_1",
			Status:             "active",
			PerPage:            10,
			Page:               1,
			CouponCode:         []string{"BLACK_FRIDAY", "CYBER_MONDAY"},
		})
		c.Assert(err == nil, qt.IsTrue)
		assertAppliedCoupontGetListResponse(c, result)
	})
}
