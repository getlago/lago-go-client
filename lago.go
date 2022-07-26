package lago

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

const baseURL string = "https://api.getlago.com/"
const apiPath string = "api/v1/"

type ClientConfig struct {
	ApiKey string
	ApiURL *string
}
type Client struct {
	Config     ClientConfig
	HttpClient *http.Client
}

func NewClient(config ClientConfig) *Client {
	if config.ApiURL == nil {
		urlValue := baseURL
		config.ApiURL = &urlValue
	}

	return &Client{
		Config:     config,
		HttpClient: &http.Client{},
	}
}

func (c *Client) Post(path string, body string) string {
	url := fmt.Sprintf("%s%s%s", *c.Config.ApiURL, apiPath, path)
	bodyReader := strings.NewReader(body)

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		panic(err)
	}

	authHeader := fmt.Sprintf("Bearer %s", c.Config.ApiKey)
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Lago Go") // TODO: Add version number

	response, err := c.HttpClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", response.StatusCode)

	buff := new(bytes.Buffer)
	buff.ReadFrom(response.Body)

	return string(buff.String())
}
