package model

type Frame struct {
	StateCode int         `json:"stateCode"  validate:"required"`
	Message   string      `json:"message"  validate:"required"`
	Details   interface{} `json:"details"  validate:"required"`
}

// ///// response
type Response struct {
	StatusCode int   `json:"statuscode"  validate:"required"`
	Message    Frame `json:"message"  validate:"required"`
}
type GenericResponse struct {
	StateCode int    `json:"stateCode"  validate:"required"`
	Message   string `json:"message"  validate:"required"`
	Details   string `json:"details"  validate:"required"`
}
