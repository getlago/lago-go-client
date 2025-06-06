package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AlertType string

const (
	CurrentUsageAmountAlertType               AlertType = "current_usage_amount"
	BillableMetricCurrentUsageAmountAlertType AlertType = "billable_metric_current_usage_amount"
	BillableMetricCurrentUsageUnitsAlertType  AlertType = "billable_metric_current_usage_units"
	LifetimeUsageAmountAlertType              AlertType = "lifetime_usage_amount"
)

type AlertRequest struct {
	client *Client
}

type AlertResult struct {
	Alert  *Alert   `json:"alert,omitempty"`
	Alerts []Alert  `json:"alerts,omitempty"`
	Meta   Metadata `json:"meta,omitempty"`
}

type AlertParams struct {
	Alert *AlertInput `json:"alert"`
}

type AlertInput struct {
	Code               string           `json:"code,omitempty"`
	BillableMetricCode string           `json:"billable_metric_code,omitempty"`
	Name               string           `json:"name,omitempty"`
	AlertType          AlertType        `json:"alert_type,omitempty"`
	Thresholds         []AlertThreshold `json:"thresholds,omitempty"`
}

type AlertThreshold struct {
	Code      string `json:"code,omitempty"`
	Value     string `json:"value"`
	Recurring bool   `json:"recurring,omitempty"`
}

type Alert struct {
	LagoID                 uuid.UUID        `json:"lago_id"`
	LagoOrganizationID     uuid.UUID        `json:"lago_organization_id"`
	SubscriptionExternalID string           `json:"subscription_external_id"`
	BillableMetric         BillableMetric   `json:"billable_metric,omitempty"`
	AlertType              AlertType        `json:"alert_type"`
	Code                   string           `json:"code"`
	Name                   string           `json:"name,omitempty"`
	PreviousValue          string           `json:"previous_value,omitempty"`
	LastProcessedAt        *time.Time       `json:"last_processed_at,omitempty"`
	Thresholds             []AlertThreshold `json:"thresholds,omitempty"`
	CreatedAt              time.Time        `json:"created_at,omitempty"`
}

// TriggeredAlert Object returned by the alert.triggered webhook.
type TriggeredAlert struct {
	LagoID                 uuid.UUID        `json:"lago_id,omitempty"`
	LagoOrganizationID     uuid.UUID        `json:"lago_organization_id"`
	LagoAlertID            uuid.UUID        `json:"lago_alert_id"`
	LagoSubscriptionID     uuid.UUID        `json:"lago_subscription_id"`
	SubscriptionExternalID string           `json:"subscription_external_id"`
	BillableMetricCode     string           `json:"billable_metric_code,omitempty"`
	AlertType              AlertType        `json:"alert_type"`
	AlertCode              string           `json:"alert_code"`
	AlertName              string           `json:"alert_name,omitempty"`
	CurrentValue           string           `json:"current_value"`
	PreviousValue          string           `json:"previous_value"`
	LastProcessedAt        *time.Time       `json:"last_processed_at"`
	CrossedThresholds      []AlertThreshold `json:"crossed_thresholds"`
	TriggeredAt            time.Time        `json:"triggered_at"`
}

func (c *Client) Alert() *AlertRequest {
	return &AlertRequest{
		client: c,
	}
}

func (ar *AlertRequest) Get(ctx context.Context, subscriptionExternalID, alertCode string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
	}

	result, err := ar.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}

func (ar *AlertRequest) GetList(ctx context.Context, subscriptionExternalID string) (*AlertResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
	}

	result, err := ar.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult, nil
}

func (ar *AlertRequest) Create(ctx context.Context, subscriptionExternalID string, alertInput *AlertInput) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
		Body:   &AlertParams{Alert: alertInput},
	}

	result, err := ar.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}

func (ar *AlertRequest) Update(ctx context.Context, subscriptionExternalID, alertCode string, alertInput *AlertInput) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
		Body:   &AlertParams{Alert: alertInput},
	}

	result, err := ar.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}

func (ar *AlertRequest) Delete(ctx context.Context, subscriptionExternalID, alertCode string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
	}

	result, err := ar.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}
