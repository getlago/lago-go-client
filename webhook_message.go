package lago

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

var WebhookObjectTypeMapping = map[string]func() any{
	"credit_note":                            func() any { return &CreditNote{} },
	"customer":                               func() any { return &Customer{} },
	"feature":                                func() any { return &Feature{} },
	"fee":                                    func() any { return &Fee{} },
	"invoice":                                func() any { return &Invoice{} },
	"payment":                                func() any { return &Payment{} },
	"payment_provider_customer_checkout_url": func() any { return &CustomerCheckoutUrl{} },
	"payment_request":                        func() any { return &PaymentRequest{} },
	"plan":                                   func() any { return &Plan{} },
	"subscription":                           func() any { return &Subscription{} },
	"triggered_alert":                        func() any { return &TriggeredAlert{} },
	"wallet":                                 func() any { return &Wallet{} },
	"wallet_transaction":                     func() any { return &WalletTransaction{} },
}

type WebhookMessage struct {
	WebhookType    string    `json:"webhook_type"`
	ObjectType     string    `json:"object_type"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Object         any
}

func ParseWebhook(data []byte) (*WebhookMessage, error) {
	// Unmarshall the common fields
	var base struct {
		WebhookType    string    `json:"webhook_type"`
		ObjectType     string    `json:"object_type"`
		OrganizationID uuid.UUID `json:"organization_id"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	constructor, ok := WebhookObjectTypeMapping[base.ObjectType]
	if !ok {
		return nil, fmt.Errorf("unknown object_type: %s", base.ObjectType)
	}

	// Unmarshall the specific object field
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	objData, ok := raw[base.ObjectType]
	if !ok {
		return nil, fmt.Errorf("missing %s attribute", base.ObjectType)
	}

	obj := constructor()
	if err := json.Unmarshal(objData, obj); err != nil {
		return nil, err
	}

	return &WebhookMessage{
		WebhookType:    base.WebhookType,
		ObjectType:     base.ObjectType,
		OrganizationID: base.OrganizationID,
		Object:         obj,
	}, nil
}
