package models

import "testing"

func TestBinanceResult_IsCurrencyUSDT(t *testing.T) {
	type fields struct {
		Symbol string
		Price  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "success",
			fields: fields{
				Symbol: "BTCUSDT",
				Price:  "123",
			},
			want: true,
		},
		{
			name: "success",
			fields: fields{
				Symbol: "BTCUSD",
				Price:  "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BinanceResult{
				Symbol: tt.fields.Symbol,
				Price:  tt.fields.Price,
			}
			if got := b.IsCurrencyUSDT(); got != tt.want {
				t.Errorf("IsCurrencyUSDT() = %v, want %v", got, tt.want)
			}
		})
	}
}
