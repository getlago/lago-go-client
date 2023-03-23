package lago

import "context"

type EventRequest struct {
	client *Client
}

type EventParams struct {
	Event *EventInput `json:"event"`
}

type EventInput struct {
	TransactionID          string            `json:"transaction_id,omitempty"`
	ExternalCustomerID     string            `json:"external_customer_id,omitempty"`
	ExternalSubscriptionID string            `json:"external_subscription_id,omitempty"`
	Code                   string            `json:"code,omitempty"`
	Timestamp              int64             `json:"timestamp,omitempty"`
	Properties             map[string]string `json:"properties,omitempty"`
}

type EventEstimateFeesParams struct {
	Event *EventEstimateFeesInput `json:"event"`
}

type EventEstimateFeesInput struct {
	ExternalCustomerID     string            `json:"external_customer_id,omitempty"`
	ExternalSubscriptionID string            `json:"external_subscription_id,omitempty"`
	Code                   string            `json:"code,omitempty"`
	Properties             map[string]string `json:"properties,omitempty"`
}

func (c *Client) Event() *EventRequest {
	return &EventRequest{
		client: c,
	}
}

func (er *EventRequest) Create(ctx context.Context, eventInput *EventInput) *Error {
	eventParams := &EventParams{
		Event: eventInput,
	}

	clientRequest := &ClientRequest{
		Path: "events",
		Body: eventParams,
	}

	err := er.client.PostWithoutResult(ctx, clientRequest)
	if err != nil {
		return err
	}

	return nil
}

func (er *EventRequest) EstimateFees(ctx context.Context, estimateInput EventEstimateFeesInput) (*FeeResult, *Error) {
	eventEstimateParams := &EventEstimateFeesParams{
		Event: &estimateInput,
	}

	clientRequest := &ClientRequest{
		Path: "events/estimate_fees",
		Body: eventEstimateParams,
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
