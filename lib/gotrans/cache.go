package gotrans

import "github.com/hashicorp/golang-lru/v2/expirable"

var Cache Caches

type Caches struct {
	Departures   *expirable.LRU[string, Departures]
	Destinations *expirable.LRU[string, Destinations]
	Timetable    *expirable.LRU[string, Timetable]
}

func InitCache() {
	Cache = Caches{
		Departures:   initDeparturesCache(),
		Destinations: initDestinationsCache(),
		Timetable:    initTimetableCache(),
	}
}
