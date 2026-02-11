package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockAlertResponse = map[string]any{
	"alert": map[string]any{
		"lago_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id":     "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"subscription_external_id": "sub_1234",
		"alert_type":               "current_usage_amount",
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

func TestAlertRequest_BatchCreate(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Alert().BatchCreate(context.Background(), "sub_1234", []AlertInput{})
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

		result, err := server.Client().Alert().BatchCreate(context.Background(), "sub_1234", []AlertInput{
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
			MockResponse(nil)
		defer server.Close()

		err := server.Client().Alert().DeleteAll(context.Background(), "sub_1234")
		c.Assert(err == nil, qt.IsTrue)
	})
}
