package repository

import (
	"strings"
	"time"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"gorm.io/gorm"
)

const stockTableName = "stocks"
const batchSize = 300

type Stock struct {
	Conn *gorm.DB
}

func (s Stock) SaveAll(stocks []model.Stock) error {
	tx := s.Conn.Table(stockTableName).CreateInBatches(stocks, batchSize)
	return tx.Error
}

func (s Stock) formatDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

func (s Stock) Find(f model.Filter) (stocks []model.Stock, err error) {
	var date time.Time
	if !f.Date.IsZero() {
		date = s.formatDate(f.Date)
	} else {
		now := time.Now()
		date = s.formatDate(now)
	}

	fields := []string{"created_at>?"}
	values := []interface{}{date}

	if f.PositivePBRatio == true {
		fields = append(fields, "pb_ratio>?")
		values = append(values, 0)
	}

	if f.GreaterThanDailyLiquidityInCurrency > 0 {
		fields = append(fields, "daily_liquidity_in_currency>?")
		values = append(values, f.GreaterThanDailyLiquidityInCurrency)
	}

	tx := s.Conn.Table(stockTableName).
		Where(strings.Join(fields, " AND "), values...).
		Order("s_ranking ASC").
		Find(&stocks)

	err = tx.Error
	return
}
