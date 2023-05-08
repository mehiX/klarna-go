package txs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mehix/klarna-go/klarna"
	"github.com/mehix/klarna-go/klarna/domain/txs"
)

type Service struct {
	KlarnaCli *klarna.Client
}

func NewService(klarnaCli *klarna.Client) *Service {
	return &Service{KlarnaCli: klarnaCli}
}

func (s *Service) FetchAll(ctx context.Context, insightsConsumerID string) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID

	return s.requestTransactions(ctx, r)
}

func (s *Service) FetchLatest(ctx context.Context, insightsConsumerID string, latest int64) ([]txs.CategorizedTransaction, error) {
	r := txs.DefaultRequest
	r.InsightsConsumerID = insightsConsumerID
	r.Size = latest

	return s.requestTransactions(ctx, r)
}

func (s *Service) requestTransactions(ctx context.Context, r txs.Request) ([]txs.CategorizedTransaction, error) {
	payload, err := json.Marshal(r)

	if err != nil {
		return nil, fmt.Errorf("[requestTransactions] marshal payload: %w", err)
	}

	b, err := s.KlarnaCli.Post(ctx, "/insights/v1/reports/categorization/create", payload)
	if err != nil {
		return nil, err
	}

	var resp txs.Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("[requestTransactions] cannot unmarshal '%v': %w", string(b), err)
	}

	if len(resp.Data.Reports) == 0 {
		return nil, fmt.Errorf("no transactions found")
	}

	return resp.Data.Reports[0].Transactions, nil
}
