package lago

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const baseURL string = "https://api.getlago.com"
const apiPath string = "/api/v1/"

var responseSuccessCodes []int = []int{200, 201, 202, 204}

type Client struct {
	HttpClient *resty.Client
}

type ClientRequest struct {
	Path   string
	Result interface{}
	Body   interface{}
}

func New() *Client {
	url := fmt.Sprintf("%s%s", baseURL, apiPath)

	restyClient := resty.New().
		SetBaseURL(url).
		SetHeader("Content-Type", "application/json")

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

func (c *Client) Get(cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetError(&Error{}).
		SetResult(cr.Result).
		Get(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
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

	if resp.IsError() {
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
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

	if resp.IsError() {
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
}

func (c *Client) Delete(cr *ClientRequest) (interface{}, *Error) {
	resp, err := c.HttpClient.R().
		SetError(&Error{}).
		SetResult(cr.Result).
		Delete(cr.Path)
	if err != nil {
		return nil, &Error{Err: err}
	}

	if resp.IsError() {
		return nil, resp.Error().(*Error)
	}

	return resp.Result(), nil
}
