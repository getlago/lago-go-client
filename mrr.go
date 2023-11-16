package lago

import (
	"context"
	"encoding/json"
)

type MrrRequest struct {
	client *Client
}

type MrrListInput struct {
	AmountCurrency  string	`json:"currency,omitempty,string"`
	Months          int 	`json:"months,omitempty,string"`
}

type MrrResult struct {
	Mrr  *Mrr   `json:"mrr,omitempty"`
	Mrrs []Mrr  `json:"mrrs,omitempty"`
}

type Mrr struct {
	Month          string   `json:"month,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"currency,omitempty"`
}

func (c *Client) Mrr() *MrrRequest {
	return &MrrRequest{
		client: c,
	}
}

func (adr *MrrRequest) GetList(ctx context.Context, MrrListInput *MrrListInput) (*MrrResult, *Error) {
	jsonQueryparams, err := json.Marshal(MrrListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/mrr",
		QueryParams: queryParams,
		Result:      &MrrResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	MrrResult, ok := result.(*MrrResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return MrrResult, nil
}
