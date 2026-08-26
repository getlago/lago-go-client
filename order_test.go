package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
	"github.com/google/uuid"
)

const orderID = "cc33cc33-cc33-cc33-cc33-cc33cc33cc33"
const orderCustomerID = "2b012b01-2b01-2b01-2b01-2b012b012b01"
const orderOwnerID = "5e345e34-5e34-5e34-5e34-5e345e345e34"

var orderBillingSnapshotMock = map[string]interface{}{
	"plans": []map[string]interface{}{
		{
			"id":      "7a567a56-7a56-7a56-7a56-7a567a567a56",
			"localId": "b5c1e2a4-4e1e-4a7f-9f0e-9c1a0c7e1f2b",
			"type":    "plan",
			"payload": map[string]interface{}{"code": "premium_plan"},
		},
	},
}

var OrderGetMockResponse = map[string]interface{}{
	"order": map[string]interface{}{
		"lago_id":        orderID,
		"number":         "OR-2026-0001",
		"status":         "created",
		"order_type":     "subscription_creation",
		"execution_mode": "execute_in_lago",
		"currency":       "EUR",
		"executed_at":    nil,
		"execution_record": map[string]interface{}{
			"executed_at":                 nil,
			"execution_mode":              "execute_in_lago",
			"invoice_id":                  nil,
			"subscription_ids":            []string{},
			"terminated_subscription_ids": []string{},
			"applied_coupon_ids":          []string{},
			"wallet_ids":                  []string{},
			"errors":                      []string{},
		},
		"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
		"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
		"lago_order_form_id":   "aa11aa11-aa11-aa11-aa11-aa11aa11aa11",
		"created_at":           "2026-05-02T10:15:00Z",
		"updated_at":           "2026-05-02T10:15:00Z",
		"billing_snapshot":     orderBillingSnapshotMock,
	},
}

var OrderExecutedMockResponse = map[string]interface{}{
	"order": map[string]interface{}{
		"lago_id":        orderID,
		"number":         "OR-2026-0001",
		"status":         "executed",
		"order_type":     "subscription_creation",
		"execution_mode": "execute_in_lago",
		"currency":       "EUR",
		"executed_at":    "2026-07-01T00:00:00Z",
		"execution_record": map[string]interface{}{
			"executed_at":                 "2026-07-01T00:00:00Z",
			"execution_mode":              "execute_in_lago",
			"invoice_id":                  nil,
			"subscription_ids":            []string{"dd44dd44-dd44-dd44-dd44-dd44dd44dd44"},
			"terminated_subscription_ids": []string{},
			"applied_coupon_ids":          []string{"ee55ee55-ee55-ee55-ee55-ee55ee55ee55"},
			"wallet_ids":                  []string{},
			"errors":                      []string{},
		},
		"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
		"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
		"lago_order_form_id":   "aa11aa11-aa11-aa11-aa11-aa11aa11aa11",
		"created_at":           "2026-05-02T10:15:00Z",
		"updated_at":           "2026-07-01T00:00:00Z",
		"billing_snapshot":     orderBillingSnapshotMock,
	},
}

var OrderGetListMockResponse = map[string]interface{}{
	"orders": []map[string]interface{}{
		{
			"lago_id":        orderID,
			"number":         "OR-2026-0001",
			"status":         "created",
			"order_type":     "subscription_creation",
			"execution_mode": "execute_in_lago",
			"currency":       "EUR",
			"executed_at":    nil,
			"execution_record": map[string]interface{}{
				"executed_at":                 nil,
				"execution_mode":              "execute_in_lago",
				"invoice_id":                  nil,
				"subscription_ids":            []string{},
				"terminated_subscription_ids": []string{},
				"applied_coupon_ids":          []string{},
				"wallet_ids":                  []string{},
				"errors":                      []string{},
			},
			"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
			"lago_order_form_id":   "aa11aa11-aa11-aa11-aa11-aa11aa11aa11",
			"created_at":           "2026-05-02T10:15:00Z",
			"updated_at":           "2026-05-02T10:15:00Z",
			"billing_snapshot":     map[string]interface{}{},
		},
		{
			"lago_id":        "ff66ff66-ff66-ff66-ff66-ff66ff66ff66",
			"number":         "OR-2026-0002",
			"status":         "failed",
			"order_type":     "one_off",
			"execution_mode": nil,
			"currency":       "EUR",
			"executed_at":    nil,
			"execution_record": map[string]interface{}{
				"executed_at":                 nil,
				"execution_mode":              nil,
				"invoice_id":                  nil,
				"subscription_ids":            []string{},
				"terminated_subscription_ids": []string{},
				"applied_coupon_ids":          []string{},
				"wallet_ids":                  []string{},
				"errors":                      []string{"plan_not_found"},
			},
			"lago_organization_id": "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"lago_customer_id":     "2b012b01-2b01-2b01-2b01-2b012b012b01",
			"lago_order_form_id":   "bb22bb22-bb22-bb22-bb22-bb22bb22bb22",
			"created_at":           "2026-05-03T10:15:00Z",
			"updated_at":           "2026-05-03T10:15:00Z",
			"billing_snapshot":     map[string]interface{}{},
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  2,
	},
}

func TestOrderRequest_Get(t *testing.T) {
	t.Run("When query for a specific order", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/orders/" + orderID).
			MockResponse(OrderGetMockResponse)
		defer server.Close()

		result, err := server.Client().Order().Get(context.Background(), orderID)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, orderID)
		c.Assert(result.Number, qt.Equals, "OR-2026-0001")
		c.Assert(result.Status, qt.Equals, OrderStatusCreated)
		c.Assert(result.OrderType, qt.Equals, QuoteOrderTypeSubscriptionCreation)
		c.Assert(result.ExecutionMode, qt.Equals, OrderExecutionModeExecuteInLago)
		c.Assert(result.Currency, qt.Equals, EUR)
		c.Assert(result.ExecutedAt, qt.IsNil)
		c.Assert(result.LagoOrderFormID.String(), qt.Equals, "aa11aa11-aa11-aa11-aa11-aa11aa11aa11")

		c.Assert(result.ExecutionRecord, qt.IsNotNil)
		c.Assert(result.ExecutionRecord.ExecutedAt, qt.IsNil)
		c.Assert(result.ExecutionRecord.InvoiceID, qt.IsNil)
		c.Assert(result.ExecutionRecord.Errors, qt.HasLen, 0)

		c.Assert(result.BillingSnapshot, qt.IsNotNil)
		c.Assert(result.BillingSnapshot.Plans, qt.HasLen, 1)
		c.Assert(result.BillingSnapshot.Plans[0].Payload["code"], qt.Equals, "premium_plan")
	})
}

