package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"

	lt "github.com/getlago/lago-go-client/testing"
)

func TestFeeRequest_Get_WithPresentationBreakdowns(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("GET").
		MatchPath("/api/v1/fees/fee_123").
		MockResponse(map[string]any{
			"fee": map[string]any{
				"presentation_breakdowns": []map[string]any{
					{"presentation_by": map[string]any{"team": "engineering"}, "units": "9.0"},
				},
			},
		})
	defer server.Close()

	result, err := server.Client().Fee().Get(context.Background(), "fee_123")
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result.PresentationBreakdowns[0].PresentationBy["team"], qt.Equals, "engineering")
	c.Assert(result.PresentationBreakdowns[0].Units, qt.Equals, "9.0")
}
