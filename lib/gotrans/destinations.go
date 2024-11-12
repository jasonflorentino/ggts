package gotrans

import (
	"encoding/json"
	"fmt"
	"gogotrainschedule/lib/log"
	"net/http"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/labstack/echo/v4"
)

func makeDestinationsCache() *lru.Cache[string, Destinations] {
	destinationCache, err := lru.New[string, Destinations](10)
	if err != nil {
		panic(fmt.Errorf("couldn't init destinations cache %s", err))
	}
	return destinationCache
}

func toDestinationsKey(destinationCode, date string) string {
	return fmt.Sprintf("%s:%s", destinationCode, date)
}

// date: "YYYY-MM-DD"
func FetchDestinations(c echo.Context, destinationCode, date string) (Destinations, error) {
	cacheKey := toDestinationsKey(destinationCode, date)
	if Cache.Destinations.Contains(cacheKey) {
		log.To(c).Infof("Destinations Cache HIT: %s", cacheKey)
		cachedDestinations, _ := Cache.Destinations.Get(cacheKey)
		return cachedDestinations, nil
	}
	log.To(c).Infof("Destinations Cache MISS: %s", cacheKey)

	req, err := Request(c, fmt.Sprintf("/v2/schedules/stops/%s/destinations?Date=%s", destinationCode, date))
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error creating http request: %s\n", err),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error sending http request: %s\n", err),
		)
	}
	log.To(c).Infof("Got response - Status: %d, ContentLength: %d", res.StatusCode, res.ContentLength)

	body, err := GetBody(res)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.To(c).Debugf("Body: %s", string(body))

	var destinations Destinations
	if err := json.Unmarshal(body, &destinations); err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not unmarshal json: %s\n", err),
		)
	}
	railDestinations := destinations.OnlyRail()
	Cache.Destinations.Add(cacheKey, railDestinations)
	return railDestinations, nil
}

// Fetches Union Station's destinations as the default list since it is
// central hub through which GO Trains connect.
// Union is a Rail station only so there will not be any bus destinations.
// This list won't include Union Station itself so we should add it to complete the list.
func FetchDestinationsDefault(c echo.Context, date string) (Destinations, error) {
	cacheKey := toDestinationsKey(StationCode.Union, date)
	if Cache.Destinations.Contains(cacheKey) {
		log.To(c).Infof("Destinations Cache HIT: %s", cacheKey)
		cachedDestinations, _ := Cache.Destinations.Get(cacheKey)
		return cachedDestinations, nil
	}
	log.To(c).Infof("Destinations Cache MISS: %s", cacheKey)

	destinations, err := FetchDestinations(c, StationCode.Union, date)
	if err != nil {
		return nil, err
	}
	unionIdx := destinations.IndexOfCode(Union.Code)
	if unionIdx == -1 {
		destinations = append(destinations, Union)
		destinations.Sort()
	}

	Cache.Destinations.Add(cacheKey, destinations)
	return destinations, nil
}
