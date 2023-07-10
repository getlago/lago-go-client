package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PlanAppliedTaxRequest struct {
	client *Client
}

type PlanAppliedTaxParams struct {
	AppliedTax *PlanAppliedTaxInput `json:"applied_tax"`
}

type PlanAppliedTaxInput struct {
	TaxCode string `json:"tax_code,omitempty"`
}

type PlanAppliedTaxResult struct {
	AppliedTax *PlanAppliedTax `json:"applied_tax,omitempty"`
}

type PlanAppliedTax struct {
	LagoID     uuid.UUID `json:"lago_id,omitempty"`
	LagoPlanID uuid.UUID `json:"lago_plan_id,omitempty"`
	LagoTaxID  uuid.UUID `json:"lago_tax_id,omitempty"`

	PlanCode string `json:"plan_code,omitempty"`
	TaxCode  string `json:"tax_code,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (c *Client) PlanAppliedTax() *PlanAppliedTaxRequest {
	return &PlanAppliedTaxRequest{
		client: c,
	}
}

func (patr *PlanAppliedTaxRequest) Create(ctx context.Context, planCode string, appliedTaxInput *PlanAppliedTaxInput) (*PlanAppliedTax, *Error) {
	appliedTaxParams := &PlanAppliedTaxParams{
		AppliedTax: appliedTaxInput,
	}

	subPath := fmt.Sprintf("plans/%s/applied_taxes", planCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanAppliedTaxResult{},
		Body:   appliedTaxParams,
	}

	result, err := patr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedTaxResult, ok := result.(*PlanAppliedTaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedTaxResult.AppliedTax, nil
}

func (patr *PlanAppliedTaxRequest) Delete(ctx context.Context, planCode string, taxCode string) (*PlanAppliedTax, *Error) {
	subPath := fmt.Sprintf("plans/%s/applied_taxes/%s", planCode, taxCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanAppliedTaxResult{},
	}

	result, err := patr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedTaxResult, ok := result.(*PlanAppliedTaxResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedTaxResult.AppliedTax, nil
}
