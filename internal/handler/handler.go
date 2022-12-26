package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"github.com/joao-fontenele/s-rank-fiis/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger

	etlService       service.ETL
	rankingService   service.Ranking
	portfolioService service.Portfolio
}

func New(logger *zap.Logger, stockRepo repository.Stock, opRepo repository.Operation) Handler {
	rankingService := service.Ranking{
		Logger: logger,
		Repo:   stockRepo,
	}
	return Handler{
		Logger: logger,
		etlService: service.ETL{
			Logger: logger,
			Repo:   stockRepo,
		},
		rankingService: rankingService,
		portfolioService: service.Portfolio{
			Logger:  logger,
			OpRepo:  opRepo,
			Ranking: rankingService,
		},
	}
}

func (h Handler) parseFilterBody(r *http.Request) (model.Filter, error) {
	var f model.Filter
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return f, fmt.Errorf("failed to read request body: %w", err)
	}

	err = json.Unmarshal(b, &f)
	if err != nil {
		return f, fmt.Errorf("failed to parse request body: %w", err)
	}

	if f.GreaterThanDailyLiquidityInCurrency < 0 {
		return f, errors.New("liquidity should be greater than 0")
	}

	return f, nil
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("request received", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch r.URL.Path {
	case "/summary":
		var f model.Filter
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"not found"}`))
			return
		}

		f, err := h.parseFilterBody(r)
		if err != nil {
			h.Logger.Error("failed to parse request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}

		folio, err := h.portfolioService.Summary(f)
		if err != nil {
			h.Logger.Error("failed to get portfolio summary", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"failed to get summary!"}`))
			return
		}

		marshalled, err := json.Marshal(folio)
		if err != nil {
			h.Logger.Error("failed to marshal portfolio summary", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"failed to marshal summary!"}`))
			return
		}
		_, _ = w.Write(marshalled)
	case "/etl":
		err := h.etlService.Run()
		if err != nil {
			h.Logger.Error("failed to run ETL", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"failed to run ETL!"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"ETL Done!"}`))
		}
	case "/stocks":
		var f model.Filter
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"not found"}`))
			return
		}

		f, err := h.parseFilterBody(r)
		if err != nil {
			h.Logger.Error("failed to parse request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}

		h.Logger.Info(
			"applying filter to rank request",
			zap.Time("filter.date", f.Date),
			zap.Bool("filter.positivePBRation", f.PositivePBRatio),
			zap.Float64("filter.greaterThanDailyLiquidityInCurrency", f.GreaterThanDailyLiquidityInCurrency),
		)
		stocks, err := h.rankingService.Rank(f)
		if err != nil {
			h.Logger.Error("failed to rank", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"failed to rank stocks"}`))
			return
		}

		marshalled, err := json.Marshal(stocks)
		if err != nil {
			h.Logger.Error("failed to marshal stocks", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"failed to rank stocks"}`))
			return
		}
		_, _ = w.Write(marshalled)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	}
}
