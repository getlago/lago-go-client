package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CustomerPaymentProvider string

const (
	PaymentProviderAdyen      CustomerPaymentProvider = "adyen"
	PaymentProviderStripe     CustomerPaymentProvider = "stripe"
	PaymentProviderGocardless CustomerPaymentProvider = "gocardless"
)

type FinalizeZeroAmountInvoice string

const (
	FinalizeInvoice FinalizeZeroAmountInvoice = "finalize"
	SkipInvoice     FinalizeZeroAmountInvoice = "skip"
	InheritInvoice  FinalizeZeroAmountInvoice = "inherit"
)

type IntegrationType string

const (
	IntegrationNetsuite IntegrationType = "netsuite"
	IntegrationAnrok    IntegrationType = "anrok"
	IntegrationXero     IntegrationType = "xero"
)

type CustomerType string

const (
	CompanyCustomerType    CustomerType = "company"
	IndividualCustomerType CustomerType = "individual"
)

type ProviderPaymentMethodType string

const (
	CardPaymentMethodType          ProviderPaymentMethodType = "card"
	SepaDebitPaymentMethodType     ProviderPaymentMethodType = "sepa_debit"
	USBankAccountPaymentMethodType ProviderPaymentMethodType = "us_bank_account"
	BacsDebitPaymentMethodType     ProviderPaymentMethodType = "bacs_debit"
	LinkPaymentMethodType          ProviderPaymentMethodType = "link"
	BoletoPaymentMethodType        ProviderPaymentMethodType = "boleto"
)

type CustomerParams struct {
	Customer *CustomerInput `json:"customer"`
}

type CustomerResult struct {
	Customer  *Customer  `json:"customer"`
	Customers []Customer `json:"customers,omitempty"`
	Meta      Metadata   `json:"meta,omitempty"`
}

type CustomerUsageResult struct {
	CustomerUsage *CustomerUsage `json:"customer_usage"`
}

type CustomerPastUsageResult struct {
	UsagePeriods []CustomerUsage `json:"usage_periods"`
	Meta         Metadata        `json:"meta"`
}

type CustomerPortalUrlResult struct {
	CustomerPortalUrl *CustomerPortalUrl `json:"customer"`
}

type CustomerCheckoutUrlResult struct {
	CustomerCheckoutUrl *CustomerCheckoutUrl `json:"customer"`
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
	ExternalID                string                            `json:"external_id,omitempty"`
	BillingEntityCode         string                            `json:"billing_entity_code,omitempty"`
	Name                      string                            `json:"name,omitempty"`
	Firstname                 string                            `json:"firstname,omitempty"`
	Lastname                  string                            `json:"lastname,omitempty"`
	CustomerType              CustomerType                      `json:"customer_type,omitempty"`
	Email                     string                            `json:"email,omitempty"`
	AddressLine1              string                            `json:"address_line1,omitempty"`
	AddressLine2              string                            `json:"address_line2,omitempty"`
	City                      string                            `json:"city,omitempty"`
	Zipcode                   string                            `json:"zipcode,omitempty"`
	State                     string                            `json:"state,omitempty"`
	Country                   string                            `json:"country,omitempty"`
	LegalName                 string                            `json:"legal_name,omitempty"`
	LegalNumber               string                            `json:"legal_number,omitempty"`
	NetPaymentTerm            int                               `json:"net_payment_term,omitempty"`
	TaxIdentificationNumber   string                            `json:"tax_identification_number,omitempty"`
	Phone                     string                            `json:"phone,omitempty"`
	URL                       string                            `json:"url,omitempty"`
	Currency                  Currency                          `json:"currency,omitempty"`
	Timezone                  string                            `json:"timezone,omitempty"`
	Metadata                  []CustomerMetadataInput           `json:"metadata,omitempty"`
	BillingConfiguration      CustomerBillingConfigurationInput `json:"billing_configuration,omitempty"`
	ShippingAddress           Address                           `json:"shipping_address,omitempty"`
	IntegrationCustomers      []IntegrationCustomer             `json:"integration_customers,omitempty"`
	TaxCodes                  []string                          `json:"tax_codes,omitempty"`
	FinalizeZeroAmountInvoice FinalizeZeroAmountInvoice         `json:"finalize_zero_amount_invoice,omitempty"`
	SkipInvoiceCustomSections bool                              `json:"skip_invoice_custom_sections,omitempty"`
	InvoiceCustomSectionCodes []string                          `json:"invoice_custom_section_codes,omitempty"`
}

type CustomerListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type CustomerBillingConfigurationInput struct {
	InvoiceGracePeriod     int                         `json:"invoice_grace_period,omitempty"`
	PaymentProvider        CustomerPaymentProvider     `json:"payment_provider,omitempty"`
	PaymentProviderCode    string                      `json:"payment_provider_code,omitempty"`
	ProviderCustomerID     string                      `json:"provider_customer_id,omitempty"`
	Sync                   bool                        `json:"sync,omitempty"`
	SyncWithProvider       bool                        `json:"sync_with_provider,omitempty"`
	DocumentLocale         string                      `json:"document_locale,omitempty"`
	ProviderPaymentMethods []ProviderPaymentMethodType `json:"provider_payment_methods,omitempty"`
}

type CustomerBillingConfiguration struct {
	InvoiceGracePeriod     int                         `json:"invoice_grace_period,omitempty"`
	PaymentProvider        CustomerPaymentProvider     `json:"payment_provider,omitempty"`
	PaymentProviderCode    string                      `json:"payment_provider_code,omitempty"`
	ProviderCustomerID     string                      `json:"provider_customer_id,omitempty"`
	SyncWithProvider       bool                        `json:"sync_with_provider,omitempty"`
	DocumentLocale         string                      `json:"document_locale,omitempty"`
	ProviderPaymentMethods []ProviderPaymentMethodType `json:"provider_payment_methods,omitempty"`
}

type Address struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	Zipcode      string `json:"zipcode,omitempty"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country,omitempty"`
}

type IntegrationCustomer struct {
	LagoID             uuid.UUID       `json:"id,omitempty"`
	ExternalCustomerId string          `json:"external_customer_id,omitempty"`
	IntegrationType    IntegrationType `json:"integration_type,omitempty"`
	IntegrationCode    string          `json:"integration_code,omitempty"`
	SubsidiaryId       string          `json:"subsidiary_id,omitempty"`
	SyncWithProvider   bool            `json:"sync_with_provider,omitempty"`
}

type IntegrationCustomersResponse struct {
	LagoID             uuid.UUID       `json:"lago_id,omitempty"`
	ExternalCustomerId string          `json:"external_customer_id,omitempty"`
	IntegrationType    IntegrationType `json:"type,omitempty"`
	IntegrationCode    string          `json:"integration_code,omitempty"`
	SubsidiaryId       string          `json:"subsidiary_id,omitempty"`
	SyncWithProvider   bool            `json:"sync_with_provider,omitempty"`
}

