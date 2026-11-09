package lago_test

import (
	"context"
	"encoding/json"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"
)

func TestPaymentTerm_Marshal(t *testing.T) {
	t.Run("When only the relevant fields are set", func(t *testing.T) {
		c := qt.New(t)

		body, err := json.Marshal(PaymentTerm{TermType: PaymentTermTypeNet, Days: Ptr(30)})
		c.Assert(err, qt.IsNil)
		c.Assert(string(body), qt.JSONEquals, map[string]any{"term_type": "net", "days": 30})
	})

	t.Run("When a field is set to zero", func(t *testing.T) {
		c := qt.New(t)

		body, err := json.Marshal(PaymentTerm{TermType: PaymentTermTypeNet, Days: Ptr(0)})
		c.Assert(err, qt.IsNil)
		c.Assert(string(body), qt.JSONEquals, map[string]any{"term_type": "net", "days": 0})
	})

	t.Run("When the term is day_of_month", func(t *testing.T) {
		c := qt.New(t)

		body, err := json.Marshal(PaymentTerm{TermType: PaymentTermTypeDayOfMonth, DayOfMonth: Ptr(15), MonthOffset: Ptr(2)})
		c.Assert(err, qt.IsNil)
		c.Assert(string(body), qt.JSONEquals, map[string]any{"term_type": "day_of_month", "day_of_month": 15, "month_offset": 2})
	})
}

