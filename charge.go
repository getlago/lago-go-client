package lago

import (
	"time"

	"github.com/google/uuid"
)

type ChargeModel string

const (
	StandardChargeModel   ChargeModel = "standard"
	GraduatedChargeModel  ChargeModel = "graduated"
	PackageChargeModel    ChargeModel = "package"
	PercentageChargeModel ChargeModel = "percentage"
	VolumneChargeModel    ChargeModel = "volume"
)

type Charge struct {
	LagoID               uuid.UUID              `json:"lago_id,omitempty"`
	LagoBillableMetricID uuid.UUID              `json:"lago_billable_metric_id,omitempty"`
	BillableMetricCode   string                 `json:"billable_metric_code,omitempty"`
	ChargeModel          ChargeModel            `json:"charge_model,omitempty"`
	CreatedAt            time.Time              `json:"created_at,omitempty"`
	PayInAdvance         bool                   `json:"pay_in_advance,omitempty"`
	Invoiceable          bool                   `json:"invoiceable,omitempty"`
	MinAmountCents       int                    `json:"min_amount_cents,omitempty"`
	Properties           map[string]interface{} `json:"properties,omitempty"`
	GroupProperties      []GroupProperties      `json:"group_properties,omitempty"`
}

type GroupProperties struct {
	GroupID uuid.UUID              `json:"group_id"`
	Values  map[string]interface{} `json:"values"`
}
