package gotrans

import "github.com/hashicorp/golang-lru/v2/expirable"

var Cache Caches

type Caches struct {
	Destinations *expirable.LRU[string, Destinations]
	Timetable    *expirable.LRU[string, Timetable]
}

func InitCache() {
	Cache = Caches{
		Destinations: initDestinationsCache(),
		Timetable:    initTimetableCache(),
	}
}
