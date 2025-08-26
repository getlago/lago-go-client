package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const ActivityLogsEndpoint string = "activity_logs"

type ActivityLogRequest struct {
	client *Client
}

type ActivitySource string

const (
	ApiSource    ActivitySource = "api"
	FrontSource  ActivitySource = "front"
	SystemSource ActivitySource = "system"
)

type ActivityLog struct {
	ActivityId             uuid.UUID              `json:"activity_id"`
	ActivityType           string                 `json:"activity_type"`
	ActivitySource         ActivitySource         `json:"activity_source"`
	ActivityObject         map[string]interface{} `json:"activity_object,omitempty"`
	ActivityObjectChanges  map[string]interface{} `json:"activity_object_changes,omitempty"`
	UserEmail              string                 `json:"user_email,omitempty"`
	ResourceId             uuid.UUID              `json:"rounding_precision"`
	ResourceType           string                 `json:"resource_type"`
	ExternalCustomerId     string                 `json:"external_customer_id"`
	ExternalSubscriptionId string                 `json:"external_subscription_id"`
	LoggedAt               time.Time              `json:"logged_at"`
	CreatedAt              time.Time              `json:"created_at"`
}

type ActivityLogListInput struct {
	PerPage *int `json:"per_page,omitempty,string"`
	Page    *int `json:"page,omitempty,string"`

	FromDate               string   `json:"from_date,omitempty"`
	ToDate                 string   `json:"to_date,omitempty"`
	ActivityTypes          []string `json:"activity_types,omitempty"`
	ActivitySources        []string `json:"activity_sources,omitempty"`
	UserEmails             []string `json:"user_emails,omitempty"`
	ExternalCustomerId     string   `json:"external_customer_id,omitempty"`
	ExternalSubscriptionId string   `json:"external_subscription_id,omitempty"`
	ResourceIds            []string `json:"resource_ids,omitempty"`
	ResourceTypes          []string `json:"resource_types,omitempty"`
}

type ActivityLogResult struct {
	ActivityLog  *ActivityLog  `json:"activity_log,omitempty"`
	ActivityLogs []ActivityLog `json:"activity_logs,omitempty"`
	Meta         Metadata      `json:"meta,omitempty"`
}

func (c *Client) ActivityLog() *ActivityLogRequest {
	return &ActivityLogRequest{
		client: c,
	}
}

func (alr *ActivityLogRequest) Get(ctx context.Context, ActivityId string) (*ActivityLog, *Error) {
	subPath := fmt.Sprintf("%s/%s", ActivityLogsEndpoint, ActivityId)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &ActivityLogResult{},
	}

	result, err := alr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	activityLogResult, ok := result.(*ActivityLogResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return activityLogResult.ActivityLog, nil
}

func (alr *ActivityLogRequest) GetList(ctx context.Context, activityLogListInput *ActivityLogListInput) (*ActivityLogResult, *Error) {
	jsonQueryParams, err := json.Marshal(activityLogListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        ActivityLogsEndpoint,
		QueryParams: queryParams,
		Result:      &ActivityLogResult{},
	}

	result, clientErr := alr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	activityLogResult, ok := result.(*ActivityLogResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return activityLogResult, nil
}
