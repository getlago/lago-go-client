package lago

import (
	"context"
	"encoding/json"
	"fmt"
)

type CustomerPaymentMethodResult struct {
	PaymentMethods []PaymentMethod `json:"payment_methods,omitempty"`
	Meta           Metadata        `json:"meta,omitempty"`
}

type CustomerPaymentMethodParams struct {
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
}

type CustomerPaymentMethodListInput struct {
	PerPage *int `json:"per_page,omitempty,string"`
	Page    *int `json:"page,omitempty,string"`
}

func (cr *CustomerRequest) GetPaymentMethodList(ctx context.Context, externalCustomerID string, listInput *CustomerPaymentMethodListInput) (*CustomerPaymentMethodResult, *Error) {
	jsonQueryParams, err := json.Marshal(listInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	subPath := fmt.Sprintf("customers/%s/payment_methods", externalCustomerID)
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &CustomerPaymentMethodResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	paymentMethodResult, ok := result.(*CustomerPaymentMethodResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentMethodResult, nil
}

func (cr *CustomerRequest) DestroyPaymentMethod(ctx context.Context, externalCustomerID string, paymentMethodID string) (*PaymentMethod, *Error) {
	subPath := fmt.Sprintf("customers/%s/payment_methods/%s", externalCustomerID, paymentMethodID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerPaymentMethodParams{},
	}

	result, err := cr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentMethodResult, ok := result.(*CustomerPaymentMethodParams)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentMethodResult.PaymentMethod, nil
}

func (cr *CustomerRequest) SetPaymentMethodAsDefault(ctx context.Context, externalCustomerID string, paymentMethodID string) (*PaymentMethod, *Error) {
	subPath := fmt.Sprintf("customers/%s/payment_methods/%s/set_as_default", externalCustomerID, paymentMethodID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &CustomerPaymentMethodParams{},
	}

	result, err := cr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	paymentMethodResult, ok := result.(*CustomerPaymentMethodParams)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return paymentMethodResult.PaymentMethod, nil
}
