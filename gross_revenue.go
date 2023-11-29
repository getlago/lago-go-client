package lago

import (
	"context"
	"encoding/json"
)

type GrossRevenueRequest struct {
	client *Client
}

type GrossRevenueListInput struct {
	AmountCurrency      string  `json:"currency,omitempty,string"`
	ExternalCustomerId  string  `json:"external_customer_id,omitempty,string"`
	Months              int     `json:"months,omitempty,string"`
}

type GrossRevenueResult struct {
	GrossRevenue  *GrossRevenue   `json:"gross_revenue,omitempty"`
	GrossRevenues []GrossRevenue  `json:"gross_revenues,omitempty"`
}

type GrossRevenue struct {
	Month          string   `json:"month,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"currency,omitempty"`
}

func (c *Client) GrossRevenue() *GrossRevenueRequest {
	return &GrossRevenueRequest{
		client: c,
	}
}

func (adr *GrossRevenueRequest) GetList(ctx context.Context, GrossRevenueListInput *GrossRevenueListInput) (*GrossRevenueResult, *Error) {
	jsonQueryparams, err := json.Marshal(GrossRevenueListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/gross_revenue",
		QueryParams: queryParams,
		Result:      &GrossRevenueResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	GrossRevenueResult, ok := result.(*GrossRevenueResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return GrossRevenueResult, nil
}
