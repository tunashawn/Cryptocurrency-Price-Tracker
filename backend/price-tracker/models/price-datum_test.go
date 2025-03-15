package models

import (
	"reflect"
	"testing"
	"time"
)

func TestNewPriceDatumFromBinanceResult(t *testing.T) {
	type args struct {
		res       BinanceResult
		timestamp time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    PriceDatum
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				res: BinanceResult{
					Symbol: "BTCUSDT",
					Price:  "123",
				},
			},
			want: PriceDatum{
				Symbol:   "BTC",
				Currency: "USDT",
				Price:    123,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				res: BinanceResult{
					Symbol: "BTCUSDT",
					Price:  "abc",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPriceDatumFromBinanceResult(tt.args.res, tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPriceDatumFromBinanceResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPriceDatumFromBinanceResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}
