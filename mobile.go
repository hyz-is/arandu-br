package br

import "errors"

// ErrInvalidMobile is what ParseMobile returns for a value that is not a
// Brazilian mobile number.
var ErrInvalidMobile = errors.New("br: not a valid mobile number")

// mobileMask is a mobile number punctuated: the area code in parentheses, then
// the nine digits split five and four.
const mobileMask = "(99) 99999-9999"

// areaCodes are the area codes (DDD) Anatel assigns, one per space-separated
// field. A number whose first two digits are not one of them is no number.
const areaCodes = " 11 12 13 14 15 16 17 18 19 21 22 24 27 28 31 32 33 34 35 37 38 41 42 43 44 45 46 47 48 49 51 53 54 55 61 62 63 64 65 66 67 68 69 71 73 74 75 77 79 81 82 83 84 85 86 87 88 89 91 92 93 94 95 96 97 98 99 "

// ParseMobile returns the eleven digits of a Brazilian mobile number, written
// bare (00000000000) or punctuated ((00) 00000-0000, or (00)00000-0000 without
// the space after the area code, as it is often typed), and ErrInvalidMobile
// for anything else.
//
// It checks the shape, that the area code is one Anatel assigns, that the
// number starts with 9 (every mobile number since 2016), and refuses the nine
// followed by eight zeros, which a form fills when somebody types nothing. It
// does not say whether the line exists: only the operator can.
//
// The digits come back bare, which is the spelling to store.
func ParseMobile(s string) (string, error) {
	digits, ok := readArea(s, mobileNumber)
	if !ok {
		return "", ErrInvalidMobile
	}
	return digits, nil
}

// ValidMobile reports whether ParseMobile accepts s.
func ValidMobile(s string) bool {
	_, err := ParseMobile(s)
	return err == nil
}

// Mobile writes a mobile number punctuated, as (00) 00000-0000.
//
// It takes the eleven digits bare or already punctuated, and returns anything
// else unchanged: it formats, and does not validate. Validate with ParseMobile
// where the value is taken in.
func Mobile(s string) string { return format(s, mobileMask) }

// containsField reports whether field is one of the space-separated fields of
// list, which starts and ends with a space.
func containsField(list, field string) bool {
	needle := " " + field + " "
	for i := 0; i+len(needle) <= len(list); i++ {
		if list[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
