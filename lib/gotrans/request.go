package gotrans

import (
	"gogotrainschedule/lib/log"
	"net/http"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
)

func Request(c echo.Context, endpoint string) (*http.Request, error) {
	if r, _ := utf8.DecodeRuneInString(endpoint); r != '/' {
		endpoint = "/" + endpoint
	}
	url := "https://" + API_URL + endpoint
	log.To(c).Infof("Creating request: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header = http.Header{
		"Accept":          {"*/*"},
		"Accept-Encoding": {"gzip"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Cache-Control":   {"no-cache"},
		"Connection":      {"keep-alive"},
		"Host":            {API_URL},
		"Origin":          {"https://www.gotransit.com"},
		"Pragma":          {"no-cache"},
		"Priority":        {"u=0"},
		"Referer":         {"https://www.gotransit.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}
