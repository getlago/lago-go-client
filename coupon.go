package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CouponExpiration string

const (
	CouponExpirationTimeLimit    CouponExpiration = "time_limit"
	CouponExpirationNoExpiration CouponExpiration = "no_expiration"
)

type CouponCalculationType string

const (
	CouponTypeFixedAmount CouponCalculationType = "fixed_amount"
	CouponTypePercentage  CouponCalculationType = "percentage"
)

type CouponFrequency string

const (
	CouponFrequencyOnce      CouponFrequency = "once"
	CouponFrequencyRecurring CouponFrequency = "recurring"
)

type AppliedCouponStatus string

const (
	AppliedCouponStatusActive     AppliedCouponStatus = "active"
	AppliedCouponStatusTerminated AppliedCouponStatus = "terminated"
)

type CouponRequest struct {
	client *Client
}

type AppliedCouponRequest struct {
	client *Client
}

type CouponResult struct {
	Coupon  *Coupon  `json:"coupon,omitempty"`
	Coupons []Coupon `json:"coupons,omitempty"`
	Meta    Metadata `json:"meta,omitempty"`
}

type CouponParams struct {
	Coupon *CouponInput `json:"coupon"`
}

type LimitationInput struct {
	PlanCodes           []string `json:"plan_codes,omitempty"`
	BillableMetricCodes []string `json:"billable_metric_codes,omitempty"`
}

type CouponInput struct {
	Name              string                `json:"name,omitempty"`
	Code              string                `json:"code,omitempty"`
	AmountCents       int                   `json:"amount_cents,omitempty"`
	AmountCurrency    Currency              `json:"amount_currency,omitempty"`
	Expiration        CouponExpiration      `json:"expiration,omitempty"`
	ExpirationAt      *time.Time            `json:"expiration_at,omitempty"`
	PercentageRate    float32               `json:"percentage_rate,omitempty"`
	CouponType        CouponCalculationType `json:"coupon_type,omitempty"`
	Frequency         CouponFrequency       `json:"frequency,omitempty"`
	Reusable          bool                  `json:"reusable,omitempty"`
	FrequencyDuration int                   `json:"frequency_duration,omitempty"`
	AppliesTo         LimitationInput       `json:"applies_to,omitempty"`
}

type CouponListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type Coupon struct {
	LagoID                 uuid.UUID             `json:"lago_id,omitempty"`
	Name                   string                `json:"name,omitempty"`
	Code                   string                `json:"code,omitempty"`
	AmountCents            int                   `json:"amount_cents,omitempty"`
	AmountCurrency         Currency              `json:"amount_currency,omitempty"`
	Expiration             CouponExpiration      `json:"expiration,omitempty"`
	ExpirationDate         string                `json:"expiration_date,omitempty"`
	PercentageRate         float32               `json:"percentage_rate,omitempty"`
	CouponType             CouponCalculationType `json:"coupon_type,omitempty"`
	Frequency              CouponFrequency       `json:"frequency,omitempty"`
	Reusable               bool                  `json:"reusable,omitempty"`
	LimitedPlans           bool                  `json:"limited_plans,omitempty"`
	PlanCodes              []string              `json:"plan_codes,omitempty"`
	LimitedBillableMetrics bool                  `json:"limited_billable_metrics,omitempty"`
	BillableMetricCodes    []string              `json:"billable_metric_codes,omitempty"`
	FrequencyDuration      int                   `json:"frequency_duration,omitempty"`
	CreatedAt              time.Time             `json:"created_at,omitempty"`
}

type AppliedCouponResult struct {
	AppliedCoupon  *AppliedCoupon  `json:"applied_coupon,omitempty"`
	AppliedCoupons []AppliedCoupon `json:"applied_coupons,omitempty"`
	Meta           Metadata        `json:"meta,omitempty"`
}

type AppliedCouponListInput struct {
	PerPage            int                 `json:"per_page,omitempty,string"`
	Page               int                 `json:"page,omitempty,string"`
	Status             AppliedCouponStatus `json:"status,omitempty,string"`
	ExternalCustomerID string              `json:"external_customer_id,omitempty,string"`
}

type ApplyCouponParams struct {
	AppliedCoupon *ApplyCouponInput `json:"applied_coupon"`
}

type ApplyCouponInput struct {
	ExternalCustomerID string          `json:"external_customer_id,omitempty"`
	CouponCode         string          `json:"coupon_code,omitempty"`
	AmountCents        int             `json:"amount_cents,omitempty"`
	AmountCurrency     Currency        `json:"amount_currency,omitempty"`
	PercentageRate     float32         `json:"percentage_rate,omitempty"`
	Frequency          CouponFrequency `json:"frequency,omitempty"`
	FrequencyDuration  int             `json:"frequency_duration,omitempty"`
}

