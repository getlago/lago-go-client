package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventRequest struct {
	client *Client
}

type EventParams struct {
	Event *EventInput `json:"event"`
}

type BatchEventParams struct {
	Events []EventInput `json:"events"`
}

type EventInput struct {
	TransactionID           string                 `json:"transaction_id,omitempty"`
	ExternalSubscriptionID  string                 `json:"external_subscription_id,omitempty"`
	Code                    string                 `json:"code,omitempty"`
	Timestamp               string                 `json:"timestamp,omitempty"`
	PreciseTotalAmountCents string                 `json:"precise_total_amount_cents,omitempty"`
	Properties              map[string]interface{} `json:"properties,omitempty"`
}

type EventEstimateFeesParams struct {
	Event *EventEstimateFeesInput `json:"event"`
}

type EventEstimateFeesInput struct {
	ExternalSubscriptionID string            `json:"external_subscription_id,omitempty"`
	Code                   string            `json:"code,omitempty"`
	Properties             map[string]string `json:"properties,omitempty"`
}

type BatchEventResult struct {
	Events []Event `json:"events"`
}

type EventResult struct {
	Event *Event `json:"event"`
}

type Event struct {
	LagoID                  uuid.UUID              `json:"lago_id"`
	TransactionID           string                 `json:"transaction_id"`
	LagoCustomerID          *uuid.UUID             `json:"lago_customer_id,omitempty"`
	Code                    string                 `json:"code,omitempty"`
	Timestamp               time.Time              `json:"timestamp"`
	PreciseTotalAmountCents string                 `json:"precise_total_amount_cents,omitempty"`
	Properties              map[string]interface{} `json:"properties,omitempty"`
	LagoSubscriptionID      *uuid.UUID             `json:"lago_subscription_id,omitempty"`
	ExternalSubscriptionID  string                 `json:"external_subscription_id,omitempty"`
	CreatedAt               time.Time              `json:"created_at"`
}

func (c *Client) Event() *EventRequest {
	return &EventRequest{
		client: c,
	}
}

func (er *EventRequest) Create(ctx context.Context, eventInput *EventInput) (*Event, *Error) {
	eventParams := &EventParams{
		Event: eventInput,
	}

	clientRequest := &ClientRequest{
		UseIngestService: true,
		Path:             "events",
		Result:           &EventResult{},
		Body:             eventParams,
	}

	result, err := er.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	eventResult, ok := result.(*EventResult)
	if !ok {
		return nil, err
	}

	return eventResult.Event, nil
}

func (er *EventRequest) EstimateFees(ctx context.Context, estimateInput EventEstimateFeesInput) (*FeeResult, *Error) {
	eventEstimateParams := &EventEstimateFeesParams{
		Event: &estimateInput,
	}

	clientRequest := &ClientRequest{
		Path:   "events/estimate_fees",
		Result: &FeeResult{},
		Body:   eventEstimateParams,
	}

	result, clientErr := er.client.Post(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	feeResult, ok := result.(*FeeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return feeResult, nil
}

func (er *EventRequest) Get(ctx context.Context, eventID string) (*Event, *Error) {
	subPath := fmt.Sprintf("%s/%s", "events", eventID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &EventResult{},
	}

	result, err := er.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	eventResult, ok := result.(*EventResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return eventResult.Event, nil
}

func (er *EventRequest) Batch(ctx context.Context, batchInput []EventInput) ([]Event, *Error) {
	eventParams := &BatchEventParams{
		Events: batchInput,
	}

	clientRequest := &ClientRequest{
		UseIngestService: true,
		Path:             "events/batch",
		Result:           &BatchEventResult{},
		Body:             eventParams,
	}

	result, err := er.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	batchEventResult, ok := result.(*BatchEventResult)
	if !ok {
		return nil, err
	}

	return batchEventResult.Events, nil
}
