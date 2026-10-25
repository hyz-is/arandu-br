package br

import "errors"

// ErrInvalidCNS is what ParseCNS returns for a value that is not a CNS.
var ErrInvalidCNS = errors.New("br: not a valid CNS")

// cnsMask is a CNS punctuated: the fifteen digits in a group of three and
// three of four.
const cnsMask = "999 9999 9999 9999"

// ParseCNS returns the fifteen digits of a Cartão Nacional de Saúde, written
// bare (000000000000000) or punctuated (000 0000 0000 0000), and ErrInvalidCNS
// for anything else.
//
// A definitive card begins with 1 or 2 and carries 00 and then 0 or 1 in the
// positions that follow the NIS it was issued from; a provisional one begins
// with 7, 8 or 9. Either way the digits weighted from 15 down to 1 add up to a
// multiple of eleven. It does not say whether the card was issued: only the
// Ministério da Saúde can.
//
// The digits come back bare, which is the spelling to store.
func ParseCNS(s string) (string, error) {
	digits, ok := read(s, cnsMask)
	if !ok {
		return "", ErrInvalidCNS
	}
	switch digits[0] {
	case '1', '2':
		if digits[11] != '0' || digits[12] != '0' || digits[13] > '1' {
			return "", ErrInvalidCNS
		}
	case '7', '8', '9':
	default:
		return "", ErrInvalidCNS
	}
	sum := 0
	for i := 0; i < len(digits); i++ {
		sum += int(digits[i]-'0') * (15 - i)
	}
	if sum%11 != 0 {
		return "", ErrInvalidCNS
	}
	return digits, nil
}

// ValidCNS reports whether ParseCNS accepts s.
func ValidCNS(s string) bool {
	_, err := ParseCNS(s)
	return err == nil
}

// CNS writes a CNS punctuated, as 000 0000 0000 0000.
//
// It takes the fifteen digits bare or already punctuated, and returns anything
// else unchanged: it formats, and does not validate. Validate with ParseCNS
// where the value is taken in.
func CNS(s string) string { return format(s, cnsMask) }
