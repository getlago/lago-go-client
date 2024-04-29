package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WalletTransactionStatus string

const (
	WalletTransactionStatusPending WalletTransactionStatus = "pending"
	WalletTransactionStatusSettled WalletTransactionStatus = "settled"
)

type TransactionStatus string

const (
	Purchased TransactionStatus = "purchased"
	Granted   TransactionStatus = "granted"
	Voided    TransactionStatus = "voided"
)

type TransactionType string

const (
	Outbound TransactionType = "outbound"
	Inbound  TransactionType = "inbound"
)

type WalletTransactionRequest struct {
	client *Client
}

type WalletTransactionListInput struct {
	PerPage           int                     `json:"per_page,omitempty,string"`
	Page              int                     `json:"page,omitempty,string"`
	WalletID          string                  `json:"wallet_id,omitempty"`
	Status            WalletTransactionStatus `json:"status,omitempty"`
	TransactionStatus TransactionStatus       `json:"transaction_status,omitempty"`
	TransactionType   TransactionType         `json:"transaction_type,omitempty"`
}

type WalletTransactionParams struct {
	WalletTransactionInput *WalletTransactionInput `json:"wallet_transaction"`
}

type WalletTransactionInput struct {
	WalletID       string `json:"wallet_id,omitempty"`
	PaidCredits    string `json:"paid_credits,omitempty"`
	GrantedCredits string `json:"granted_credits,omitempty"`
	VoidedCredits  string `json:"voided_credits,omitempty"`
}

type WalletTransactionResult struct {
	WalletTransactions []WalletTransaction `json:"wallet_transactions,omitempty"`
	Meta               Metadata            `json:"meta,omitempty"`
}

type WalletTransaction struct {
	LagoID          uuid.UUID               `json:"lago_id,omitempty"`
	LagoWalletID    uuid.UUID               `json:"lago_wallet_id,omitempty"`
	Status          WalletTransactionStatus `json:"status,omitempty"`
	TransactionType TransactionType         `json:"transaction_type,omitempty"`
	Amount          string                  `json:"amount,omitempty"`
	CreditAmount    string                  ` json:"credit_amount,omitempty"`
	CreatedAt       time.Time               `json:"created_at,omitempty"`
	SettledAt       time.Time               `json:"settled_at,omitempty"`
}

func (c *Client) WalletTransaction() *WalletTransactionRequest {
	return &WalletTransactionRequest{
		client: c,
	}
}

func (wtr *WalletTransactionRequest) Create(ctx context.Context, walletTransactionInput *WalletTransactionInput) (*WalletTransactionResult, *Error) {
	walletTransactionParams := &WalletTransactionParams{
		WalletTransactionInput: walletTransactionInput,
	}

	clientRequest := &ClientRequest{
		Path:   "wallet_transactions",
		Result: &WalletTransactionResult{},
		Body:   walletTransactionParams,
	}
	result, err := wtr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletTransactionResult, ok := result.(*WalletTransactionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletTransactionResult, nil
}

func (wtr *WalletTransactionRequest) GetList(ctx context.Context, walletTransactionListInput *WalletTransactionListInput) (*WalletTransactionResult, *Error) {
	jsonQueryParams, err := json.Marshal(walletTransactionListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	subPath := fmt.Sprintf("%s/%s/%s", "wallets", walletTransactionListInput.WalletID, "wallet_transactions")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &WalletTransactionResult{},
	}

	result, clientErr := wtr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	walletTransactionResult, ok := result.(*WalletTransactionResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletTransactionResult, nil
}
