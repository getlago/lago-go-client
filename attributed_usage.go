package lago

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

type AttributedUsageBasis string

const (
	AttributedUsageBasisUnits  AttributedUsageBasis = "units"
	AttributedUsageBasisAmount AttributedUsageBasis = "amount"
)

type AttributedUsageInput struct {
	ExternalSubscriptionID string
	// GroupBy is the code of the usage attribution type to group the usage by.
	GroupBy string
	// Filters maps usage attribution type codes to the values to filter on.
	// Each filter must have at least one value.
	Filters      map[string][]string
	FromDatetime string
	ToDatetime   string
	// ChargeCodes restricts the usage to the given charges. A nil slice selects
	// all charges; an empty non-nil slice is rejected.
	ChargeCodes     []string
	SplitChargeCode string
	SearchTerm      string
	Basis           AttributedUsageBasis
	Page            *int
	PerPage         *int
}

type AttributedUsageResult struct {
	AttributedUsage *AttributedUsage `json:"attributed_usage,omitempty"`
	Meta            Metadata         `json:"meta,omitempty"`
}

type AttributedUsage struct {
	LagoSubscriptionID     uuid.UUID                `json:"lago_subscription_id,omitempty"`
	ExternalSubscriptionID string                   `json:"external_subscription_id,omitempty"`
	GroupBy                string                   `json:"group_by,omitempty"`
	Basis                  AttributedUsageBasis     `json:"basis,omitempty"`
	FromDatetime           string                   `json:"from_datetime,omitempty"`
	ToDatetime             string                   `json:"to_datetime,omitempty"`
	Currency               Currency                 `json:"currency,omitempty"`
	Rows                   []AttributedUsageRow     `json:"rows,omitempty"`
	Unattributed           AttributedUsageAggregate `json:"unattributed,omitempty"`
	Totals                 AttributedUsageAggregate `json:"totals,omitempty"`
}

type AttributedUsageAggregate struct {
	AmountCents        *int                         `json:"amount_cents"`
	PreciseAmountCents *string                      `json:"precise_amount_cents"`
	EventsCount        int                          `json:"events_count"`
	ChargesUsage       []AttributedUsageChargeUsage `json:"charges_usage,omitempty"`
}

type AttributedUsageRow struct {
	Value string `json:"value,omitempty"`
	Rank  int    `json:"rank,omitempty"`
	AttributedUsageAggregate
}

type AttributedUsageChargeUsage struct {
	LagoChargeID                   uuid.UUID           `json:"lago_charge_id,omitempty"`
	ChargeCode                     string              `json:"charge_code,omitempty"`
	BillableMetricCode             string              `json:"billable_metric_code,omitempty"`
	LagoChargeFilterID             *uuid.UUID          `json:"lago_charge_filter_id,omitempty"`
	ChargeFilterValues             map[string][]string `json:"charge_filter_values,omitempty"`
	ChargeFilterInvoiceDisplayName *string             `json:"charge_filter_invoice_display_name,omitempty"`
	Units                          string              `json:"units,omitempty"`
	AmountCents                    *int                `json:"amount_cents"`
	PreciseAmountCents             *string             `json:"precise_amount_cents"`
	EventsCount                    int                 `json:"events_count"`
}

func (cr *CustomerRequest) AttributedUsage(ctx context.Context, externalCustomerID string, attributedUsageInput *AttributedUsageInput) (*AttributedUsageResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", externalCustomerID, "attributed_usage")

	urlValues := url.Values{}
	urlValues.Set("external_subscription_id", attributedUsageInput.ExternalSubscriptionID)
	urlValues.Set("group_by", attributedUsageInput.GroupBy)

	for k, values := range attributedUsageInput.Filters {
		if len(values) == 0 {
			return nil, &Error{Err: fmt.Errorf("attributed usage filter %q requires at least one value", k)}
		}
		for _, v := range values {
			urlValues.Add(fmt.Sprintf("filters[%s][]", k), v)
		}
	}
	if attributedUsageInput.FromDatetime != "" {
		urlValues.Set("from_datetime", attributedUsageInput.FromDatetime)
	}
	if attributedUsageInput.ToDatetime != "" {
		urlValues.Set("to_datetime", attributedUsageInput.ToDatetime)
	}
	if attributedUsageInput.ChargeCodes != nil && len(attributedUsageInput.ChargeCodes) == 0 {
		return nil, &Error{Err: fmt.Errorf("attributed usage charge codes require at least one value")}
	}
	for _, code := range attributedUsageInput.ChargeCodes {
		urlValues.Add("charge_codes[]", code)
	}
	if attributedUsageInput.SplitChargeCode != "" {
		urlValues.Set("split_charge_code", attributedUsageInput.SplitChargeCode)
	}
	if attributedUsageInput.SearchTerm != "" {
		urlValues.Set("search_term", attributedUsageInput.SearchTerm)
	}
	if attributedUsageInput.Basis != "" {
		urlValues.Set("basis", string(attributedUsageInput.Basis))
	}
	if attributedUsageInput.Page != nil {
		urlValues.Set("page", strconv.Itoa(*attributedUsageInput.Page))
	}
	if attributedUsageInput.PerPage != nil {
		urlValues.Set("per_page", strconv.Itoa(*attributedUsageInput.PerPage))
	}

	clientRequest := &ClientRequest{
		Path:      subPath,
		UrlValues: urlValues,
		Result:    &AttributedUsageResult{},
	}

	result, clientErr := cr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	attributedUsageResult, ok := result.(*AttributedUsageResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return attributedUsageResult, nil
}
