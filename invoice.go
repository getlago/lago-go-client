package lago

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

type InvoiceType string
type InvoiceStatus string
type InvoicePaymentStatus string
type InvoiceCreditItemType string

const (
	SubscriptionInvoiceType       InvoiceType = "subscription"
	AddOnInvoiceType              InvoiceType = "add_on"
	CreditInvoiceType             InvoiceType = "credit"
	OneOffInvoiceType             InvoiceType = "one_off"
	ProgressiveBillingInvoiceType InvoiceType = "progressive_billing"
)

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusFinalized InvoiceStatus = "finalized"
	InvoiceStatusFailed    InvoiceStatus = "failed"
	InvoiceStatusVoided    InvoiceStatus = "voided"
	InvoiceStatusPending   InvoiceStatus = "pending"
)

const (
	InvoicePaymentStatusPending   InvoicePaymentStatus = "pending"
	InvoicePaymentStatusSucceeded InvoicePaymentStatus = "succeeded"
	InvoicePaymentStatusFailed    InvoicePaymentStatus = "failed"
)

const (
	InvoiceCreditItemCoupon     InvoiceCreditItemType = "coupon"
	InvoiceCreditItemCreditNote InvoiceCreditItemType = "credit_note"
	InvoiceCreditItemInvoice    InvoiceCreditItemType = "invoice"
)

type InvoiceRequest struct {
	client *Client
}

type InvoiceResult struct {
	Invoice  *Invoice  `json:"invoice,omitempty"`
	Invoices []Invoice `json:"invoices,omitempty"`
	Meta     Metadata  `json:"meta,omitempty"`
}

type InvoicePaymentDetailsResult struct {
	InvoicePaymentDetails *InvoicePaymentDetails `json:"invoice_payment_details,omitempty"`
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
	AddOnCode          string     `json:"add_on_code,omitempty"`
	InvoiceDisplayName string     `json:"invoice_display_name,omitempty"`
	UnitAmountCents    int        `json:"unit_amount_cents,omitempty"`
	Description        string     `json:"description,omitempty"`
	FromDatetime       *time.Time `json:"from_datetime,omitempty"`
	ToDatetime         *time.Time `json:"to_datetime,omitempty"`
	Units              float32    `json:"units,omitempty"`
	TaxCodes           []string   `json:"tax_codes,omitempty"`
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
	SkipPsp            bool               `json:"skip_psp,omitempty"`
}

type InvoicePreviewInput struct {
	PlanCode          string              `json:"plan_code,omitempty"`
	BillingTime       string              `json:"billing_time,omitempty"`
	SubscriptionAt    string              `json:"subscription_at,omitempty"`
	Coupons           []CouponInput       `json:"coupons,omitempty"`
	Customer          *CustomerInput      `json:"customer,omitempty"`
	Subscriptions     *SubscriptionsInput `json:"subscriptions,omitempty"`
	BillingEntityCode string              `json:"billing_entity_code,omitempty"`
}

type InvoiceListInputMetadata map[string]any
type InvoiceListInput struct {
	PerPage *int `url:"per_page,omitempty"`
	Page    *int `url:"page,omitempty"`

	IssuingDateFrom string `url:"issuing_date_from,omitempty"`
	IssuingDateTo   string `url:"issuing_date_to,omitempty"`

	// NOTE: Expose the fields as ExternalCustomerID to keep consistency with other endpoints
	ExternalCustomerID string `url:"customer_external_id,omitempty"`

	InvoiceType   InvoiceType          `url:"invoice_type,omitempty"`
	Status        InvoiceStatus        `url:"status,omitempty"`
	PaymentStatus InvoicePaymentStatus `url:"payment_status,omitempty"`

	PaymentOverdue     *bool `url:"payment_overdue,omitempty,string"`
	PartiallyPaid      *bool `url:"partially_paid,omitempty,string"`
	SelfBilled         *bool `url:"self_billed,omitempty,string"`
	PaymentDisputeLost *bool `url:"payment_dispute_lost,omitempty,string"`

	AmountFrom *int `url:"amount_from,omitempty"`
	AmountTo   *int `url:"amount_to,omitempty"`

	SearchTerm       string      `url:"search_term,omitempty"`
	BillingEntityIDs []uuid.UUID `url:"billing_entity_ids[],omitempty"`
	Currency         Currency    `url:"currency,omitempty"`

	Metadata *InvoiceListInputMetadata `url:"metadata,omitempty"`
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

	LagoItemID     uuid.UUID `json:"lago_item_id,omitempty"`
	AmountCents    int       `json:"amount_cents,omitempty"`
	AmountCurrency Currency  `json:"amount_currency,omitempty"`
	BeforeTaxes    bool      `json:"before_taxes,omitempty"`
}

