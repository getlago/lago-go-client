package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreditNoteCreditStatus string
type CreditNoteReason string

const (
	CreditNoteCreditStatusAvailable CreditNoteCreditStatus = "available"
	CreditNoteCreditStatusConsumed  CreditNoteCreditStatus = "consumed"
)

const (
	CreditNoteReasonOverpaid CreditNoteReason = "overpaid"
)

type CreditNoteRequest struct {
	client *Client
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
	LagoID         uuid.UUID  `json:"lago_id,omitempty"`
	AmountCents    int        `json:"amount_cents,omitempty"`
	AmountCurrency Currency   `json:"amount_currency,omitempty"`
	Fee            InvoiceFee `json:"fee,omitempty"`
}

type CreditNote struct {
	LagoID        uuid.UUID        `json:"lago_id,omitempty"`
	SequentialID  int              `json:"sequential_id,omitempty"`
	Number        string           `json:"number,omitempty"`
	LagoInvoiceID uuid.UUID        `json:"lago_invoice_id,omitempty"`
	InvoiceNumber string           `json:"invoice_number,omitempty"`
	Reason        CreditNoteReason `json:"reason,omitempty"`

	Status CreditNoteCreditStatus `json:"status,omitempty"`

	TotalAmountCents      int      `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency   Currency `json:"total_amount_currency,omitempty"`
	CreditAmountCents     int      `json:"credit_amount_cents,omitempty"`
	CreditAmountCurrency  Currency `json:"credit_amount_currency,omitempty"`
	BalanceAmountCents    int      `json:"balance_amount_cents,omitempty"`
	BalanceAmountCurrency Currency `json:"balance_amount_currency,omitempty"`

	FileURL string `json:"file_url,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Items []CreditNoteItem `json:"items,omitempty"`
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

	creditNoteResult := result.(*CreditNoteResult)

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
		creditNoteResult := result.(*CreditNoteResult)

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

	creditNoteResult := result.(*CreditNoteResult)

	return creditNoteResult, nil
}
