package usecase

import (
	"github.com/RacoWireless/iot-gw-mqtt-stresser/implementation/utils"
	"github.com/RacoWireless/iot-gw-mqtt-stresser/model"
)

func (i *Usecase) ExecuteStresser(Arguments model.Stresser) (dr model.Response, err error) {
	err = i.StresserService.ExecuteStresser(Arguments, "Pepsi")
	err = i.StresserService.ExecuteStresser(Arguments, "Cola")
	if err != nil {
		dr = utils.FrameGenericResponse(500, model.SERVERERROR, "")
		return dr, err

	}
	dr = utils.FrameGenericResponse(200, "Success", "")
	return dr, nil
}
