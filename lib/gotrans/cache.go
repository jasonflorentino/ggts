package gotrans

import lru "github.com/hashicorp/golang-lru/v2"

var Cache Caches

type Caches struct {
	Destinations *lru.Cache[string, Destinations]
	Timetable    *lru.Cache[string, Timetable]
}

func InitCache() {
	Cache = Caches{
		Destinations: makeDestinationsCache(),
		Timetable:    makeTimetableCache(),
	}
}
