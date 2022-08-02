package lago

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	Active     SubscriptionStatus = "active"
	Pending    SubscriptionStatus = "pending"
	Terminated SubscriptionStatus = "terminated"
	Canceled   SubscriptionStatus = "canceled"
)

type SubscriptionRequest struct {
	client *Client
}

type SubscriptionResult struct {
	Subscription *Subscription `json:"subscription,omitempty"`
}

type SubscriptionParams struct {
	Subscription *SubscriptionInput `json:"subscription"`
}

type SubscriptionInput struct {
	CustomerID string `json:"customer_id,omitempty"`
	PlanCode   string `json:"plan_code,omitempty"`
}

type Subscription struct {
	LagoID         uuid.UUID `json:"lago_id"`
	LagoCustomerID uuid.UUID `json:"lago_customer_id"`
	CustomerID     string    `json:"customer_id"`

	PlanCode string `json:"plan_code"`

	Status SubscriptionStatus `json:"status"`

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

func (sr *SubscriptionRequest) Create(subscriptionInput *SubscriptionInput) (*Subscription, *Error) {
	subscriptionParam := &SubscriptionParams{
		Subscription: subscriptionInput,
	}

	clientRequest := &ClientRequest{
		Path:   "subscriptions",
		Result: &SubscriptionResult{},
		Body:   subscriptionParam,
	}

	result, err := sr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	subscriptionResult := result.(*SubscriptionResult)

	return subscriptionResult.Subscription, nil
}

func (sr *SubscriptionRequest) Terminate(customerID string) (*Subscription, *Error) {
	subscriptionInput := &SubscriptionInput{
		CustomerID: customerID,
	}

	clientRequest := &ClientRequest{
		Path:   "subscriptions",
		Result: &SubscriptionResult{},
		Body:   subscriptionInput,
	}

	result, err := sr.client.Delete(clientRequest)
	if err != nil {
		return nil, err
	}

	subscriptionResult := result.(*SubscriptionResult)

	return subscriptionResult.Subscription, nil
}
