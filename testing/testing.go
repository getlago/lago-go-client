package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"

	qt "github.com/frankban/quicktest"
	"github.com/getlago/lago-go-client"
)

func handlerFuncWithResponse(c *qt.C, responseFunc func() interface{}, assertRequestFunc func(*qt.C, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestFunc(c, r)
		mockResponse := responseFunc()
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
	})
}
func ServerWithAssertions(c *qt.C, mockResponse interface{}, assertRequestFunc func(*qt.C, *http.Request)) *httptest.Server {
	responseFunc := func() interface{} {
		return mockResponse
	}
	return httptest.NewServer(handlerFuncWithResponse(c, responseFunc, assertRequestFunc))
}

type MockServer struct {
	c              *qt.C
	server         *httptest.Server
	called         bool
	expectedMethod string
	expectedPath   string
	expectedQuery  *string
	mockResponse   interface{}
}

func NewMockServer(c *qt.C) *MockServer {
	mockServer := &MockServer{c: c}
	responseFunc := func() interface{} {
		return mockServer.mockResponse
	}
	mockServer.server = httptest.NewServer(handlerFuncWithResponse(c, responseFunc, func(c *qt.C, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		mockServer.c.Assert(apiKey, qt.Equals, "Bearer test_api_key")
		mockServer.called = true
		if mockServer.expectedMethod != "" {
			c.Assert(r.Method, qt.Equals, mockServer.expectedMethod)
		}
		if mockServer.expectedPath != "" {
			c.Assert(r.URL.Path, qt.Equals, mockServer.expectedPath)
		}
		if mockServer.expectedQuery != nil {
			parsedQuery, err := url.ParseQuery(*mockServer.expectedQuery)
			mockServer.c.Assert(err, qt.IsNil)
			c.Assert(r.URL.Query(), qt.DeepEquals, parsedQuery)
		}
	}))
	return mockServer
}

func (m *MockServer) MatchMethod(method string) *MockServer {
	m.expectedMethod = method
	return m
}

func (m *MockServer) MatchPath(path string) *MockServer {
	m.expectedPath = path
	return m
}

func (m *MockServer) MatchQuery(queryParams interface{}) *MockServer {
	switch mapQueryParams := queryParams.(type) {
	case map[string]string:
	case map[string][]string:
	case map[string]interface{}:
		urlValues := url.Values{}
		for key, value := range mapQueryParams {
			switch stringOrArrayValue := value.(type) {
			case string:
				urlValues.Add(key, stringOrArrayValue)
			case []string:
				for _, v := range stringOrArrayValue {
					urlValues.Add(key, v)
				}
			}
		}
		str := urlValues.Encode()
		m.expectedQuery = &str
		return m
	case string:
		str := queryParams.(string)
		m.expectedQuery = &str
		return m
	default:
		m.c.Fatalf("Invalid query params type: %T", queryParams)
	}

	if queryString, ok := queryParams.(string); ok {
		m.expectedQuery = &queryString
		return m
	}
	m.c.Fatalf("Invalid query params type: %T", queryParams)
	return m
}

func (m *MockServer) MockResponse(mockResponse interface{}) *MockServer {
	m.mockResponse = mockResponse
	return m
}

func (m *MockServer) Close() {
	m.server.Close()
	m.c.Assert(m.called, qt.IsTrue)
}

func (m MockServer) Client() *lago.Client {
	return lago.New().SetBaseURL(m.server.URL).SetApiKey("test_api_key")
}
