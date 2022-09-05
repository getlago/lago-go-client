package lago

import (
	"encoding/json"
	"fmt"

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
	Name            string          `json:"name,omitempty"`
	Code            string          `json:"code,omitempty"`
	Description     string          `json:"description,omitempty"`
	AggregationType AggregationType `json:"aggregation_type,omitempty"`
	FieldName       string          `json:"field_name"`
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
	LagoID uuid.UUID `json:"lago_id"`
}

func (c *Client) BillableMetric() *BillableMetricRequest {
	return &BillableMetricRequest{
		client: c,
	}
}

func (bmr *BillableMetricRequest) Get(billableMetricCode string) (*BillableMetric, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billable_metrics", billableMetricCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillableMetricResult{},
	}

	result, err := bmr.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult := result.(*BillableMetricResult)

	return billableMetricResult.BillableMetric, nil
}

func (bmr *BillableMetricRequest) GetList(billableMetricListInput *BillableMetricListInput) (*BillableMetricResult, *Error) {
	jsonQueryParams, err := json.Marshal(billableMetricListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	json.Unmarshal(jsonQueryParams, &queryParams)

	clientRequest := &ClientRequest{
		Path:        "billable_metrics",
		QueryParams: queryParams,
		Result:      &BillableMetricResult{},
	}

	result, clientErr := bmr.client.Get(clientRequest)
	if err != nil {
		return nil, clientErr
	}

	billableMetricResult := result.(*BillableMetricResult)

	return billableMetricResult, nil
}

func (bmr *BillableMetricRequest) Create(billableMetricInput *BillableMetricInput) (*BillableMetric, *Error) {
	clientRequest := &ClientRequest{
		Path:   "billable_metrics",
		Result: &BillableMetricResult{},
		Body:   billableMetricInput,
	}

	result, err := bmr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult := result.(*BillableMetricResult)

	return billableMetricResult.BillableMetric, nil
}

func (bmr *BillableMetricRequest) Update(billableMetricInput *BillableMetricInput) (*BillableMetric, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billable_metrics", billableMetricInput.Code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillableMetricResult{},
		Body:   billableMetricInput,
	}

	result, err := bmr.client.Put(clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult := result.(*BillableMetricResult)

	return billableMetricResult.BillableMetric, nil
}

func (bmr *BillableMetricRequest) Delete(billableMetricCode string) (*BillableMetric, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billable_metrics", billableMetricCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillableMetricResult{},
	}

	result, err := bmr.client.Delete(clientRequest)
	if err != nil {
		return nil, err
	}

	billableMetricResult := result.(*BillableMetricResult)

	return billableMetricResult.BillableMetric, nil
}
