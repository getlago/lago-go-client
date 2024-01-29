package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type InvoiceType string
type InvoiceStatus string
type InvoicePaymentStatus string
type InvoiceCreditItemType string

const (
	SubscriptionInvoiceType InvoiceType = "subscription"
	AddOnInvoiceType        InvoiceType = "add_on"
	CreditInvoiceType       InvoiceType = "credit"
	OneOffInvoiceType       InvoiceType = "one_off"
)

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusFinalized InvoiceStatus = "finalized"
)

const (
	InvoicePaymentStatusPending   InvoicePaymentStatus = "pending"
	InvoicePaymentStatusSucceeded InvoicePaymentStatus = "succeeded"
	InvoicePaymentStatusFailed    InvoicePaymentStatus = "failed"
)

const (
	InvoiceCreditItemCredit InvoiceCreditItemType = "coupon"
)

type InvoiceRequest struct {
	client *Client
}

type InvoiceResult struct {
	Invoice  *Invoice  `json:"invoice,omitempty"`
	Invoices []Invoice `json:"invoices,omitempty"`
	Meta     Metadata  `json:"meta,omitempty"`
}

type InvoicePaymentUrlResult struct {
	InvoicePaymentUrl *InvoicePaymentUrl `json:"invoice_payment_url"`
}

type InvoiceParams struct {
	Invoice *InvoiceInput `json:"invoice"`
}

type InvoiceOneOffParams struct {
	Invoice *InvoiceOneOffInput `json:"invoice"`
}

type InvoiceMetadataInput struct {
	LagoID *uuid.UUID `json:"id,omitempty"`
	Key    string     `json:"key,omitempty"`
	Value  string     `json:"value,omitempty"`
}

type InvoiceFeesInput struct {
	AddOnCode          string   `json:"add_on_code,omitempty"`
	InvoiceDisplayName string   `json:"invoice_display_name,omitempty"`
	UnitAmountCents    int      `json:"unit_amount_cents,omitempty"`
	Description        string   `json:"description,omitempty"`
	Units              float32  `json:"units,omitempty"`
	TaxCodes           []string `json:"tax_codes,omitempty"`
}