func TestOrderRequest_GetList(t *testing.T) {
	t.Run("When query for all orders", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/orders").
			MockResponse(OrderGetListMockResponse)
		defer server.Close()

		result, err := server.Client().Order().GetList(context.Background(), &OrderListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Orders, qt.HasLen, 2)
		c.Assert(result.Orders[0].Number, qt.Equals, "OR-2026-0001")
		c.Assert(result.Orders[1].Status, qt.Equals, OrderStatusFailed)
		c.Assert(string(result.Orders[1].ExecutionMode), qt.Equals, "")
		c.Assert(result.Orders[1].ExecutionRecord.Errors, qt.DeepEquals, []string{"plan_not_found"})
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
	})

	t.Run("When filtering orders", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/orders").
			MatchQuery(map[string]interface{}{
				"status[]":            []string{"created", "executed"},
				"order_type[]":        []string{"one_off"},
				"execution_mode[]":    []string{"execute_in_lago"},
				"order_form_number[]": []string{"OF-2026-0001"},
				"quote_number[]":      []string{"QT-2026-0001"},
				"customer_id[]":       []string{orderCustomerID},
				"owner_id[]":          []string{orderOwnerID},
				"search_term":         "OR-2026",
			}).
			MockResponse(OrderGetListMockResponse)
		defer server.Close()

		result, err := server.Client().Order().GetList(context.Background(), &OrderListInput{
			Status:          []OrderStatus{OrderStatusCreated, OrderStatusExecuted},
			OrderType:       []QuoteOrderType{QuoteOrderTypeOneOff},
			ExecutionMode:   []OrderExecutionMode{OrderExecutionModeExecuteInLago},
			OrderFormNumber: []string{"OF-2026-0001"},
			QuoteNumber:     []string{"QT-2026-0001"},
			CustomerIDs:     []uuid.UUID{uuid.MustParse(orderCustomerID)},
			OwnerIDs:        []uuid.UUID{uuid.MustParse(orderOwnerID)},
			SearchTerm:      "OR-2026",
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Orders, qt.HasLen, 2)
	})
}

func TestOrderRequest_Execute(t *testing.T) {
	t.Run("When executing without a payload", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/orders/" + orderID + "/execute").
			MockResponse(OrderExecutedMockResponse)
		defer server.Close()

		result, err := server.Client().Order().Execute(context.Background(), orderID, nil)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Status, qt.Equals, OrderStatusExecuted)
		c.Assert(result.ExecutedAt, qt.IsNotNil)
		c.Assert(result.ExecutedAt.Format(time.RFC3339), qt.Equals, "2026-07-01T00:00:00Z")
		c.Assert(result.ExecutionRecord.SubscriptionIDs, qt.HasLen, 1)
		c.Assert(result.ExecutionRecord.SubscriptionIDs[0].String(), qt.Equals, "dd44dd44-dd44-dd44-dd44-dd44dd44dd44")
		c.Assert(result.ExecutionRecord.AppliedCouponIDs, qt.HasLen, 1)
		c.Assert(result.ExecutionRecord.TerminatedSubscriptionIDs, qt.HasLen, 0)
	})

	t.Run("When restating the execution mode", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/orders/" + orderID + "/execute").
			MatchJSONBody(map[string]interface{}{
				"order": map[string]interface{}{"execution_mode": "execute_in_lago"},
			}).
			MockResponse(OrderExecutedMockResponse)
		defer server.Close()

		result, err := server.Client().Order().Execute(context.Background(), orderID, &OrderExecuteInput{
			ExecutionMode: OrderExecutionModeExecuteInLago,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Status, qt.Equals, OrderStatusExecuted)
	})
}
