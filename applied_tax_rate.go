package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AppliedTaxRateRequest struct {
	client *Client
}

type AppliedTaxRateParams struct {
	AppliedTaxRate *AppliedTaxRateInput `json:"applied_tax_rate"`
}

type AppliedTaxRateInput struct {
	TaxRateCode string `json:"tax_rate_code,omitempty"`
}

type AppliedTaxRateResult struct {
	AppliedTaxRate *AppliedTaxRate `json:"applied_tax_rate,omitempty"`
	Meta           Metadata        `json:"meta,omitempty"`
}

type AppliedTaxRate struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	LagoCustomerID uuid.UUID `json:"lago_customer_id,omitempty"`
	LagoTaxRateID  uuid.UUID `json:"lago_tax_rate_id,omitempty"`

	TaxRateCode string `json:"tax_rate_code,omitempty"`

	ExternalCustomerID string `json:"external_customer_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (c *Client) AppliedTaxRate() *AppliedTaxRateRequest {
	return &AppliedTaxRateRequest{
		client: c,
	}
}

func (adr *AppliedTaxRateRequest) Create(ctx context.Context, externalCustomerID string, appliedTaxRateInput *AppliedTaxRateInput) (*AppliedTaxRate, *Error) {
	appliedTaxRateParams := &AppliedTaxRateParams{
		AppliedTaxRate: appliedTaxRateInput,
	}

	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "applied_tax_rates")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &TaxRateResult{},
		Body:   appliedTaxRateParams,
	}

	result, err := adr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedTaxRateResult, ok := result.(*AppliedTaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedTaxRateResult.AppliedTaxRate, nil
}

func (adr *AppliedTaxRateRequest) Delete(ctx context.Context, externalCustomerID string, taxRateCode string) (*AppliedTaxRate, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "customers", externalCustomerID, "applied_tax_rates", taxRateCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AppliedTaxRateResult{},
	}

	result, err := adr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedTaxRateResult, ok := result.(*AppliedTaxRateResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedTaxRateResult.AppliedTaxRate, nil
}
