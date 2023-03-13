package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreditNoteCreditStatus string
type CreditNoteRefundStatus string
type CreditNoteReason string

const (
	CreditNoteCreditStatusAvailable CreditNoteCreditStatus = "available"
	CreditNoteCreditStatusConsumed  CreditNoteCreditStatus = "consumed"
)

const (
	CreditNoteRefundStatusPending  CreditNoteRefundStatus = "pending"
	CreditNoteRefundStatusRefunded CreditNoteRefundStatus = "refunded"
)

const (
	CreditNoteReasonDuplicatedCharge      CreditNoteReason = "duplicated_charge"
	CreditNoteReasonProductUnsatisfactory CreditNoteReason = "product_unsatisfactory"
	CreditNoteReasonOrderChange           CreditNoteReason = "order_change"
	CreditNoteReasonOrderCancellation     CreditNoteReason = "order_cancellation"
	CreditNoteReasonFraudulentCharge      CreditNoteReason = "fraudulent_charge"
	CreditNoteReasonOther                 CreditNoteReason = "other"
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

type CreditListInput struct {
	PerPage            int    `json:"per_page,omitempty,string"`
	Page               int    `json:"page,omitempty,string"`
	ExternalCustomerID string `json:"external_customer_id,omitempty"`
}

type CreditNoteItem struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	AmountCents    int       `json:"amount_cents,omitempty"`
	AmountCurrency Currency  `json:"amount_currency,omitempty"`
	Fee            Fee       `json:"fee,omitempty"`
}

type CreditNote struct {
	LagoID        uuid.UUID        `json:"lago_id,omitempty"`
	SequentialID  int              `json:"sequential_id,omitempty"`
	Number        string           `json:"number,omitempty"`
	LagoInvoiceID uuid.UUID        `json:"lago_invoice_id,omitempty"`
	InvoiceNumber string           `json:"invoice_number,omitempty"`
	Reason        CreditNoteReason `json:"reason,omitempty"`

	CreditStatus CreditNoteCreditStatus `json:"credit_status,omitempty"`
	RefundStatus CreditNoteRefundStatus `json:"refund_status,omitempty"`

	TotalAmountCents                  int      `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency               Currency `json:"total_amount_currency,omitempty"`
	CreditAmountCents                 int      `json:"credit_amount_cents,omitempty"`
	CreditAmountCurrency              Currency `json:"credit_amount_currency,omitempty"`
	BalanceAmountCents                int      `json:"balance_amount_cents,omitempty"`
	BalanceAmountCurrency             Currency `json:"balance_amount_currency,omitempty"`
	RefundAmountCents                 int      `json:"refund_amount_cents,omitempty"`
	RefundAmountCurrency              Currency `json:"refund_amount_currency,omitempty"`
	VatAmountCents                    int      `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency                 Currency `json:"vat_amount_currency,omitempty"`
	SubTotalVatExcludedAmountCents    int      `json:"sub_total_vat_excluded_amount_cents,omitempty"`
	SubTotalVatExcludedAmountCurrency Currency `json:"sub_total_vat_excluded_amount_currency,omitempty"`

	FileURL string `json:"file_url,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Items []CreditNoteItem `json:"items,omitempty"`
}

type CreditNoteItemInput struct {
	LagoFeeID   uuid.UUID `json:"fee_id,omitempty"`
	AmountCents int       `json:"amount_cents,omitempty"`
}

type CreditNoteInput struct {
	LagoInvoiceID     uuid.UUID             `json:"invoice_id,omitempty"`
	Reason            CreditNoteReason      `json:"reason,omitempty"`
	Items             []CreditNoteItemInput `json:"items,omitempty"`
	CreditAmountCents int                   `json:"refund_amount_cents,omitempty"`
	RefundAmountCents int                   `json:"credit_amount_cents,omitempty"`
}

type CreditNoteUpdateInput struct {
	LagoID       string                 `json:"id,omitempty"`
	RefundStatus CreditNoteRefundStatus `json:"refund_status,omitempty"`
}

type CreditNoteUpdateParams struct {
	CreditNote *CreditNoteUpdateInput `json:"credit_note"`
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

func (cr *CreditNoteRequest) Download(ctx context.Context, creditNoteID string) (*CreditNote, *Error) {
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

func (cr *CreditNoteRequest) GetList(ctx context.Context, creditNoteListInput *CreditListInput) (*CreditNoteResult, *Error) {
	jsonQueryParams, err := json.Marshal(creditNoteListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "credit_notes",
		QueryParams: queryParams,
		Result:      &CreditNoteResult{},
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
		Result: &PlanResult{},
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

func (cr *CreditNoteRequest) Void(ctx context.Context, creditNoteID string) (*CreditNote, *Error) {
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
