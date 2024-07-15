package lago

import (
	"context"
	"encoding/json"
)

type OverdueBalanceRequest struct {
	client *Client
}

type OverdueBalanceListInput struct {
	AmountCurrency     string `json:"currency,omitempty"`
	ExternalCustomerId string `json:"external_customer_id,omitempty"`
	Months             int    `json:"months,omitempty,string"`
}

type OverdueBalanceResult struct {
	OverdueBalance  *OverdueBalance  `json:"overdue_balance,omitempty"`
	OverdueBalances []OverdueBalance `json:"overdue_balances,omitempty"`
}

type OverdueBalance struct {
	Month          string   `json:"month,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"currency,omitempty"`
}

func (c *Client) OverdueBalance() *OverdueBalanceRequest {
	return &OverdueBalanceRequest{
		client: c,
	}
}

func (adr *OverdueBalanceRequest) GetList(ctx context.Context, OverdueBalanceListInput *OverdueBalanceListInput) (*OverdueBalanceResult, *Error) {
	jsonQueryparams, err := json.Marshal(OverdueBalanceListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/overdue_balance",
		QueryParams: queryParams,
		Result:      &OverdueBalanceResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	OverdueBalanceResult, ok := result.(*OverdueBalanceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return OverdueBalanceResult, nil
}
