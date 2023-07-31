package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CustomerPaymentProvider string

const (
	PaymentProviderStripe     CustomerPaymentProvider = "stripe"
	PaymentProviderGocardless CustomerPaymentProvider = "gocardless"
)

type CustomerParams struct {
	Customer *CustomerInput `json:"customer"`
}

type CustomerResult struct {
	Customer  *Customer  `json:"customer"`
	Customers []Customer `json:"customers,omitempty"`
	Meta      Metadata   `json:"medata,omitempty"`
}

type CustomerUsageResult struct {
	CustomerUsage *CustomerUsage `json:"customer_usage"`
}

type CustomerPortalUrlResult struct {
	CustomerPortalUrl *CustomerPortalUrl `json:"customer"`
}

type CustomerMetadataInput struct {
	LagoID           *uuid.UUID `json:"id,omitempty"`
	Key              string     `json:"key,omitempty"`
	Value            string     `json:"value,omitempty"`
	DisplayInInvoice bool       `json:"display_in_invoice,omitempty"`
}

type MetadataResponse struct {
	LagoID           uuid.UUID `json:"lago_id,omitempty"`
	Key              string    `json:"key,omitempty"`
	Value            string    `json:"value,omitempty"`
	DisplayInInvoice bool      `json:"display_in_invoice,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type CustomerInput struct {
	ExternalID              string                            `json:"external_id,omitempty"`
	Name                    string                            `json:"name,omitempty"`
	Email                   string                            `json:"email,omitempty"`
	AddressLine1            string                            `json:"address_line1,omitempty"`
	AddressLine2            string                            `json:"address_line2,omitempty"`
	City                    string                            `json:"city,omitempty"`
	Zipcode                 string                            `json:"zipcode,omitempty"`
	State                   string                            `json:"state,omitempty"`
	Country                 string                            `json:"country,omitempty"`
	LegalName               string                            `json:"legal_name,omitempty"`
	LegalNumber             string                            `json:"legal_number,omitempty"`
	NetPaymentTerm          int                               `json:"net_payment_term,omitempty"`
	TaxIdentificationNumber string                            `json:"tax_identification_number,omitempty"`
	Phone                   string                            `json:"phone,omitempty"`
	URL                     string                            `json:"url,omitempty"`
	Currency                Currency                          `json:"currency,omitempty"`
	Timezone                string                            `json:"timezone,omitempty"`
	Metadata                []CustomerMetadataInput           `json:"metadata,omitempty"`
	BillingConfiguration    CustomerBillingConfigurationInput `json:"billing_configuration,omitempty"`
	TaxCodes                []string                          `json:"tax_codes,omitempty"`
}

type CustomerListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type CustomerBillingConfigurationInput struct {
	InvoiceGracePeriod int                     `json:"invoice_grace_period,omitempty"`
	PaymentProvider    CustomerPaymentProvider `json:"payment_provider,omitempty"`
	ProviderCustomerID string                  `json:"provider_customer_id,omitempty"`
	Sync               bool                    `json:"sync,omitempty"`
	SyncWithProvider   bool                    `json:"sync_with_provider,omitempty"`
	DocumentLocale     string                  `json:"document_locale,omitempty"`
}

type CustomerBillingConfiguration struct {
	InvoiceGracePeriod int                     `json:"invoice_grace_period,omitempty"`
	PaymentProvider    CustomerPaymentProvider `json:"payment_provider,omitempty"`
	ProviderCustomerID string                  `json:"provider_customer_id,omitempty"`
	SyncWithProvider   bool                    `json:"sync_with_provider,omitempty"`
	DocumentLocale     string                  `json:"document_locale,omitempty"`
}

type CustomerChargeUsage struct {
	Units          string   `json:"units,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	Charge         *Charge         `json:"charge,omitempty"`
	BillableMetric *BillableMetric `json:"billable_metric,omitempty"`
}

type CustomerUsage struct {
	FromDatetime     time.Time `json:"from_datetime,omitempty"`
	ToDatetime       time.Time `json:"to_datetime,omitempty"`
	IssuingDate      string    `json:"issuing_date,omitempty"`
	Currency         Currency  `json:"currency,omitempty"`
	AmountCents      int       `json:"amount_cents,omitempty"`
	TotalAmountCents int       `json:"total_amount_cents,omitempty"`
	TaxesAmountCents int       `json:"taxes_amount_cents,omitempty"`

	ChargesUsage []CustomerChargeUsage `json:"charges_usage,omitempty"`
}

