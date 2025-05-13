package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentReceiptRequest struct {
	client *Client
}

type PaymentReceiptResult struct {
	PaymentReceipt  *PaymentReceipt  `json:"payment_receipt,omitempty"`
	PaymentReceipts []PaymentReceipt `json:"payment_receipts,omitempty"`
	Meta            Metadata         `json:"meta,omitempty"`
}

type PaymentReceiptListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	InvoiceID string `json:"invoice_id,omitempty"`
}

type PaymentReceipt struct {
	LagoID    uuid.UUID `json:"lago_id,omitempty"`
	Number    string    `json:"number,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Payment   *Payment  `json:"payment,omitempty"`
}

func (c *Client) PaymentReceipt() *PaymentReceiptRequest {
	return &PaymentReceiptRequest{
		client: c,
	}
}

func (adr *PaymentReceiptRequest) Get(ctx context.Context, id string) (*PaymentReceipt, *Error) {
	subPath := fmt.Sprintf("%s/%s", "payment_receipts", id)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PaymentReceiptResult{},
	}

	result, err := adr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult, ok := result.(*PaymentReceiptResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return addOnResult.PaymentReceipt, nil
}

func (ir *PaymentReceiptRequest) GetList(ctx context.Context, paymentReceiptListInput *PaymentReceiptListInput) (*PaymentReceiptResult, *Error) {
	jsonQueryParams, err := json.Marshal(paymentReceiptListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "payment_receipts",
		QueryParams: queryParams,
		Result:      &PaymentReceiptResult{},
	}

	result, clientErr := ir.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	paymentReceiptResult, ok := result.(*PaymentReceiptResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentReceiptResult, nil
}
