package lago

import (
	"context"
	"encoding/json"
	"fmt"
)

type CustomerWalletRequest struct {
	client *Client
}

func (c *Client) CustomerWallet() *CustomerWalletRequest {
	return &CustomerWalletRequest{
		client: c,
	}
}

func (cwr *CustomerWalletRequest) Get(ctx context.Context, customerExternalID string, walletCode string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
	}

	result, err := cwr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}

func (cwr *CustomerWalletRequest) GetList(ctx context.Context, customerExternalID string, walletListInput *WalletListInput) (*WalletResult, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", customerExternalID, "wallets")

	jsonQueryParams, err := json.Marshal(walletListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &WalletResult{},
	}

	result, clientErr := cwr.client.Get(ctx, clientRequest)
	if clientErr != nil {
		return nil, clientErr
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult, nil
}

func (cwr *CustomerWalletRequest) Create(ctx context.Context, customerExternalID string, walletInput *WalletInput) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s", "customers", customerExternalID, "wallets")

	walletParams := &WalletParams{
		WalletInput: walletInput,
	}

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
		Body:   walletParams,
	}

	result, err := cwr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}

func (cwr *CustomerWalletRequest) Update(ctx context.Context, customerExternalID string, walletCode string, walletInput *WalletInput) (*Wallet, *Error) {
	walletParams := &WalletParams{
		WalletInput: walletInput,
	}

	subPath := fmt.Sprintf("%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
		Body:   walletParams,
	}

	result, err := cwr.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}

func (cwr *CustomerWalletRequest) Delete(ctx context.Context, customerExternalID string, walletCode string) (*Wallet, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletResult{},
	}

	result, err := cwr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	walletResult, ok := result.(*WalletResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return walletResult.Wallet, nil
}
