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

type Subscription struct {
	LagoID         uuid.UUID `json:"lago_id"`
	LagoCustomerID uuid.UUID `json:"lago_customer_id"`
	CustomerID     string    `json:"customer_id"`
	PlanCode       string    `json:"plan_code"`

	Status SubscriptionStatus `json:"status"`

	CreatedAt    *time.Time `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	CanceledAt   *time.Time `json:"canceled_at"`
	TerminatedAt *time.Time `json:"terminated_at"`
}

type SubscriptionInput struct {
	CustomerID string `json:"customer_id"`
	PlanCode   string `json:"plan_code"`
}
