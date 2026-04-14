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

var mockCustomerWalletResponse = `{
	"wallet": {
		"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
		"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
		"external_customer_id": "customer_id",
		"status": "active",
		"currency": "USD",
		"name": "wallet name",
		"code": "wallet_code",
		"priority": 30,
		"rate_amount": "1.00",
		"credits_balance": "100.00",
		"balance_cents": 10000,
		"paid_top_up_max_amount_cents": 1000,
		"paid_top_up_min_amount_cents": 200,
		"expiration_at": "2022-07-07T23:59:59Z",
		"created_at": "2022-04-29T08:59:51Z",
		"recurring_transaction_rules": [{
			"lago_id": "c1c2d3d4-e5f6-7890-1234-56789abcdef0",
			"interval": "monthly",
			"method": "fixed",
			"started_at": null,
			"expiration_at": "2026-12-31T23:59:59Z",
			"status": "active",
			"target_ongoing_balance": "0.00",
			"threshold_credits": "0.00",
			"trigger": "interval",
			"paid_credits": "105.00",
			"granted_credits": "105.00",
			"created_at": "2022-04-29T08:59:51Z",
			"invoice_requires_successful_payment": false,
			"transaction_metadata": [],
			"transaction_name": "Recurring Transaction Rule",
			"ignore_paid_top_up_limits": false,
			"payment_method": {
				"payment_method_type": "card",
				"payment_method_id": "pm_rule_123"
			},
			"applied_invoice_custom_sections": [{
				"lago_id": "d1d2e3e4-f5f6-7890-1234-56789abcdef0",
				"invoice_custom_section_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"created_at": "2022-04-29T08:59:51Z",
				"invoice_custom_section": {
					"lago_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
					"code": "rule_section_code",
					"name": "Rule Section Name"
				}
			}]
		}],
		"applies_to": {
			"fee_types": ["charge"],
			"billable_metric_codes": ["bm1"]
		},
		"payment_method": {
			"payment_method_type": "card",
			"payment_method_id": "pm_wallet_123"
		},
		"applied_invoice_custom_sections": [{
			"lago_id": "e1e2f3f4-a5b6-7890-1234-56789abcdef0",
			"invoice_custom_section_id": "b2b3c4c5-d6e7-8901-2345-6789abcdef01",
			"created_at": "2022-04-29T08:59:51Z",
			"invoice_custom_section": {
				"lago_id": "b2b3c4c5-d6e7-8901-2345-6789abcdef01",
				"code": "wallet_section_code",
				"name": "Wallet Section Name"
			}
		}]
	}
}`

var mockCustomerWalletListResponse = `{
	"wallets": [{
		"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
		"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
		"external_customer_id": "customer_id",
		"status": "active",
		"currency": "USD",
		"name": "wallet name",
		"code": "wallet_code",
		"priority": 30,
		"rate_amount": "1.00",
		"credits_balance": "100.00",
		"balance_cents": 10000,
		"paid_top_up_max_amount_cents": 1000,
		"paid_top_up_min_amount_cents": 200,
		"expiration_at": "2022-07-07T23:59:59Z",
		"created_at": "2022-04-29T08:59:51Z",
		"recurring_transaction_rules": [{
			"lago_id": "c1c2d3d4-e5f6-7890-1234-56789abcdef0",
			"interval": "monthly",
			"method": "fixed",
			"started_at": null,
			"expiration_at": "2026-12-31T23:59:59Z",
			"status": "active",
			"target_ongoing_balance": "0.00",
			"threshold_credits": "0.00",
			"trigger": "interval",
			"paid_credits": "105.00",
			"granted_credits": "105.00",
			"created_at": "2022-04-29T08:59:51Z",
			"invoice_requires_successful_payment": false,
			"transaction_metadata": [],
			"transaction_name": "Recurring Transaction Rule",
			"ignore_paid_top_up_limits": false,
			"payment_method": {
				"payment_method_type": "card",
				"payment_method_id": "pm_rule_123"
			}
		}],
		"applies_to": {
			"fee_types": ["charge"],
			"billable_metric_codes": ["bm1"]
		},
		"payment_method": {
			"payment_method_type": "card",
			"payment_method_id": "pm_wallet_123"
		}
	}],
	"meta": {
		"current_page": 1,
		"next_page": 2,
		"prev_page": null,
		"total_pages": 7,
		"total_count": 63
	}
}`

