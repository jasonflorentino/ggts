package main

import (
	"gogotrainschedule/lib/gotrans"
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
	DestinationsFrom gotrans.Destinations
	DestinationsTo   gotrans.Destinations
	ErrorCode        int
	ErrorMessage     string
	Timetable        gotrans.Timetable
	From             gotrans.Destination
	To               gotrans.Destination
}

func NewPage() Page {
	return Page{
		Timetable: gotrans.Timetable{},
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

func defaultFrom() gotrans.Destination {
	defaultStop := gotrans.Union
	if isStartOfDay() {
		defaultStop = gotrans.WestHarbour
	}
	return defaultStop
}

func defaultTo() gotrans.Destination {
	defaultStop := gotrans.WestHarbour
	if isStartOfDay() {
		defaultStop = gotrans.Union
	}
	return defaultStop
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := NewPage()
	e.Renderer = newTemplate()

	// destinations := make(gotrans.Destinations, 1)
	destinations, err := gotrans.FetchDestinationsDefault(defaultDate())
	if err != nil {
		e.Logger.Fatal(err)
	}
	page.DestinationsFrom = destinations
	page.DestinationsTo = destinations

	e.GET("/to", func(c echo.Context) error {
		fromStop := c.QueryParam("fromStop")
		dests, err := gotrans.FetchDestinations(fromStop, defaultDate())
		if err != nil {
			e.Logger.Fatal(err)
		}
		page := NewPage()
		page.DestinationsTo = dests.SetSelected(dests[0].Code)
		return c.Render(http.StatusOK, "selectTo", page)
	})

	e.GET("/", func(c echo.Context) error {
		fromStop := defaultIfEmpty(c.QueryParam("fromStop"), defaultFrom().Code)
		toStop := defaultIfEmpty(c.QueryParam("toStop"), defaultTo().Code)
		date := defaultIfEmpty(c.QueryParam("date"), defaultDate())

		timetable, err := gotrans.FetchTimetable(fromStop, toStop, date)
		if err != nil {
			return err
		}
		page.Timetable = timetable

		var from gotrans.Destination
		fromIdx := destinations.IndexOfCode(fromStop)
		if fromIdx == -1 {
			from = defaultFrom()
		} else {
			from = destinations[fromIdx]
		}
		page.From = from

		var to gotrans.Destination
		toIdx := destinations.IndexOfCode(toStop)
		if toIdx == -1 {
			to = defaultTo()
		} else {
			to = destinations[toIdx]
		}
		page.To = to

		page.DestinationsFrom = page.DestinationsFrom.SetSelected(from.Code)
		page.DestinationsTo = page.DestinationsTo.SetSelected(to.Code)

		return c.Render(http.StatusOK, "index", page)
	})

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.Fatal(e.Start("localhost:42069"))
}
