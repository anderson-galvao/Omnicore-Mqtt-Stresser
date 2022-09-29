package utils

import (
	"github.com/RacoWireless/iot-gw-mqtt-stresser/model"
)

func FrameGenericResponse(statusCode int, msg string, details string) model.Response {

	frame := model.Frame{StateCode: statusCode, Message: msg, Details: details}
	return model.Response{StatusCode: statusCode, Message: frame}
}
