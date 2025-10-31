package integrationtest

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/getlago/lago-go-client"
)

func skipIntegrationTest(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS_ENABLED") != "true" {
		t.Skip("Integration tests are not enabled")
	}
}

func client() *lago.Client {
	baseUrl := os.Getenv("TEST_LAGO_API_URL")
	apiKey := os.Getenv("TEST_LAGO_API_KEY")
	if baseUrl == "" || apiKey == "" {
		log.Fatalf("TEST_LAGO_API_URL and TEST_LAGO_API_KEY must be set")
	}
	return lago.New().SetBaseURL(baseUrl).SetApiKey(apiKey)
}

func prefix() string {
	return strings.ReplaceAll(time.Now().Format("2006-01-02T15-04-05-.999"), ".", "")
}

func uniqueId() string {
	return fmt.Sprintf("go-%s", prefix())
}

type CustomerCreateOption func(*lago.CustomerInput)

func CustomerUsPreset() CustomerCreateOption {
	return func(input *lago.CustomerInput) {
		input.Currency = "USD"
		input.Country = "US"
		input.AddressLine1 = "123 Main St"
		input.AddressLine2 = "Apt 1"
		input.City = "San Francisco"
		input.Zipcode = "90001"
		input.State = "CA"
		input.Phone = "0601020304"
		input.LegalNumber = "US1234567890"
		input.Timezone = "America/New_York"
	}
}

func CustomerFrPreset() CustomerCreateOption {
	return func(input *lago.CustomerInput) {
		input.Currency = "EUR"
		input.Country = "FR"
		input.AddressLine1 = "123 Main St"
		input.AddressLine2 = "Apt 1"
		input.City = "Paris"
		input.Zipcode = "75001"
		input.State = "Paris"
		input.Phone = "0601020304"
		input.LegalNumber = "FR1234567890"
		input.Timezone = "Europe/Paris"
	}
}

func CustomerWithParams(fn func(input *lago.CustomerInput)) CustomerCreateOption {
	return func(input *lago.CustomerInput) {
		fn(input)
	}
}

func CustomerWithMetadata(key string, value string) CustomerCreateOption {
	return func(input *lago.CustomerInput) {
		if input.Metadata == nil {
			input.Metadata = []lago.CustomerMetadataInput{}
		}
		input.Metadata = append(input.Metadata, lago.CustomerMetadataInput{
			Key:   key,
			Value: value,
		})
	}
}

func createCustomer(t *testing.T, options ...CustomerCreateOption) *lago.Customer {
	externalID := uniqueId()
	input := &lago.CustomerInput{
		Email:      fmt.Sprintf("yohan+%s@getlago.com", externalID),
		Name:       fmt.Sprintf("Name | %s", externalID),
		ExternalID: fmt.Sprintf("ExternalID %s", externalID),
		Firstname:  fmt.Sprintf("Firstname %s", externalID),
		Lastname:   fmt.Sprintf("Lastname %s", externalID),
		LegalName:  fmt.Sprintf("LegalName %s", externalID),
	}
	for _, option := range options {
		option(input)
	}
	customer, err := client().Customer().Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Error creating customer: %v", err)
	}
	return customer
}
