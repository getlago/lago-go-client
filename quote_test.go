package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
	"github.com/google/uuid"
)

const quoteOwnerID = "5e345e34-5e34-5e34-5e34-5e345e345e34"

var quoteCurrentVersionMock = map[string]interface{}{
	"lago_id":              "4d234d23-4d23-4d23-4d23-4d234d234d23",
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
}

var QuoteGetMockResponse = map[string]interface{}{
	"quote": map[string]interface{}{
		"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"number":               "QT-2026-0001",
		"order_type":           "subscription_creation",
		"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
		"lago_subscription_id": nil,
		"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
		"created_at":           "2026-04-29T08:59:51Z",
		"updated_at":           "2026-04-29T08:59:51Z",
		"current_version":      quoteCurrentVersionMock,
		"owners": []map[string]interface{}{
			{
				"lago_id": "5e345e34-5e34-5e34-5e34-5e345e345e34",
				"email":   "sales@getlago.com",
			},
		},
	},
}

var QuoteGetListMockResponse = map[string]interface{}{
	"quotes": []map[string]interface{}{
		{
			"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"number":               "QT-2026-0001",
			"order_type":           "subscription_creation",
			"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
			"lago_subscription_id": nil,
			"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"created_at":           "2026-04-29T08:59:51Z",
			"updated_at":           "2026-04-29T08:59:51Z",
			"current_version":      quoteCurrentVersionMock,
		},
		{
			"lago_id":              "6f456f45-6f45-6f45-6f45-6f456f456f45",
			"number":               "QT-2026-0002",
			"order_type":           "one_off",
			"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
			"lago_subscription_id": nil,
			"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"created_at":           "2026-04-30T08:59:51Z",
			"updated_at":           "2026-04-30T08:59:51Z",
			"current_version":      nil,
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

func TestQuoteRequest_Get(t *testing.T) {
	t.Run("When query for a specific quote", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/quotes/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(QuoteGetMockResponse)
		defer server.Close()

		result, err := server.Client().Quote().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(result.Number, qt.Equals, "QT-2026-0001")
		c.Assert(result.OrderType, qt.Equals, QuoteOrderTypeSubscriptionCreation)
		c.Assert(result.LagoSubscriptionID, qt.IsNil)
		c.Assert(result.CreatedAt.Format(time.RFC3339), qt.Equals, "2026-04-29T08:59:51Z")
		c.Assert(result.CurrentVersion, qt.IsNotNil)
		c.Assert(result.CurrentVersion.Version, qt.Equals, 1)
		c.Assert(result.CurrentVersion.Status, qt.Equals, QuoteVersionStatusDraft)
		c.Assert(result.Owners, qt.HasLen, 1)
		c.Assert(result.Owners[0].Email, qt.Equals, "sales@getlago.com")
	})
}

func TestQuoteRequest_GetList(t *testing.T) {
	t.Run("When query for all quotes", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/quotes").
			MockResponse(QuoteGetListMockResponse)
		defer server.Close()

		result, err := server.Client().Quote().GetList(context.Background(), &QuoteListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Quotes, qt.HasLen, 2)
		c.Assert(result.Quotes[0].Number, qt.Equals, "QT-2026-0001")
		c.Assert(result.Quotes[1].OrderType, qt.Equals, QuoteOrderTypeOneOff)
		c.Assert(result.Quotes[1].CurrentVersion, qt.IsNil)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
	})

	t.Run("When filtering quotes", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/quotes").
			MatchQuery(map[string]interface{}{
				"status[]":               []string{"draft", "approved"},
				"order_type[]":           []string{"one_off"},
				"number[]":               []string{"QT-2026-0001"},
				"owner_id[]":             []string{quoteOwnerID},
				"external_customer_id[]": []string{"ext_customer_1"},
				"from_date":              "2026-01-01",
				"to_date":                "2026-12-31",
			}).
			MockResponse(QuoteGetListMockResponse)
		defer server.Close()

		result, err := server.Client().Quote().GetList(context.Background(), &QuoteListInput{
			Status:             []QuoteVersionStatus{QuoteVersionStatusDraft, QuoteVersionStatusApproved},
			OrderType:          []QuoteOrderType{QuoteOrderTypeOneOff},
			Number:             []string{"QT-2026-0001"},
			OwnerIDs:           []uuid.UUID{uuid.MustParse(quoteOwnerID)},
			ExternalCustomerID: []string{"ext_customer_1"},
			FromDate:           "2026-01-01",
			ToDate:             "2026-12-31",
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Quotes, qt.HasLen, 2)
	})
}

func TestQuoteRequest_GetVersionList(t *testing.T) {
	t.Run("When query for the versions of a quote", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/quotes/1a901a90-1a90-1a90-1a90-1a901a901a90/versions").
			MockResponse(QuoteVersionGetListMockResponse)
		defer server.Close()

		result, err := server.Client().Quote().GetVersionList(
			context.Background(),
			"1a901a90-1a90-1a90-1a90-1a901a901a90",
			&QuoteVersionListInput{},
		)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.QuoteVersions, qt.HasLen, 2)
		c.Assert(result.QuoteVersions[0].Version, qt.Equals, 2)
		c.Assert(result.QuoteVersions[0].Status, qt.Equals, QuoteVersionStatusDraft)
		c.Assert(result.QuoteVersions[1].Status, qt.Equals, QuoteVersionStatusVoided)
		c.Assert(result.QuoteVersions[1].VoidReason, qt.Equals, QuoteVersionVoidReasonSuperseded)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
	})
}
