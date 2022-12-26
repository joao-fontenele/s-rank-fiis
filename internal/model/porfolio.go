package model

import "sort"

type ByTopSRank []PortfolioElement

func (a ByTopSRank) Len() int           { return len(a) }
func (a ByTopSRank) Less(i, j int) bool { return a[i].TopSRankPosition < a[j].TopSRankPosition }
func (a ByTopSRank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type PortfolioElement struct {
	AvgPrice float64 `json:"avgPrice"`
	Amount   int
	Code     string

	CurrentPrice             float64 `json:"currentPrice"`
	DividendYield            float64 `json:"dividendYield"`
	TopSRankPosition         int     `json:"topSRankPosition"`
	SRanking                 int     `json:"sRanking"`
	DailyLiquidityInCurrency float64 `json:"dailyLiquidityInCurrency"`
}

func SortBySRankPosition(folio []PortfolioElement) []PortfolioElement {
	sort.Sort(ByTopSRank(folio))
	return folio
}
