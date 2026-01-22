package lago

import "github.com/google/uuid"

type TaxProviderCustomerError struct {
	LagoCustomerID     uuid.UUID     `json:"lago_customer_id"`
	ExternalCustomerID string        `json:"external_customer_id"`
	TaxProvider        string        `json:"tax_provider"`
	TaxProviderCode    string        `json:"tax_provider_code"`
	ProviderError      ProviderError `json:"provider_error"`
}

type TaxProviderFeeError struct {
	TaxProviderCode    string        `json:"tax_provider_code"`
	LagoChargeID       *uuid.UUID    `json:"lago_charge_id,omitempty"`
	EventTransactionID *string       `json:"event_transaction_id,omitempty"`
	ProviderError      ProviderError `json:"provider_error,omitempty"`
}
