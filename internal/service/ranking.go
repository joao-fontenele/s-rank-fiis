package service

import (
	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"go.uber.org/zap"
)

type Ranking struct {
	Logger *zap.Logger
	Repo   repository.Stock
}

func (r Ranking) Rank(f model.Filter) ([]model.Stock, error) {
	return r.Repo.Find(f)
}
