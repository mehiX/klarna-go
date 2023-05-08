package account

import "github.com/mehix/klarna-go/klarna/domain"

type RequestBalances struct {
	InsightsConsumerID         string   `json:"insights_consumer_id"`          // String
	InsightsAccountIDs         []string `json:"insights_account_ids"`          // ?Array<String>
	ExcludedInsightsAccountIDs []string `json:"excluded_insights_account_ids"` // ?Array<String>
	RequiredDataAvailability   string   `json:"required_data_availability"`    // ?Enum<'NONE', 'BALANCE_OR_AVAILABLE'>
}

type ResponseBalances struct {
	Data DataBalances `json:"data"`
}

type DataBalances struct {
	Reports []Balances `json:"reports"`
}

type Balances struct {
	Type               string        `json:"type"`                 // "BALANCES"
	InsightsConsumerID string        `json:"insights_consumer_id"` // String
	InsightsAccountIDs []string      `json:"insights_account_ids"` // Array<String>
	Balance            domain.Amount `json:"balance"`              // ?Amount
	Available          domain.Amount `json:"available"`            // ?Amount
	Limit              domain.Amount `json:"limit"`                // ?Amount
	Reserved           domain.Amount `json:"reserved"`             // ?Amount,
	RefreshedAt        string        `json:"refreshed_at"`         // DateTime
}
