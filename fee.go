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
	PerPage *int `json:"per_page,omitempty,string"`
	Page    *int `json:"page,omitempty,string"`

	FeeType       FeeType          `json:"fee_type,omitempty"`
	PaymentStatus FeePaymentStatus `json:"payment_status,omitempty"`

	EventTransactionID     string `json:"event_transaction_id,omitempty"`
	ExternalSubscriptionID string `json:"external_subscription_id,omitempty"`
	ExternalCustomerID     string `json:"external_customer_id,omitempty"`

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
	Type                     FeeType                `json:"type,omitempty"`
	Code                     string                 `json:"code,omitempty"`
	Name                     string                 `json:"name,omitempty"`
	InvoiceDisplayName       string                 `json:"invoice_display_name,omitempty"`
	FilterInvoiceDisplayName string                 `json:"filter_invoice_display_name,omitempty"`
	Filters                  map[string]interface{} `json:"filters,omitempty"`
	LagoItemID               uuid.UUID              `json:"lago_item_id,omitempty"`
	ItemType                 FeeItemType            `json:"item_type,omitempty"`
	GroupedBy                map[string]interface{} `json:"grouped_by,omitempty"`
}

type FeeAppliedTax struct {
	LagoId         uuid.UUID `json:"lago_id,omitempty"`
	LagoFeeId      uuid.UUID `json:"lago_fee_id,omitempty"`
	LagoTaxId      uuid.UUID `json:"lago_tax_id,omitempty"`
	TaxName        string    `json:"tax_name,omitempty"`
	TaxCode        string    `json:"tax_code,omitempty"`
	TaxRate        float32   `json:"tax_rate,omitempty"`
	TaxDescription string    `json:"tax_description,omitempty"`
	AmountCents    int       `json:"amount_cents,omitempty"`
	AmountCurrency Currency  `json:"amount_currency,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

type PricingUnitDetails struct {
	LagoPricingUnitID  uuid.UUID `json:"lago_pricing_unit_id,omitempty"`
	PricingUnitCode    string    `json:"pricing_unit_code,omitempty"`
	ShortName          string    `json:"short_name,omitempty"`
	AmountCents        int       `json:"amount_cents,omitempty"`
	PreciseAmountCents string    `json:"precise_amount_cents,omitempty"`
	UnitAmountCents    int       `json:"unit_amount_cents,omitempty"`
	PreciseUnitAmount  string    `json:"precise_unit_amount,omitempty"`
	ConversionRate     float64   `json:"conversion_rate,omitempty"`
}

type Fee struct {
	LagoID                 uuid.UUID `json:"lago_id,omitempty"`
	LagoChargeID           uuid.UUID `json:"lago_charge_id,omitempty"`
	LagoChargeFilterID     uuid.UUID `json:"lago_charge_filter_id,omitempty"`
	LagoInvoiceID          uuid.UUID `json:"lago_invoice_id,omitempty"`
	LagoTrueUpFeeID        uuid.UUID `json:"lago_true_up_fee_id,omitempty"`
	LagoTrueUpParentFeeID  uuid.UUID `json:"lago_true_up_parent_fee_id,omitempty"`
	ExternalSubscriptionID string    `json:"external_subscription_id,omitempty"`

	AmountCents         int                    `json:"amount_cents,omitempty"`
	AmountDetails       map[string]interface{} `json:"amount_details,omitempty"`
	PreciseUnitAmount   string                 `json:"precise_unit_amount,omitempty"`
	PreciseAmount       string                 `json:"precise_amount,omitempty"`
	PreciseTotalAmount  string                 `json:"precise_total_amount,omitempty"`
	AmountCurrency      string                 `json:"amount_currency,omitempty"`
	TaxesAmountCents    int                    `json:"taxes_amount_cents,omitempty"`
	TaxesPreciseAmount  string                 `json:"taxes_precise_amount,omitempty"`
	TaxesRate           float32                `json:"taxes_rate,omitempty"`
	TotalAmountCents    int                    `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency string                 `json:"total_amount_currency,omitempty"`
	PayInAdvance        bool                   `json:"pay_in_advance,omitempty"`
	Invoiceable         bool                   `json:"invoiceable,omitempty"`
	FromDate            string                 `json:"from_date,omitempty"`
	ToDate              string                 `json:"to_date,omitempty"`
	InvoiceDisplayName  string                 `json:"invoice_display_name,omitempty"`

	TotalAggregatedUnits string `json:"total_aggregated_units,omitempty"`
	Units                string `json:"units,omitempty"`
	Description          string `json:"description,omitempty"`
	EventsCount          int    `json:"events_count,omitempty"`

	PaymentStatus FeePaymentStatus `json:"payment_status,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty"`
	SucceededAt time.Time `json:"succeeded_at,omitempty"`
	FailedAt    time.Time `json:"failed_at,omitempty"`
	RefundedAt  time.Time `json:"refunded_at,omitempty"`

	Item               FeeItem             `json:"item,omitempty"`
	AppliedTaxes       []FeeAppliedTax     `json:"applied_taxes,omitempty"`
	PricingUnitDetails *PricingUnitDetails `json:"pricing_unit_details,omitempty"`
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

func (fr *FeeRequest) Delete(ctx context.Context, feeID string) (*Fee, *Error) {
	subPath := fmt.Sprintf("%s/%s", "fees", feeID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeeResult{},
	}

	result, err := fr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	feeResult, ok := result.(*FeeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return feeResult.Fee, nil
}
