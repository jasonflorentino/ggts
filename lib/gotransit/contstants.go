package gotransit

const API_URL = "api.gotransit.com"

var StationCode = struct {
	Union       string
	WestHarbour string
}{
	Union:       "UN",
	WestHarbour: "WR",
}

var Union Destination = Destination{
	Code:        StationCode.Union,
	Name:        "Union Station GO",
	TransitType: 1,
}

var WestHarbour Destination = Destination{
	Code:        StationCode.WestHarbour,
	Name:        "West Harbour GO",
	TransitType: 1,
}
