package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockAttributedUsageResponse = map[string]any{
	"attributed_usage": map[string]any{
		"lago_subscription_id":     "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"external_subscription_id": "sub_1234",
		"group_by":                 "team",
		"basis":                    "amount",
		"from_datetime":            "2026-09-01T00:00:00Z",
		"to_datetime":              "2026-09-30T23:59:59Z",
		"currency":                 "EUR",
		"rows": []map[string]any{
			{
				"value":                "eng",
				"rank":                 1,
				"amount_cents":         75001,
				"precise_amount_cents": "75000.5",
				"events_count":         3,
				"charges_usage": []map[string]any{
					{
						"lago_charge_id":                     "2b902b90-2b90-2b90-2b90-2b902b902b90",
						"charge_code":                        "tokens_charge",
						"billable_metric_code":               "tokens",
						"lago_charge_filter_id":              "3c903c90-3c90-3c90-3c90-3c903c903c90",
						"charge_filter_values":               map[string]any{"model": []string{"opus"}},
						"charge_filter_invoice_display_name": "Opus",
						"units":                              "1000.0",
						"amount_cents":                       50001,
						"precise_amount_cents":               "50000.5",
						"events_count":                       2,
					},
				},
			},
		},
		"unattributed": map[string]any{
			"amount_cents":         nil,
			"precise_amount_cents": nil,
			"events_count":         0,
			"charges_usage":        []map[string]any{},
		},
		"totals": map[string]any{
			"amount_cents":         75001,
			"precise_amount_cents": "75000.5",
			"events_count":         3,
			"charges_usage":        []map[string]any{},
		},
	},
	"meta": map[string]any{
		"current_page": 1,
		"next_page":    nil,
		"prev_page":    nil,
		"total_pages":  1,
		"total_count":  1,
	},
}

func TestCustomerRequest_AttributedUsage(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().AttributedUsage(context.Background(), "cust_1234", &AttributedUsageInput{
			ExternalSubscriptionID: "sub_1234",
			GroupBy:                "team",
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When a filter or the charge codes are empty", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().AttributedUsage(context.Background(), "cust_1234", &AttributedUsageInput{
			ExternalSubscriptionID: "sub_1234",
			GroupBy:                "team",
			Filters:                map[string][]string{"model": {}},
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Err, qt.ErrorMatches, `attributed usage filter "model" requires at least one value`)

		result, err = client.Customer().AttributedUsage(context.Background(), "cust_1234", &AttributedUsageInput{
			ExternalSubscriptionID: "sub_1234",
			GroupBy:                "team",
			ChargeCodes:            []string{},
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Err, qt.ErrorMatches, "attributed usage charge codes require at least one value")
	})

	t.Run("When all parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/cust_1234/attributed_usage").
			MatchQuery(map[string]interface{}{
				"external_subscription_id": "sub_1234",
				"group_by":                 "team",
				"filters[department][]":    "rnd",
				"filters[model][]":         []string{"opus", "haiku"},
				"from_datetime":            "2026-09-01T00:00:00Z",
				"to_datetime":              "2026-09-30T23:59:59Z",
				"charge_codes[]":           []string{"tokens_charge", "seats_charge"},
				"split_charge_code":        "tokens_charge",
				"search_term":              "eng",
				"basis":                    "amount",
				"page":                     "2",
				"per_page":                 "50",
			}).
			MockResponse(mockAttributedUsageResponse)
		defer server.Close()

		page, perPage := 2, 50
		result, err := server.Client().Customer().AttributedUsage(context.Background(), "cust_1234", &AttributedUsageInput{
			ExternalSubscriptionID: "sub_1234",
			GroupBy:                "team",
			Filters: map[string][]string{
				"department": {"rnd"},
				"model":      {"opus", "haiku"},
			},
			FromDatetime:    "2026-09-01T00:00:00Z",
			ToDatetime:      "2026-09-30T23:59:59Z",
			ChargeCodes:     []string{"tokens_charge", "seats_charge"},
			SplitChargeCode: "tokens_charge",
			SearchTerm:      "eng",
			Basis:           AttributedUsageBasisAmount,
			Page:            &page,
			PerPage:         &perPage,
		})
		c.Assert(err == nil, qt.IsTrue)

		usage := result.AttributedUsage
		c.Assert(usage.ExternalSubscriptionID, qt.Equals, "sub_1234")
		c.Assert(usage.Basis, qt.Equals, AttributedUsageBasisAmount)
		c.Assert(usage.Currency, qt.Equals, EUR)
		c.Assert(usage.Rows, qt.HasLen, 1)
		c.Assert(usage.Rows[0].Value, qt.Equals, "eng")
		c.Assert(usage.Rows[0].Rank, qt.Equals, 1)
		c.Assert(*usage.Rows[0].PreciseAmountCents, qt.Equals, "75000.5")
		c.Assert(usage.Rows[0].ChargesUsage, qt.HasLen, 1)

		chargeUsage := usage.Rows[0].ChargesUsage[0]
		c.Assert(chargeUsage.ChargeCode, qt.Equals, "tokens_charge")
		c.Assert(chargeUsage.LagoChargeFilterID.String(), qt.Equals, "3c903c90-3c90-3c90-3c90-3c903c903c90")
		c.Assert(chargeUsage.ChargeFilterValues, qt.DeepEquals, map[string][]string{"model": {"opus"}})
		c.Assert(*chargeUsage.ChargeFilterInvoiceDisplayName, qt.Equals, "Opus")
		c.Assert(chargeUsage.Units, qt.Equals, "1000.0")
		c.Assert(*chargeUsage.AmountCents, qt.Equals, 50001)

		c.Assert(usage.Unattributed.AmountCents, qt.IsNil)
		c.Assert(usage.Unattributed.PreciseAmountCents, qt.IsNil)
		c.Assert(usage.Unattributed.EventsCount, qt.Equals, 0)
		c.Assert(*usage.Totals.AmountCents, qt.Equals, 75001)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})
}
