package lago_test

import (
	"context"
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	. "github.com/getlago/lago-go-client"
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
    "plan_amount_cents": 10000,
    "plan_amount_currency": "USD",
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
    "on_termination_invoice": "skip",
    "payment_method": {
      "payment_method_type": "card",
      "payment_method_id": "pm_123456"
    },
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
	c.Assert(subscription.OnTerminationInvoice, qt.Equals, OnTerminationInvoiceSkip)
}

func terminateSubscription(c *qt.C, server *lt.MockServer, input SubscriptionTerminateInput) *Subscription {
	subscription, err := server.Client().Subscription().Terminate(context.Background(), input)
	c.Assert(err == nil, qt.IsTrue)
	return subscription
}

func TestSubscriptionRequest_Terminate(t *testing.T) {
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

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MatchQuery("").
			MockResponse(mockSubscriptionResponse)
		defer server.Close()

		subscription := terminateSubscription(c, server, SubscriptionTerminateInput{
			ExternalID: "1a901a90-1a90-1a90-1a90-1a901a901a90",
		})

		assertSubscriptionTerminateResponse(c, subscription)
	})

	for _, onTerminationCreditNote := range []OnTerminationCreditNote{
		OnTerminationCreditNoteCredit,
		OnTerminationCreditNoteRefund,
		OnTerminationCreditNoteOffset,
		OnTerminationCreditNoteSkip,
	} {
		title := fmt.Sprintf("When providing the on_termination_credit_note=%s parameter", onTerminationCreditNote)
		t.Run(title, func(t *testing.T) {
			c := qt.New(t)

			server := lt.NewMockServer(c).
				MatchMethod("DELETE").
				MatchPath("/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90").
				MatchQuery(map[string]string{"on_termination_credit_note": string(onTerminationCreditNote)}).
				MockResponse(mockSubscriptionResponse)

			defer server.Close()

			input := SubscriptionTerminateInput{
				ExternalID:              "1a901a90-1a90-1a90-1a90-1a901a901a90",
				OnTerminationCreditNote: onTerminationCreditNote,
			}
			subscription := terminateSubscription(c, server, input)

			assertSubscriptionTerminateResponse(c, subscription)
		})
	}

	for _, onTerminationInvoice := range []OnTerminationInvoice{
		OnTerminationInvoiceGenerate,
		OnTerminationInvoiceSkip,
	} {
		title := fmt.Sprintf("When providing the on_termination_invoice=%s parameter", onTerminationInvoice)
		t.Run(title, func(t *testing.T) {
			c := qt.New(t)

			server := lt.NewMockServer(c).
				MatchMethod("DELETE").
				MatchPath("/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90").
				MatchQuery(map[string]string{"on_termination_invoice": string(onTerminationInvoice)}).
				MockResponse(mockSubscriptionResponse)

			defer server.Close()

			input := SubscriptionTerminateInput{
				ExternalID:           "1a901a90-1a90-1a90-1a90-1a901a901a90",
				OnTerminationInvoice: onTerminationInvoice,
			}
			subscription := terminateSubscription(c, server, input)

			assertSubscriptionTerminateResponse(c, subscription)
		})
	}
}

func TestSubscriptionRequest_CreateWithPaymentMethod(t *testing.T) {
	t.Run("When creating a subscription with a not found payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/subscriptions").
			MatchJSONBody(`{
				"subscription": {
					"external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
					"plan_code": "premium",
					"external_id": "my_sub_1",
					"name": "Repository A",
					"billing_time": "anniversary",
					"payment_method": {
						"payment_method_type": "provider",
						"payment_method_id": "pm_not_found"
					}
				}
			}`).
			MockResponseWithCode(404, map[string]any{
				"status": 404,
				"error":  "Not Found",
				"code":   "resource_not_found",
			})
		defer server.Close()

		result, err := server.Client().Subscription().Create(context.Background(), &SubscriptionInput{
			ExternalCustomerID: "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
			PlanCode:           "premium",
			ExternalID:         "my_sub_1",
			Name:               "Repository A",
			BillingTime:        Anniversary,
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "provider",
				PaymentMethodID:   "pm_not_found",
			},
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 404)
		c.Assert(err.Message, qt.Equals, "Not Found")
		c.Assert(err.ErrorCode, qt.Equals, "resource_not_found")
	})

	t.Run("When creating a subscription with a payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/subscriptions").
			MatchJSONBody(`{
				"subscription": {
					"external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
					"plan_code": "premium",
					"external_id": "my_sub_1",
					"name": "Repository A",
					"billing_time": "anniversary",
					"payment_method": {
						"payment_method_type": "card",
						"payment_method_id": "pm_123456"
					}
				}
			}`).
			MockResponse(mockSubscriptionResponse)
		defer server.Close()

		subscription, err := server.Client().Subscription().Create(context.Background(), &SubscriptionInput{
			ExternalCustomerID: "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
			PlanCode:           "premium",
			ExternalID:         "my_sub_1",
			Name:               "Repository A",
			BillingTime:        Anniversary,
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "card",
				PaymentMethodID:   "pm_123456",
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(subscription.PaymentMethod, qt.IsNotNil)
		c.Assert(subscription.PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(subscription.PaymentMethod.PaymentMethodID, qt.Equals, "pm_123456")
	})
}

func TestSubscriptionRequest_UpdateWithPaymentMethod(t *testing.T) {
	t.Run("When updating a subscription with an invalid payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/subscriptions/my_sub_1").
			MatchJSONBody(`{
				"subscription": {
					"external_id": "my_sub_1",
					"name": "Repository A",
					"payment_method": {
						"payment_method_type": "invalid",
						"payment_method_id": "pm_invalid"
					}
				}
			}`).
			MockResponseWithCode(422, map[string]any{
				"status": 422,
				"error":  "Unprocessable Entity",
				"code":   "validation_errors",
				"error_details": map[string]any{
					"payment_method": []string{"invalid_payment_method"},
				},
			})
		defer server.Close()

		result, err := server.Client().Subscription().Update(context.Background(), &SubscriptionInput{
			ExternalID: "my_sub_1",
			Name:       "Repository A",
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "invalid",
				PaymentMethodID:   "pm_invalid",
			},
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 422)
		c.Assert(err.Message, qt.Equals, "Unprocessable Entity")
		c.Assert(err.ErrorCode, qt.Equals, "validation_errors")
		c.Assert(err.ErrorDetail, qt.IsNotNil)
		details, detailErr := err.ErrorDetail.Details()
		c.Assert(detailErr, qt.IsNil)
		c.Assert(details["payment_method"], qt.DeepEquals, []string{"invalid_payment_method"})
	})

	t.Run("When updating a subscription with a payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/subscriptions/my_sub_1").
			MatchJSONBody(`{
				"subscription": {
					"external_id": "my_sub_1",
					"name": "Repository A",
					"payment_method": {
						"payment_method_type": "card",
						"payment_method_id": "pm_789"
					}
				}
			}`).
			MockResponse(mockSubscriptionResponse)
		defer server.Close()

		subscription, err := server.Client().Subscription().Update(context.Background(), &SubscriptionInput{
			ExternalID: "my_sub_1",
			Name:       "Repository A",
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "card",
				PaymentMethodID:   "pm_789",
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(subscription.PaymentMethod, qt.IsNotNil)
		c.Assert(subscription.PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(subscription.PaymentMethod.PaymentMethodID, qt.Equals, "pm_123456")
	})
}
