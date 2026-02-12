package lago_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
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

func TestWallet_Create(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			RateAmount:         "1.00",
			Name:               "wallet name",
			Code:               Ptr("wallet_code"),
			Priority:           Ptr(int(30)),
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
			MockResponse(mockWalletResponse)
		defer server.Close()

		result, err := server.Client().Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID:      "12345",
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
			Priority:    Ptr(int(40)),
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

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{
			Name:        "updated wallet name",
			Priority:    Ptr(int(40)),
			RateAmount:  "1.50",
			PaidCredits: "200.00",
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

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

		result, err := server.Client().Wallet().Delete(context.Background(), "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

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

var mockWalletMetadataResponse = map[string]interface{}{
	"metadata": map[string]interface{}{
		"foo": "bar",
		"baz": nil,
	},
}

var mockWalletNullMetadataResponse = map[string]interface{}{
	"metadata": nil,
}

func TestWalletRequest_ReplaceMetadata(t *testing.T) {
	t.Run("When replace metadata is called", func(t *testing.T) {
		c := qt.New(t)

		walletID := "b1b2c3d4-e5f6-7890-1234-56789abcdef0"

		server := lt.ServerWithAssertions(c, mockWalletMetadataResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0/metadata")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			c.Assert(requestData["metadata"], qt.IsNotNil)
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		bar := "bar"
		result, err := client.Wallet().ReplaceMetadata(context.Background(), walletID, map[string]*string{
			"foo": &bar,
			"baz": nil,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
		c.Assert(result["baz"], qt.IsNil)
	})
}

func TestWalletRequest_MergeMetadata(t *testing.T) {
	t.Run("When merge metadata is called", func(t *testing.T) {
		c := qt.New(t)

		walletID := "b1b2c3d4-e5f6-7890-1234-56789abcdef0"

		server := lt.ServerWithAssertions(c, mockWalletMetadataResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "PATCH")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0/metadata")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			c.Assert(requestData["metadata"], qt.IsNotNil)
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		qux := "qux"
		result, err := client.Wallet().MergeMetadata(context.Background(), walletID, map[string]*string{
			"foo": &qux,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
		c.Assert(result["baz"], qt.IsNil)
	})
}

func TestWalletRequest_DeleteAllMetadata(t *testing.T) {
	t.Run("When delete all metadata is called", func(t *testing.T) {
		c := qt.New(t)

		walletID := "b1b2c3d4-e5f6-7890-1234-56789abcdef0"

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0/metadata").
			MockResponse(mockWalletNullMetadataResponse)
		defer server.Close()

		result, err := server.Client().Wallet().DeleteAllMetadata(context.Background(), walletID)
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNil)
	})
}

func TestWalletRequest_DeleteMetadataKey(t *testing.T) {
	t.Run("When delete metadata key is called", func(t *testing.T) {
		c := qt.New(t)

		walletID := "b1b2c3d4-e5f6-7890-1234-56789abcdef0"

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0/metadata/foo").
			MockResponse(mockWalletMetadataResponse)
		defer server.Close()

		result, err := server.Client().Wallet().DeleteMetadataKey(context.Background(), walletID, "foo")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
	})
}

func TestWallet_CreateWithPaymentMethod(t *testing.T) {
	t.Run("When creating a wallet with an invalid payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/wallets").
			MatchJSONBody(`{
				"wallet": {
					"external_customer_id": "12345",
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

		result, err := server.Client().Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			RateAmount:         "1.00",
			Name:               "wallet name",
			Currency:           "USD",
			PaidCredits:        "100.00",
			GrantedCredits:     "100.00",
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
			MatchPath("/api/v1/wallets").
			MatchJSONBody(`{
				"wallet": {
					"external_customer_id": "12345",
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

		result, err := server.Client().Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			RateAmount:         "1.00",
			Name:               "wallet name",
			Currency:           "USD",
			PaidCredits:        "100.00",
			GrantedCredits:     "100.00",
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
			MatchPath("/api/v1/wallets").
			MatchJSONBody(`{
				"wallet": {
					"external_customer_id": "12345",
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
			MockResponse(mockWalletResponse)
		defer server.Close()

		result, err := server.Client().Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			RateAmount:         "1.00",
			Name:               "wallet name",
			Currency:           "USD",
			PaidCredits:        "100.00",
			GrantedCredits:     "100.00",
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

func TestWallet_UpdateWithPaymentMethod(t *testing.T) {
	t.Run("When updating a wallet with an invalid payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0").
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

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{
			Name:       "updated wallet name",
			RateAmount: "1.50",
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "invalid",
				PaymentMethodID:   "pm_invalid",
			},
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")
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
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0").
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

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{
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
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")
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
				"external_customer_id": "12345",
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
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0").
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

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{
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
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")

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
