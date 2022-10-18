package router

import (
	"context"

	"github.com/labstack/echo/v4"
	"hypefast.io/services/config"
	apiShorter "hypefast.io/services/domain/shorter/handler/api"
	"hypefast.io/services/domain/shorter/usecase"
)

func SetupRouterShorter(ctx context.Context, e *echo.Echo, conf config.Configuration) {

	shorterUseCase := usecase.NewShorterUseCase(ctx, conf)

	shorter := apiShorter.InitHandlerAPI(ctx, conf, shorterUseCase)

	gr := e.Group("/shorter")
	gr.POST("/set", shorter.CreateURL)
	gr.GET("/:key", shorter.GetURL)
}
