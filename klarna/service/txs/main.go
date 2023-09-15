package txs

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mehix/klarna-go/klarna"
	"github.com/mehix/klarna-go/klarna/domain/report"
	"github.com/mehix/klarna-go/klarna/domain/txs"
	"golang.org/x/exp/slices"
)

type Service struct {
	klarnaCli *klarna.Client
}

func NewService(klarnaCli *klarna.Client) *Service {
	return &Service{klarnaCli: klarnaCli}
}

func (s *Service) FetchAll(ctx context.Context, insightsConsumerID string) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchAllCredit(ctx context.Context, insightsConsumerID string) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.TransactionType = "CREDIT"

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchAllDebit(ctx context.Context, insightsConsumerID string) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.TransactionType = "DEBIT"

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchLatest(ctx context.Context, insightsConsumerID string, latest int64) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.Size = latest

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchLatestCredit(ctx context.Context, insightsConsumerID string, latest int64) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.Size = latest
	r.TransactionType = "CREDIT"

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchLatestDebit(ctx context.Context, insightsConsumerID string, latest int64) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.Size = latest
	r.TransactionType = "DEBIT"

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchDebitForPeriod(ctx context.Context, insightsConsumerID string, fromDate, toDate string) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.TransactionType = "DEBIT"
	r.ReportDays = 0
	r.FromDate = fromDate
	r.ToDate = toDate

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchLastDays(ctx context.Context, insightsConsumerID string, days int64) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.ReportDays = days

	return s.requestTransactions(ctx, r)
}

// ReportDailySpending returns aggregated spendings per day.
// Fetches transactions from Klarna and then processes them.
// It will ignore from the calculations transactions towards any IBAN in `ignoreIbans`
func (s *Service) ReportDailySpending(
	ctx context.Context,
	insightsConsumerID string,
	ignoreIbans ...string) ([]report.DailySpending, error) {

	transactions, err := s.FetchAllDebit(ctx, insightsConsumerID)
	if err != nil {
		return nil, err
	}

	return dailySpending(transactions, ignoreIbans...)
}

func (s *Service) ReportDailySpendingForPeriod(
	ctx context.Context,
	insightsConsumerID string,
	fromDate, toDate string,
	ignoreIbans ...string,
) ([]report.DailySpending, error) {

	transactions, err := s.FetchDebitForPeriod(ctx, insightsConsumerID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return dailySpending(transactions, ignoreIbans...)
}

func dailySpending(transactions []txs.CategorizedTransaction, ignoreIbans ...string) ([]report.DailySpending, error) {

	daily := make(map[string][5]int64)
	dailyTxs := make(map[string]map[string][]txs.CategorizedTransaction)

	for _, t := range transactions {
		if slices.Contains(ignoreIbans, t.CounterParty.IBAN) {
			continue
		}
		amounts, ok := daily[t.BookingDate]
		if !ok {
			amounts = [5]int64{}
		}
		dt, ok := dailyTxs[t.BookingDate]
		if !ok {
			dt = make(map[string][]txs.CategorizedTransaction)
		}

		if t.IsDebit() {
			amounts[0] += t.Amount.Amount
			part, ok := t.PartOfDay()
			if ok {
				switch part {
				case "morning":
					amounts[1] += t.Amount.Amount
				case "afternoon":
					amounts[2] += t.Amount.Amount
				case "evening":
					amounts[3] += t.Amount.Amount
				case "night":
					amounts[4] += t.Amount.Amount
				}
			}
			dt[part] = append(dt[part], t)
		}

		daily[t.BookingDate] = amounts
		dailyTxs[t.BookingDate] = dt
	}

	var rep []report.DailySpending
	for date, amounts := range daily {
		rep = append(rep, report.DailySpending{
			Date:         date,
			Debit:        amounts[0],
			Mornings:     amounts[1],
			Afternoons:   amounts[2],
			Evenings:     amounts[3],
			Nights:       amounts[4],
			Transactions: dailyTxs[date],
		})
	}

	sort.Slice(rep, func(i, j int) bool {
		return rep[i].Date > rep[j].Date
	})

	return rep, nil
}

// ReportMonthlyCreditBalance returns aggregated spendings per month.
// Fetches transactions from Klarna and then processes them.
// It will ignore from the calculations transactions towards any IBAN in `ignoreIbans`
func (s *Service) ReportMonthlyCreditBalance(
	ctx context.Context,
	filter txs.Filter,
	ignoreIbans ...string) ([]report.MonthlyCreditDebit, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = filter.InsightsConsumerID
	r.FromDate = filter.FromDate
	r.ToDate = filter.ToDate

	transactions, err := s.requestTransactions(ctx, r)
	if err != nil {
		return nil, err
	}

	monthly := make(map[string][2]int64)

	for _, t := range transactions {
		if slices.Contains(ignoreIbans, t.CounterParty.IBAN) {
			continue
		}
		key := string(t.BookingDate[:7])
		amounts, ok := monthly[key]
		if !ok {
			amounts = [2]int64{}
		}
		if t.IsCredit() {
			amounts[0] += t.Amount.Amount
		} else if t.IsDebit() {
			amounts[1] += t.Amount.Amount
		}
		monthly[key] = amounts
	}

	var rep []report.MonthlyCreditDebit
	for month, amounts := range monthly {
		rep = append(rep, report.MonthlyCreditDebit{
			Month:   month,
			Credit:  amounts[0],
			Debit:   amounts[1],
			Balance: amounts[0] - amounts[1],
		})
	}

	sort.Slice(rep, func(i, j int) bool {
		return rep[i].Month > rep[j].Month
	})

	return rep, nil
}

func (s *Service) requestTransactions(ctx context.Context, r txs.Request) ([]txs.CategorizedTransaction, error) {
	payload, err := json.Marshal(r)

	if err != nil {
		return nil, fmt.Errorf("[requestTransactions] marshal payload: %w", err)
	}

	b, err := s.klarnaCli.Post(ctx, "/insights/v1/reports/categorization/create", payload)
	if err != nil {
		return nil, err
	}

	var resp txs.Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("[requestTransactions] cannot unmarshal '%v': %w", string(b), err)
	}

	if len(resp.Data.Reports) == 0 {
		if resp.Error.Code != "" {
			return nil, resp.Error
		}
		return nil, fmt.Errorf("no transactions found")
	}

	return resp.Data.Reports[0].Transactions, nil
}
