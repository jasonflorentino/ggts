package gotrans

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

// Fetches Union Station's destinations as the default list since it is
// central hub through which GO Trains connect.
// Union is a Rail station only so there will not be any bus destinations.
// This list won't include Union Station itself so we should add it to complete the list.
func FetchDestinationsDefault(date string) (Destinations, error) {
	destinations, err := FetchDestinations(StationCode.Union, date)
	if err != nil {
		return nil, err
	}
	unionIdx := destinations.IndexOfCode(Union.Code)
	if unionIdx == -1 {
		destinations = append(destinations, Union)
		destinations.Sort()
	}
	return destinations, nil
}
