package lago

import (
	"context"
	"encoding/json"
	"fmt"
)

// Charges

func (pr *PlanRequest) GetCharge(ctx context.Context, planCode string, chargeCode string) (*Charge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "charges", chargeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeResult{},
	}

	result, err := pr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult.Charge, nil
}

func (pr *PlanRequest) GetChargeList(ctx context.Context, planCode string, chargeListInput *ChargeListInput) (*ChargeResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "plans", planCode, "charges")

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

	result, clientErr := pr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult, nil
}

func (pr *PlanRequest) CreateCharge(ctx context.Context, planCode string, chargeInput *ChargeInput) (*Charge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "plans", planCode, "charges")

	chargeParams := &ChargeParams{
		Charge: chargeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeResult{},
		Body:   chargeParams,
	}

	result, err := pr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult.Charge, nil
}

func (pr *PlanRequest) UpdateCharge(ctx context.Context, planCode string, chargeCode string, chargeInput *ChargeInput) (*Charge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "charges", chargeCode)

	chargeParams := &ChargeParams{
		Charge: chargeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeResult{},
		Body:   chargeParams,
	}

	result, err := pr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	chargeResult, ok := result.(*ChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return chargeResult.Charge, nil
}

func (pr *PlanRequest) DeleteCharge(ctx context.Context, planCode string, chargeCode string) (*Charge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "charges", chargeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeResult{},
	}

	result, err := pr.client.Delete(ctx, clientRequest)
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

func (pr *PlanRequest) GetFixedCharge(ctx context.Context, planCode string, fixedChargeCode string) (*FixedCharge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "fixed_charges", fixedChargeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
	}

	result, err := pr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult.FixedCharge, nil
}

func (pr *PlanRequest) GetFixedChargeList(ctx context.Context, planCode string, fixedChargeListInput *FixedChargeListInput) (*FixedChargeResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "plans", planCode, "fixed_charges")

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

	result, clientErr := pr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult, nil
}

func (pr *PlanRequest) CreateFixedCharge(ctx context.Context, planCode string, fixedChargeInput *FixedChargeInput) (*FixedCharge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "plans", planCode, "fixed_charges")

	fixedChargeParams := &FixedChargeParams{
		FixedCharge: fixedChargeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
		Body:   fixedChargeParams,
	}

	result, err := pr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult.FixedCharge, nil
}

func (pr *PlanRequest) UpdateFixedCharge(ctx context.Context, planCode string, fixedChargeCode string, fixedChargeInput *FixedChargeInput) (*FixedCharge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "fixed_charges", fixedChargeCode)

	fixedChargeParams := &FixedChargeParams{
		FixedCharge: fixedChargeInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
		Body:   fixedChargeParams,
	}

	result, err := pr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	fixedChargeResult, ok := result.(*FixedChargeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return fixedChargeResult.FixedCharge, nil
}

func (pr *PlanRequest) DeleteFixedCharge(ctx context.Context, planCode string, fixedChargeCode string) (*FixedCharge, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "fixed_charges", fixedChargeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FixedChargeResult{},
	}

	result, err := pr.client.Delete(ctx, clientRequest)
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

func (pr *PlanRequest) GetChargeFilter(ctx context.Context, planCode string, chargeCode string, filterID string) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "plans", planCode, "charges", chargeCode, "filters", filterID)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
	}

	result, err := pr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}

func (pr *PlanRequest) GetChargeFilterList(ctx context.Context, planCode string, chargeCode string, filterListInput *ChargeFilterListInput) (*ChargeFilterResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "plans", planCode, "charges", chargeCode, "filters")

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

	result, clientErr := pr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult, nil
}

func (pr *PlanRequest) CreateChargeFilter(ctx context.Context, planCode string, chargeCode string, filterInput *ChargeFilterInput) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "plans", planCode, "charges", chargeCode, "filters")

	filterParams := &ChargeFilterParams{
		Filter: filterInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
		Body:   filterParams,
	}

	result, err := pr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}

func (pr *PlanRequest) UpdateChargeFilter(ctx context.Context, planCode string, chargeCode string, filterID string, filterInput *ChargeFilterInput) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "plans", planCode, "charges", chargeCode, "filters", filterID)

	filterParams := &ChargeFilterParams{
		Filter: filterInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
		Body:   filterParams,
	}

	result, err := pr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}

func (pr *PlanRequest) DeleteChargeFilter(ctx context.Context, planCode string, chargeCode string, filterID string) (*ChargeFilterResponse, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "plans", planCode, "charges", chargeCode, "filters", filterID)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ChargeFilterResult{},
	}

	result, err := pr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	filterResult, ok := result.(*ChargeFilterResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return filterResult.Filter, nil
}
