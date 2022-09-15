package lago

import (
	"context"
	"encoding/json"
	"time"

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

type SubscriptionInput struct {
	ExternalCustomerID string      `json:"external_customer_id,omitempty"`
	PlanCode           string      `json:"plan_code,omitempty"`
	BillingTime        BillingTime `json:"billing_time,omitempty"`
}

type SubscriptionListInput struct {
	ExternalCustomerID string `json:"external_customer_id,omitempty"`
	PerPage            int    `json:"per_page,omitempty,string"`
	Page               int    `json:"page,omitempty,string"`
}

type Subscription struct {
	LagoID             uuid.UUID `json:"lago_id"`
	LagoCustomerID     uuid.UUID `json:"lago_customer_id"`
	ExternalCustomerID string    `json:"external_customer_id"`

	PlanCode string `json:"plan_code"`

	Status           SubscriptionStatus `json:"status"`
	BillingTime      BillingTime        `json:"billing_time"`
	SubscriptionDate string             `json:"subscription_date"`

	PreviousPlanCode  string `json:"previous_plan_code"`
	NextPlanCode      string `json:"next_plan_code"`
	DowngradePlanDate string `json:"downgrade_plan_date"`

	CreatedAt    *time.Time `json:"created_at"`
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

	subscriptionResult := result.(*SubscriptionResult)

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) Terminate(ctx context.Context, externalCustomerID string) (*Subscription, *Error) {
	subscriptionInput := &SubscriptionInput{
		ExternalCustomerID: externalCustomerID,
	}

	clientRequest := &ClientRequest{
		Path:   "subscriptions",
		Result: &SubscriptionResult{},
		Body:   subscriptionInput,
	}

	result, err := sr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	subscriptionResult := result.(*SubscriptionResult)

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) GetList(ctx context.Context, subscriptionListInput SubscriptionListInput) (*SubscriptionResult, *Error) {
	jsonQueryParams, err := json.Marshal(subscriptionListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "subscriptions",
		QueryParams: queryParams,
		Result:      &SubscriptionResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	subscriptionResult, ok := result.(*SubscriptionResult)
	if !ok {
		return nil, &Error{Err: ErrorTypeAssert}
	}

	return subscriptionResult, nil
}
