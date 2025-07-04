package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

type CreditNoteCreditStatus string
type CreditNoteRefundStatus string
type CreditNoteReason string
type CreditNoteErrorCode string

const (
	CreditNoteCreditStatusAvailable CreditNoteCreditStatus = "available"
	CreditNoteCreditStatusConsumed  CreditNoteCreditStatus = "consumed"
	CreditNoteCreditStatusVoided    CreditNoteCreditStatus = "voided"
)

const (
	CreditNoteRefundStatusPending   CreditNoteRefundStatus = "pending"
	CreditNoteRefundStatusSucceeded CreditNoteRefundStatus = "succeeded"
	CreditNoteRefundStatusFailed    CreditNoteRefundStatus = "failed"
)

const (
	CreditNoteReasonDuplicatedCharge      CreditNoteReason = "duplicated_charge"
	CreditNoteReasonProductUnsatisfactory CreditNoteReason = "product_unsatisfactory"
	CreditNoteReasonOrderChange           CreditNoteReason = "order_change"
	CreditNoteReasonOrderCancellation     CreditNoteReason = "order_cancellation"
	CreditNoteReasonFraudulentCharge      CreditNoteReason = "fraudulent_charge"
	CreditNoteReasonOther                 CreditNoteReason = "other"
)

const (
	CreditNoteErrorCodeNotProvided            CreditNoteErrorCode = "not_provided"
	CreditNoteErrorCodeTaxError               CreditNoteErrorCode = "tax_error"
	CreditNoteErrorCodeTaxVoidingError        CreditNoteErrorCode = "tax_voiding_error"
	CreditNoteErrorCodeInvoiceGenerationError CreditNoteErrorCode = "invoice_generation_error"
)

type CreditNoteRequest struct {
	client *Client
}

type CreditNoteParams struct {
	CreditNote *CreditNoteInput `json:"credit_note"`
}

type CreditNoteResult struct {
	CreditNote  *CreditNote  `json:"credit_note,omitempty"`
	CreditNotes []CreditNote `json:"credit_notes,omitempty"`
	Meta        Metadata     `json:"meta,omitempty"`
}

type EstimatedCreditNoteResult struct {
	CreditNoteEstimated *EstimatedCreditNote `json:"estimated_credit_note"`
}

type CreditNoteListInput struct {
	PerPage int `url:"per_page,omitempty,string"`
	Page    int `url:"page,omitempty,string"`

	ExternalCustomerID string `url:"external_customer_id,omitempty"`
	IssuingDateFrom    string `url:"issuing_date_from,omitempty"`
	IssuingDateTo      string `url:"issuing_date_to,omitempty"`

	AmountFrom int `url:"amount_from,omitempty,string"`
	AmountTo   int `url:"amount_to,omitempty,string"`

	SearchTerm       string                 `url:"search_term,omitempty"`
	BillingEntityIDs []uuid.UUID            `url:"billing_entity_ids[],omitempty"`
	CreditStatus     CreditNoteCreditStatus `url:"credit_status,omitempty"`
	Currency         Currency               `url:"currency,omitempty"`
	InvoiceNumber    string                 `url:"invoice_number,omitempty"`
	Reason           CreditNoteReason       `url:"reason,omitempty"`
	RefundStatus     CreditNoteRefundStatus `url:"refund_status,omitempty"`
	SelfBilled       *bool                  `url:"self_billed,omitempty,string"`
}

type CreditNoteItem struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	AmountCents    int       `json:"amount_cents,omitempty"`
	AmountCurrency Currency  `json:"amount_currency,omitempty"`
	Fee            Fee       `json:"fee,omitempty"`
}

