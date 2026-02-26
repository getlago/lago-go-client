package lago_test

import (
	"encoding/json"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/getlago/lago-go-client"
)

func TestProviderError_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals a JSON object", func(t *testing.T) {
		c := qt.New(t)

		jsonData := []byte(`{"error_message":"message","error_code":"code"}`)

		var pe lago.ProviderError
		err := json.Unmarshal(jsonData, &pe)
		c.Assert(err, qt.IsNil)

		c.Assert(pe["error_message"], qt.Equals, "message")
		c.Assert(pe["error_code"], qt.Equals, "code")
	})

	t.Run("unmarshals a plain string into message key", func(t *testing.T) {
		c := qt.New(t)

		jsonData := []byte(`"something went wrong"`)

		var pe lago.ProviderError
		err := json.Unmarshal(jsonData, &pe)
		c.Assert(err, qt.IsNil)

		c.Assert(pe["message"], qt.Equals, "something went wrong")
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		c := qt.New(t)

		jsonData := []byte(`not valid json`)

		var pe lago.ProviderError
		err := json.Unmarshal(jsonData, &pe)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err.Error(), qt.Contains, "invalid character")
	})

	t.Run("unmarshals nested object", func(t *testing.T) {
		c := qt.New(t)

		jsonData := []byte(`{"error":{"code":"invalid_card","message":"Your card was declined"}}`)

		var pe lago.ProviderError
		err := json.Unmarshal(jsonData, &pe)
		c.Assert(err, qt.IsNil)

		nested, ok := pe["error"].(map[string]any)
		c.Assert(ok, qt.IsTrue)

		c.Assert(nested["code"], qt.Equals, "invalid_card")
		c.Assert(nested["message"], qt.Equals, "Your card was declined")
	})

	t.Run("unmarshals empty object", func(t *testing.T) {

		c := qt.New(t)

		jsonData := []byte(`{}`)

		var pe lago.ProviderError
		err := json.Unmarshal(jsonData, &pe)
		c.Assert(err, qt.IsNil)

		c.Assert(len(pe), qt.Equals, 0)
	})

	t.Run("unmarshals empty string", func(t *testing.T) {

		c := qt.New(t)

		jsonData := []byte(`""`)

		var pe lago.ProviderError
		err := json.Unmarshal(jsonData, &pe)
		c.Assert(err, qt.IsNil)

		c.Assert(pe["message"], qt.Equals, "")
	})
}
