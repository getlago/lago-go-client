package lago

import (
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
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`

	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`
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
	LagoID      uuid.UUID `json:"lago_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`

	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AppliedAddOnResult struct {
	AppliedAddOn *AppliedAddOn `json:"applied_add_on,omitempty"`
}

type ApplyAddOnParams struct {
	AppliedAddOn *ApplyAddOnInput `json:"applied_add_on"`
}

type ApplyAddOnInput struct {
	CustomerID     string   `json:"customer_id,omitempty"`
	AddOnCode      string   `json:"add_on_code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`
}

type AppliedAddOn struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	LagoAddOnID    uuid.UUID `json:"lago_add_on_id,omitempty"`
	LagoCustomerID uuid.UUID `json:"lago_customer_id,omitempty"`
	CustomerID     string    `json:"customer_id,omitempty"`

	AddOnCode      string   `json:"add_on_code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (c *Client) AddOn() *AddOnRequest {
	return &AddOnRequest{
		client: c,
	}
}

func (adr *AddOnRequest) Get(addOnCode string) (*AddOn, *Error) {
	subPath := fmt.Sprintf("%s/%s", "add_ons", addOnCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AddOnResult{},
	}

	result, err := adr.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult := result.(*AddOnResult)

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) GetList(addOnListInput *AddOnListInput) (*AddOnResult, *Error) {
	jsonQueryparams, err := json.Marshal(addOnListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	json.Unmarshal(jsonQueryparams, &queryParams)

	clientRequest := &ClientRequest{
		Path:        "add_ons",
		QueryParams: queryParams,
		Result:      &AddOnResult{},
	}

	result, clientErr := adr.client.Get(clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	addOnResult := result.(*AddOnResult)

	return addOnResult, nil
}

func (adr *AddOnRequest) Create(addOnInput *AddOnInput) (*AddOn, *Error) {
	addOnParams := &AddOnParams{
		AddOn: addOnInput,
	}

	clientRequest := &ClientRequest{
		Path:   "add_ons",
		Result: &AddOnResult{},
		Body:   addOnParams,
	}

	result, err := adr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult := result.(*AddOnResult)

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) Update(addOnInput *AddOnInput) (*AddOn, *Error) {
	subPath := fmt.Sprintf("%s/%s", "add_ons", addOnInput.Code)
	addOnParams := &AddOnParams{
		AddOn: addOnInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AddOnResult{},
		Body:   addOnParams,
	}

	result, err := adr.client.Put(clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult := result.(*AddOnResult)

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) Delete(addOnCode string) (*AddOn, *Error) {
	subPath := fmt.Sprintf("%s/%s", "add_ons", addOnCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AddOnResult{},
	}

	result, err := adr.client.Delete(clientRequest)
	if err != nil {
		return nil, err
	}

	addOnResult := result.(*AddOnResult)

	return addOnResult.AddOn, nil
}

func (adr *AddOnRequest) ApplyToCustomer(applyAddOnInput *ApplyAddOnInput) (*AppliedAddOn, *Error) {
	applyAddOnParams := &ApplyAddOnParams{
		AppliedAddOn: applyAddOnInput,
	}

	clientRequest := &ClientRequest{
		Path:   "applied_add_ons",
		Result: &AppliedAddOnResult{},
		Body:   applyAddOnParams,
	}

	result, err := adr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	appliedAddOnResult := result.(*AppliedAddOnResult)

	return appliedAddOnResult.AppliedAddOn, nil
}
