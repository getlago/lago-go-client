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
	LagoID          uuid.UUID  `json:"lago_id,omitempty"`
	Email           string     `json:"email,omitempty"`
	AmountCurrency  Currency   `json:"amount_currency,omitempty"`
	AmountCents     int        `json:"amount_cents,omitempty"`
	PaymentStatus   string     `json:"payment_status,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`

	Customer        *Customer  `json:"customer,omitempty"`
	Invoices        []Invoice  `json:"fees,omitempty"`
}

type PaymentRequestParams struct {
	PaymentRequest *PaymentRequestInput `json:"payment_request"`
}

type PaymentRequestInput struct {
	Email                 string     `json:"email,omitempty"`
	CustomerExternalId    string     `json:"customer_external_id,omitempty"`
	LagoInvoiceIds        []string   `json:"lago_invoice_ids,omitempty"`
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
