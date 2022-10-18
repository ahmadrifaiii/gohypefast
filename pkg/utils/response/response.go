package response

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Response struct {
	LogId    string      `json:"logId"`
	Status   int         `json:"status"`
	Message  interface{} `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Paginate interface{} `json:"paginate,omitempty"`
	Sort     interface{} `json:"sort,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

// handling response for success
func Success(c echo.Context, r Response) error {
	return c.JSON(r.Status, r)
}

// handling response for success
func Redirect(c echo.Context, r Response) error {
	return c.JSON(r.Status, r)
}

// handling response for error
func Error(c echo.Context, r Response) error {
	return c.JSON(r.Status, r)
}

func HttpErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				report.Message = fmt.Sprintf("%s is required",
					err.Field())
			case "email":
				report.Message = fmt.Sprintf("%s is not valid email",
					err.Field())
			case "gte":
				report.Message = fmt.Sprintf("%s value must be greater than %s",
					err.Field(), err.Param())
			case "lte":
				report.Message = fmt.Sprintf("%s value must be lower than %s",
					err.Field(), err.Param())
			}

			break
		}
	}

	c.Logger().Error(report)
	c.JSON(report.Code, report)
}
