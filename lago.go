package lago

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

// executeWithRetry runs the given request function with automatic retry on HTTP 429 responses.
// It respects the client's RetryPolicy settings including max attempts, backoff, and
// the x-ratelimit-reset header from the server.
//
// On a non-429 response, if the RetryPolicy.OnRateLimitInfo callback is set,
// the parsed x-ratelimit-* headers are delivered to it for observability.
func (c *Client) executeWithRetry(ctx context.Context, fn func() (*resty.Response, error)) (*resty.Response, error) {
	for attempt := 0; ; attempt++ {
		resp, err := fn()
		if err != nil {
			return resp, err
		}

		// If not rate limited, emit observability info and return
		if resp.StatusCode() != 429 {
			c.emitRateLimitInfo(resp)
			return resp, nil
		}

		// Rate limited but retries disabled or max attempts reached: return as-is
		if c.RetryPolicy == nil ||
			!c.RetryPolicy.EnableRetry ||
			attempt >= c.RetryPolicy.MaxAttempts-1 {
			return resp, nil
		}

		// Calculate wait duration from headers or exponential backoff
		waitDuration := c.RetryPolicy.waitDuration(resp.RawResponse, attempt)
		timer := time.NewTimer(waitDuration)
		select {
		case <-timer.C:
			// continue to next attempt
		case <-ctx.Done():
			timer.Stop()
			return resp, ctx.Err()
		}
	}
}

// emitRateLimitInfo invokes the configured OnRateLimitInfo callback (if any)
// with parsed x-ratelimit-* headers from the given response.
func (c *Client) emitRateLimitInfo(resp *resty.Response) {
	if c.RetryPolicy == nil || c.RetryPolicy.OnRateLimitInfo == nil || resp == nil || resp.RawResponse == nil {
		return
	}
	method := ""
	url := ""
	if resp.Request != nil && resp.Request.RawRequest != nil {
		method = resp.Request.RawRequest.Method
		if resp.Request.RawRequest.URL != nil {
			url = resp.Request.RawRequest.URL.String()
		}
	}
	c.RetryPolicy.emitRateLimitInfo(resp.RawResponse, method, url)
}

const baseURL string = "https://api.getlago.com"
const baseIngestURL string = "https://ingest.getlago.com"
const apiPath string = "/api/v1/"

type Client struct {
	BaseUrl          string
	BaseIngestUrl    string
	UseIngestService bool
	Debug            bool
	HttpClient       *resty.Client
	IngestHttpClient *resty.Client
	RetryPolicy      *RetryPolicy
}

type ClientRequest struct {
	UseIngestService bool
	Path             string
	QueryParams      map[string]string
	UrlValues        url.Values
	Result           interface{}
	Body             interface{}
}

type Metadata struct {
	CurrentPage int `json:"current_page,omitempty"`
	NextPage    int `json:"next_page,omitempty"`
	PrevPage    int `json:"prev_page,omitempty"`
	TotalPages  int `json:"total_pages,omitempty"`
	TotalCount  int `json:"total_count,omitempty"`
}

func New() *Client {
	url := fmt.Sprintf("%s%s", baseURL, apiPath)

	restyClient := resty.New().
		SetBaseURL(url).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "lago-go-client github.com/getlago/lago-go-client/v1")

	ingestRestyClient := resty.New().
		SetBaseURL(url).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "lago-go-client github.com/getlago/lago-go-client/v1")

	retryPolicy := DefaultRetryPolicy()

	return &Client{
		BaseUrl:          url,
		BaseIngestUrl:    url,
		HttpClient:       restyClient,
		IngestHttpClient: ingestRestyClient,
		RetryPolicy:      retryPolicy,
	}
}

func (c *Client) SetApiKey(apiKey string) *Client {
	c.HttpClient = c.HttpClient.SetAuthToken(apiKey)
	c.IngestHttpClient = c.IngestHttpClient.SetAuthToken(apiKey)

	return c
}

func (c *Client) SetBaseURL(url string) *Client {
	c.BaseUrl = url
	c.BaseIngestUrl = url

	customURL := fmt.Sprintf("%s%s", url, apiPath)
	c.HttpClient = c.HttpClient.SetBaseURL(customURL)
	c.IngestHttpClient = c.IngestHttpClient.SetBaseURL(customURL)

	return c
}

func (c *Client) SetDebug(debug bool) *Client {
	c.Debug = debug

	return c
}

