package br

import (
	"errors"
	"strings"
)

// ErrInvalidRenavam is what ParseRenavam returns for a value that is not a
// RENAVAM.
var ErrInvalidRenavam = errors.New("br: not a valid RENAVAM")

// renavamWeights are the weights of the ten digits that precede the check
// digit.
var renavamWeights = [10]int{3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

// ParseRenavam returns the eleven digits of a Registro Nacional de Veículos
// Automotores number, written bare, and ErrInvalidRenavam for anything else.
//
// The number issued before 2013 has nine digits and the ones between that and
// the eleven have ten; both are read as the eleven they were always meant to
// be, with the zeros in front, and come back as those eleven. The last digit
// is the check digit: the ten before it weighted by 3, 2, 9, 8, 7, 6, 5, 4, 3
// and 2, the sum taken modulo eleven, eleven minus the remainder, and zero
// when that is ten or more. A number of one digit repeated is refused. It
// does not say whether the vehicle exists: only the Detran can.
func ParseRenavam(s string) (string, error) {
	var digits string
	switch len(s) {
	case 9, 10, 11:
		var ok bool
		digits, ok = read(strings.Repeat("0", 11-len(s))+s, "99999999999")
		if !ok {
			return "", ErrInvalidRenavam
		}
	default:
		return "", ErrInvalidRenavam
	}
	if repeated(digits) {
		return "", ErrInvalidRenavam
	}
	sum := 0
	for i, weight := range renavamWeights {
		sum += int(digits[i]-'0') * weight
	}
	check := 11 - sum%11
	if check >= 10 {
		check = 0
	}
	if int(digits[10]-'0') != check {
		return "", ErrInvalidRenavam
	}
	return digits, nil
}

// ValidRenavam reports whether ParseRenavam accepts s.
func ValidRenavam(s string) bool {
	_, err := ParseRenavam(s)
	return err == nil
}
