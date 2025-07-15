package lago

import (
	"context"
	"encoding/json"
)

type UsageRequest struct {
	client *Client
}

type UsageCustomerType string

const (
	UsageCompanyCustomerType    UsageCustomerType = "company"
	UsageIndividualCustomerType UsageCustomerType = "individual"
)

type UsageListInput struct {
	Currency					string				`json:"currency,omitempty"`
	CustomerCountry     		string 				`json:"customer_country,omitempty,string"`
	CustomerType     			UsageCustomerType	`json:"customer_type,omitempty,string"`
	FromDate            		string          	`json:"from_date,omitempty"`
	ToDate              		string          	`json:"to_date,omitempty"`
	IsBillableMetricRecurring 	bool        	    `json:"is_billable_metric_recurring,omitempty"`
	TimeGranularity    			string   	    	`json:"time_granularity,omitempty"`
	ExternalCustomerID 			string      	    `json:"external_customer_id,omitempty"`
	ExternalSubscriptionID		string				`json:"external_subscription_id,omitempty"`
	BillableMetricCode			string				`json:"billable_metric_code,omitempty"`
	PlanCode					string				`json:"plan_code,omitempty"`
}

type UsageResult struct {
	Usage  *Usage  `json:"usage,omitempty"`
	Usages []Usage `json:"usages,omitempty"`
}

type Usage struct {
	OrganizationID			string   `json:"organization_id,omitempty"`
	StartOfPeriodDt 		string   `json:"start_of_period_dt,omitempty"`
	EndOfPeriodDt   		string   `json:"end_of_period_dt,omitempty"`
	AmountCurrency			Currency `json:"amount_currency,omitempty"`
	AmountCents				int      `json:"amount_cents,omitempty"`
	BillableMetricCode		string   `json:"billable_metric_code,omitempty"`
	Units					string      `json:"units,omitempty"`
	IsBillableMetricDeleted bool     `json:"is_billable_metric_deleted,omitempty"`
}

func (c *Client) Usage() *UsageRequest {
	return &UsageRequest{
		client: c,
	}
}

func (adr *UsageRequest) GetList(ctx context.Context, UsageListInput *UsageListInput) (*UsageResult, *Error) {
	jsonQueryparams, err := json.Marshal(UsageListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryparams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "analytics/usage",
		QueryParams: queryParams,
		Result:      &UsageResult{},
	}

	result, clientErr := adr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	UsageResult, ok := result.(*UsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return UsageResult, nil
}
