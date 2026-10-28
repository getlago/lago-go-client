package lago_test

import (
	"context"
	"testing"
	"time"

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

func TestPlanRequest_Get_WithResponseFields(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("GET").
		MatchPath("/api/v1/plans/my-plan").
		MockResponse(map[string]any{
			"plan": map[string]any{
				"code":             "my-plan",
				"created_at":       "2026-10-01T12:30:00Z",
				"trial_period":     14.0,
				"pending_deletion": true,
				"applicable_usage_thresholds": []map[string]any{
					{"threshold_display_name": "Alert", "amount_cents": 10000, "recurring": false},
					{"threshold_display_name": nil, "amount_cents": 500, "recurring": true},
				},
			},
		})
	defer server.Close()

	result, err := server.Client().Plan().Get(context.Background(), "my-plan")
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result.CreatedAt, qt.Equals, time.Date(2026, 10, 1, 12, 30, 0, 0, time.UTC))
	c.Assert(result.TrialPeriod, qt.IsNotNil)
	c.Assert(*result.TrialPeriod, qt.Equals, float32(14.0))
	c.Assert(result.PendingDeletion, qt.IsTrue)
	c.Assert(result.ApplicableUsageThresholds, qt.HasLen, 2)
	c.Assert(result.ApplicableUsageThresholds[0].ThresholdDisplayName, qt.IsNotNil)
	c.Assert(*result.ApplicableUsageThresholds[0].ThresholdDisplayName, qt.Equals, "Alert")
	c.Assert(result.ApplicableUsageThresholds[0].AmountCents, qt.Equals, 10000)
	c.Assert(result.ApplicableUsageThresholds[0].Recurring, qt.IsFalse)
	c.Assert(result.ApplicableUsageThresholds[1].ThresholdDisplayName, qt.IsNil)
	c.Assert(result.ApplicableUsageThresholds[1].AmountCents, qt.Equals, 500)
	c.Assert(result.ApplicableUsageThresholds[1].Recurring, qt.IsTrue)
}

func TestPlanRequest_Get_WithNullTrialPeriod(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("GET").
		MatchPath("/api/v1/plans/my-plan").
		MockResponse(map[string]any{
			"plan": map[string]any{
				"code":                        "my-plan",
				"trial_period":                nil,
				"pending_deletion":            false,
				"applicable_usage_thresholds": []map[string]any{},
			},
		})
	defer server.Close()

	result, err := server.Client().Plan().Get(context.Background(), "my-plan")
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result.TrialPeriod, qt.IsNil)
	c.Assert(result.PendingDeletion, qt.IsFalse)
	c.Assert(result.ApplicableUsageThresholds, qt.HasLen, 0)
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
