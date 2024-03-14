package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SignatureAlgo string

const (
	JWT  SignatureAlgo = "jwt"
	HMac SignatureAlgo = "hmac"
)

type WebhookEndpointRequest struct {
	client *Client
}

type WebhookEndpointParams struct {
	WebhookEndpointInput *WebhookEndpointInput `json:"webhook_endpoint"`
}

type WebhookEndpointInput struct {
	WebhookURL    string        `json:"webhook_url,omitempty"`
	SignatureAlgo SignatureAlgo `json:"signature_algo,omitempty"`
}

type WebhookEndpointListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type WebhookEndpointResult struct {
	WebhookEndpoint  *WebhookEndpoint  `json:"webhook_endpoint,omitempty"`
	WebhookEndpoints []WebhookEndpoint `json:"webhook_endpoints,omitempty"`
	Meta             Metadata          `json:"meta,omitempty"`
}

type WebhookEndpoint struct {
	LagoID             uuid.UUID     `json:"lago_id,omitempty"`
	LagoOrganizationID uuid.UUID     `json:"lago_organization_id,omitempty"`
	WebhookURL         string        `json:"webhook_url,omitempty"`
	SignatureAlgo      SignatureAlgo `json:"signature_algo,omitempty"`
	CreatedAt          time.Time     `json:"created_at,omitempty"`
}

func (c *Client) WebhookEndpoint() *WebhookEndpointRequest {
	return &WebhookEndpointRequest{
		client: c,
	}
}

func (wer *WebhookEndpointRequest) Get(ctx context.Context, webhookEndpointID string) (*WebhookEndpoint, *Error) {
	subPath := fmt.Sprintf("%s/%s", "webhook_endpoints", webhookEndpointID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WebhookEndpointResult{},
	}

	result, err := wer.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}

func (wer *WebhookEndpointRequest) GetList(ctx context.Context, webhookEndpointListInput *WebhookEndpointListInput) (*WebhookEndpointResult, *Error) {
	jsonQueryParams, err := json.Marshal(webhookEndpointListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "webhook_endpoints",
		QueryParams: queryParams,
		Result:      &WebhookEndpointResult{},
	}

	result, clientErr := wer.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult, nil
}

func (wer *WebhookEndpointRequest) Create(ctx context.Context, webhookEndpointInput *WebhookEndpointInput) (*WebhookEndpoint, *Error) {
	webhookEndpointParams := &WebhookEndpointParams{
		WebhookEndpointInput: webhookEndpointInput,
	}

	clientRequest := &ClientRequest{
		Path:   "webhook_endpoints",
		Result: &WebhookEndpointResult{},
		Body:   webhookEndpointParams,
	}

	result, err := wer.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}

func (wer *WebhookEndpointRequest) Update(ctx context.Context, webhookEndpointInput *WebhookEndpointInput, webhookEndpointID string) (*WebhookEndpoint, *Error) {
	subPath := fmt.Sprintf("%s/%s", "webhook_endpoints", webhookEndpointID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WebhookEndpointResult{},
		Body:   webhookEndpointInput,
	}

	result, err := wer.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}

func (wer *WebhookEndpointRequest) Delete(ctx context.Context, webhookEndpointID string) (*WebhookEndpoint, *Error) {
	subPath := fmt.Sprintf("%s/%s", "webhook_endpoints", webhookEndpointID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WebhookEndpointResult{},
	}

	result, err := wer.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}
