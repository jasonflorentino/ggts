package main

import (
	"bytes"
	"fmt"
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

func NewPage() Page {
	return Page{
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
	page := NewPage()
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
	e.GET("/trips", handleTrips)
	e.GET("/to", handleTo)
	e.GET("/", handleRoot)

	// Run
	e.Logger.Fatal(e.Start(":" + env.Port()))
}

func handleTrips(c echo.Context) error {
	page := NewPage()
	fromStop := c.QueryParam("fromStop")
	toStop := c.QueryParam("toStop")
	date := defaultIfEmpty(c.QueryParam("date"), defaultDate())
	c.Response().Header().Add(
		"HX-Push-Url",
		fmt.Sprintf("?fromStop=%s&toStop=%s", fromStop, toStop),
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

	document, err := renderTemplates(c, []string{"timetable", "otherway"}, page)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, document)
}

func handleTo(c echo.Context) error {
	page := NewPage()
	fromStop := c.QueryParam("fromStop")
	date := defaultIfEmpty(c.QueryParam("date"), defaultDate())
	c.Response().Header().Add(
		"HX-Push-Url",
		fmt.Sprintf("?fromStop=%s", fromStop),
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

func handleRoot(c echo.Context) error {
	page := NewPage()

	fromStop := defaultIfEmpty(c.QueryParam("fromStop"), defaultFrom().Code)
	toStop := defaultIfEmpty(c.QueryParam("toStop"), defaultTo().Code)
	date := defaultIfEmpty(c.QueryParam("date"), defaultDate())

	// Fetch destination list for FROM and TO drop downs

	// Always fetch desinations from Union for our default list of stations
	destinationsDefault, err := gotrans.FetchDestinationsDefault(c, date)
	if err != nil {
		return err
	}
	page.DestinationsFrom = destinationsDefault

	if fromStop == gotrans.StationCode.Union {
		page.DestinationsTo = destinationsDefault
	} else {
		destinations, err := gotrans.FetchDestinations(c, fromStop, date)
		if err != nil {
			return err
		}
		page.DestinationsTo = destinations
	}

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

	// Fetch timetable for the from-to combination

	timetable, err := gotrans.FetchTimetable(c, fromStop, toStop, date)
	if err != nil {
		return err
	}
	timetable, err = gotrans.TransformTimetableForClient(timetable)
	if err != nil {
		return err
	}
	page.Timetable = timetable

	return c.Render(http.StatusOK, "index", page)
}
