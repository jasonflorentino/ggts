package gotrans

import "sort"

type Destination struct {
	Code         string      `json:"code"`
	Name         string      `json:"name"`
	TransitType  TransitType `json:"transitType"`
	X_isSelected bool
}

type Destinations []Destination

func (dests Destinations) IndexOfCode(code string) int {
	for i, d := range dests {
		if d.Code == code {
			return i
		}
	}
	return -1
}

type destsSorter struct {
	dests Destinations
	by    func(a, b *Destination) bool
}

// Len is part of sort.Interface.
func (s *destsSorter) Len() int {
	return len(s.dests)
}

// Swap is part of sort.Interface.
func (s *destsSorter) Swap(i, j int) {
	s.dests[i], s.dests[j] = s.dests[j], s.dests[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *destsSorter) Less(i, j int) bool {
	return s.by(&s.dests[i], &s.dests[j])
}

func (dests Destinations) Sort() {
	sorter := &destsSorter{
		dests: dests,
		by: func(a, b *Destination) bool {
			return a.Name < b.Name
		},
	}
	sort.Sort(sorter)
}

func (dests Destinations) SetSelected(code string) Destinations {
	newDests := make(Destinations, len(dests))
	var dest Destination
	for i, d := range dests {
		if d.Code == code {
			dest.Code = code
			dest.Name = d.Name
			dest.X_isSelected = true
			newDests[i] = dest
		} else {
			newDests[i] = d
		}
	}
	return newDests
}

func (dests Destinations) OnlyRail() Destinations {
	i := 0
	for _, d := range dests {
		if d.TransitType == 1 || d.TransitType == 2 {
			dests[i] = d
			i++
		}
	}
	return dests[:i]

}

type Line struct {
	BlockNumber      string        `json:"blockNumber"`
	FromNotesImages  []interface{} `json:"fromNotesImages"`
	FromStopCode     string        `json:"fromStopCode"`
	FromStopTime     string        `json:"fromStopTime"`
	FromStopDisplay  string        `json:"fromStopDisplay"`
	HeadSign         string        `json:"headSign"`
	IsExpress        bool          `json:"isExpress"`
	IsTransfer       bool          `json:"isTransfer"`
	LineDisplay      string        `json:"lineDisplay"`
	ServiceLineName  string        `json:"serviceLineName"`
	Stops            []Stop        `json:"stops"`
	ToNotesImages    []interface{} `json:"toNotesImages"`
	ToStopCode       string        `json:"toStopCode"`
	ToStopDisplay    string        `json:"toStopDisplay"`
	ToStopTime       string        `json:"toStopTime"`
	TransferStopCode string        `json:"transferStopCode"`
	TransferDisplay  string        `json:"transferDisplay"`
	TransitType      TransitType   `json:"transitType"`
	TransitTypeName  string        `json:"transitTypeName"`
	TripNumber       string        `json:"tripNumber"`
}

type Note struct {
	English string `json:"english"`
	French  string `json:"french"`
	Image   string `json:"image"`
	Type    string `json:"type"`
}

type Stop struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Time        string      `json:"time"` // "14:44"
	TransitType TransitType `json:"transitType"`
}

type Timetable struct {
	ArrivalDisplay       string        `json:"arrivalDisplay"`
	ArrivalNotesImages   []interface{} `json:"arrivalNotesImages"`
	ArrivalStopId        string        `json:"arrivalStopId"`
	Date                 string        `json:"date"`
	DepartureDisplay     string        `json:"departureDisplay"`
	DepartureNotesImages []interface{} `json:"departureNotesImages"`
	DepartureStopId      string        `json:"departureStopId"`
	Notes                []Note        `json:"notes"`
	ServiceCode          string        `json:"serviceCode"`
	ServiceName          string        `json:"serviceName"`
	Trips                []Trip        `json:"trips"`
}

type TransitType int

type Trip struct {
	ArrivalTimeDisplay   string      `json:"arrivalTimeDisplay"`
	DepartureTimeDisplay string      `json:"departureTimeDisplay"`
	Duration             string      `json:"duration"`
	DurationMinutes      int         `json:"durationMinutes"`
	Lines                []Line      `json:"lines"`
	OrderTime            string      `json:"orderTime"`
	ServiceCode          string      `json:"serviceCode"`
	ServiceName          string      `json:"serviceName"` // Line Name
	Transfers            int         `json:"transfers"`
	TransitType          TransitType `json:"transitType"`
}

type Trips = []Trip
