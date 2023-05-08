package account

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mehix/klarna-go/klarna/domain/account"
)

func (s *Service) requestInfo(ctx context.Context, r account.RequestInfo) ([]account.Info, error) {
	payload, err := json.Marshal(r)

	if err != nil {
		return nil, fmt.Errorf("[requestInfo] marshal payload: %w", err)
	}

	b, err := s.KlarnaCli.Post(ctx, "/insights/v1/reports/account-info/create", payload)
	if err != nil {
		return nil, err
	}

	var resp account.ResponseInfo
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("[requestInfo] cannot unmarshal '%v': %w", string(b), err)
	}

	if len(resp.Data.Reports) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}

	return resp.Data.Reports, nil
}

func (s *Service) requestBalances(ctx context.Context, r account.RequestBalances) ([]account.Balances, error) {
	payload, err := json.Marshal(r)

	if err != nil {
		return nil, fmt.Errorf("[requestBalances] marshal payload: %w", err)
	}

	b, err := s.KlarnaCli.Post(ctx, "/insights/v1/reports/balances/create", payload)
	if err != nil {
		return nil, err
	}

	var resp account.ResponseBalances
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("[requestBalances] cannot unmarshal '%v': %w", string(b), err)
	}

	if len(resp.Data.Reports) == 0 {
		return nil, fmt.Errorf("no balances data found")
	}

	return resp.Data.Reports, nil
}

func (s *Service) requestBalanceOverTime(ctx context.Context, r account.RequestBalanceOverTime) ([]account.BalanceOverTime, error) {
	payload, err := json.Marshal(r)

	if err != nil {
		return nil, fmt.Errorf("[requestBalanceOverTime] marshal payload: %w", err)
	}

	b, err := s.KlarnaCli.Post(ctx, "/insights/v1/reports/balance-over-time/create", payload)
	if err != nil {
		return nil, err
	}

	var resp account.ResponseBalanceOverTime
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("[requestBalanceOverTime] cannot unmarshal '%v': %w", string(b), err)
	}

	if len(resp.Data.Reports) == 0 {
		return nil, fmt.Errorf("no balance over time data found")
	}

	return resp.Data.Reports, nil
}
