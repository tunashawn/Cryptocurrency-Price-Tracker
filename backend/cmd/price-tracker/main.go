package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/price-tracker/controllers"
	"backend/price-tracker/repositories"
	"backend/price-tracker/services"
	ws2 "backend/price-tracker/services/ws"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func main() {
	// Config
	cache := sync.Map{}

	var sqliteCfg config.SqliteConfig
	err := config.GetConfig(&sqliteCfg)
	if err != nil {
		log.Fatal(err)
	}

	var cfg config.PriceTrackerConfig
	err = config.GetConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Repository
	sqliteDB, err := db.NewSqliteDB(sqliteCfg)
	if err != nil {
		log.Fatalf("connecting to sqlite database: %v", err)
	}
	repository := repositories.NewSqliteRepository(sqliteDB)

	// Worker
	ws := ws2.NewBinanceWebsocket(&cache, cfg.BinanceWebSocketURL, cfg.WsFetchInterval, repository)
	worker := controllers.NewPriceTrackingWorker(ws)

	// Run worker to fetch data automatically
	worker.Run()

	// Service
	priceTrackingService := services.NewPriceTrackingService(&cache, repository)

	// Controller
	controller := controllers.NewPriceTrackerController(priceTrackingService)

	// Router
	r := gin.Default()

	r.GET("/list/name", controller.GetCryptoList)

	price := r.Group("/price/")
	{
		price.GET("/latest", controller.GetLatestPrice)
		price.GET("/interval", controller.GetPriceHistory)
	}

	r.Run()
}
