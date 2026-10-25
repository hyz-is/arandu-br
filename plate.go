package br

import "errors"

// ErrInvalidPlate is what ParsePlate returns for a value that is not a
// vehicle plate.
var ErrInvalidPlate = errors.New("br: not a valid vehicle plate")

// ErrInvalidPassport is what ParsePassport returns for a value that is not a
// Brazilian passport number.
var ErrInvalidPassport = errors.New("br: not a valid passport number")

// plateShapes are the shapes of the plates Brazil has issued, with L for a
// letter and N for a digit: the first, from 1915, a letter and five digits;
// then seven digits; two letters and four digits; three letters and four
// digits; and, since 2018, the Mercosul plate, three letters, a digit, a
// letter and two digits, and the variant for motorcycles, three letters, two
// digits, a letter and a digit.
var plateShapes = [...]string{"LNNNNN", "NNNNNNN", "LLNNNN", "LLLNNNN", "LLLNLNN", "LLLNNLN"}

// ParsePlate returns a vehicle plate in capitals and without punctuation, and
// ErrInvalidPlate for anything else.
//
// It reads every shape in plateShapes, a letter in either case, and the
// hyphen after the three letters that the plate of 1990 to 2018 is written
// with -- ABC-1234, and ABC-1D23 too -- which does not come back. It checks
// the shape and nothing else: only the Detran knows whether a plate was
// issued.
func ParsePlate(s string) (string, error) {
	if len(s) == 8 && s[3] == '-' {
		s = s[:3] + s[4:]
	}
	if len(s) != 6 && len(s) != 7 {
		return "", ErrInvalidPlate
	}
	out := make([]byte, len(s))
	shape := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c >= '0' && c <= '9':
			out[i], shape[i] = c, 'N'
		case c >= 'A' && c <= 'Z':
			out[i], shape[i] = c, 'L'
		case c >= 'a' && c <= 'z':
			out[i], shape[i] = c-'a'+'A', 'L'
		default:
			return "", ErrInvalidPlate
		}
	}
	for _, want := range plateShapes {
		if string(shape) == want {
			return string(out), nil
		}
	}
	return "", ErrInvalidPlate
}

// ValidPlate reports whether ParsePlate accepts s.
func ValidPlate(s string) bool {
	_, err := ParsePlate(s)
	return err == nil
}

// ParsePassport returns a Brazilian passport number in capitals, and
// ErrInvalidPassport for anything else. The number is two letters, in either
// case, and six digits. It checks the shape: the passport carries no check
// digit this package knows.
func ParsePassport(s string) (string, error) {
	if len(s) != 8 {
		return "", ErrInvalidPassport
	}
	out := make([]byte, 8)
	for i := 0; i < 8; i++ {
		c := s[i]
		switch {
		case i < 2 && c >= 'A' && c <= 'Z':
			out[i] = c
		case i < 2 && c >= 'a' && c <= 'z':
			out[i] = c - 'a' + 'A'
		case i >= 2 && c >= '0' && c <= '9':
			out[i] = c
		default:
			return "", ErrInvalidPassport
		}
	}
	return string(out), nil
}

// ValidPassport reports whether ParsePassport accepts s.
func ValidPassport(s string) bool {
	_, err := ParsePassport(s)
	return err == nil
}
