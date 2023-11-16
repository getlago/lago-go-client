package lago

import (
	"context"
	"encoding/json"
)

type InvoicedUsageRequest struct {
	client *Client
}

type InvoicedUsageListInput struct {
	AmountCurrency  string	`json:"currency,omitempty,string"`
	Months          int 	`json:"months,omitempty,string"`
}

type InvoicedUsageResult struct {
	InvoicedUsage  *InvoicedUsage   `json:"invoiced_usage,omitempty"`
	InvoicedUsages []InvoicedUsage  `json:"invoiced_usages,omitempty"`
}

type InvoicedUsage struct {
	Month          string   `json:"month,omitempty"`
	Code           string   `json:"code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"currency,omitempty"`
}

func (c *Client) InvoicedUsage() *InvoicedUsageRequest {
	return &InvoicedUsageRequest{
    	client: c,
  	}
}

func (adr *InvoicedUsageRequest) GetList(ctx context.Context, InvoicedUsageListInput *InvoicedUsageListInput) (*InvoicedUsageResult, *Error) {
	jsonQueryparams, err := json.Marshal(InvoicedUsageListInput)
	if err != nil {
    	return nil, &Error{Err: err}
  	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

  clientRequest := &ClientRequest{
    Path:        "analytics/invoiced_usage",
    QueryParams: queryParams,
    Result:      &InvoicedUsageResult{},
  }

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	InvoicedUsageResult, ok := result.(*InvoicedUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

  	return InvoicedUsageResult, nil
}
