package lib

import (
	"fmt"
	"ggts/lib/env"
	"time"
)

type DatePartOption struct {
	IsSelected bool
	Meta       string
	Name       string
	Val        string
}

type Day = DatePartOption
type Month = DatePartOption
type Year = DatePartOption

type DatePicker struct {
	Days         []Day
	RenderedDate string
	Months       []Month
	Years        []Year
}

func parseDateOnly(s string) (time.Time, error) {
	t, err := time.ParseInLocation(time.DateOnly, s, env.Location())
	if err != nil {
		return t, err
	}
	return t, nil
}

func lastDayOfMonth(t time.Time) time.Time {
	nextMonth := t.AddDate(0, 1, 0)
	lastDay := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 0, 0, 0, 0, t.Location())
	return lastDay
}

func NewDatePicker(today string, selected string) DatePicker {
	todayTime, err := parseDateOnly(today)
	if err != nil {
		panic(fmt.Sprintf("bad today date: %v", today))
	}
	selectedTime, err := parseDateOnly(selected)
	if err != nil {
		panic(fmt.Sprintf("bad selected date: %v", selected))
	}

	datePicker := DatePicker{}
	datePicker.RenderedDate = selected

	years := []Year{}
	for y := todayTime.Year(); y <= todayTime.AddDate(1, 0, 0).Year(); y += 1 {
		year := fmt.Sprintf("%d", y)
		years = append(years, Year{Name: year, Val: year, IsSelected: y == selectedTime.Year()})
	}
	datePicker.Years = years

	months := []Month{}
	for m := 1; m <= 12; m += 1 {
		month := fmt.Sprintf("%02d", m)
		months = append(months, Month{Name: time.Month(m).String(), Val: month, IsSelected: time.Month(m) == selectedTime.Month()})
	}
	datePicker.Months = months

	days := []Day{}
	for d := 1; d <= lastDayOfMonth(selectedTime).Day(); d += 1 {
		name := fmt.Sprintf("%d", d)
		val := fmt.Sprintf("%02d", d)
		day := Day{Name: name, Val: val, IsSelected: d == selectedTime.Day()}
		dTime := time.Date(selectedTime.Year(), selectedTime.Month(), d, 0, 0, 0, 0, selectedTime.Location())
		// Mark Sundays
		if dTime.Weekday() == time.Weekday(0) {
			day.Meta = dTime.Format("(Monday)")
		}
		days = append(days, day)
	}
	datePicker.Days = days

	return datePicker
}
