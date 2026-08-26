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

const orderFormID = "aa11aa11-aa11-aa11-aa11-aa11aa11aa11"
const orderFormCustomerID = "2b012b01-2b01-2b01-2b01-2b012b012b01"
const orderFormOwnerID = "5e345e34-5e34-5e34-5e34-5e345e345e34"

var OrderFormGetMockResponse = map[string]interface{}{
	"order_form": map[string]interface{}{
		"lago_id":               orderFormID,
		"number":                "OF-2026-0001",
		"status":                "generated",
		"void_reason":           nil,
		"expires_at":            "2026-06-30T23:59:59Z",
		"signed_at":             nil,
		"voided_at":             nil,
		"signed_document_url":   nil,
		"lago_organization_id":  "3c123c12-3c12-3c12-3c12-3c123c123c12",
		"lago_customer_id":      "2b012b01-2b01-2b01-2b01-2b012b012b01",
		"lago_quote_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_quote_version_id": "4d234d23-4d23-4d23-4d23-4d234d234d23",
		"created_at":            "2026-04-29T08:59:51Z",
		"updated_at":            "2026-04-29T08:59:51Z",
	},
}

var OrderFormSignedMockResponse = map[string]interface{}{
	"order_form": map[string]interface{}{
		"lago_id":               orderFormID,
		"number":                "OF-2026-0001",
		"status":                "signed",
		"void_reason":           nil,
		"expires_at":            "2026-06-30T23:59:59Z",
		"signed_at":             "2026-05-02T10:15:00Z",
		"voided_at":             nil,
		"signed_document_url":   "https://api.getlago.com/rails/active_storage/blobs/redirect/eyJfcmFpbHMi/OF-2026-0001",
		"lago_organization_id":  "3c123c12-3c12-3c12-3c12-3c123c123c12",
		"lago_customer_id":      "2b012b01-2b01-2b01-2b01-2b012b012b01",
		"lago_quote_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"lago_quote_version_id": "4d234d23-4d23-4d23-4d23-4d234d234d23",
		"created_at":            "2026-04-29T08:59:51Z",
		"updated_at":            "2026-05-02T10:15:00Z",
	},
}

