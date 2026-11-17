package lago_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	qt "github.com/frankban/quicktest"

	. "github.com/getlago/lago-go-client"
)

func TestUsageGetListRecurringFilter(t *testing.T) {
	c := qt.New(t)

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("is_billable_metric_recurring")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"usages":[]}`))
	}))
	defer server.Close()

	client := New().SetBaseURL(server.URL)
	result, err := client.Usage().GetList(context.Background(), &UsageListInput{IsBillableMetricRecurring: true})

	c.Assert(err == nil, qt.IsTrue, qt.Commentf("%+v", err))
	c.Assert(result, qt.IsNotNil)
	c.Assert(gotQuery, qt.Equals, "true")
}
