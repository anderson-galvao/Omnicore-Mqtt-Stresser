package usecase

import (
	"github.com/RacoWireless/iot-gw-mqtt-stresser/model"
)

type Usecase struct {
	StresserService model.StresserService
}

func NewUsecase(Service model.StresserService) model.Usecase {
	return &Usecase{
		StresserService: Service,
	}
}