type InvoiceAppliedInvoiceCustomSection struct {
	LagoId        uuid.UUID `json:"lago_id,omitempty"`
	LagoInvoiceId uuid.UUID `json:"lago_invoice_id,omitempty"`
	Code          string    `json:"code,omitempty"`
	Details       string    `json:"details,omitempty"`
	DisplayName   string    `json:"display_name,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
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

type InvoiceErrorDetail struct {
	LagoId    uuid.UUID      `json:"lago_id,omitempty"`
	ErrorCode string         `json:"error_code,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}

type Invoice struct {
	LagoID            uuid.UUID `json:"lago_id,omitempty"`
	SequentialID      int       `json:"sequential_id,omitempty"`
	BillingEntityCode string    `json:"billing_entity_code,omitempty"`
	Number            string    `json:"number,omitempty"`

	IssuingDate          string    `json:"issuing_date,omitempty"`
	PaymentDisputeLostAt time.Time `json:"payment_dispute_lost_at,omitempty"`
	PaymentDueDate       string    `json:"payment_due_date,omitempty"`
	PaymentOverdue       bool      `json:"payment_overdue,omitempty"`

	InvoiceType   InvoiceType          `json:"invoice_type,omitempty"`
	Status        InvoiceStatus        `json:"status,omitempty"`
	PaymentStatus InvoicePaymentStatus `json:"payment_status,omitempty"`

	Currency Currency `json:"currency,omitempty"`

	FeesAmountCents                     int `json:"fees_amount_cents,omitempty"`
	TaxesAmountCents                    int `json:"taxes_amount_cents,omitempty"`
	CouponsAmountCents                  int `json:"coupons_amount_cents,omitempty"`
	CreditNotesAmountCents              int `json:"credit_notes_amount_cents,omitempty"`
	SubTotalExcludingTaxesAmountCents   int `json:"sub_total_excluding_taxes_amount_cents,omitempty"`
	SubTotalIncludingTaxesAmountCents   int `json:"sub_total_including_taxes_amount_cents,omitempty"`
	TotalAmountCents                    int `json:"total_amount_cents,omitempty"`
	TotalDueAmountCents                 int `json:"total_due_amount_cents,omitempty"`
	PrepaidCreditAmountCents            int `json:"prepaid_credit_amount_cents,omitempty"`
	ProgressiveBillingCreditAmountCents int `json:"progressive_billing_credit_amount_cents"`
	NetPaymentTerm                      int `json:"net_payment_term,omitempty"`

	FileURL       string                    `json:"file_url,omitempty"`
	Metadata      []InvoiceMetadataResponse `json:"metadata,omitempty"`
	VersionNumber int                       `json:"version_number,omitempty"`

	Customer       *Customer       `json:"customer,omitempty"`
	BillingPeriods []BillingPeriod `json:"billing_periods,omitempty"`
	Subscriptions  []Subscription  `json:"subscriptions,omitempty"`

	Fees                         []Fee                                `json:"fees,omitempty"`
	Credits                      []InvoiceCredit                      `json:"credits,omitempty"`
	AppliedInvoiceCustomSections []InvoiceAppliedInvoiceCustomSection `json:"applied_invoice_custom_sections,omitempty"`
	AppliedTaxes                 []InvoiceAppliedTax                  `json:"applied_taxes,omitempty"`
	ErrorDetails                 []InvoiceErrorDetail                 `json:"error_details,omitempty"`
	AppliedUsageThreshold        []AppliedUsageThreshold              `json:"applied_usage_threshold,omitempty"`
}

type InvoicePaymentDetails struct {
	LagoCustomerID     uuid.UUID `json:"lago_customer_id,omitempty"`
	LagoInvoiceID      uuid.UUID `json:"lago_invoice_id,omitempty"`
	ExternalCustomerID string    `json:"external_customer_id,omitempty"`
	PaymentProvider    string    `json:"payment_provider,omitempty"`
	PaymentUrl         string    `json:"payment_url,omitempty"`
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
	urlValues, err := query.Values(invoiceListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      "invoices",
		UrlValues: urlValues,
		Result:    &InvoiceResult{},
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

func (ir *InvoiceRequest) Preview(ctx context.Context, invoicePreviewInput *InvoicePreviewInput) (*Invoice, *Error) {
	clientRequest := &ClientRequest{
		Path:   "invoices/preview",
		Result: &InvoiceResult{},
		Body:   invoicePreviewInput,
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

func (ir *InvoiceRequest) Retry(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "retry")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	result, err := ir.client.Post(ctx, clientRequest)
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

type VoidInvoiceOptions struct {
	GenerateCreditNote bool `json:"generate_credit_note,omitempty"`
	RefundAmount       int  `json:"refund_amount,omitempty"`
	CreditAmount       int  `json:"credit_amount,omitempty"`
}

func (ir *InvoiceRequest) Void(ctx context.Context, invoiceID string, opts *VoidInvoiceOptions) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "void")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	if opts != nil {
		clientRequest.Body = opts
	}

	result, err := ir.client.Post(ctx, clientRequest)
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

func (ir *InvoiceRequest) LoseDispute(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "lose_dispute")
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

// We have Invoice as a possible return to be consitent with other endpoints, but no Invoice will be returned.
func (ir *InvoiceRequest) RetryPayment(ctx context.Context, invoiceID string) (*Invoice, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "retry_payment")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoiceResult{},
	}

	// We don't return an invoice here due to async retry payment processing
	_, err := ir.client.PostWithoutBody(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ir *InvoiceRequest) PaymentUrl(ctx context.Context, invoiceID string) (*InvoicePaymentDetails, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "invoices", invoiceID, "payment_url")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &InvoicePaymentDetailsResult{},
	}

	result, err := ir.client.PostWithoutBody(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentUrlResult, ok := result.(*InvoicePaymentDetailsResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentUrlResult.InvoicePaymentDetails, nil
}

// URL encoder for metadata type on GetList
// it will convert a map[string]any to a query string metadata[key1]=value1&metadata[key2]=value2
func (m InvoiceListInputMetadata) EncodeValues(key string, values *url.Values) error {
	for k, v := range m {
		metadataKey := fmt.Sprintf("metadata[%s]", k)
		values.Set(metadataKey, fmt.Sprintf("%v", v))
	}
	return nil
}
