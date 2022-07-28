package lago

import (
	"github.com/google/uuid"
)

type BillableMetricRequest struct {
	client *Client
}

type AggregationType string

const (
	CountAggregation       AggregationType = "count_agg"
	SumAggregation         AggregationType = "sum_agg"
	MaxAggregation         AggregationType = "max_agg"
	UniqueCountAggregation AggregationType = "unique_count_agg"
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

type BillableMetricResult struct {
	BillableMetric *BillableMetric `json:"billable_metric"`
}

type BillableMetric struct {
	LagoID uuid.UUID `json:"lago_id"`
}

func (c *Client) BillableMetric() *BillableMetricRequest {
	return &BillableMetricRequest{
		client: c,
	}
}

func (bmr *BillableMetricRequest) Create(billableMetricInput *BillableMetricInput) (*BillableMetric, *Error) {
	resp, err := bmr.client.HttpClient.
		R().
		SetError(&Error{}).
		SetResult(&BillableMetricResult{}).
		SetBody(billableMetricInput).
		Post("billable_metrics")
	if err != nil {
		return nil, &Error{Err: err}
	}

	if resp.IsError() {
		return nil, resp.Error().(*Error)
	}

	billableMetricResult := resp.Result().(*BillableMetricResult)

	return billableMetricResult.BillableMetric, nil
}
