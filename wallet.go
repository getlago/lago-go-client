package lago

import (
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

type WalletRequest struct {
	client *Client
}

type WalletParams struct {
	WalletInput *WalletInput
}

type WalletInput struct {
	RateAmount         string `json:"rate_amount,omitempty"`
	Name               string `json:"name,omitempty"`
	PaidCredits        string `json:"paid_credits,omitempty"`
	GrantedCredits     string `json:"granted_credits,omitempty"`
	ExpirationDate     string `json:"expiration_date,omitempty"`
	ExternalCustomerId string `json:"external_customer_id,omitempty"`
}

type WalletListInput struct {
	PerPage            int `json:"per_page,omitempty,string"`
	Page               int `json:"page,omitempty,string"`
	ExternalCustomerID int `json:"external_customer_id,omitempty,string"`
}

type WalletResult struct {
	Wallet  *Wallet  `json:"wallet,omitempty"`
	Wallets []Wallet `json:"wallets,omitempty"`
	Meta    Metadata `json:"meta,omitempty"`
}

type Wallet struct {
	LagoID               uuid.UUID `json:"lago_id,omitempty"`
	LagoCustomerID       uuid.UUID `json:"lago_customer_id,omitempty"`
	ExternalCustomerID   string    `json:"external_customer_id,omitempty"`
	Status               Status    `json:"status,omitempty"`
	Currency             string    `json:"currency,omitempty"`
	Name                 string    `json:"name,omitempty"`
	RateAmount           string    `json:"rate_amount,omitempty"`
	CreditsBalance       string    `json:"credits_balance,omitempty"`
	Balance              string    `json:"balance,omitempty"`
	ConsumedCredits      string    `json:"consumed_credits,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	ExpirationDate       time.Time `json:"expiration_date,omitempty"`
	LastBalanceSyncAt    time.Time `json:"last_balance_sync_at,omitempty"`
	LastConsumedCreditAt time.Time `json:"last_consumed_credit_at,omitempty"`
	TerminatedAt         time.Time `json:"terminated_at,omitempty"`
}

func (c *Client) Wallet() *WalletRequest {
	return &WalletRequest{
		client: c,
	}
}

func (bmr *WalletRequest) Get(walletId string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s", "wallets", walletId)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
	}

	result, err := bmr.client.Get(clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult := result.(*WalletResult)

	return walletResult.Wallet, nil
}

func (bmr *WalletRequest) GetList(walletListInput *WalletListInput) (*WalletResult, *Error) {
	jsonQueryParams, err := json.Marshal(walletListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	json.Unmarshal(jsonQueryParams, &queryParams)

	clientRequest := &ClientRequest{
		Path:        "wallets",
		QueryParams: queryParams,
		Result:      &WalletResult{},
	}

	result, clientErr := bmr.client.Get(clientRequest)
	if err != nil {
		return nil, clientErr
	}

	walletResult := result.(*WalletResult)

	return walletResult, nil
}

func (bmr *WalletRequest) Create(walletInput *WalletInput) (*Wallet, *Error) {
	clientRequest := &ClientRequest{
		Path:   "wallets",
		Result: &WalletResult{},
		Body:   walletInput,
	}

	result, err := bmr.client.Post(clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult := result.(*WalletResult)

	return walletResult.Wallet, nil
}

func (bmr *WalletRequest) Update(walletInput *WalletInput, walletId string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s", "wallets", walletId)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
		Body:   walletInput,
	}

	result, err := bmr.client.Put(clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult := result.(*WalletResult)

	return walletResult.Wallet, nil
}

func (bmr *WalletRequest) Delete(walletId string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s", "wallets", walletId)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
	}

	result, err := bmr.client.Delete(clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult := result.(*WalletResult)

	return walletResult.Wallet, nil
}
