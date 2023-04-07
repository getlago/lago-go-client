package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FeeType string
type FeePaymentStatus string
type FeeItemType string

const (
	FeeItemSubscription FeeType = "subscription"
	FeeItemCharge       FeeType = "charge"
	FeeItemAddOn        FeeType = "add_on"
)

const (
	FeePaymentStatusPending   FeePaymentStatus = "pending"
	FeePaymentStatusSucceeded FeePaymentStatus = "succeeded"
	FeePaymentStatusFailed    FeePaymentStatus = "failed"
	FeePaymentStatusRefunded  FeePaymentStatus = "refunded"
)

const (
	FeeAddOnType         FeeItemType = "AddOn"
	FeeBillableMetric    FeeItemType = "BillableMetric"
	FeeSubscription      FeeItemType = "Subscription"
	FeeWalletTransaction FeeItemType = "WalletTransaction"
)

type FeeRequest struct {
	client *Client
}

type FeeResult struct {
	Fee  *Fee     `json:"fee,omitempty"`
	Fees []Fee    `json:"fees,omitempty"`
	Meta Metadata `json:"meta,omitempty"`
}

type FeeUpdateParams struct {
	Fee *FeeUpdateInput `json:"fee"`
}

type FeeUpdateInput struct {
	LagoID        uuid.UUID        `json:"lago_id,omitempty"`
	PaymentStatus FeePaymentStatus `json:"payment_status,omitempty"`
}

type FeeListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`

	FeeType       FeeType          `json:"fee_type,omitempty"`
	PaymentStatus FeePaymentStatus `json:"payment_status,omitempty"`

	ExternalSubscriptionId string `json:"external_subscription_id,omitempty"`
	ExternalCustomerId     string `json:"external_customer_id,omitempty"`

	BillableMetricCode string `json:"billable_metric_code,omitempty"`

	Currency Currency `json:"currency"`

	CreatedAtFrom   string `json:"created_at_from,omitempty"`
	CreatedAtTo     string `json:"created_at_to,omitempty"`
	FailedAtFrom    string `json:"failed_at_from,omitempty"`
	FailedAtTo      string `json:"failed_at_to,omitempty"`
	SucceededAtFrom string `json:"succeeded_at_from,omitempty"`
	SucceededAtTo   string `json:"succeeded_at_to,omitempty"`
	RefundedAtFrom  string `json:"refunded_at_from,omitempty"`
	RefundedAtTo    string `json:"refunded_at_to,omitempty"`
}

type FeeItem struct {
	Type       FeeType     `json:"type,omitempty"`
	Code       string      `json:"code,omitempty"`
	Name       string      `json:"name,omitempty"`
	LagoItemID uuid.UUID   `json:"lago_item_id,omitempty"`
	ItemType   FeeItemType `json:"item_type,omitempty"`
}

type Fee struct {
	LagoID                 uuid.UUID `json:"lago_id,omitempty"`
	LagoGroupID            uuid.UUID `json:"lago_group_id,omitempty"`
	LagoInvoiceID          uuid.UUID `json:"lago_invoice_id,omitempty"`
	ExternalSubscriptionID string    `json:"external_subscription_id,omitempty"`

	AmountCents         int    `json:"amount_cents,omitempty"`
	AmountCurrency      string `json:"amount_currenty,omitempty"`
	VatAmountCents      int    `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency   string `json:"vat_amount_currency,omitempty"`
	TotalAmountCents    int    `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency string `json:"total_amount_currency,omitempty"`
	FromDate            string `json:"from_date,omitempty"`
	ToDate              string `json:"to_date,omitempty"`

	Units       string `json:"units,omitempty"`
	EventsCount int    `json:"events_count,omitempty"`

	PaymentStatus FeePaymentStatus `json:"payment_status,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty"`
	SucceededAt time.Time `json:"succeeded_at,omitempty"`
	FailedAt    time.Time `json:"failed_at,omitempty"`
	RefundedAt  time.Time `json:"refunded_at,omitempty"`

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

func (fr *FeeRequest) Update(ctx context.Context, feeInput *FeeUpdateInput) (*Fee, *Error) {
	subPath := fmt.Sprintf("%s/%s", "fees", feeInput.LagoID)
	feeParams := &FeeUpdateParams{
		Fee: feeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeeResult{},
		Body:   feeParams,
	}

	result, err := fr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	feeResult, ok := result.(*FeeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return feeResult.Fee, nil
}

func (fr *FeeRequest) GetList(ctx context.Context, feeListInput *FeeListInput) (*FeeResult, *Error) {
	jsonQueryParams, err := json.Marshal(feeListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "fees",
		QueryParams: queryParams,
		Result:      &FeeResult{},
	}

	result, clientErr := fr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	feeResult, ok := result.(*FeeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return feeResult, nil
}
