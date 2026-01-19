package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive     SubscriptionStatus = "active"
	SubscriptionStatusPending    SubscriptionStatus = "pending"
	SubscriptionStatusTerminated SubscriptionStatus = "terminated"
	SubscriptionStatusCanceled   SubscriptionStatus = "canceled"
)

type BillingTime string

const (
	Anniversary BillingTime = "anniversary"
	Calendar    BillingTime = "calendar"
)

type OnTerminationCreditNote string

const (
	OnTerminationCreditNoteCredit OnTerminationCreditNote = "credit"
	OnTerminationCreditNoteRefund OnTerminationCreditNote = "refund"
	OnTerminationCreditNoteSkip   OnTerminationCreditNote = "skip"
)

type OnTerminationInvoice string

const (
	OnTerminationInvoiceGenerate OnTerminationInvoice = "generate"
	OnTerminationInvoiceSkip     OnTerminationInvoice = "skip"
)

type SubscriptionRequest struct {
	client *Client
}

type SubscriptionResult struct {
	Subscription  *Subscription  `json:"subscription,omitempty"`
	Subscriptions []Subscription `json:"subscriptions,omitempty"`
	Meta          Metadata       `json:"meta,omitempty"`
}

type SubscriptionParams struct {
	Subscription *SubscriptionInput `json:"subscription"`
}

type ChargeOverridesInput struct {
	ID                 *uuid.UUID             `json:"id,omitempty"`
	AmountCurrency     Currency               `json:"amount_currency,omitempty"`
	InvoiceDisplayName string                 `json:"invoice_display_name,omitempty"`
	MinAmountCents     int                    `json:"min_amount_cents,omitempty"`
	Properties         map[string]interface{} `json:"properties"`
	Filters            []ChargeFilter         `json:"filters,omitempty"`
	AppliedPricingUnit *AppliedPricingUnit    `json:"applied_pricing_unit,omitempty"`
	TaxCodes           []string               `json:"tax_codes,omitempty"`
}

type MinimumCommitmentOverridesInput struct {
	AmountCents        int      `json:"amount_cents,omitempty"`
	InvoiceDisplayName string   `json:"invoice_display_name,omitempty"`
	TaxCodes           []string `json:"tax_codes,omitempty"`
}

type PlanOverridesInput struct {
	Name               string                           `json:"name,omitempty"`
	InvoiceDisplayName string                           `json:"invoice_display_name,omitempty"`
	Code               string                           `json:"code,omitempty"`
	Description        string                           `json:"description,omitempty"`
	AmountCents        int                              `json:"amount_cents"`
	AmountCurrency     Currency                         `json:"amount_currency,omitempty"`
	TrialPeriod        float32                          `json:"trial_period"`
	Charges            []ChargeOverridesInput           `json:"charges,omitempty"`
	FixedCharges       []FixedChargeOverridesInput      `json:"fixed_charges,omitempty"`
	MinimumCommitment  *MinimumCommitmentOverridesInput `json:"minimum_commitment"`
	TaxCodes           []string                         `json:"tax_codes,omitempty"`
	UsageThresholds    []UsageThreshold                 `json:"usage_thresholds,omitempty"`
}

type SubscriptionInput struct {
	ExternalCustomerID string              `json:"external_customer_id,omitempty"`
	PlanCode           string              `json:"plan_code,omitempty"`
	SubscriptionAt     *time.Time          `json:"subscription_at,omitempty"`
	EndingAt           *time.Time          `json:"ending_at,omitempty"`
	BillingTime        BillingTime         `json:"billing_time,omitempty"`
	PlanOverrides      *PlanOverridesInput `json:"plan_overrides,omitempty"`
	ExternalID         string              `json:"external_id"`
	Name               string              `json:"name"`
}

type SubscriptionsInput struct {
	ExternalIds  []string   `json:"external_ids,omitempty"`
	TerminatedAt *time.Time `json:"terminated_at,omitempty"`
	PlanCode     string     `json:"plan_code,omitempty"`
}

type SubscriptionTerminateInput struct {
	ExternalID              string                  `json:"-"`
	Status                  string                  `json:"status,omitempty"`
	OnTerminationCreditNote OnTerminationCreditNote `json:"on_termination_credit_note,omitempty"`
	OnTerminationInvoice    OnTerminationInvoice    `json:"on_termination_invoice,omitempty"`
}

