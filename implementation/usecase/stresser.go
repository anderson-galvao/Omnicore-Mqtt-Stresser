package usecase

import (
	"github.com/RacoWireless/Omnicore-Mqtt-Stresser/implementation/utils"
	"github.com/RacoWireless/Omnicore-Mqtt-Stresser/model"
)

func (i *Usecase) ExecuteStresser(Arguments model.Stresser) (dr model.Response, err error) {
	err = i.StresserService.ExecuteStresser(Arguments, "epsi")
	err = i.StresserService.ExecuteStresser(Arguments, "eliance")
	err = i.StresserService.ExecuteStresser(Arguments, "ooing")
	if err != nil {
		dr = utils.FrameGenericResponse(500, model.SERVERERROR, "")
		return dr, err

	}
	dr = utils.FrameGenericResponse(200, "Success", "")
	return dr, nil
}
