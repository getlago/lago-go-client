package lago

import (
	"context"
	"time"
)

type OrganizationRequest struct {
	client *Client
}

type OrganizationParams struct {
	Organization *OrganizationInput `json:"organization"`
}

type OrganizationBillingConfigurationInput struct {
	InvoiceGracePeriod int     `json:"invoice_grace_period,omitempty"`
	InvoiceFooter      string  `json:"invoice_footer,omitempty"`
	VatRate            float32 `json:"vat_rate,omitempty"`
	DocumentLocale     string  `json:"document_locale,omitempty"`
}

type OrganizationBillingConfiguration struct {
	InvoiceGracePeriod int     `json:"invoice_grace_period,omitempty"`
	InvoiceFooter      string  `json:"invoice_footer,omitempty"`
	VatRate            float32 `json:"vat_rate,omitempty"`
	DocumentLocale     string  `json:"document_locale,omitempty"`
}

type OrganizationInput struct {
	Name string `json:"name,omitempty"`

	Email                            string   `json:"email,omitempty"`
	AddressLine1                     string   `json:"address_line1,omitempty"`
	AddressLine2                     string   `json:"address_line2,omitempty"`
	City                             string   `json:"city,omitempty"`
	Zipcode                          string   `json:"zipcode,omitempty"`
	State                            string   `json:"state,omitempty"`
	Country                          string   `json:"country,omitempty"`
	LegalName                        string   `json:"legal_name,omitempty"`
	LegalNumber                      string   `json:"legal_number,omitempty"`
	TaxIdentificationNumber          string   `json:"tax_identification_number,omitempty"`
	WebhookURL                       string   `json:"webhook_url,omitempty"`
	Timezone                         string   `json:"timezone,omitempty"`
	EmailSettings                    []string `json:"email_settings,omitempty"`

	BillingConfiguration OrganizationBillingConfigurationInput `json:"billing_configuration,omitempty"`
}

type OrganizationResult struct {
	Organization *Organization `json:"organization,omitempty"`
}

type Organization struct {
	Name string `json:"name,omitempty"`

	Email                            string                           `json:"email,omitempty"`
	AddressLine1                     string                           `json:"address_line1,omitempty"`
	AddressLine2                     string                           `json:"address_line2,omitempty"`
	City                             string                           `json:"city,omitempty"`
	Zipcode                          string                           `json:"zipcode,omitempty"`
	State                            string                           `json:"state,omitempty"`
	Country                          string                           `json:"country,omitempty"`
	LegalName                        string                           `json:"legal_name,omitempty"`
	LegalNumber                      string                           `json:"legal_number,omitempty"`
	TaxIdentificationNumber          string                           `json:"tax_identification_number,omitempty"`
	Timezone                         string                           `json:"timezone,omitempty"`
	EmailSettings                    []string                         `json:"email_settings,omitempty"`

	BillingConfiguration             OrganizationBillingConfiguration `json:"billing_configuration,omitempty"`

	Taxes []Tax `json:"taxes,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (c *Client) Organization() *OrganizationRequest {
	return &OrganizationRequest{
		client: c,
	}
}

func (or *OrganizationRequest) Update(ctx context.Context, organizationInput *OrganizationInput) (*Organization, *Error) {
	organizationParams := &OrganizationParams{
		Organization: organizationInput,
	}

	clientRequest := &ClientRequest{
		Path:   "organizations",
		Result: &OrganizationResult{},
		Body:   organizationParams,
	}

	result, err := or.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	organizationResult, ok := result.(*OrganizationResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return organizationResult.Organization, nil
}