func TestPaymentTerm_CustomerCreate(t *testing.T) {
	t.Run("When only net_payment_term is set, payment_term is not sent", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers").
			MatchJSONBody(`{
				"customer": {
					"external_id": "CUSTOMER_1",
					"net_payment_term": 30,
					"billing_configuration": {
						"subscription_invoice_issuing_date_anchor": null,
						"subscription_invoice_issuing_date_adjustment": null
					},
					"shipping_address": {}
				}
			}`).
			MockResponse(mockCustomerCreateResponse)
		defer server.Close()

		result, err := server.Client().Customer().Update(context.Background(), &CustomerInput{
			ExternalID:     "CUSTOMER_1",
			NetPaymentTerm: 30,
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
	})

	t.Run("When payment_term is set", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/customers").
			MatchJSONBody(`{
				"customer": {
					"external_id": "CUSTOMER_1",
					"payment_term": {"term_type": "net_end_of_month", "days": 10},
					"billing_configuration": {
						"subscription_invoice_issuing_date_anchor": null,
						"subscription_invoice_issuing_date_adjustment": null
					},
					"shipping_address": {}
				}
			}`).
			MockResponse(mockCustomerCreateResponse)
		defer server.Close()

		result, err := server.Client().Customer().Create(context.Background(), &CustomerInput{
			ExternalID:  "CUSTOMER_1",
			PaymentTerm: &PaymentTerm{TermType: PaymentTermTypeNetEndOfMonth, Days: Ptr(10)},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
	})
}

func TestPaymentTerm_BillingEntity(t *testing.T) {
	t.Run("When creating with payment_term", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/billing_entities").
			MatchJSONBody(`{"billing_entity":{"name":"Acme FR","code":"acme_fr","payment_term":{"term_type":"due_on_receipt"},"billing_configuration":{}}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		_, err := server.Client().BillingEntity().Create(context.Background(), &BillingEntityCreateInput{
			Name:        "Acme FR",
			Code:        "acme_fr",
			PaymentTerm: &PaymentTerm{TermType: PaymentTermTypeDueOnReceipt},
		})
		c.Assert(err == nil, qt.IsTrue)
	})

	t.Run("When updating with only net_payment_term, payment_term is not sent", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MatchJSONBody(`{"billing_entity":{"net_payment_term":15,"billing_configuration":{}}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		_, err := server.Client().BillingEntity().Update(context.Background(), "acme_fr", &BillingEntityUpdateInput{
			NetPaymentTerm: 15,
		})
		c.Assert(err == nil, qt.IsTrue)
	})

	t.Run("When updating with payment_term", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MatchJSONBody(`{"billing_entity":{"payment_term":{"term_type":"end_of_month"},"billing_configuration":{}}}`).
			MockResponse(mockBillingEntityResponse)
		defer server.Close()

		_, err := server.Client().BillingEntity().Update(context.Background(), "acme_fr", &BillingEntityUpdateInput{
			PaymentTerm: &PaymentTerm{TermType: PaymentTermTypeEndOfMonth},
		})
		c.Assert(err == nil, qt.IsTrue)
	})
}

func TestPaymentTerm_SubscriptionUpdate(t *testing.T) {
	t.Run("When payment_term is not set, it is not sent", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/subscriptions/my_sub_1").
			MatchJSONBody(`{"subscription":{"external_id":"my_sub_1","name":"Repository A"}}`).
			MockResponse(mockSubscriptionResponse)
		defer server.Close()

		subscription, err := server.Client().Subscription().Update(context.Background(), &SubscriptionInput{
			ExternalID: "my_sub_1",
			Name:       "Repository A",
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(subscription, qt.IsNotNil)
	})

	t.Run("When payment_term is set", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/subscriptions/my_sub_1").
			MatchJSONBody(`{"subscription":{"external_id":"my_sub_1","name":"Repository A","payment_term":{"term_type":"day_of_month","day_of_month":5,"month_offset":1}}}`).
			MockResponse(mockSubscriptionResponse)
		defer server.Close()

		subscription, err := server.Client().Subscription().Update(context.Background(), &SubscriptionInput{
			ExternalID: "my_sub_1",
			Name:       "Repository A",
			PaymentTerm: &PaymentTerm{
				TermType:    PaymentTermTypeDayOfMonth,
				DayOfMonth:  Ptr(5),
				MonthOffset: Ptr(1),
			},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(subscription, qt.IsNotNil)
	})
}

func TestPaymentTerm_SubscriptionCreate(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("POST").
		MatchPath("/api/v1/subscriptions").
		MatchJSONBody(`{"subscription":{"external_id":"my_sub_1","name":"Repository A","payment_term":{"term_type":"net","days":30}}}`).
		MockResponse(mockSubscriptionResponse)
	defer server.Close()

	subscription, err := server.Client().Subscription().Create(context.Background(), &SubscriptionInput{
		ExternalID:  "my_sub_1",
		Name:        "Repository A",
		PaymentTerm: &PaymentTerm{TermType: PaymentTermTypeNet, Days: Ptr(30)},
	})
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(subscription, qt.IsNotNil)
}

func TestPaymentTerm_InvoicePreview(t *testing.T) {
	t.Run("When payment_term is not set, it is not sent", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/invoices/preview").
			MatchJSONBody(`{"plan_code":"plan_1"}`).
			MockResponse(map[string]any{"invoice": mockInvoice})
		defer server.Close()

		_, err := server.Client().Invoice().Preview(context.Background(), &InvoicePreviewInput{PlanCode: "plan_1"})
		c.Assert(err == nil, qt.IsTrue)
	})

	t.Run("When payment_term is set at the top level", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/invoices/preview").
			MatchJSONBody(`{"plan_code":"plan_1","payment_term":{"term_type":"days_end_of_month","days":15}}`).
			MockResponse(map[string]any{"invoice": mockInvoice})
		defer server.Close()

		_, err := server.Client().Invoice().Preview(context.Background(), &InvoicePreviewInput{
			PlanCode:    "plan_1",
			PaymentTerm: &PaymentTerm{TermType: PaymentTermTypeDaysEndOfMonth, Days: Ptr(15)},
		})
		c.Assert(err == nil, qt.IsTrue)
	})
}

func TestPaymentTerm_InvoiceOneOffCreate(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("POST").
		MatchPath("/api/v1/invoices").
		MatchJSONBody(`{
			"invoice": {
				"external_customer_id": "CUSTOMER_1",
				"currency": "USD",
				"fees": [{"add_on_code": "setup_fee", "units": 1}],
				"payment_term": {"term_type": "net", "days": 7}
			}
		}`).
		MockResponse(map[string]any{"invoice": mockInvoice})
	defer server.Close()

	result, err := server.Client().Invoice().Create(context.Background(), &InvoiceOneOffInput{
		ExternalCustomerId: "CUSTOMER_1",
		Currency:           "USD",
		Fees:               []InvoiceFeesInput{{AddOnCode: "setup_fee", Units: 1}},
		PaymentTerm:        &PaymentTerm{TermType: PaymentTermTypeNet, Days: Ptr(7)},
	})
	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result, qt.IsNotNil)
}

func TestPaymentTerm_WalletRequests(t *testing.T) {
	t.Run("When creating with payment_term", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/wallets").
			MatchJSONBody(`{
				"wallet": {
					"external_customer_id": "12345",
					"currency": "USD",
					"recurring_transaction_rules": null,
					"applies_to": {},
					"payment_term": {"term_type": "net", "days": 14}
				}
			}`).
			MockResponse(mockWalletResponse)
		defer server.Close()

		result, err := server.Client().Wallet().Create(context.Background(), &WalletInput{
			ExternalCustomerID: "12345",
			Currency:           "USD",
			PaymentTerm:        &PaymentTerm{TermType: PaymentTermTypeNet, Days: Ptr(14)},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
	})

	t.Run("When updating without payment_term, it is not sent", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0").
			MatchJSONBody(`{"wallet":{"name":"renamed","recurring_transaction_rules":null,"applies_to":{}}}`).
			MockResponse(mockWalletResponse)
		defer server.Close()

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{Name: "renamed"}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
	})

	t.Run("When updating with payment_term", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0").
			MatchJSONBody(`{"wallet":{"recurring_transaction_rules":null,"applies_to":{},"payment_term":{"term_type":"due_on_receipt"}}}`).
			MockResponse(mockWalletResponse)
		defer server.Close()

		result, err := server.Client().Wallet().Update(context.Background(), &WalletInput{
			PaymentTerm: &PaymentTerm{TermType: PaymentTermTypeDueOnReceipt},
		}, "b1b2c3d4-e5f6-7890-1234-56789abcdef0")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.IsNotNil)
	})
}

