package lago

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req interface{}) (interface{}, error) {
	args := m.Called(req)
	return args.Get(0), args.Error(1)
}

func TestBillingEntityCreate(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{HttpClient: mockClient}
	billingEntityReq := client.BillingEntity()

	expectedBillingEntity := &BillingEntity{
		LagoID:          uuid.New(),
		Name:            "Test Company",
		Code:            "TEST123",
		Email:           "test@example.com",
		AddressLine1:    "123 Test St",
		City:            "Test City",
		Country:         "US",
		DefaultCurrency: USD,
		LegalName:       "Test Company Legal",
		LegalNumber:     "123456789",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockClient.On("Do", mock.Anything).Return(&BillingEntityResult{
		BillingEntity: expectedBillingEntity,
	}, nil)

	result, err := billingEntityReq.Create(context.Background(), &BillingEntityCreateInput{
		Name:            "Test Company",
		Code:            "TEST123",
		Email:           "test@example.com",
		AddressLine1:    "123 Test St",
		City:            "Test City",
		Country:         "US",
		DefaultCurrency: USD,
		LegalName:       "Test Company Legal",
		LegalNumber:     "123456789",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBillingEntity.Name, result.Name)
	assert.Equal(t, expectedBillingEntity.Code, result.Code)
	mockClient.AssertExpectations(t)
}

func TestBillingEntityGet(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{HttpClient: mockClient}
	billingEntityReq := client.BillingEntity()

	expectedBillingEntity := &BillingEntity{
		LagoID:          uuid.New(),
		Name:            "Test Company",
		Code:            "TEST123",
		Email:           "test@example.com",
		AddressLine1:    "123 Test St",
		City:            "Test City",
		Country:         "US",
		DefaultCurrency: USD,
	}

	mockClient.On("Do", mock.Anything).Return(&BillingEntityResult{
		BillingEntity: expectedBillingEntity,
	}, nil)

	result, err := billingEntityReq.Get(context.Background(), "TEST123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBillingEntity.Name, result.Name)
	assert.Equal(t, expectedBillingEntity.Code, result.Code)
	mockClient.AssertExpectations(t)
}

func TestBillingEntityGetList(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{HttpClient: mockClient}
	billingEntityReq := client.BillingEntity()

	expectedBillingEntities := []BillingEntity{
		{
			LagoID:          uuid.New(),
			Name:            "Test Company 1",
			Code:            "TEST123",
			Email:           "test1@example.com",
			AddressLine1:    "123 Test St",
			City:            "Test City",
			Country:         "US",
			DefaultCurrency: USD,
		},
		{
			LagoID:          uuid.New(),
			Name:            "Test Company 2",
			Code:            "TEST456",
			Email:           "test2@example.com",
			AddressLine1:    "456 Test St",
			City:            "Test City",
			Country:         "US",
			DefaultCurrency: EUR,
		},
	}

	mockClient.On("Do", mock.Anything).Return(&BillingEntityResult{
		BillingEntities: expectedBillingEntities,
		Meta: Metadata{
			CurrentPage: 1,
			TotalPages:  1,
			TotalCount:  2,
		},
	}, nil)

	result, err := billingEntityReq.GetList(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.BillingEntities, 2)
	assert.Equal(t, expectedBillingEntities[0].Name, result.BillingEntities[0].Name)
	assert.Equal(t, expectedBillingEntities[1].Name, result.BillingEntities[1].Name)
	mockClient.AssertExpectations(t)
}

func TestBillingEntityUpdate(t *testing.T) {
	mockClient := new(MockHTTPClient)
	client := &Client{HttpClient: mockClient}
	billingEntityReq := client.BillingEntity()

	expectedBillingEntity := &BillingEntity{
		LagoID:          uuid.New(),
		Name:            "Updated Company",
		Code:            "TEST123",
		Email:           "updated@example.com",
		AddressLine1:    "456 Updated St",
		City:            "Updated City",
		Country:         "US",
		DefaultCurrency: USD,
	}

	mockClient.On("Do", mock.Anything).Return(&BillingEntityResult{
		BillingEntity: expectedBillingEntity,
	}, nil)

	updateInput := &BillingEntityUpdateInput{
		Name:         "Updated Company",
		Email:        "updated@example.com",
		AddressLine1: "456 Updated St",
		City:         "Updated City",
	}

	result, err := billingEntityReq.Update(context.Background(), "TEST123", updateInput)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBillingEntity.Name, result.Name)
	assert.Equal(t, expectedBillingEntity.Email, result.Email)
	assert.Equal(t, expectedBillingEntity.AddressLine1, result.AddressLine1)
	mockClient.AssertExpectations(t)
} 