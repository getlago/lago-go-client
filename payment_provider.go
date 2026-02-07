package lago

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ProviderError map[string]any

type PaymentProviderError struct {
	LagoPaymentProviderID uuid.UUID               `json:"lago_payment_provider_id"`
	PaymentProviderName   string                  `json:"payment_provider_name"`
	PaymentProviderCode   string                  `json:"payment_provider_code"`
	Source                CustomerPaymentProvider `json:"source"`
	Action                string                  `json:"action"`
	ProviderError         ProviderError           `json:"provider_error"`
}

type PaymentProviderCustomerError struct {
	LagoCustomerID      uuid.UUID               `json:"lago_customer_id"`
	ExternalCustomerID  string                  `json:"external_customer_id"`
	PaymentProvider     CustomerPaymentProvider `json:"payment_provider"`
	PaymentProviderCode string                  `json:"payment_provider_code"`
	ProviderError       ProviderError           `json:"provider_error"`
}

type PaymentProviderCreditNoteRefundError struct {
	LagoCreditNoteID    uuid.UUID               `json:"lago_credit_note_id"`
	LagoCustomerID      uuid.UUID               `json:"lago_customer_id"`
	ExternalCustomerID  string                  `json:"external_customer_id"`
	ProviderCustomerID  string                  `json:"provider_customer_id"`
	PaymentProvider     CustomerPaymentProvider `json:"payment_provider"`
	PaymentProviderCode string                  `json:"payment_provider_code"`
	ProviderError       ProviderError           `json:"provider_error"`
}

type PaymentProviderInvoiceError struct {
	LagoInvoiceID       uuid.UUID               `json:"lago_invoice_id"`
	LagoCustomerID      uuid.UUID               `json:"lago_customer_id"`
	ExternalCustomerID  string                  `json:"external_customer_id"`
	ProviderCustomerID  string                  `json:"provider_customer_id"`
	PaymentProvider     CustomerPaymentProvider `json:"payment_provider"`
	PaymentProviderCode string                  `json:"payment_provider_code"`
	ProviderError       ProviderError           `json:"provider_error"`
}

type PaymentProviderPaymentRequestError struct {
	LagoPaymentRequestID uuid.UUID               `json:"lago_payment_request_id"`
	LagoInvoiceIDs       []uuid.UUID             `json:"lago_invoice_ids"`
	LagoCustomerID       uuid.UUID               `json:"lago_customer_id"`
	ExternalCustomerID   string                  `json:"external_customer_id"`
	ProviderCustomerID   string                  `json:"provider_customer_id"`
	PaymentProvider      CustomerPaymentProvider `json:"payment_provider"`
	PaymentProviderCode  string                  `json:"payment_provider_code"`
	ProviderError        ProviderError           `json:"provider_error"`
}

type PaymentProviderWalletTransactionError struct {
	LagoWalletTransactionID uuid.UUID               `json:"lago_wallet_transaction_id"`
	LagoCustomerID          uuid.UUID               `json:"lago_customer_id"`
	ExternalCustomerID      string                  `json:"external_customer_id"`
	ProviderCustomerID      string                  `json:"provider_customer_id"`
	PaymentProvider         CustomerPaymentProvider `json:"payment_provider"`
	PaymentProviderCode     string                  `json:"payment_provider_code"`
	ProviderError           ProviderError           `json:"provider_error"`
}

func (p *ProviderError) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as an object first
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err == nil {
		*p = obj
		return nil
	}

	// Then, try to unmarshal as a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*p = map[string]any{"message": str}
	return nil
}
