package lago

import "time"

type PaymentMethodInput struct {
	PaymentMethodType string `json:"payment_method_type,omitempty"`
	PaymentMethodID   string `json:"payment_method_id,omitempty"`
}

type PaymentMethod struct {
	LagoID              string    `json:"lago_id,omitempty"`
	IsDefault           bool      `json:"is_default"`
	PaymentProviderCode string    `json:"payment_provider_code,omitempty"`
	PaymentProviderName string    `json:"payment_provider_name,omitempty"`
	PaymentProviderType string    `json:"payment_provider_type,omitempty"`
	ProviderMethodID    string    `json:"provider_method_id,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
}