type SubscriptionListInput struct {
	ExternalCustomerID string               `url:"external_customer_id,omitempty"`
	PlanCode           string               `url:"plan_code,omitempty"`
	PerPage            *int                 `url:"per_page,omitempty"`
	Page               *int                 `url:"page,omitempty"`
	Status             []SubscriptionStatus `url:"status[],omitempty"`
}

type Subscription struct {
	LagoID             uuid.UUID `json:"lago_id"`
	LagoCustomerID     uuid.UUID `json:"lago_customer_id"`
	ExternalCustomerID string    `json:"external_customer_id"`
	ExternalID         string    `json:"external_id"`

	PlanCode string `json:"plan_code"`

	Name string `json:"name"`

	Status                  SubscriptionStatus      `json:"status"`
	BillingTime             BillingTime             `json:"billing_time"`
	SubscriptionAt          time.Time               `json:"subscription_at"`
	EndingAt                *time.Time              `json:"ending_at"`
	TrialEndedAt            *time.Time              `json:"trial_ended_at"`
	OnTerminationCreditNote OnTerminationCreditNote `json:"on_termination_credit_note,omitempty"`
	OnTerminationInvoice    OnTerminationInvoice    `json:"on_termination_invoice,omitempty"`

	PreviousPlanCode  string `json:"previous_plan_code"`
	NextPlanCode      string `json:"next_plan_code"`
	DowngradePlanDate string `json:"downgrade_plan_date"`

	CurrentBillingPeriodStartedAt *time.Time `json:"current_billing_period_started_at"`
	CurrentBillingPeriodEndingAt  *time.Time `json:"current_billing_period_ending_at"`

	Plan *Plan `json:"plan,omitempty"`

	Entitlements []SubscriptionEntitlement `json:"entitlements,omitempty"`

	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	CanceledAt   *time.Time `json:"canceled_at"`
	TerminatedAt *time.Time `json:"terminated_at"`
}

func (c *Client) Subscription() *SubscriptionRequest {
	return &SubscriptionRequest{
		client: c,
	}
}

func (sr *SubscriptionRequest) Create(ctx context.Context, subscriptionInput *SubscriptionInput) (*Subscription, *Error) {
	subscriptionParam := &SubscriptionParams{
		Subscription: subscriptionInput,
	}

	clientRequest := &ClientRequest{
		Path:   "subscriptions",
		Result: &SubscriptionResult{},
		Body:   subscriptionParam,
	}

	result, err := sr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	subscriptionResult, ok := result.(*SubscriptionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) Terminate(ctx context.Context, subscriptionTerminateInput SubscriptionTerminateInput) (*Subscription, *Error) {
	jsonQueryParams, err := json.Marshal(subscriptionTerminateInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	subPath := fmt.Sprintf("%s/%s", "subscriptions", subscriptionTerminateInput.ExternalID)

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &SubscriptionResult{},
	}

	result, clientErr := sr.client.Delete(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	subscriptionResult, ok := result.(*SubscriptionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) Get(ctx context.Context, subscriptionExternalId string) (*Subscription, *Error) {
	subPath := fmt.Sprintf("%s/%s", "subscriptions", subscriptionExternalId)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &SubscriptionResult{},
	}

	result, err := sr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	subscriptionResult, ok := result.(*SubscriptionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) GetList(ctx context.Context, subscriptionListInput SubscriptionListInput) (*SubscriptionResult, *Error) {
	urlValues, err := query.Values(subscriptionListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      "subscriptions",
		UrlValues: urlValues,
		Result:    &SubscriptionResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	subscriptionResult, ok := result.(*SubscriptionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionResult, nil
}

func (sr *SubscriptionRequest) Update(ctx context.Context, subscriptionInput *SubscriptionInput) (*Subscription, *Error) {
	subPath := fmt.Sprintf("%s/%s", "subscriptions", subscriptionInput.ExternalID)
	subscriptionParam := &SubscriptionParams{
		Subscription: subscriptionInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &SubscriptionResult{},
		Body:   subscriptionParam,
	}

	result, err := sr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	subscriptionResult, ok := result.(*SubscriptionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) GetFixedCharges(ctx context.Context, subscriptionExternalID string) (*FixedChargeResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalID, "fixed_charges")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
	}

	result, err := sr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult, nil
}
