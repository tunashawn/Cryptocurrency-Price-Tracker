package controllers

import (
	"backend/price-tracker/services/ws"
	"testing"
)

func TestNewPriceTrackingWorker(t *testing.T) {
	b := ws.BinanceWebsocketImpl{}
	NewPriceTrackingWorker(&b)
}
