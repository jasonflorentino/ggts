package gotrans

import (
	"testing"
	"time"
)

func TestFilterTrips(t *testing.T) {
	trips := Trips{
		{OrderTime: "2025-09-03T09:00:00", TransitType: TransitTypes.Rail, Transfers: 0},
		{OrderTime: "2025-09-03T12:00:00", TransitType: TransitTypes.Rail, Transfers: 0},
		{OrderTime: "2025-09-03T15:00:00", TransitType: TransitTypes.Rail, Transfers: 0},
		{OrderTime: "2025-09-03T18:00:00", TransitType: TransitTypes.Bus, Transfers: 0},
		{OrderTime: "2025-09-03T19:00:00", TransitType: TransitTypes.Rail, Transfers: 1},
	}

	now := time.Date(2025, 9, 3, 13, 0, 0, 0, time.Local)

	got, err := FilterTrips(trips, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Trips{
		{OrderTime: "2025-09-03T12:00:00", TransitType: TransitTypes.Rail, Transfers: 0},
		{OrderTime: "2025-09-03T15:00:00", TransitType: TransitTypes.Rail, Transfers: 0},
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d trips, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i].OrderTime != want[i].OrderTime {
			t.Errorf("trip %d: expected %s, got %s", i, want[i].OrderTime, got[i].OrderTime)
		}
	}
}
