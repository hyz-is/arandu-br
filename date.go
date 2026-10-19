package br

import "time"

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
