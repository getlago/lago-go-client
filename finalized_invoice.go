package lago

import (
	"context"
	"encoding/json"
)

type PaymentStatus string

type FinalizedInvoiceRequest struct {
	client *Client
}

type FinalizedInvoiceListInput struct {
	AmountCurrency  string  `json:"currency,omitempty,string"`
	Months          int     `json:"months,omitempty,string"`
}

type FinalizedInvoiceResult struct {
	FinalizedInvoice  *FinalizedInvoice   `json:"finalized_invoice,omitempty"`
	FinalizedInvoices []FinalizedInvoice  `json:"finalized_invoices,omitempty"`
}

type FinalizedInvoice struct {
	Month			string   			  `json:"month,omitempty"`
	PaymentStatus	InvoicePaymentStatus  `json:"payment_status,omitempty"`
	InvoicesCount   int      			  `json:"invoices_count,omitempty"`
	AmountCents    	int      			  `json:"amount_cents,omitempty"`
	AmountCurrency	Currency 			  `json:"currency,omitempty"`
}

func (c *Client) FinalizedInvoice() *FinalizedInvoiceRequest {
	return &FinalizedInvoiceRequest{
		client: c,
	}
}

func (adr *FinalizedInvoiceRequest) GetList(ctx context.Context, FinalizedInvoiceListInput *FinalizedInvoiceListInput) (*FinalizedInvoiceResult, *Error) {
	jsonQueryparams, err := json.Marshal(FinalizedInvoiceListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/finalized_invoices",
		QueryParams: queryParams,
		Result:      &FinalizedInvoiceResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	FinalizedInvoiceResult, ok := result.(*FinalizedInvoiceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return FinalizedInvoiceResult, nil
}
