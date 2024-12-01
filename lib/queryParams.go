package lib

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func toDateOnly(year, month, day string) string {
	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

// GetChangeDate returns the requested change date specified in the query params.
// It will return successfully with the empty string if there is any part missing.
func GetChangeDate(c echo.Context) (string, error) {
	year := c.QueryParam("year")
	month := c.QueryParam("month")
	day := c.QueryParam("day")

	if year == "" || month == "" || day == "" {
		c.Echo().Logger.Warnf("GetChangeDate: missing param year:%s month:%s day:%s", year, month, day)
		return "", nil
	}

	selectedDate := toDateOnly(year, month, day)
	if len([]rune(selectedDate)) != len([]rune(time.DateOnly)) {
		return "", echo.NewHTTPError(http.StatusBadRequest, "selectedDate len != DateOnly")
	}

	if _, err := time.Parse(time.DateOnly, selectedDate); err != nil {
		c.Echo().Logger.Warnf("GetChangeDate: %v", err)
		selectedDate = toDateOnly(year, month, "01")
	}

	return selectedDate, nil
}
