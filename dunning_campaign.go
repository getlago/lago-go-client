package lago

type DunningCampaign struct {
	CustomerExternalID     string                          `json:"customer_external_id"`
	DunningCampaignCode    string                          `json:"dunning_campaign_code"`
	OverdueBalanceCents    int                             `json:"overdue_balance_cents"`
	OverdueBalanceCurrency Currency                        `json:"overdue_balance_currency"`
	OverdueBalances        []DunningCampaignOverdueBalance `json:"overdue_balances,omitempty"`
}

type DunningCampaignOverdueBalance struct {
	AmountCents int      `json:"amount_cents"`
	Currency    Currency `json:"currency"`
}
