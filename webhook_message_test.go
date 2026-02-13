package lago

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

var tests = []struct {
	fixture string
	test    func(object any) bool
}{
	{
		fixture: "alert_triggered",
		test: func(object any) bool {
			_, ok := object.(*TriggeredAlert)
			return ok
		},
	},
	{
		fixture: "credit_note_created",
		test: func(object any) bool {
			_, ok := object.(*CreditNote)
			return ok
		},
	},
	{
		fixture: "credit_note_generated",
		test: func(object any) bool {
			_, ok := object.(*CreditNote)
			return ok
		},
	},
	{
		fixture: "customer_accounting_provider_created",
		test: func(object any) bool {
			_, ok := object.(*Customer)
			return ok
		},
	},
	{
		fixture: "customer_checkout_url_generated",
		test: func(object any) bool {
			_, ok := object.(*CustomerCheckoutUrl)
			return ok
		},
	},
	{
		fixture: "customer_created",
		test: func(object any) bool {
			_, ok := object.(*Customer)
			return ok
		},
	},
	{
		fixture: "customer_crm_provider_created",
		test: func(object any) bool {
			_, ok := object.(*Customer)
			return ok
		},
	},
	{
		fixture: "customer_payment_provider_created",
		test: func(object any) bool {
			_, ok := object.(*Customer)
			return ok
		},
	},
	{
		fixture: "customer_updated",
		test: func(object any) bool {
			_, ok := object.(*Customer)
			return ok
		},
	},
	{
		fixture: "customer_vies_check",
		test: func(object any) bool {
			_, ok := object.(*Customer)
			return ok
		},
	},
	{
		fixture: "feature_created",
		test: func(object any) bool {
			_, ok := object.(*Feature)
			return ok
		},
	},
	{
		fixture: "feature_deleted",
		test: func(object any) bool {
			_, ok := object.(*Feature)
			return ok
		},
	},
	{
		fixture: "feature_updated",
		test: func(object any) bool {
			_, ok := object.(*Feature)
			return ok
		},
	},
	{
		fixture: "fee_created",
		test: func(object any) bool {
			_, ok := object.(*Fee)
			return ok
		},
	},
	{
		fixture: "invoice_add_on_added",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_created",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_drafted",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_generated",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_one_off_created",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_paid_credit_added",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_payment_overdue",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_payment_status_updated",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_resynced",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "invoice_voided",
		test: func(object any) bool {
			_, ok := object.(*Invoice)
			return ok
		},
	},
	{
		fixture: "payment_request_created",
		test: func(object any) bool {
			_, ok := object.(*PaymentRequest)
			return ok
		},
	},
	{
		fixture: "payment_request_payment_status_updated",
		test: func(object any) bool {
			_, ok := object.(*PaymentRequest)
			return ok
		},
	},
	{
		fixture: "payment_requires_action",
		test: func(object any) bool {
			_, ok := object.(*Payment)
			return ok
		},
	},
	{
		fixture: "plan_created",
		test: func(object any) bool {
			_, ok := object.(*Plan)
			return ok
		},
	},
	{
		fixture: "plan_deleted",
		test: func(object any) bool {
			_, ok := object.(*Plan)
			return ok
		},
	},
	{
		fixture: "plan_updated",
		test: func(object any) bool {
			_, ok := object.(*Plan)
			return ok
		},
	},
	{
		fixture: "subscription_started",
		test: func(object any) bool {
			_, ok := object.(*Subscription)
			return ok
		},
	},
	{
		fixture: "subscription_terminated",
		test: func(object any) bool {
			_, ok := object.(*Subscription)
			return ok
		},
	},
	{
		fixture: "subscription_termination_alert",
		test: func(object any) bool {
			_, ok := object.(*Subscription)
			return ok
		},
	},
	{
		fixture: "subscription_trial_ended",
		test: func(object any) bool {
			_, ok := object.(*Subscription)
			return ok
		},
	},
	{
		fixture: "subscription_updated",
		test: func(object any) bool {
			_, ok := object.(*Subscription)
			return ok
		},
	},
	{
		fixture: "subscription_usage_threshold_reached",
		test: func(object any) bool {
			_, ok := object.(*Subscription)
			return ok
		},
	},
	{
		fixture: "wallet_depleted_ongoing_balance",
		test: func(object any) bool {
			_, ok := object.(*Wallet)
			return ok
		},
	},
	{
		fixture: "wallet_transaction_created",
		test: func(object any) bool {
			_, ok := object.(*WalletTransaction)
			return ok
		},
	},
	{
		fixture: "wallet_transaction_updated",
		test: func(object any) bool {
			_, ok := object.(*WalletTransaction)
			return ok
		},
	},
}

func TestWebhookMessage_ParseWebhook_AllFixtures(t *testing.T) {
	for _, tt := range tests {
		t.Run("Testing_webhook_"+tt.fixture, func(t *testing.T) {
			c := qt.New(t)

			// Load the fixture file
			jsonData, err := os.ReadFile("testing/fixtures/webhooks/" + tt.fixture + ".json")
			c.Assert(err, qt.IsNil)

			// Parse the webhook
			msg, err := ParseWebhook(jsonData)
			c.Assert(err, qt.IsNil)

			// Check the test function returns true for the parsed object
			c.Assert(tt.test(msg.Object), qt.IsTrue)
		})
	}
}

func TestWebhookMessage_ParseWebhook_MissingObjectField(t *testing.T) {
	c := qt.New(t)

	// JSON where the object_type is "fee" but no "fee" field exists
	jsonData := []byte(`{
		"webhook_type": "fee.created",
		"object_type": "fee",
		"organization_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	}`)

	_, err := ParseWebhook(jsonData)

	// Should not error, but Object should be zero value
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "missing fee attribute")
}

func TestWebhookMessage_UnmarshalJSON_InvalidJSON(t *testing.T) {
	c := qt.New(t)

	jsonData := []byte(`{invalid json}`)

	_, err := ParseWebhook(jsonData)

	c.Assert(err, qt.IsNotNil)
}

func TestWebhookMessage_ParseWebhook_InvalidObjectJSON(t *testing.T) {
	c := qt.New(t)

	// The "fee" field contains invalid data for TestFee (lago_id should be UUID)
	jsonData := []byte(`{
		"webhook_type": "fee.created",
		"object_type": "fee",
		"organization_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"fee": {
			"lago_id": "not-a-valid-uuid",
			"amount_cents": 1000
		}
	}`)

	_, err := ParseWebhook(jsonData)

	c.Assert(err, qt.IsNotNil)
}
