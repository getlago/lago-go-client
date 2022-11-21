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
	InvoiceFooter string  `json:"invoice_footer,omitempty"`
	VatRate       float32 `json:"vat_rate,omitempty"`
}

type OrganizationBillingConfiguration struct {
	InvoiceFooter string  `json:"invoice_footer,omitempty"`
	VatRate       float32 `json:"vat_rate,omitempty"`
}

type OrganizationInput struct {
	Name string `json:"name,omitempty"`

	Email                string                                `json:"email,omitempty"`
	AddressLine1         string                                `json:"address_line_1,omitempty"`
	AddressLine2         string                                `json:"address_line_2,omitempty"`
	City                 string                                `json:"city,omitempty"`
	Zipcode              string                                `json:"zipcode,omitempty"`
	State                string                                `json:"state,omitempty"`
	Country              string                                `json:"country,omitempty"`
	LegalName            string                                `json:"legal_name,omitempty"`
	LegalNumber          string                                `json:"legal_number,omitempty"`
	WebhookURL           string                                `json:"webhook_url,omitempty"`
	BillingConfiguration OrganizationBillingConfigurationInput `json:"billing_configuration,omitempty"`
}

type OrganizationResult struct {
	Organization *Organization `json:"organization,omitempty"`
}

type Organization struct {
	Name string `json:"name,omitempty"`

	Email                string                           `json:"email,omitempty"`
	AddressLine1         string                           `json:"address_line_1,omitempty"`
	AddressLine2         string                           `json:"address_line_2,omitempty"`
	City                 string                           `json:"city,omitempty"`
	Zipcode              string                           `json:"zipcode,omitempty"`
	State                string                           `json:"state,omitempty"`
	Country              string                           `json:"country,omitempty"`
	LegalName            string                           `json:"legal_name,omitempty"`
	LegalNumber          string                           `json:"legal_number,omitempty"`
	BillingConfiguration OrganizationBillingConfiguration `json:"billing_configuration,omitempty"`

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

	organizationResult := result.(*OrganizationResult)

	return organizationResult.Organization, nil
}
