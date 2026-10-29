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
	BillableMetricLifetimeUsageUnitsAlertType AlertType = "billable_metric_lifetime_usage_units"
	WalletBalanceAmountAlertType              AlertType = "wallet_balance_amount"
	WalletCreditsBalanceAlertType             AlertType = "wallet_credits_balance"
	WalletOngoingBalanceAmountAlertType       AlertType = "wallet_ongoing_balance_amount"
	WalletCreditsOngoingBalanceAlertType      AlertType = "wallet_credits_ongoing_balance"
)

type AlertDirection string

const (
	IncreasingAlertDirection AlertDirection = "increasing"
	DecreasingAlertDirection AlertDirection = "decreasing"
)

type AlertResolutionReason string

const (
	AlertResolutionReasonRecovered        AlertResolutionReason = "recovered"
	AlertResolutionReasonThresholdUpdated AlertResolutionReason = "threshold_updated"
	AlertResolutionReasonThresholdRemoved AlertResolutionReason = "threshold_removed"
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

type AlertListParams struct {
	Alerts []AlertInput `json:"alerts"`
}

type AlertInput struct {
	Code               string           `json:"code,omitempty"`
	BillableMetricCode string           `json:"billable_metric_code,omitempty"`
	Name               string           `json:"name,omitempty"`
	AlertType          AlertType        `json:"alert_type,omitempty"`
	Thresholds         []AlertThreshold `json:"thresholds,omitempty"`
}

type AlertListInput struct {
	PerPage *int `json:"per_page,omitempty,string"`
	Page    *int `json:"page,omitempty,string"`
}

type AlertThreshold struct {
	Code      string   `json:"code,omitempty"`
	Value     string   `json:"value"`
	Recurring bool     `json:"recurring,omitempty"`
	NotifyOn  []string `json:"notify_on,omitempty"`
}

type Alert struct {
	LagoID                 uuid.UUID        `json:"lago_id"`
	LagoOrganizationID     uuid.UUID        `json:"lago_organization_id"`
	SubscriptionExternalID string           `json:"subscription_external_id"`
	ExternalSubscriptionID *string          `json:"external_subscription_id"`
	LagoWalletID           *uuid.UUID       `json:"lago_wallet_id"`
	WalletCode             *string          `json:"wallet_code"`
	BillableMetric         BillableMetric   `json:"billable_metric,omitempty"`
	AlertType              AlertType        `json:"alert_type"`
	Direction              AlertDirection   `json:"direction,omitempty"`
	Code                   string           `json:"code"`
	Name                   string           `json:"name,omitempty"`
	PreviousValue          string           `json:"previous_value,omitempty"`
	LastProcessedAt        *time.Time       `json:"last_processed_at,omitempty"`
	Thresholds             []AlertThreshold `json:"thresholds,omitempty"`
	CreatedAt              time.Time        `json:"created_at,omitempty"`
}

// AlertCrossedThreshold is an item of crossed_thresholds in the alert.triggered webhook.
// Code is nil for a recurring threshold that has no code.
type AlertCrossedThreshold struct {
	Code      *string `json:"code"`
	Value     string  `json:"value"`
	Recurring bool    `json:"recurring"`
}

// TriggeredAlert Object returned by the alert.triggered webhook.
type TriggeredAlert struct {
	LagoID                 uuid.UUID               `json:"lago_id,omitempty"`
	LagoOrganizationID     uuid.UUID               `json:"lago_organization_id"`
	LagoAlertID            uuid.UUID               `json:"lago_alert_id"`
	LagoSubscriptionID     uuid.UUID               `json:"lago_subscription_id"`
	ExternalSubscriptionID *string                 `json:"external_subscription_id"`
	LagoWalletID           *uuid.UUID              `json:"lago_wallet_id"`
	ExternalCustomerID     string                  `json:"external_customer_id"`
	SubscriptionExternalID string                  `json:"subscription_external_id"`
	CustomerExternalID     string                  `json:"customer_external_id"`
	BillableMetricCode     string                  `json:"billable_metric_code,omitempty"`
	AlertType              AlertType               `json:"alert_type"`
	AlertCode              string                  `json:"alert_code"`
	AlertName              string                  `json:"alert_name,omitempty"`
	CurrentValue           string                  `json:"current_value"`
	PreviousValue          string                  `json:"previous_value"`
	CrossedThresholds      []AlertCrossedThreshold `json:"crossed_thresholds"`
	TriggeredAt            time.Time               `json:"triggered_at"`
}

// ResolvedAlert Object returned by the alert.resolved webhook.
type ResolvedAlert struct {
	LagoID                 uuid.UUID             `json:"lago_id"`
	LagoOrganizationID     uuid.UUID             `json:"lago_organization_id"`
	LagoAlertID            uuid.UUID             `json:"lago_alert_id"`
	LagoSubscriptionID     *uuid.UUID            `json:"lago_subscription_id"`
	ExternalSubscriptionID *string               `json:"external_subscription_id"`
	LagoWalletID           *uuid.UUID            `json:"lago_wallet_id"`
	ExternalCustomerID     string                `json:"external_customer_id"`
	BillableMetricCode     *string               `json:"billable_metric_code"`
	AlertType              AlertType             `json:"alert_type"`
	AlertCode              string                `json:"alert_code"`
	AlertName              *string               `json:"alert_name"`
	CurrentValue           string                `json:"current_value"`
	PreviousValue          string                `json:"previous_value"`
	RecoveredThresholds    []AlertThreshold      `json:"recovered_thresholds"`
	InAlarmThresholds      []string              `json:"in_alarm_thresholds"`
	FullyResolved          bool                  `json:"fully_resolved"`
	Reason                 AlertResolutionReason `json:"reason"`
	ResolvedAt             time.Time             `json:"resolved_at"`
}

func (c *Client) Alert() *AlertRequest {
	return &AlertRequest{
		client: c,
	}
}

func (ar *AlertRequest) Get(ctx context.Context, subscriptionExternalID, alertCode string, subscriptionStatus ...string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
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

func (ar *AlertRequest) GetList(ctx context.Context, subscriptionExternalID string, subscriptionStatus ...string) (*AlertResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
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

func (ar *AlertRequest) Create(ctx context.Context, subscriptionExternalID string, alertInput *AlertInput, subscriptionStatus ...string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
		Body:        &AlertParams{Alert: alertInput},
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

func (ar *AlertRequest) Update(ctx context.Context, subscriptionExternalID, alertCode string, alertInput *AlertInput, subscriptionStatus ...string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
		Body:        &AlertParams{Alert: alertInput},
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

func (ar *AlertRequest) CreateList(ctx context.Context, subscriptionExternalID string, alertInputs []AlertInput, subscriptionStatus ...string) ([]Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
		Body:        &AlertListParams{Alerts: alertInputs},
	}

	result, err := ar.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alerts, nil
}

func (ar *AlertRequest) Delete(ctx context.Context, subscriptionExternalID, alertCode string, subscriptionStatus ...string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
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

func (ar *AlertRequest) DeleteAll(ctx context.Context, subscriptionExternalID string, subscriptionStatus ...string) *Error {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "alerts")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: statusQueryParams(subscriptionStatus),
		Result:      &AlertResult{},
	}

	_, err := ar.client.Delete(ctx, clientRequest)
	if err != nil {
		return err
	}

	return nil
}
