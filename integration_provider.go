package lago

import "github.com/google/uuid"

type IntegrationProviderError struct {
	LagoIntegrationID uuid.UUID       `json:"lago_integration_id"`
	Provider          IntegrationType `json:"provider"`
	ProviderCode      string          `json:"provider_code"`
	ProviderError     ProviderError   `json:"provider_error"`
}

type IntegrationCustomerError struct {
	LagoCustomerID         uuid.UUID       `json:"lago_customer_id"`
	ExternalCustomerID     string          `json:"external_customer_id"`
	AccountingProvider     IntegrationType `json:"accounting_provider"`
	AccountingProviderCode string          `json:"accounting_provider_code"`
	ProviderError          ProviderError   `json:"provider_error"`
}
