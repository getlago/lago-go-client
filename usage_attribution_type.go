package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UsageAttributionTypeRole string

const (
	UsageAttributionTypeRoleHierarchical UsageAttributionTypeRole = "hierarchical"
	UsageAttributionTypeRoleFlat         UsageAttributionTypeRole = "flat"
)

type UsageAttributionTypeRequest struct {
	client *Client
}

type UsageAttributionTypeParams struct {
	UsageAttributionType *UsageAttributionTypeInput `json:"usage_attribution_type"`
}

type UsageAttributionTypeInput struct {
	Code            string                   `json:"code,omitempty"`
	Name            string                   `json:"name,omitempty"`
	Description     *string                  `json:"description,omitempty"`
	AttributionKeys []string                 `json:"attribution_keys,omitempty"`
	Role            UsageAttributionTypeRole `json:"role,omitempty"`
	// ParentCode set to an empty string detaches the type from its parent.
	ParentCode *string `json:"parent_code,omitempty"`
}

type UsageAttributionTypeListInput struct {
	PerPage    *int                     `json:"per_page,omitempty,string"`
	Page       *int                     `json:"page,omitempty,string"`
	Role       UsageAttributionTypeRole `json:"role,omitempty"`
	SearchTerm string                   `json:"search_term,omitempty"`
}

type UsageAttributionTypeResult struct {
	UsageAttributionType  *UsageAttributionType  `json:"usage_attribution_type,omitempty"`
	UsageAttributionTypes []UsageAttributionType `json:"usage_attribution_types,omitempty"`
	Meta                  Metadata               `json:"meta,omitempty"`
}

type UsageAttributionType struct {
	LagoID             uuid.UUID                `json:"lago_id,omitempty"`
	LagoOrganizationID uuid.UUID                `json:"lago_organization_id,omitempty"`
	Code               string                   `json:"code,omitempty"`
	Name               string                   `json:"name,omitempty"`
	Description        *string                  `json:"description,omitempty"`
	AttributionKeys    []string                 `json:"attribution_keys,omitempty"`
	Role               UsageAttributionTypeRole `json:"role,omitempty"`
	LagoParentID       *uuid.UUID               `json:"lago_parent_id,omitempty"`
	ParentCode         *string                  `json:"parent_code,omitempty"`
	CreatedAt          time.Time                `json:"created_at,omitempty"`
	UpdatedAt          time.Time                `json:"updated_at,omitempty"`
}

func (c *Client) UsageAttributionType() *UsageAttributionTypeRequest {
	return &UsageAttributionTypeRequest{
		client: c,
	}
}

func (uatr *UsageAttributionTypeRequest) Get(ctx context.Context, code string) (*UsageAttributionType, *Error) {
	subPath := fmt.Sprintf("%s/%s", "usage_attribution_types", code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &UsageAttributionTypeResult{},
	}

	result, err := uatr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	usageAttributionTypeResult, ok := result.(*UsageAttributionTypeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return usageAttributionTypeResult.UsageAttributionType, nil
}

func (uatr *UsageAttributionTypeRequest) GetList(ctx context.Context, usageAttributionTypeListInput *UsageAttributionTypeListInput) (*UsageAttributionTypeResult, *Error) {
	jsonQueryParams, err := json.Marshal(usageAttributionTypeListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "usage_attribution_types",
		QueryParams: queryParams,
		Result:      &UsageAttributionTypeResult{},
	}

	result, clientErr := uatr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	usageAttributionTypeResult, ok := result.(*UsageAttributionTypeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return usageAttributionTypeResult, nil
}

func (uatr *UsageAttributionTypeRequest) Create(ctx context.Context, usageAttributionTypeInput *UsageAttributionTypeInput) (*UsageAttributionType, *Error) {
	clientRequest := &ClientRequest{
		Path:   "usage_attribution_types",
		Result: &UsageAttributionTypeResult{},
		Body:   &UsageAttributionTypeParams{UsageAttributionType: usageAttributionTypeInput},
	}

	result, err := uatr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	usageAttributionTypeResult, ok := result.(*UsageAttributionTypeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return usageAttributionTypeResult.UsageAttributionType, nil
}

func (uatr *UsageAttributionTypeRequest) Update(ctx context.Context, code string, usageAttributionTypeInput *UsageAttributionTypeInput) (*UsageAttributionType, *Error) {
	subPath := fmt.Sprintf("%s/%s", "usage_attribution_types", code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &UsageAttributionTypeResult{},
		Body:   &UsageAttributionTypeParams{UsageAttributionType: usageAttributionTypeInput},
	}

	result, err := uatr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	usageAttributionTypeResult, ok := result.(*UsageAttributionTypeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return usageAttributionTypeResult.UsageAttributionType, nil
}

func (uatr *UsageAttributionTypeRequest) Delete(ctx context.Context, code string) (*UsageAttributionType, *Error) {
	subPath := fmt.Sprintf("%s/%s", "usage_attribution_types", code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &UsageAttributionTypeResult{},
	}

	result, err := uatr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	usageAttributionTypeResult, ok := result.(*UsageAttributionTypeResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return usageAttributionTypeResult.UsageAttributionType, nil
}
