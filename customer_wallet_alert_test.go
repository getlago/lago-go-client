package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockCustomerWalletAlertResponse = map[string]any{
	"alert": map[string]any{
		"lago_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"subscription_external_id": "sub_1234",
		"alert_type":               "wallet_balance_amount",
		"code":                     "wallet_balance_alert",
		"name":                     "Wallet Balance Alert",
		"previous_value":           "0.0",
		"last_processed_at":        nil,
		"thresholds": []map[string]any{
			{
				"code":      "warn",
				"value":     "1000.0",
				"recurring": false,
			},
		},
		"created_at": "2025-03-20T10:00:00Z",
	},
}

var mockCustomerWalletAlertListResponse = map[string]any{
	"alerts": []map[string]any{
		{
			"lago_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"subscription_external_id": "sub_1234",
			"alert_type":               "wallet_balance_amount",
			"code":                     "wallet_balance_alert",
			"name":                     "Wallet Balance Alert",
			"previous_value":           "0.0",
			"last_processed_at":        nil,
			"thresholds": []map[string]any{
				{
					"code":      "warn",
					"value":     "1000.0",
					"recurring": false,
				},
			},
			"created_at": "2025-03-20T10:00:00Z",
		},
		{
			"lago_id":                  "3c903c90-3c90-3c90-3c90-3c903c903c90",
			"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"subscription_external_id": "sub_1234",
			"alert_type":               "wallet_credits_balance",
			"code":                     "wallet_credits_alert",
			"name":                     "Wallet Credits Alert",
			"previous_value":           "0.0",
			"last_processed_at":        nil,
			"thresholds": []map[string]any{
				{
					"code":      "",
					"value":     "2000.0",
					"recurring": false,
				},
			},
			"created_at": "2025-03-20T10:00:00Z",
		},
	},
}

func TestCustomerWalletAlertRequest_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWalletAlert().Get(context.Background(), "customer_id", "wallet_code", "wallet_balance_alert")
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alert is found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts/wallet_balance_alert").
			MockResponse(mockCustomerWalletAlertResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().Get(context.Background(), "customer_id", "wallet_code", "wallet_balance_alert")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(result.Code, qt.Equals, "wallet_balance_alert")
		c.Assert(result.AlertType, qt.Equals, WalletBalanceAmountAlertType)
		c.Assert(result.Thresholds, qt.HasLen, 1)
		c.Assert(result.Thresholds[0].Value, qt.Equals, "1000.0")
	})
}

func TestCustomerWalletAlertRequest_GetList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWalletAlert().GetList(context.Background(), "customer_id", "wallet_code", &AlertListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts").
			MockResponse(mockCustomerWalletAlertListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().GetList(context.Background(), "customer_id", "wallet_code", &AlertListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Alerts, qt.HasLen, 2)
		c.Assert(result.Alerts[0].Code, qt.Equals, "wallet_balance_alert")
		c.Assert(result.Alerts[1].Code, qt.Equals, "wallet_credits_alert")
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts").
			MockResponse(mockCustomerWalletAlertListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().GetList(context.Background(), "customer_id", "wallet_code", &AlertListInput{
			PerPage: Ptr(10),
			Page:    Ptr(1),
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Alerts, qt.HasLen, 2)
		c.Assert(result.Alerts[0].Code, qt.Equals, "wallet_balance_alert")
		c.Assert(result.Alerts[1].Code, qt.Equals, "wallet_credits_alert")
	})
}

func TestCustomerWalletAlertRequest_Create(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWalletAlert().Create(context.Background(), "customer_id", "wallet_code", &AlertInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alert is created", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts").
			MockResponse(mockCustomerWalletAlertResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().Create(context.Background(), "customer_id", "wallet_code", &AlertInput{
			Code:      "wallet_balance_alert",
			AlertType: WalletBalanceAmountAlertType,
			Thresholds: []AlertThreshold{
				{Code: "warn", Value: "1000"},
			},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "wallet_balance_alert")
	})
}

func TestCustomerWalletAlertRequest_CreateList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWalletAlert().CreateList(context.Background(), "customer_id", "wallet_code", []AlertInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alerts are created", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts").
			MockResponse(mockCustomerWalletAlertListResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().CreateList(context.Background(), "customer_id", "wallet_code", []AlertInput{
			{
				Code:      "wallet_balance_alert",
				AlertType: WalletBalanceAmountAlertType,
				Thresholds: []AlertThreshold{
					{Code: "warn", Value: "1000"},
				},
			},
			{
				Code:      "wallet_credits_alert",
				AlertType: WalletCreditsBalanceAlertType,
				Thresholds: []AlertThreshold{
					{Code: "warn", Value: "2000"},
				},
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.HasLen, 2)
		c.Assert(result[0].Code, qt.Equals, "wallet_balance_alert")
		c.Assert(result[0].AlertType, qt.Equals, WalletBalanceAmountAlertType)
		c.Assert(result[1].Code, qt.Equals, "wallet_credits_alert")
		c.Assert(result[1].AlertType, qt.Equals, WalletCreditsBalanceAlertType)
	})
}

func TestCustomerWalletAlertRequest_Update(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWalletAlert().Update(context.Background(), "customer_id", "wallet_code", "wallet_balance_alert", &AlertInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alert is updated", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts/wallet_balance_alert").
			MockResponse(mockCustomerWalletAlertResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().Update(context.Background(), "customer_id", "wallet_code", "wallet_balance_alert", &AlertInput{
			Name: "Updated Alert",
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "wallet_balance_alert")
	})
}

func TestCustomerWalletAlertRequest_Delete(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.CustomerWalletAlert().Delete(context.Background(), "customer_id", "wallet_code", "wallet_balance_alert")
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alert is deleted", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts/wallet_balance_alert").
			MockResponse(mockCustomerWalletAlertResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletAlert().Delete(context.Background(), "customer_id", "wallet_code", "wallet_balance_alert")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "wallet_balance_alert")
	})
}

func TestCustomerWalletAlertRequest_DeleteAll(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		err := client.CustomerWalletAlert().DeleteAll(context.Background(), "customer_id", "wallet_code")

		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alerts are deleted", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/customer_id/wallets/wallet_code/alerts").
			MockResponse(mockCustomerWalletAlertResponse)
		defer server.Close()

		err := server.Client().CustomerWalletAlert().DeleteAll(context.Background(), "customer_id", "wallet_code")

		c.Assert(err == nil, qt.IsTrue)
	})
}
