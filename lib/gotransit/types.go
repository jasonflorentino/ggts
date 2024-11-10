package gotransit

type Destination struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	TransitType TransitType `json:"transitType"`
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

const (
	Bus  TransitType = 0
	Rail TransitType = 1
	All  TransitType = 2
)

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
