package lago

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type InvoiceStatus string
type InvoiceFeeItemType string
type InvoiceCreditItemType string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusSucceeded InvoiceStatus = "succeeded"
	InvoiceStatusFailed    InvoiceStatus = "failed"
)

const (
	InvoiceFeeItemSubscription InvoiceFeeItemType = "subscription"
	InvoiceFeeItemCharge       InvoiceFeeItemType = "charge"
	InvoiceFeeItemAddOn        InvoiceFeeItemType = "add_on"
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

type InvoiceInput struct {
	LagoID uuid.UUID     `json:"lago_id,omitempty"`
	Status InvoiceStatus `json:"status,omitempty"`
}

type InvoiceListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	IssuingDateFrom string `json:"issuing_date_from,omitempty"`
	IssuingDateTo   string `json:"issuing_date_to,omitempty"`
}

type InvoiceFeeItem struct {
	Type InvoiceFeeItemType `json:"type,omitempty"`
	Code string             `json:"code,omitempty"`
	Name string             `json:"name,omitempty"`
}

type InvoiceFee struct {
	Item InvoiceFeeItem `json:"item,omitempty"`

	AmountCents       int      `json:"amount_cents,omitempty"`
	AmountCurrency    Currency `json:"amount_currency,omitempty"`
	VatAmountCents    int      `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency Currency `json:"vat_amount_currency,omitempty"`
}

type InvoiceCreditItem struct {
	Type InvoiceCreditItemType `json:"type,omitempty"`
	Code string                `json:"code,omitempty"`
	Name string                `json:"name,omitempty"`
}

type InvoiceCredit struct {
	Item InvoiceCreditItem `json:"item,omitempty"`

	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`
}

type Invoice struct {
	LagoID       uuid.UUID `json:"lago_id,omitempty"`
	SequentialID int       `json:"sequential_id,omitempty"`
	Number       string    `json:"number,omitempty"`

	Status InvoiceStatus `json:"status,omitempty"`

	AmountCents       int      `json:"amount_cents,omitempty"`
	AmountCurrency    Currency `json:"amount_currency,omitempty"`
	VatAmountCents    int      `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency Currency `json:"vat_amount_currency,omitempty"`

	FileURL string `json:"file_url,omitempty"`

	FromDate        string `json:"from_date,omitempty"`
	ToDate          string `json:"to_date,omitempty"`
	ChargesFromDate string `json:"charges_from_date,omitempty"`
	IssuingDate     string `json:"issuing_date,omitempty"`

	Customer     *Customer      `json:"customer,omitempty"`
	Subscription []Subscription `json:"subscriptions,omitempty"`

	Fees    []InvoiceFee    `json:"fees,omitempty"`
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

	invoiceResult := result.(*InvoiceResult)

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

	invoiceResult := result.(*InvoiceResult)

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

	invoiceResult := result.(*InvoiceResult)

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
		invoiceResult := result.(*InvoiceResult)

		return invoiceResult.Invoice, nil
	}

	return nil, nil
}
