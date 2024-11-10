package gotransit

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gogotrainschedule/lib/log"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// date: "YYYY-MM-DD"
func FetchDestinations(destinationCode, date string) (Destinations, error) {
	req, err := Request(fmt.Sprintf("/v2/schedules/stops/%s/destinations?Date=%s", destinationCode, date))
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error creating http request: %s\n", err),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error sending http request: %s\n", err),
		)
	}
	log.Infof("Got response - Status: %d, ContentLength: %d", res.StatusCode, res.ContentLength)

	var body []byte
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	} else {
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	log.Infof("Body: %s", string(body))

	var destinations Destinations
	if err := json.Unmarshal(body, &destinations); err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not unmarshal json: %s\n", err),
		)
	}
	return destinations.OnlyRail(), nil
}
