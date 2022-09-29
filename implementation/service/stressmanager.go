package service

import (
	"github.com/RacoWireless/iot-gw-stresser/model"
)

type StresserService struct {
	BrokerUrl string
}

func NewStresserService(BrokerUrl string) model.StresserService {
	return &StresserService{
		BrokerUrl: BrokerUrl,
	}
}