type CreditNoteAppliedTax struct {
	LagoID           uuid.UUID `json:"lago_id,omitempty"`
	LagoCreditNoteID uuid.UUID `json:"lago_credit_note_id,omitempty"`
	LagoTaxID        uuid.UUID `json:"lago_tax_id,omitempty"`
	TaxName          string    `json:"tax_name,omitempty"`
	TaxCode          string    `json:"tax_code,omitempty"`
	TaxRate          float64   `json:"tax_rate,omitempty"`
	TaxDescription   string    `json:"tax_description,omitempty"`
	AmountCents      int       `json:"amount_cents,omitempty"`
	AmountCurrency   Currency  `json:"amount_currency,omitempty"`
	BaseAmountCents  int       `json:"base_amount_cents,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type CreditNoteErrorDetails struct {
	LagoID    uuid.UUID              `json:"lago_id,omitempty"`
	ErrorCode CreditNoteErrorCode    `json:"error_code,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type CreditNote struct {
	LagoID            uuid.UUID        `json:"lago_id,omitempty"`
	SequentialID      int              `json:"sequential_id,omitempty"`
	BillingEntityCode string           `json:"billing_entity_code,omitempty"`
	Number            string           `json:"number,omitempty"`
	LagoInvoiceID     uuid.UUID        `json:"lago_invoice_id,omitempty"`
	InvoiceNumber     string           `json:"invoice_number,omitempty"`
	Reason            CreditNoteReason `json:"reason,omitempty"`
	Description       string           `json:"description,omitempty"`

	SelfBilled   bool                   `json:"self_billed,omitempty"`
	CreditStatus CreditNoteCreditStatus `json:"credit_status,omitempty"`
	RefundStatus CreditNoteRefundStatus `json:"refund_status,omitempty"`

	Currency                          Currency `json:"currency,omitempty"`
	TotalAmountCents                  int      `json:"total_amount_cents,omitempty"`
	CreditAmountCents                 int      `json:"credit_amount_cents,omitempty"`
	BalanceAmountCents                int      `json:"balance_amount_cents,omitempty"`
	RefundAmountCents                 int      `json:"refund_amount_cents,omitempty"`
	TaxesAmountCents                  int      `json:"taxes_amount_cents,omitempty"`
	TaxesRate                         float64  `json:"taxes_rate,omitempty"`
	SubTotalExcludingTaxesAmountCents int      `json:"sub_total_excluding_taxes_amount_cents,omitempty"`
	CouponsAdjustmentAmountCents      int      `json:"coupons_adjustment_amount_cents,omitempty"`

	FileURL string `json:"file_url,omitempty"`

	IssuingDate string    `json:"issuing_date,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`

	Items        []CreditNoteItem         `json:"items,omitempty"`
	AppliedTaxes []CreditNoteAppliedTax   `json:"applied_taxes,omitempty"`
	ErrorDetails []CreditNoteErrorDetails `json:"error_details,omitempty"`
}

type EstimatedCreditNote struct {
	LagoInvoiceID uuid.UUID `json:"lago_invoice_id,omitempty"`
	InvoiceNumber string    `json:"invoice_number,omitempty"`

	Currency                            Currency `json:"currency,omitempty"`
	MaxCreditableAmountCents            int      `json:"max_creditable_amount_cents,omitempty"`
	MaxRefundableAmountCents            int      `json:"max_refundable_amount_cents,omitempty"`
	TaxesAmountCents                    int      `json:"taxes_amount_cents,omitempty"`
	PreciseTaxesAmountCents             float64  `json:"precise_taxes_amount_cents,omitempty"`
	TaxesRate                           float64  `json:"taxes_rate,omitempty"`
	SubTotalExcludingTaxesAmountCents   int      `json:"sub_total_excluding_taxes_amount_cents,omitempty"`
	CouponsAdjustmentAmountCents        int      `json:"coupons_adjustment_amount_cents,omitempty"`
	PreciseCouponsAdjustmentAmountCents float64  `json:"precise_coupons_adjustment_amount_cents,omitempty"`

	Items []EstimatedCreditNoteItem `json:"items,omitempty"`

	AppliedTaxes []EstimatedCreditNoteAppliedTax `json:"applied_taxes,omitempty"`
}

type EstimatedCreditNoteItem struct {
	AmountCents int       `json:"amount_cents,omitempty"`
	LagoFeeID   uuid.UUID `json:"lago_fee_id,omitempty"`
}

type EstimatedCreditNoteAppliedTax struct {
	LagoTaxID       uuid.UUID `json:"lago_tax_id,omitempty"`
	TaxName         string    `json:"tax_name,omitempty"`
	TaxCode         string    `json:"tax_code,omitempty"`
	TaxRate         float64   `json:"tax_rate,omitempty"`
	TaxDescription  string    `json:"tax_description,omitempty"`
	AmountCents     int       `json:"amount_cents,omitempty"`
	AmountCurrency  Currency  `json:"amount_currency,omitempty"`
	BaseAmountCents int       `json:"base_amount_cents,omitempty"`
}

type CreditNoteItemInput struct {
	LagoFeeID   uuid.UUID `json:"fee_id,omitempty"`
	AmountCents int       `json:"amount_cents,omitempty"`
}

