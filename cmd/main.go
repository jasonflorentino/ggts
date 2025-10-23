package main

import (
	"bytes"
	"fmt"
	"ggts/lib"
	"ggts/lib/env"
	"ggts/lib/gotrans"
	"ggts/lib/log"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Page struct {
	DatePicker       lib.DatePicker
	DestinationsFrom gotrans.Destinations
	DestinationsTo   gotrans.Destinations
	ErrorCode        int
	ErrorMessage     string
	Timetable        gotrans.Timetable
	From             gotrans.Destination
	To               gotrans.Destination
	GGTS_TITLE       string
	GGTS_URL         string
}

func (p Page) String() string {
	return fmt.Sprintf(
		"DestinationsFrom: %v, DestinationsTo: %v, ErrorCode: %d, ErrorMessage: %s, Timetable: %v, From: %v, To: %v, GGTS_TITLE: %s, GGTS_URL: %s",
		p.DestinationsFrom,
		p.DestinationsTo,
		p.ErrorCode,
		p.ErrorMessage,
		p.Timetable,
		p.From,
		p.To,
		p.GGTS_TITLE,
		p.GGTS_URL,
	)
}

func NewPage(today string, selected string) Page {
	return Page{
		DatePicker: lib.NewDatePicker(today, selected),
		Timetable:  gotrans.Timetable{},
		GGTS_TITLE: env.Title(),
		GGTS_URL:   env.URL(),
	}
}

func defaultDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func isStartOfDay() bool {
	now := time.Now()
	return now.In(env.Location()).Hour() < 13
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

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	log.To(c).Error(err)
	page := NewPage(defaultDate(), defaultDate())
	page.ErrorCode = code
	page.ErrorMessage = http.StatusText(code)
	c.Render(code, "error", page)
}

func renderTemplates(c echo.Context, blocks []string, page Page) (string, error) {
	var buf bytes.Buffer
	var builder strings.Builder
	for _, block := range blocks {
		buf.Reset()
		if err := c.Echo().Renderer.Render(&buf, block, page, c); err != nil {
			return "", err
		}
		builder.WriteString(buf.String())
	}
	return builder.String(), nil
}

func main() {
	// Load runtime vars
	env.LoadEnv()
	gotrans.InitCache()

	// Init echo
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	// TODO: Ensure API won't be affected before implementing this:
	// e.Use(middleware.BodyLimit("2M"))

	if env.IsProd() {
		e.HideBanner = true
		e.Logger.SetOutput(log.ToFile(env.LogFile()))
	} else {
		e.Logger.SetOutput(os.Stdout)
	}
	e.Logger.SetLevel(log.Lvl())

	e.Logger.Infof("SERVER START")
	e.Renderer = newTemplate()
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Routes
	e.Static("/", "static")
	e.GET("/date-picker", handleDatePicker)
	e.GET("/trips", handleTrips)
	e.GET("/to", handleTo)
	e.GET("/", handleRoot)

	// Run
	e.Logger.Fatal(e.Start(":" + env.Port()))
}

// handleDatePicker responds with an updated datepicker.
func handleDatePicker(c echo.Context) error {
	changeDate, err := lib.GetChangeDate(c)
	if err != nil {
		return err
	}
	if changeDate == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing date parts")
	}

	page := NewPage(defaultDate(), changeDate)

	return c.Render(http.StatusOK, "datePicker", page)
}

// handleTrips responds with the timetable of trips.
func handleTrips(c echo.Context) error {
	fromStop := c.QueryParam("from")
	toStop := c.QueryParam("to")
	date := defaultIfEmpty(c.QueryParam("date"), defaultDate())

	if fromStop == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing from")
	}
	if toStop == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing to")
	}

	changeDate, err := lib.GetChangeDate(c)
	if err != nil {
		return err
	}
	if changeDate != "" {
		date = changeDate
	}

	page := NewPage(defaultDate(), date)
	c.Response().Header().Add(
		"HX-Push-Url",
		fmt.Sprintf("?from=%s&to=%s", fromStop, toStop),
	)
	timetable, err := gotrans.FetchTimetable(c, fromStop, toStop, date)
	if err != nil {
		return err
	}
	timetable, err = gotrans.TransformTimetableForClient(timetable)
	if err != nil {
		return err
	}
	page.Timetable = timetable
	page.To.Code = toStop
	page.From.Code = fromStop

	document, err := renderTemplates(c, []string{"datePicker", "otherway", "timetable"}, page)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, document)
}

// handleTo responds with the list of destinations and clears the timetable.
func handleTo(c echo.Context) error {
	fromStop := c.QueryParam("from")
	date := defaultIfEmpty(c.QueryParam("date"), defaultDate())
	page := NewPage(defaultDate(), date)
	c.Response().Header().Add(
		"HX-Push-Url",
		fmt.Sprintf("?from=%s", fromStop),
	)
	dests, err := gotrans.FetchDestinations(c, fromStop, date)
	if err != nil {
		return err
	}
	dests = append(gotrans.Destinations{gotrans.Destination{Code: "", Name: "Where to?"}}, dests...)
	page.DestinationsTo = dests.SetSelected(dests[0].Code)

	document, err := renderTemplates(c, []string{"selectTo", "otherway", "timetable"}, page)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, document)
}

// handleRoot responds with the full document for the default options.
func handleRoot(c echo.Context) error {
	fromStop := defaultIfEmpty(c.QueryParam("from"), defaultFrom().Code)
	toStop := defaultIfEmpty(c.QueryParam("to"), defaultTo().Code)
	date := defaultIfEmpty(c.QueryParam("date"), defaultDate())

	page := NewPage(defaultDate(), date)
	c.Echo().Logger.Debugf("datepicker: %v", page.DatePicker)
	// Fetch destination list for FROM and TO drop downs

	// Always fetch desinations from Union for our default list of stations
	destinationsDefault, err := gotrans.FetchDestinationsDefault(c, date)
	if err != nil {
		return err
	}
	page.DestinationsFrom = destinationsDefault

	destinations, err := gotrans.FetchDestinations(c, fromStop, date)
	if err != nil {
		return err
	}
	page.DestinationsTo = destinations

	// Set "selected" for the drop down

	var from gotrans.Destination
	fromIdx := destinationsDefault.IndexOfCode(fromStop)
	if fromIdx == -1 {
		from = defaultFrom()
	} else {
		from = destinationsDefault[fromIdx]
	}
	page.From = from
	page.DestinationsFrom = page.DestinationsFrom.SetSelected(from.Code)

	// TODO: evaluate if this is really necessary for TO
	var to gotrans.Destination
	toIdx := page.DestinationsTo.IndexOfCode(toStop)
	if toIdx == -1 {
		to = defaultTo()
	} else {
		to = page.DestinationsTo[toIdx]
	}
	page.To = to
	page.DestinationsTo = page.DestinationsTo.SetSelected(to.Code)

	// Fetch upcoming departures for the departure stop

	departures, err := gotrans.FetchDepartures(c, fromStop)
	if err != nil {
		return err
	}

	// Fetch timetable for the from-to combination

	timetable, err := gotrans.FetchTimetable(c, fromStop, toStop, date)
	if err != nil {
		return err
	}
	timetable, err = gotrans.TransformTimetableForClient(timetable)
	if err != nil {
		return err
	}
	timetable.AddPlatforms(departures.ToTripNumberXPlatform())

	page.Timetable = timetable

	return c.Render(http.StatusOK, "index", page)
}
