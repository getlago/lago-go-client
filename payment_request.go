package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentRequestRequest struct {
	client *Client
}

type PaymentRequestResult struct {
	PaymentRequest  *PaymentRequest  `json:"payment_request,omitempty"`
	PaymentRequests []PaymentRequest `json:"payment_requests,omitempty"`
	Meta            Metadata         `json:"meta,omitempty"`
}

type PaymentRequestListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	ExternalCustomerID string `json:"external_customer_id,omitempty"`
	PaymentStatus      string `json:"payment_status,omitempty"`
}

type PaymentRequest struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	Email          string    `json:"email,omitempty"`
	AmountCurrency Currency  `json:"amount_currency,omitempty"`
	AmountCents    int       `json:"amount_cents,omitempty"`
	PaymentStatus  string    `json:"payment_status,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`

	Customer *Customer `json:"customer,omitempty"`
	Invoices []Invoice `json:"fees,omitempty"`
}

type PaymentRequestParams struct {
	PaymentRequest *PaymentRequestInput `json:"payment_request"`
}

type PaymentRequestInput struct {
	Email              string   `json:"email,omitempty"`
	ExternalCustomerId string   `json:"external_customer_id,omitempty"`
	LagoInvoiceIds     []string `json:"lago_invoice_ids,omitempty"`
}

func (c *Client) PaymentRequest() *PaymentRequestRequest {
	return &PaymentRequestRequest{
		client: c,
	}
}

func (adr *PaymentRequestRequest) Get(ctx context.Context, paymentRequestID string) (*PaymentRequest, *Error) {
	subPath := fmt.Sprintf("%s/%s", "payment_requests", paymentRequestID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PaymentRequestResult{},
	}

	result, err := adr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentRequestResult, ok := result.(*PaymentRequestResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentRequestResult.PaymentRequest, nil
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

func (cr *PaymentRequestRequest) Create(ctx context.Context, paymentRequestInput *PaymentRequestInput) (*PaymentRequest, *Error) {
	paymentRequestParams := &PaymentRequestParams{
		PaymentRequest: paymentRequestInput,
	}

	clientRequest := &ClientRequest{
		Path:   "payment_requests",
		Result: &PaymentRequestResult{},
		Body:   paymentRequestParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentRequestResult, ok := result.(*PaymentRequestResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentRequestResult.PaymentRequest, nil
}
