package usecase

import (
	"context"

	"github.com/RacoWireless/iot-gw-stresser/model"
)

func (i *healthUsecase) ExecuteStresser(ctx context.Context) (model.Response, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()
	dr, err := i.healthRepo.HealthCheck(ctx)
	if err != nil {

		return dr, err

	}
	return dr, nil
}
