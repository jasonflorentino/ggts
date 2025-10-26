package main

import (
	"bytes"
	"fmt"
	"ggts/lib"
	"ggts/lib/env"
	"ggts/lib/gotrans"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

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
