package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"github.com/joao-fontenele/s-rank-fiis/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	Logger *zap.Logger

	etlService     service.ETL
	rankingService service.Ranking
}

func New(logger *zap.Logger, stockRepo repository.Stock) Handler {
	return Handler{
		Logger: logger,
		etlService: service.ETL{
			Logger: logger,
			Repo:   stockRepo,
		},
		rankingService: service.Ranking{
			Logger: logger,
			Repo:   stockRepo,
		},
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("request received", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
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

		b, err := io.ReadAll(r.Body)
		if err != nil {
			h.Logger.Error("failed to read request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"failed to read request body"}`))
			return
		}

		err = json.Unmarshal(b, &f)
		if err != nil {
			h.Logger.Error("failed to parse request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"failed to parse request body"}`))
			return
		}

		if f.GreaterThanDailyLiquidityInCurrency < 0 {
			h.Logger.Error("bad request liquidity is < 0", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"liquidity should be greater than 0"}`))
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
