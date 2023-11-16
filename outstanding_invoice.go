package lago

import (
	"context"
	"encoding/json"
)

type PaymentStatus string

type OutstandingInvoiceRequest struct {
	client *Client
}

type OutstandingInvoiceListInput struct {
	AmountCurrency  string  `json:"currency,omitempty,string"`
	Months          int     `json:"months,omitempty,string"`
}

type OutstandingInvoiceResult struct {
	OutstandingInvoice  *OutstandingInvoice   `json:"outstanding_invoice,omitempty"`
	OutstandingInvoices []OutstandingInvoice  `json:"outstanding_invoices,omitempty"`
}

type OutstandingInvoice struct {
	Month			string   			  `json:"month,omitempty"`
	PaymentStatus	InvoicePaymentStatus  `json:"payment_status,omitempty"`
	InvoicesCount   int      			  `json:"invoices_count,omitempty"`
	AmountCents    	int      			  `json:"amount_cents,omitempty"`
	AmountCurrency	Currency 			  `json:"currency,omitempty"`
}

func (c *Client) OutstandingInvoice() *OutstandingInvoiceRequest {
	return &OutstandingInvoiceRequest{
		client: c,
	}
}

func (adr *OutstandingInvoiceRequest) GetList(ctx context.Context, OutstandingInvoiceListInput *OutstandingInvoiceListInput) (*OutstandingInvoiceResult, *Error) {
	jsonQueryparams, err := json.Marshal(OutstandingInvoiceListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/outstanding_invoices",
		QueryParams: queryParams,
		Result:      &OutstandingInvoiceResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	OutstandingInvoiceResult, ok := result.(*OutstandingInvoiceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return OutstandingInvoiceResult, nil
}
