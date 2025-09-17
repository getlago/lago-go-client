package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Active     Status = "active"
	Terminated Status = "terminated"
)

type RecurringTransactionRuleInput struct {
	LagoID                           uuid.UUID                   `json:"lago_id,omitempty"`
	Interval                         string                      `json:"interval,omitempty"`
	Method                           string                      `json:"method,omitempty"`
	StartedAt                        *time.Time                  `json:"started_at,omitempty"`
	ExpirationAt                     *time.Time                  `json:"expiration_at,omitempty"`
	TargetOngoingBalance             string                      `json:"target_ongoing_balance,omitempty"`
	ThresholdCredits                 string                      `json:"threshold_credits,omitempty"`
	Trigger                          string                      `json:"trigger,omitempty"`
	PaidCredits                      string                      `json:"paid_credits,omitempty"`
	GrantedCredits                   string                      `json:"granted_credits,omitempty"`
	InvoiceRequiresSuccessfulPayment bool                        `json:"invoice_requires_successful_payment,omitempty"`
	TransactionMetadata              []WalletTransactionMetadata `json:"transaction_metadata,omitempty"`
	TransactionName                  string                      `json:"transaction_name,omitempty"`
}

type RecurringTransactionRuleResponse struct {
	LagoID                           uuid.UUID                   `json:"lago_id,omitempty"`
	Interval                         string                      `json:"interval,omitempty"`
	Method                           string                      `json:"method,omitempty"`
	StartedAt                        *time.Time                  `json:"started_at,omitempty"`
	ExpirationAt                     *time.Time                  `json:"expiration_at,omitempty"`
	Status                           Status                      `json:"status,omitempty"`
	TargetOngoingBalance             string                      `json:"target_ongoing_balance,omitempty"`
	ThresholdCredits                 string                      `json:"threshold_credits,omitempty"`
	Trigger                          string                      `json:"trigger,omitempty"`
	PaidCredits                      string                      `json:"paid_credits,omitempty"`
	GrantedCredits                   string                      `json:"granted_credits,omitempty"`
	CreatedAt                        time.Time                   `json:"created_at,omitempty"`
	InvoiceRequiresSuccessfulPayment bool                        `json:"invoice_requires_successful_payment,omitempty"`
	TransactionMetadata              []WalletTransactionMetadata `json:"transaction_metadata,omitempty"`
	TransactionName                  string                      `json:"transaction_name,omitempty"`
}

type WalletRequest struct {
	client *Client
}

type WalletParams struct {
	WalletInput *WalletInput `json:"wallet"`
}

type AppliesTo struct {
	FeeTypes            []string `json:"fee_types,omitempty"`
	BillableMetricCodes []string `json:"billable_metric_codes,omitempty"`
}

type WalletInput struct {
	RateAmount                       string                          `json:"rate_amount,omitempty"`
	Currency                         Currency                        `json:"currency,omitempty"`
	Name                             string                          `json:"name,omitempty"`
	PaidCredits                      string                          `json:"paid_credits,omitempty"`
	GrantedCredits                   string                          `json:"granted_credits,omitempty"`
	ExpirationAt                     *time.Time                      `json:"expiration_at,omitempty"`
	ExternalCustomerID               string                          `json:"external_customer_id,omitempty"`
	InvoiceRequiresSuccessfulPayment bool                            `json:"invoice_requires_successful_payment,omitempty"`
	TransactionMetadata              []WalletTransactionMetadata     `json:"transaction_metadata,omitempty"`
	TransactionName                  string                          `json:"transaction_name,omitempty"`
	RecurringTransactionRules        []RecurringTransactionRuleInput `json:"recurring_transaction_rules"`
	AppliesTo                        AppliesTo                       `json:"applies_to,omitempty"`
}

type WalletListInput struct {
	PerPage            *int   `json:"per_page,omitempty,string"`
	Page               *int   `json:"page,omitempty,string"`
	ExternalCustomerID string `json:"external_customer_id,omitempty"`
}

type WalletResult struct {
	Wallet  *Wallet  `json:"wallet,omitempty"`
	Wallets []Wallet `json:"wallets,omitempty"`
	Meta    Metadata `json:"meta,omitempty"`
}

type Wallet struct {
	LagoID                           uuid.UUID                          `json:"lago_id,omitempty"`
	LagoCustomerID                   uuid.UUID                          `json:"lago_customer_id,omitempty"`
	ExternalCustomerID               string                             `json:"external_customer_id,omitempty"`
	Status                           Status                             `json:"status,omitempty"`
	Currency                         Currency                           `json:"currency,omitempty"`
	Name                             string                             `json:"name,omitempty"`
	RateAmount                       string                             `json:"rate_amount,omitempty"`
	CreditsBalance                   string                             `json:"credits_balance,omitempty"`
	BalanceCents                     int                                `json:"balance_cents,omitempty"`
	ConsumedCredits                  string                             `json:"consumed_credits,omitempty"`
	InvoiceRequiresSuccessfulPayment bool                               `json:"invoice_requires_successful_payment,omitempty"`
	CreatedAt                        time.Time                          `json:"created_at,omitempty"`
	ExpirationAt                     time.Time                          `json:"expiration_at,omitempty"`
	LastBalanceSyncAt                time.Time                          `json:"last_balance_sync_at,omitempty"`
	LastConsumedCreditAt             time.Time                          `json:"last_consumed_credit_at,omitempty"`
	TerminatedAt                     time.Time                          `json:"terminated_at,omitempty"`
	RecurringTransactionRules        []RecurringTransactionRuleResponse `json:"recurring_transaction_rules,omitempty"`
	OngoingBalanceCents              int                                `json:"ongoing_balance_cents,omitempty"`
	OngoingUsageBalanceCents         int                                `json:"ongoing_usage_balance_cents,omitempty"`
	CreditsOngoingBalance            string                             `json:"credits_ongoing_balance,omitempty"`
	CreditsOngoingUsageBalance       string                             `json:"credits_ongoing_usage_balance,omitempty"`
	AppliesTo                        AppliesTo                          `json:"applies_to,omitempty"`
}

func (c *Client) Wallet() *WalletRequest {
	return &WalletRequest{
		client: c,
	}
}

func (bmr *WalletRequest) Get(ctx context.Context, walletID string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s", "wallets", walletID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
	}

	result, err := bmr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}

func (bmr *WalletRequest) GetList(ctx context.Context, walletListInput *WalletListInput) (*WalletResult, *Error) {
	jsonQueryParams, err := json.Marshal(walletListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        "wallets",
		QueryParams: queryParams,
		Result:      &WalletResult{},
	}

	result, clientErr := bmr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult, nil
}

func (bmr *WalletRequest) Create(ctx context.Context, walletInput *WalletInput) (*Wallet, *Error) {
	walletParams := &WalletParams{
		WalletInput: walletInput,
	}

	clientRequest := &ClientRequest{
		Path:   "wallets",
		Result: &WalletResult{},
		Body:   walletParams,
	}

	result, err := bmr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}

func (bmr *WalletRequest) Update(ctx context.Context, walletInput *WalletInput, walletID string) (*Wallet, *Error) {
	walletParams := &WalletParams{
		WalletInput: walletInput,
	}

	subPath := fmt.Sprintf("%s/%s", "wallets", walletID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
		Body:   walletParams,
	}

	result, err := bmr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}

func (bmr *WalletRequest) Delete(ctx context.Context, walletID string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s", "wallets", walletID)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
	}

	result, err := bmr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}
