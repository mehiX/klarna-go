package account

import "github.com/mehix/klarna-go/klarna/domain"

type RequestBalanceOverTime struct {
	InsightsConsumerID         string   `json:"insights_consumer_id"`          // String
	InsightsAccountIDs         []string `json:"insights_account_ids"`          // ?Array<String>
	ExcludedInsightsAccountIDs []string `json:"excluded_insights_account_ids"` // ?Array<String>
	ReportDays                 int64    `json:"report_days,omitempty"`         // ?Integer
	FromDate                   string   `json:"from_date,omitempty"`           // ?Date
	ToDate                     string   `json:"to_date,omitempty"`             // ?Date
}

type ResponseBalanceOverTime struct {
	Data DataBalanceOverTime `json:"data"`
}

type DataBalanceOverTime struct {
	Reports []BalanceOverTime `json:"reports"`
}

type BalanceOverTime struct {
	Type               string        `json:"type"`                 // "BALANCE_OVER_TIME"
	InsightsConsumerID string        `json:"insights_consumer_id"` // String
	InsightsAccountIDs []string      `json:"insights_account_ids"` // Array<String>
	Balances           []BalanceInfo `json:"balances"`             // Array<BalanceInfo>
	FromDate           string        `json:"from_date"`            // Date
	ToDate             string        `json:"to_date"`              // Date
	RefreshedAt        string        `json:"refreshed_at"`         // DateTime
}

type BalanceInfo struct {
	Date    string        `json:"date"` // YYYY-MM-DD
	Balance domain.Amount `json:"balance"`
}
