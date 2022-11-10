package service

import (
	"github.com/RacoWireless/Omnicore-Mqtt-Stresser/model"
)

type StresserService struct {
	BrokerUrl string
}

func NewStresserService(BrokerUrl string) model.StresserService {
	return &StresserService{
		BrokerUrl: BrokerUrl,
	}
}
