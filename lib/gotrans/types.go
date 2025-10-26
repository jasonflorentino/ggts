package gotrans

import (
	"sort"
)

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

// Sort sorts dests by name
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

type Departures struct {
	StationCode     string            `json:"stationCode"`
	AllDepartures   TransitDepartures `json:"allDepartures,omitempty"`
	TrainDepartures TransitDepartures `json:"trainDepartures,omitempty"`
	BusDepartures   TransitDepartures `json:"busDepartures,omitempty"`
}

type TripNumber = string
type Platform = string
type PlatformMap = map[TripNumber]Platform

func (d Departures) ToPlatformMap() PlatformMap {
	deps := d.AllDepartures // I think with our use we'll always get this field from the API, but let's check it to be sure.
	if deps.IsEmpty() {
		return map[string]string{}
	}

	tripXPlatform := make(map[string]string)
	for _, item := range deps.Items {
		tripXPlatform[item.TripNumber] = item.Platform
	}

	return tripXPlatform
}

type TransitDepartures struct {
	Items          []Departure `json:"items"`
	Page           int         `json:"page"`
	PageSize       int         `json:"pageSize"`
	TotalItemCount int         `json:"totalItemCount"`
}

func (td TransitDepartures) IsEmpty() bool {
	return td.Items == nil && td.Page == 0 && td.PageSize == 0 && td.TotalItemCount == 0
}

type Departure struct {
	AllDepartureStops    AllDepartureStops `json:"allDepartureStops"`
	DelayedDepartureTime string            `json:"delayedDepartureTime,omitempty"`
	DelayMessage         string            `json:"delayMessage,omitempty"`
	DelaySeconds         int               `json:"delaySeconds,omitempty"`
	Gate                 *string           `json:"gate"` // null
	Info                 string            `json:"info"` // "Proceed / Avancez"
	LineCode             string            `json:"lineCode"`
	LineColour           string            `json:"lineColour"` // "#00853e"
	LineMessageEn        string            `json:"lineMessageEn,omitempty"`
	LineMessageFr        string            `json:"lineMessageFr,omitempty"`
	Platform             string            `json:"platform"`          // "5 & 6"
	ScheduledDateTime    string            `json:"scheduledDateTime"` // "2025-10-18T18:48:00"
	ScheduledPlatform    *string           `json:"scheduledPlatform"` // null
	ScheduledTime        string            `json:"scheduledTime"`     // "18:48"
	Service              string            `json:"service"`
	Status               string            `json:"status,omitempty"` // "ontime"
	StopsDisplay         string            `json:"stopsDisplay"`     // "Bloor-Weston-Malton"
	TransitType          int               `json:"transitType"`
	TransitTypeName      string            `json:"transitTypeName"`
	TripNumber           string            `json:"tripNumber"`
	Zone                 *string           `json:"zone"` // null
}

type AllDepartureStops struct {
	StayInTrain          bool              `json:"stayInTrain"` // false
	TripNumbers          []string          `json:"tripNumbers"` // ["1732"]
	DepartureDetailsList []DepartureDetail `json:"departureDetailsList"`
}

type DepartureDetail struct {
	StopName      string `json:"stopName"`      // "West Harbour GO"
	DepartureTime string `json:"departureTime"` // "19:24"
	StopCode      string `json:"stopCode"`      // "WR"
	IsMajorStop   bool   `json:"isMajorStop"`   // true
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
	X_DateDisplay        string
	X_DateOnly           string
}

func (t *Timetable) AddPlatforms(platforms PlatformMap) {
	for i := range t.Trips {
		trip := &t.Trips[i]
		line := trip.Lines[0] // We only support direct trips; there will only be one "line" in this slice
		trip.X_Platform = platforms[line.TripNumber]
	}
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
	X_Platform           string
}

type Trips []Trip

func (ts *Trips) Map(fn func(Trip) Trip) {
	for i, t := range *ts {
		(*ts)[i] = fn(t)
	}
}

func (ts *Trips) Sort() {
	sort.Slice(*ts, func(i int, j int) bool {
		return (*ts)[i].OrderTime < (*ts)[j].OrderTime
	})
}
