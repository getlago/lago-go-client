package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AddOnRequest struct {
	client *Client
}

type AddOnParams struct {
	AddOn *AddOnInput `json:"add_on"`
}

type AddOnInput struct {
	Name								string `json:"name,omitempty"`
	InvoiceDisplayName	string `json:"invoice_display_name,omitempty"`
	Code								string `json:"code,omitempty"`
	Description					string `json:"description,omitempty"`

	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	TaxCodes []string `json:"tax_codes,omitempty"`
}

type AddOnListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type AddOnResult struct {
	AddOn  *AddOn   `json:"add_on,omitempty"`
	AddOns []AddOn  `json:"add_ons,omitempty"`
	Meta   Metadata `json:"meta,omitempty"`
}

type AddOn struct {
	LagoID							uuid.UUID `json:"lago_id,omitempty"`
	Name								string    `json:"name,omitempty"`
	InvoiceDisplayName	string    `json:"invoice_display_name,omitempty"`
	Code								string    `json:"code,omitempty"`
	Description					string    `json:"description,omitempty"`

	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	Taxes []Tax `json:"tax,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (c *Client) AddOn() *AddOnRequest {
	return &AddOnRequest{
		client: c,
	}
}

func (adr *AddOnRequest) Get(ctx context.Context, addOnCode string) (*AddOn, *Error) {
	subPath := fmt.Sprintf("%s/%s", "add_ons", addOnCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AddOnResult{},
	}

	result, err := adr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult, ok := result.(*AddOnResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) GetList(ctx context.Context, addOnListInput *AddOnListInput) (*AddOnResult, *Error) {
	jsonQueryparams, err := json.Marshal(addOnListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "add_ons",
		QueryParams: queryParams,
		Result:      &AddOnResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	addOnResult, ok := result.(*AddOnResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return addOnResult, nil
}

func (adr *AddOnRequest) Create(ctx context.Context, addOnInput *AddOnInput) (*AddOn, *Error) {
	addOnParams := &AddOnParams{
		AddOn: addOnInput,
	}

	clientRequest := &ClientRequest{
		Path:   "add_ons",
		Result: &AddOnResult{},
		Body:   addOnParams,
	}

	result, err := adr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult, ok := result.(*AddOnResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) Update(ctx context.Context, addOnInput *AddOnInput) (*AddOn, *Error) {
	subPath := fmt.Sprintf("%s/%s", "add_ons", addOnInput.Code)
	addOnParams := &AddOnParams{
		AddOn: addOnInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AddOnResult{},
		Body:   addOnParams,
	}

	result, err := adr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult, ok := result.(*AddOnResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) Delete(ctx context.Context, addOnCode string) (*AddOn, *Error) {
	subPath := fmt.Sprintf("%s/%s", "add_ons", addOnCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AddOnResult{},
	}

	result, err := adr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult, ok := result.(*AddOnResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return addOnResult.AddOn, nil
}
