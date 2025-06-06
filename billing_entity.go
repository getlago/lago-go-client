package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BillingEntityDocumentNumbering string

const (
	BillingEntityDocumentNumberingPerCustomer      BillingEntityDocumentNumbering = "per_customer"
	BillingEntityDocumentNumberingPerBillingEntity BillingEntityDocumentNumbering = "per_billing_entity"
)

type BillingEntityRequest struct {
	client *Client
}

type BillingEntityCreateParams struct {
	BillingEntity *BillingEntityCreateInput `json:"billing_entity"`
}

type BillingEntityUpdateParams struct {
	BillingEntity *BillingEntityUpdateInput `json:"billing_entity"`
}

type BillingEntityResult struct {
	BillingEntity   *BillingEntity  `json:"billing_entity,omitempty"`
	BillingEntities []BillingEntity `json:"billing_entities,omitempty"`
	Meta            Metadata        `json:"meta,omitempty"`
}

type BillingEntityBillingConfiguration struct {
	InvoiceGracePeriod int    `json:"invoice_grace_period,omitempty"`
	InvoiceFooter      string `json:"invoice_footer,omitempty"`
	DocumentLocale     string `json:"document_locale,omitempty"`
}

type BillingEntity struct {
	LagoID                    uuid.UUID                      `json:"lago_id,omitempty"`
	Name                      string                         `json:"name,omitempty"`
	Code                      string                         `json:"code,omitempty"`
	Email                     string                         `json:"email,omitempty"`
	AddressLine1              string                         `json:"address_line1,omitempty"`
	AddressLine2              string                         `json:"address_line2,omitempty"`
	City                      string                         `json:"city,omitempty"`
	Zipcode                   string                         `json:"zipcode,omitempty"`
	State                     string                         `json:"state,omitempty"`
	Country                   string                         `json:"country,omitempty"`
	DefaultCurrency           Currency                       `json:"default_currency,omitempty"`
	LegalName                 string                         `json:"legal_name,omitempty"`
	LegalNumber               string                         `json:"legal_number,omitempty"`
	DocumentNumbering         BillingEntityDocumentNumbering `json:"document_numbering,omitempty"`
	DocumentNumberPrefix      string                         `json:"document_number_prefix,omitempty"`
	NetPaymentTerm            int                            `json:"net_payment_term,omitempty"`
	CreatedAt                 time.Time                      `json:"created_at,omitempty"`
	UpdatedAt                 time.Time                      `json:"updated_at,omitempty"`
	IsDefault                 bool                           `json:"is_default,omitempty"`
	Timezone                  string                         `json:"timezone,omitempty"`
	EmailSettings             []string                       `json:"email_settings,omitempty"`
	TaxIdentificationNumber   string                         `json:"tax_identification_number,omitempty"`
	FinalizeZeroAmountInvoice bool                           `json:"finalize_zero_amount_invoice,omitempty"`
	EuTaxManagement           bool                           `json:"eu_tax_management,omitempty"`
	LogoURL                   string                         `json:"logo_url,omitempty"`
	InvoiceFooter             string                         `json:"invoice_footer,omitempty"`
	InvoiceGracePeriod        int                            `json:"invoice_grace_period,omitempty"`
	DocumentLocale            string                         `json:"document_locale,omitempty"`
	Taxes                     []Tax                          `json:"taxes,omitempty"`
}

type BillingEntityCreateInput struct {
	Name                      string                            `json:"name,omitempty"`
	Code                      string                            `json:"code,omitempty"`
	Email                     string                            `json:"email,omitempty"`
	AddressLine1              string                            `json:"address_line1,omitempty"`
	AddressLine2              string                            `json:"address_line2,omitempty"`
	City                      string                            `json:"city,omitempty"`
	Zipcode                   string                            `json:"zipcode,omitempty"`
	State                     string                            `json:"state,omitempty"`
	Country                   string                            `json:"country,omitempty"`
	DefaultCurrency           Currency                          `json:"default_currency,omitempty"`
	LegalName                 string                            `json:"legal_name,omitempty"`
	LegalNumber               string                            `json:"legal_number,omitempty"`
	DocumentNumbering         BillingEntityDocumentNumbering    `json:"document_numbering,omitempty"`
	DocumentNumberPrefix      string                            `json:"document_number_prefix,omitempty"`
	NetPaymentTerm            int                               `json:"net_payment_term,omitempty"`
	Timezone                  string                            `json:"timezone,omitempty"`
	EmailSettings             []string                          `json:"email_settings,omitempty"`
	TaxIdentificationNumber   string                            `json:"tax_identification_number,omitempty"`
	FinalizeZeroAmountInvoice bool                              `json:"finalize_zero_amount_invoice,omitempty"`
	EuTaxManagement           bool                              `json:"eu_tax_management,omitempty"`
	BillingConfiguration      BillingEntityBillingConfiguration `json:"billing_configuration,omitempty"`
	LogoBase64                string                            `json:"logo,omitempty"`
}

