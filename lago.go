package lago

import (
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

func (c *Client) Get(cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetError(&Error{}).
		SetResult(cr.Result).
		SetQueryParams(cr.QueryParams).
		Get(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
}

func (c *Client) Post(cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
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
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
}

func (c *Client) PostWithoutResult(cr *ClientRequest) *Error {
	resp, err := c.HttpClient.R().
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
		return resp.Error().(*Error)
	}

	return nil
}

func (c *Client) Put(cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
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
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
}

func (c *Client) Delete(cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetError(&Error{}).
		SetResult(cr.Result).
		SetBody(cr.Body).
		Delete(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if c.Debug {
		fmt.Println("REQUEST: ", resp.Request.RawRequest)
		fmt.Println("RESPONSE: ", resp.String())
	}

	if resp.IsError() {
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
}
