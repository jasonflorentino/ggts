package gotrans

import (
	"encoding/json"
	"fmt"
	"ggts/lib/log"
	"net/http"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/labstack/echo/v4"
)

func initDeparturesCache() *expirable.LRU[string, Departures] {
	const MAX_ITEMS = 10
	cache := expirable.NewLRU[string, Departures](MAX_ITEMS, nil, time.Second*30)
	return cache
}

func toDeparturesKey(destinationCode string) string {
	return fmt.Sprintf("%s", destinationCode)
}

func FetchDepartures(c echo.Context, destinationCode string) (Departures, error) {
	cacheKey := toDeparturesKey(destinationCode)
	if Cache.Departures.Contains(cacheKey) {
		log.To(c).Infof("Departures Cache HIT: %s", cacheKey)
		cachedDepartures, _ := Cache.Departures.Get(cacheKey)
		return cachedDepartures, nil
	}
	log.To(c).Infof("Departures Cache MISS: %s", cacheKey)

	req, err := Request(c, fmt.Sprintf("/external/go/departures/stops/%s/status/departures?page=1&transitTypeName=All&pageLimit=5", destinationCode))
	if err != nil {
		return Departures{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error creating http request: %s\n", err),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Departures{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error sending http request: %s\n", err),
		)
	}
	log.To(c).Infof("Got response - Status: %d, ContentLength: %d", res.StatusCode, res.ContentLength)

	body, err := GetBody(res)
	if err != nil {
		return Departures{}, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.To(c).Debugf("Body: %s", string(body))

	var departures Departures
	if err := json.Unmarshal(body, &departures); err != nil {
		return Departures{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not unmarshal json: %s\n", err),
		)
	}
	Cache.Departures.Add(cacheKey, departures)
	return departures, nil
}
