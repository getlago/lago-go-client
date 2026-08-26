package lago

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const QuoteVersionsEndpoint string = "quote_versions"

type QuoteVersionStatus string

const (
	QuoteVersionStatusDraft    QuoteVersionStatus = "draft"
	QuoteVersionStatusApproved QuoteVersionStatus = "approved"
	QuoteVersionStatusVoided   QuoteVersionStatus = "voided"
)

type QuoteVersionVoidReason string

const (
	QuoteVersionVoidReasonManual           QuoteVersionVoidReason = "manual"
	QuoteVersionVoidReasonSuperseded       QuoteVersionVoidReason = "superseded"
	QuoteVersionVoidReasonCascadeOfExpired QuoteVersionVoidReason = "cascade_of_expired"
	QuoteVersionVoidReasonCascadeOfVoided  QuoteVersionVoidReason = "cascade_of_voided"
)

type QuoteBillingItemType string

const (
	QuoteBillingItemTypeAddOn        QuoteBillingItemType = "add_on"
	QuoteBillingItemTypePlan         QuoteBillingItemType = "plan"
	QuoteBillingItemTypeCoupon       QuoteBillingItemType = "coupon"
	QuoteBillingItemTypeWalletCredit QuoteBillingItemType = "wallet_credit"
)

type QuoteVersionRequest struct {
	client *Client
}

// The billing items are stored and returned as they are authored in the Lago user
// interface, so their keys are camelCased, unlike the rest of the API. Every Payload
// pins the catalog record the entry was built from, and may carry catalog properties
// beyond the documented ones, so it is left untyped.

type QuoteBillingItemAddOn struct {
	LagoID    uuid.UUID              `json:"id,omitempty"`
	LocalID   string                 `json:"localId,omitempty"`
	Type      QuoteBillingItemType   `json:"type,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

type QuoteBillingItemPlan struct {
	LagoID    uuid.UUID              `json:"id,omitempty"`
	LocalID   string                 `json:"localId,omitempty"`
	Type      QuoteBillingItemType   `json:"type,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

type QuoteBillingItemCoupon struct {
	LagoID    uuid.UUID              `json:"id,omitempty"`
	LocalID   string                 `json:"localId,omitempty"`
	Type      QuoteBillingItemType   `json:"type,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

// A wallet credit is not built from an existing record, so it carries no id and no overrides.
type QuoteBillingItemWalletCredit struct {
	LocalID string                 `json:"localId,omitempty"`
	Type    QuoteBillingItemType   `json:"type,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type QuoteBillingItems struct {
	AddOns        []QuoteBillingItemAddOn        `json:"addOns,omitempty"`
	Plans         []QuoteBillingItemPlan         `json:"plans,omitempty"`
	Coupons       []QuoteBillingItemCoupon       `json:"coupons,omitempty"`
	WalletCredits []QuoteBillingItemWalletCredit `json:"walletCredits,omitempty"`
}

type QuoteVersion struct {
	LagoID             uuid.UUID              `json:"lago_id,omitempty"`
	LagoQuoteID        uuid.UUID              `json:"lago_quote_id,omitempty"`
	LagoOrganizationID uuid.UUID              `json:"lago_organization_id,omitempty"`
	Version            int                    `json:"version,omitempty"`
	Status             QuoteVersionStatus     `json:"status,omitempty"`
	Currency           Currency               `json:"currency,omitempty"`
	BillingEntityCode  string                 `json:"billing_entity_code,omitempty"`
	VoidReason         QuoteVersionVoidReason `json:"void_reason,omitempty"`
	ApprovedAt         *time.Time             `json:"approved_at,omitempty"`
	VoidedAt           *time.Time             `json:"voided_at,omitempty"`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`

	// Only returned when a single version is retrieved.
	Content      string             `json:"content,omitempty"`
	BillingItems *QuoteBillingItems `json:"billing_items,omitempty"`
}

type QuoteVersionApproveInput struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type QuoteVersionResult struct {
	QuoteVersion  *QuoteVersion  `json:"quote_version,omitempty"`
	QuoteVersions []QuoteVersion `json:"quote_versions,omitempty"`
	Meta          Metadata       `json:"meta,omitempty"`
}

func (c *Client) QuoteVersion() *QuoteVersionRequest {
	return &QuoteVersionRequest{
		client: c,
	}
}

func (qvr *QuoteVersionRequest) Get(ctx context.Context, quoteVersionID string) (*QuoteVersion, *Error) {
	subPath := fmt.Sprintf("%s/%s", QuoteVersionsEndpoint, quoteVersionID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &QuoteVersionResult{},
	}

	result, err := qvr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	quoteVersionResult, ok := result.(*QuoteVersionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return quoteVersionResult.QuoteVersion, nil
}

// Approve approves a draft version and generates the order form to send for signature.
func (qvr *QuoteVersionRequest) Approve(ctx context.Context, quoteVersionID string, approveInput *QuoteVersionApproveInput) (*QuoteVersion, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", QuoteVersionsEndpoint, quoteVersionID, "approve")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &QuoteVersionResult{},
	}

	if approveInput != nil {
		clientRequest.Body = approveInput
	}

	return qvr.postAndUnwrap(ctx, clientRequest)
}

// Void voids a draft version, which makes it definitive.
func (qvr *QuoteVersionRequest) Void(ctx context.Context, quoteVersionID string) (*QuoteVersion, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", QuoteVersionsEndpoint, quoteVersionID, "void")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &QuoteVersionResult{},
	}

	return qvr.postAndUnwrap(ctx, clientRequest)
}

// Clone copies a version into a new draft version of the same quote.
func (qvr *QuoteVersionRequest) Clone(ctx context.Context, quoteVersionID string) (*QuoteVersion, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", QuoteVersionsEndpoint, quoteVersionID, "clone")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &QuoteVersionResult{},
	}

	return qvr.postAndUnwrap(ctx, clientRequest)
}

func (qvr *QuoteVersionRequest) postAndUnwrap(ctx context.Context, clientRequest *ClientRequest) (*QuoteVersion, *Error) {
	result, err := qvr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	quoteVersionResult, ok := result.(*QuoteVersionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return quoteVersionResult.QuoteVersion, nil
}
