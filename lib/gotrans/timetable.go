package gotrans

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gogotrainschedule/lib/log"
	"io"
	"net/http"
	"net/url"
	"time"

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

	trips, err := filterTrips(timetable.Trips)
	if err != nil {
		return Timetable{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error filtering trips: %s\n", err),
		)
	}
	timetable.Trips = trips

	return timetable, nil
}

// Filters only trips
// - that haven't happened yet
// - are rail
// - are direct
func filterTrips(trips []Trip) ([]Trip, error) {
	now := time.Now()
	i := 0
	for _, trip := range trips {
		tripTime, err := time.ParseInLocation("2006-01-02T15:04:05", trip.OrderTime, time.Local)
		if err != nil {
			return nil, err
		}
		if tripTime.After(now) &&
			trip.TransitType == TransitTypes.Rail &&
			trip.Transfers == 0 {
			trips[i] = trip
			i++
		}
	}
	return trips[:i], nil
}
