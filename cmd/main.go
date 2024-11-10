package main

import (
	"gogotrainschedule/lib/gotransit"
	"gogotrainschedule/lib/log"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Page struct {
	Destinations gotransit.Destinations
	ErrorCode    int
	ErrorMessage string
	Timetable    gotransit.Timetable
	From         gotransit.Destination
	To           gotransit.Destination
}

func NewPage() Page {
	return Page{
		Timetable: gotransit.Timetable{},
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func defaultDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	log.Error(err)
	page := NewPage()
	page.ErrorCode = code
	page.ErrorMessage = http.StatusText(code)
	c.Render(code, "error", page)
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func isStartOfDay() bool {
	now := time.Now()
	hour := now.Hour()
	return hour < 13
}

func defaultFrom() gotransit.Destination {
	defaultStop := gotransit.Union
	if isStartOfDay() {
		defaultStop = gotransit.WestHarbour
	}
	return defaultStop
}

func defaultTo() gotransit.Destination {
	defaultStop := gotransit.WestHarbour
	if isStartOfDay() {
		defaultStop = gotransit.Union
	}
	return defaultStop
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := NewPage()
	e.Renderer = newTemplate()

	// destinations := make(gotransit.Destinations, 1)
	destinations, err := gotransit.FetchDestinations(gotransit.StationCode.Union, defaultDate())
	if err != nil {
		e.Logger.Fatal(err)
	}
	page.Destinations = destinations

	e.GET("/", func(c echo.Context) error {
		fromStop := defaultIfEmpty(c.QueryParam("fromStop"), defaultFrom().Code)
		toStop := defaultIfEmpty(c.QueryParam("toStop"), defaultTo().Code)
		date := defaultIfEmpty(c.QueryParam("date"), defaultDate())

		timetable, err := gotransit.FetchTimetable(fromStop, toStop, date)
		if err != nil {
			return err
		}
		timetable.Trips, err = filterTrips(timetable.Trips)
		if err != nil {
			return err
		}
		page.Timetable = timetable

		var from gotransit.Destination
		fromIdx := destinations.IndexOfCode(fromStop)
		if fromIdx == -1 {
			from = defaultFrom()
		} else {
			from = destinations[fromIdx]
		}
		page.From = from

		var to gotransit.Destination
		toIdx := destinations.IndexOfCode(toStop)
		if toIdx == -1 {
			to = defaultTo()
		} else {
			to = destinations[toIdx]
		}
		page.To = to

		return c.Render(http.StatusOK, "index", page)
	})

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.Fatal(e.Start("localhost:42069"))
}

func filterTrips(trips []gotransit.Trip) ([]gotransit.Trip, error) {
	now := time.Now()
	i := 0
	for _, trip := range trips {
		tripTime, err := time.ParseInLocation("2006-01-02T15:04:05", trip.OrderTime, time.Local)
		if err != nil {
			return nil, err
		}
		if tripTime.After(now) {
			trips[i] = trip
			i++
		}
	}
	return trips[:i], nil
}