func TestCustomerWallet_Create(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().Create(context.Background(), "customer_id", &WalletInput{
			RateAmount:     "1.00",
			Name:           "wallet name",
			Code:           Ptr("wallet_code"),
			Priority:       Ptr(int(30)),
			PaidCredits:    "100.00",
			GrantedCredits: "100.00",
			ExpirationAt:   Ptr(time.Date(2022, 7, 7, 23, 59, 59, 0, time.UTC)),
			RecurringTransactionRules: []RecurringTransactionRuleInput{
				{
					PaidCredits:      "105.00",
					GrantedCredits:   "105.00",
					ThresholdCredits: "0.00",
					Trigger:          "interval",
					Interval:         "monthly",
					Method:           "fixed",
					StartedAt:        nil,
					ExpirationAt:     Ptr(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)),
					TransactionName:  "Recurring Transaction Rule",
				},
			},
			AppliesTo: AppliesTo{
				FeeTypes:            []string{"charge"},
				BillableMetricCodes: []string{"bm1"},
			},
		})

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers/customer_id/wallets").
			MatchJSONBody(`{
				"wallet": {
					"rate_amount": "1.00",
					"name": "wallet name",
					"code": "wallet_code",
					"priority": 30,
					"paid_credits": "100.00",
					"granted_credits": "100.00",
                    "transaction_name": "wallet transaction name",
					"expiration_at": "2022-07-07T23:59:59Z",
					"paid_top_up_max_amount_cents": 1000,
					"paid_top_up_min_amount_cents": 200,
					"recurring_transaction_rules": [
						{
							"paid_credits": "105.00",
							"granted_credits": "105.00",
							"threshold_credits": "0.00",
							"lago_id": "00000000-0000-0000-0000-000000000000",
							"trigger": "interval",
							"interval": "monthly",
							"method": "fixed",
							"expiration_at": "2026-12-31T23:59:59Z",
							"transaction_name": "Recurring Transaction Rule",
							"ignore_paid_top_up_limits": true
						}
					],
					"applies_to": {
						"fee_types": ["charge"],
						"billable_metric_codes": ["bm1"]
					}
				}
			}`).
			MockResponse(mockCustomerWalletResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().Create(context.Background(), "customer_id", &WalletInput{
			RateAmount:              "1.00",
			Name:                    "wallet name",
			Code:                    Ptr("wallet_code"),
			Priority:                Ptr(int(30)),
			PaidCredits:             "100.00",
			GrantedCredits:          "100.00",
			TransactionName:         "wallet transaction name",
			ExpirationAt:            Ptr(time.Date(2022, 7, 7, 23, 59, 59, 0, time.UTC)),
			PaidTopUpMaxAmountCents: Ptr(int(1000)),
			PaidTopUpMinAmountCents: Ptr(int(200)),
			RecurringTransactionRules: []RecurringTransactionRuleInput{
				{
					PaidCredits:           "105.00",
					GrantedCredits:        "105.00",
					ThresholdCredits:      "0.00",
					Trigger:               "interval",
					Interval:              "monthly",
					Method:                "fixed",
					StartedAt:             nil,
					ExpirationAt:          Ptr(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)),
					TransactionName:       "Recurring Transaction Rule",
					IgnorePaidTopUpLimits: Ptr(true),
				},
			},
			AppliesTo: AppliesTo{
				FeeTypes:            []string{"charge"},
				BillableMetricCodes: []string{"bm1"},
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.Code, qt.DeepEquals, Ptr("wallet_code"))
		c.Assert(result.Priority, qt.Equals, int(30))
		c.Assert(result.ExternalCustomerID, qt.Equals, "customer_id")
		c.Assert(result.RateAmount, qt.Equals, "1.00")
		c.Assert(result.CreditsBalance, qt.Equals, "100.00")
		c.Assert(result.Status, qt.Equals, Status("active"))
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
		c.Assert(result.PaidTopUpMaxAmountCents, qt.IsNotNil)
		c.Assert(*result.PaidTopUpMaxAmountCents, qt.Equals, int(1000))
		c.Assert(result.PaidTopUpMinAmountCents, qt.IsNotNil)
		c.Assert(*result.PaidTopUpMinAmountCents, qt.Equals, int(200))
		c.Assert(len(result.RecurringTransactionRules), qt.Equals, 1)
		c.Assert(result.RecurringTransactionRules[0].Trigger, qt.Equals, "interval")
		c.Assert(result.RecurringTransactionRules[0].Interval, qt.Equals, "monthly")
		c.Assert(result.RecurringTransactionRules[0].TransactionName, qt.Equals, "Recurring Transaction Rule")
		c.Assert(result.RecurringTransactionRules[0].IgnorePaidTopUpLimits, qt.Equals, false)
		c.Assert(result.RecurringTransactionRules[0].ExpirationAt, qt.DeepEquals, Ptr(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)))
		c.Assert(result.AppliesTo.FeeTypes, qt.DeepEquals, []string{"charge"})
		c.Assert(result.AppliesTo.BillableMetricCodes, qt.DeepEquals, []string{"bm1"})
	})
}

