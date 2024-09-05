package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LifetimeUsage struct {
	LagoID                 uuid.UUID `json:"lago_id"`
	LagoSubscriptionID     uuid.UUID `json:"lago_subscription_id"`
	ExternalSubscriptionID string    `json:"external_subscription_id"`

	ExternalHistoricalUsageAmountCents *int `json:"external_historical_usage_amount_cents,omitempty"`
	InvoicedUsageAmountCents           *int `json:"invoiced_usage_amount_cents,omitempty"`
	CurrentUsageAmountCents            *int `json:"current_usage_amount_cents,omitempty"`

	FromDatetime time.Time `json:"from_datetime"`
	ToDatetime   time.Time `json:"to_datetime"`

	UsageThresholds []LifetimeUsageThreshold `json:"usage_thresholds,omitempty"`
}

type LifetimeUsageThreshold struct {
	AmountCents     int        `json:"amount_cents"`
	CompletionRatio float32    `json:"completion_ratio"`
	ReachedAt       *time.Time `json:"reached_at"`
}

type LifetimeUsageResult struct {
	LifetimeUsage *LifetimeUsage `json:"lifetime_usage"`
}

type LifetimeUsageParams struct {
	LifetimeUsage *LifetimeUsageInput `json:"lifetime_usage"`
}

type LifetimeUsageInput struct {
	ExternalSubscriptionID             string `json:"external_subscription_id"`
	ExternalHistoricalUsageAmountCents int    `json:"external_historical_usage_amount_cents"`
}

func (sr *SubscriptionRequest) GetLifetimeUsage(ctx context.Context, externalSubscriptionID string) (*LifetimeUsage, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", externalSubscriptionID, "lifetime_usage")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &LifetimeUsageResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	lifetimeUsageResult, ok := result.(*LifetimeUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return lifetimeUsageResult.LifetimeUsage, nil
}

func (sr *SubscriptionRequest) UpdateLifetimeUsage(ctx context.Context, lifetimeUsageInput *LifetimeUsageInput) (*LifetimeUsage, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", lifetimeUsageInput.ExternalSubscriptionID, "lifetime_usage")

	lifetimeUsageParams := &LifetimeUsageParams{
		LifetimeUsage: lifetimeUsageInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &LifetimeUsageResult{},
		Body:   lifetimeUsageParams,
	}

	result, clientErr := sr.client.Put(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	lifetimeUsageResult, ok := result.(*LifetimeUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return lifetimeUsageResult.LifetimeUsage, nil
}
