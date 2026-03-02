package lago_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockCustomerWalletMetadataResponse = map[string]interface{}{
	"metadata": map[string]interface{}{
		"foo": "bar",
		"baz": nil,
	},
}

var mockCustomerWalletNullMetadataResponse = map[string]interface{}{
	"metadata": nil,
}

func TestCustomerWalletRequest_Replace(t *testing.T) {
	t.Run("When replace metadata is called", func(t *testing.T) {
		c := qt.New(t)

		server := lt.ServerWithAssertions(c, mockCustomerWalletMetadataResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "POST")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/customers/12345/wallets/wallet_code/metadata")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			c.Assert(requestData["metadata"], qt.IsNotNil)
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		bar := "bar"
		result, err := client.CustomerWalletMetadata().Replace(context.Background(), "12345", "wallet_code", map[string]*string{
			"foo": &bar,
			"baz": nil,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
		c.Assert(result["baz"], qt.IsNil)
	})
}

func TestCustomerWalletRequest_Merge(t *testing.T) {
	t.Run("When merge metadata is called", func(t *testing.T) {
		c := qt.New(t)

		server := lt.ServerWithAssertions(c, mockCustomerWalletMetadataResponse, func(c *qt.C, r *http.Request) {
			c.Assert(r.Method, qt.Equals, "PATCH")
			c.Assert(r.URL.Path, qt.Equals, "/api/v1/customers/12345/wallets/wallet_code/metadata")

			body, err := io.ReadAll(r.Body)
			c.Assert(err, qt.IsNil)

			var requestData map[string]interface{}
			err = json.Unmarshal(body, &requestData)
			c.Assert(err, qt.IsNil)

			c.Assert(requestData["metadata"], qt.IsNotNil)
		})
		defer server.Close()

		client := New().SetBaseURL(server.URL).SetApiKey("test_api_key")
		qux := "qux"
		result, err := client.CustomerWalletMetadata().Merge(context.Background(), "12345", "wallet_code", map[string]*string{
			"foo": &qux,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
		c.Assert(result["baz"], qt.IsNil)
	})
}

func TestCustomerWalletRequest_DeleteAll(t *testing.T) {
	t.Run("When delete all metadata is called", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/12345/wallets/wallet_code/metadata").
			MockResponse(mockCustomerWalletNullMetadataResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletMetadata().DeleteAll(context.Background(), "12345", "wallet_code")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNil)
	})
}

func TestCustomerWalletRequest_DeleteKey(t *testing.T) {
	t.Run("When delete metadata key is called", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/customers/12345/wallets/wallet_code/metadata/foo").
			MockResponse(mockCustomerWalletMetadataResponse)
		defer server.Close()

		result, err := server.Client().CustomerWalletMetadata().DeleteKey(context.Background(), "12345", "wallet_code", "foo")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result["foo"], qt.IsNotNil)
		c.Assert(*result["foo"], qt.Equals, "bar")
	})
}