type InvoiceMetadataResponse struct {
	LagoID    uuid.UUID `json:"lago_id,omitempty"`
	Key       string    `json:"key,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type InvoiceInput struct {
	LagoID        uuid.UUID              `json:"lago_id,omitempty"`
	PaymentStatus InvoicePaymentStatus   `json:"payment_status,omitempty"`
	Metadata      []InvoiceMetadataInput `json:"metadata,omitempty"`
}

type InvoiceOneOffInput struct {
	ExternalCustomerId string             `json:"external_customer_id,omitempty"`
	Currency           string             `json:"currency,omitempty"`
	Fees               []InvoiceFeesInput `json:"fees,omitempty"`
}

type InvoiceListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	IssuingDateFrom string `json:"issuing_date_from,omitempty"`
	IssuingDateTo   string `json:"issuing_date_to,omitempty"`

	ExternalCustomerID string               `json:"external_customer_id,omitempty"`
	Status             InvoiceStatus        `json:"status,omitempty"`
	PaymentStatus      InvoicePaymentStatus `json:"payment_status,omitempty"`
}

type InvoiceCreditItem struct {
	LagoID uuid.UUID             `json:"lago_id,omitempty"`
	Type   InvoiceCreditItemType `json:"type,omitempty"`
	Code   string                `json:"code,omitempty"`
	Name   string                `json:"name,omitempty"`
}

type InvoiceSummary struct {
	LagoID        uuid.UUID            `json:"lago_id,omitempty"`
	PaymentStatus InvoicePaymentStatus `json:"payment_status,omitempty"`
}

type InvoiceCredit struct {
	Item InvoiceCreditItem `json:"item,omitempty"`

	Invoice InvoiceSummary `json:"invoice,omitempty"`

	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	AmountCents    int       `json:"amount_cents,omitempty"`
	AmountCurrency Currency  `json:"amount_currency,omitempty"`
	BeforeVAT      bool      `json:"before_vat,omitempty"`
}

type InvoiceAppliedTax struct {
	LagoId          uuid.UUID `json:"lago_id,omitempty"`
	LagoInvoiceId   uuid.UUID `json:"lago_invoice_id,omitempty"`
	LagoTaxId       uuid.UUID `json:"lago_tax_id,omitempty"`
	TaxName         string    `json:"tax_name,omitempty"`
	TaxCode         string    `json:"tax_code,omitempty"`
	TaxRate         float32   `json:"tax_rate,omitempty"`
	TaxDescription  string    `json:"tax_description,omitempty"`
	AmountCents     int       `json:"amount_cents,omitempty"`
	AmountCurrency  Currency  `json:"amount_currency,omitempty"`
	FeesAmountCents int       `json:"fees_amount_cents,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

type Invoice struct {
	LagoID       uuid.UUID `json:"lago_id,omitempty"`
	SequentialID int       `json:"sequential_id,omitempty"`
	Number       string    `json:"number,omitempty"`

	IssuingDate    string `json:"issuing_date,omitempty"`
	PaymentDueDate string `json:"payment_due_date,omitempty"`

	InvoiceType   InvoiceType          `json:"invoice_type,omitempty"`
	Status        InvoiceStatus        `json:"status,omitempty"`
	PaymentStatus InvoicePaymentStatus `json:"payment_status,omitempty"`

	Currency Currency `json:"currency,omitempty"`

	FeesAmountCents                   int `json:"fees_amount_cents,omitempty"`
	TaxesAmountCents                  int `json:"taxes_amount_cents,omitempty"`
	CouponsAmountCents                int `json:"coupons_amount_cents,omitempty"`
	CreditNotesAmountCents            int `json:"credit_notes_amount_cents,omitempty"`
	SubTotalExcludingTaxesAmountCents int `json:"sub_total_excluding_taxes_amount_cents,omitempty"`
	SubTotalIncludingTaxesAmountCents int `json:"sub_total_including_taxes_amount_cents,omitempty"`
	TotalAmountCents                  int `json:"total_amount_cents,omitempty"`
	PrepaidCreditAmountCents          int `json:"prepaid_credit_amount_cents,omitempty"`
	NetPaymentTerm                    int `json:"net_payment_term,omitempty"`

	FileURL       string                    `json:"file_url,omitempty"`
	Metadata      []InvoiceMetadataResponse `json:"metadata,omitempty"`
	VersionNumber int                       `json:"version_number,omitempty"`

	Customer      *Customer      `json:"customer,omitempty"`
	Subscriptions []Subscription `json:"subscriptions,omitempty"`

	Fees         []Fee               `json:"fees,omitempty"`
	Credits      []InvoiceCredit     `json:"credits,omitempty"`
	AppliedTaxes []InvoiceAppliedTax `json:"applied_taxes,omitempty"`

	// Deprecated: Will be removed in the future
	Legacy                         bool     `json:"legacy,omitempty"`
	AmountCurrency                 Currency `json:"amount_currency,omitempty"`
	AmountCents                    int      `json:"amount_cents,omitempty"`
	CreditAmountCents              int      `json:"credit_amount_cents,omitempty"`
	CreditAmountCurrency           Currency `json:"credit_amount_currency,omitempty"`
	TotalAmountCurrency            Currency `json:"total_amount_currency,omitempty"`
	VatAmountCents                 int      `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency              Currency `json:"vat_amount_currency,omitempty"`
	SubTotalVatExcludedAmountCents int      `json:"sub_total_vat_excluded_amount_cents,omitempty"`
	SubTotalVatIncludedAmountCents int      `json:"sub_total_vat_included_amount_cents,omitempty"`
}

type InvoicePaymentUrl struct {
	PaymentUrl string `json:"payment_url,omitempty"`
}

func (c *Client) Invoice() *InvoiceRequest {
	return &InvoiceRequest{
		client: c,
	}
}

func (ir *InvoiceRequest) Get(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s", "invoices", invoiceID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	invoiceResult, ok := result.(*InvoiceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return invoiceResult.Invoice, nil
}

func (ir *InvoiceRequest) GetList(ctx context.Context, invoiceListInput *InvoiceListInput) (*InvoiceResult, *Error) {
	jsonQueryParams, err := json.Marshal(invoiceListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "invoices",
		QueryParams: queryParams,
		Result:      &InvoiceResult{},
	}

	result, clientErr := ir.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	invoiceResult, ok := result.(*InvoiceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return invoiceResult, nil
}

func (ir *InvoiceRequest) Create(ctx context.Context, oneOffInput *InvoiceOneOffInput) (*Invoice, *Error) {
	invoiceOneOffParams := &InvoiceOneOffParams{
		Invoice: oneOffInput,
	}

	clientRequest := &ClientRequest{
		Path:   "invoices",
		Result: &InvoiceResult{},
		Body:   invoiceOneOffParams,
	}

	result, err := ir.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	invoiceResult, ok := result.(*InvoiceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return invoiceResult.Invoice, nil
}

func (ir *InvoiceRequest) Update(ctx context.Context, invoiceInput *InvoiceInput) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s", "invoices", invoiceInput.LagoID)
	invoiceParams := &InvoiceParams{
		Invoice: invoiceInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
		Body:   invoiceParams,
	}

	result, err := ir.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	invoiceResult, ok := result.(*InvoiceResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return invoiceResult.Invoice, nil
}

func (ir *InvoiceRequest) Download(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "download")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.PostWithoutBody(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result != nil {
		invoiceResult, ok := result.(*InvoiceResult)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return invoiceResult.Invoice, nil
	}

	return nil, nil
}

func (ir *InvoiceRequest) Refresh(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "refresh")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result != nil {
		invoiceResult, ok := result.(*InvoiceResult)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return invoiceResult.Invoice, nil
	}

	return nil, nil
}

func (ir *InvoiceRequest) Finalize(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "finalize")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result != nil {
		invoiceResult, ok := result.(*InvoiceResult)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return invoiceResult.Invoice, nil
	}

	return nil, nil
}

func (ir *InvoiceRequest) RetryPayment(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "retry_payment")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.PostWithoutBody(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result != nil {
		invoiceResult, ok := result.(*InvoiceResult)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return invoiceResult.Invoice, nil
	}

	return nil, nil
}

func (ir *InvoiceRequest) PaymentUrl(ctx context.Context, invoiceID string) (*InvoicePaymentUrl, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "payment_url")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoicePaymentUrlResult{},
	}

	result, err := ir.client.PostWithoutBody(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result != nil {
		paymentUrlResult, ok := result.(*InvoicePaymentUrlResult)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return paymentUrlResult.InvoicePaymentUrl, nil
	}

	return nil, nil
}
