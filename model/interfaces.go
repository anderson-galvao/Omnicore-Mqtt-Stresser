package model

import "context"

type Usecase interface {
	HealthCheck(ctx context.Context) (Response, error)
}
type Service interface {
	HealthCheck(ctx context.Context) (Response, error)
}
