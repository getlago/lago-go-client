package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PlanRequest struct {
	client *Client
}

type PlanResult struct {
	Plan  *Plan    `json:"plan,omitempty"`
	Plans []Plan   `json:"plans,omitempty"`
	Meta  Metadata `json:"meta,omitempty"`
}

type PlanParams struct {
	Plan *PlanInput `json:"plan"`
}

type PlanInterval string

const (
	PlanWeekly    PlanInterval = "weekly"
	PlanMonthly   PlanInterval = "monthly"
	PlanQuarterly PlanInterval = "quarterly"
	PlanYearly    PlanInterval = "yearly"
)

type PlanChargeInput struct {
	LagoID           *uuid.UUID             `json:"id,omitempty"`
	BillableMetricID uuid.UUID              `json:"billable_metric_id,omitempty"`
	AmountCurrency   Currency               `json:"amount_currency,omitempty"`
	ChargeModel      ChargeModel            `json:"charge_model,omitempty"`
	PayInAdvance     bool                   `json:"pay_in_advance,omitempty"`
	Invoiceable      bool                   `json:"invoiceable,omitempty"`
	RegroupPaidFees  string                 `json:"regroup_paid_fees,omitempty"`
	Prorated         bool                   `json:"prorated,omitempty"`
	MinAmountCents   int                    `json:"min_amount_cents,omitempty"`
	Properties       map[string]interface{} `json:"properties"`
	Filters          []ChargeFilter         `json:"filters,omitempty"`

	TaxCodes []string `json:"tax_codes,omitempty"`
}

type MinimumCommitmentInput struct {
	AmountCents        int      `json:"amount_cents,omitempty"`
	InvoiceDisplayName string   `json:"invoice_display_name,omitempty"`
	TaxCodes           []string `json:"tax_codes,omitempty"`
}

type PlanInput struct {
	Name               string                  `json:"name,omitempty"`
	InvoiceDisplayName string                  `json:"invoice_display_name,omitempty"`
	Code               string                  `json:"code,omitempty"`
	Interval           PlanInterval            `json:"interval,omitempty"`
	Description        string                  `json:"description,omitempty"`
	AmountCents        int                     `json:"amount_cents"`
	AmountCurrency     Currency                `json:"amount_currency,omitempty"`
	PayInAdvance       bool                    `json:"pay_in_advance"`
	BillChargesMonthly bool                    `json:"bill_charges_monthly"`
	TrialPeriod        float32                 `json:"trial_period"`
	Charges            []PlanChargeInput       `json:"charges,omitempty"`
	MinimumCommitment  *MinimumCommitmentInput `json:"minimum_commitment,omitempty"`
	TaxCodes           []string                `json:"tax_codes,omitempty"`
	UsageThresholds    []UsageThresholdInput   `json:"usage_thresholds,omitempty"`
	CascadeUpdates     bool                    `json:"cascade_updates"`
}

type PlanListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type MinimumCommitment struct {
	LagoID             uuid.UUID    `json:"lago_id"`
	PlanCode           string       `json:"plan_code,omitempty"`
	InvoiceDisplayName string       `json:"invoice_display_name,omitempty"`
	AmountCents        int          `json:"amount_cents"`
	Interval           PlanInterval `json:"interval,omitempty"`
	CreatedAt          time.Time    `json:"created_at,omitempty"`
	UpdatedAt          time.Time    `json:"updated_at,omitempty"`

	Taxes []Tax `json:"tax,omitempty"`
}

type Plan struct {
	LagoID             uuid.UUID          `json:"lago_id"`
	Name               string             `json:"name,omitempty"`
	InvoiceDisplayName string             `json:"invoice_display_name,omitempty"`
	Code               string             `json:"code,omitempty"`
	Interval           PlanInterval       `json:"interval,omitempty"`
	Description        string             `json:"description,omitempty"`
	AmountCents        int                `json:"amount_cents,omitempty"`
	AmountCurrency     Currency           `json:"amount_currency,omitempty"`
	PayInAdvance       bool               `json:"pay_in_advance,omitempty"`
	BillChargesMonthly bool               `json:"bill_charges_monthly,omitempty"`
	Charges            []Charge           `json:"charges,omitempty"`
	MinimumCommitment  *MinimumCommitment `json:"minimum_commitment"`

	Taxes           []Tax            `json:"taxes,omitempty"`
	UsageThresholds []UsageThreshold `json:"usage_thresholds,omitempty"`
}

func (c *Client) Plan() *PlanRequest {
	return &PlanRequest{
		client: c,
	}
}

func (pr *PlanRequest) Get(ctx context.Context, planCode string) (*Plan, *Error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanResult{},
	}

	result, err := pr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	planResult, ok := result.(*PlanResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planResult.Plan, nil
}

func (pr *PlanRequest) GetList(ctx context.Context, planListInput *PlanListInput) (*PlanResult, *Error) {
	jsonQueryParams, err := json.Marshal(planListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "plans",
		QueryParams: queryParams,
		Result:      &PlanResult{},
	}

	result, clientErr := pr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	planResult, ok := result.(*PlanResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planResult, nil
}

func (pr *PlanRequest) Create(ctx context.Context, planInput *PlanInput) (*Plan, *Error) {
	planParams := &PlanParams{
		Plan: planInput,
	}

	clientRequest := &ClientRequest{
		Path:   "plans",
		Result: &PlanResult{},
		Body:   planParams,
	}

	result, err := pr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	planResult, ok := result.(*PlanResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planResult.Plan, nil
}

func (pr *PlanRequest) Update(ctx context.Context, planInput *PlanInput) (*Plan, *Error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planInput.Code)
	planParams := &PlanParams{
		Plan: planInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanResult{},
		Body:   planParams,
	}

	result, err := pr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	planResult, ok := result.(*PlanResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planResult.Plan, nil
}

func (pr *PlanRequest) Delete(ctx context.Context, planCode string) (*Plan, *Error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanResult{},
	}

	result, err := pr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	planResult, ok := result.(*PlanResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planResult.Plan, nil
}
