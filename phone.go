package br

import (
	"errors"
	"strings"
)

// ErrInvalidLandline is what the Landline parsers return for a value that is
// not a Brazilian landline number.
var ErrInvalidLandline = errors.New("br: not a valid landline number")

// ErrInvalidAreaCode is what ParseAreaCode returns for a value that is not an
// area code (DDD) Anatel assigns.
var ErrInvalidAreaCode = errors.New("br: not a valid area code")

// countryCode is the calling code of Brazil, which the "with country code"
// spellings of a number begin with.
const countryCode = "55"

// numberKind tells a mobile number from a landline one. They share the area
// codes and the punctuation and differ in the subscriber number: nine digits
// that start with 9, or eight digits that do not start with 0 or 1.
type numberKind int

const (
	mobileNumber numberKind = iota
	landlineNumber
)

// subscriberMask is the subscriber number punctuated.
func (k numberKind) subscriberMask() string {
	if k == mobileNumber {
		return "99999-9999"
	}
	return "9999-9999"
}

// subscriber reports whether n is a subscriber number of this kind: the nine
// digits of a mobile line, starting with 9 and not the nine followed by eight
// zeros a form fills when nobody types anything, or the eight of a landline,
// whose first digit is not 0 or 1 -- those start the emergency and the service
// numbers.
func (k numberKind) subscriber(n string) bool {
	if k == mobileNumber {
		return len(n) == 9 && n[0] == '9' && n[1:] != "00000000"
	}
	return len(n) == 8 && n[0] >= '2'
}

// readLocal reads a subscriber number written without an area code, bare or
// punctuated, and returns its digits.
func readLocal(s string, kind numberKind) (string, bool) {
	digits, ok := read(s, kind.subscriberMask())
	if !ok || !kind.subscriber(digits) {
		return "", false
	}
	return digits, true
}

// readArea reads a number with its area code -- bare, or punctuated as
// (00) 0000-0000 or (00)0000-0000 -- and returns the digits.
func readArea(s string, kind numberKind) (string, bool) {
	mask := "(99) " + kind.subscriberMask()
	// The area code closed with no space before the number is the same
	// punctuation, one character shorter: read it with the space put back.
	if len(s) == len(mask)-1 && s[0] == '(' && s[3] == ')' && s[4] != ' ' {
		s = s[:4] + " " + s[4:]
	}
	digits, ok := read(s, mask)
	if !ok || !ValidAreaCode(digits[:2]) || !kind.subscriber(digits[2:]) {
		return "", false
	}
	return digits, true
}

// readCountry reads a number that begins with the country code, +55, and
// returns the digits with it: 55, the area code and the subscriber number.
//
// Both spellings are read. The punctuated one, +55 (00) 00000-0000, is the
// area code in parentheses and the subscriber number split by its hyphen, with
// one optional space after the country code and one after the parenthesis.
// The unpunctuated one, +5500000000000, is the digits alone, with one optional
// space after the country code and one after the area code. The second reports
// bare.
func readCountry(s string, kind numberKind) (digits string, bare, ok bool) {
	rest, found := strings.CutPrefix(s, "+"+countryCode)
	if !found {
		return "", false, false
	}
	rest = strings.TrimPrefix(rest, " ")

	if strings.HasPrefix(rest, "(") {
		if len(rest) < 4 || rest[3] != ')' {
			return "", false, false
		}
		area, aok := read(rest[1:3], "99")
		number := strings.TrimPrefix(rest[4:], " ")
		mask := kind.subscriberMask()
		if !aok || len(number) != len(mask) {
			return "", false, false
		}
		sub, sok := read(number, mask)
		if !sok || !ValidAreaCode(area) || !kind.subscriber(sub) {
			return "", false, false
		}
		return countryCode + area + sub, false, true
	}

	if len(rest) < 2 {
		return "", false, false
	}
	area, aok := read(rest[:2], "99")
	number := strings.TrimPrefix(rest[2:], " ")
	sub, sok := read(number, strings.Repeat("9", len(kind.subscriberMask())-1))
	if !aok || !sok || !ValidAreaCode(area) || !kind.subscriber(sub) {
		return "", false, false
	}
	return countryCode + area + sub, true, true
}

// ParseAreaCode returns the two digits of an area code (DDD) written bare, and
// ErrInvalidAreaCode for anything else, including a code Anatel does not
// assign. It accepts nothing around the digits: not a leading zero, not
// parentheses.
func ParseAreaCode(s string) (string, error) {
	if len(s) != 2 || !containsField(areaCodes, s) {
		return "", ErrInvalidAreaCode
	}
	return s, nil
}