func TestPaymentTerm_ResponseParsing(t *testing.T) {
	paymentTerm := map[string]any{"term_type": "net", "days": 30}
	expected := &PaymentTerm{TermType: PaymentTermTypeNet, Days: Ptr(30)}

	t.Run("Customer", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1").
			MockResponse(map[string]any{"customer": map[string]any{
				"external_id":      "CUSTOMER_1",
				"net_payment_term": 30,
				"payment_term":     paymentTerm,
			}})
		defer server.Close()

		result, err := server.Client().Customer().Get(context.Background(), "CUSTOMER_1")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.NetPaymentTerm, qt.Equals, 30)
		c.Assert(result.PaymentTerm, qt.DeepEquals, expected)
	})

	t.Run("Customer with a null payment_term", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/customers/CUSTOMER_1").
			MockResponse(map[string]any{"customer": map[string]any{
				"external_id":      "CUSTOMER_1",
				"net_payment_term": nil,
				"payment_term":     nil,
			}})
		defer server.Close()

		result, err := server.Client().Customer().Get(context.Background(), "CUSTOMER_1")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.NetPaymentTerm, qt.Equals, 0)
		c.Assert(result.PaymentTerm, qt.IsNil)
	})

	t.Run("BillingEntity", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/billing_entities/acme_fr").
			MockResponse(map[string]any{"billing_entity": map[string]any{
				"code":             "acme_fr",
				"net_payment_term": nil,
				"payment_term":     map[string]any{"term_type": "end_of_month"},
			}})
		defer server.Close()

		result, err := server.Client().BillingEntity().Get(context.Background(), "acme_fr")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.NetPaymentTerm, qt.Equals, 0)
		c.Assert(result.PaymentTerm, qt.DeepEquals, &PaymentTerm{TermType: PaymentTermTypeEndOfMonth})
	})

	t.Run("Organization", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("PUT").
			MatchPath("/api/v1/organizations").
			MockResponse(map[string]any{"organization": map[string]any{
				"name":             "Lago",
				"net_payment_term": 30,
				"payment_term":     paymentTerm,
			}})
		defer server.Close()

		result, err := server.Client().Organization().Update(context.Background(), &OrganizationInput{NetPaymentTerm: 30})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentTerm, qt.DeepEquals, expected)
	})

	t.Run("Subscription", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/subscriptions/my_sub_1").
			MockResponse(map[string]any{"subscription": map[string]any{
				"external_id":  "my_sub_1",
				"payment_term": map[string]any{"term_type": "day_of_month", "day_of_month": 5, "month_offset": 1},
			}})
		defer server.Close()

		result, err := server.Client().Subscription().Get(context.Background(), "my_sub_1")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentTerm, qt.DeepEquals, &PaymentTerm{
			TermType:    PaymentTermTypeDayOfMonth,
			DayOfMonth:  Ptr(5),
			MonthOffset: Ptr(1),
		})
	})

	t.Run("Invoice", func(t *testing.T) {
		c := qt.New(t)

		invoice := map[string]any{}
		for k, v := range mockInvoice {
			invoice[k] = v
		}
		invoice["net_payment_term"] = nil
		invoice["payment_term"] = map[string]any{"term_type": "due_on_receipt"}

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/invoices/1a901a90-1a90-1a90-1a90-1a901a901a90").
			MockResponse(map[string]any{"invoice": invoice})
		defer server.Close()

		result, err := server.Client().Invoice().Get(context.Background(), "1a901a90-1a90-1a90-1a90-1a901a901a90")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.NetPaymentTerm, qt.Equals, 0)
		c.Assert(result.PaymentTerm, qt.DeepEquals, &PaymentTerm{TermType: PaymentTermTypeDueOnReceipt})
	})

	t.Run("Wallet", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("GET").
			MatchPath("/api/v1/wallets/b1b2c3d4-e5f6-7890-1234-56789abcdef0").
			MockResponse(map[string]any{"wallet": map[string]any{
				"lago_id":      "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"payment_term": paymentTerm,
			}})
		defer server.Close()

		result, err := server.Client().Wallet().Get(context.Background(), "b1b2c3d4-e5f6-7890-1234-56789abcdef0")
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.PaymentTerm, qt.DeepEquals, expected)
	})
}
