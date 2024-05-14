package lago

import (
	"context"
	"encoding/json"
)

type PaymentStatus string

type InvoiceCollectionRequest struct {
	client *Client
}

type InvoiceCollectionListInput struct {
	AmountCurrency string `json:"currency,omitempty,string"`
	Months         int    `json:"months,omitempty,string"`
}

type InvoiceCollectionResult struct {
	InvoiceCollection  *InvoiceCollection  `json:"invoice_collection,omitempty"`
	InvoiceCollections []InvoiceCollection `json:"invoice_collections,omitempty"`
}

type InvoiceCollection struct {
	Month          string               `json:"month,omitempty"`
	PaymentStatus  InvoicePaymentStatus `json:"payment_status,omitempty"`
	InvoicesCount  int                  `json:"invoices_count,omitempty"`
	AmountCents    int                  `json:"amount_cents,omitempty"`
	AmountCurrency Currency             `json:"currency,omitempty"`
}

func (c *Client) InvoiceCollection() *InvoiceCollectionRequest {
	return &InvoiceCollectionRequest{
		client: c,
	}
}

func (adr *InvoiceCollectionRequest) GetList(ctx context.Context, InvoiceCollectionListInput *InvoiceCollectionListInput) (*InvoiceCollectionResult, *Error) {
	jsonQueryparams, err := json.Marshal(InvoiceCollectionListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/invoice_collection",
		QueryParams: queryParams,
		Result:      &InvoiceCollectionResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	InvoiceCollectionResult, ok := result.(*InvoiceCollectionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return InvoiceCollectionResult, nil
}
