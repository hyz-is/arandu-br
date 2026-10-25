package br

import "errors"

// ErrInvalidNIS is what ParseNIS returns for a value that is not a NIS.
var ErrInvalidNIS = errors.New("br: not a valid NIS")

// nisMask is a NIS punctuated.
const nisMask = "999.99999.99-9"

// nisWeights are the weights of the ten digits that precede the check digit.
var nisWeights = [10]int{3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

// ParseNIS returns the eleven digits of a NIS -- the number the PIS, the PASEP,
// the NIT and the Cadastro Único all use -- written bare (00000000000) or
// punctuated (000.00000.00-0), and ErrInvalidNIS for anything else.
//
// It checks the shape and the check digit, and refuses eleven repetitions of
// one digit. A number with fewer than eleven digits is refused, not padded: a
// leading zero is a digit of the number, and the person who left it out typed
// another number.
//
// The digits come back bare, which is the spelling to store.
func ParseNIS(s string) (string, error) {
	digits, ok := read(s, nisMask)
	if !ok || repeated(digits) {
		return "", ErrInvalidNIS
	}
	sum := 0
	for i, weight := range nisWeights {
		sum += int(digits[i]-'0') * weight
	}
	check := 0
	if rest := sum % 11; rest >= 2 {
		check = 11 - rest
	}
	if int(digits[10]-'0') != check {
		return "", ErrInvalidNIS
	}
	return digits, nil
}

// ValidNIS reports whether ParseNIS accepts s.
func ValidNIS(s string) bool {
	_, err := ParseNIS(s)
	return err == nil
}

// NIS writes a NIS punctuated, as 000.00000.00-0.
//
// It takes the eleven digits bare or already punctuated, and returns anything
// else unchanged: it formats, and does not validate. Validate with ParseNIS
// where the value is taken in.
func NIS(s string) string { return format(s, nisMask) }
