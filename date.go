package br

import (
	"errors"
	"strconv"
	"time"
)

// Date writes the calendar day of t as a Brazilian page prints it: 31/12/2026,
// day, month and four-digit year.
//
// The day is the one t holds in its own location. An instant read from a
// database is usually in UTC, where a late evening in Brasília is already the
// next day, so convert it first -- br.Date(t.In(location)) -- with the location
// the reader lives in.
func Date(t time.Time) string { return t.Format("02/01/2006") }

// DateTime writes the calendar day and the time of day of t as a Brazilian
// page prints them: 31/12/2026 15:04, the day as Date writes it, a space, and
// the hours and minutes on a twenty-four hour clock.
//
// Like Date it reads t in its own location, so an instant from a database is
// converted first -- br.DateTime(t.In(location)) -- with the location the
// reader lives in, and a page that keeps UTC says so beside it.
func DateTime(t time.Time) string { return t.Format("02/01/2006 15:04") }

// ErrInvalidDate is what ParseDate returns for a value that is not a calendar
// day written as 31/12/2026.
var ErrInvalidDate = errors.New("br: not a valid date")

// ErrInvalidTimeOfDay is what ParseTimeOfDay returns for a value that is not
// a time of day.
var ErrInvalidTimeOfDay = errors.New("br: not a valid time of day")

// ParseDate reads a calendar day written as a Brazilian page prints it,
// 31/12/2026 -- two digits of the day, two of the month and four of the
// year -- and returns midnight of that day in UTC. It returns ErrInvalidDate
// for anything else, including a day the month does not have: 29/02/2027,
// 31/04/2026, and the year 0000.
//
// A date has no time zone, so what it returns is the day and nothing else;
// compare its Year, Month and Day, or put the day in a location with
// time.Date.
func ParseDate(s string) (time.Time, error) {
	digits, ok := read(s, "99/99/9999")
	if !ok || len(s) != 10 {
		return time.Time{}, ErrInvalidDate
	}
	day, _ := strconv.Atoi(digits[:2])
	month, _ := strconv.Atoi(digits[2:4])
	year, _ := strconv.Atoi(digits[4:])
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if year < 1 || t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, ErrInvalidDate
	}
	return t, nil
}

// ValidDate reports whether ParseDate accepts s.
func ValidDate(s string) bool {
	_, err := ParseDate(s)
	return err == nil
}

// ParseTimeOfDay reads a time of day written 15:04 or 15:04:05, on a
// twenty-four hour clock, and returns how long after midnight it is. It
// returns ErrInvalidTimeOfDay for anything else: hours above 23, minutes or
// seconds above 59, a missing zero, or a time with a meridiem.
func ParseTimeOfDay(s string) (time.Duration, error) {
	mask := "99:99:99"
	if len(s) == 5 {
		mask = "99:99"
	}
	digits, ok := read(s, mask)
	if !ok || len(s) != len(mask) {
		return 0, ErrInvalidTimeOfDay
	}
	hours, _ := strconv.Atoi(digits[:2])
	minutes, _ := strconv.Atoi(digits[2:4])
	seconds := 0
	if len(digits) == 6 {
		seconds, _ = strconv.Atoi(digits[4:])
	}
	if hours > 23 || minutes > 59 || seconds > 59 {
		return 0, ErrInvalidTimeOfDay
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}

// ValidTimeOfDay reports whether ParseTimeOfDay accepts s.
func ValidTimeOfDay(s string) bool {
	_, err := ParseTimeOfDay(s)
	return err == nil
}
