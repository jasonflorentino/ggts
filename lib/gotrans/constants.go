package gotrans

const (
	bus  TransitType = 0
	rail TransitType = 1
	all  TransitType = 2
)

var TransitTypes = struct {
	Bus  TransitType
	Rail TransitType
	All  TransitType
}{
	Bus:  bus,
	Rail: rail,
	All:  all,
}

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
