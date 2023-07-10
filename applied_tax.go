package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AppliedTaxRequest struct {
	client *Client
}

type AppliedTaxParams struct {
	AppliedTax *AppliedTaxInput `json:"applied_tax"`
}

type AppliedTaxInput struct {
	TaxCode string `json:"tax_code,omitempty"`
}

type AppliedTaxResult struct {
	AppliedTax *AppliedTax `json:"applied_tax,omitempty"`
	Meta       Metadata    `json:"meta,omitempty"`
}

type AppliedTax struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	LagoCustomerID uuid.UUID `json:"lago_customer_id,omitempty"`
	LagoTaxID      uuid.UUID `json:"lago_tax_id,omitempty"`

	TaxCode string `json:"tax_code,omitempty"`

	ExternalCustomerID string `json:"external_customer_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (c *Client) AppliedTax() *AppliedTaxRequest {
	return &AppliedTaxRequest{
		client: c,
	}
}

func (adr *AppliedTaxRequest) Create(ctx context.Context, externalCustomerID string, appliedTaxInput *AppliedTaxInput) (*AppliedTax, *Error) {
	appliedTaxParams := &AppliedTaxParams{
		AppliedTax: appliedTaxInput,
	}

	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "applied_taxes")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AppliedTaxResult{},
		Body:   appliedTaxParams,
	}

	result, err := adr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedTaxResult, ok := result.(*AppliedTaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedTaxResult.AppliedTax, nil
}

func (adr *AppliedTaxRequest) Delete(ctx context.Context, externalCustomerID string, taxCode string) (*AppliedTax, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "customers", externalCustomerID, "applied_taxes", taxCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AppliedTaxResult{},
	}

	result, err := adr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedTaxResult, ok := result.(*AppliedTaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedTaxResult.AppliedTax, nil
}
