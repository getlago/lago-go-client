package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaxRequest struct {
	client *Client
}

type TaxParams struct {
	Tax *TaxInput `json:"tax"`
}

type TaxInput struct {
	Name                  string   `json:"name,omitempty"`
	Code                  string   `json:"code,omitempty"`
	Rate                  *float32 `json:"rate,omitempty"`
	Description           string   `json:"description,omitempty"`
	AppliedToOrganization bool     `json:"applied_to_organization,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

type TaxListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type TaxResult struct {
	Tax   *Tax     `json:"tax,omitempty"`
	Taxes []Tax    `json:"taxes,omitempty"`
	Meta  Metadata `json:"meta,omitempty"`
}

type Tax struct {
	LagoID                uuid.UUID `json:"lago_id,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Code                  string    `json:"code,omitempty"`
	Rate                  float32   `json:"rate,omitempty"`
	Description           string    `json:"description,omitempty"`
	CustomersCount        int       `json:"customers_count,omitempty"`
	PlansCount            int       `json:"plans_count,omitempty"`
	ChargesCount          int       `json:"charges_count,omitempty"`
	AppliedToOrganization bool      `json:"applied_to_organization,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
}

func (c *Client) Tax() *TaxRequest {
	return &TaxRequest{
		client: c,
	}
}

func (adr *TaxRequest) Get(ctx context.Context, taxCode string) (*Tax, *Error) {
	subPath := fmt.Sprintf("%s/%s", "taxes", taxCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxResult{},
	}

	result, err := adr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxResult, ok := result.(*TaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxResult.Tax, nil
}

func (adr *TaxRequest) GetList(ctx context.Context, taxListInput *TaxListInput) (*TaxResult, *Error) {
	jsonQueryparams, err := json.Marshal(taxListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "taxes",
		QueryParams: queryParams,
		Result:      &TaxResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	taxResult, ok := result.(*TaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxResult, nil
}

func (adr *TaxRequest) Create(ctx context.Context, taxInput *TaxInput) (*Tax, *Error) {
	taxParams := &TaxParams{
		Tax: taxInput,
	}

	clientRequest := &ClientRequest{
		Path:   "taxes",
		Result: &TaxResult{},
		Body:   taxParams,
	}

	result, err := adr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxResult, ok := result.(*TaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxResult.Tax, nil
}

func (adr *TaxRequest) Update(ctx context.Context, taxInput *TaxInput) (*Tax, *Error) {
	subPath := fmt.Sprintf("%s/%s", "taxes", taxInput.Code)
	taxParams := &TaxParams{
		Tax: taxInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxResult{},
		Body:   taxParams,
	}

	result, err := adr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxResult, ok := result.(*TaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxResult.Tax, nil
}

func (adr *TaxRequest) Delete(ctx context.Context, taxCode string) (*Tax, *Error) {
	subPath := fmt.Sprintf("%s/%s", "taxes", taxCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxResult{},
	}

	result, err := adr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	taxResult, ok := result.(*TaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return taxResult.Tax, nil
}
