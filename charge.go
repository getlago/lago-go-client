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
	DynamicChargeModel             ChargeModel = "dynamic"
)

type AppliedPricingUnit struct {
	Code           string  `json:"code,omitempty"`
	ConversionRate float64 `json:"conversion_rate,omitempty"`
}

type ChargeFilter struct {
	InvoiceDisplayName string                 `json:"invoice_display_name,omitempty"`
	Properties         map[string]interface{} `json:"properties,omitempty"`
	Values             map[string]interface{} `json:"values,omitempty"`
	CascadeUpdates     *bool                  `json:"cascade_updates,omitempty"`
}

type Charge struct {
	LagoID               uuid.UUID              `json:"lago_id,omitempty"`
	LagoBillableMetricID uuid.UUID              `json:"lago_billable_metric_id,omitempty"`
	BillableMetricCode   string                 `json:"billable_metric_code,omitempty"`
	Code                 string                 `json:"code,omitempty"`
	ChargeModel          ChargeModel            `json:"charge_model,omitempty"`
	CreatedAt            time.Time              `json:"created_at,omitempty"`
	PayInAdvance         bool                   `json:"pay_in_advance,omitempty"`
	Invoiceable          bool                   `json:"invoiceable,omitempty"`
	RegroupPaidFees      string                 `json:"regroup_paid_fees,omitempty"`
	InvoiceDisplayName   string                 `json:"invoice_display_name,omitempty"`
	Prorated             bool                   `json:"prorated,omitempty"`
	MinAmountCents       int                    `json:"min_amount_cents,omitempty"`
	Properties           map[string]interface{} `json:"properties,omitempty"`
	Filters              []ChargeFilter         `json:"filters,omitempty"`
	AppliedPricingUnit   *AppliedPricingUnit    `json:"applied_pricing_unit,omitempty"`
	AcceptsTargetWallet  *bool                  `json:"accepts_target_wallet,omitempty"`

	Taxes []Tax `json:"tax,omitempty"`
}

// Standalone charge CRUD types

type ChargeInput struct {
	BillableMetricID    string                 `json:"billable_metric_id,omitempty"`
	Code                string                 `json:"code,omitempty"`
	ChargeModel         ChargeModel            `json:"charge_model,omitempty"`
	PayInAdvance        bool                   `json:"pay_in_advance,omitempty"`
	Invoiceable         bool                   `json:"invoiceable,omitempty"`
	RegroupPaidFees     string                 `json:"regroup_paid_fees,omitempty"`
	InvoiceDisplayName  string                 `json:"invoice_display_name,omitempty"`
	Prorated            bool                   `json:"prorated,omitempty"`
	MinAmountCents      int                    `json:"min_amount_cents,omitempty"`
	Properties          map[string]interface{} `json:"properties"`
	Filters             []ChargeFilter         `json:"filters,omitempty"`
	TaxCodes            []string               `json:"tax_codes,omitempty"`
	AppliedPricingUnit  *AppliedPricingUnit    `json:"applied_pricing_unit,omitempty"`
	AcceptsTargetWallet *bool                  `json:"accepts_target_wallet,omitempty"`
	CascadeUpdates      *bool                  `json:"cascade_updates,omitempty"`
}

type ChargeParams struct {
	Charge *ChargeInput `json:"charge"`
}

type ChargeResult struct {
	Charge  *Charge  `json:"charge,omitempty"`
	Charges []Charge `json:"charges,omitempty"`
	Meta    Metadata `json:"meta,omitempty"`
}

type ChargeListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

// Standalone charge filter CRUD types

type ChargeFilterInput struct {
	InvoiceDisplayName string                 `json:"invoice_display_name,omitempty"`
	Properties         map[string]interface{} `json:"properties,omitempty"`
	Values             map[string]interface{} `json:"values,omitempty"`
	CascadeUpdates     *bool                  `json:"cascade_updates,omitempty"`
}

type ChargeFilterParams struct {
	Filter *ChargeFilterInput `json:"filter"`
}

type ChargeFilterResponse struct {
	LagoID             uuid.UUID              `json:"lago_id,omitempty"`
	InvoiceDisplayName string                 `json:"invoice_display_name,omitempty"`
	Properties         map[string]interface{} `json:"properties,omitempty"`
	Values             map[string]interface{} `json:"values,omitempty"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
}

type ChargeFilterResult struct {
	Filter  *ChargeFilterResponse  `json:"filter,omitempty"`
	Filters []ChargeFilterResponse `json:"filters,omitempty"`
	Meta    Metadata               `json:"meta,omitempty"`
}

type ChargeFilterListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}
