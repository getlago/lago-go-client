package lago_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"

	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockAlertResponse = map[string]any{
	"alert": map[string]any{
		"lago_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"subscription_external_id": "sub_1234",
		"external_subscription_id": "sub_1234",
		"lago_wallet_id":           "5e905e90-5e90-5e90-5e90-5e905e905e90",
		"wallet_code":              "main_wallet",
		"alert_type":               "current_usage_amount",
		"direction":                "decreasing",
		"code":                     "usage_alert",
		"name":                     "Usage Alert",
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

var mockAlertsResponse = map[string]any{
	"alerts": []map[string]any{
		{
			"lago_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"subscription_external_id": "sub_1234",
			"alert_type":               "current_usage_amount",
			"code":                     "alert1",
			"name":                     "First Alert",
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
			"alert_type":               "billable_metric_current_usage_amount",
			"code":                     "alert2",
			"name":                     "Second Alert",
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

func TestAlertRequest_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Alert().Get(context.Background(), "sub_1234", "usage_alert")
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alert is found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/usage_alert").
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Get(context.Background(), "sub_1234", "usage_alert")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(result.Code, qt.Equals, "usage_alert")
		c.Assert(result.AlertType, qt.Equals, CurrentUsageAmountAlertType)
		c.Assert(result.Thresholds, qt.HasLen, 1)
		c.Assert(result.Thresholds[0].Value, qt.Equals, "1000.0")
		c.Assert(*result.ExternalSubscriptionID, qt.Equals, "sub_1234")
		c.Assert(result.Direction, qt.Equals, DecreasingAlertDirection)
		c.Assert(result.LagoWalletID.String(), qt.Equals, "5e905e90-5e90-5e90-5e90-5e905e905e90")
		c.Assert(*result.WalletCode, qt.Equals, "main_wallet")
	})

	t.Run("When the alert is a billable metric lifetime usage units alert", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/units_alert").
			MockResponse(map[string]any{
				"alert": map[string]any{
					"lago_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
					"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
					"subscription_external_id": "sub_1234",
					"external_subscription_id": "sub_1234",
					"lago_wallet_id":           nil,
					"wallet_code":              nil,
					"alert_type":               "billable_metric_lifetime_usage_units",
					"direction":                "increasing",
					"code":                     "units_alert",
					"previous_value":           "0.0",
					"thresholds":               []map[string]any{},
					"created_at":               "2025-03-20T10:00:00Z",
				},
			})
		defer server.Close()

		result, err := server.Client().Alert().Get(context.Background(), "sub_1234", "units_alert")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.AlertType, qt.Equals, BillableMetricLifetimeUsageUnitsAlertType)
		c.Assert(result.Direction, qt.Equals, IncreasingAlertDirection)
		c.Assert(result.LagoWalletID, qt.IsNil)
		c.Assert(result.WalletCode, qt.IsNil)
	})
}

func TestAlertRequest_GetList(t *testing.T) {
	t.Run("When alerts are found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MockResponse(mockAlertsResponse)
		defer server.Close()

		result, err := server.Client().Alert().GetList(context.Background(), "sub_1234")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Alerts, qt.HasLen, 2)
		c.Assert(result.Alerts[0].Code, qt.Equals, "alert1")
		c.Assert(result.Alerts[1].Code, qt.Equals, "alert2")
	})
}

