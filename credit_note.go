package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreditNoteStatus string
type CreditNoteReason string

const (
	CreditNoteStatusAvailable CreditNoteStatus = "available"
	CreditNoteStatusConsumed  CreditNoteStatus = "consumed"
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
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type CreditNoteItem struct {
	LagoID         uuid.UUID  `json:"lago_id,omitempty"`
	AmountCents    int        `json:"amount_cents,omitempty"`
	AmountCurrency Currency   `json:"amount_currency,omitempty"`
	Fee            InvoiceFee `json:"fee,omitempty"`
}

type CreditNote struct {
	LagoID                  uuid.UUID        `json:"lago_id,omitempty"`
	SequentialID            int              `json:"sequential_id,omitempty"`
	Number                  string           `json:"number,omitempty"`
	LagoInvoiceID           uuid.UUID        `json:"lago_invoice_id,omitempty"`
	InvoiceNumber           string           `json:"invoice_number,omitempty"`
	Status                  CreditNoteStatus `json:"status,omitempty"`
	Reason                  CreditNoteReason `json:"reason,omitempty"`
	AmountCents             int              `json:"amount_cents,omitempty"`
	AmountCurrency          Currency         `json:"amount_currency,omitempty"`
	RemainingAmountCents    int              `json:"remaining_amount_cents,omitempty"`
	RemainingAmountCurrency Currency         `json:"remaining_amount_currency,omitempty"`
	CreatedAt               time.Time        `json:"created_at,omitempty"`
	UpdatedAt               time.Time        `json:"updated_at,omitempty"`

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