// ValidAreaCode reports whether ParseAreaCode accepts s.
func ValidAreaCode(s string) bool {
	_, err := ParseAreaCode(s)
	return err == nil
}

// ParseLocalMobile returns the nine digits of a mobile number written without
// the area code, bare (900000000) or punctuated (90000-0000), and
// ErrInvalidMobile for anything else.
//
// The number starts with 9 -- every mobile number has since 2016, and an
// eight-digit number is a landline, so ParseLocalLandline reads it -- and it
// refuses the nine followed by eight zeros, which a form fills when somebody
// types nothing.
func ParseLocalMobile(s string) (string, error) {
	digits, ok := readLocal(s, mobileNumber)
	if !ok {
		return "", ErrInvalidMobile
	}
	return digits, nil
}

// ValidLocalMobile reports whether ParseLocalMobile accepts s.
func ValidLocalMobile(s string) bool {
	_, err := ParseLocalMobile(s)
	return err == nil
}

// ParseMobileWithCountry returns the thirteen digits of a Brazilian mobile
// number written with its country code -- 55, the area code and the nine
// digits -- and ErrInvalidMobile for anything else.
//
// It reads +55 (00) 00000-0000, with the optional spaces after the country
// code and after the parenthesis, and +5500000000000, with the optional
// spaces after the country code and after the area code. The code has to be
// Brazil's: the area code and the number are checked against Brazil's plan,
// which no other country's number follows. The area code is one Anatel
// assigns, as in ParseMobile.
func ParseMobileWithCountry(s string) (string, error) {
	digits, _, ok := readCountry(s, mobileNumber)
	if !ok {
		return "", ErrInvalidMobile
	}
	return digits, nil
}

// ValidMobileWithCountry reports whether ParseMobileWithCountry accepts s.
func ValidMobileWithCountry(s string) bool {
	_, err := ParseMobileWithCountry(s)
	return err == nil
}

// validUnpunctuatedMobileWithCountry reports whether s is a mobile number with
// its country code written with no parentheses and no hyphen: +5500000000000.
func validUnpunctuatedMobileWithCountry(s string) bool {
	_, bare, ok := readCountry(s, mobileNumber)
	return ok && bare
}

// ParseLocalLandline returns the eight digits of a landline number written
// without the area code, bare (30000000) or punctuated (3000-0000), and
// ErrInvalidLandline for anything else. The first digit is not 0 or 1.
func ParseLocalLandline(s string) (string, error) {
	digits, ok := readLocal(s, landlineNumber)
	if !ok {
		return "", ErrInvalidLandline
	}
	return digits, nil
}

// ValidLocalLandline reports whether ParseLocalLandline accepts s.
func ValidLocalLandline(s string) bool {
	_, err := ParseLocalLandline(s)
	return err == nil
}

// ParseLandline returns the ten digits of a landline number with its area
// code, written bare (1130000000) or punctuated ((11) 3000-0000, or
// (11)3000-0000 without the space after the area code), and
// ErrInvalidLandline for anything else.
//
// The area code is one Anatel assigns. It checks the shape and the first digit
// of the number, and does not say whether the line exists: only the operator
// can.
func ParseLandline(s string) (string, error) {
	digits, ok := readArea(s, landlineNumber)
	if !ok {
		return "", ErrInvalidLandline
	}
	return digits, nil
}

// ValidLandline reports whether ParseLandline accepts s.
func ValidLandline(s string) bool {
	_, err := ParseLandline(s)
	return err == nil
}

// Landline writes a landline number with its area code punctuated, as
// (00) 0000-0000.
//
// It takes the ten digits bare or already punctuated, and returns anything
// else unchanged: it formats, and does not validate. Validate with
// ParseLandline where the value is taken in.
func Landline(s string) string { return format(s, "(99) 9999-9999") }

// ParseLandlineWithCountry returns the twelve digits of a Brazilian landline
// number written with its country code -- 55, the area code and the eight
// digits -- and ErrInvalidLandline for anything else. It reads the spellings
// ParseMobileWithCountry reads.
func ParseLandlineWithCountry(s string) (string, error) {
	digits, _, ok := readCountry(s, landlineNumber)
	if !ok {
		return "", ErrInvalidLandline
	}
	return digits, nil
}

// ValidLandlineWithCountry reports whether ParseLandlineWithCountry accepts s.
func ValidLandlineWithCountry(s string) bool {
	_, err := ParseLandlineWithCountry(s)
	return err == nil
}
