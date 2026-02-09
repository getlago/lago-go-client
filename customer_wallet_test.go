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
		"external_customer_id": "12345",
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
			"ignore_paid_top_up_limits": false
		}],
		"applies_to": {
			"fee_types": ["charge"],
			"billable_metric_codes": ["bm1"]
		}
	}
}`

var mockCustomerWalletListResponse = `{
	"wallets": [{
		"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
		"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
		"external_customer_id": "12345",
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
			"ignore_paid_top_up_limits": false
		}],
		"applies_to": {
			"fee_types": ["charge"],
			"billable_metric_codes": ["bm1"]
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
		result, err := client.CustomerWallet().Create(context.Background(), "12345", &WalletInput{
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
			MatchPath("/api/v1/customers/12345/wallets").
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

		result, err := server.Client().CustomerWallet().Create(context.Background(), "12345", &WalletInput{
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
		c.Assert(result.ExternalCustomerID, qt.Equals, "12345")
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
		result, err := client.CustomerWallet().Get(context.Background(), "12345", "wallet_code")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/12345/wallets/wallet_code").
			MockResponse(mockCustomerWalletResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().Get(context.Background(), "12345", "wallet_code")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.Code, qt.DeepEquals, Ptr("wallet_code"))
		c.Assert(result.Priority, qt.Equals, int(30))
		c.Assert(result.ExternalCustomerID, qt.Equals, "12345")
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
		c.Assert(result.AppliesTo.FeeTypes, qt.DeepEquals, []string{"charge"})
		c.Assert(result.AppliesTo.BillableMetricCodes, qt.DeepEquals, []string{"bm1"})
	})
}

func TestCustomerWallet_GetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().GetList(context.Background(), "12345", &WalletListInput{})

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/12345/wallets").
			MockResponse(mockCustomerWalletListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().GetList(context.Background(), "12345", &WalletListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Wallets, qt.HasLen, 1)
		c.Assert(result.Wallets[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Wallets[0].Name, qt.Equals, "wallet name")
		c.Assert(result.Wallets[0].Priority, qt.Equals, 30)
		c.Assert(result.Wallets[0].ExternalCustomerID, qt.Equals, "12345")
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
			MatchPath("/api/v1/customers/12345/wallets").
			MockResponse(mockCustomerWalletListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWallet().GetList(context.Background(), "12345", &WalletListInput{
			PerPage: Ptr(10),
			Page:    Ptr(1),
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Wallets, qt.HasLen, 1)
		c.Assert(result.Wallets[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Wallets[0].ExternalCustomerID, qt.Equals, "12345")
	})
}

func TestCustomerWallet_Update(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWallet().Update(context.Background(), "12345", "wallet_code", &WalletInput{
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
			MatchPath("/api/v1/customers/12345/wallets/wallet_code")
		defer server.Close()

		updatedWalletResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "12345",
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

		result, err := server.Client().CustomerWallet().Update(context.Background(), "12345", "wallet_code", &WalletInput{
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
		result, err := client.CustomerWallet().Delete(context.Background(), "12345", "wallet_code")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/12345/wallets/wallet_code")
		defer server.Close()

		deletedWalletResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "12345",
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

		result, err := server.Client().CustomerWallet().Delete(context.Background(), "12345", "wallet_code")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Status, qt.Equals, Status("terminated"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.Priority, qt.Equals, int(30))
		c.Assert(result.ExternalCustomerID, qt.Equals, "12345")
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
	})
}
