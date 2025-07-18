package lago

import (
	"context"
	"net/http"
	"testing"

	qt "github.com/frankban/quicktest"

	lt "github.com/getlago/lago-go-client/testing"
)

// Mock JSON response structure
var mockSubscriptionResponse = `{
  "subscription": {
    "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
    "external_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
    "lago_customer_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
    "external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
    "billing_time": "anniversary",
    "name": "Repository A",
    "plan_code": "premium",
    "status": "terminated",
    "created_at": "2022-08-08T00:00:00Z",
    "canceled_at": "2022-09-14T16:35:31Z",
    "started_at": "2022-08-08T00:00:00Z",
    "ending_at": "2022-10-08T00:00:00Z",
    "subscription_at": "2022-08-08T00:00:00Z",
    "terminated_at": "2022-09-14T16:35:31Z",
    "previous_plan_code": null,
    "next_plan_code": null,
    "downgrade_plan_date": "2022-04-30",
    "trial_ended_at": "2022-08-08T00:00:00Z",
    "on_termination_credit_note": "skip",
    "current_billing_period_started_at": "2022-08-08T00:00:00Z",
    "current_billing_period_ending_at": "2022-09-08T00:00:00Z",
    "plan": {
      "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
      "name": "Startup",
      "invoice_display_name": "Startup plan",
      "created_at": "2023-06-27T19:43:42Z",
      "code": "startup",
      "interval": "monthly",
      "description": "",
      "amount_cents": 10000,
      "amount_currency": "USD",
      "trial_period": 5,
      "pay_in_advance": true,
      "bill_charges_monthly": null,
      "minimum_commitment": {
        "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
        "plan_code": "premium",
        "amount_cents": 100000,
        "invoice_display_name": "Minimum Commitment (C1)",
        "interval": "monthly",
        "created_at": "2022-04-29T08:59:51Z",
        "updated_at": "2022-04-29T08:59:51Z",
        "taxes": [
          {
            "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
            "name": "TVA",
            "code": "french_standard_vat",
            "description": "French standard VAT",
            "rate": 20,
            "applied_to_organization": true,
            "created_at": "2023-07-06T14:35:58Z"
          }
        ]
      },
      "charges": [
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a91",
          "lago_billable_metric_id": "1a901a90-1a90-1a90-1a90-1a901a901a91",
          "billable_metric_code": "requests",
          "created_at": "2023-06-27T19:43:42Z",
          "charge_model": "package",
          "invoiceable": true,
          "invoice_display_name": "Setup",
          "pay_in_advance": false,
          "regroup_paid_fees": null,
          "prorated": false,
          "min_amount_cents": 3000,
          "properties": {
            "amount": "30",
            "free_units": 100,
            "package_size": 1000
          },
          "filters": []
        },
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a92",
          "lago_billable_metric_id": "1a901a90-1a90-1a90-1a90-1a901a901a92",
          "billable_metric_code": "cpu",
          "created_at": "2023-06-27T19:43:42Z",
          "charge_model": "graduated",
          "invoiceable": true,
          "invoice_display_name": "Setup",
          "pay_in_advance": false,
          "prorated": false,
          "min_amount_cents": 0,
          "properties": {
            "graduated_ranges": [
              {
                "from_value": 0,
                "to_value": 10,
                "flat_amount": "10",
                "per_unit_amount": "0.5"
              },
              {
                "from_value": 11,
                "to_value": null,
                "flat_amount": "0",
                "per_unit_amount": "0.4"
              }
            ]
          },
          "filters": []
        },
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a93",
          "lago_billable_metric_id": "1a901a90-1a90-1a90-1a90-1a901a901a93",
          "billable_metric_code": "seats",
          "created_at": "2023-06-27T19:43:42Z",
          "charge_model": "standard",
          "invoiceable": true,
          "invoice_display_name": "Setup",
          "pay_in_advance": true,
          "prorated": false,
          "min_amount_cents": 0,
          "properties": {},
          "filters": [
            {
              "invoice_display_name": "Europe",
              "properties": {
                "amount": "10"
              },
              "values": {
                "region": [
                  "Europe"
                ]
              }
            },
            {
              "invoice_display_name": "USA",
              "properties": {
                "amount": "5"
              },
              "values": {
                "region": [
                  "USA"
                ]
              }
            },
            {
              "invoice_display_name": "Africa",
              "properties": {
                "amount": "8"
              },
              "values": {
                "region": [
                  "Africa"
                ]
              }
            }
          ]
        },
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a94",
          "lago_billable_metric_id": "1a901a90-1a90-1a90-1a90-1a901a901a94",
          "billable_metric_code": "storage",
          "created_at": "2023-06-27T19:43:42Z",
          "charge_model": "volume",
          "invoiceable": true,
          "invoice_display_name": "Setup",
          "pay_in_advance": false,
          "prorated": false,
          "min_amount_cents": 0,
          "properties": {
            "volume_ranges": [
              {
                "from_value": 0,
                "to_value": 100,
                "flat_amount": "0",
                "per_unit_amount": "0"
              },
              {
                "from_value": 101,
                "to_value": null,
                "flat_amount": "0",
                "per_unit_amount": "0.5"
              }
            ]
          },
          "filters": []
        },
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a95",
          "lago_billable_metric_id": "1a901a90-1a90-1a90-1a90-1a901a901a95",
          "billable_metric_code": "payments",
          "created_at": "2023-06-27T19:43:42Z",
          "charge_model": "percentage",
          "invoiceable": false,
          "invoice_display_name": "Setup",
          "pay_in_advance": true,
          "regroup_paid_fees": "invoice",
          "prorated": false,
          "min_amount_cents": 0,
          "properties": {
            "rate": "1",
            "fixed_amount": "0.5",
            "free_units_per_events": 5,
            "free_units_per_total_aggregation": "500"
          },
          "filters": []
        }
      ],
      "taxes": [
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
          "name": "TVA",
          "code": "french_standard_vat",
          "description": "French standard VAT",
          "rate": 20,
          "applied_to_organization": true,
          "created_at": "2023-07-06T14:35:58Z"
        }
      ],
      "usage_thresholds": [
        {
          "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
          "threshold_display_name": "Threshold 1",
          "amount_cents": 10000,
          "recurring": true,
          "created_at": "2023-06-27T19:43:42Z",
          "updated_at": "2023-06-27T19:43:42Z"
        }
      ]
    }
  }
}`

