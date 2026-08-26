package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

const OrdersEndpoint string = "orders"

type OrderStatus string

const (
	OrderStatusCreated  OrderStatus = "created"
	OrderStatusExecuted OrderStatus = "executed"
	OrderStatusFailed   OrderStatus = "failed"
)

type OrderRequest struct {
	client *Client
}

// A trace of what the execution produced. Every key is always present, whichever the
// order type, so the ones that do not apply keep their empty value.
type OrderExecutionRecord struct {
	ExecutedAt    *time.Time         `json:"executed_at,omitempty"`
	ExecutionMode OrderExecutionMode `json:"execution_mode,omitempty"`

	InvoiceID                 *uuid.UUID  `json:"invoice_id,omitempty"`
	SubscriptionIDs           []uuid.UUID `json:"subscription_ids,omitempty"`
	TerminatedSubscriptionIDs []uuid.UUID `json:"terminated_subscription_ids,omitempty"`
	AppliedCouponIDs          []uuid.UUID `json:"applied_coupon_ids,omitempty"`
	WalletIDs                 []uuid.UUID `json:"wallet_ids,omitempty"`

	Errors []string `json:"errors,omitempty"`
}

type Order struct {
	LagoID        uuid.UUID          `json:"lago_id,omitempty"`
	Number        string             `json:"number,omitempty"`
	Status        OrderStatus        `json:"status,omitempty"`
	OrderType     QuoteOrderType     `json:"order_type,omitempty"`
	ExecutionMode OrderExecutionMode `json:"execution_mode,omitempty"`
	Currency      Currency           `json:"currency,omitempty"`

	ExecutedAt      *time.Time            `json:"executed_at,omitempty"`
	ExecutionRecord *OrderExecutionRecord `json:"execution_record,omitempty"`

	LagoOrganizationID uuid.UUID `json:"lago_organization_id,omitempty"`
	LagoCustomerID     uuid.UUID `json:"lago_customer_id,omitempty"`
	LagoOrderFormID    uuid.UUID `json:"lago_order_form_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// The billing items of the quote version the order comes from, as they stood when it
	// was approved. Omitted from the webhook payloads, being a heavy blob.
	BillingSnapshot *QuoteBillingItems `json:"billing_snapshot,omitempty"`
}

type OrderListInput struct {
	PerPage *int `url:"per_page,omitempty"`
	Page    *int `url:"page,omitempty"`

	Status          []OrderStatus        `url:"status[],omitempty"`
	OrderType       []QuoteOrderType     `url:"order_type[],omitempty"`
	ExecutionMode   []OrderExecutionMode `url:"execution_mode[],omitempty"`
	CustomerIDs     []uuid.UUID          `url:"customer_id[],omitempty"`
	Number          []string             `url:"number[],omitempty"`
	OrderFormNumber []string             `url:"order_form_number[],omitempty"`
	QuoteNumber     []string             `url:"quote_number[],omitempty"`
	OwnerIDs        []uuid.UUID          `url:"owner_id[],omitempty"`

	SearchTerm string `url:"search_term,omitempty"`

	ExecutedAtFrom string `url:"executed_at_from,omitempty"`
	ExecutedAtTo   string `url:"executed_at_to,omitempty"`
}

type OrderExecuteInput struct {
	ExecutionMode OrderExecutionMode `json:"execution_mode,omitempty"`
}

type OrderExecuteParams struct {
	Order *OrderExecuteInput `json:"order"`
}

type OrderResult struct {
	Order  *Order   `json:"order,omitempty"`
	Orders []Order  `json:"orders,omitempty"`
	Meta   Metadata `json:"meta,omitempty"`
}

func (c *Client) Order() *OrderRequest {
	return &OrderRequest{
		client: c,
	}
}

func (or *OrderRequest) Get(ctx context.Context, orderID string) (*Order, *Error) {
	subPath := fmt.Sprintf("%s/%s", OrdersEndpoint, orderID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &OrderResult{},
	}

	result, err := or.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	orderResult, ok := result.(*OrderResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return orderResult.Order, nil
}

func (or *OrderRequest) GetList(ctx context.Context, orderListInput *OrderListInput) (*OrderResult, *Error) {
	urlValues, err := query.Values(orderListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      OrdersEndpoint,
		UrlValues: urlValues,
		Result:    &OrderResult{},
	}

	result, clientErr := or.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	orderResult, ok := result.(*OrderResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return orderResult, nil
}

// Execute carries out an order on demand, without waiting for its execute_at schedule.
// It runs synchronously and is idempotent: an order already executed is returned untouched.
func (or *OrderRequest) Execute(ctx context.Context, orderID string, executeInput *OrderExecuteInput) (*Order, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", OrdersEndpoint, orderID, "execute")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &OrderResult{},
	}

	if executeInput != nil {
		clientRequest.Body = &OrderExecuteParams{Order: executeInput}
	}

	result, err := or.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	orderResult, ok := result.(*OrderResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return orderResult.Order, nil
}
