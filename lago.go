package lago

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const baseURL string = "https://api.getlago.com"
const apiPath string = "/api/v1/"

type Client struct {
	Debug      bool
	HttpClient *resty.Client
}

type ClientRequest struct {
	Path        string
	QueryParams map[string]string
	Result      interface{}
	Body        interface{}
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

	return &Client{
		HttpClient: restyClient,
	}
}

func (c *Client) SetApiKey(apiKey string) *Client {
	c.HttpClient = c.HttpClient.SetAuthToken(apiKey)

	return c
}

func (c *Client) SetBaseURL(url string) *Client {
	customURL := fmt.Sprintf("%s%s", url, apiPath)
	c.HttpClient = c.HttpClient.SetBaseURL(customURL)

	return c
}

func (c *Client) SetDebug(debug bool) *Client {
	c.Debug = debug

	return c
}

func (c *Client) Get(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	hasResult := cr.Result != nil

	request := c.HttpClient.R().
		SetContext(ctx).
		SetError(&Error{}).
		SetQueryParams(cr.QueryParams)

	if hasResult {
		request.SetResult(cr.Result)
	}

	resp, err := request.
		Get(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		err, ok := resp.Error().(*Error)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return nil, err
	}

	if hasResult {
		return resp.Result(), nil
	}

	return resp.String(), nil
}

func (c *Client) Post(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetContext(ctx).
		SetError(&Error{}).
		SetResult(cr.Result).
		SetBody(cr.Body).
		Post(cr.Path)

	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		err, ok := resp.Error().(*Error)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return nil, err
	}

	return resp.Result(), nil
}

func (c *Client) PostWithoutResult(ctx context.Context, cr *ClientRequest) *Error {
	resp, err := c.HttpClient.R().
		SetContext(ctx).
		SetError(&Error{}).
		SetBody(cr.Body).
		Post(cr.Path)
	if err != nil {
		return &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		err, ok := resp.Error().(*Error)
		if !ok {
			return &ErrorTypeAssert
		}

		return err
	}

	return nil
}

func (c *Client) PostWithoutBody(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetContext(ctx).
		SetError(&Error{}).
		Post(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		err, ok := resp.Error().(*Error)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return nil, err
	}

	return resp.Result(), nil
}

func (c *Client) Put(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetContext(ctx).
		SetError(&Error{}).
		SetResult(cr.Result).
		SetBody(cr.Body).
		Put(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		err, ok := resp.Error().(*Error)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return nil, err
	}

	return resp.Result(), nil
}

func (c *Client) Delete(ctx context.Context, cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetContext(ctx).
		SetError(&Error{}).
		SetResult(cr.Result).
		SetBody(cr.Body).
		SetQueryParams(cr.QueryParams).
		Delete(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		err, ok := resp.Error().(*Error)
		if !ok {
			return nil, &ErrorTypeAssert
		}

		return nil, err
	}

	return resp.Result(), nil
}
