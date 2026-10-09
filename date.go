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
