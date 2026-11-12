package lago

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// WebhookEndpointV2Expand names an association the v2 webhook endpoints API
// can add to its responses through the `expands[]` query param.
type WebhookEndpointV2Expand string

// WebhookEndpointV2ExpandSigningKeys asks for the `signing_keys` of the
// endpoint. They are secrets, returned only to an API key with write access.
const WebhookEndpointV2ExpandSigningKeys WebhookEndpointV2Expand = "signing_keys"

type WebhookEndpointV2Request struct {
	client *Client
}

// WebhookEndpointV2ListInput holds the cursor pagination params of the v2
// list. Only the params that are set are sent.
type WebhookEndpointV2ListInput struct {
	Limit             *int
	After             *string
	Before            *string
	IncludeTotalCount *bool
}

// CursorMetadata is the `meta` of a v2 list response.
type CursorMetadata struct {
	NextCursor          *string `json:"next_cursor"`
	PrevCursor          *string `json:"prev_cursor"`
	TotalCount          *int    `json:"total_count,omitempty"`
	TotalCountEstimated bool    `json:"total_count_estimated,omitempty"`
}

// WebhookEndpointV2Result is the body of a v2 webhook endpoints response. It
// reuses WebhookEndpoint, whose SigningKeys holds secrets and is filled only
// when WebhookEndpointV2ExpandSigningKeys is asked with a write API key.
type WebhookEndpointV2Result struct {
	WebhookEndpoint  *WebhookEndpoint  `json:"webhook_endpoint,omitempty"`
	WebhookEndpoints []WebhookEndpoint `json:"webhook_endpoints,omitempty"`
	Meta             CursorMetadata    `json:"meta"`
}

// WebhookEndpointV2 returns the client of the v2 webhook endpoints API
// (`/api/v2/webhook_endpoints`). Other resources stay on v1.
func (c *Client) WebhookEndpointV2() *WebhookEndpointV2Request {
	return &WebhookEndpointV2Request{
		client: c,
	}
}

func (wer *WebhookEndpointV2Request) path(webhookEndpointID string) string {
	return wer.client.apiV2URL(fmt.Sprintf("%s/%s", "webhook_endpoints", webhookEndpointID))
}

// expandsQueryParams builds the `expands[]` query param of a POST or PUT,
// which only accept one value per key. Only one expansion exists, so the
// last value given is sent.
func expandsQueryParams(expands []WebhookEndpointV2Expand) map[string]string {
	queryParams := make(map[string]string)
	for _, expand := range expands {
		queryParams["expands[]"] = string(expand)
	}

	return queryParams
}

func (wer *WebhookEndpointV2Request) GetList(ctx context.Context, listInput *WebhookEndpointV2ListInput) (*WebhookEndpointV2Result, *Error) {
	urlValues := url.Values{}
	if listInput != nil {
		if listInput.Limit != nil {
			urlValues.Set("limit", strconv.Itoa(*listInput.Limit))
		}
		if listInput.After != nil {
			urlValues.Set("after", *listInput.After)
		}
		if listInput.Before != nil {
			urlValues.Set("before", *listInput.Before)
		}
		if listInput.IncludeTotalCount != nil {
			urlValues.Set("include_total_count", strconv.FormatBool(*listInput.IncludeTotalCount))
		}
	}

	clientRequest := &ClientRequest{
		Path:      wer.client.apiV2URL("webhook_endpoints"),
		UrlValues: urlValues,
		Result:    &WebhookEndpointV2Result{},
	}

	result, err := wer.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointV2Result)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult, nil
}

// Get returns a webhook endpoint. Pass WebhookEndpointV2ExpandSigningKeys to
// get its signing keys.
func (wer *WebhookEndpointV2Request) Get(ctx context.Context, webhookEndpointID string, expands ...WebhookEndpointV2Expand) (*WebhookEndpoint, *Error) {
	urlValues := url.Values{}
	for _, expand := range expands {
		urlValues.Add("expands[]", string(expand))
	}

	clientRequest := &ClientRequest{
		Path:      wer.path(webhookEndpointID),
		UrlValues: urlValues,
		Result:    &WebhookEndpointV2Result{},
	}

	result, err := wer.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointV2Result)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}

// Create creates a webhook endpoint. Pass WebhookEndpointV2ExpandSigningKeys
// to get the signing keys of a `standard_webhooks` endpoint.
func (wer *WebhookEndpointV2Request) Create(ctx context.Context, webhookEndpointInput *WebhookEndpointInput, expands ...WebhookEndpointV2Expand) (*WebhookEndpoint, *Error) {
	clientRequest := &ClientRequest{
		Path:        wer.client.apiV2URL("webhook_endpoints"),
		QueryParams: expandsQueryParams(expands),
		Result:      &WebhookEndpointV2Result{},
		Body:        &WebhookEndpointParams{WebhookEndpointInput: webhookEndpointInput},
	}

	result, err := wer.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointV2Result)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}

// Update updates a webhook endpoint. Pass WebhookEndpointV2ExpandSigningKeys
// to get the signing keys of a `standard_webhooks` endpoint.
func (wer *WebhookEndpointV2Request) Update(ctx context.Context, webhookEndpointInput *WebhookEndpointInput, webhookEndpointID string, expands ...WebhookEndpointV2Expand) (*WebhookEndpoint, *Error) {
	clientRequest := &ClientRequest{
		Path:        wer.path(webhookEndpointID),
		QueryParams: expandsQueryParams(expands),
		Result:      &WebhookEndpointV2Result{},
		Body:        &WebhookEndpointParams{WebhookEndpointInput: webhookEndpointInput},
	}

	result, err := wer.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointV2Result)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}

func (wer *WebhookEndpointV2Request) Delete(ctx context.Context, webhookEndpointID string) (*WebhookEndpoint, *Error) {
	clientRequest := &ClientRequest{
		Path:   wer.path(webhookEndpointID),
		Result: &WebhookEndpointV2Result{},
	}

	result, err := wer.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	webhookEndpointResult, ok := result.(*WebhookEndpointV2Result)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return webhookEndpointResult.WebhookEndpoint, nil
}
