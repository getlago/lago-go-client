package lago

import (
	"context"
	"fmt"
)

type SubscriptionEntitlementRequest struct {
	client *Client
}

type SubscriptionEntitlement struct {
	Code        string                             `json:"code"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Privileges  []SubscriptionEntitlementPrivilege `json:"privileges"`
	Overrides   map[string]any                     `json:"overrides"`
}

type SubscriptionEntitlementPrivilege struct {
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	ValueType     string          `json:"value_type"`
	Config        PrivilegeConfig `json:"config"`
	Value         any             `json:"value"`
	PlanValue     any             `json:"plan_value"`
	OverrideValue any             `json:"override_value"`
}

type SubscriptionEntitlementResult struct {
	Entitlement  *SubscriptionEntitlement  `json:"entitlement,omitempty"`
	Entitlements []SubscriptionEntitlement `json:"entitlements,omitempty"`
}

func (c *Client) SubscriptionEntitlement() *SubscriptionEntitlementRequest {
	return &SubscriptionEntitlementRequest{
		client: c,
	}
}

func (sr *SubscriptionEntitlementRequest) GetList(ctx context.Context, subscriptionExternalId string) (*SubscriptionEntitlementResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalId, "entitlements")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &SubscriptionEntitlementResult{},
	}

	result, clientErr := sr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	subscriptionEntitlementResult, ok := result.(*SubscriptionEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionEntitlementResult, nil
}

func (sr *SubscriptionEntitlementRequest) Delete(ctx context.Context, subscriptionExternalId string, featureCode string) (*SubscriptionEntitlementResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "subscriptions", subscriptionExternalId, "entitlements", featureCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &SubscriptionEntitlementResult{},
	}

	result, clientErr := sr.client.Delete(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	subscriptionEntitlementResult, ok := result.(*SubscriptionEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionEntitlementResult, nil
}

func (sr *SubscriptionEntitlementRequest) DeletePrivilege(ctx context.Context, subscriptionExternalId string, featureCode string, privilegeCode string) (*SubscriptionEntitlementResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "subscriptions", subscriptionExternalId, "entitlements", featureCode, "privileges", privilegeCode)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &SubscriptionEntitlementResult{},
	}

	result, clientErr := sr.client.Delete(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	subscriptionEntitlementResult, ok := result.(*SubscriptionEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionEntitlementResult, nil
}
func (sr *SubscriptionEntitlementRequest) Update(ctx context.Context, subscriptionExternalId string, input []EntitlementInput) (*SubscriptionEntitlementResult, *Error) {
	return sr.update(ctx, subscriptionExternalId, input, true)
}

func (sr *SubscriptionEntitlementRequest) update(ctx context.Context, subscriptionExternalId string, input []EntitlementInput, partial bool) (*SubscriptionEntitlementResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "subscriptions", subscriptionExternalId, "entitlements")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &SubscriptionEntitlementResult{},
		Body:   EntitlementsInput{Entitlements: input},
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

	subscriptionEntitlementResult, ok := result.(*SubscriptionEntitlementResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return subscriptionEntitlementResult, nil
}