// SetRetryPolicy updates the retry policy for both HTTP clients.
// This allows customization of rate limit retry behavior.
// Pass nil to disable retries.
func (c *Client) SetRetryPolicy(policy *RetryPolicy) *Client {
	if policy == nil {
		// Disable retries by creating a no-op policy
		c.RetryPolicy = &RetryPolicy{EnableRetry: false}
		return c
	}

	c.RetryPolicy = policy
	return c
}

func (c *Client) SetUseIngestService(useIngestService bool) *Client {
	c.UseIngestService = useIngestService

	if useIngestService {
		c = c.SetBaseIngestUrl(baseIngestURL)
	} else {
		c = c.SetBaseURL(c.BaseUrl)
	}

	return c
}

func (c *Client) SetBaseIngestUrl(url string) *Client {
	c.BaseIngestUrl = url

	customURL := fmt.Sprintf("%s%s", url, apiPath)
	c.IngestHttpClient = c.IngestHttpClient.SetBaseURL(customURL)

	return c
}

// handleErrorResponse converts an error HTTP response into the appropriate error type.
// For 429 responses, it returns a RateLimitError with parsed rate limit headers.
// For all other errors, it returns the standard Error type.
func (c *Client) handleErrorResponse(resp *resty.Response) *Error {
	errObj, ok := resp.Error().(*Error)
	if !ok {
		return &ErrorTypeAssert
	}

	if resp.StatusCode() == 429 {
		rlErr := ParseRateLimitError(errObj, resp.RawResponse.Header)
		return &Error{
			Err:            rlErr,
			HTTPStatusCode: rlErr.HTTPStatusCode,
			Message:        rlErr.Message,
			ErrorCode:      rlErr.ErrorCode,
			ErrorDetail:    rlErr.ErrorDetail,
		}
	}

	return errObj
}

func (c *Client) Get(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	hasResult := cr.Result != nil

	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		request := c.HttpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetQueryParams(cr.QueryParams).
			SetQueryParamsFromValues(cr.UrlValues)

		if hasResult {
			request.SetResult(cr.Result)
		}

		return request.Get(cr.Path)
	})
	if retryErr != nil {
		return nil, &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, c.handleErrorResponse(resp)
	}

	if hasResult {
		return resp.Result(), nil
	}

	return resp.String(), nil
}

func (c *Client) Patch(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	httpClient := c.HttpClient
	if cr.UseIngestService {
		httpClient = c.IngestHttpClient
	}

	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		return httpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetResult(cr.Result).
			SetBody(cr.Body).
			SetQueryParams(cr.QueryParams).
			Patch(cr.Path)
	})
	if retryErr != nil {
		return nil, &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, c.handleErrorResponse(resp)
	}

	return resp.Result(), nil
}

func (c *Client) Post(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	httpClient := c.HttpClient
	if cr.UseIngestService {
		httpClient = c.IngestHttpClient
	}

	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		return httpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetResult(cr.Result).
			SetBody(cr.Body).
			SetQueryParams(cr.QueryParams).
			Post(cr.Path)
	})
	if retryErr != nil {
		return nil, &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, c.handleErrorResponse(resp)
	}

	return resp.Result(), nil
}

func (c *Client) PostWithoutResult(ctx context.Context, cr *ClientRequest) *Error {
	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		return c.HttpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetBody(cr.Body).
			Post(cr.Path)
	})
	if retryErr != nil {
		return &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return c.handleErrorResponse(resp)
	}

	return nil
}

func (c *Client) PostWithoutBody(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		return c.HttpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetResult(cr.Result).
			Post(cr.Path)
	})
	if retryErr != nil {
		return nil, &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, c.handleErrorResponse(resp)
	}

	return resp.Result(), nil
}

func (c *Client) Put(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		return c.HttpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetResult(cr.Result).
			SetBody(cr.Body).
			SetQueryParams(cr.QueryParams).
			Put(cr.Path)
	})
	if retryErr != nil {
		return nil, &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, c.handleErrorResponse(resp)
	}

	return resp.Result(), nil
}

func (c *Client) Delete(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	hasResult := cr.Result != nil

	resp, retryErr := c.executeWithRetry(ctx, func() (*resty.Response, error) {
		request := c.HttpClient.R().
			SetContext(ctx).
			SetError(&Error{}).
			SetBody(cr.Body).
			SetQueryParams(cr.QueryParams)

		if hasResult {
			request.SetResult(cr.Result)
		}

		return request.Delete(cr.Path)
	})
	if retryErr != nil {
		return nil, &Error{Err: retryErr}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, c.handleErrorResponse(resp)
	}

	return resp.Result(), nil
}