func TestAlertRequest_Create(t *testing.T) {
	t.Run("When alert is created", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Create(context.Background(), "sub_1234", &AlertInput{
			Code:      "usage_alert",
			AlertType: CurrentUsageAmountAlertType,
			Thresholds: []AlertThreshold{
				{Code: "warn", Value: "1000"},
			},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_CreateList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Alert().CreateList(context.Background(), "sub_1234", []AlertInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When alerts are created", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MockResponse(mockAlertsResponse)
		defer server.Close()

		result, err := server.Client().Alert().CreateList(context.Background(), "sub_1234", []AlertInput{
			{
				Code:      "alert1",
				AlertType: CurrentUsageAmountAlertType,
				Thresholds: []AlertThreshold{
					{Code: "warn", Value: "1000"},
				},
			},
			{
				Code:               "alert2",
				AlertType:          BillableMetricCurrentUsageAmountAlertType,
				BillableMetricCode: "storage",
				Thresholds: []AlertThreshold{
					{Value: "2000"},
				},
			},
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.HasLen, 2)
		c.Assert(result[0].Code, qt.Equals, "alert1")
		c.Assert(result[0].AlertType, qt.Equals, CurrentUsageAmountAlertType)
		c.Assert(result[1].Code, qt.Equals, "alert2")
		c.Assert(result[1].AlertType, qt.Equals, BillableMetricCurrentUsageAmountAlertType)
	})
}

func TestAlertRequest_Update(t *testing.T) {
	t.Run("When alert is updated", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/usage_alert").
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Update(context.Background(), "sub_1234", "usage_alert", &AlertInput{
			Name: "Updated Alert",
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_Delete(t *testing.T) {
	t.Run("When alert is deleted", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/usage_alert").
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Delete(context.Background(), "sub_1234", "usage_alert")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_DeleteAll(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		err := client.Alert().DeleteAll(context.Background(), "sub_1234")
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When all alerts are deleted", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MockEmptyResponse()
		defer server.Close()

		err := server.Client().Alert().DeleteAll(context.Background(), "sub_1234")
		c.Assert(err == nil, qt.IsTrue)
	})
}

// --- Status parameter tests ---

func TestAlertRequest_GetWithStatus(t *testing.T) {
	t.Run("When alert is found with status parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/usage_alert").
			MatchQuery(map[string]string{"subscription_status": "pending"}).
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Get(context.Background(), "sub_1234", "usage_alert", "pending")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_GetListWithStatus(t *testing.T) {
	t.Run("When alerts are found with status parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MatchQuery(map[string]string{"subscription_status": "pending"}).
			MockResponse(mockAlertsResponse)
		defer server.Close()

		result, err := server.Client().Alert().GetList(context.Background(), "sub_1234", "pending")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Alerts, qt.HasLen, 2)
	})
}

func TestAlertRequest_CreateWithStatus(t *testing.T) {
	t.Run("When alert is created with status parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MatchQuery(map[string]string{"subscription_status": "pending"}).
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Create(context.Background(), "sub_1234", &AlertInput{
			Code:      "usage_alert",
			AlertType: CurrentUsageAmountAlertType,
			Thresholds: []AlertThreshold{
				{Code: "warn", Value: "1000"},
			},
		}, "pending")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_UpdateWithStatus(t *testing.T) {
	t.Run("When alert is updated with status parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/usage_alert").
			MatchQuery(map[string]string{"subscription_status": "pending"}).
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Update(context.Background(), "sub_1234", "usage_alert", &AlertInput{
			Name: "Updated Alert",
		}, "pending")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_DeleteWithStatus(t *testing.T) {
	t.Run("When alert is deleted with status parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts/usage_alert").
			MatchQuery(map[string]string{"subscription_status": "pending"}).
			MockResponse(mockAlertResponse)
		defer server.Close()

		result, err := server.Client().Alert().Delete(context.Background(), "sub_1234", "usage_alert", "pending")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "usage_alert")
	})
}

func TestAlertRequest_DeleteAllWithStatus(t *testing.T) {
	t.Run("When all alerts are deleted with status parameter", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/subscriptions/sub_1234/alerts").
			MatchQuery(map[string]string{"subscription_status": "pending"}).
			MockEmptyResponse()
		defer server.Close()

		err := server.Client().Alert().DeleteAll(context.Background(), "sub_1234", "pending")
		c.Assert(err == nil, qt.IsTrue)
	})
}

func TestResolvedAlert_Unmarshal(t *testing.T) {
	t.Run("When the alert is attached to a subscription", func(t *testing.T) {
		c := qt.New(t)

		payload := `{
			"lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"lago_alert_id": "3c903c90-3c90-3c90-3c90-3c903c903c90",
			"lago_subscription_id": "4d904d90-4d90-4d90-4d90-4d904d904d90",
			"external_subscription_id": "sub_1234567890",
			"lago_wallet_id": null,
			"external_customer_id": "cust_1234567890",
			"billable_metric_code": "storage_usage",
			"alert_type": "billable_metric_current_usage_amount",
			"alert_code": "storage_threshold_alert",
			"alert_name": "Storage Usage Alert",
			"current_value": "700.0",
			"previous_value": "1200.0",
			"recovered_thresholds": [
				{"code": "warn", "value": "1000.0", "recurring": false, "notify_on": ["triggered", "resolved"]}
			],
			"in_alarm_thresholds": ["warning"],
			"fully_resolved": false,
			"reason": "recovered",
			"resolved_at": "2025-05-19T10:04:21Z"
		}`

		var resolved ResolvedAlert
		err := json.Unmarshal([]byte(payload), &resolved)
		c.Assert(err, qt.IsNil)

		c.Assert(resolved.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(resolved.LagoAlertID.String(), qt.Equals, "3c903c90-3c90-3c90-3c90-3c903c903c90")
		c.Assert(resolved.LagoSubscriptionID.String(), qt.Equals, "4d904d90-4d90-4d90-4d90-4d904d904d90")
		c.Assert(*resolved.ExternalSubscriptionID, qt.Equals, "sub_1234567890")
		c.Assert(resolved.LagoWalletID, qt.IsNil)
		c.Assert(resolved.ExternalCustomerID, qt.Equals, "cust_1234567890")
		c.Assert(*resolved.BillableMetricCode, qt.Equals, "storage_usage")
		c.Assert(resolved.AlertType, qt.Equals, BillableMetricCurrentUsageAmountAlertType)
		c.Assert(*resolved.AlertName, qt.Equals, "Storage Usage Alert")
		c.Assert(resolved.CurrentValue, qt.Equals, "700.0")
		c.Assert(resolved.PreviousValue, qt.Equals, "1200.0")
		c.Assert(resolved.RecoveredThresholds, qt.HasLen, 1)
		c.Assert(resolved.RecoveredThresholds[0].NotifyOn, qt.DeepEquals, []string{"triggered", "resolved"})
		c.Assert(resolved.InAlarmThresholds, qt.DeepEquals, []string{"warning"})
		c.Assert(resolved.FullyResolved, qt.IsFalse)
		c.Assert(resolved.Reason, qt.Equals, AlertResolutionReasonRecovered)
		c.Assert(resolved.ResolvedAt.Format(time.RFC3339), qt.Equals, "2025-05-19T10:04:21Z")
	})

	t.Run("When the alert is attached to a wallet", func(t *testing.T) {
		c := qt.New(t)

		payload := `{
			"lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"lago_alert_id": "3c903c90-3c90-3c90-3c90-3c903c903c90",
			"lago_subscription_id": null,
			"external_subscription_id": null,
			"lago_wallet_id": "5e905e90-5e90-5e90-5e90-5e905e905e90",
			"external_customer_id": "cust_1234567890",
			"billable_metric_code": null,
			"alert_type": "wallet_balance_amount",
			"alert_code": "wallet_balance_alert",
			"alert_name": null,
			"current_value": "50.0",
			"previous_value": "10.0",
			"recovered_thresholds": [],
			"in_alarm_thresholds": [],
			"fully_resolved": true,
			"reason": "threshold_removed",
			"resolved_at": "2025-05-19T10:04:21Z"
		}`

		var resolved ResolvedAlert
		err := json.Unmarshal([]byte(payload), &resolved)
		c.Assert(err, qt.IsNil)

		c.Assert(resolved.LagoSubscriptionID, qt.IsNil)
		c.Assert(resolved.ExternalSubscriptionID, qt.IsNil)
		c.Assert(resolved.LagoWalletID.String(), qt.Equals, "5e905e90-5e90-5e90-5e90-5e905e905e90")
		c.Assert(resolved.BillableMetricCode, qt.IsNil)
		c.Assert(resolved.AlertName, qt.IsNil)
		c.Assert(resolved.AlertType, qt.Equals, WalletBalanceAmountAlertType)
		c.Assert(resolved.FullyResolved, qt.IsTrue)
		c.Assert(resolved.Reason, qt.Equals, AlertResolutionReasonThresholdRemoved)
	})
}

func TestTriggeredAlert_Unmarshal(t *testing.T) {
	t.Run("When the alert is attached to a subscription", func(t *testing.T) {
		c := qt.New(t)

		payload := `{
			"lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"lago_alert_id": "3c903c90-3c90-3c90-3c90-3c903c903c90",
			"lago_subscription_id": "4d904d90-4d90-4d90-4d90-4d904d904d90",
			"external_subscription_id": "ext-id",
			"lago_wallet_id": null,
			"external_customer_id": "cust-ext-id",
			"billable_metric_code": null,
			"alert_name": "General Alert",
			"alert_code": "first",
			"alert_type": "current_usage_amount",
			"current_value": "3000.0",
			"previous_value": "1000.0",
			"crossed_thresholds": [
				{"code": "warn", "value": "2000.0", "recurring": false},
				{"code": "repeat", "value": "2500.0", "recurring": true}
			],
			"triggered_at": "2000-01-01T12:00:00Z",
			"subscription_external_id": "ext-id",
			"customer_external_id": "cust-ext-id"
		}`

		var triggered TriggeredAlert
		err := json.Unmarshal([]byte(payload), &triggered)
		c.Assert(err, qt.IsNil)

		c.Assert(triggered.LagoSubscriptionID.String(), qt.Equals, "4d904d90-4d90-4d90-4d90-4d904d904d90")
		c.Assert(*triggered.ExternalSubscriptionID, qt.Equals, "ext-id")
		c.Assert(triggered.LagoWalletID, qt.IsNil)
		c.Assert(triggered.ExternalCustomerID, qt.Equals, "cust-ext-id")
		c.Assert(triggered.SubscriptionExternalID, qt.Equals, "ext-id")
		c.Assert(triggered.CustomerExternalID, qt.Equals, "cust-ext-id")
		c.Assert(triggered.AlertType, qt.Equals, CurrentUsageAmountAlertType)
		c.Assert(triggered.AlertName, qt.Equals, "General Alert")
		c.Assert(triggered.TriggeredAt, qt.Equals, time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC))
		c.Assert(triggered.CrossedThresholds, qt.HasLen, 2)
		c.Assert(*triggered.CrossedThresholds[0].Code, qt.Equals, "warn")
		c.Assert(triggered.CrossedThresholds[0].Value, qt.Equals, "2000.0")
		c.Assert(triggered.CrossedThresholds[0].Recurring, qt.IsFalse)
		c.Assert(*triggered.CrossedThresholds[1].Code, qt.Equals, "repeat")
		c.Assert(triggered.CrossedThresholds[1].Value, qt.Equals, "2500.0")
		c.Assert(triggered.CrossedThresholds[1].Recurring, qt.IsTrue)
	})

	t.Run("When the alert is attached to a wallet", func(t *testing.T) {
		c := qt.New(t)

		payload := `{
			"lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"lago_alert_id": "3c903c90-3c90-3c90-3c90-3c903c903c90",
			"lago_subscription_id": null,
			"external_subscription_id": null,
			"lago_wallet_id": "5e905e90-5e90-5e90-5e90-5e905e905e90",
			"external_customer_id": "wallet-cust-ext-id",
			"billable_metric_code": null,
			"alert_name": null,
			"alert_code": "wallet-alert",
			"alert_type": "wallet_balance_amount",
			"current_value": "50.0",
			"previous_value": "100.0",
			"crossed_thresholds": [
				{"code": null, "value": "60.0", "recurring": true}
			],
			"triggered_at": "2000-01-01T12:00:00Z",
			"subscription_external_id": null,
			"customer_external_id": "wallet-cust-ext-id"
		}`

		var triggered TriggeredAlert
		err := json.Unmarshal([]byte(payload), &triggered)
		c.Assert(err, qt.IsNil)

		c.Assert(triggered.LagoSubscriptionID, qt.Equals, uuid.Nil)
		c.Assert(triggered.ExternalSubscriptionID, qt.IsNil)
		c.Assert(triggered.LagoWalletID.String(), qt.Equals, "5e905e90-5e90-5e90-5e90-5e905e905e90")
		c.Assert(triggered.ExternalCustomerID, qt.Equals, "wallet-cust-ext-id")
		c.Assert(triggered.BillableMetricCode, qt.Equals, "")
		c.Assert(triggered.AlertName, qt.Equals, "")
		c.Assert(triggered.AlertType, qt.Equals, WalletBalanceAmountAlertType)
		c.Assert(triggered.CrossedThresholds, qt.HasLen, 1)
		c.Assert(triggered.CrossedThresholds[0].Code, qt.IsNil)
		c.Assert(triggered.CrossedThresholds[0].Value, qt.Equals, "60.0")
		c.Assert(triggered.CrossedThresholds[0].Recurring, qt.IsTrue)
	})
}
