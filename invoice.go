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

type InvoiceParams struct {
	Invoice *InvoiceInput `json:"invoice"`
}

type InvoiceMetadataInput struct {
	LagoID *uuid.UUID `json:"id,omitempty"`
	Key    string     `json:"key,omitempty"`
	Value  string     `json:"value,omitempty"`
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

type InvoiceListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	IssuingDateFrom string `json:"issuing_date_from,omitempty"`
	IssuingDateTo   string `json:"issuing_date_to,omitempty"`

	ExternalCustomerId string               `json:"external_customer_id,omitempty"`
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
}

type Invoice struct {
	LagoID       uuid.UUID `json:"lago_id,omitempty"`
	SequentialID int       `json:"sequential_id,omitempty"`
	Number       string    `json:"number,omitempty"`

	IssuingDate string `json:"issuing_date,omitempty"`

	InvoiceType   InvoiceType          `json:"invoice_type,omitempty"`
	Status        InvoiceStatus        `json:"status,omitempty"`
	PaymentStatus InvoicePaymentStatus `json:"payment_status,omitempty"`

	AmountCents          int      `json:"amount_cents,omitempty"`
	AmountCurrency       Currency `json:"amount_currency,omitempty"`
	VatAmountCents       int      `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency    Currency `json:"vat_amount_currency,omitempty"`
	CreditAmountCents    int      `json:"credit_amount_cents,omitempty"`
	CreditAmountCurrency Currency `json:"credit_amount_currency,omitempty"`
	TotalAmountCents     int      `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency  Currency `json:"total_amount_currency,omitempty"`

	FileURL  string                    `json:"file_url,omitempty"`
	Metadata []InvoiceMetadataResponse `json:"metadata,omitempty"`

	FromDate        string `json:"from_date,omitempty"`
	ToDate          string `json:"to_date,omitempty"`
	ChargesFromDate string `json:"charges_from_date,omitempty"`

	Customer      *Customer      `json:"customer,omitempty"`
	Subscriptions []Subscription `json:"subscriptions,omitempty"`

	Fees    []Fee           `json:"fees,omitempty"`
	Credits []InvoiceCredit `json:"credits,omitempty"`
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
