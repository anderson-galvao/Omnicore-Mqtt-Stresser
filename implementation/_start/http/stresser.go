package http

import (
	"net/http"

	"github.com/RacoWireless/iot-gw-stresser/implementation/utils"
	"github.com/RacoWireless/iot-gw-stresser/model"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Stress godoc
// @Summary      Stress Mqtt
// @Description  Api For Stressing
// @Tags         stress
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GenericResponse
// @Failure      400  {object}  model.GenericResponse
// @Failure      404  {object}  model.GenericResponse
// @Failure      500  {object}  model.GenericResponse
// @Router       /stress [post]
func (r *Handler) ExecuteStresser(c echo.Context) error {
	//ctx := c.Request().Context()
	req := new(model.Stresser)
	if err := c.Bind(req); err != nil {
		log.Error().Err(err).Msg("Error in Binding Request")
		r := utils.FrameGenericResponse(400, model.INVALIDJSON, err.Error())
		return c.JSON(r.StatusCode, r.Message)
	}
	if err := c.Validate(req); err != nil {
		r := utils.FrameGenericResponse(400, model.INVALIDJSON, err.Error())
		return c.JSON(http.StatusBadRequest, r.Message)
	}
	mResponse, err := r.Usecase.ExecuteStresser(*req)

	if mResponse.StatusCode != 200 {
		log.Error().Err(err).Msg("")
		return c.JSON(mResponse.StatusCode, mResponse.Message)
	}
	return c.JSON(http.StatusOK, mResponse.Message)
}
