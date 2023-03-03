package lago

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type FeeItemType string

const (
	FeeItemSubscription FeeItemType = "subscription"
	FeeItemCharge       FeeItemType = "charge"
	FeeItemAddOn        FeeItemType = "add_on"
)

type FeeRequest struct {
	client *Client
}

type FeeResult struct {
	Fee  *Fee     `json:"fee,omitempty"`
	Meta Metadata `json:"meta,omitempty"`
}

type FeeItem struct {
	Type FeeItemType `json:"type,omitempty"`
	Code string      `json:"code,omitempty"`
	Name string      `json:"name,omitempty"`
}

type Fee struct {
	LagoID      uuid.UUID `json:"lago_id,omitempty"`
	LagoGroupID uuid.UUID `json:"lago_group_id,omitempty"`

	AmountCents       int    `json:"amount_cents,omitempty"`
	AmountCurrency    string `json:"amount_currenty,omitempty"`
	VatAmountCents    int    `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency string `json:"vat_amount_currency,omitempty"`

	Units       float32 `json:"units,omitempty"`
	EventsCount int     `json:"events_count,omitempty"`

	Item FeeItem `json:"item,omitempty"`
}

func (c *Client) Fee() *FeeRequest {
	return &FeeRequest{
		client: c,
	}
}

func (fr *FeeRequest) Get(ctx context.Context, feeID string) (*Fee, *Error) {
	subPath := fmt.Sprintf("%s/%s", "fees", feeID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeeResult{},
	}

	result, err := fr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	feeResult, ok := result.(*FeeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return feeResult.Fee, nil
}
