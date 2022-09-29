package service

import (
	"github.com/RacoWireless/iot-gw-stresser/model"
)

type Service struct {
	Tenant string
}

func NewService(Tenant string) model.Service {
	return &Service{
		Tenant: Tenant,
	}
}
