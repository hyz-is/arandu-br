package br

import "errors"

// ErrInvalidVoterTitle is what ParseVoterTitle returns for a value that is not
// a título de eleitor.
var ErrInvalidVoterTitle = errors.New("br: not a valid voter registration number")

// voterTitleMask is a título de eleitor of twelve digits punctuated: the
// eight of the sequence and the two of the state and of the check, in groups
// of four.
const voterTitleMask = "9999 9999 9999"

// ParseVoterTitle returns the digits of a título de eleitor, written bare or,
// for the common twelve digits, punctuated (0000 0000 0000), and
// ErrInvalidVoterTitle for anything else.
//
// The number is a sequence, then two digits for the state that issued it -- 01
// to 28, São Paulo to the title issued abroad -- and two check digits. It reads
// from five to thirteen digits, because the sequence of the older titles is
// shorter and a few are longer. The first check digit weighs the sequence from
// 9 down to 2 and over again from the right; the second weighs the state and
// the first check digit the same way; each is the sum modulo eleven, with ten
// written as zero and, for São Paulo and Minas Gerais, a zero written as one.
// A number of one digit repeated is refused. It does not say whether the
// voter is registered: only the Justiça Eleitoral can.
func ParseVoterTitle(s string) (string, error) {
	digits, ok := read(s, voterTitleMask)
	if !ok {
		digits, ok = readDigits(s, 5, 13)
	}
	if !ok || repeated(digits) {
		return "", ErrInvalidVoterTitle
	}
	n := len(digits)
	state := digits[n-4 : n-2]
	number := int(state[0]-'0')*10 + int(state[1]-'0')
	if number < 1 || number > 28 {
		return "", ErrInvalidVoterTitle
	}
	first := voterCheckDigit(digits[:n-4], number)
	second := voterCheckDigit(state+string(first), number)
	if digits[n-2] != first || digits[n-1] != second {
		return "", ErrInvalidVoterTitle
	}
	return digits, nil
}

// voterCheckDigit is the check digit of a sequence of digits: weights 9 down
// to 2 from the right, and over again, modulo eleven.
func voterCheckDigit(sequence string, state int) byte {
	weight, sum := 9, 0
	for i := len(sequence) - 1; i >= 0; i-- {
		sum += int(sequence[i]-'0') * weight
		if weight == 2 {
			weight = 10
		}
		weight--
	}
	switch rest := sum % 11; {
	case rest == 0 && state < 3:
		return '1'
	case rest == 10:
		return '0'
	default:
		return byte('0' + rest)
	}
}

// ValidVoterTitle reports whether ParseVoterTitle accepts s.
func ValidVoterTitle(s string) bool {
	_, err := ParseVoterTitle(s)
	return err == nil
}

// readDigits returns s when it is from min to max ASCII digits and nothing
// else.
func readDigits(s string, min, max int) (string, bool) {
	if len(s) < min || len(s) > max {
		return "", false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return "", false
		}
	}
	return s, true
}
