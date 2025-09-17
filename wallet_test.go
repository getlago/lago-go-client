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

var mockWalletResponse = `{
	"wallet": {
		"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
		"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
		"external_customer_id": "12345",
		"status": "active",
		"currency": "USD",
		"name": "wallet name",
		"rate_amount": "1.00",
		"credits_balance": "100.00",
		"balance_cents": 10000,
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
			"transaction_name": "Recurring Transaction Rule"
		}],
		"applies_to": {
			"fee_types": ["charge"],
			"billable_metric_codes": ["bm1"]
		}
	}
}`

var mockWalletListResponse = `{
	"wallets": [{
		"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
		"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
		"external_customer_id": "12345",
		"status": "active",
		"currency": "USD",
		"name": "wallet name",
		"rate_amount": "1.00",
		"credits_balance": "100.00",
		"balance_cents": 10000,
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
			"transaction_name": "Recurring Transaction Rule"
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

func TestWallet_Create(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			RateAmount:         "1.00",
			Name:               "wallet name",
			PaidCredits:        "100.00",
			GrantedCredits:     "100.00",
			ExpirationAt:       Ptr(time.Date(2022, 7, 7, 23, 59, 59, 0, time.UTC)),
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
			MatchPath("/api/v1/wallets").
			MatchJSONBody(`{
				"wallet": {
					"external_customer_id": "12345",
					"rate_amount": "1.00",
					"name": "wallet name",
					"paid_credits": "100.00",
					"granted_credits": "100.00",
					"transaction_name": "wallet transaction name",
					"expiration_at": "2022-07-07T23:59:59Z",
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
							"transaction_name": "Recurring Transaction Rule"
						}
					],
					"applies_to": {
						"fee_types": ["charge"],
						"billable_metric_codes": ["bm1"]
					}
				}
			}`).
			MockResponse(mockWalletResponse)
		defer server.Close()

		result, err := server.Client().Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			RateAmount:         "1.00",
			Name:               "wallet name",
			PaidCredits:        "100.00",
			GrantedCredits:     "100.00",
			TransactionName:    "wallet transaction name",
			ExpirationAt:       Ptr(time.Date(2022, 7, 7, 23, 59, 59, 0, time.UTC)),
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

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.ExternalCustomerID, qt.Equals, "12345")
		c.Assert(result.RateAmount, qt.Equals, "1.00")
		c.Assert(result.CreditsBalance, qt.Equals, "100.00")
		c.Assert(result.Status, qt.Equals, Status("active"))
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
		c.Assert(len(result.RecurringTransactionRules), qt.Equals, 1)
		c.Assert(result.RecurringTransactionRules[0].Trigger, qt.Equals, "interval")
		c.Assert(result.RecurringTransactionRules[0].Interval, qt.Equals, "monthly")
		c.Assert(result.RecurringTransactionRules[0].TransactionName, qt.Equals, "Recurring Transaction Rule")
		c.Assert(result.RecurringTransactionRules[0].ExpirationAt, qt.DeepEquals, Ptr(time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)))
		c.Assert(result.AppliesTo.FeeTypes, qt.DeepEquals, []string{"charge"})
		c.Assert(result.AppliesTo.BillableMetricCodes, qt.DeepEquals, []string{"bm1"})
	})
}

func TestWallet_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Wallet().Get(context.Background(), "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c)
		defer server.Close()

		server.MockResponse(mockWalletResponse)

		result, err := server.Client().Wallet().Get(context.Background(), "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.ExternalCustomerID, qt.Equals, "12345")
		c.Assert(result.RateAmount, qt.Equals, "1.00")
		c.Assert(result.CreditsBalance, qt.Equals, "100.00")
		c.Assert(result.Status, qt.Equals, Status("active"))
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
		c.Assert(len(result.RecurringTransactionRules), qt.Equals, 1)
		c.Assert(result.RecurringTransactionRules[0].Trigger, qt.Equals, "interval")
		c.Assert(result.RecurringTransactionRules[0].Interval, qt.Equals, "monthly")
		c.Assert(result.RecurringTransactionRules[0].TransactionName, qt.Equals, "Recurring Transaction Rule")
		c.Assert(result.AppliesTo.FeeTypes, qt.DeepEquals, []string{"charge"})
		c.Assert(result.AppliesTo.BillableMetricCodes, qt.DeepEquals, []string{"bm1"})
	})
}

func TestWallet_GetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Wallet().GetList(context.Background(), &WalletListInput{})

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c)
		defer server.Close()

		server.MockResponse(mockWalletListResponse)

		result, err := server.Client().Wallet().GetList(context.Background(), &WalletListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Wallets, qt.HasLen, 1)
		c.Assert(result.Wallets[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Wallets[0].Name, qt.Equals, "wallet name")
		c.Assert(result.Wallets[0].ExternalCustomerID, qt.Equals, "12345")
		c.Assert(result.Wallets[0].Status, qt.Equals, Status("active"))
		c.Assert(result.Wallets[0].Currency, qt.Equals, Currency("USD"))
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.NextPage, qt.Equals, 2)
		c.Assert(result.Meta.TotalPages, qt.Equals, 7)
		c.Assert(result.Meta.TotalCount, qt.Equals, 63)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c)
		defer server.Close()

		server.MockResponse(mockWalletListResponse)

		result, err := server.Client().Wallet().GetList(context.Background(), &WalletListInput{
			ExternalCustomerID: "12345",
			PerPage:            Ptr(10),
			Page:               Ptr(1),
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Wallets, qt.HasLen, 1)
		c.Assert(result.Wallets[0].LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Wallets[0].ExternalCustomerID, qt.Equals, "12345")
	})
}

func TestWallet_Update(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Wallet().Update(context.Background(), &WalletInput{
			Name:        "updated wallet name",
			RateAmount:  "1.50",
			PaidCredits: "200.00",
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c)
		defer server.Close()

		updatedWalletResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "12345",
				"status": "active",
				"currency": "USD",
				"name": "updated wallet name",
				"rate_amount": "1.50",
				"credits_balance": "200.00",
				"balance_cents": 20000,
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

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{
			Name:        "updated wallet name",
			RateAmount:  "1.50",
			PaidCredits: "200.00",
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Name, qt.Equals, "updated wallet name")
		c.Assert(result.RateAmount, qt.Equals, "1.50")
		c.Assert(result.CreditsBalance, qt.Equals, "200.00")
		c.Assert(result.BalanceCents, qt.Equals, 20000)
		c.Assert(result.Status, qt.Equals, Status("active"))
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
	})
}

func TestWallet_Delete(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Wallet().Delete(context.Background(), "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

		c.Assert(err, qt.IsNotNil)
		c.Assert(result, qt.IsNil)
	})

	t.Run("When the server returns a successful response", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c)
		defer server.Close()

		deletedWalletResponse := `{
			"wallet": {
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_customer_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"external_customer_id": "12345",
				"status": "terminated",
				"currency": "USD",
				"name": "wallet name",
				"rate_amount": "1.00",
				"credits_balance": "100.00",
				"balance_cents": 10000,
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

		result, err := server.Client().Wallet().Delete(context.Background(), "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID, qt.Equals, uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"))
		c.Assert(result.Status, qt.Equals, Status("terminated"))
		c.Assert(result.Name, qt.Equals, "wallet name")
		c.Assert(result.ExternalCustomerID, qt.Equals, "12345")
		c.Assert(result.Currency, qt.Equals, Currency("USD"))
	})
}