func TestCustomerWallet_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().Get(context.Background(), "customer_id", "wallet_code")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code").
			MockResponse(mockCustomerWalletResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().Get(context.Background(), "customer_id", "wallet_code")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.Code, qt.DeepEquals, Ptr("wallet_code"))
		c.Assert(result.Priority, qt.Equals, int(30))
		c.Assert(result.ExternalCustomerID, qt.Equals, "customer_id")
		c.Assert(result.RateAmount, qt.Equals, "1.00")
		c.Assert(result.CreditsBalance, qt.Equals, "100.00")
		c.Assert(result.Status, qt.Equals, Status("active"))
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
		c.Assert(result.PaidTopUpMaxAmountCents, qt.IsNotNil)
		c.Assert(*result.PaidTopUpMaxAmountCents, qt.Equals, int(1000))
		c.Assert(result.PaidTopUpMinAmountCents, qt.IsNotNil)
		c.Assert(*result.PaidTopUpMinAmountCents, qt.Equals, int(200))
		c.Assert(len(result.RecurringTransactionRules), qt.Equals, 1)
		c.Assert(result.RecurringTransactionRules[0].Trigger, qt.Equals, "interval")
		c.Assert(result.RecurringTransactionRules[0].Interval, qt.Equals, "monthly")
		c.Assert(result.RecurringTransactionRules[0].TransactionName, qt.Equals, "Recurring Transaction Rule")
		c.Assert(result.RecurringTransactionRules[0].IgnorePaidTopUpLimits, qt.Equals, false)
		c.Assert(result.RecurringTransactionRules[0].AppliedInvoiceCustomSections, qt.HasLen, 1)
		c.Assert(result.RecurringTransactionRules[0].AppliedInvoiceCustomSections[0].LagoId, qt.Equals, uuid.MustParse("d1d2e3e4-f5f6-7890-1234-56789abcdef0"))
		c.Assert(result.RecurringTransactionRules[0].AppliedInvoiceCustomSections[0].InvoiceCustomSectionId, qt.Equals, uuid.MustParse("a1a2b3b4-c5d6-7890-1234-56789abcdef0"))
		c.Assert(result.AppliesTo.FeeTypes, qt.DeepEquals, []string{"charge"})
		c.Assert(result.AppliesTo.BillableMetricCodes, qt.DeepEquals, []string{"bm1"})
		c.Assert(result.AppliedInvoiceCustomSections, qt.HasLen, 1)
		c.Assert(result.AppliedInvoiceCustomSections[0].LagoId, qt.Equals, uuid.MustParse("e1e2f3f4-a5b6-7890-1234-56789abcdef0"))
		c.Assert(result.AppliedInvoiceCustomSections[0].InvoiceCustomSectionId, qt.Equals, uuid.MustParse("b2b3c4c5-d6e7-8901-2345-6789abcdef01"))
		c.Assert(result.AppliedInvoiceCustomSections[0].InvoiceCustomSection.LagoId, qt.Equals, uuid.MustParse("b2b3c4c5-d6e7-8901-2345-6789abcdef01"))
		c.Assert(result.AppliedInvoiceCustomSections[0].InvoiceCustomSection.Name, qt.Equals, "Wallet Section Name")
	})
}

