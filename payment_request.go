package lago

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PaymentRequestRequest struct {
	client *Client
}

type PaymentRequestResult struct {
	PaymentRequest  *PaymentRequest `json:"payment_request,omitempty"`
	PaymentRequests []PaymentRequest `json:"payment_requests,omitempty"`
}

type PaymentRequestListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	ExternalCustomerID string               `json:"external_customer_id,omitempty"`
}

type PaymentRequest struct {
	LagoID       uuid.UUID  `json:"lago_id,omitempty"`
	Email        string     `json:"email,omitempty"`
	Currency     Currency   `json:"currency,omitempty"`
	AmountCents  int        `json:"amount_cents,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`

	Customer     *Customer  `json:"customer,omitempty"`
	Invoices     []Invoice  `json:"fees,omitempty"`
}

func (c *Client) PaymentRequest() *PaymentRequestRequest {
	return &PaymentRequestRequest{
		client: c,
	}
}

func (ir *PaymentRequestRequest) GetList(ctx context.Context, paymentRequestListInput *PaymentRequestListInput) (*PaymentRequestResult, *Error) {
	jsonQueryParams, err := json.Marshal(paymentRequestListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "payment_requests",
		QueryParams: queryParams,
		Result:      &PaymentRequestResult{},
	}

	result, clientErr := ir.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	paymentRequestResult, ok := result.(*PaymentRequestResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentRequestResult, nil
}
