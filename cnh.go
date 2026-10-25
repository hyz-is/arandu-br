package br

import "errors"

// ErrInvalidCNH is what ParseCNH returns for a value that is not a CNH.
var ErrInvalidCNH = errors.New("br: not a valid CNH")

// ParseCNH returns the eleven digits of a Carteira Nacional de Habilitação
// register number, written bare, and ErrInvalidCNH for anything else.
//
// The last two digits are the check digits the Denatran computes: the first
// is the sum of the nine digits before it, weighted from 9 down to 1, modulo
// eleven, with ten written as zero; the second is the sum of the same nine
// weighted from 1 up to 9, modulo eleven, less 2 when the first came from a
// ten, brought back above zero by adding eleven, with ten again written as
// zero. Eleven repetitions of one digit are refused. The number is not
// punctuated, and a leading zero is part of it.
//
// It does not say whether the licence was issued or is in force: only the
// Detran can.
func ParseCNH(s string) (string, error) {
	digits, ok := read(s, "99999999999")
	if !ok || repeated(digits) {
		return "", ErrInvalidCNH
	}
	first, second, discount := 0, 0, 0
	for i := 0; i < 9; i++ {
		first += int(digits[i]-'0') * (9 - i)
		second += int(digits[i]-'0') * (i + 1)
	}
	if first %= 11; first >= 10 {
		first, discount = 0, 2
	}
	second = (second%11 - discount + 11) % 11
	if second >= 10 {
		second = 0
	}
	if int(digits[9]-'0') != first || int(digits[10]-'0') != second {
		return "", ErrInvalidCNH
	}
	return digits, nil
}

// ValidCNH reports whether ParseCNH accepts s.
func ValidCNH(s string) bool {
	_, err := ParseCNH(s)
	return err == nil
}
