package lago

import (
	"time"

	"github.com/google/uuid"
)

type ChargeModel string

const (
	StandardChargeModel            ChargeModel = "standard"
	GraduatedChargeModel           ChargeModel = "graduated"
	GraduatedPercentageChargeModel ChargeModel = "graduated_percentage"
	PackageChargeModel             ChargeModel = "package"
	PercentageChargeModel          ChargeModel = "percentage"
	VolumeChargeModel              ChargeModel = "volume"
)

type ChargeFilter struct {
	InvoiceDisplayName string                 `json:"invoice_display_name,omitempty"`
	Properties         map[string]interface{} `json:"properties,omitempty"`
	Values             map[string]interface{} `json:"values,omitempty"`
}

type Charge struct {
	LagoID               uuid.UUID              `json:"lago_id,omitempty"`
	LagoBillableMetricID uuid.UUID              `json:"lago_billable_metric_id,omitempty"`
	BillableMetricCode   string                 `json:"billable_metric_code,omitempty"`
	ChargeModel          ChargeModel            `json:"charge_model,omitempty"`
	CreatedAt            time.Time              `json:"created_at,omitempty"`
	PayInAdvance         bool                   `json:"pay_in_advance,omitempty"`
	Invoiceable          bool                   `json:"invoiceable,omitempty"`
	InvoiceDisplayName   string                 `json:"invoice_display_name,omitempty"`
	Prorated             bool                   `json:"prorated,omitempty"`
	MinAmountCents       int                    `json:"min_amount_cents,omitempty"`
	Properties           map[string]interface{} `json:"properties,omitempty"`
	Filters              []ChargeFilter         `json:"filters,omitempty"`

	Taxes []Tax `json:"tax,omitempty"`
}
