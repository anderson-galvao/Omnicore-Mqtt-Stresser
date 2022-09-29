package utils

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func QueryParamInt(c echo.Context, name string) (result int) {
	param := c.QueryParam(name)
	result, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return result
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
