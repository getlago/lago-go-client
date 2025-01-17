package lago

import (
	"time"

	"github.com/google/uuid"
)

type InvoicingReason string

const (
	BillingPeriodSubscriptionStarting    InvoicingReason = "subscription_starting"
	BillingPeriodSubscriptionPeriodic    InvoicingReason = "subscription_periodic"
	BillingPeriodSubscriptionTerminating InvoicingReason = "subscription_terminating"
	BillingPeriodSInAdvanceCharge        InvoicingReason = "in_advance_charge"
	BillingPeriodInAdvanceChargePeriodic InvoicingReason = "in_advance_charge_periodic"
	BillingPeriodSProgressiveBilling     InvoicingReason = "progressive_billing"
)

type BillingPeriod struct {
	LagoSubscriptionId       uuid.UUID       `json:"lago_subscription_id"`
	ExternalSubscriptionId   string          `json:"external_subscription_id"`
	LagoPlanId               uuid.UUID       `json:"lago_plan_id"`
	SubscriptionFromDatetime time.Time       `json:"subscription_from_datetime"`
	SubscriptionToDatetime   time.Time       `json:"subscription_to_datetime"`
	ChargesFromDatetime      time.Time       `json:"charges_from_datetime"`
	ChargesToDatetime        time.Time       `json:"charges_to_datetime"`
	InvoicingReason          InvoicingReason `json:"invoicing_reason"`
}