type CreditNoteInput struct {
	LagoInvoiceID     uuid.UUID             `json:"invoice_id,omitempty"`
	Reason            CreditNoteReason      `json:"reason,omitempty"`
	Description       *string               `json:"description,omitempty"`
	Items             []CreditNoteItemInput `json:"items,omitempty"`
	CreditAmountCents int                   `json:"credit_amount_cents,omitempty"`
	RefundAmountCents int                   `json:"refund_amount_cents,omitempty"`
}

type CreditNoteUpdateInput struct {
	LagoID       uuid.UUID
	RefundStatus CreditNoteRefundStatus `json:"refund_status,omitempty"`
}

type CreditNoteUpdateParams struct {
	CreditNote *CreditNoteUpdateInput `json:"credit_note"`
}

type EstimateCreditNoteInput struct {
	LagoInvoiceID uuid.UUID             `json:"invoice_id,omitempty"`
	Items         []CreditNoteItemInput `json:"items,omitempty"`
}

type EstimateCreditNoteParams struct {
	CreditNote *EstimateCreditNoteInput `json:"credit_note"`
}

func (c *Client) CreditNote() *CreditNoteRequest {
	return &CreditNoteRequest{
		client: c,
	}
}

func (cr *CreditNoteRequest) Get(ctx context.Context, creditNoteID uuid.UUID) (*CreditNote, *Error) {
	subPath := fmt.Sprintf("%s/%s", "credit_notes", creditNoteID)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CreditNoteResult{},
	}

	result, err := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	creditNoteResult, ok := result.(*CreditNoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return creditNoteResult.CreditNote, nil
}

func (cr *CreditNoteRequest) Download(ctx context.Context, creditNoteID uuid.UUID) (*CreditNote, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "credit_notes", creditNoteID, "download")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CreditNoteResult{},
	}

	result, err := cr.client.PostWithoutBody(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result != nil {
		creditNoteResult, ok := result.(*CreditNoteResult)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return creditNoteResult.CreditNote, nil
	}

	return nil, nil
}

func (cr *CreditNoteRequest) GetList(ctx context.Context, creditNoteListInput *CreditNoteListInput) (*CreditNoteResult, *Error) {
	urlValues, err := query.Values(creditNoteListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      "credit_notes",
		UrlValues: urlValues,
		Result:    &CreditNoteResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	creditNoteResult, ok := result.(*CreditNoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return creditNoteResult, nil
}

func (cr *CreditNoteRequest) Create(ctx context.Context, creditNoteInput *CreditNoteInput) (*CreditNote, *Error) {
	creditNoteParams := &CreditNoteParams{
		CreditNote: creditNoteInput,
	}

	clientRequest := &ClientRequest{
		Path:   "credit_notes",
		Result: &CreditNoteResult{},
		Body:   creditNoteParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	creditNoteResult, ok := result.(*CreditNoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return creditNoteResult.CreditNote, nil
}

func (cr *CreditNoteRequest) Update(ctx context.Context, creditNoteUpdateInput *CreditNoteUpdateInput) (*CreditNote, *Error) {
	subPath := fmt.Sprintf("%s/%s", "credit_notes", creditNoteUpdateInput.LagoID)
	creditNoteParams := &CreditNoteUpdateParams{
		CreditNote: creditNoteUpdateInput,
	}

	ClientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CreditNoteResult{},
		Body:   creditNoteParams,
	}

	result, err := cr.client.Put(ctx, ClientRequest)
	if err != nil {
		return nil, err
	}

	creditNoteResult, ok := result.(*CreditNoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return creditNoteResult.CreditNote, nil
}

func (cr *CreditNoteRequest) Void(ctx context.Context, creditNoteID uuid.UUID) (*CreditNote, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "credit_notes", creditNoteID, "void")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CreditNoteResult{},
	}

	result, err := cr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	creditNoteResult, ok := result.(*CreditNoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return creditNoteResult.CreditNote, nil
}

func (cr *CreditNoteRequest) Estimate(ctx context.Context, estimateCreditNoteInput *EstimateCreditNoteInput) (*EstimatedCreditNote, *Error) {
	estimateCreditNoteParams := &EstimateCreditNoteParams{
		CreditNote: estimateCreditNoteInput,
	}

	clientRequest := &ClientRequest{
		Path:   "credit_notes/estimate",
		Result: &EstimatedCreditNoteResult{},
		Body:   estimateCreditNoteParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	estimatedCreditNoteResult, ok := result.(*EstimatedCreditNoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return estimatedCreditNoteResult.CreditNoteEstimated, nil
}
