package model

import (
	"sort"
	"time"
)

type Filter struct {
	PositivePBRatio                     bool      `json:"positivePBRatio"`
	GreaterThanDailyLiquidityInCurrency float64   `json:"greaterThanDailyLiquidityInCurrency"`
	Date                                time.Time `json:"date"`
}

type Stock struct {
	Code                     string
	Sector                   string
	DividendYield            float64
	DailyLiquidityInCurrency float64
	LastDividend             float64
	CurrentPrice             float64
	PBRatio                  float64
	DailyNegotiations        int
	AmountOfProperties       int

	PBRatioRanking int
	DYRanking      int
	SRanking       int
	CreatedAt      time.Time
}

type ByPBRatio []Stock
type ByDY []Stock
type BySRank []Stock

func (a ByDY) Len() int           { return len(a) }
func (a ByDY) Less(i, j int) bool { return a[j].DividendYield < a[i].DividendYield }
func (a ByDY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a ByPBRatio) Len() int           { return len(a) }
func (a ByPBRatio) Less(i, j int) bool { return a[i].PBRatio < a[j].PBRatio }
func (a ByPBRatio) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a BySRank) Len() int           { return len(a) }
func (a BySRank) Less(i, j int) bool { return a[i].SRanking < a[j].SRanking }
func (a BySRank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func Rank(stocks []Stock) []Stock {
	sort.Sort(ByDY(stocks))
	for i := 0; i < len(stocks); i++ {
		stocks[i].DYRanking = i + 1
	}

	sort.Sort(ByPBRatio(stocks))
	for i := 0; i < len(stocks); i++ {
		stocks[i].PBRatioRanking = i + 1
		stocks[i].SRanking = stocks[i].DYRanking + stocks[i].PBRatioRanking
	}

	//sort.Sort(BySRank(stocks))
	return stocks
}
