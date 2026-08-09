package lago_test

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/getlago/lago-go-client"
	lt "github.com/getlago/lago-go-client/testing"

	qt "github.com/frankban/quicktest"
)

func TestOverdueBalanceRequest_GetList(t *testing.T) {
	c := qt.New(t)

	server := lt.NewMockServer(c).
		MatchMethod("GET").
		MatchPath("/api/v1/analytics/overdue_balance").
		MatchQuery(map[string]string{}).
		MockResponse(`{
			"overdue_balances": [{
				"month": "2026-07-01T00:00:00.000Z",
				"amount_cents": "500.0",
				"currency": "USD",
				"lago_invoice_ids": ["275e091b-2e7e-4d1e-ae12-646994c9f8d7"],
				"billing_entity_id": "bbea5aa6-197b-4ff1-a983-1884aea618e9"
			}]
		}`)
	defer server.Close()

	result, err := server.Client().OverdueBalance().GetList(context.Background(), &OverdueBalanceListInput{})

	c.Assert(err == nil, qt.IsTrue)
	c.Assert(result.OverdueBalances, qt.HasLen, 1)
	c.Assert(result.OverdueBalances[0].Month, qt.Equals, "2026-07-01T00:00:00.000Z")
	c.Assert(result.OverdueBalances[0].AmountCents, qt.Equals, 500)
	c.Assert(result.OverdueBalances[0].AmountCurrency, qt.Equals, USD)
}

func TestOverdueBalance_UnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		name        string
		amountCents string
		want        int
		wantErr     bool
	}{
		{name: "integer", amountCents: "500", want: 500},
		{name: "decimal", amountCents: "500.0", want: 500},
		{name: "quoted integer", amountCents: `"500"`, want: 500},
		{name: "quoted decimal", amountCents: `"500.0"`, want: 500},
		{name: "fractional", amountCents: "500.5", wantErr: true},
		{name: "quoted fractional", amountCents: `"500.5"`, wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var balance OverdueBalance
			err := json.Unmarshal([]byte(`{"amount_cents":`+test.amountCents+`}`), &balance)

			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if balance.AmountCents != test.want {
				t.Fatalf("AmountCents = %d, want %d", balance.AmountCents, test.want)
			}
		})
	}
}
