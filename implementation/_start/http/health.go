package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Api For Health Check
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GenericResponse
// @Failure      400  {object}  model.GenericResponse
// @Failure      404  {object}  model.GenericResponse
// @Failure      500  {object}  model.GenericResponse
// @Router       /health [get]
func (r *Handler) HealthCheck(c echo.Context) error {
	ctx := c.Request().Context()

	mResponse, err := r.hUsecase.HealthCheck()

	if mResponse.StatusCode != 200 {
		log.Error().Err(err).Msg("")
		return c.JSON(mResponse.StatusCode, mResponse.Message)
	}
	return c.JSON(http.StatusOK, mResponse.Message)
}
