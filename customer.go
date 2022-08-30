package lago

import (
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
	Customer *Customer `json:"customer"`
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
	BillingConfiguration CustomerBillingConfigurationInput `json:"billing_configuration,omitempty"`
	VatRate              float32                           `json:"vat_rate,omitempty"`
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

func (cr *CustomerRequest) Create(customerInput *CustomerInput) (*Customer, *Error) {
	customerParams := &CustomerParams{
		Customer: customerInput,
	}

	clientRequest := &ClientRequest{
		Path:   "customers",
		Result: &CustomerResult{},
		Body:   customerParams,
	}

	result, err := cr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	customerResult := result.(*CustomerResult)

	return customerResult.Customer, nil
}

// NOTE: Update endpoint does not exists, actually we use the create endpoint with the
// same externalID to update a customer
func (cr *CustomerRequest) Update(customerInput *CustomerInput) (*Customer, *Error) {
	return cr.Create(customerInput)
}

func (cr *CustomerRequest) CurrentUsage(externalCustomerID string) (*CustomerUsage, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "current_usage")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerUsageResult{},
	}

	result, err := cr.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	currentUsageResult := result.(*CustomerUsageResult)

	return currentUsageResult.CustomerUsage, nil
}
