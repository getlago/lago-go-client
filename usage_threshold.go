package lago

import (
	"time"

	"github.com/google/uuid"
)

type UsageThresholdInput struct {
	LagoId               *uuid.UUID `json:"id,omitempty"`
	ThresholdDisplayName string     `json:"threshold_display_name,omitempty"`
	AmountCents          int        `json:"amount_cents"`
	Recurring            bool       `json:"recurring"`
}

type UsageThreshold struct {
	LagoID               uuid.UUID `json:"lago_id"`
	ThresholdDisplayName string    `json:"threshold_display_name,omitempty"`
	AmountCents          int       `json:"amount_cents"`
	Recurring            bool      `json:"recurring"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type AppliedUsageThreshold struct {
	LifetimeUsageAmountCents int            `json:"lifetime_usage_amount_cents"`
	CreatedAt                time.Time      `json:"created_at"`
	UsageThreshold           UsageThreshold `json:"usage_threshold"`
}
