package model

type Usecase interface {
	ExecuteStresser(Stresser) (Response, error)
}
type StresserService interface {
	ExecuteStresser(Stresser) error
}
