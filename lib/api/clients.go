package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Gotransit(c echo.Context, endpoint string) (*http.Request, error) {
	return request(c, GOTRANSIT_URL, endpoint)
}

func Metrolinx(c echo.Context, endpoint string) (*http.Request, error) {
	return request(c, METROLINX_URL, endpoint)
}
