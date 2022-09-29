package usecase

import (
	"github.com/RacoWireless/iot-gw-stresser/implementation/utils"
	"github.com/RacoWireless/iot-gw-stresser/model"
)

func (i *Usecase) ExecuteStresser(Arguments model.Stresser) (model.Response, error) {
	err := i.StresserService.ExecuteStresser(Arguments)
	dr := utils.FrameGenericResponse(400, model.INVALIDCA, "")
	if err != nil {

		return dr, err

	}
	return dr, nil
}
