package integrationtest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/getlago/lago-go-client"
	"github.com/stretchr/testify/assert"
)

func getCustomers(t *testing.T, input *lago.CustomerListInput) []lago.Customer {
	client := client()
	customers, err := client.Customer().GetList(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to get customers: %v", err)
	}
	return customers.Customers
}

func TestCustomer_GetAllCustomers(t *testing.T) {
	skipIntegrationTest(t)

	usCustomer := createCustomer(
		t,
		CustomerUsPreset(),
		CustomerWithMetadata("is_synced", "false"),
		CustomerWithParams(func(input *lago.CustomerInput) {
			input.TaxIdentificationNumber = "US1234567890"
			input.CustomerType = lago.CompanyCustomerType
		}),
	)

	frenchCustomer := createCustomer(
		t,
		CustomerFrPreset(),
		CustomerWithMetadata("is_synced", "true"),
		CustomerWithMetadata("last_synced_at", "2025-01-01"),
	)

	suffix := strings.Split(frenchCustomer.ExternalID, " ")[1]

	t.Run("without filters", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{})
		assert.Greater(t, len(customers), 1)
	})

	t.Run("by name", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			SearchTerm: fmt.Sprintf("e | %s", suffix),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by firstname", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			SearchTerm: fmt.Sprintf("rstname %s", suffix),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by lastname", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			SearchTerm: fmt.Sprintf("astname %s", suffix),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by legalname", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			SearchTerm: fmt.Sprintf("egalname %s", suffix),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by externalid", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			SearchTerm: fmt.Sprintf("ID %s", suffix),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by email", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			SearchTerm: fmt.Sprintf("n+%s@getlago.com", suffix),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by country", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			Countries: []string{"US"},
			Page:      lago.Ptr(1),
			PerPage:   lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Countries: []string{"US", "FR"},
			Page:      lago.Ptr(1),
			PerPage:   lago.Ptr(2),
		})
		assert.Len(t, customers, 2)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
		assert.Equal(t, usCustomer.LagoID, customers[1].LagoID)
	})

	t.Run("by state", func(t *testing.T) {

		customers := getCustomers(t, &lago.CustomerListInput{
			States:  []string{"CA"},
			Page:    lago.Ptr(1),
			PerPage: lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			States:  []string{"CA", "Paris"},
			Page:    lago.Ptr(1),
			PerPage: lago.Ptr(2),
		})
		assert.Len(t, customers, 2)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
		assert.Equal(t, usCustomer.LagoID, customers[1].LagoID)
	})

	t.Run("by zipcode", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			Zipcodes: []string{"90001"},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Zipcodes: []string{"90001", "75001"},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(2),
		})
		assert.Len(t, customers, 2)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
		assert.Equal(t, usCustomer.LagoID, customers[1].LagoID)
	})

	t.Run("by currency", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			Currencies: []lago.Currency{lago.USD},
			Page:       lago.Ptr(1),
			PerPage:    lago.Ptr(1),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Currencies: []lago.Currency{lago.USD, lago.EUR},
			Page:       lago.Ptr(1),
			PerPage:    lago.Ptr(2),
		})
		assert.Len(t, customers, 2)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
		assert.Equal(t, usCustomer.LagoID, customers[1].LagoID)
	})

	t.Run("by has_tax_identification_number", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			HasTaxIdentificationNumber: lago.Ptr(true),
			Page:                       lago.Ptr(1),
			PerPage:                    lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			HasTaxIdentificationNumber: lago.Ptr(false),
			Page:                       lago.Ptr(1),
			PerPage:                    lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by metadata", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			Metadata: map[string]string{"is_synced": "true"},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Metadata: map[string]string{"is_synced": "true", "last_synced_at": "2025-01-01"},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Metadata: map[string]string{"is_synced": "true", "last_synced_at": "2025-01-01", "first_synced_at": ""},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Metadata: map[string]string{"is_synced": "false"},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			Metadata: map[string]string{"last_synced_at": ""},
			Page:     lago.Ptr(1),
			PerPage:  lago.Ptr(1),
		})
		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by customer_type", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			CustomerType: lago.CompanyCustomerType,
			Page:         lago.Ptr(1),
			PerPage:      lago.Ptr(1),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)
	})

	t.Run("by has_customer_type", func(t *testing.T) {
		customers := getCustomers(t, &lago.CustomerListInput{
			HasCustomerType: lago.Ptr(true),
			Page:            lago.Ptr(1),
			PerPage:         lago.Ptr(1),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, usCustomer.LagoID, customers[0].LagoID)

		customers = getCustomers(t, &lago.CustomerListInput{
			HasCustomerType: lago.Ptr(false),
			Page:            lago.Ptr(1),
			PerPage:         lago.Ptr(1),
		})

		assert.Len(t, customers, 1)
		assert.Equal(t, frenchCustomer.LagoID, customers[0].LagoID)
	})
}
