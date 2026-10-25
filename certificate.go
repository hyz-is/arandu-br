package br

import "errors"

// ErrInvalidCivilCertificate is what ParseCivilCertificate returns for a value
// that is not the number of a civil registry certificate.
var ErrInvalidCivilCertificate = errors.New("br: not a valid civil registry certificate number")

// certificateGroups are the lengths of the nine groups of the number: the
// registry, the service, the year of the registration, the book, the page and
// the term, and the check digits.
var certificateGroups = [...]int{6, 2, 2, 4, 1, 5, 3, 7, 2}

// readCertificate returns the thirty-two digits of a certificate number,
// written bare or in nine groups. A group ends with a dot or a space, and the
// last group follows a hyphen or a space: 000000.00.00.0000.0.00000.000.0000000-00
// and 000000 00 00 0000 0 00000 000 0000000 00 are both read, and so is a
// mix of the two.
func readCertificate(s string) (string, bool) {
	if len(s) == 32 {
		return readDigits(s, 32, 32)
	}
	out := make([]byte, 0, 32)
	for group, size := range certificateGroups {
		if group > 0 {
			if len(s) == 0 {
				return "", false
			}
			separator := s[0]
			last := group == len(certificateGroups)-1
			if separator != ' ' && !(last && separator == '-') && !(!last && separator == '.') {
				return "", false
			}
			s = s[1:]
		}
		if len(s) < size {
			return "", false
		}
		digits, ok := readDigits(s[:size], size, size)
		if !ok {
			return "", false
		}
		out = append(out, digits...)
		s = s[size:]
	}
	return string(out), len(s) == 0
}

// ParseCivilCertificate returns the thirty-two digits of the number of a birth,
// marriage or death certificate, written bare or in its nine groups, and
// ErrInvalidCivilCertificate for anything else.
//
// The groups are told apart by a dot or a space, with a hyphen or a space
// before the last: 000000.00.00.0000.0.00000.000.0000000-00 and
// 000000 00 00 0000 0 00000 000 0000000 00. The last two digits are check
// digits: each is the sum of the digits before it, weighted from 2 upwards for
// the first and from 1 upwards for the second, to 10 and then from 0, modulo
// eleven, with ten written as one. Thirty-two repetitions of one digit are
// refused: their check digits compute and the number belongs to nobody. It
// does not say whether the certificate exists: only the registry can.
//
// The digits come back bare, which is the spelling to store.
func ParseCivilCertificate(s string) (string, error) {
	digits, ok := readCertificate(s)
	if !ok || repeated(digits) {
		return "", ErrInvalidCivilCertificate
	}
	first := certificateCheckDigit(digits[:30])
	second := certificateCheckDigit(digits[:30] + string(first))
	if digits[30] != first || digits[31] != second {
		return "", ErrInvalidCivilCertificate
	}
	return digits, nil
}

// certificateCheckDigit is the check digit that follows digits.
func certificateCheckDigit(digits string) byte {
	weight, sum := 32-len(digits), 0
	for i := 0; i < len(digits); i++ {
		sum += int(digits[i]-'0') * weight
		weight++
		if weight > 10 {
			weight = 0
		}
	}
	if rest := sum % 11; rest <= 9 {
		return byte('0' + rest)
	}
	return '1'
}

// ValidCivilCertificate reports whether ParseCivilCertificate accepts s.
func ValidCivilCertificate(s string) bool {
	_, err := ParseCivilCertificate(s)
	return err == nil
}
