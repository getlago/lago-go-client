package lago

import (
	"context"
	"encoding/json"
	"fmt"
)

type PlanEntitlementRequest struct {
	client *Client
}

type PlanEntitlement struct {
	Code        string                     `json:"code"`
	Name        string                     `json:"name,omitempty"`
	Description string                     `json:"description,omitempty"`
	Privileges  []PlanEntitlementPrivilege `json:"privileges"`
}

type PlanEntitlementPrivilege struct {
	Code      string `json:"code"`
	Name      string `json:"name,omitempty"`
	ValueType string `json:"value_type,omitempty"`
	// TODO: Add config
	Value any `json:"value"`
}

type PlanEntitlementsInput struct {
	Entitlements []PlanEntitlementInput `json:"entitlements"`
}

type PlanEntitlementsInputMap map[string]map[string]any

func (p PlanEntitlementsInput) MarshalJSON() ([]byte, error) {
	result := make(PlanEntitlementsInputMap, len(p.Entitlements))
	for _, ent := range p.Entitlements {
		privilegeMap := make(map[string]any, len(ent.Privileges))

		for _, privilege := range ent.Privileges {
			privilegeMap[privilege.Code] = privilege.Value
		}

		result[ent.Code] = privilegeMap
	}

	return json.Marshal(map[string]PlanEntitlementsInputMap{"entitlements": result})
}

type PlanEntitlementPrivilegeInput struct {
	Code  string `json:"code"`
	Value any    `json:"value"`
}

type PlanEntitlementInput struct {
	Code       string                          `json:"code"`
	Privileges []PlanEntitlementPrivilegeInput `json:"privileges"`
}

type PlanEntitlementResult struct {
	Entitlement  *PlanEntitlement  `json:"entitlement,omitempty"`
	Entitlements []PlanEntitlement `json:"entitlements,omitempty"`
}

func (c *Client) PlanEntitlement() *PlanEntitlementRequest {
	return &PlanEntitlementRequest{
		client: c,
	}
}

func (sr *PlanEntitlementRequest) GetList(ctx context.Context, planCode string) (*PlanEntitlementResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "plans", planCode, "entitlements")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanEntitlementResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	planEntitlementResult, ok := result.(*PlanEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planEntitlementResult, nil
}

func (sr *PlanEntitlementRequest) Get(ctx context.Context, planCode string, featureCode string) (*PlanEntitlement, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "entitlements", featureCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanEntitlementResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	planEntitlementResult, ok := result.(*PlanEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planEntitlementResult.Entitlement, nil
}

func (sr *PlanEntitlementRequest) Delete(ctx context.Context, planCode string, featureCode string) (*PlanEntitlement, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "plans", planCode, "entitlements", featureCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanEntitlementResult{},
	}

	result, clientErr := sr.client.Delete(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	planEntitlementResult, ok := result.(*PlanEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planEntitlementResult.Entitlement, nil
}

func (sr *PlanEntitlementRequest) DeletePrivilege(ctx context.Context, planCode string, featureCode string, privilegeCode string) (*PlanEntitlement, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "plans", planCode, "entitlements", featureCode, "privileges", privilegeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanEntitlementResult{},
	}

	result, clientErr := sr.client.Delete(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	planEntitlementResult, ok := result.(*PlanEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planEntitlementResult.Entitlement, nil
}

func (sr *PlanEntitlementRequest) Replace(ctx context.Context, planCode string, input []PlanEntitlementInput) (*PlanEntitlementResult, *Error) {
	return sr.update(ctx, planCode, input, false)
}

func (sr *PlanEntitlementRequest) Update(ctx context.Context, planCode string, input []PlanEntitlementInput) (*PlanEntitlementResult, *Error) {
	return sr.update(ctx, planCode, input, true)
}

func (sr *PlanEntitlementRequest) update(ctx context.Context, planCode string, input []PlanEntitlementInput, partial bool) (*PlanEntitlementResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "plans", planCode, "entitlements")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &PlanEntitlementResult{},
		Body:   PlanEntitlementsInput{Entitlements: input},
	}

	var result interface{}
	var clientErr *Error

	if partial {
		result, clientErr = sr.client.Patch(ctx, clientRequest)
	} else {
		result, clientErr = sr.client.Post(ctx, clientRequest)
	}
	if clientErr != nil {
		return nil, clientErr
	}

	planEntitlementResult, ok := result.(*PlanEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return planEntitlementResult, nil
}
