package lago

import (
	"context"
	"encoding/json"
	"fmt"
)

type CustomerWalletAlertRequest struct {
	client *Client
}

func (c *Client) CustomerWalletAlert() *CustomerWalletAlertRequest {
	return &CustomerWalletAlertRequest{
		client: c,
	}
}

func (ar *CustomerWalletAlertRequest) Get(ctx context.Context, customerExternalID string, walletCode string, alertCode string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
	}

	result, err := ar.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}

func (ar *CustomerWalletAlertRequest) GetList(ctx context.Context, customerExternalID string, walletCode string, alertListInput *AlertListInput) (*AlertResult, *Error) {
	jsonQueryParams, err := json.Marshal(alertListInput)
	if err != nil {
		return nil, &Error{Err: err}
	}

	queryParams := make(map[string]string)
	if err = json.Unmarshal(jsonQueryParams, &queryParams); err != nil {
		return nil, &Error{Err: err}
	}

	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode, "alerts")
	clientRequest := &ClientRequest{
		Path:        subPath,
		QueryParams: queryParams,
		Result:      &AlertResult{},
	}

	result, clientErr := ar.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, clientErr
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult, nil
}

func (ar *CustomerWalletAlertRequest) Create(ctx context.Context, customerExternalID string, walletCode string, alertInput *AlertInput) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode, "alerts")
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
		Body:   &AlertParams{Alert: alertInput},
	}

	result, err := ar.client.Post(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}

func (ar *CustomerWalletAlertRequest) Update(ctx context.Context, customerExternalID string, walletCode string, alertCode string, alertInput *AlertInput) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
		Body:   &AlertParams{Alert: alertInput},
	}

	result, err := ar.client.Put(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}

func (ar *CustomerWalletAlertRequest) Delete(ctx context.Context, customerExternalID string, walletCode string, alertCode string) (*Alert, *Error) {
	subPath := fmt.Sprintf("%s/%s/%s/%s/%s/%s", "customers", customerExternalID, "wallets", walletCode, "alerts", alertCode)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &AlertResult{},
	}

	result, err := ar.client.Delete(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	alertResult, ok := result.(*AlertResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return alertResult.Alert, nil
}
