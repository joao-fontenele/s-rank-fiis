package translator

import (
	"log"
	"strconv"
	"strings"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
)

func transformNumberToFloat(s string) float64 {
	b := strings.Builder{}
	for _, c := range s {
		switch c {
		case ' ':
		case '%':
		case 'R':
		case '$':
		case '.':

		case ',':
			b.WriteRune('.')
		default:
			b.WriteRune(c)
		}
	}

	val, err := strconv.ParseFloat(b.String(), 64)
	if err != nil {
		log.Printf("failed to parse '%s' to float64: %v", s, err)
	}

	return val
}

func transformNumberToInt(s string) int {
	val, err := strconv.Atoi(strings.TrimSuffix(s, ".0"))
	if err != nil {
		log.Printf("failed to parse '%s' to int: %v", s, err)
	}

	return val
}

func Translate(raw map[string]string) (s model.Stock) {
	s.Code = raw["Código do fundo"]
	s.Sector = raw["Setor"]
	s.CurrentPrice = transformNumberToFloat(raw["Preço Atual"])
	s.DividendYield = transformNumberToFloat(raw["Dividend Yield"])
	s.DailyNegotiations = transformNumberToInt(raw["Liquidez Diária"])
	s.LastDividend = transformNumberToFloat(raw["Dividendo"])
	s.PBRatio = transformNumberToFloat(raw["P/VPA"])
	s.AmountOfProperties = transformNumberToInt(raw["Quantidade Ativos"])
	s.DailyLiquidityInCurrency = float64(s.DailyNegotiations) * s.CurrentPrice

	return
}
