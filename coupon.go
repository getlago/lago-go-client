package lago

import (
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
	Name               string           `json:"name,omitempty"`
	Code               string           `json:"code,omitempty"`
	AmountCents        int              `json:"amount_cents,omitempty"`
	AmountCurrency     Currency         `json:"amount_currency,omitempty"`
	Expiration         CouponExpiration `json:"expiration,omitempty"`
	ExpirationDuration int              `json:"expiration_duration,omitempty"`
}

type CouponListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string`
}

type Coupon struct {
	LagoID             uuid.UUID        `json:"lago_id,omitempty"`
	Name               string           `json:"name,omitempty"`
	Code               string           `json:"code,omitempty"`
	AmountCents        int              `json:"amount_cents,omitempty"`
	AmountCurrency     Currency         `json:"amount_currency,omitempty"`
	Expiration         CouponExpiration `json:"expiration,omitempty"`
	ExpirationDuration int              `json:"expiration_duration,omitempty"`
	CreatedAt          time.Time        `json:"created_at,omitempty"`
}

type AppliedCouponResult struct {
	AppliedCoupon *AppliedCoupon `json:"applied_coupon,omitempty"`
}

type ApplyCouponParams struct {
	AppliedCoupon *ApplyCouponInput `json:"applied_coupon"`
}

type ApplyCouponInput struct {
	CustomerID     string   `json:"customer_id,omitempty"`
	CouponCode     string   `json:"coupon_code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`
}

type AppliedCoupon struct {
	LagoID         uuid.UUID `json:"lago_id,omitempty"`
	LagoCouponID   uuid.UUID `json:"lago_coupon_id,omitempty"`
	CustomerID     string    `json:"customer_id,omitempty"`
	LagoCustomerID uuid.UUID `json:"lago_customer_id,omitempty"`

	CouponCode     string   `json:"coupon_code,omitempty"`
	AmountCents    int      `json:"amount_cents,omitempty"`
	AmountCurrency Currency `json:"amount_currency,omitempty"`

	ExpirationDate time.Time `json:"expiration_date,omitempty`
	TerminatedAt   time.Time `json:"terminated_at,omitempty"`
}

func (c *Client) Coupon() *CouponRequest {
	return &CouponRequest{
		client: c,
	}
}

func (cr *CouponRequest) Get(couponCode string) (*Coupon, *Error) {
	subPath := fmt.Sprintf("%s/%s", "coupons", couponCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CouponResult{},
	}

	result, err := cr.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult := result.(*CouponResult)

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) GetList(couponListInput *CouponListInput) (*CouponResult, *Error) {
	jsonQueryParams, err := json.Marshal(couponListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	json.Unmarshal(jsonQueryParams, &queryParams)

	clientRequest := &ClientRequest{
		Path:        "plans",
		QueryParams: queryParams,
		Result:      &CouponResult{},
	}

	result, clientErr := cr.client.Get(clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	couponResult := result.(*CouponResult)

	return couponResult, nil
}

func (cr *CouponRequest) Create(couponInput *CouponInput) (*Coupon, *Error) {
	couponParams := &CouponParams{
		Coupon: couponInput,
	}

	clientRequest := &ClientRequest{
		Path:   "coupons",
		Result: &CouponResult{},
		Body:   couponParams,
	}

	result, err := cr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult := result.(*CouponResult)

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) Update(couponInput *CouponInput) (*Coupon, *Error) {
	subPath := fmt.Sprintf("%s/%s", "coupons", couponInput.Code)
	couponParams := &CouponParams{
		Coupon: couponInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CouponResult{},
		Body:   couponParams,
	}

	result, err := cr.client.Put(clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult := result.(*CouponResult)

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) Delete(couponCode string) (*Coupon, *Error) {
	subPath := fmt.Sprintf("%s/%s", "coupons", couponCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CouponResult{},
	}

	result, err := cr.client.Delete(clientRequest)
	if err != nil {
		return nil, err
	}

	couponResult := result.(*CouponResult)

	return couponResult.Coupon, nil
}

func (cr *CouponRequest) ApplyToCustomer(applyCouponInput *ApplyCouponInput) (*AppliedCoupon, *Error) {
	applyCouponParams := &ApplyCouponParams{
		AppliedCoupon: applyCouponInput,
	}

	clientRequest := &ClientRequest{
		Path:   "applied_coupons",
		Result: &AppliedCouponResult{},
		Body:   applyCouponParams,
	}

	result, err := cr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	appliedCouponResult := result.(*AppliedCouponResult)

	return appliedCouponResult.AppliedCoupon, nil
}
