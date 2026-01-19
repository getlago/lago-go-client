package lago

import (
	"time"

	"github.com/google/uuid"
)

type FixedChargeModel string

const (
	StandardFixedChargeModel  FixedChargeModel = "standard"
	GraduatedFixedChargeModel FixedChargeModel = "graduated"
	VolumeFixedChargeModel    FixedChargeModel = "volume"
)

type GraduatedRange struct {
	FromValue     int    `json:"from_value"`
	ToValue       *int   `json:"to_value"`
	FlatAmount    string `json:"flat_amount,omitempty"`
	PerUnitAmount string `json:"per_unit_amount,omitempty"`
}

type VolumeRange struct {
	FromValue     int    `json:"from_value"`
	ToValue       *int   `json:"to_value"`
	FlatAmount    string `json:"flat_amount,omitempty"`
	PerUnitAmount string `json:"per_unit_amount,omitempty"`
}

type FixedChargeProperties struct {
	Amount          *string          `json:"amount,omitempty"`
	GraduatedRanges []GraduatedRange `json:"graduated_ranges,omitempty"`
	VolumeRanges    []VolumeRange    `json:"volume_ranges,omitempty"`
}

type FixedCharge struct {
	LagoID             uuid.UUID              `json:"lago_id"`
	LagoAddOnID        uuid.UUID              `json:"lago_add_on_id"`
	InvoiceDisplayName string                 `json:"invoice_display_name"`
	AddOnCode          string                 `json:"add_on_code"`
	CreatedAt          time.Time              `json:"created_at"`
	Code               string                 `json:"code,omitempty"`
	ChargeModel        FixedChargeModel       `json:"charge_model"`
	PayInAdvance       bool                   `json:"pay_in_advance"`
	Prorated           bool                   `json:"prorated"`
	Properties         *FixedChargeProperties `json:"properties"`
	Units              float64                `json:"units"`
	LagoParentID       *uuid.UUID             `json:"lago_parent_id,omitempty"`
	Taxes              []Tax                  `json:"taxes,omitempty"`
}

type FixedChargeInput struct {
	LagoID                *uuid.UUID             `json:"id,omitempty"`
	AddOnID               uuid.UUID              `json:"add_on_id,omitempty"`
	ChargeModel           FixedChargeModel       `json:"charge_model,omitempty"`
	Code                  string                 `json:"code,omitempty"`
	InvoiceDisplayName    string                 `json:"invoice_display_name,omitempty"`
	Units                 float64                `json:"units"`
	PayInAdvance          bool                   `json:"pay_in_advance"`
	Prorated              bool                   `json:"prorated"`
	Properties            *FixedChargeProperties `json:"properties,omitempty"`
	TaxCodes              []string               `json:"tax_codes,omitempty"`
	ApplyUnitsImmediately bool                   `json:"apply_units_immediately,omitempty"`
}

type FixedChargeOverridesInput struct {
	ID                    *uuid.UUID             `json:"id,omitempty"`
	InvoiceDisplayName    string                 `json:"invoice_display_name,omitempty"`
	Units                 *float64               `json:"units,omitempty"`
	ApplyUnitsImmediately bool                   `json:"apply_units_immediately,omitempty"`
	Properties            *FixedChargeProperties `json:"properties,omitempty"`
	TaxCodes              []string               `json:"tax_codes,omitempty"`
}

type FixedChargeResult struct {
	FixedCharges []FixedCharge `json:"fixed_charges,omitempty"`
}
