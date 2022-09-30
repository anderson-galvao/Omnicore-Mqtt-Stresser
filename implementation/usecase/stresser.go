package usecase

import (
	"github.com/RacoWireless/iot-gw-mqtt-stresser/implementation/utils"
	"github.com/RacoWireless/iot-gw-mqtt-stresser/model"
)

func (i *Usecase) ExecuteStresser(Arguments model.Stresser) (dr model.Response, err error) {
	go i.StresserService.ExecuteStresser(Arguments, "epsi")
	go i.StresserService.ExecuteStresser(Arguments, "eliance")
	go i.StresserService.ExecuteStresser(Arguments, "ooing")
	if err != nil {
		dr = utils.FrameGenericResponse(500, model.SERVERERROR, "")
		return dr, err

	}
	dr = utils.FrameGenericResponse(200, "Success", "")
	return dr, nil
}
