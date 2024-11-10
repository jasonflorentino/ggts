package gotransit

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gogotrainschedule/lib/log"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// date: "YYYY-MM-DD"
func FetchTimetable(fromStop, toStop, date string) (Timetable, error) {
	params := url.Values{}
	params.Add("fromStop", fromStop)
	params.Add("toStop", toStop)
	params.Add("date", date)
	queryString := params.Encode()
	req, err := Request(fmt.Sprintf("/v2/schedules/en/timetable/all?%s", queryString))
	if err != nil {
		return Timetable{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error creating http request: %s", err),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Timetable{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error sending http request: %s", err),
		)
	}
	log.Infof("Got response - Status: %d, ContentLength: %d", res.StatusCode, res.ContentLength)

	var body []byte
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			return Timetable{}, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			return Timetable{}, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	} else {
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return Timetable{}, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	log.Infof("Body: %s", string(body))

	var timetable Timetable
	if err := json.Unmarshal(body, &timetable); err != nil {
		return Timetable{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not unmarshal json: %s\n", err),
		)
	}
	return timetable, nil
}
