package gotrans

import (
	"encoding/json"
	"fmt"
	"gogotrainschedule/lib/log"
	"net/http"
	"net/url"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/labstack/echo/v4"
)

func makeTimetableCache() *lru.Cache[string, Timetable] {
	const MAX_ITEMS = 10
	timetableCache, err := lru.New[string, Timetable](MAX_ITEMS)
	if err != nil {
		panic(fmt.Errorf("couldn't init timetable cache %s", err))
	}
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

// Filters only trips
// - that haven't happened yet
// - are rail
// - are direct
func FilterTrips(trips Trips) (Trips, error) {
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

func ToDurationDisplay(d string) string {
	parts := strings.Split(d, ":")
	parts[0], _ = strings.CutPrefix(parts[0], "0")
	return fmt.Sprintf("%sh%sm", parts[0], parts[1])
}
