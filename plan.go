package lago

import (
	"fmt"

	"github.com/google/uuid"
)

type PlanRequest struct {
	client *Client
}

type PlanResult struct {
	Plan *Plan
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
	Name           string            `json:"name,omitempty"`
	Code           string            `json:"code,omitempty"`
	Interval       PlanInterval      `json:"interval,omitempty"`
	Description    string            `json:"description,omitempty"`
	AmountCents    int               `json:"amount_cents,omitempty"`
	AmountCurrency Currency          `json:"currency,omitempty"`
	PayInAdvance   bool              `json:"pay_in_advance,omitempty"`
	Charges        []PlanChargeInput `json:"charges,omitempty"`
}

type Plan struct {
	LagoID uuid.UUID `json:"lago_id"`
	Name   string    `json:"name"`
}

func (c *Client) Plan() *PlanRequest {
	return &PlanRequest{
		client: c,
	}
}

func (pr *PlanRequest) Get(planCode string) (*Plan, error) {
	subPath := fmt.Sprintf("%s/%s", "plans", planCode)

	resp, err := pr.client.HttpClient.
		R().
		SetResult(&PlanResult{}).
		Get(subPath)
	if err != nil {
		return nil, err
	}

	planResult := resp.Result().(*PlanResult)

	return planResult.Plan, nil
}

func (pr *PlanRequest) Create(planInput PlanInput) (*Plan, error) {
	resp, err := pr.client.HttpClient.
		R().
		SetResult(&PlanResult{}).
		Post("plans")
	if err != nil {
		return nil, err
	}

	planResult := resp.Result().(*PlanResult)

	return planResult.Plan, nil
}
