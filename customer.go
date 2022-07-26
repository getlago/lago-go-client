package lago

import "github.com/google/uuid"

type CustomerParams struct {
	Customer CustomerInput `json:"customer"`
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
