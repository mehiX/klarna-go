package report

import "github.com/mehix/klarna-go/klarna/domain/txs"

type MonthlyCreditDebit struct {
	Month   string `json:"month"` // YYYY-MM
	Credit  int64  `json:"credit"`
	Debit   int64  `json:"debit"`
	Balance int64  `json:"balance"`
}

type DailySpending struct {
	Date         string                                  `json:"date"` // YYYY-MM-DD
	Debit        int64                                   `json:"debit"`
	Mornings     int64                                   `json:"mornings"`
	Afternoons   int64                                   `json:"afternoons"`
	Evenings     int64                                   `json:"evenings"`
	Nights       int64                                   `json:"nights"`
	Transactions map[string][]txs.CategorizedTransaction `json:"transactions"`
}
