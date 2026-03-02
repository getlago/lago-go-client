package lago

import (
	"context"
	"fmt"
)

type CustomerWalletMetadataRequest struct {
	client *Client
}

func (c *Client) CustomerWalletMetadata() *CustomerWalletMetadataRequest {
	return &CustomerWalletMetadataRequest{
		client: c,
	}
}

func (cwmr *CustomerWalletMetadataRequest) Replace(ctx context.Context, customerID string, walletCode string, metadata map[string]*string) (map[string]*string, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "customers", customerID, "wallets", walletCode, "metadata")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletMetadataResult{},
		Body:   &WalletMetadataParams{Metadata: metadata},
	}

	result, err := cwmr.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	metadataResult, ok := result.(*WalletMetadataResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return metadataResult.Metadata, nil
}

func (cwmr *CustomerWalletMetadataRequest) Merge(ctx context.Context, customerID string, walletCode string, metadata map[string]*string) (map[string]*string, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "customers", customerID, "wallets", walletCode, "metadata")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletMetadataResult{},
		Body:   &WalletMetadataParams{Metadata: metadata},
	}

	result, err := cwmr.client.Patch(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	metadataResult, ok := result.(*WalletMetadataResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return metadataResult.Metadata, nil
}

func (cwmr *CustomerWalletMetadataRequest) DeleteAll(ctx context.Context, customerID string, walletCode string) (map[string]*string, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "customers", customerID, "wallets", walletCode, "metadata")

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletMetadataResult{},
	}

	result, err := cwmr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	metadataResult, ok := result.(*WalletMetadataResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return metadataResult.Metadata, nil
}

func (cwmr *CustomerWalletMetadataRequest) DeleteKey(ctx context.Context, customerID string, walletCode string, key string) (map[string]*string, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "customers", customerID, "wallets", walletCode, "metadata", key)

	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &WalletMetadataResult{},
	}

	result, err := cwmr.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	metadataResult, ok := result.(*WalletMetadataResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return metadataResult.Metadata, nil
}
