package model

type Stresser struct {
	Clients  int `json:"clients"  validate:"required"`
	Messages int `json:"messages"  validate:"required"`
}
