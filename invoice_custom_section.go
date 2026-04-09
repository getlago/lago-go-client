package lago

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceCustomSection struct {
	LagoId      uuid.UUID `json:"lago_id,omitempty"`
	Code        string    `json:"code,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Details     string    `json:"details,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
}

type InvoiceCustomSectionInput struct {
	InvoiceCustomSectionCodes []string `json:"invoice_custom_section_codes,omitempty"`
	SkipInvoiceCustomSections bool     `json:"skip_invoice_custom_sections,omitempty"`
}

type AppliedInvoiceCustomSection struct {
	LagoId               uuid.UUID            `json:"lago_id,omitempty"`
	CreatedAt            time.Time            `json:"created_at,omitempty"`
	InvoiceCustomSection InvoiceCustomSection `json:"invoice_custom_section,omitempty"`
}
