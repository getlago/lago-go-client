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
	PaymentProviderStripe CustomerPaymentProvider = "stripe"
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

type CustomerInput struct {
	ExternalID           string                            `json:"external_id,omitempty"`
	Name                 string                            `json:"name,omitempty"`
	Email                string                            `json:"email,omitempty"`
	AddressLine1         string                            `json:"address_line_1,omitempty"`
	AddressLine2         string                            `json:"address_line_2,omitempty"`
	City                 string                            `json:"city,omitempty"`
	Zipcode              string                            `json:"zipcode,omitempty"`
	State                string                            `json:"state,omitempty"`
	Country              string                            `json:"country,omitempty"`
	LegalName            string                            `json:"legal_name,omitempty"`
	LegalNumber          string                            `json:"legal_number,omitempty"`
	Phone                string                            `json:"phone,omitempty"`
	URL                  string                            `json:"url,omitempty"`
	Currency             Currency                          `json:"currency,omitempty"`
	BillingConfiguration CustomerBillingConfigurationInput `json:"billing_configuration,omitempty"`
	VatRate              float32                           `json:"vat_rate,omitempty"`
}

type CustomerListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type CustomerBillingConfigurationInput struct {
	PaymentProvider    CustomerPaymentProvider `json:"payment_provider,omitempty"`
	ProviderCustomerID string                  `json:"provider_customer_id,omitempty"`
	Sync               bool                    `json:"sync,omitempty"`
}

type CustomerBillingConfiguration struct {
	PaymentProvider    CustomerPaymentProvider `json:"payment_provider,omitempty"`
	ProviderCustomerID string                  `json:"provider_customer_id,omitempty"`
}

type CustomerChargeUsage struct {
	Units          string   `json:"units,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	Charge         *Charge         `json:"charge,omitempty"`
	BillableMetric *BillableMetric `json:"billable_metric,omitempty"`
}

type CustomerUsage struct {
	FromDate    string `json:"from_date,omitempty"`
	ToDate      string `json:"to_date,omitempty"`
	IssuingDate string `json:"issuing_date,omitempty"`

	AmountCents         int      `json:"amount_cents,omitempty"`
	AmountCurrency      Currency `json:"amount_currency,omitempty"`
	TotalAmountCents    int      `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency Currency `json:"total_amount_currency,omitempty"`
	VatAmountCents      int      `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency   Currency `json:"vat_amount_currency,omitempty"`

	ChargesUsage []CustomerChargeUsage `json:"charges_usage,omitempty"`
}

type Customer struct {
	LagoID       uuid.UUID `json:"lago_id,omitempty"`
	SequentialID int       `json:"sequential_id,omitempty"`
	ExternalID   string    `json:"external_id,omitempty"`
	Slug         string    `json:"slug,omitempty"`

	Name                 string                       `json:"name,omitempty"`
	Email                string                       `json:"email,omitempty"`
	AddressLine1         string                       `json:"address_line1,omitempty"`
	AddressLine2         string                       `json:"address_line2,omitempty"`
	City                 string                       `json:"city,omitempty"`
	State                string                       `json:"state,omitempty"`
	Zipcode              string                       `json:"zipcode,omitempty"`
	Country              string                       `json:"country,omitempty"`
	LegalName            string                       `json:"legal_name,omitempty"`
	LegalNumber          string                       `json:"legal_number,omitempty"`
	LogoURL              string                       `json:"logo_url,omitempty"`
	Phone                string                       `json:"phone,omitempty"`
	URL                  string                       `json:"url,omitempty"`
	BillingConfiguration CustomerBillingConfiguration `json:"billing_configuration,omitempty"`
	VatRate              float32                      `json:"vat_rate,omitempty"`
	Currency             Currency                     `json:"currency,omitempty"`

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

	customerResult := result.(*CustomerResult)

	return customerResult.Customer, nil
}

// NOTE: Update endpoint does not exists, actually we use the create endpoint with the
// same externalID to update a customer
func (cr *CustomerRequest) Update(ctx context.Context, customerInput *CustomerInput) (*Customer, *Error) {
	return cr.Create(ctx, customerInput)
}

func (cr *CustomerRequest) CurrentUsage(ctx context.Context, externalCustomerID string) (*CustomerUsage, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "current_usage")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerUsageResult{},
	}

	result, err := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	currentUsageResult := result.(*CustomerUsageResult)

	return currentUsageResult.CustomerUsage, nil
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

	customerResult := result.(*CustomerResult)

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

	customerResult := result.(*CustomerResult)

	return customerResult, nil
}