type BillingEntityUpdateInput struct {
	Name                      string                            `json:"name,omitempty"`
	Email                     string                            `json:"email,omitempty"`
	AddressLine1              string                            `json:"address_line1,omitempty"`
	AddressLine2              string                            `json:"address_line2,omitempty"`
	City                      string                            `json:"city,omitempty"`
	Zipcode                   string                            `json:"zipcode,omitempty"`
	State                     string                            `json:"state,omitempty"`
	Country                   string                            `json:"country,omitempty"`
	DefaultCurrency           Currency                          `json:"default_currency,omitempty"`
	LegalName                 string                            `json:"legal_name,omitempty"`
	LegalNumber               string                            `json:"legal_number,omitempty"`
	DocumentNumbering         BillingEntityDocumentNumbering    `json:"document_numbering,omitempty"`
	DocumentNumberPrefix      string                            `json:"document_number_prefix,omitempty"`
	NetPaymentTerm            int                               `json:"net_payment_term,omitempty"`
	Timezone                  string                            `json:"timezone,omitempty"`
	EmailSettings             []string                          `json:"email_settings,omitempty"`
	TaxIdentificationNumber   string                            `json:"tax_identification_number,omitempty"`
	FinalizeZeroAmountInvoice bool                              `json:"finalize_zero_amount_invoice,omitempty"`
	EuTaxManagement           bool                              `json:"eu_tax_management,omitempty"`
	BillingConfiguration      BillingEntityBillingConfiguration `json:"billing_configuration,omitempty"`
	LogoBase64                string                            `json:"logo,omitempty"`
	TaxCodes                  []string                          `json:"tax_codes,omitempty"`
}

func (c *Client) BillingEntity() *BillingEntityRequest {
	return &BillingEntityRequest{
		client: c,
	}
}

func (ber *BillingEntityRequest) Create(ctx context.Context, billingEntityInput *BillingEntityCreateInput) (*BillingEntity, *Error) {
	billingEntityParams := &BillingEntityCreateParams{
		BillingEntity: billingEntityInput,
	}

	clientRequest := &ClientRequest{
		Path:   "billing_entities",
		Result: &BillingEntityResult{},
		Body:   billingEntityParams,
	}

	result, err := ber.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billingEntityResult, ok := result.(*BillingEntityResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billingEntityResult.BillingEntity, nil
}

func (ber *BillingEntityRequest) Get(ctx context.Context, billingEntityCode string) (*BillingEntity, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billing_entities", billingEntityCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillingEntityResult{},
	}

	result, err := ber.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billingEntityResult, ok := result.(*BillingEntityResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billingEntityResult.BillingEntity, nil
}

func (ber *BillingEntityRequest) GetList(ctx context.Context) (*BillingEntityResult, *Error) {
	clientRequest := &ClientRequest{
		Path:   "billing_entities",
		Result: &BillingEntityResult{},
	}

	result, err := ber.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billingEntityResult, ok := result.(*BillingEntityResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billingEntityResult, nil
}

func (ber *BillingEntityRequest) Update(ctx context.Context, billingEntityCode string, billingEntityInput *BillingEntityUpdateInput) (*BillingEntity, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billing_entities", billingEntityCode)

	billingEntityParams := &BillingEntityUpdateParams{
		BillingEntity: billingEntityInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillingEntityResult{},
		Body:   billingEntityParams,
	}

	result, err := ber.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billingEntityResult, ok := result.(*BillingEntityResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billingEntityResult.BillingEntity, nil
}
