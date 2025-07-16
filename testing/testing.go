package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	qt "github.com/frankban/quicktest"
)

func HandlerFunc(c *qt.C, mockResponse interface{}, assertRequestFunc func(*qt.C, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestFunc(c, r)
		if mockResponse == nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")

		if mockResponseMap, ok := mockResponse.(map[string]interface{}); ok {
			_ = json.NewEncoder(w).Encode(mockResponseMap)
			return
		}
		if mockResponseString, ok := mockResponse.(string); ok {
			_, _ = w.Write([]byte(mockResponseString))
			return
		}

		c.Fatalf("Invalid mock response type: %T", mockResponse)
	}))
}
