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
)

type Charge struct {
	LagoID               uuid.UUID              `json:"lago_id,omitempty"`
	LagoBillableMetricID uuid.UUID              `json:"lago_billable_metric_id,omitempty"`
	ChargeModel          ChargeModel            `json:"charge_model,omitempty"`
	CreatedAt            time.Time              `json:"created_at,omitempty"`
	Properties           map[string]interface{} `json:"properties,omitempty"`
	GroupProperties      []GroupProperties      `json:"group_properties,omitempty"`
}

type GroupProperties struct {
	GroupId uuid.UUID              `json:"group_id"`
	Values  map[string]interface{} `json:"values"`
}
