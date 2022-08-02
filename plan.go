package lago

import (
	"encoding/json"
	"fmt"

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
	PlanWeekly  PlanInterval = "weekly"
	PlanMonthly PlanInterval = "monthly"
	PlanYearly  PlanInterval = "yearly"
)

type PlanChargeInput struct {
	BillableMetricID uuid.UUID         `json:"billable_metric_id,omitempty"`
	AmountCurrency   Currency          `json:"amount_currency,omitempty"`
	ChargeModel      ChargeModel       `json:"charge_model,omitempty"`
	Properties       map[string]string `json:"properties"`
}

type PlanInput struct {
	Name              string            `json:"name,omitempty"`
	Code              string            `json:"code,omitempty"`
	Interval          PlanInterval      `json:"interval,omitempty"`
	Description       string            `json:"description,omitempty"`
	AmountCents       int               `json:"amount_cents,omitempty"`
	AmountCurrency    Currency          `json:"amount_currency,omitempty"`
	PayInAdvance      bool              `json:"pay_in_advance,omitempty"`
	BillChargeMonthly bool              `json:"bill_charge_monthly,omitempty"`
	Charges           []PlanChargeInput `json:"charges,omitempty"`
}

type PlanListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type Plan struct {
	LagoID            uuid.UUID    `json:"lago_id"`
	Name              string       `json:"name,omitempty"`
	Code              string       `json:"code,omitempty"`
	Interval          PlanInterval `json:"interval,omitempty"`
	Description       string       `json:"description,omitempty"`
	AmountCents       int          `json:"amount_cents,omitempty"`
	AmountCurrency    Currency     `json:"amount_currency,omitempty"`
	PayInAdvance      bool         `json:"pay_in_advance,omitempty"`
	BillChargeMonthly bool         `json:"bill_charge_monthly,omitempty"`
	Charges           []Charge     `json:"charges,omitempty"`
}

func (c *Client) Plan() *PlanRequest {
	return &PlanRequest{
		client: c,
	}
}

func (pr *PlanRequest) Get(planCode string) (*Plan, *Error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanResult{},
	}

	result, err := pr.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	planResult := result.(*PlanResult)

	return planResult.Plan, nil
}

func (pr *PlanRequest) GetList(planListInput *PlanListInput) (*PlanResult, *Error) {
	jsonQueryParams, err := json.Marshal(planListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	json.Unmarshal(jsonQueryParams, &queryParams)

	clientRequest := &ClientRequest{
		Path:        "plans",
		QueryParams: queryParams,
		Result:      &PlanResult{},
	}

	result, clientErr := pr.client.Get(clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	planResult := result.(*PlanResult)

	return planResult, nil
}

func (pr *PlanRequest) Create(planInput *PlanInput) (*Plan, *Error) {
	planParams := &PlanParams{
		Plan: planInput,
	}

	clientRequest := &ClientRequest{
		Path:   "plans",
		Result: &PlanResult{},
		Body:   planParams,
	}

	result, err := pr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	planResult := result.(*PlanResult)

	return planResult.Plan, nil
}

func (pr *PlanRequest) Update(planInput *PlanInput) (*Plan, *Error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planInput.Code)
	planParams := &PlanParams{
		Plan: planInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanResult{},
		Body:   planParams,
	}

	result, err := pr.client.Put(clientRequest)
	if err != nil {
		return nil, err
	}

	planResult := result.(*PlanResult)

	return planResult.Plan, nil
}

func (pr *PlanRequest) Delete(planCode string) (*Plan, *Error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanResult{},
	}

	result, err := pr.client.Delete(clientRequest)
	if err != nil {
		return nil, err
	}

	planResult := result.(*PlanResult)

	return planResult.Plan, nil
}
