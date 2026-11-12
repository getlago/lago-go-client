package lago

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

type WebhookRequest struct {
	client *Client
}

func (c *Client) Webhook() *WebhookRequest {
	return &WebhookRequest{
		client: c,
	}
}

func (wr *WebhookRequest) GetPublicKey(ctx context.Context) (*rsa.PublicKey, *Error) {
	clientRequest := &ClientRequest{
		Path: "webhooks/public_key",
	}

	result, err := wr.client.Get(ctx, clientRequest)
	if err != nil {
		return nil, err
	}

	validatedResult, ok := result.(string)
	if !ok {
		return nil, &Error{
			Err:            errors.New("response is not a string"),
			HTTPStatusCode: http.StatusInternalServerError,
			Message:        "response is not a string",
		}
	}

	// Decode the base64-encoded key
	bytesResult, decodeErr := base64.StdEncoding.DecodeString(validatedResult)
	if decodeErr != nil {
		return nil, &Error{
			Err:            decodeErr,
			HTTPStatusCode: http.StatusInternalServerError,
			Message:        "cannot decode the key",
		}
	}

	// Parse the PEM block
	block, _ := pem.Decode(bytesResult)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, &Error{
			Err:            errors.New("failed to decode pem block containing public key"),
			HTTPStatusCode: http.StatusInternalServerError,
			Message:        "Failed to decode PEM block containing public key",
		}
	}

	// Parse the DER-encoded public key
	publicKey, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
	if parseErr != nil {
		return nil, &Error{
			Err:            parseErr,
			HTTPStatusCode: http.StatusInternalServerError,
			Message:        "Failed to to parse the public key",
		}
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, &Error{
			Err:            errors.New("unexpected type of public key"),
			HTTPStatusCode: http.StatusInternalServerError,
			Message:        "Unexpected type of public key",
		}
	}

	return rsaPublicKey, nil
}

func (wr *WebhookRequest) parseSignature(ctx context.Context, signature string) (*jwt.Token, *Error) {
	publicKey, err := wr.GetPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	token, parseErr := jwt.Parse(signature, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if parseErr != nil {
		return nil, &Error{
			Err:            parseErr,
			HTTPStatusCode: http.StatusInternalServerError,
			Message:        "cannot parse token",
		}
	}

	return token, nil
}

func (wr *WebhookRequest) ValidateSignature(ctx context.Context, signature string) (bool, *Error) {
	if token, err := wr.parseSignature(ctx, signature); err == nil && token.Valid {
		return true, nil
	} else {
		return false, err
	}
}

func (wr *WebhookRequest) ValidateBody(ctx context.Context, signature string, body string) (bool, *Error) {
	if token, err := wr.parseSignature(ctx, signature); err == nil && token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return false, &Error{
				Err:            errors.New("error casting claims"),
				HTTPStatusCode: http.StatusInternalServerError,
				Message:        "cannot parse token",
			}
		}

		return claims["data"] == body, nil
	} else {
		return false, err
	}
}

// VerifyStandardWebhook checks that a delivery sent by a webhook endpoint
// whose SignatureAlgo is StandardWebhooks was signed by Lago.
//
// The secret is a `whsec_` signing key from WebhookEndpoint.SigningKeys. The
// headers are those of the incoming request, and the payload must be the raw
// request body bytes, read before any JSON decoding. A delivery whose
// webhook-timestamp is more than 5 minutes in the past or in the future is
// rejected. It returns nil when the delivery is valid.
//
// See https://www.standardwebhooks.com for the specification.
func VerifyStandardWebhook(secret string, headers http.Header, payload []byte) error {
	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, "whsec_"))
	if err != nil {
		return fmt.Errorf("standard webhooks: the signing secret is not valid base64: %w", err)
	}
	if len(key) == 0 {
		return errors.New("standard webhooks: the signing secret is empty")
	}

	for _, header := range []string{
		standardwebhooks.HeaderWebhookID,
		standardwebhooks.HeaderWebhookTimestamp,
		standardwebhooks.HeaderWebhookSignature,
	} {
		if headers.Get(header) == "" {
			return fmt.Errorf("standard webhooks: the %s header is missing", header)
		}
	}

	webhook, err := standardwebhooks.NewWebhookRaw(key)
	if err != nil {
		return err
	}

	return webhook.Verify(payload, headers)
}
