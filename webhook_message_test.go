package lago

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"
)

func TestWebhookMessage_ParseWebhook_Fee(t *testing.T) {
	c := qt.New(t)

	jsonData := []byte(`{
		"webhook_type": "fee.created",
		"object_type": "fee",
		"organization_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"fee": {
	    "lago_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
	    "billing_entity_code": "test-company",
	    "lago_charge_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
	    "lago_charge_filter_id": "1a2d9c8d-5875-4688-9854-5ccfd414bc5e",
	    "lago_invoice_id": "1a2d9c8d-5875-4688-9854-5ccfd414bc5e",
	    "lago_true_up_fee_id": "1a2d9c8d-5875-4688-9854-5ccfd414bc5e",
	    "lago_true_up_parent_fee_id": null,
	    "invoice_display_name": "fee_invoice_display_name",
	    "item": {
	      "type": "charge",
	      "code": "fee_code",
	      "name": "Fee Code",
	      "invoice_display_name": "charge_invoice_display_name",
	      "charge_filter_invoice_display_name": "charge_filter_invoice_display_name",
	      "grouped_by": {
	        "agent_name": "aragorn"
	      }
	    },
	    "amount_cents": 120,
	    "amount_currency": "EUR",
	    "taxes_amount_cents": 20,
	    "taxes_rate": 20.0,
	    "total_aggregated_units": "10.0",
	    "total_amount_cents": 140,
	    "total_amount_currency": "EUR",
	    "units": "10.0",
	    "events_count": 10,
	    "pricing_unit_details": {
	      "lago_pricing_unit_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
	      "pricing_unit_code": "unit_code",
	      "short_name": "UC",
	      "amount_cents": 120,
	      "precise_amount_cents": "120.00",
	      "unit_amount_cents": 12,
	      "precise_unit_amount": "12.00",
	      "conversion_rate": 1.0
	    },
	    "applied_taxes": []
		}
	}`)

	msg, err := ParseWebhook(jsonData)

	c.Assert(err, qt.IsNil)
	c.Assert(msg.WebhookType, qt.Equals, "fee.created")
	c.Assert(msg.ObjectType, qt.Equals, "fee")
	c.Assert(msg.OrganizationID, qt.Equals, uuid.MustParse("1a901a90-1a90-1a90-1a90-1a901a901a90"))

	c.Assert(msg.Object, qt.IsNotNil)

	object, ok := msg.Object.(*Fee)
	c.Assert(ok, qt.IsTrue)
	c.Assert(object.LagoID, qt.Equals, uuid.MustParse("1a901a90-1a90-1a90-1a90-1a901a901a90"))
	c.Assert(object.AmountCents, qt.Equals, 120)
}

func TestWebhookMessage_ParseWebhook_CustomerAccountingProviderError(t *testing.T) {
	c := qt.New(t)

	jsonData := []byte(`{
		"webhook_type": "customer.accounting_provider_error",
		"object_type": "accounting_provider_customer_error",
		"organization_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"accounting_provider_customer_error": {
  		"lago_customer_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
	    "external_customer_id": "customer_1234",
	    "accounting_provider": "netsuite",
	    "accounting_provider_code": "Netsuite Prod",
	    "provider_error": {
     		"error_message": "message",
        "error_code": "code"
			}
		}
	}`)

	msg, err := ParseWebhook(jsonData)

	c.Assert(err, qt.IsNil)
	c.Assert(msg.WebhookType, qt.Equals, "customer.accounting_provider_error")
	c.Assert(msg.ObjectType, qt.Equals, "accounting_provider_customer_error")
	c.Assert(msg.OrganizationID, qt.Equals, uuid.MustParse("1a901a90-1a90-1a90-1a90-1a901a901a90"))

	c.Assert(msg.Object, qt.IsNotNil)

	object, ok := msg.Object.(*IntegrationCustomerError)
	c.Assert(ok, qt.IsTrue)
	c.Assert(object.LagoCustomerID, qt.Equals, uuid.MustParse("1a901a90-1a90-1a90-1a90-1a901a901a90"))

	c.Assert(object.ProviderError["error_code"], qt.Equals, "code")
	c.Assert(object.ProviderError["error_message"], qt.Equals, "message")
}

func TestWebhookMessage_ParseWebhook_CustomerAccountingProviderError_WithSingleMessage(t *testing.T) {
	c := qt.New(t)

	jsonData := []byte(`{
		"webhook_type": "customer.accounting_provider_error",
		"object_type": "accounting_provider_customer_error",
		"organization_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
		"accounting_provider_customer_error": {
  		"lago_customer_id": "1a901a90-1a90-1a90-1a90-1a901a901a90",
	    "external_customer_id": "customer_1234",
	    "accounting_provider": "netsuite",
	    "accounting_provider_code": "Netsuite Prod",
	    "provider_error": "error_message"
		}
	}`)

	msg, err := ParseWebhook(jsonData)

	c.Assert(err, qt.IsNil)
	c.Assert(msg.WebhookType, qt.Equals, "customer.accounting_provider_error")
	c.Assert(msg.ObjectType, qt.Equals, "accounting_provider_customer_error")
	c.Assert(msg.OrganizationID, qt.Equals, uuid.MustParse("1a901a90-1a90-1a90-1a90-1a901a901a90"))

	c.Assert(msg.Object, qt.IsNotNil)

	object, ok := msg.Object.(*IntegrationCustomerError)
	c.Assert(ok, qt.IsTrue)
	c.Assert(object.LagoCustomerID, qt.Equals, uuid.MustParse("1a901a90-1a90-1a90-1a90-1a901a901a90"))

	c.Assert(object.ProviderError["message"], qt.Equals, "error_message")
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
