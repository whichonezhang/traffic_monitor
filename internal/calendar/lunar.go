package calendar

import (
	"fmt"
	"time"
)

// LunarCalendar handles lunar calendar related operations
type LunarCalendar struct {
	festivalMap map[string]time.Time
}

// NewLunarCalendar creates a new lunar calendar instance
func NewLunarCalendar() *LunarCalendar {
	return &LunarCalendar{
		festivalMap: map[string]time.Time{
			"春节":  time.Date(2024, 2, 10, 0, 0, 0, 0, time.Local),
			"元宵节": time.Date(2024, 2, 24, 0, 0, 0, 0, time.Local),
			"端午节": time.Date(2024, 6, 10, 0, 0, 0, 0, time.Local),
			"中秋节": time.Date(2024, 9, 17, 0, 0, 0, 0, time.Local),
		},
	}
}

// GetFestival returns the festival name if the given date is a lunar festival
func (c *LunarCalendar) GetFestival(date time.Time) (string, bool) {
	for festival, festivalDate := range c.festivalMap {
		if date.Year() == festivalDate.Year() &&
			date.Month() == festivalDate.Month() &&
			date.Day() == festivalDate.Day() {
			return festival, true
		}
	}
	return "", false
}

// GetPreviousFestivalDate returns the date of the same festival from the previous year
func (c *LunarCalendar) GetPreviousFestivalDate(currentDate time.Time, festival string) (time.Time, error) {
	festivalDate, exists := c.festivalMap[festival]
	if !exists {
		return time.Time{}, fmt.Errorf("festival %s not found", festival)
	}

	// Simply subtract one year from the festival date
	return festivalDate.AddDate(-1, 0, 0), nil
}

// GetNextFestival returns the next upcoming festival and its date
func (c *LunarCalendar) GetNextFestival(currentDate time.Time) (string, time.Time, error) {
	var nextFestival string
	var nextDate time.Time
	minDiff := time.Duration(1<<63 - 1)

	for festival, festivalDate := range c.festivalMap {
		// If the festival date is in the past, get next year's date
		if festivalDate.Before(currentDate) {
			festivalDate = festivalDate.AddDate(1, 0, 0)
		}

		diff := festivalDate.Sub(currentDate)
		if diff < minDiff {
			minDiff = diff
			nextFestival = festival
			nextDate = festivalDate
		}
	}

	if nextFestival == "" {
		return "", time.Time{}, fmt.Errorf("no upcoming festivals found")
	}

	return nextFestival, nextDate, nil
}
