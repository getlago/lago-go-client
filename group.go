package lago

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type GroupRequest struct {
	client *Client
}

type GroupResult struct {
	Groups []Group  `json:"groups,omitempty"`
	Meta   Metadata `json:"meta,omitempty"`
}

type GroupListInput struct {
	PerPage int    `json:"per_page,omitempty,string"`
	Page    int    `json:"page,omitempty,string"`
	Code    string `json:"code,omitempty"`
}

type Group struct {
	LagoID uuid.UUID `json:"lago_id,omitempty"`
	Key    string    `json:"key,omitempty"`
	Value  string    `json:"value,omitempty"`
}

func (c *Client) Group() *GroupRequest {
	return &GroupRequest{
		client: c,
	}
}

func (cr *GroupRequest) GetList(ctx context.Context, groupListInput *GroupListInput) (*GroupResult, *Error) {
	jsonQueryParams, err := json.Marshal(groupListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	subPath := fmt.Sprintf("%s/%s/%s", "billable_metrics", groupListInput.Code, "groups")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &GroupResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	groupResult := result.(*GroupResult)

	return groupResult, nil
}
