package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

const quoteVersionID = "4d234d23-4d23-4d23-4d23-4d234d234d23"

var QuoteVersionGetMockResponse = map[string]interface{}{
	"quote_version": map[string]interface{}{
		"lago_id":              quoteVersionID,
		"lago_quote_id":        "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
		"version":              1,
		"status":               "draft",
		"currency":             "EUR",
		"billing_entity_code":  "acme_corp",
		"void_reason":          nil,
		"approved_at":          nil,
		"voided_at":            nil,
		"created_at":           "2026-04-29T08:59:51Z",
		"updated_at":           "2026-04-29T08:59:51Z",
		"content":              "<h1>Quote QT-2026-0001</h1>",
		"billing_items": map[string]interface{}{
			"plans": []map[string]interface{}{
				{
					"id":      "7a567a56-7a56-7a56-7a56-7a567a567a56",
					"localId": "b5c1e2a4-4e1e-4a7f-9f0e-9c1a0c7e1f2b",
					"type":    "plan",
					"payload": map[string]interface{}{
						"code":                   "premium_plan",
						"subscriptionExternalId": "sub_1234567890",
						"endDate":                "2027-01-01T00:00:00Z",
					},
					"overrides": map[string]interface{}{
						"amountCents":    50000,
						"amountCurrency": "EUR",
					},
				},
			},
			"coupons": []map[string]interface{}{
				{
					"id":      "8b678b67-8b67-8b67-8b67-8b678b678b67",
					"localId": "c6d2f3b5-5f2f-5b8a-af1f-ad2b1d8f2a3c",
					"type":    "coupon",
					"payload": map[string]interface{}{
						"code":        "welcome_offer",
						"type":        "fixed_amount",
						"amountCents": 10000,
					},
				},
			},
			"walletCredits": []map[string]interface{}{
				{
					"localId": "d7e3a4c6-6a3a-6c9b-ba2a-be3c2e9a3b4d",
					"type":    "wallet_credit",
					"payload": map[string]interface{}{
						"paidCredits":    "100.0",
						"grantedCredits": "20.0",
						"rateAmount":     "1.0",
					},
				},
			},
		},
	},
}

var QuoteVersionGetListMockResponse = map[string]interface{}{
	"quote_versions": []map[string]interface{}{
		{
			"lago_id":              quoteVersionID,
			"lago_quote_id":        "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"version":              2,
			"status":               "draft",
			"currency":             "EUR",
			"billing_entity_code":  "acme_corp",
			"void_reason":          nil,
			"approved_at":          nil,
			"voided_at":            nil,
			"created_at":           "2026-04-30T08:59:51Z",
			"updated_at":           "2026-04-30T08:59:51Z",
		},
		{
			"lago_id":              "9c789c78-9c78-9c78-9c78-9c789c789c78",
			"lago_quote_id":        "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"version":              1,
			"status":               "voided",
			"currency":             "EUR",
			"billing_entity_code":  "acme_corp",
			"void_reason":          "superseded",
			"approved_at":          nil,
			"voided_at":            "2026-04-30T08:59:51Z",
			"created_at":           "2026-04-29T08:59:51Z",
			"updated_at":           "2026-04-30T08:59:51Z",
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  2,
	},
}

func TestQuoteVersionRequest_Get(t *testing.T) {
	t.Run("When query for a specific quote version", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/quote_versions/" + quoteVersionID).
			MockResponse(QuoteVersionGetMockResponse)
		defer server.Close()

		result, err := server.Client().QuoteVersion().Get(context.Background(), quoteVersionID)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, quoteVersionID)
		c.Assert(result.Version, qt.Equals, 1)
		c.Assert(result.Status, qt.Equals, QuoteVersionStatusDraft)
		c.Assert(result.Currency, qt.Equals, EUR)
		c.Assert(result.BillingEntityCode, qt.Equals, "acme_corp")
		c.Assert(string(result.VoidReason), qt.Equals, "")
		c.Assert(result.ApprovedAt, qt.IsNil)
		c.Assert(result.VoidedAt, qt.IsNil)
		c.Assert(result.Content, qt.Equals, "<h1>Quote QT-2026-0001</h1>")

		c.Assert(result.BillingItems, qt.IsNotNil)
		c.Assert(result.BillingItems.AddOns, qt.HasLen, 0)
		c.Assert(result.BillingItems.Plans, qt.HasLen, 1)
		c.Assert(result.BillingItems.Plans[0].Type, qt.Equals, QuoteBillingItemTypePlan)
		c.Assert(result.BillingItems.Plans[0].LagoID.String(), qt.Equals, "7a567a56-7a56-7a56-7a56-7a567a567a56")
		c.Assert(result.BillingItems.Plans[0].Payload["code"], qt.Equals, "premium_plan")
		c.Assert(result.BillingItems.Plans[0].Overrides["amountCurrency"], qt.Equals, "EUR")
		c.Assert(result.BillingItems.Coupons, qt.HasLen, 1)
		c.Assert(result.BillingItems.Coupons[0].Type, qt.Equals, QuoteBillingItemTypeCoupon)
		c.Assert(result.BillingItems.WalletCredits, qt.HasLen, 1)
		c.Assert(result.BillingItems.WalletCredits[0].Type, qt.Equals, QuoteBillingItemTypeWalletCredit)
		c.Assert(result.BillingItems.WalletCredits[0].Payload["paidCredits"], qt.Equals, "100.0")
	})
}

func TestQuoteVersionRequest_Approve(t *testing.T) {
	t.Run("When approving without an expiration", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/quote_versions/" + quoteVersionID + "/approve").
			MockResponse(QuoteVersionGetMockResponse)
		defer server.Close()

		result, err := server.Client().QuoteVersion().Approve(context.Background(), quoteVersionID, nil)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, quoteVersionID)
	})

	t.Run("When approving with an expiration", func(t *testing.T) {
		c := qt.New(t)

		expiresAt, parseErr := time.Parse(time.RFC3339, "2026-06-30T23:59:59Z")
		c.Assert(parseErr, qt.IsNil)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/quote_versions/" + quoteVersionID + "/approve").
			MatchJSONBody(map[string]interface{}{"expires_at": "2026-06-30T23:59:59Z"}).
			MockResponse(QuoteVersionGetMockResponse)
		defer server.Close()

		result, err := server.Client().QuoteVersion().Approve(
			context.Background(),
			quoteVersionID,
			&QuoteVersionApproveInput{ExpiresAt: &expiresAt},
		)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
	})
}

func TestQuoteVersionRequest_Void(t *testing.T) {
	t.Run("When voiding a draft version", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/quote_versions/" + quoteVersionID + "/void").
			MockResponse(QuoteVersionGetMockResponse)
		defer server.Close()

		result, err := server.Client().QuoteVersion().Void(context.Background(), quoteVersionID)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, quoteVersionID)
	})
}

func TestQuoteVersionRequest_Clone(t *testing.T) {
	t.Run("When cloning a version", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/quote_versions/" + quoteVersionID + "/clone").
			MockResponse(QuoteVersionGetMockResponse)
		defer server.Close()

		result, err := server.Client().QuoteVersion().Clone(context.Background(), quoteVersionID)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, quoteVersionID)
	})
}
