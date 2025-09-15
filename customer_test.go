package lago_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"

	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

// Mock response for customer invoice list
var mockCustomerInvoiceListResponse = map[string]any{
	"invoices": []map[string]interface{}{
		mockInvoice,
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

// Mock response for customer credit note list
var mockCustomerCreditNoteListResponse = map[string]any{
	"credit_notes": []map[string]interface{}{
		{
			"lago_id":                                 "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"sequential_id":                           2,
			"number":                                  "LAG-1234-CN-001-002",
			"lago_invoice_id":                         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"invoice_number":                          "LAG-1234-001-002",
			"issuing_date":                            "2022-04-30",
			"credit_status":                           "available",
			"refund_status":                           "pending",
			"reason":                                  "duplicated_charge",
			"description":                             "Duplicated charge",
			"currency":                                "EUR",
			"total_amount_cents":                      120,
			"credit_amount_cents":                     100,
			"balance_amount_cents":                    100,
			"refund_amount_cents":                     0,
			"coupons_adjustment_amount_cents":         0,
			"taxes_amount_cents":                      20,
			"sub_total_excluding_taxes_amount_cents":  100,
			"max_creditable_amount_cents":             100,
			"max_refundable_amount_cents":             100,
			"precise_coupons_adjustment_amount_cents": "0.0",
			"created_at":                              "2022-04-29T08:59:51Z",
			"updated_at":                              "2022-04-29T08:59:51Z",
			"file_url":                                "https://getlago.com/credit_note/file",
			"voided_at":                               nil,
			"self_billed":                             false,
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

// Mock response for customer payment list
var mockCustomerPaymentListResponse = map[string]any{
	"payments": []map[string]interface{}{
		{
			"lago_id":              "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"amount_cents":         1200,
			"amount_currency":      "EUR",
			"payment_status":       "succeeded",
			"type":                 "manual",
			"reference":            "REF-123456",
			"external_payment_id":  "ext_payment_123",
			"created_at":           "2022-04-29T08:59:51Z",
			"invoice_ids":          []string{"1a901a90-1a90-1a90-1a90-1a901a901a90"},
			"external_customer_id": "CUSTOMER_1",
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

// Mock response for customer payment request list
var mockCustomerPaymentRequestListResponse = map[string]any{
	"payment_requests": []map[string]interface{}{
		{
			"lago_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"email":           "customer@example.com",
			"amount_cents":    1200,
			"amount_currency": "EUR",
			"payment_status":  "pending",
			"created_at":      "2022-04-29T08:59:51Z",
			"customer": map[string]interface{}{
				"lago_id":     "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"external_id": "CUSTOMER_1",
				"name":        "John Doe",
				"email":       "customer@example.com",
			},
			"invoices": []map[string]interface{}{
				{
					"lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
					"number":  "LAG-1234-001-002",
				},
			},
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

// Mock response for customer applied coupon list
var mockCustomerAppliedCouponListResponse = map[string]any{
	"applied_coupons": []map[string]interface{}{
		{
			"lago_id":                      "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"lago_coupon_id":               "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"coupon_code":                  "APPLIED_COUPON",
			"coupon_name":                  "Startup Deal",
			"lago_customer_id":             "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"external_customer_id":         "CUSTOMER_1",
			"status":                       "active",
			"amount_cents":                 2000,
			"amount_cents_remaining":       50,
			"amount_currency":              "EUR",
			"percentage_rate":              nil,
			"frequency":                    "recurring",
			"frequency_duration":           3,
			"frequency_duration_remaining": 1,
			"expiration_at":                "2022-04-29T08:59:51Z",
			"created_at":                   "2022-04-29T08:59:51Z",
			"terminated_at":                "2022-04-29T08:59:51Z",
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

// Mock response for customer subscription list
var mockCustomerSubscriptionListResponse = map[string]any{
	"subscriptions": []map[string]interface{}{
		{
			"lago_id":                           "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"external_id":                       "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
			"lago_customer_id":                  "1a901a90-1a90-1a90-1a90-1a901a901a90",
			"external_customer_id":              "CUSTOMER_1",
			"billing_time":                      "anniversary",
			"name":                              "Repository A",
			"plan_code":                         "premium",
			"status":                            "active",
			"created_at":                        "2022-08-08T00:00:00Z",
			"started_at":                        "2022-08-08T00:00:00Z",
			"subscription_at":                   "2022-08-08T00:00:00Z",
			"current_billing_period_started_at": "2022-08-08T00:00:00Z",
			"current_billing_period_ending_at":  "2022-09-08T00:00:00Z",
			"on_termination_credit_note":        "skip",
			"on_termination_invoice":            "skip",
			"plan": map[string]interface{}{
				"lago_id":         "1a901a90-1a90-1a90-1a90-1a901a901a90",
				"name":            "Premium Plan",
				"code":            "premium",
				"interval":        "monthly",
				"amount_cents":    10000,
				"amount_currency": "USD",
			},
		},
	},
	"meta": map[string]interface{}{
		"current_page": 1,
		"next_page":    0,
		"prev_page":    0,
		"total_pages":  1,
		"total_count":  1,
	},
}

func TestCustomerGetInvoiceList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().GetInvoiceList(context.Background(), "CUSTOMER_1", &CustomerInvoiceListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/customers/CUSTOMER_1/invoices\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/invoices").
			MatchQuery("").
			MockResponse(mockCustomerInvoiceListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetInvoiceList(context.Background(), "CUSTOMER_1", &CustomerInvoiceListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Invoices, qt.HasLen, 1)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/invoices").
			MatchQuery("per_page=10&page=1&issuing_date_from=2022-01-01&issuing_date_to=2022-12-31&invoice_type=subscription&status=finalized&payment_status=succeeded&payment_overdue=true&partially_paid=false&self_billed=true&payment_dispute_lost=false&amount_from=100&amount_to=1000&search_term=test&currency=EUR").
			MockResponse(mockCustomerInvoiceListResponse)
		defer server.Close()

		perPage := 10
		page := 1
		paymentOverdue := true
		partiallyPaid := false
		selfBilled := true
		paymentDisputeLost := false
		amountFrom := 100
		amountTo := 1000
		result, err := server.Client().Customer().GetInvoiceList(context.Background(), "CUSTOMER_1", &CustomerInvoiceListInput{
			PerPage:            &perPage,
			Page:               &page,
			IssuingDateFrom:    "2022-01-01",
			IssuingDateTo:      "2022-12-31",
			InvoiceType:        SubscriptionInvoiceType,
			Status:             InvoiceStatusFinalized,
			PaymentStatus:      InvoicePaymentStatusSucceeded,
			PaymentOverdue:     &paymentOverdue,
			PartiallyPaid:      &partiallyPaid,
			SelfBilled:         &selfBilled,
			PaymentDisputeLost: &paymentDisputeLost,
			AmountFrom:         &amountFrom,
			AmountTo:           &amountTo,
			SearchTerm:         "test",
			Currency:           EUR,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Invoices, qt.HasLen, 1)
	})
}

func TestCustomerGetCreditNoteList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().GetCreditNoteList(context.Background(), "CUSTOMER_1", &CustomerCreditNoteListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/customers/CUSTOMER_1/credit_notes\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/credit_notes").
			MatchQuery("").
			MockResponse(mockCustomerCreditNoteListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetCreditNoteList(context.Background(), "CUSTOMER_1", &CustomerCreditNoteListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.CreditNotes, qt.HasLen, 1)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/credit_notes").
			MatchQuery("per_page=10" +
				"&page=1" +
				"&issuing_date_from=2022-01-01" +
				"&issuing_date_to=2022-12-31" +
				"&amount_from=100" +
				"&amount_to=1000" +
				"&search_term=test" +
				"&credit_status=available" +
				"&currency=EUR" +
				"&invoice_number=INV-001" +
				"&reason=duplicated_charge" +
				"&refund_status=pending" +
				"&self_billed=true").
			MockResponse(mockCustomerCreditNoteListResponse)
		defer server.Close()

		perPage := 10
		page := 1
		selfBilled := true
		result, err := server.Client().Customer().GetCreditNoteList(context.Background(), "CUSTOMER_1", &CustomerCreditNoteListInput{
			PerPage:         &perPage,
			Page:            &page,
			IssuingDateFrom: "2022-01-01",
			IssuingDateTo:   "2022-12-31",
			AmountFrom:      100,
			AmountTo:        1000,
			SearchTerm:      "test",
			CreditStatus:    CreditNoteCreditStatusAvailable,
			Currency:        EUR,
			InvoiceNumber:   "INV-001",
			Reason:          CreditNoteReasonDuplicatedCharge,
			RefundStatus:    CreditNoteRefundStatusPending,
			SelfBilled:      &selfBilled,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.CreditNotes, qt.HasLen, 1)
	})
}

func TestCustomerGetPaymentList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().GetPaymentList(context.Background(), "CUSTOMER_1", &CustomerPaymentListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/customers/CUSTOMER_1/payments\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/payments").
			MatchQuery("").
			MockResponse(mockCustomerPaymentListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetPaymentList(context.Background(), "CUSTOMER_1", &CustomerPaymentListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Payments, qt.HasLen, 1)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/payments").
			MatchQuery("per_page=10&page=1&invoice_id=invoice_123").
			MockResponse(mockCustomerPaymentListResponse)
		defer server.Close()

		perPage := 10
		page := 1
		result, err := server.Client().Customer().GetPaymentList(context.Background(), "CUSTOMER_1", &CustomerPaymentListInput{
			PerPage:   &perPage,
			Page:      &page,
			InvoiceID: "invoice_123",
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Payments, qt.HasLen, 1)
	})
}

func TestCustomerGetPaymentRequestList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().GetPaymentRequestList(context.Background(), "CUSTOMER_1", &CustomerPaymentRequestListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/customers/CUSTOMER_1/payment_requests\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_requests").
			MatchQuery("").
			MockResponse(mockCustomerPaymentRequestListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetPaymentRequestList(context.Background(), "CUSTOMER_1", &CustomerPaymentRequestListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentRequests, qt.HasLen, 1)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/payment_requests").
			MatchQuery("per_page=10&page=1&payment_status=pending").
			MockResponse(mockCustomerPaymentRequestListResponse)
		defer server.Close()

		perPage := 10
		page := 1
		result, err := server.Client().Customer().GetPaymentRequestList(context.Background(), "CUSTOMER_1", &CustomerPaymentRequestListInput{
			PerPage:       &perPage,
			Page:          &page,
			PaymentStatus: "pending",
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentRequests, qt.HasLen, 1)
	})
}

func TestCustomerGetAppliedCouponList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().GetAppliedCouponList(context.Background(), "CUSTOMER_1", &CustomerAppliedCouponListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/customers/CUSTOMER_1/applied_coupons\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/applied_coupons").
			MatchQuery("").
			MockResponse(mockCustomerAppliedCouponListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetAppliedCouponList(context.Background(), "CUSTOMER_1", &CustomerAppliedCouponListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.AppliedCoupons, qt.HasLen, 1)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/applied_coupons").
			MatchQuery("per_page=10&page=1&status=active&coupon_code[]=COUPON1&coupon_code[]=COUPON2").
			MockResponse(mockCustomerAppliedCouponListResponse)
		defer server.Close()

		perPage := 10
		page := 1
		result, err := server.Client().Customer().GetAppliedCouponList(context.Background(), "CUSTOMER_1", &CustomerAppliedCouponListInput{
			PerPage:    &perPage,
			Page:       &page,
			Status:     AppliedCouponStatusActive,
			CouponCode: []string{"COUPON1", "COUPON2"},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.AppliedCoupons, qt.HasLen, 1)
	})
}

func TestCustomerGetSubscriptionList(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Customer().GetSubscriptionList(context.Background(), "CUSTOMER_1", &CustomerSubscriptionListInput{})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Get \"http://localhost:88888/api/v1/customers/CUSTOMER_1/subscriptions\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When no parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/subscriptions").
			MatchQuery("").
			MockResponse(mockCustomerSubscriptionListResponse)
		defer server.Close()

		result, err := server.Client().Customer().GetSubscriptionList(context.Background(), "CUSTOMER_1", &CustomerSubscriptionListInput{})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Subscriptions, qt.HasLen, 1)
		c.Assert(result.Meta.CurrentPage, qt.Equals, 1)
		c.Assert(result.Meta.TotalCount, qt.Equals, 1)
	})

	t.Run("When parameters are provided", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1/subscriptions").
			MatchQuery("per_page=10&page=1&plan_code=premium&status[]=active&status[]=terminated").
			MockResponse(mockCustomerSubscriptionListResponse)
		defer server.Close()

		perPage := 10
		page := 1
		result, err := server.Client().Customer().GetSubscriptionList(context.Background(), "CUSTOMER_1", &CustomerSubscriptionListInput{
			PerPage:  &perPage,
			Page:     &page,
			PlanCode: "premium",
			Status:   []SubscriptionStatus{SubscriptionStatusActive, SubscriptionStatusTerminated},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.Subscriptions, qt.HasLen, 1)
	})
}