func assertSubscriptionTerminateResponse(c *qt.C, subscription *Subscription) {
	c.Assert(subscription.OnTerminationCreditNote, qt.Equals, OnTerminationCreditNoteSkip)
}

func TestSubscriptionTerminate(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Subscription().Terminate(context.Background(), SubscriptionTerminateInput{
			ExternalID: "1a901a90-1a90-1a90-1a90-1a901a901a90",
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Delete \"http://localhost:88888/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameter is provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.HandlerFunc(c, mockSubscriptionResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "DELETE")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90")
			c.Assert(r.URL.Query().Encode(), qt.Equals, "")
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		subscription, err := client.Subscription().Terminate(context.Background(), SubscriptionTerminateInput{
			ExternalID: "1a901a90-1a90-1a90-1a90-1a901a901a90",
		})
		// The method interface should return `error` and not `*Error` but that would break the API.
		// See https://go.dev/doc/faq#nil_error.
		c.Assert(err == nil, qt.IsTrue)
		assertSubscriptionTerminateResponse(c, subscription)
	})

	t.Run("When providing the on_termination_credit_note parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.HandlerFunc(c, mockSubscriptionResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "DELETE")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90")
			c.Assert(r.URL.Query().Encode(), qt.Equals, "on_termination_credit_note=skip")
		})

		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		subscription, err := client.Subscription().Terminate(context.Background(), SubscriptionTerminateInput{
			ExternalID:              "1a901a90-1a90-1a90-1a90-1a901a901a90",
			OnTerminationCreditNote: OnTerminationCreditNoteSkip,
		})
		// The method interface should return `error` and not `*Error` but that would break the API.
		// See https://go.dev/doc/faq#nil_error.
		c.Assert(err == nil, qt.IsTrue)
		assertSubscriptionTerminateResponse(c, subscription)
	})
}
