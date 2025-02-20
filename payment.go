package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ManualPaymentRequest struct {
	client *Client
}

type PaymentResult struct {
	Payment  *Payment  `json:"payment,omitempty"`
	Payments []Payment `json:"payments,omitempty"`
}

type PaymentListInput struct {
	PerPage 			int `json:"per_page,omitempty,string"`
	Page				int `json:"page,omitempty,string"`

	ExternalCustomerID	string `json:"external_customer_id,omitempty"`
	InvoiceID 			string `json:"invoice_id,omitempty"`
}

type Payment struct {
	LagoID         		uuid.UUID `json:"lago_id,omitempty"`
	AmountCurrency		Currency  `json:"amount_currency,omitempty"`
	AmountCents			int       `json:"amount_cents,omitempty"`
	PaymentStatus		string    `json:"payment_status,omitempty"`
	Type  		   		string    `json:"type,omitempty"`
	Reference			string    `json:"reference,omitempty"`
	ExternalPaymentID	string    `json:"external_payment_id,omitempty"`
	CreatedAt      		time.Time `json:"created_at,omitempty"`
	InvoiceIds 	   		[]string  `json:"invoice_ids,omitempty"`
}

type PaymentParams struct {
	Payment *PaymentInput `json:"payment"`
}

type PaymentInput struct {
	InvoiceId	string		`json:"invoice_id,omitempty"`
	AmountCents int   		`json:"amount_cents,omitempty"`
	Reference	string		`json:"reference,omitempty"`
	PaidAt		string		`json:"paid_at,omitempty"`
}

func (c *Client) Payment() *ManualPaymentRequest {
	return &ManualPaymentRequest{
		client: c,
	}
}

func (adr *ManualPaymentRequest) Get(ctx context.Context, paymentID string) (*Payment, *Error) {
	subPath := fmt.Sprintf("%s/%s", "payments", paymentID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PaymentResult{},
	}

	result, err := adr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentResult, ok := result.(*PaymentResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentResult.Payment, nil
}

func (ir *ManualPaymentRequest) GetList(ctx context.Context, paymentListInput *PaymentListInput) (*PaymentResult, *Error) {
	jsonQueryParams, err := json.Marshal(paymentListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "payments",
		QueryParams: queryParams,
		Result:      &PaymentResult{},
	}

	result, clientErr := ir.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	paymentResult, ok := result.(*PaymentResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentResult, nil
}

func (cr *ManualPaymentRequest) Create(ctx context.Context, paymentInput *PaymentInput) (*Payment, *Error) {
	paymentParams := &PaymentParams{
		Payment: paymentInput,
	}

	clientRequest := &ClientRequest{
		Path:   "payments",
		Result: &PaymentResult{},
		Body:   paymentParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentResult, ok := result.(*PaymentResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentResult.Payment, nil
}