func TestCustomerWallet_GetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().GetList(context.Background(), "customer_id", &WalletListInput{})

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/customer_id/wallets").
			MockResponse(mockCustomerWalletListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().GetList(context.Background(), "customer_id", &WalletListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Wallets, qt.HasLen, 1)
		c.Assert(result.Wallets[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Wallets[0].Name, qt.Equals, "wallet name")
		c.Assert(result.Wallets[0].Priority, qt.Equals, 30)
		c.Assert(result.Wallets[0].ExternalCustomerID, qt.Equals, "customer_id")
		c.Assert(result.Wallets[0].Status, qt.Equals, Status("active"))
		c.Assert(result.Wallets[0].Currency, qt.Equals, Currency("USD"))
		c.Assert(result.Wallets[0].PaidTopUpMaxAmountCents, qt.IsNotNil)
		c.Assert(*result.Wallets[0].PaidTopUpMaxAmountCents, qt.Equals, int(1000))
		c.Assert(result.Wallets[0].PaidTopUpMinAmountCents, qt.IsNotNil)
		c.Assert(*result.Wallets[0].PaidTopUpMinAmountCents, qt.Equals, int(200))
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.NextPage, qt.Equals, 2)
		c.Assert(result.Meta.TotalPages, qt.Equals, 7)
		c.Assert(result.Meta.TotalCount, qt.Equals, 63)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/customer_id/wallets").
			MockResponse(mockCustomerWalletListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().GetList(context.Background(), "customer_id", &WalletListInput{
			PerPage: Ptr(10),
			Page:    Ptr(1),
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Wallets, qt.HasLen, 1)
		c.Assert(result.Wallets[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Wallets[0].ExternalCustomerID, qt.Equals, "customer_id")
	})
}

func TestCustomerWallet_Update(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().Update(context.Background(), "customer_id", "wallet_code", &WalletInput{
			Name:        "updated wallet name",
			Priority:    Ptr(int(40)),
			RateAmount:  "1.50",
			PaidCredits: "200.00",
		})

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code")
		defer server.Close()

		updatedWalletResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "customer_id",
				"status": "active",
				"currency": "USD",
				"name": "updated wallet name",
				"priority": 40,
				"rate_amount": "1.50",
				"credits_balance": "200.00",
    "balance_cents": 20000,
				"paid_top_up_max_amount_cents": 1500,
				"paid_top_up_min_amount_cents": 300,
				"expiration_at": "2022-07-07T23:59:59Z",
				"created_at": "2022-04-29T08:59:51Z",
				"recurring_transaction_rules": [],
				"applies_to": {
					"fee_types": ["charge"],
					"billable_metric_codes": ["bm1"]
				}
			}
		}`

		server.MockResponse(updatedWalletResponse)

		result, err := server.Client().CustomerWallet().Update(context.Background(), "customer_id", "wallet_code", &WalletInput{
			Name:        "updated wallet name",
			Priority:    Ptr(int(40)),
			RateAmount:  "1.50",
			PaidCredits: "200.00",
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "updated wallet name")
		c.Assert(result.Priority, qt.Equals, int(40))
		c.Assert(result.RateAmount, qt.Equals, "1.50")
		c.Assert(result.CreditsBalance, qt.Equals, "200.00")
		c.Assert(result.BalanceCents, qt.Equals, 20000)
		c.Assert(result.Status, qt.Equals, Status("active"))
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
	})
}

func TestCustomerWallet_Delete(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().Delete(context.Background(), "customer_id", "wallet_code")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code")
		defer server.Close()

		deletedWalletResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "customer_id",
				"status": "terminated",
				"currency": "USD",
				"name": "wallet name",
				"priority": 30,
				"rate_amount": "1.00",
				"credits_balance": "100.00",
    "balance_cents": 10000,
				"paid_top_up_max_amount_cents": 1000,
				"paid_top_up_min_amount_cents": 200,
				"expiration_at": "2022-07-07T23:59:59Z",
				"created_at": "2022-04-29T08:59:51Z",
				"terminated_at": "2022-07-08T10:00:00Z",
				"recurring_transaction_rules": [],
				"applies_to": {
					"fee_types": ["charge"],
					"billable_metric_codes": ["bm1"]
				}
			}
		}`

		server.MockResponse(deletedWalletResponse)

		result, err := server.Client().CustomerWallet().Delete(context.Background(), "customer_id", "wallet_code")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Status, qt.Equals, Status("terminated"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.Priority, qt.Equals, int(30))
		c.Assert(result.ExternalCustomerID, qt.Equals, "customer_id")
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
	})
}

