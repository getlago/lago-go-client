package lago

import (
	"time"

	"github.com/google/uuid"
)

type ChargeModel string

const (
	Standard   ChargeModel = "standard"
	Graduated  ChargeModel = "graduated"
	Package    ChargeModel = "package"
	Percentage ChargeModel = "percentage"
)

type Charge struct {
	LagoID               uuid.UUID         `json:"lago_id,omitempty"`
	LagoBillableMetricID uuid.UUID         `json:"lago_billable_metric_id,omitempty"`
	ChargeModel          ChargeModel       `json:"charge_model,omitempty"`
	CreatedAt            time.Time         `json:"created_at,omitempty"`
	Properties           map[string]string `json:"properties,omitempty"`
}
