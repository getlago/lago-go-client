package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

const OrderFormsEndpoint string = "order_forms"

type OrderFormStatus string

const (
	OrderFormStatusGenerated OrderFormStatus = "generated"
	OrderFormStatusSigned    OrderFormStatus = "signed"
	OrderFormStatusExpired   OrderFormStatus = "expired"
	OrderFormStatusVoided    OrderFormStatus = "voided"
)

type OrderFormVoidReason string

const (
	OrderFormVoidReasonManual  OrderFormVoidReason = "manual"
	OrderFormVoidReasonExpired OrderFormVoidReason = "expired"
	OrderFormVoidReasonInvalid OrderFormVoidReason = "invalid"
)

type OrderExecutionMode string

const (
	OrderExecutionModeExecuteInLago OrderExecutionMode = "execute_in_lago"
	OrderExecutionModeOrderOnly     OrderExecutionMode = "order_only"
)

type OrderFormRequest struct {
	client *Client
}

type OrderForm struct {
	LagoID     uuid.UUID           `json:"lago_id,omitempty"`
	Number     string              `json:"number,omitempty"`
	Status     OrderFormStatus     `json:"status,omitempty"`
	VoidReason OrderFormVoidReason `json:"void_reason,omitempty"`

	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	SignedAt          *time.Time `json:"signed_at,omitempty"`
	VoidedAt          *time.Time `json:"voided_at,omitempty"`
	SignedDocumentUrl string     `json:"signed_document_url,omitempty"`

	LagoOrganizationID uuid.UUID `json:"lago_organization_id,omitempty"`
	LagoCustomerID     uuid.UUID `json:"lago_customer_id,omitempty"`
	LagoQuoteID        uuid.UUID `json:"lago_quote_id,omitempty"`
	LagoQuoteVersionID uuid.UUID `json:"lago_quote_version_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type OrderFormListInput struct {
	PerPage *int `url:"per_page,omitempty"`
	Page    *int `url:"page,omitempty"`

	Status      []OrderFormStatus `url:"status[],omitempty"`
	CustomerIDs []uuid.UUID       `url:"customer_id[],omitempty"`
	Number      []string          `url:"number[],omitempty"`
	QuoteNumber []string          `url:"quote_number[],omitempty"`
	OwnerIDs    []uuid.UUID       `url:"owner_id[],omitempty"`

	SearchTerm string `url:"search_term,omitempty"`

	CreatedAtFrom string `url:"created_at_from,omitempty"`
	CreatedAtTo   string `url:"created_at_to,omitempty"`
	ExpiresAtFrom string `url:"expires_at_from,omitempty"`
	ExpiresAtTo   string `url:"expires_at_to,omitempty"`
}

// The signed document is sent as a base64 data URI (data:<content_type>;base64,<data>),
// accepting application/pdf, image/jpeg and image/png up to 10 MB.
type OrderFormMarkAsSignedInput struct {
	SignedDocument string             `json:"signed_document,omitempty"`
	ExecutionMode  OrderExecutionMode `json:"execution_mode,omitempty"`
	ExecuteAt      *time.Time         `json:"execute_at,omitempty"`
}

type OrderFormMarkAsSignedParams struct {
	OrderForm *OrderFormMarkAsSignedInput `json:"order_form"`
}

type OrderFormResult struct {
	OrderForm  *OrderForm  `json:"order_form,omitempty"`
	OrderForms []OrderForm `json:"order_forms,omitempty"`
	Meta       Metadata    `json:"meta,omitempty"`
}

func (c *Client) OrderForm() *OrderFormRequest {
	return &OrderFormRequest{
		client: c,
	}
}

func (ofr *OrderFormRequest) Get(ctx context.Context, orderFormID string) (*OrderForm, *Error) {
	subPath := fmt.Sprintf("%s/%s", OrderFormsEndpoint, orderFormID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &OrderFormResult{},
	}

	result, err := ofr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	orderFormResult, ok := result.(*OrderFormResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return orderFormResult.OrderForm, nil
}

func (ofr *OrderFormRequest) GetList(ctx context.Context, orderFormListInput *OrderFormListInput) (*OrderFormResult, *Error) {
	urlValues, err := query.Values(orderFormListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      OrderFormsEndpoint,
		UrlValues: urlValues,
		Result:    &OrderFormResult{},
	}

	result, clientErr := ofr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	orderFormResult, ok := result.(*OrderFormResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return orderFormResult, nil
}

// MarkAsSigned records the customer's signature and creates the order carrying the deal out.
func (ofr *OrderFormRequest) MarkAsSigned(ctx context.Context, orderFormID string, markAsSignedInput *OrderFormMarkAsSignedInput) (*OrderForm, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", OrderFormsEndpoint, orderFormID, "mark_as_signed")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &OrderFormResult{},
	}

	if markAsSignedInput != nil {
		clientRequest.Body = &OrderFormMarkAsSignedParams{OrderForm: markAsSignedInput}
	}

	return ofr.postAndUnwrap(ctx, clientRequest)
}

// Void voids a generated order form, cascading to the quote version it was generated from.
func (ofr *OrderFormRequest) Void(ctx context.Context, orderFormID string) (*OrderForm, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", OrderFormsEndpoint, orderFormID, "void")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &OrderFormResult{},
	}

	return ofr.postAndUnwrap(ctx, clientRequest)
}

func (ofr *OrderFormRequest) postAndUnwrap(ctx context.Context, clientRequest *ClientRequest) (*OrderForm, *Error) {
	result, err := ofr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	orderFormResult, ok := result.(*OrderFormResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return orderFormResult.OrderForm, nil
}
