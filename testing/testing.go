package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	qt "github.com/frankban/quicktest"
)

func HandlerFunc(c *qt.C, mockResponse map[string]interface{}, assertRequestFunc func(*qt.C, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestFunc(c, r)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
}
