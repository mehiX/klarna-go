package account

import (
	"context"
	"time"

	"github.com/mehix/klarna-go/klarna"
	"github.com/mehix/klarna-go/klarna/domain/account"
)

type Service struct {
	klarnaCli *klarna.Client
}

func NewService(kc *klarna.Client) *Service {
	return &Service{kc}
}

func (s *Service) Info(ctx context.Context, insightsConsumerID string, accountsID ...string) ([]account.Info, error) {

	r := account.RequestInfo{
		InsightsConsumerID: insightsConsumerID,
		InsightsAccountIDs: accountsID,
	}

	return s.requestInfo(ctx, r)
}

func (s *Service) Balances(ctx context.Context, insightsConsumerID string, accountsID ...string) ([]account.Balances, error) {
	r := account.RequestBalances{
		InsightsConsumerID:       insightsConsumerID,
		InsightsAccountIDs:       accountsID,
		RequiredDataAvailability: "NONE",
	}

	return s.requestBalances(ctx, r)
}

func (s *Service) BalanceOverTime(ctx context.Context, insightsConsumerID string, from, to time.Time, accountsID ...string) ([]account.BalanceOverTime, error) {
	r := account.RequestBalanceOverTime{
		InsightsConsumerID: insightsConsumerID,
		InsightsAccountIDs: accountsID,
		FromDate:           from.Format("2006-01-02"),
		ToDate:             to.Format("2006-01-02"),
	}

	return s.requestBalanceOverTime(ctx, r)
}
