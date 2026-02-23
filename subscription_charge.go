package lago

import (
	"context"
	"encoding/json"
	"fmt"
)

// Charges

func (sr *SubscriptionRequest) GetCharge(ctx context.Context, externalID string, chargeCode string) (*Charge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeResult{},
	}

	result, err := sr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult.Charge, nil
}

func (sr *SubscriptionRequest) GetChargeList(ctx context.Context, externalID string, chargeListInput *ChargeListInput) (*ChargeResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", externalID, "charges")

	jsonQueryParams, err := json.Marshal(chargeListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &ChargeResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult, nil
}

func (sr *SubscriptionRequest) UpdateCharge(ctx context.Context, externalID string, chargeCode string, chargeInput *ChargeInput) (*Charge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode)

	chargeParams := &ChargeParams{
		Charge: chargeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeResult{},
		Body:   chargeParams,
	}

	result, err := sr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult.Charge, nil
}

// Fixed Charges

func (sr *SubscriptionRequest) GetFixedCharge(ctx context.Context, externalID string, fixedChargeCode string) (*FixedCharge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", externalID, "fixed_charges", fixedChargeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
	}

	result, err := sr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult.FixedCharge, nil
}

func (sr *SubscriptionRequest) GetFixedChargeList(ctx context.Context, externalID string, fixedChargeListInput *FixedChargeListInput) (*FixedChargeResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", externalID, "fixed_charges")

	jsonQueryParams, err := json.Marshal(fixedChargeListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &FixedChargeResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult, nil
}

func (sr *SubscriptionRequest) UpdateFixedCharge(ctx context.Context, externalID string, fixedChargeCode string, fixedChargeInput *FixedChargeInput) (*FixedCharge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", externalID, "fixed_charges", fixedChargeCode)

	fixedChargeParams := &FixedChargeParams{
		FixedCharge: fixedChargeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
		Body:   fixedChargeParams,
	}

	result, err := sr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult.FixedCharge, nil
}

// Charge Filters

func (sr *SubscriptionRequest) GetChargeFilter(ctx context.Context, externalID string, chargeCode string, filterID string) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode, "filters", filterID)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
	}

	result, err := sr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}

func (sr *SubscriptionRequest) GetChargeFilterList(ctx context.Context, externalID string, chargeCode string, filterListInput *ChargeFilterListInput) (*ChargeFilterResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode, "filters")

	jsonQueryParams, err := json.Marshal(filterListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &ChargeFilterResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult, nil
}

func (sr *SubscriptionRequest) CreateChargeFilter(ctx context.Context, externalID string, chargeCode string, filterInput *ChargeFilterInput) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode, "filters")

	filterParams := &ChargeFilterParams{
		Filter: filterInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
		Body:   filterParams,
	}

	result, err := sr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}

func (sr *SubscriptionRequest) UpdateChargeFilter(ctx context.Context, externalID string, chargeCode string, filterID string, filterInput *ChargeFilterInput) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode, "filters", filterID)

	filterParams := &ChargeFilterParams{
		Filter: filterInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
		Body:   filterParams,
	}

	result, err := sr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}

func (sr *SubscriptionRequest) DeleteChargeFilter(ctx context.Context, externalID string, chargeCode string, filterID string) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "subscriptions", externalID, "charges", chargeCode, "filters", filterID)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
	}

	result, err := sr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}
