package model

type Gateway struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
