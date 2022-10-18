package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"hypefast.io/services/config"
	shorterModel "hypefast.io/services/domain/shorter/model"
	shorterUserCase "hypefast.io/services/domain/shorter/usecase"
	"hypefast.io/services/pkg/utils/logging"
	"hypefast.io/services/pkg/utils/response"
)

type HandlerAPI struct {
	Ctx            context.Context
	Config         config.Configuration
	ShorterService shorterUserCase.ShorterUseCase
}

func InitHandlerAPI(ctx context.Context, conf config.Configuration, shorter shorterUserCase.ShorterUseCase) *HandlerAPI {
	return &HandlerAPI{
		Ctx:            ctx,
		Config:         conf,
		ShorterService: shorter,
	}
}

func (h *HandlerAPI) CreateURL(c echo.Context) error {
	var (
		requestId = c.Get("request_id").(string)
		payload   = shorterModel.Payload{}
	)

	err := c.Bind(&payload)
	if err != nil {
		return response.Error(c, response.Response{
			LogId:  requestId,
			Status: http.StatusBadRequest,
			Error:  err,
		})
	}

	result, err := h.ShorterService.SetURL(payload)
	if err != nil {
		return response.Error(c, response.Response{
			LogId:  requestId,
			Status: http.StatusBadRequest,
			Error:  err,
		})
	}

	return response.Success(c, response.Response{
		LogId:   requestId,
		Status:  http.StatusOK,
		Message: nil,
		Data:    result,
	})
}

func (h *HandlerAPI) GetURL(c echo.Context) error {
	var (
		requestId  = c.Get("request_id").(string)
		requestKey = c.Param("key")
	)

	url, err := h.ShorterService.GetURL(requestKey)
	if err != nil {
		logging.Error(err, zap.String("request_id", requestId))
		return response.Error(c, response.Response{
			LogId:   requestId,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		})
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}