func TestCustomerWallet_CreateWithPaymentMethod(t *testing.T) {
	t.Run("When creating a wallet with an invalid payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers/customer_id/wallets").
			MatchJSONBody(`{
				"wallet": {
					"rate_amount": "1.00",
					"name": "wallet name",
					"currency": "USD",
					"paid_credits": "100.00",
					"granted_credits": "100.00",
					"recurring_transaction_rules": null,
					"applies_to": {},
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

		result, err := server.Client().CustomerWallet().Create(context.Background(), "customer_id", &WalletInput{
			RateAmount:     "1.00",
			Name:           "wallet name",
			Currency:       "USD",
			PaidCredits:    "100.00",
			GrantedCredits: "100.00",
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

	t.Run("When creating a wallet with an invalid payment method in recurring transaction rules", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers/customer_id/wallets").
			MatchJSONBody(`{
				"wallet": {
					"rate_amount": "1.00",
					"name": "wallet name",
					"currency": "USD",
					"paid_credits": "100.00",
					"granted_credits": "100.00",
					"applies_to": {},
					"recurring_transaction_rules": [
						{
							"paid_credits": "105.00",
							"granted_credits": "105.00",
							"threshold_credits": "0.00",
							"lago_id": "00000000-0000-0000-0000-000000000000",
							"trigger": "interval",
							"interval": "monthly",
							"method": "fixed",
							"transaction_name": "Recurring Transaction Rule",
							"payment_method": {
								"payment_method_type": "invalid",
								"payment_method_id": "pm_invalid"
							}
						}
					]
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

		result, err := server.Client().CustomerWallet().Create(context.Background(), "customer_id", &WalletInput{
			RateAmount:     "1.00",
			Name:           "wallet name",
			Currency:       "USD",
			PaidCredits:    "100.00",
			GrantedCredits: "100.00",
			RecurringTransactionRules: []RecurringTransactionRuleInput{
				{
					PaidCredits:      "105.00",
					GrantedCredits:   "105.00",
					ThresholdCredits: "0.00",
					Trigger:          "interval",
					Interval:         "monthly",
					Method:           "fixed",
					TransactionName:  "Recurring Transaction Rule",
					PaymentMethod: &PaymentMethodInput{
						PaymentMethodType: "invalid",
						PaymentMethodID:   "pm_invalid",
					},
				},
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

	t.Run("When creating a wallet with a payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers/customer_id/wallets").
			MatchJSONBody(`{
				"wallet": {
					"rate_amount": "1.00",
					"name": "wallet name",
					"currency": "USD",
					"paid_credits": "100.00",
					"granted_credits": "100.00",
					"applies_to": {},
					"recurring_transaction_rules": [
						{
							"paid_credits": "105.00",
							"granted_credits": "105.00",
							"threshold_credits": "0.00",
							"lago_id": "00000000-0000-0000-0000-000000000000",
							"trigger": "interval",
							"interval": "monthly",
							"method": "fixed",
							"transaction_name": "Recurring Transaction Rule",
							"payment_method": {
								"payment_method_type": "card",
								"payment_method_id": "pm_rule_123"
							}
						}
					],
					"payment_method": {
						"payment_method_type": "card",
						"payment_method_id": "pm_wallet_123"
					}
				}
			}`).
			MockResponse(mockCustomerWalletResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().Create(context.Background(), "customer_id", &WalletInput{
			RateAmount:     "1.00",
			Name:           "wallet name",
			Currency:       "USD",
			PaidCredits:    "100.00",
			GrantedCredits: "100.00",
			RecurringTransactionRules: []RecurringTransactionRuleInput{
				{
					PaidCredits:      "105.00",
					GrantedCredits:   "105.00",
					ThresholdCredits: "0.00",
					Trigger:          "interval",
					Interval:         "monthly",
					Method:           "fixed",
					TransactionName:  "Recurring Transaction Rule",
					PaymentMethod: &PaymentMethodInput{
						PaymentMethodType: "card",
						PaymentMethodID:   "pm_rule_123",
					},
				},
			},
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "card",
				PaymentMethodID:   "pm_wallet_123",
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.PaymentMethod, qt.IsNotNil)
		c.Assert(result.PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(result.PaymentMethod.PaymentMethodID, qt.Equals, "pm_wallet_123")
		c.Assert(len(result.RecurringTransactionRules), qt.Equals, 1)
		c.Assert(result.RecurringTransactionRules[0].PaymentMethod, qt.IsNotNil)
		c.Assert(result.RecurringTransactionRules[0].PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(result.RecurringTransactionRules[0].PaymentMethod.PaymentMethodID, qt.Equals, "pm_rule_123")
	})
}

func TestCustomerWallet_UpdateWithPaymentMethod(t *testing.T) {
	t.Run("When updating a wallet with an invalid payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code").
			MatchJSONBody(`{
				"wallet": {
					"name": "updated wallet name",
					"rate_amount": "1.50",
					"recurring_transaction_rules": null,
					"applies_to": {},
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

		result, err := server.Client().CustomerWallet().Update(context.Background(), "customer_id", "wallet_code", &WalletInput{
			Name:       "updated wallet name",
			RateAmount: "1.50",
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

	t.Run("When updating a wallet with an invalid payment method in recurring transaction rules", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code").
			MatchJSONBody(`{
				"wallet": {
					"name": "updated wallet name",
					"rate_amount": "1.50",
					"applies_to": {},
					"recurring_transaction_rules": [
						{
							"lago_id": "c1c2d3d4-e5f6-7890-1234-56789abcdef0",
							"trigger": "interval",
							"interval": "monthly",
							"method": "fixed",
							"paid_credits": "105.00",
							"granted_credits": "105.00",
							"threshold_credits": "0.00",
							"transaction_name": "Recurring Transaction Rule",
							"payment_method": {
								"payment_method_type": "invalid",
								"payment_method_id": "pm_invalid"
							}
						}
					]
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

		result, err := server.Client().CustomerWallet().Update(context.Background(), "customer_id", "wallet_code", &WalletInput{
			Name:       "updated wallet name",
			RateAmount: "1.50",
			RecurringTransactionRules: []RecurringTransactionRuleInput{
				{
					LagoID:           uuid.MustParse("c1c2d3d4-e5f6-7890-1234-56789abcdef0"),
					Trigger:          "interval",
					Interval:         "monthly",
					Method:           "fixed",
					PaidCredits:      "105.00",
					GrantedCredits:   "105.00",
					ThresholdCredits: "0.00",
					TransactionName:  "Recurring Transaction Rule",
					PaymentMethod: &PaymentMethodInput{
						PaymentMethodType: "invalid",
						PaymentMethodID:   "pm_invalid",
					},
				},
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

	t.Run("When updating a wallet with a payment method", func(t *testing.T) {
		c := qt.New(t)

		updatedWalletWithPaymentMethodResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "customer_id",
				"status": "active",
				"currency": "USD",
				"name": "updated wallet name",
				"priority": 40,
				"rate_amount": "1.50",
				"credits_balance": "200.00",
				"balance_cents": 20000,
				"expiration_at": "2022-07-07T23:59:59Z",
				"created_at": "2022-04-29T08:59:51Z",
				"recurring_transaction_rules": [{
					"lago_id": "c1c2d3d4-e5f6-7890-1234-56789abcdef0",
					"interval": "monthly",
					"method": "fixed",
					"status": "active",
					"target_ongoing_balance": "0.00",
					"threshold_credits": "0.00",
					"trigger": "interval",
					"paid_credits": "105.00",
					"granted_credits": "105.00",
					"created_at": "2022-04-29T08:59:51Z",
					"invoice_requires_successful_payment": false,
					"transaction_metadata": [],
					"transaction_name": "Recurring Transaction Rule",
					"ignore_paid_top_up_limits": false,
					"payment_method": {
						"payment_method_type": "card",
						"payment_method_id": "pm_rule_updated"
					}
				}],
				"applies_to": {
					"fee_types": ["charge"],
					"billable_metric_codes": ["bm1"]
				},
				"payment_method": {
					"payment_method_type": "card",
					"payment_method_id": "pm_wallet_updated"
				}
			}
		}`

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code").
			MatchJSONBody(`{
				"wallet": {
					"name": "updated wallet name",
					"rate_amount": "1.50",
					"applies_to": {},
					"recurring_transaction_rules": [
						{
							"lago_id": "c1c2d3d4-e5f6-7890-1234-56789abcdef0",
							"trigger": "interval",
							"interval": "monthly",
							"method": "fixed",
							"paid_credits": "105.00",
							"granted_credits": "105.00",
							"threshold_credits": "0.00",
							"transaction_name": "Recurring Transaction Rule",
							"payment_method": {
								"payment_method_type": "card",
								"payment_method_id": "pm_rule_updated"
							}
						}
					],
					"payment_method": {
						"payment_method_type": "card",
						"payment_method_id": "pm_wallet_updated"
					}
				}
			}`).
			MockResponse(updatedWalletWithPaymentMethodResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().Update(context.Background(), "customer_id", "wallet_code", &WalletInput{
			Name:       "updated wallet name",
			RateAmount: "1.50",
			RecurringTransactionRules: []RecurringTransactionRuleInput{
				{
					LagoID:           uuid.MustParse("c1c2d3d4-e5f6-7890-1234-56789abcdef0"),
					Trigger:          "interval",
					Interval:         "monthly",
					Method:           "fixed",
					PaidCredits:      "105.00",
					GrantedCredits:   "105.00",
					ThresholdCredits: "0.00",
					TransactionName:  "Recurring Transaction Rule",
					PaymentMethod: &PaymentMethodInput{
						PaymentMethodType: "card",
						PaymentMethodID:   "pm_rule_updated",
					},
				},
			},
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "card",
				PaymentMethodID:   "pm_wallet_updated",
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.PaymentMethod, qt.IsNotNil)
		c.Assert(result.PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(result.PaymentMethod.PaymentMethodID, qt.Equals, "pm_wallet_updated")
		c.Assert(len(result.RecurringTransactionRules), qt.Equals, 1)
		c.Assert(result.RecurringTransactionRules[0].PaymentMethod, qt.IsNotNil)
		c.Assert(result.RecurringTransactionRules[0].PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(result.RecurringTransactionRules[0].PaymentMethod.PaymentMethodID, qt.Equals, "pm_rule_updated")
	})
}
