package lago

import (
	"fmt"

	"github.com/google/uuid"
)

type CustomerParams struct {
	Customer CustomerInput `json:"customer"`
}

type CustomerResult struct {
	Customer Customer `json:"customer"`
}

type CustomerInput struct {
	CustomerID string `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type Customer struct {
	LagoID     uuid.UUID `json:"lago_id"`
	CustomerID string    `json:"customer_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
}

type CustomerRequest struct {
	client *Client
}

func (c *Client) Customer() *CustomerRequest {
	return &CustomerRequest{
		client: c,
	}
}

func (cr *CustomerRequest) Create(ci *CustomerInput) *Customer {
	cp := &CustomerParams{
		Customer: *ci,
	}

	resp, err := cr.client.HttpClient.R().
		SetBody(*cp).
		SetResult(&Customer{}).
		Post("customers")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return nil
}
