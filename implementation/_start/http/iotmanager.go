package http

import (
	"github.com/RacoWireless/iot-gw-stresser/model"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Usecase model.Usecase
}

func NewIoTtHandler(e *echo.Echo, Usecase model.Usecase) {
	Handler := &Handler{
		Usecase: Usecase,
	}

	e.GET("/health", Handler.HealthCheck)

}
