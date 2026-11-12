package lago_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	. "github.com/getlago/lago-go-client"
)

const (
	standardWebhookSecret      = "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	otherStandardWebhookSecret = "whsec_C2FVsBQIhrscChlQIMV+b5sSYspob7oD"
	standardWebhookID          = "msg_p5jXN8AQM9LWM0D4loKWxJek"
	standardWebhookBody        = `{"test": 2432232314}`
)

// signStandardWebhook computes a Standard Webhooks v1 signature by hand,
// without the library, so the helper is checked against the spec itself.
func signStandardWebhook(c *qt.C, secret, id, timestamp string, body []byte) string {
	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, "whsec_"))
	c.Assert(err, qt.IsNil)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id + "." + timestamp + "." + string(body)))

	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func standardWebhookTimestamp(offset time.Duration) string {
	return strconv.FormatInt(time.Now().Add(offset).Unix(), 10)
}

func standardWebhookHeaders(id, timestamp, signature string) http.Header {
	headers := http.Header{}
	headers.Set("webhook-id", id)
	headers.Set("webhook-timestamp", timestamp)
	headers.Set("webhook-signature", signature)

	return headers
}

func validStandardWebhookHeaders(c *qt.C, offset time.Duration) http.Header {
	timestamp := standardWebhookTimestamp(offset)
	signature := signStandardWebhook(c, standardWebhookSecret, standardWebhookID, timestamp, []byte(standardWebhookBody))

	return standardWebhookHeaders(standardWebhookID, timestamp, signature)
}

func TestVerifyStandardWebhook(t *testing.T) {
	t.Run("Returns nil for a valid delivery", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)

		c.Assert(VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody)), qt.IsNil)
	})

	t.Run("Returns nil when a wrong v1 signature comes before the right one", func(t *testing.T) {
		c := qt.New(t)

		timestamp := standardWebhookTimestamp(0)
		right := signStandardWebhook(c, standardWebhookSecret, standardWebhookID, timestamp, []byte(standardWebhookBody))
		wrong := signStandardWebhook(c, otherStandardWebhookSecret, standardWebhookID, timestamp, []byte(standardWebhookBody))
		headers := standardWebhookHeaders(standardWebhookID, timestamp, wrong+" "+right)

		c.Assert(VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody)), qt.IsNil)
	})

	t.Run("Returns an error when the body changed by one byte", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)

		err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(`{"test": 2432232315}`))
		c.Assert(err, qt.ErrorMatches, ".*no matching signature found")
	})

	t.Run("Returns an error with another valid secret", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)

		err := VerifyStandardWebhook(otherStandardWebhookSecret, headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, ".*no matching signature found")
	})

	t.Run("Returns an error with a different webhook-id", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)
		headers.Set("webhook-id", "msg_another")

		err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, ".*no matching signature found")
	})

	t.Run("Accepts a secret without the whsec_ prefix", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)

		err := VerifyStandardWebhook(strings.TrimPrefix(standardWebhookSecret, "whsec_"), headers, []byte(standardWebhookBody))
		c.Assert(err, qt.IsNil)
	})

	t.Run("Returns an error for a secret that is not base64", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)

		err := VerifyStandardWebhook("whsec_not*base64", headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, "standard webhooks: the signing secret is not valid base64.*")
	})

	for _, secret := range []string{"", "whsec_"} {
		t.Run("Returns an error for the empty secret "+strconv.Quote(secret), func(t *testing.T) {
			c := qt.New(t)

			// Signed with an empty key, which anyone can compute.
			timestamp := standardWebhookTimestamp(0)
			signature := signStandardWebhook(c, secret, standardWebhookID, timestamp, []byte(standardWebhookBody))
			headers := standardWebhookHeaders(standardWebhookID, timestamp, signature)

			err := VerifyStandardWebhook(secret, headers, []byte(standardWebhookBody))
			c.Assert(err, qt.ErrorMatches, "standard webhooks: the signing secret is empty")
		})
	}

	t.Run("Returns an error for a timestamp 6 minutes in the past", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, -6*time.Minute)

		err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, ".*message timestamp too old")
	})

	t.Run("Returns an error for a timestamp 6 minutes in the future", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 6*time.Minute)

		err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, ".*message timestamp too new")
	})

	t.Run("Returns nil for a timestamp 4 minutes in the past", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, -4*time.Minute)

		c.Assert(VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody)), qt.IsNil)
	})

	t.Run("Returns nil for a timestamp 4 minutes in the future", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 4*time.Minute)

		c.Assert(VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody)), qt.IsNil)
	})

	for _, header := range []string{"webhook-id", "webhook-timestamp", "webhook-signature"} {
		t.Run("Returns an error when "+header+" is absent", func(t *testing.T) {
			c := qt.New(t)

			headers := validStandardWebhookHeaders(c, 0)
			headers.Del(header)

			err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
			c.Assert(err, qt.ErrorMatches, "standard webhooks: the "+header+" header is missing")
		})

		t.Run("Returns an error when "+header+" is empty", func(t *testing.T) {
			c := qt.New(t)

			headers := validStandardWebhookHeaders(c, 0)
			headers.Set(header, "")

			err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
			c.Assert(err, qt.ErrorMatches, "standard webhooks: the "+header+" header is missing")
		})
	}

	t.Run("Returns an error when webhook-timestamp is not an integer", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)
		headers.Set("webhook-timestamp", "abc")

		err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, "(?s).*invalid signature headers.*")
	})

	t.Run("Returns an error when webhook-signature holds no v1 entry", func(t *testing.T) {
		c := qt.New(t)

		headers := validStandardWebhookHeaders(c, 0)
		signature := headers.Get("webhook-signature")
		headers.Set("webhook-signature", "v2,"+strings.TrimPrefix(signature, "v1,"))

		err := VerifyStandardWebhook(standardWebhookSecret, headers, []byte(standardWebhookBody))
		c.Assert(err, qt.ErrorMatches, ".*no matching signature found")
	})
}
