package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"

	lt "github.com/getlago/lago-go-client/testing"
)

func TestPlanRequest_Get_WithFixedCharges(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("GET").
		MatchPath("/api/v1/plans/my-plan").
		MockResponse(map[string]any{
			"plan": map[string]any{
				"code": "my-plan",
				"fixed_charges": []map[string]any{
					{"code": "licenses", "units": "0.0"},
				},
			},
		})
	defer server.Close()

	result, err := server.Client().Plan().Get(context.Background(), "my-plan")
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result.FixedCharges[0].Code, qt.Equals, "licenses")
	c.Assert(result.FixedCharges[0].Units, qt.Equals, "0.0")
}

func TestPlanRequest_GetFixedCharge(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("GET").
		MatchPath("/api/v1/plans/my-plan/fixed_charges/licenses").
		MockResponse(map[string]any{
			"fixed_charge": map[string]any{"code": "licenses", "units": "12.5"},
		})
	defer server.Close()

	result, err := server.Client().Plan().GetFixedCharge(context.Background(), "my-plan", "licenses")
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result.Units, qt.Equals, "12.5")
}
