package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type FeatureRequest struct {
	client *Client
}

type PrivilegeConfig struct {
	SelectOptions []string `json:"select_options,omitempty"`
}

type Privilege struct {
	Code      string          `json:"code"`
	Name      string          `json:"name,omitempty"`
	ValueType string          `json:"value_type"` // TODO: add enum
	Config    PrivilegeConfig `json:"config"`
}

type Feature struct {
	Name        string      `json:"name,omitempty"`
	Code        string      `json:"code"`
	Description string      `json:"description,omitempty"`
	Privileges  []Privilege `json:"privileges,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
}

type FeatureResult struct {
	Feature  *Feature  `json:"feature,omitempty"`
	Features []Feature `json:"features,omitempty"`
	Meta     Metadata  `json:"meta,omitempty"`
}

type FeatureListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type FeatureParams struct {
	FeatureInput *FeatureInput `json:"feature,omitempty"`
}

type FeatureInput struct {
	Name        string           `json:"name,omitempty"`
	Code        string           `json:"code,omitempty"`
	Description string           `json:"description,omitempty"`
	Privileges  []PrivilegeInput `json:"privileges,omitempty"`
}

type ConfigInput struct {
	SelectOptions []string `json:"select_options,omitempty"`
}

type PrivilegeInput struct {
	Code      string      `json:"code,omitempty"`
	Name      string      `json:"name,omitempty"`
	ValueType string      `json:"value_type,omitempty"`
	Config    ConfigInput `json:"config,omitempty"`
}

type EntitlementsInput struct {
	Entitlements []EntitlementInput `json:"entitlements"`
}

type EntitlementsInputMap map[string]map[string]any

func (p EntitlementsInput) MarshalJSON() ([]byte, error) {
	result := make(EntitlementsInputMap, len(p.Entitlements))
	for _, ent := range p.Entitlements {
		privilegeMap := make(map[string]any, len(ent.Privileges))

		for _, privilege := range ent.Privileges {
			privilegeMap[privilege.Code] = privilege.Value
		}

		result[ent.Code] = privilegeMap
	}

	return json.Marshal(map[string]EntitlementsInputMap{"entitlements": result})
}

type EntitlementPrivilegeInput struct {
	Code  string `json:"code"`
	Value any    `json:"value"`
}

type EntitlementInput struct {
	Code       string                      `json:"code"`
	Privileges []EntitlementPrivilegeInput `json:"privileges"`
}

func (c *Client) Feature() *FeatureRequest {
	return &FeatureRequest{
		client: c,
	}
}

func (bmr *FeatureRequest) Get(ctx context.Context, featureCode string) (*Feature, *Error) {
	subPath := fmt.Sprintf("%s/%s", "features", featureCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeatureResult{},
	}

	result, err := bmr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	featureResult, ok := result.(*FeatureResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return featureResult.Feature, nil
}

func (bmr *FeatureRequest) GetList(ctx context.Context, featureListInput *FeatureListInput) (*FeatureResult, *Error) {
	jsonQueryParams, err := json.Marshal(featureListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "features",
		QueryParams: queryParams,
		Result:      &FeatureResult{},
	}

	result, clientErr := bmr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	featureResult, ok := result.(*FeatureResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return featureResult, nil
}

func (bmr *FeatureRequest) Create(ctx context.Context, featureInput *FeatureInput) (*Feature, *Error) {

	clientRequest := &ClientRequest{
		Path:   "features",
		Result: &FeatureResult{},
		Body: &FeatureParams{
			FeatureInput: featureInput,
		},
	}

	result, err := bmr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	featureResult, ok := result.(*FeatureResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return featureResult.Feature, nil
}

func (bmr *FeatureRequest) Update(ctx context.Context, featureInput *FeatureInput) (*Feature, *Error) {
	subPath := fmt.Sprintf("%s/%s", "features", featureInput.Code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeatureResult{},
		Body: &FeatureParams{
			FeatureInput: featureInput,
		},
	}

	result, err := bmr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	featureResult, ok := result.(*FeatureResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return featureResult.Feature, nil
}

func (bmr *FeatureRequest) Delete(ctx context.Context, featureCode string) (*Feature, *Error) {
	subPath := fmt.Sprintf("%s/%s", "features", featureCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeatureResult{},
	}

	result, err := bmr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	featureResult, ok := result.(*FeatureResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return featureResult.Feature, nil
}

func (bmr *FeatureRequest) DeletePrivilege(ctx context.Context, featureCode string, privilegeCode string) (*Feature, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "features", featureCode, "privileges", privilegeCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &FeatureResult{},
	}

	result, err := bmr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	featureResult, ok := result.(*FeatureResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return featureResult.Feature, nil
}
