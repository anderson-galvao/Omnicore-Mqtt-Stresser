package utils

import (
	"github.com/RacoWireless/Omnicore-Mqtt-Stresser/model"
)

func FrameGenericResponse(statusCode int, msg string, details string) model.Response {

	frame := model.Frame{StateCode: statusCode, Message: msg, Details: details}
	return model.Response{StatusCode: statusCode, Message: frame}
}
