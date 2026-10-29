package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

var mockUsageAttributionTypeResponse = map[string]any{
	"usage_attribution_type": map[string]any{
		"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
		"code":                 "team",
		"name":                 "Team",
		"description":          "Engineering teams",
		"attribution_keys":     []string{"team_id"},
		"role":                 "hierarchical",
		"lago_parent_id":       "3c903c90-3c90-3c90-3c90-3c903c903c90",
		"parent_code":          "department",
		"created_at":           "2026-09-01T10:00:00Z",
		"updated_at":           "2026-09-01T10:00:00Z",
	},
}

var mockUsageAttributionTypesResponse = map[string]any{
	"usage_attribution_types": []map[string]any{
		{
			"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_organization_id": "2b902b90-2b90-2b90-2b90-2b902b902b90",
			"code":                 "team",
			"name":                 "Team",
			"description":          nil,
			"attribution_keys":     []string{"team_id"},
			"role":                 "flat",
			"lago_parent_id":       nil,
			"parent_code":          nil,
			"created_at":           "2026-09-01T10:00:00Z",
			"updated_at":           "2026-09-01T10:00:00Z",
		},
	},
	"meta": map[string]any{
		"current_page": 1,
		"next_page":    nil,
		"prev_page":    nil,
		"total_pages":  1,
		"total_count":  1,
	},
}

func TestUsageAttributionTypeRequest_Get(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.UsageAttributionType().Get(context.Background(), "team")
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
	})

	t.Run("When the usage attribution type is found", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/usage_attribution_types/team").
			MockResponse(mockUsageAttributionTypeResponse)
		defer server.Close()

		result, err := server.Client().UsageAttributionType().Get(context.Background(), "team")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.LagoID.String(), qt.Equals, "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(result.Code, qt.Equals, "team")
		c.Assert(*result.Description, qt.Equals, "Engineering teams")
		c.Assert(result.AttributionKeys, qt.DeepEquals, []string{"team_id"})
		c.Assert(result.Role, qt.Equals, UsageAttributionTypeRoleHierarchical)
		c.Assert(result.LagoParentID.String(), qt.Equals, "3c903c90-3c90-3c90-3c90-3c903c903c90")
		c.Assert(*result.ParentCode, qt.Equals, "department")
	})
}

func TestUsageAttributionTypeRequest_GetList(t *testing.T) {
	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/usage_attribution_types").
			MatchQuery(map[string]interface{}{
				"page":        "1",
				"per_page":    "10",
				"role":        "flat",
				"search_term": "team",
			}).
			MockResponse(mockUsageAttributionTypesResponse)
		defer server.Close()

		page, perPage := 1, 10
		result, err := server.Client().UsageAttributionType().GetList(context.Background(), &UsageAttributionTypeListInput{
			Page:       &page,
			PerPage:    &perPage,
			Role:       UsageAttributionTypeRoleFlat,
			SearchTerm: "team",
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.UsageAttributionTypes, qt.HasLen, 1)
		c.Assert(result.UsageAttributionTypes[0].Role, qt.Equals, UsageAttributionTypeRoleFlat)
		c.Assert(result.UsageAttributionTypes[0].LagoParentID, qt.IsNil)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})
}

func TestUsageAttributionTypeRequest_Create(t *testing.T) {
	t.Run("Wraps the body under the usage_attribution_type key", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/usage_attribution_types").
			MatchJSONBody(`{"usage_attribution_type":{"code":"team","name":"Team","description":"Engineering teams","attribution_keys":["team_id"],"role":"hierarchical","parent_code":"department"}}`).
			MockResponse(mockUsageAttributionTypeResponse)
		defer server.Close()

		description, parentCode := "Engineering teams", "department"
		result, err := server.Client().UsageAttributionType().Create(context.Background(), &UsageAttributionTypeInput{
			Code:            "team",
			Name:            "Team",
			Description:     &description,
			AttributionKeys: []string{"team_id"},
			Role:            UsageAttributionTypeRoleHierarchical,
			ParentCode:      &parentCode,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "team")
	})
}

func TestUsageAttributionTypeRequest_Update(t *testing.T) {
	t.Run("Sends an empty parent_code to detach the parent", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/usage_attribution_types/team").
			MatchJSONBody(`{"usage_attribution_type":{"name":"Team","parent_code":""}}`).
			MockResponse(mockUsageAttributionTypeResponse)
		defer server.Close()

		parentCode := ""
		result, err := server.Client().UsageAttributionType().Update(context.Background(), "team", &UsageAttributionTypeInput{
			Name:       "Team",
			ParentCode: &parentCode,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "team")
	})
}

func TestUsageAttributionTypeRequest_Delete(t *testing.T) {
	t.Run("When the usage attribution type is deleted", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("DELETE").
			MatchPath("/api/v1/usage_attribution_types/team").
			MockResponse(mockUsageAttributionTypeResponse)
		defer server.Close()

		result, err := server.Client().UsageAttributionType().Delete(context.Background(), "team")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Code, qt.Equals, "team")
	})
}
