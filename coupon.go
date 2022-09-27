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
	CouponTypeFixedAmount    CouponCalculationType = "fixed_amount"
	CouponTypePercentage     CouponCalculationType = "percentage"
)

type CouponFrequency string

const (
	CouponFrequencyOnce      CouponFrequency = "once"
	CouponFrequencyRecurring CouponFrequency = "recurring"
)

type CouponRequest struct {
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

type CouponInput struct {
	Name               string                   `json:"name,omitempty"`
	Code               string                   `json:"code,omitempty"`
	AmountCents        int                      `json:"amount_cents,omitempty"`
	AmountCurrency     Currency                 `json:"amount_currency,omitempty"`
	Expiration         CouponExpiration         `json:"expiration,omitempty"`
	ExpirationDate     string                   `json:"expiration_date,omitempty"`
	PercentageRate     float32                  `json:"percentage_rate,omitempty"`
	CouponType         CouponCalculationType    `json:"coupon_type,omitempty"`
	Frequency          CouponFrequency          `json:"frequency,omitempty"`
	FrequencyDuration  int                      `json:"frequency_duration,omitempty"`
}

type CouponListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type Coupon struct {
	LagoID             uuid.UUID             `json:"lago_id,omitempty"`
	Name               string                `json:"name,omitempty"`
	Code               string                `json:"code,omitempty"`
	AmountCents        int                   `json:"amount_cents,omitempty"`
	AmountCurrency     Currency              `json:"amount_currency,omitempty"`
	Expiration         CouponExpiration      `json:"expiration,omitempty"`
	ExpirationDate     string                `json:"expiration_date,omitempty"`
	PercentageRate     float32               `json:"percentage_rate,omitempty"`
	CouponType         CouponCalculationType `json:"coupon_type,omitempty"`
	Frequency          CouponFrequency       `json:"frequency,omitempty"`
	FrequencyDuration  int                   `json:"frequency_duration,omitempty"`
	CreatedAt          time.Time             `json:"created_at,omitempty"`
}

type AppliedCouponResult struct {
	AppliedCoupon *AppliedCoupon `json:"applied_coupon,omitempty"`
}

type ApplyCouponParams struct {
	AppliedCoupon *ApplyCouponInput `json:"applied_coupon"`
}

type ApplyCouponInput struct {
	ExternalCustomerID string           `json:"external_customer_id,omitempty"`
	CouponCode         string           `json:"coupon_code,omitempty"`
	AmountCents        int              `json:"amount_cents,omitempty"`
	AmountCurrency     Currency         `json:"amount_currency,omitempty"`
	PercentageRate     float32          `json:"percentage_rate,omitempty"`
	Frequency          CouponFrequency  `json:"frequency,omitempty"`
	FrequencyDuration  int              `json:"frequency_duration,omitempty"`
}

type AppliedCoupon struct {
	LagoID             uuid.UUID `json:"lago_id,omitempty"`
	LagoCouponID       uuid.UUID `json:"lago_coupon_id,omitempty"`
	ExternalCustomerID string    `json:"external_customer_id,omitempty"`
	LagoCustomerID     uuid.UUID `json:"lago_customer_id,omitempty"`

	CouponCode     string   `json:"coupon_code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	ExpirationDate string    `json:"expiration_date,omitempty"`
	TerminatedAt   time.Time `json:"terminated_at,omitempty"`

	PercentageRate     float32          `json:"percentage_rate,omitempty"`
	Frequency          CouponFrequency  `json:"frequency,omitempty"`
	FrequencyDuration  int              `json:"frequency_duration,omitempty"`
}

func (c *Client) Coupon() *CouponRequest {
	return &CouponRequest{
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

	couponResult := result.(*CouponResult)

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
		Path:        "plans",
		QueryParams: queryParams,
		Result:      &CouponResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	couponResult := result.(*CouponResult)

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

	couponResult := result.(*CouponResult)

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

	couponResult := result.(*CouponResult)

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

	couponResult := result.(*CouponResult)

	return couponResult.Coupon, nil
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

	appliedCouponResult := result.(*AppliedCouponResult)

	return appliedCouponResult.AppliedCoupon, nil
}
