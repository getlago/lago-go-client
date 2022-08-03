package lago

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string
type InvoiceFeeItemType string
type InvoiceCreditItemType string

const (
	InvoicePending   InvoiceStatus = "pending"
	InvoiceSucceeded InvoiceStatus = "succeeded"
	InvoiceFailed    InvoiceStatus = "failed"
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

type InvoiceListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	IssuingDateFrom *time.Time `json:"issuing_date_from,omitempty"`
	IssuingDateTo   *time.Time `json:"issuing_date_to,omitempty"`
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

	Customer     *Customer     `json:"customer,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`

	Fees    []InvoiceFee    `json:"fees,omitempty"`
	Credits []InvoiceCredit `json:"credits,omitempty"`
}

func (c *Client) Invoice() *InvoiceRequest {
	return &InvoiceRequest{
		client: c,
	}
}

func (ir *InvoiceRequest) Get(invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s", "invoices", invoiceID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	invoiceResult := result.(*InvoiceResult)

	return invoiceResult.Invoice, nil
}

func (ir *InvoiceRequest) GetList(invoiceListInput *InvoiceListInput) (*InvoiceResult, *Error) {
	jsonQueryParams, err := json.Marshal(invoiceListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	json.Unmarshal(jsonQueryParams, &queryParams)

	clientRequest := &ClientRequest{
		Path:        "invoices",
		QueryParams: queryParams,
		Result:      &InvoiceResult{},
	}

	result, clientErr := ir.client.Get(clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	invoiceResult := result.(*InvoiceResult)

	return invoiceResult, nil
}
