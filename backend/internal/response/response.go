package response

import (
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
)

type CustomResponse struct {
}

type ResponseData struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *CustomResponse) Success(data interface{}, ctx *gin.Context) {
	res := ResponseData{
		Meta: Meta{
			Code:    http.StatusOK,
			Message: "ok",
		},
		Data: data,
	}

	ctx.JSON(http.StatusOK, res)

	c.logging(res, slog.LevelInfo, ctx.Request)
}

func (c *CustomResponse) BadRequest(err error, ctx *gin.Context) {
	res := ResponseData{
		Meta: Meta{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		},
	}

	ctx.JSON(http.StatusBadRequest, res)

	c.logging(res, slog.LevelError, ctx.Request)
}

func (c *CustomResponse) InternalServerError(err error, ctx *gin.Context) {
	res := ResponseData{
		Meta: Meta{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}

	ctx.JSON(http.StatusInternalServerError, res)

	c.logging(res, slog.LevelError, ctx.Request)
}

func (c *CustomResponse) Unauthorized(err error, ctx *gin.Context) {
	res := ResponseData{
		Meta: Meta{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		},
	}

	ctx.JSON(http.StatusUnauthorized, res)

	c.logging(res, slog.LevelError, ctx.Request)
}

func (c *CustomResponse) logging(res ResponseData, logLevel slog.Level, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	slog.LogAttrs(
		context.Background(),
		logLevel,
		res.Meta.Message,
		slog.Group("request",
			"method", r.Method,
			"uri", r.RequestURI,
			"request_body", body),
		slog.Group("response",
			"code", res.Meta.Code,
			"message", res.Meta.Message,
			"data", res.Data),
	)
}