type CustomerPortalUrl struct {
	PortalUrl string `json:"portal_url,omitempty"`
}

type CustomerUsageInput struct {
	ExternalSubscriptionID string `json:"external_subscription_id,omitempty"`
}

type Customer struct {
	LagoID       uuid.UUID `json:"lago_id,omitempty"`
	SequentialID int       `json:"sequential_id,omitempty"`
	ExternalID   string    `json:"external_id,omitempty"`
	Slug         string    `json:"slug,omitempty"`

	Name                    string                       `json:"name,omitempty"`
	Email                   string                       `json:"email,omitempty"`
	AddressLine1            string                       `json:"address_line1,omitempty"`
	AddressLine2            string                       `json:"address_line2,omitempty"`
	City                    string                       `json:"city,omitempty"`
	State                   string                       `json:"state,omitempty"`
	Zipcode                 string                       `json:"zipcode,omitempty"`
	Country                 string                       `json:"country,omitempty"`
	LegalName               string                       `json:"legal_name,omitempty"`
	LegalNumber             string                       `json:"legal_number,omitempty"`
	NetPaymentTerm          int                          `json:"net_payment_term,omitempty"`
	TaxIdentificationNumber string                       `json:"tax_identification_number,omitempty"`
	LogoURL                 string                       `json:"logo_url,omitempty"`
	Phone                   string                       `json:"phone,omitempty"`
	URL                     string                       `json:"url,omitempty"`
	BillingConfiguration    CustomerBillingConfiguration `json:"billing_configuration,omitempty"`
	Metadata                []MetadataResponse           `json:"metadata,omitempty"`
	Currency                Currency                     `json:"currency,omitempty"`
	Timezone                string                       `json:"timezone,omitempty"`
	ApplicableTimezone      string                       `json:"applicable_timezone,omitempty"`

	Taxes []Tax `json:"taxes,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CustomerRequest struct {
	client *Client
}

func (c *Client) Customer() *CustomerRequest {
	return &CustomerRequest{
		client: c,
	}
}

func (cr *CustomerRequest) Create(ctx context.Context, customerInput *CustomerInput) (*Customer, *Error) {
	customerParams := &CustomerParams{
		Customer: customerInput,
	}

	clientRequest := &ClientRequest{
		Path:   "customers",
		Result: &CustomerResult{},
		Body:   customerParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	customerResult, ok := result.(*CustomerResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return customerResult.Customer, nil
}

// NOTE: Update endpoint does not exists, actually we use the create endpoint with the
// same externalID to update a customer
func (cr *CustomerRequest) Update(ctx context.Context, customerInput *CustomerInput) (*Customer, *Error) {
	return cr.Create(ctx, customerInput)
}

func (cr *CustomerRequest) CurrentUsage(ctx context.Context, externalCustomerID string, customerUsageInput *CustomerUsageInput) (*CustomerUsage, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "current_usage")

	jsonQueryParams, err := json.Marshal(customerUsageInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &CustomerUsageResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, clientErr
	}

	currentUsageResult, ok := result.(*CustomerUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return currentUsageResult.CustomerUsage, nil
}

func (cr *CustomerRequest) PortalUrl(ctx context.Context, externalCustomerID string) (*CustomerPortalUrl, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "portal_url")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerPortalUrlResult{},
	}

	result, err := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	portalUrlResult, ok := result.(*CustomerPortalUrlResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return portalUrlResult.CustomerPortalUrl, nil
}

func (cr *CustomerRequest) Delete(ctx context.Context, externalCustomerID string) (*Customer, *Error) {
	subPath := fmt.Sprintf("%s/%s", "customers", externalCustomerID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerResult{},
	}

	result, err := cr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	customerResult, ok := result.(*CustomerResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return customerResult.Customer, nil
}

func (cr *CustomerRequest) Get(ctx context.Context, externalCustomerID string) (*Customer, *Error) {
	subPath := fmt.Sprintf("%s/%s", "customers", externalCustomerID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerResult{},
	}

	result, err := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	customerResult, ok := result.(*CustomerResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return customerResult.Customer, nil
}

func (cr *CustomerRequest) GetList(ctx context.Context, customerListInput *CustomerListInput) (*CustomerResult, *Error) {
	jsonQueryParams, err := json.Marshal(customerListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "customers",
		QueryParams: queryParams,
		Result:      &CustomerResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, clientErr
	}

	customerResult, ok := result.(*CustomerResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return customerResult, nil
}
