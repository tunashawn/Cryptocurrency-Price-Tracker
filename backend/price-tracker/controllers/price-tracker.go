package controllers

import (
	"backend/internal/response"
	"backend/price-tracker/models"
	"backend/price-tracker/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type PriceTrackerController interface {
	// GetLatestPrice handles requests to get the latest price of an symbol
	GetLatestPrice(ctx *gin.Context)
	// GetPriceHistory handles requests to get the price history
	// current supports only 24h interval
	GetPriceHistory(ctx *gin.Context)
}

type PriceTrackerControllerImpl struct {
	priceTrackingService services.PriceTrackingService
	httpResponse         response.CustomResponse
}

func NewPriceTrackerController(priceTrackingService services.PriceTrackingService) PriceTrackerController {
	return &PriceTrackerControllerImpl{
		priceTrackingService: priceTrackingService,
	}
}

// GetLatestPrice handles requests to get the latest price of an symbol
func (p *PriceTrackerControllerImpl) GetLatestPrice(ctx *gin.Context) {
	symbol, ok := p.getSymbol(ctx)
	if !ok {
		return
	}

	req := models.PriceDatum{
		Symbol: symbol,
	}

	res, err := p.priceTrackingService.GetLatestPrice(req)
	if err != nil {
		p.httpResponse.InternalServerError(err, ctx)
		return
	}

	p.httpResponse.Success(res, ctx)
}

// GetPriceHistory handles requests to get the price history
// current supports only 24h interval
func (p *PriceTrackerControllerImpl) GetPriceHistory(ctx *gin.Context) {
	symbol, ok := p.getSymbol(ctx)
	if !ok {
		return
	}

	if _, ok = p.getInterval(ctx); !ok {
		return
	}

	req := models.PriceDatum{
		Symbol: symbol,
	}

	res, err := p.priceTrackingService.GetPriceOfTheLast24h(req)
	if err != nil {
		p.httpResponse.InternalServerError(err, ctx)
		return
	}

	p.httpResponse.Success(res, ctx)
}

// getSymbol returns symbol value from query params if exist
func (p *PriceTrackerControllerImpl) getSymbol(ctx *gin.Context) (string, bool) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		p.httpResponse.BadRequest(fmt.Errorf("missing symbol value '%s'", symbol), ctx)
		return "", false
	}

	if !p.priceTrackingService.IsSymbolValid(symbol) {
		p.httpResponse.BadRequest(fmt.Errorf("invalid symbol value '%s'", symbol), ctx)
		return "", false
	}

	return symbol, true
}

// getInterval returns interval value from query params if exist
func (p *PriceTrackerControllerImpl) getInterval(ctx *gin.Context) (string, bool) {
	interval := ctx.Query("interval")
	if interval == "" {
		p.httpResponse.BadRequest(errors.New("interval only support '24' hours"), ctx)
		return "", false
	}

	return interval, true
}
