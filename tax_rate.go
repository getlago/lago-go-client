package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaxRateRequest struct {
	client *Client
}

type TaxRateParams struct {
	TaxRate *TaxRateInput `json:"tax_rate"`
}

type TaxRateInput struct {
	Name             string  `json:"name,omitempty"`
	Code             string  `json:"code,omitempty"`
	Value            float32 `json:"value,omitempty"`
	Description      string  `json:"description,omitempty"`
	AppliedByDefault bool    `json:"applied_by_default,omitempty"`
}

type TaxRateListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type TaxRateResult struct {
	TaxRate  *TaxRate  `json:"tax_rate,omitempty"`
	TaxRates []TaxRate `json:"tax_rates,omitempty"`
	Meta     Metadata  `json:"meta,omitempty"`
}

type TaxRate struct {
	LagoID           uuid.UUID `json:"lago_id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Code             string    `json:"code,omitempty"`
	Value            float32   `json:"value,omitempty"`
	Description      string    `json:"description,omitempty"`
	CustomersCount   int       `json:"customers_count,omitempty"`
	AppliedByDefault bool      `json:"applied_by_default,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

func (c *Client) TaxRate() *TaxRateRequest {
	return &TaxRateRequest{
		client: c,
	}
}

func (adr *TaxRateRequest) Get(ctx context.Context, taxRateCode string) (*TaxRate, *Error) {
	subPath := fmt.Sprintf("%s/%s", "tax_rates", taxRateCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxRateResult{},
	}

	result, err := adr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxRateResult, ok := result.(*TaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxRateResult.TaxRate, nil
}

func (adr *TaxRateRequest) GetList(ctx context.Context, taxRateListInput *TaxRateListInput) (*TaxRateResult, *Error) {
	jsonQueryparams, err := json.Marshal(taxRateListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "tax_rates",
		QueryParams: queryParams,
		Result:      &TaxRateResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	taxRateResult, ok := result.(*TaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxRateResult, nil
}

func (adr *TaxRateRequest) Create(ctx context.Context, taxRateInput *TaxRateInput) (*TaxRate, *Error) {
	taxRateParams := &TaxRateParams{
		TaxRate: taxRateInput,
	}

	clientRequest := &ClientRequest{
		Path:   "tax_rates",
		Result: &TaxRateResult{},
		Body:   taxRateParams,
	}

	result, err := adr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxRateResult, ok := result.(*TaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxRateResult.TaxRate, nil
}

func (adr *TaxRateRequest) Update(ctx context.Context, taxRateInput *TaxRateInput) (*TaxRate, *Error) {
	subPath := fmt.Sprintf("%s/%s", "tax_rates", taxRateInput.Code)
	taxRateParams := &TaxRateParams{
		TaxRate: taxRateInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxRateResult{},
		Body:   taxRateParams,
	}

	result, err := adr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxRateResult, ok := result.(*TaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxRateResult.TaxRate, nil
}

func (adr *TaxRateRequest) Delete(ctx context.Context, taxRateCode string) (*TaxRate, *Error) {
	subPath := fmt.Sprintf("%s/%s", "tax_rates", taxRateCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxRateResult{},
	}

	result, err := adr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxRateResult, ok := result.(*TaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxRateResult.TaxRate, nil
}
