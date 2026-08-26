package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

const QuotesEndpoint string = "quotes"

type QuoteOrderType string

const (
	QuoteOrderTypeSubscriptionCreation  QuoteOrderType = "subscription_creation"
	QuoteOrderTypeSubscriptionAmendment QuoteOrderType = "subscription_amendment"
	QuoteOrderTypeOneOff                QuoteOrderType = "one_off"
)

type QuoteRequest struct {
	client *Client
}

type QuoteOwner struct {
	LagoID uuid.UUID `json:"lago_id,omitempty"`
	Email  string    `json:"email,omitempty"`
}

type Quote struct {
	LagoID             uuid.UUID      `json:"lago_id,omitempty"`
	Number             string         `json:"number,omitempty"`
	OrderType          QuoteOrderType `json:"order_type,omitempty"`
	LagoCustomerID     uuid.UUID      `json:"lago_customer_id,omitempty"`
	LagoSubscriptionID *uuid.UUID     `json:"lago_subscription_id,omitempty"`
	LagoOrganizationID uuid.UUID      `json:"lago_organization_id,omitempty"`
	CreatedAt          time.Time      `json:"created_at,omitempty"`
	UpdatedAt          time.Time      `json:"updated_at,omitempty"`

	CurrentVersion *QuoteVersion `json:"current_version,omitempty"`

	// Only returned when a single quote is retrieved.
	Owners []QuoteOwner `json:"owners,omitempty"`
}

// QuoteWithVersion is the payload of the quote lifecycle webhooks. It names the version the event
// happened to, which is not necessarily the current version of the quote: a quote.voided triggered
// by a clone names the version that was voided, while the quote already carries its replacement.
// Unlike Quote, it carries no owners.
type QuoteWithVersion struct {
	LagoID             uuid.UUID      `json:"lago_id,omitempty"`
	Number             string         `json:"number,omitempty"`
	OrderType          QuoteOrderType `json:"order_type,omitempty"`
	LagoCustomerID     uuid.UUID      `json:"lago_customer_id,omitempty"`
	LagoSubscriptionID *uuid.UUID     `json:"lago_subscription_id,omitempty"`
	LagoOrganizationID uuid.UUID      `json:"lago_organization_id,omitempty"`
	CreatedAt          time.Time      `json:"created_at,omitempty"`
	UpdatedAt          time.Time      `json:"updated_at,omitempty"`

	// The Content and BillingItems of the version are omitted; retrieve the version to get them.
	Version *QuoteVersion `json:"version,omitempty"`
}

type QuoteListInput struct {
	PerPage *int `url:"per_page,omitempty"`
	Page    *int `url:"page,omitempty"`

	Status             []QuoteVersionStatus `url:"status[],omitempty"`
	OrderType          []QuoteOrderType     `url:"order_type[],omitempty"`
	Number             []string             `url:"number[],omitempty"`
	OwnerIDs           []uuid.UUID          `url:"owner_id[],omitempty"`
	ExternalCustomerID []string             `url:"external_customer_id[],omitempty"`

	FromDate string `url:"from_date,omitempty"`
	ToDate   string `url:"to_date,omitempty"`
}

type QuoteVersionListInput struct {
	PerPage *int `url:"per_page,omitempty"`
	Page    *int `url:"page,omitempty"`
}

type QuoteResult struct {
	Quote  *Quote   `json:"quote,omitempty"`
	Quotes []Quote  `json:"quotes,omitempty"`
	Meta   Metadata `json:"meta,omitempty"`
}

func (c *Client) Quote() *QuoteRequest {
	return &QuoteRequest{
		client: c,
	}
}

func (qr *QuoteRequest) Get(ctx context.Context, quoteID string) (*Quote, *Error) {
	subPath := fmt.Sprintf("%s/%s", QuotesEndpoint, quoteID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &QuoteResult{},
	}

	result, err := qr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	quoteResult, ok := result.(*QuoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return quoteResult.Quote, nil
}

func (qr *QuoteRequest) GetList(ctx context.Context, quoteListInput *QuoteListInput) (*QuoteResult, *Error) {
	urlValues, err := query.Values(quoteListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:      QuotesEndpoint,
		UrlValues: urlValues,
		Result:    &QuoteResult{},
	}

	result, clientErr := qr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	quoteResult, ok := result.(*QuoteResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return quoteResult, nil
}

// GetVersionList lists the versions of a quote, from the most recent to the oldest.
// The Content and BillingItems of each version are omitted; retrieve a version to get them.
func (qr *QuoteRequest) GetVersionList(ctx context.Context, quoteID string, quoteVersionListInput *QuoteVersionListInput) (*QuoteVersionResult, *Error) {
	urlValues, err := query.Values(quoteVersionListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	subPath := fmt.Sprintf("%s/%s/%s", QuotesEndpoint, quoteID, "versions")
	clientRequest := &ClientRequest{
		Path:      subPath,
		UrlValues: urlValues,
		Result:    &QuoteVersionResult{},
	}

	result, clientErr := qr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	quoteVersionResult, ok := result.(*QuoteVersionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return quoteVersionResult, nil
}