type CustomerChargeUsage struct {
	Units          string   `json:"units,omitempty"`
	EventsCount    int      `json:"events_count"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	Charge             *Charge                      `json:"charge,omitempty"`
	BillableMetric     *BillableMetric              `json:"billable_metric,omitempty"`
	PricingUnitDetails *CustomerPricingUnitDetails  `json:"pricing_unit_details,omitempty"`
	Filters            []CustomerChargeFilterUsage  `json:"filters,omitempty"`
	GroupedUsage       []CustomerChargeGroupedUsage `json:"grouped_usage,omitempty"`
}

type CustomerPricingUnitDetails struct {
	AmountCents    int     `json:"amount_cents,omitempty"`
	ShortName      string  `json:"short_name,omitempty"`
	ConversionRate float64 `json:"conversion_rate,omitempty"`
}

type CustomerChargeFilterUsage struct {
	InvoiceDisplayName string                      `json:"invoice_display_name,omitempty"`
	Values             map[string]interface{}      `json:"value,omitempty"`
	AmountCents        int                         `json:"amount_cents,omitempty"`
	EventsCount        int                         `json:"events_count,omitempty"`
	Units              string                      `json:"units,omitempty"`
	PricingUnitDetails *CustomerPricingUnitDetails `json:"pricing_unit_details,omitempty"`
}

type CustomerChargeGroupedUsage struct {
	AmountCents        int                         `json:"amount_cents,omitempty"`
	EventsCount        int                         `json:"events_count,omitempty"`
	Units              string                      `json:"units,omitempty"`
	GroupedBy          map[string]interface{}      `json:"grouped_by,omitempty"`
	Filters            []CustomerChargeFilterUsage `json:"filters,omitempty"`
	PricingUnitDetails *CustomerPricingUnitDetails `json:"pricing_unit_details,omitempty"`
}

type CustomerUsage struct {
	FromDatetime     time.Time `json:"from_datetime,omitempty"`
	ToDatetime       time.Time `json:"to_datetime,omitempty"`
	IssuingDate      string    `json:"issuing_date,omitempty"`
	LagoInvoiceID    string    `json:"lago_invoice_id,omitempty"`
	Currency         Currency  `json:"currency,omitempty"`
	AmountCents      int       `json:"amount_cents,omitempty"`
	TotalAmountCents int       `json:"total_amount_cents,omitempty"`
	TaxesAmountCents int       `json:"taxes_amount_cents,omitempty"`

	ChargesUsage []CustomerChargeUsage `json:"charges_usage,omitempty"`
}

type CustomerPortalUrl struct {
	PortalUrl string `json:"portal_url,omitempty"`
}

type CustomerCheckoutUrl struct {
	CheckoutUrl string `json:"checkout_url,omitempty"`
}

type CustomerUsageInput struct {
	ExternalSubscriptionID string `json:"external_subscription_id,omitempty"`
	ApplyTaxes             bool   `json:"apply_taxes,omitempty"`
}

type CustomerPastUsageInput struct {
	ExternalSubscriptionID string `json:"external_subscription_id"`
	BillableMetricCode     string `json:"billable_metric_code,omitempty"`
	PeriodsCount           int    `json:"periods_count,omitempty,string"`
}

type Customer struct {
	LagoID            uuid.UUID `json:"lago_id,omitempty"`
	SequentialID      int       `json:"sequential_id,omitempty"`
	ExternalID        string    `json:"external_id,omitempty"`
	Slug              string    `json:"slug,omitempty"`
	BillingEntityCode string    `json:"billing_entity_code,omitempty"`

	Name                      string                         `json:"name,omitempty"`
	Firstname                 string                         `json:"firstname,omitempty"`
	Lastname                  string                         `json:"lastname,omitempty"`
	CustomerType              string                         `json:"customer_type,omitempty"`
	Email                     string                         `json:"email,omitempty"`
	AddressLine1              string                         `json:"address_line1,omitempty"`
	AddressLine2              string                         `json:"address_line2,omitempty"`
	City                      string                         `json:"city,omitempty"`
	State                     string                         `json:"state,omitempty"`
	Zipcode                   string                         `json:"zipcode,omitempty"`
	Country                   string                         `json:"country,omitempty"`
	LegalName                 string                         `json:"legal_name,omitempty"`
	LegalNumber               string                         `json:"legal_number,omitempty"`
	NetPaymentTerm            int                            `json:"net_payment_term,omitempty"`
	TaxIdentificationNumber   string                         `json:"tax_identification_number,omitempty"`
	LogoURL                   string                         `json:"logo_url,omitempty"`
	Phone                     string                         `json:"phone,omitempty"`
	URL                       string                         `json:"url,omitempty"`
	FinalizeZeroAmountInvoice FinalizeZeroAmountInvoice      `json:"finalize_zero_amount_invoice,omitempty"`
	BillingConfiguration      CustomerBillingConfiguration   `json:"billing_configuration,omitempty"`
	ShippingAddress           Address                        `json:"shipping_address,omitempty"`
	IntegrationCustomers      []IntegrationCustomersResponse `json:"integration_customers,omitempty"`
	Metadata                  []MetadataResponse             `json:"metadata,omitempty"`
	Currency                  Currency                       `json:"currency,omitempty"`
	Timezone                  string                         `json:"timezone,omitempty"`
	ApplicableTimezone        string                         `json:"applicable_timezone,omitempty"`
	SkipInvoiceCustomSections bool                           `json:"skip_invoice_custom_sections,omitempty"`

	Taxes                           []Tax                  `json:"taxes,omitempty"`
	ApplicableInvoiceCustomSections []InvoiceCustomSection `json:"applicable_invoice_custom_sections,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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

	queryParams := map[string]string{
		"external_subscription_id": customerUsageInput.ExternalSubscriptionID,
		"apply_taxes":              strconv.FormatBool(customerUsageInput.ApplyTaxes),
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &CustomerUsageResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	currentUsageResult, ok := result.(*CustomerUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return currentUsageResult.CustomerUsage, nil
}

func (cr *CustomerRequest) PastUsage(ctx context.Context, externalCustomerID string, customerPastUsageInput *CustomerPastUsageInput) (*CustomerPastUsageResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "past_usage")

	jsonQueryParams, err := json.Marshal(customerPastUsageInput)
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
		Result:      &CustomerPastUsageResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	pastUsageResult, ok := result.(*CustomerPastUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return pastUsageResult, nil
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

func (cr *CustomerRequest) CheckoutUrl(ctx context.Context, externalCustomerID string) (*CustomerCheckoutUrl, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "checkout_url")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerCheckoutUrlResult{},
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	checkoutUrlResult, ok := result.(*CustomerCheckoutUrlResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return checkoutUrlResult.CustomerCheckoutUrl, nil
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
	if clientErr != nil {
		return nil, clientErr
	}

	customerResult, ok := result.(*CustomerResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return customerResult, nil
}
