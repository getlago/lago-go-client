package lago_test

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"

	"github.com/google/uuid"
)

var mockWalletTransactionListResponse = `{
			"wallet_transactions": [{
				"lago_id": "b1b2c3d4-e5f6-7890-1234-56789abcdef0",
				"lago_wallet_id": "a1a2b3b4-c5d6-7890-1234-56789abcdef0",
				"status": "settled",
				"transaction_type": "outbound",
				"transaction_status": "granted",
				"amount": "50.00",
				"credit_amount": "50.00",
				"invoice_requires_successful_payment": true,
				"created_at": "2024-06-01T12:00:00Z",
				"settled_at": "2024-06-01T12:05:00Z",
				"failed_at": "2024-06-01T12:10:00Z",
				"metadata": [{
					"key": "source",
					"value": "test"
				}],
				"name": "Test Transaction",
				"payment_method": {
					"payment_method_type": "card",
					"payment_method_id": "pm_wt_123"
				}
			}],
			"meta": {
				"current_page": 1,
				"next_page": null,
				"prev_page": null,
				"total_pages": 1,
				"total_count": 1
			}
		}`

func TestWalletTransactionRequest_Create(t *testing.T) {
	t.Run("When the server is not reachable", func(t *testing.T) {
		c := qt.New(t)

		client := New().SetBaseURL("http://localhost:88888").SetApiKey("test_api_key")
		result, err := client.Subscription().Terminate(context.Background(), SubscriptionTerminateInput{
			ExternalID: "1a901a90-1a90-1a90-1a90-1a901a901a90",
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err.Error(), qt.Equals, `{"status":0,"error":"","code":"","err":"Delete \"http://localhost:88888/api/v1/subscriptions/1a901a90-1a90-1a90-1a90-1a901a901a90\": dial tcp: address 88888: invalid port"}`)
	})

	t.Run("When sending all the parameters", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/wallet_transactions").
			MatchJSONBody(`{
				"wallet_transaction": {
					"granted_credits":                     "0.00",
					"invoice_requires_successful_payment": true,
					"metadata":                            [{"key": "source", "value": "test"}],
					"name":                                "Test Transaction",
					"paid_credits":                        "50.00",
					"voided_credits":                      "0.00",
					"wallet_id":                           "1a901a90-1a90-1a90-1a90-1a901a901a90"
				}
			}`).
			MockResponse(mockWalletTransactionListResponse)
		defer server.Close()

		result, err := server.Client().WalletTransaction().Create(context.Background(), &WalletTransactionInput{
			WalletID:                         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			Name:                             "Test Transaction",
			PaidCredits:                      "50.00",
			GrantedCredits:                   "0.00",
			VoidedCredits:                    "0.00",
			InvoiceRequiresSuccessfulPayment: true,
			Metadata: []WalletTransactionMetadata{
				{Key: "source", Value: "test"},
			},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result, qt.DeepEquals, &WalletTransactionResult{
			WalletTransactions: []WalletTransaction{
				{
					LagoID:                           uuid.MustParse("b1b2c3d4-e5f6-7890-1234-56789abcdef0"),
					LagoWalletID:                     uuid.MustParse("a1a2b3b4-c5d6-7890-1234-56789abcdef0"),
					Status:                           "settled",
					TransactionType:                  "outbound",
					TransactionStatus:                "granted",
					Amount:                           "50.00",
					CreditAmount:                     "50.00",
					InvoiceRequiresSuccessfulPayment: true,
					CreatedAt:                        time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					SettledAt:                        time.Date(2024, 6, 1, 12, 5, 0, 0, time.UTC),
					FailedAt:                         time.Date(2024, 6, 1, 12, 10, 0, 0, time.UTC),
					Metadata:                         []WalletTransactionMetadata{{Key: "source", Value: "test"}},
					Name:                             "Test Transaction",
					PaymentMethod: &PaymentMethodInput{
						PaymentMethodType: "card",
						PaymentMethodID:   "pm_wt_123",
					},
				},
			},
			Meta: Metadata{
				CurrentPage: 1,
				NextPage:    0,
				PrevPage:    0,
				TotalPages:  1,
				TotalCount:  1,
			},
		})
	})
}

func TestWalletTransactionRequest_CreateWithPaymentMethod(t *testing.T) {
	t.Run("When creating a wallet transaction with an invalid payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/wallet_transactions").
			MatchJSONBody(`{
				"wallet_transaction": {
					"granted_credits":                     "0.00",
					"invoice_requires_successful_payment": true,
					"metadata":                            [{"key": "source", "value": "test"}],
					"name":                                "Test Transaction",
					"paid_credits":                        "50.00",
					"voided_credits":                      "0.00",
					"wallet_id":                           "1a901a90-1a90-1a90-1a90-1a901a901a90",
					"payment_method": {
						"payment_method_type": "invalid",
						"payment_method_id": "pm_invalid"
					}
				}
			}`).
			MockResponseWithCode(422, map[string]any{
				"status": 422,
				"error":  "Unprocessable Entity",
				"code":   "validation_errors",
				"error_details": map[string]any{
					"payment_method": []string{"invalid_payment_method"},
				},
			})
		defer server.Close()

		result, err := server.Client().WalletTransaction().Create(context.Background(), &WalletTransactionInput{
			WalletID:                         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			Name:                             "Test Transaction",
			PaidCredits:                      "50.00",
			GrantedCredits:                   "0.00",
			VoidedCredits:                    "0.00",
			InvoiceRequiresSuccessfulPayment: true,
			Metadata: []WalletTransactionMetadata{
				{Key: "source", Value: "test"},
			},
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "invalid",
				PaymentMethodID:   "pm_invalid",
			},
		})
		c.Assert(result, qt.IsNil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.HTTPStatusCode, qt.Equals, 422)
		c.Assert(err.Message, qt.Equals, "Unprocessable Entity")
		c.Assert(err.ErrorCode, qt.Equals, "validation_errors")
		c.Assert(err.ErrorDetail, qt.IsNotNil)
		details, detailErr := err.ErrorDetail.Details()
		c.Assert(detailErr, qt.IsNil)
		c.Assert(details["payment_method"], qt.DeepEquals, []string{"invalid_payment_method"})
	})

	t.Run("When creating a wallet transaction with a payment method", func(t *testing.T) {
		c := qt.New(t)

		server := lt.NewMockServer(c).
			MatchMethod("POST").
			MatchPath("/api/v1/wallet_transactions").
			MatchJSONBody(`{
				"wallet_transaction": {
					"granted_credits":                     "0.00",
					"invoice_requires_successful_payment": true,
					"metadata":                            [{"key": "source", "value": "test"}],
					"name":                                "Test Transaction",
					"paid_credits":                        "50.00",
					"voided_credits":                      "0.00",
					"wallet_id":                           "1a901a90-1a90-1a90-1a90-1a901a901a90",
					"payment_method": {
						"payment_method_type": "card",
						"payment_method_id": "pm_wt_123"
					}
				}
			}`).
			MockResponse(mockWalletTransactionListResponse)
		defer server.Close()

		result, err := server.Client().WalletTransaction().Create(context.Background(), &WalletTransactionInput{
			WalletID:                         "1a901a90-1a90-1a90-1a90-1a901a901a90",
			Name:                             "Test Transaction",
			PaidCredits:                      "50.00",
			GrantedCredits:                   "0.00",
			VoidedCredits:                    "0.00",
			InvoiceRequiresSuccessfulPayment: true,
			Metadata: []WalletTransactionMetadata{
				{Key: "source", Value: "test"},
			},
			PaymentMethod: &PaymentMethodInput{
				PaymentMethodType: "card",
				PaymentMethodID:   "pm_wt_123",
			},
		})
		c.Assert(err == nil, qt.IsTrue)
		c.Assert(result.WalletTransactions, qt.HasLen, 1)
		c.Assert(result.WalletTransactions[0].PaymentMethod, qt.IsNotNil)
		c.Assert(result.WalletTransactions[0].PaymentMethod.PaymentMethodType, qt.Equals, "card")
		c.Assert(result.WalletTransactions[0].PaymentMethod.PaymentMethodID, qt.Equals, "pm_wt_123")
	})
}
