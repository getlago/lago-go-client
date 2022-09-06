package lago

import (
	"time"

	"github.com/google/uuid"
)

type WalletStatus string

const (
	WalletStatusPending WalletStatus = "pending"
	WalletStatusSettled WalletStatus = "settled"
)

type TransactionType string

const (
	Outbound TransactionType = "outbound"
	Inbound  TransactionType = "inbound"
)

type WalletTransactionRequest struct {
	client *Client
}

type WalletTransactionParams struct {
	WalletTransactionInput *WalletTransactionInput
}

type WalletTransactionInput struct {
	WalletId       string `json:"wallet_id,omitempty"`
	PaidCredits    string `json:"paid_credits,omitempty"`
	GrantedCredits string `json:"granted_credits,omitempty"`
}

type WalletTransactionResult struct {
	WalletTransactions []WalletTransaction `json:"wallet_transactions,omitempty"`
}

type WalletTransaction struct {
	LagoID          uuid.UUID       `json:"lago_id,omitempty"`
	LagoWalletID    uuid.UUID       `json:"lago_wallet_id,omitempty"`
	Status          WalletStatus    `json:"status,omitempty"`
	TransactionType TransactionType `json:"transaction_type,omitempty"`
	Amount          string          `json:"amount,omitempty"`
	CreditAmount    string          ` json:"credit_amount,omitempty"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	SettledAt       time.Time       `json:"settled_at,omitempty"`
}

func (c *Client) WalletTransaction() *WalletTransactionRequest {
	return &WalletTransactionRequest{
		client: c,
	}
}

func (bmr *WalletTransactionRequest) Create(walletTransactionInput *WalletTransactionInput) (*WalletTransaction, *Error) {
	clientRequest := &ClientRequest{
		Path:   "wallet_transactions",
		Result: &WalletTransactionResult{},
		Body:   walletTransactionInput,
	}

	result, err := bmr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	walletTransactionResult := result.(*WalletTransaction)

	return walletTransactionResult, nil
}