type AppliedCoupon struct {
	LagoID             uuid.UUID           `json:"lago_id,omitempty"`
	LagoCouponID       uuid.UUID           `json:"lago_coupon_id,omitempty"`
	ExternalCustomerID string              `json:"external_customer_id,omitempty"`
	LagoCustomerID     uuid.UUID           `json:"lago_customer_id,omitempty"`
	Status             AppliedCouponStatus `json:"status,omitempty"`

	CouponCode     string   `json:"coupon_code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	ExpirationDate string    `json:"expiration_date,omitempty"`
	TerminatedAt   time.Time `json:"terminated_at,omitempty"`

	PercentageRate    float32         `json:"percentage_rate,omitempty"`
	Frequency         CouponFrequency `json:"frequency,omitempty"`
	FrequencyDuration int             `json:"frequency_duration,omitempty"`

	AmountCentsRemaining       int `json:"amount_cents_remaining,omitempty"`
	FrequencyDurationRemaining int `json:"frequency_duration_remaining,omitempty"`

	Credits []InvoiceCredit `json:"credits,omitempty"`
}

func (c *Client) Coupon() *CouponRequest {
	return &CouponRequest{
		client: c,
	}
}

func (c *Client) AppliedCoupon() *AppliedCouponRequest {
	return &AppliedCouponRequest{
		client: c,
	}
}

func (cr *CouponRequest) Get(ctx context.Context, couponCode string) (*Coupon, *Error) {
	subPath := fmt.Sprintf("%s/%s", "coupons", couponCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CouponResult{},
	}

	result, err := cr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult, ok := result.(*CouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) GetList(ctx context.Context, couponListInput *CouponListInput) (*CouponResult, *Error) {
	jsonQueryParams, err := json.Marshal(couponListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "coupons",
		QueryParams: queryParams,
		Result:      &CouponResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	couponResult, ok := result.(*CouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return couponResult, nil
}

func (cr *CouponRequest) Create(ctx context.Context, couponInput *CouponInput) (*Coupon, *Error) {
	couponParams := &CouponParams{
		Coupon: couponInput,
	}

	clientRequest := &ClientRequest{
		Path:   "coupons",
		Result: &CouponResult{},
		Body:   couponParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult, ok := result.(*CouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) Update(ctx context.Context, couponInput *CouponInput) (*Coupon, *Error) {
	subPath := fmt.Sprintf("%s/%s", "coupons", couponInput.Code)
	couponParams := &CouponParams{
		Coupon: couponInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CouponResult{},
		Body:   couponParams,
	}

	result, err := cr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult, ok := result.(*CouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) Delete(ctx context.Context, couponCode string) (*Coupon, *Error) {
	subPath := fmt.Sprintf("%s/%s", "coupons", couponCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CouponResult{},
	}

	result, err := cr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult, ok := result.(*CouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return couponResult.Coupon, nil
}

func (cr *AppliedCouponRequest) GetList(ctx context.Context, appliedCouponListInput *AppliedCouponListInput) (*AppliedCouponResult, *Error) {
	jsonQueryParams, err := json.Marshal(appliedCouponListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "applied_coupons",
		QueryParams: queryParams,
		Result:      &AppliedCouponResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	appliedCouponResult, ok := result.(*AppliedCouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedCouponResult, nil
}

func (cr *CouponRequest) ApplyToCustomer(ctx context.Context, applyCouponInput *ApplyCouponInput) (*AppliedCoupon, *Error) {
	applyCouponParams := &ApplyCouponParams{
		AppliedCoupon: applyCouponInput,
	}

	clientRequest := &ClientRequest{
		Path:   "applied_coupons",
		Result: &AppliedCouponResult{},
		Body:   applyCouponParams,
	}

	result, err := cr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedCouponResult, ok := result.(*AppliedCouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedCouponResult.AppliedCoupon, nil
}

func (acr *AppliedCouponRequest) AppliedCouponDelete(ctx context.Context, externalCustomerID string, appliedCouponID string) (*AppliedCoupon, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "customers", externalCustomerID, "applied_coupons", appliedCouponID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AppliedCouponResult{},
	}

	result, err := acr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	appliedCouponResult, ok := result.(*AppliedCouponResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return appliedCouponResult.AppliedCoupon, nil
}
