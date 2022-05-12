package service

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"github.com/joao-fontenele/s-rank-fiis/internal/parser"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"github.com/joao-fontenele/s-rank-fiis/internal/translator"
	"go.uber.org/zap"
)

const (
	RequestTimeoutMilliseconds = 10000
	RankingURL                 = "https://www.fundsexplorer.com.br/ranking"
)

func requestHTML() ([]byte, error) {
	c := http.Client{
		Timeout: time.Duration(RequestTimeoutMilliseconds) * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, RankingURL, nil)
	if err != nil {
		return []byte{}, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type ETL struct {
	Logger *zap.Logger
	Repo   repository.Stock
}

func (e ETL) Run() (err error) {
	html, err := requestHTML()
	if err != nil {
		return err
	}

	p, err := parser.NewRankingTableParser(bytes.NewReader(html))
	if err != nil {
		return err
	}

	rawStocks, err := p.Parse()
	if err != nil {
		return err
	}

	stocks := make([]model.Stock, len(rawStocks))
	for i := 0; i < len(rawStocks); i++ {
		stocks[i] = translator.Translate(rawStocks[i])
	}

	// fill Rankings of stocks
	stocks = model.Rank(stocks)

	return e.Repo.SaveAll(stocks)
}
