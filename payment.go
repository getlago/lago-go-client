package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

type ManualPaymentRequest struct {
	client *Client
}

type PaymentResult struct {
	Payment  *Payment  `json:"payment,omitempty"`
	Payments []Payment `json:"payments,omitempty"`
	Meta     Metadata  `json:"meta,omitempty"`
}

type PaymentListInput struct {
	PerPage *int `json:"per_page,omitempty,string" url:"per_page,omitempty"`
	Page    *int `json:"page,omitempty,string" url:"page,omitempty"`

	ExternalCustomerID  string   `json:"external_customer_id,omitempty" url:"external_customer_id,omitempty"`
	InvoiceID           string   `json:"invoice_id,omitempty" url:"invoice_id,omitempty"`
	PaymentStatus       []string `json:"payment_status,omitempty" url:"payment_status[],omitempty"`
	PaymentStatuses     []string `json:"payment_statuses,omitempty" url:"payment_statuses[],omitempty"`
	AmountFrom          *int64   `json:"amount_from,omitempty" url:"amount_from,omitempty"`
	AmountTo            *int64   `json:"amount_to,omitempty" url:"amount_to,omitempty"`
	ReceiptNumber       string   `json:"receipt_number,omitempty" url:"receipt_number,omitempty"`
	CreatedAtFrom       string   `json:"created_at_from,omitempty" url:"created_at_from,omitempty"`
	CreatedAtTo         string   `json:"created_at_to,omitempty" url:"created_at_to,omitempty"`
	PaymentProviderType []string `json:"payment_provider_type,omitempty" url:"payment_provider_type[],omitempty"`
	Currency            Currency `json:"currency,omitempty" url:"currency,omitempty"`
	InvoiceNumber       string   `json:"invoice_number,omitempty" url:"invoice_number,omitempty"`
	PaymentType         []string `json:"payment_type,omitempty" url:"payment_type[],omitempty"`
	PayableType         []string `json:"payable_type,omitempty" url:"payable_type[],omitempty"`
	SearchTerm          string   `json:"search_term,omitempty" url:"search_term,omitempty"`
}

type NextAction struct {
	Type          string         `json:"type,omitempty"`
	RedirectToURL *RedirectToURL `json:"redirect_to_url,omitempty"`
}

type RedirectToURL struct {
	URL       string `json:"url,omitempty"`
	ReturnURL string `json:"return_url,omitempty"`
}

type Payment struct {
	LagoID             uuid.UUID   `json:"lago_id,omitempty"`
	AmountCurrency     Currency    `json:"amount_currency,omitempty"`
	AmountCents        int         `json:"amount_cents,omitempty"`
	PaymentStatus      string      `json:"payment_status,omitempty"`
	Type               string      `json:"type,omitempty"`
	Reference          string      `json:"reference,omitempty"`
	ExternalPaymentID  string      `json:"external_payment_id,omitempty"`
	CreatedAt          time.Time   `json:"created_at,omitempty"`
	InvoiceIds         []string    `json:"invoice_ids,omitempty"`
	InvoiceNumbers     []string    `json:"invoice_numbers,omitempty"`
	ExternalCustomerID string      `json:"external_customer_id,omitempty"`
	NextAction         *NextAction `json:"next_action,omitempty"`
}

type PaymentParams struct {
	Payment *PaymentInput `json:"payment"`
}

type PaymentInput struct {
	InvoiceId   string `json:"invoice_id,omitempty"`
	AmountCents int    `json:"amount_cents,omitempty"`
	Reference   string `json:"reference,omitempty"`
	PaidAt      string `json:"paid_at,omitempty"`
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
	urlValues, err := query.Values(paymentListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      "payments",
		UrlValues: urlValues,
		Result:    &PaymentResult{},
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
