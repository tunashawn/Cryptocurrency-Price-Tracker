package controllers

import (
	"backend/price-tracker/services/ws"
	"log/slog"
	"time"
)

type PriceTrackingWorker interface {
	// Run starts the worker and keeps retrying if errors occurs
	Run()
}

type PriceTrackingWorkerImpl struct {
	ws ws.WebSocketFetcher
}

func NewPriceTrackingWorker(ws ws.WebSocketFetcher) PriceTrackingWorker {
	return &PriceTrackingWorkerImpl{
		ws: ws,
	}
}

// Run starts the worker and keeps retrying if errors occurs
func (p *PriceTrackingWorkerImpl) Run() {
	go func() {
		for {
			time.Sleep(5 * time.Second)

			connection, err := p.ws.Connect()
			if err != nil {
				slog.Error("connect to websocket:", "error", err)
				continue
			}

			err = p.ws.Fetch(connection)
			if err != nil {
				slog.Error("worker fetch:", "error", err)
				continue
			}
		}
	}()
}
