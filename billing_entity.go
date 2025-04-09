package lago

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BillingEntityResult struct {
	BillingEntity  *BillingEntity  `json:"billing_entity,omitempty"`
	BillingEntities []BillingEntity `json:"billing_entities,omitempty"`
	Meta           Metadata         `json:"meta,omitempty"`
}

type BillingEntity struct {
	LagoID                uuid.UUID `json:"lago_id,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Code                  string    `json:"code,omitempty"`
	Description           string    `json:"description,omitempty"`
	AddressLine1          string    `json:"address_line1,omitempty"`
	AddressLine2          string    `json:"address_line2,omitempty"`
	City                  string    `json:"city,omitempty"`
	State                 string    `json:"state,omitempty"`
	Zipcode               string    `json:"zipcode,omitempty"`
	Country               string    `json:"country,omitempty"`
	Email                 string    `json:"email,omitempty"`
	Phone                 string    `json:"phone,omitempty"`
	LegalName             string    `json:"legal_name,omitempty"`
	LegalNumber           string    `json:"legal_number,omitempty"`
	TaxIdentificationNumber string    `json:"tax_identification_number,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
	IsDefault             bool      `json:"is_default,omitempty"`
}

type BillingEntityRequest struct {
	client *Client
}

func (c *Client) BillingEntity() *BillingEntityRequest {
	return &BillingEntityRequest{
		client: c,
	}
}

func (ber *BillingEntityRequest) Get(ctx context.Context, code string) (*BillingEntity, *Error) {
	subPath := fmt.Sprintf("%s/%s", "billing_entities", code)
	clientRequest := &ClientRequest{
		Path:   subPath,
		Result: &BillingEntityResult{},
	}

	result, err := ber.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billingEntityResult, ok := result.(*BillingEntityResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billingEntityResult.BillingEntity, nil
}

func (ber *BillingEntityRequest) GetList(ctx context.Context) (*BillingEntityResult, *Error) {
	clientRequest := &ClientRequest{
		Path:   "billing_entities",
		Result: &BillingEntityResult{},
	}

	result, err := ber.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	billingEntityResult, ok := result.(*BillingEntityResult)
	if !ok {
		return nil, &ErrorTypeAssert
	}

	return billingEntityResult, nil
}
