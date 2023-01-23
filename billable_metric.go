package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BillableMetricRequest struct {
	client *Client
}

type AggregationType string

const (
	CountAggregation          AggregationType = "count_agg"
	SumAggregation            AggregationType = "sum_agg"
	MaxAggregation            AggregationType = "max_agg"
	UniqueCountAggregation    AggregationType = "unique_count_agg"
	RecurringCountAggregation AggregationType = "recurring_count_agg"
)

type BillableMetricParams struct {
	BillableMetricInput *BillableMetricInput
}

type BillableMetricInput struct {
	Name            string                 `json:"name,omitempty"`
	Code            string                 `json:"code,omitempty"`
	Description     string                 `json:"description,omitempty"`
	AggregationType AggregationType        `json:"aggregation_type,omitempty"`
	FieldName       string                 `json:"field_name"`
	Group           map[string]interface{} `json:"group,omitempty"`
}

type BillableMetricListInput struct {
	PerPage int `json:"per_page,omitempty,string"`
	Page    int `json:"page,omitempty,string"`
}

type BillableMetricResult struct {
	BillableMetric  *BillableMetric  `json:"billable_metric,omitempty"`
	BillableMetrics []BillableMetric `json:"billable_metrics,omitempty"`
	Meta            Metadata         `json:"meta,omitempty"`
}

type BillableMetric struct {
	LagoID                   uuid.UUID              `json:"lago_id"`
	Name                     string                 `json:"name,omitempty"`
	Code                     string                 `json:"code,omitempty"`
	Description              string                 `json:"description,omitempty"`
	AggregationType          AggregationType        `json:"aggregation_type,omitempty"`
	FieldName                string                 `json:"field_name"`
	CreatedAt                time.Time              `json:"created_at,omitempty"`
	Group                    map[string]interface{} `json:"group,omitempty"`
	ActiveSubscriptionsCount int                    `json:"active_subscriptions_count,omitempty"`
	DraftInvoicesCount       int                    `json:"draft_invoices_count,omitempty"`
}

func (c *Client) BillableMetric() *BillableMetricRequest {
	return &BillableMetricRequest{
		client: c,
	}
}

func (bmr *BillableMetricRequest) Get(ctx context.Context, billableMetricCode string) (*BillableMetric, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billable_metrics", billableMetricCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillableMetricResult{},
	}

	result, err := bmr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult, ok := result.(*BillableMetricResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billableMetricResult.BillableMetric, nil
}

func (bmr *BillableMetricRequest) GetList(ctx context.Context, billableMetricListInput *BillableMetricListInput) (*BillableMetricResult, *Error) {
	jsonQueryParams, err := json.Marshal(billableMetricListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "billable_metrics",
		QueryParams: queryParams,
		Result:      &BillableMetricResult{},
	}

	result, clientErr := bmr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, clientErr
	}

	billableMetricResult, ok := result.(*BillableMetricResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billableMetricResult, nil
}

func (bmr *BillableMetricRequest) Create(ctx context.Context, billableMetricInput *BillableMetricInput) (*BillableMetric, *Error) {
	clientRequest := &ClientRequest{
		Path:   "billable_metrics",
		Result: &BillableMetricResult{},
		Body:   billableMetricInput,
	}

	result, err := bmr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult, ok := result.(*BillableMetricResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billableMetricResult.BillableMetric, nil
}

func (bmr *BillableMetricRequest) Update(ctx context.Context, billableMetricInput *BillableMetricInput) (*BillableMetric, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billable_metrics", billableMetricInput.Code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillableMetricResult{},
		Body:   billableMetricInput,
	}

	result, err := bmr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult, ok := result.(*BillableMetricResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billableMetricResult.BillableMetric, nil
}

func (bmr *BillableMetricRequest) Delete(ctx context.Context, billableMetricCode string) (*BillableMetric, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billable_metrics", billableMetricCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillableMetricResult{},
	}

	result, err := bmr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult, ok := result.(*BillableMetricResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billableMetricResult.BillableMetric, nil
}
