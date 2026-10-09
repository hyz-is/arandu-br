package br

import "errors"

// ErrInvalidCEP is what ParseCEP returns for a value that is not a CEP.
var ErrInvalidCEP = errors.New("br: not a valid CEP")

// lowestCEP is where the ranges the Correios allocate begin. Every one of them
// starts at 01000-000 or above, so a code below it -- 00000-000 among them,
// the placeholder a form is most often given -- addresses nowhere.
const lowestCEP = "01000000"

// ParseCEP returns the eight digits of a CEP, written bare (00000000) or
// punctuated (00000-000), and ErrInvalidCEP for anything else or for a code
// below 01000-000, where no range begins.
//
// A CEP carries no check digit, so a code that passes may still address
// nothing: only a lookup in the Correios base can say it exists.
//
// The digits come back bare, which is the spelling to store.
func ParseCEP(s string) (string, error) {
	digits, ok := read(s, cepMask)
	if !ok || digits < lowestCEP {
		return "", ErrInvalidCEP
	}
	return digits, nil
}

// ValidCEP reports whether ParseCEP accepts s.
func ValidCEP(s string) bool {
	_, err := ParseCEP(s)
	return err == nil
}

// CEP writes a CEP punctuated, as 00000-000.
//
// It takes the eight digits bare or already punctuated, and returns anything
// else unchanged: it formats, and does not validate. Validate with ParseCEP
// where the value is taken in.
func CEP(s string) string { return format(s, cepMask) }