var OrderFormGetListMockResponse = map[string]interface{}{
	"order_forms": []map[string]interface{}{
		{
			"lago_id":               orderFormID,
			"number":                "OF-2026-0001",
			"status":                "generated",
			"void_reason":           nil,
			"expires_at":            "2026-06-30T23:59:59Z",
			"signed_at":             nil,
			"voided_at":             nil,
			"signed_document_url":   nil,
			"lago_organization_id":  "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"lago_customer_id":      "2b012b01-2b01-2b01-2b01-2b012b012b01",
			"lago_quote_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_quote_version_id": "4d234d23-4d23-4d23-4d23-4d234d234d23",
			"created_at":            "2026-04-29T08:59:51Z",
			"updated_at":            "2026-04-29T08:59:51Z",
		},
		{
			"lago_id":               "bb22bb22-bb22-bb22-bb22-bb22bb22bb22",
			"number":                "OF-2026-0002",
			"status":                "voided",
			"void_reason":           "manual",
			"expires_at":            nil,
			"signed_at":             nil,
			"voided_at":             "2026-05-01T09:00:00Z",
			"signed_document_url":   nil,
			"lago_organization_id":  "3c123c12-3c12-3c12-3c12-3c123c123c12",
			"lago_customer_id":      "2b012b01-2b01-2b01-2b01-2b012b012b01",
			"lago_quote_id":         "6f456f45-6f45-6f45-6f45-6f456f456f45",
			"lago_quote_version_id": "9c789c78-9c78-9c78-9c78-9c789c789c78",
			"created_at":            "2026-04-30T08:59:51Z",
			"updated_at":            "2026-05-01T09:00:00Z",
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

func TestOrderFormRequest_Get(t *testing.T) {
	t.Run("When query for a specific order form", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/order_forms/" + orderFormID).
			MockResponse(OrderFormGetMockResponse)
		defer server.Close()

		result, err := server.Client().OrderForm().Get(context.Background(), orderFormID)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, orderFormID)
		c.Assert(result.Number, qt.Equals, "OF-2026-0001")
		c.Assert(result.Status, qt.Equals, OrderFormStatusGenerated)
		c.Assert(string(result.VoidReason), qt.Equals, "")
		c.Assert(result.ExpiresAt, qt.IsNotNil)
		c.Assert(result.ExpiresAt.Format(time.RFC3339), qt.Equals, "2026-06-30T23:59:59Z")
		c.Assert(result.SignedAt, qt.IsNil)
		c.Assert(result.SignedDocumentUrl, qt.Equals, "")
		c.Assert(result.LagoQuoteVersionID.String(), qt.Equals, "4d234d23-4d23-4d23-4d23-4d234d234d23")
	})
}

func TestOrderFormRequest_GetList(t *testing.T) {
	t.Run("When query for all order forms", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/order_forms").
			MockResponse(OrderFormGetListMockResponse)
		defer server.Close()

		result, err := server.Client().OrderForm().GetList(context.Background(), &OrderFormListInput{})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.OrderForms, qt.HasLen, 2)
		c.Assert(result.OrderForms[0].Number, qt.Equals, "OF-2026-0001")
		c.Assert(result.OrderForms[1].Status, qt.Equals, OrderFormStatusVoided)
		c.Assert(result.OrderForms[1].VoidReason, qt.Equals, OrderFormVoidReasonManual)
		c.Assert(result.OrderForms[1].ExpiresAt, qt.IsNil)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
	})

	t.Run("When filtering order forms", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/order_forms").
			MatchQuery(map[string]interface{}{
				"status[]":       []string{"generated", "signed"},
				"number[]":       []string{"OF-2026-0001"},
				"quote_number[]": []string{"QT-2026-0001"},
				"customer_id[]":  []string{orderFormCustomerID},
				"owner_id[]":     []string{orderFormOwnerID},
				"search_term":    "OF-2026",
				"expires_at_to":  "2026-12-31T23:59:59Z",
			}).
			MockResponse(OrderFormGetListMockResponse)
		defer server.Close()

		result, err := server.Client().OrderForm().GetList(context.Background(), &OrderFormListInput{
			Status:      []OrderFormStatus{OrderFormStatusGenerated, OrderFormStatusSigned},
			Number:      []string{"OF-2026-0001"},
			QuoteNumber: []string{"QT-2026-0001"},
			CustomerIDs: []uuid.UUID{uuid.MustParse(orderFormCustomerID)},
			OwnerIDs:    []uuid.UUID{uuid.MustParse(orderFormOwnerID)},
			SearchTerm:  "OF-2026",
			ExpiresAtTo: "2026-12-31T23:59:59Z",
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.OrderForms, qt.HasLen, 2)
	})
}

func TestOrderFormRequest_MarkAsSigned(t *testing.T) {
	t.Run("When signing without a payload", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/order_forms/" + orderFormID + "/mark_as_signed").
			MockResponse(OrderFormSignedMockResponse)
		defer server.Close()

		result, err := server.Client().OrderForm().MarkAsSigned(context.Background(), orderFormID, nil)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.Status, qt.Equals, OrderFormStatusSigned)
	})

	t.Run("When signing with an execution mode and date", func(t *testing.T) {
		c := qt.New(t)

		executeAt, parseErr := time.Parse(time.RFC3339, "2026-07-01T00:00:00Z")
		c.Assert(parseErr, qt.IsNil)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/order_forms/" + orderFormID + "/mark_as_signed").
			MatchJSONBody(map[string]interface{}{
				"order_form": map[string]interface{}{
					"signed_document": "data:application/pdf;base64,JVBERi0xLjQKJcfs",
					"execution_mode":  "execute_in_lago",
					"execute_at":      "2026-07-01T00:00:00Z",
				},
			}).
			MockResponse(OrderFormSignedMockResponse)
		defer server.Close()

		result, err := server.Client().OrderForm().MarkAsSigned(context.Background(), orderFormID, &OrderFormMarkAsSignedInput{
			SignedDocument: "data:application/pdf;base64,JVBERi0xLjQKJcfs",
			ExecutionMode:  OrderExecutionModeExecuteInLago,
			ExecuteAt:      &executeAt,
		})

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.SignedAt, qt.IsNotNil)
		c.Assert(result.SignedDocumentUrl, qt.Contains, "OF-2026-0001")
	})
}

func TestOrderFormRequest_Void(t *testing.T) {
	t.Run("When voiding a generated order form", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/order_forms/" + orderFormID + "/void").
			MockResponse(OrderFormGetMockResponse)
		defer server.Close()

		result, err := server.Client().OrderForm().Void(context.Background(), orderFormID)

		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
		c.Assert(result.LagoID.String(), qt.Equals, orderFormID)
	})
}
