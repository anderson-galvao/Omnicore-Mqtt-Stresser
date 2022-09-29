package usecase

import (
	"time"

	"github.com/RacoWireless/iot-gw-stresser/model"
)

type healthUsecase struct {
	healthRepo     model.Service
	contextTimeout time.Duration
}

func NewHealthUsecase(Service model.Service, timeout time.Duration) model.Usecase {
	return &healthUsecase{
		Service:        Service,
		contextTimeout: timeout,
	}
}
