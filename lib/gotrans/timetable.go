package gotrans

import (
	"encoding/json"
	"fmt"
	"ggts/lib/env"
	"ggts/lib/log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/labstack/echo/v4"
)

func initTimetableCache() *expirable.LRU[string, Timetable] {
	const MAX_ITEMS = 10
	timetableCache := expirable.NewLRU[string, Timetable](MAX_ITEMS, nil, time.Hour*1)
	return timetableCache
}

func toTimetableKey(fromStop, toStop, date string) string {
	return fmt.Sprintf("%s:%s:%s", fromStop, toStop, date)
}

// date: "YYYY-MM-DD"
func FetchTimetable(c echo.Context, fromStop, toStop, date string) (Timetable, error) {
	cacheKey := toTimetableKey(fromStop, toStop, date)
	if Cache.Timetable.Contains(cacheKey) {
		log.To(c).Infof("Timetable Cache HIT: %s", cacheKey)
		cached, _ := Cache.Timetable.Get(cacheKey)
		return cached, nil
	}
	log.To(c).Infof("Timetable Cache MISS: %s", cacheKey)

	params := url.Values{}
	params.Add("fromStop", fromStop)
	params.Add("toStop", toStop)
	params.Add("date", date)
	queryString := params.Encode()
	req, err := Request(c, fmt.Sprintf("/v2/schedules/en/timetable/all?%s", queryString))
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
	log.To(c).Infof("Got response - Status: %d, ContentLength: %d", res.StatusCode, res.ContentLength)

	body, err := GetBody(res)
	if err != nil {
		return Timetable{}, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.To(c).Debugf("Body: %s", string(body))

	var timetable Timetable
	if err := json.Unmarshal(body, &timetable); err != nil {
		return Timetable{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not unmarshal json: %s\n", err),
		)
	}

	Cache.Timetable.Add(cacheKey, timetable)
	return timetable, nil
}

func TransformTimetableForClient(timetable Timetable) (Timetable, error) {
	datestr := ToDatestring(timetable.Date)
	timetable.Date = datestr
	trips, err := FilterTrips(timetable.Trips)
	if err != nil {
		return timetable, err
	}
	trips.Map(func(t Trip) Trip {
		t.Duration = ToDurationDisplay(t.Duration)
		return t
	})
	timetable.Trips = trips
	return timetable, nil
}

// Filters only trips
// - that haven't happened yet
// - are rail
// - are direct
func FilterTrips(trips Trips) (Trips, error) {
	now := time.Now()
	out := make(Trips, 0)
	for _, trip := range trips {
		tripTime, err := time.ParseInLocation("2006-01-02T15:04:05", trip.OrderTime, time.Local)
		if err != nil {
			return nil, err
		}
		if tripTime.After(now) &&
			trip.TransitType == TransitTypes.Rail &&
			trip.Transfers == 0 {
			out = append(out, trip)
		}
	}
	return out, nil
}

func ToDurationDisplay(d string) string {
	parts := strings.Split(d, ":")
	if len(parts) != 3 {
		return d
	}
	// First char is always `0`.
	// There are no GO trips that take 10+ hours
	hour, _ := strings.CutPrefix(parts[0], "0")
	if hour == "0" {
		return fmt.Sprintf("%sm", parts[1])
	} else {
		return fmt.Sprintf("%sh%sm", hour, parts[1])
	}
}

func ToDatestring(s string) string {
	t, err := time.ParseInLocation("2006-01-02T15:04:05", s, env.Location())
	if err != nil {
		t, err = time.ParseInLocation("2006-01-02T15:04:05-07:00", s, env.Location())
	}
	if err != nil {
		return s
	}
	return t.Format("Mon Jan 2, 2006")
}
