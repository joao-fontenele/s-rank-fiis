package service

import (
	"math"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"go.uber.org/zap"
)

type Portfolio struct {
	Logger  *zap.Logger
	OpRepo  repository.Operation
	Ranking Ranking
}

func (p Portfolio) Summary(filter model.Filter) ([]model.PortfolioElement, error) {
	ops, err := p.OpRepo.FindAll()
	var folio []model.PortfolioElement
	if err != nil {
		return folio, err
	}

	m := make(map[string]model.PortfolioElement)
	for _, op := range ops {
		curr, ok := m[op.Code]
		if !ok {
			m[op.Code] = model.PortfolioElement{Code: op.Code}
			curr = m[op.Code]
		}

		curr.AvgPrice += op.Price * float64(op.Amount)
		curr.Amount += op.Amount
		m[op.Code] = curr
	}

	stocks, err := p.Ranking.Rank(filter)
	if err != nil {
		return folio, err
	}

	for i, s := range stocks {
		f, ok := m[s.Code]
		if !ok {
			continue
		}

		f.DividendYield = s.DividendYield
		f.CurrentPrice = s.CurrentPrice
		f.SRanking = s.SRanking
		f.TopSRankPosition = i + 1
		f.DailyLiquidityInCurrency = s.DailyLiquidityInCurrency

		m[s.Code] = f
	}

	for _, op := range m {
		if op.Amount <= 0 {
			p.Logger.Warn("division by 0 in calculation of avg price, skipping", zap.String("code", op.Code))
			continue
		}

		op.AvgPrice /= float64(op.Amount)
		if op.TopSRankPosition == 0 {
			op.TopSRankPosition = math.MaxInt
		}

		folio = append(folio, op)
	}

	folio = model.SortBySRankPosition(folio)

	return folio, nil
}
