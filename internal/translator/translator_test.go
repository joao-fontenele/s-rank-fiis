package translator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joao-fontenele/s-rank-fiis/internal/model"
)

func Test_transforNumberToFloat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "works with percentages",
			args: args{s: "1,4%"},
			want: 1.4,
		},
		{
			name: "works with numbers",
			args: args{s: "9,95"},
			want: 9.95,
		},
		{
			name: "works with money",
			args: args{s: "R$ 992.187.188.537,90"},
			want: 992187188537.90,
		},
		{
			name: "works with spaces",
			args: args{s: " 537,90 "},
			want: 537.90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformNumberToFloat(tt.args.s); got != tt.want {
				t.Errorf("transforPercentageToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslate(t *testing.T) {
	type args struct {
		raw map[string]string
	}
	tests := []struct {
		name  string
		args  args
		wantS model.Stock
	}{
		{
			name: "works",
			args: args{
				raw: map[string]string{
					"Código do fundo":          "HGRU11",
					"DY (12M) Acumulado":       "7,94%",
					"DY (12M) Média":           "0,66%",
					"DY (3M) Acumulado":        "2,12%",
					"DY (3M) Média":            "0,71%",
					"DY (6M) Acumulado":        "4,09%",
					"DY (6M) Média":            "0,68%",
					"DY Ano":                   "2,78%",
					"DY Patrimonial":           "0,69%",
					"Dividend Yield":           "0,71%",
					"Dividendo":                "R$ 0,82",
					"Liquidez Diária":          "34863.0",
					"P/VPA":                    "0,96",
					"Patrimônio Líq.":          "R$ 2.187.188.537,90",
					"Preço Atual":              "R$ 114,02",
					"Quantidade Ativos":        "17",
					"Rentab. Acumulada":        "3,43%",
					"Rentab. Patr. Acumulada":  "2,06%",
					"Rentab. Patr. no Período": "0,84%",
					"Rentab. Período":          "3,95%",
					"Setor":                    "Híbrido",
					"VPA":                      "R$ 118,83",
					"Vacância Financeira":      "0,00%",
					"Vacância Física":          "0,00%",
					"Variação Patrimonial":     "0,15%",
					"Variação Preço":           "3,22%",
				},
			},
			wantS: model.Stock{
				Code:                     "HGRU11",
				DividendYield:            0.71,
				DailyNegotiations:        34863.0,
				DailyLiquidityInCurrency: 3975079.26,
				LastDividend:             0.82,
				CurrentPrice:             114.02,
				Sector:                   "Híbrido",
				PBRatio:                  0.96,
				AmountOfProperties:       17,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS := Translate(tt.args.raw)
			if diff := cmp.Diff(tt.wantS, gotS); diff != "" {
				t.Errorf("Found link different from expected (-want +got):\n%s", diff)
			}
		})
	}
}
