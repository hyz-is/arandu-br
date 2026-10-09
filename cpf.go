package br

import "errors"

// ErrInvalidCPF is what ParseCPF and CPFCheckDigits return for a value that is
// not a CPF. It says nothing about which part was wrong, because the person who
// typed the value is the only one who can tell a typo from a wrong number.
var ErrInvalidCPF = errors.New("br: not a valid CPF")

// ParseCPF returns the eleven digits of a CPF, written bare (00000000000) or
// punctuated (000.000.000-00), and ErrInvalidCPF for anything else.
//
// It checks the shape and the two check digits the Receita Federal computes,
// and refuses the eleven repetitions of one digit, whose check digits compute
// but which belong to nobody. It does not say whether the CPF was issued or is
// regular: only the Receita Federal can.
//
// The digits come back bare, which is the spelling to store: a column that
// holds both spellings cannot be searched for either.
func ParseCPF(s string) (string, error) {
	digits, ok := read(s, cpfMask)
	if !ok || repeated(digits) {
		return "", ErrInvalidCPF
	}
	check, err := CPFCheckDigits(digits[:9])
	if err != nil || digits[9:] != check {
		return "", ErrInvalidCPF
	}
	return digits, nil
}

// ValidCPF reports whether ParseCPF accepts s.
func ValidCPF(s string) bool {
	_, err := ParseCPF(s)
	return err == nil
}

// CPFCheckDigits returns the two check digits of a CPF whose first nine digits
// are base, and ErrInvalidCPF when base is not nine digits.
//
// Each check digit is the sum of the digits before it, weighted from 10 (or 11)
// down to 2, times ten, modulo eleven -- with ten written as zero. It is what a
// seeder or a factory needs to write a CPF that passes, without borrowing one
// that belongs to somebody.
func CPFCheckDigits(base string) (string, error) {
	if len(base) != 9 {
		return "", ErrInvalidCPF
	}
	for i := 0; i < len(base); i++ {
		if base[i] < '0' || base[i] > '9' {
			return "", ErrInvalidCPF
		}
	}
	first := cpfDigit(base)
	second := cpfDigit(string(append([]byte(base), first)))
	return string([]byte{first, second}), nil
}

// cpfDigit is the check digit that follows digits: weights from len+1 down to
// 2, the sum times ten modulo eleven, and ten written as zero.
func cpfDigit(digits string) byte {
	sum := 0
	for i := 0; i < len(digits); i++ {
		sum += int(digits[i]-'0') * (len(digits) + 1 - i)
	}
	rest := sum * 10 % 11
	if rest == 10 {
		rest = 0
	}
	return byte('0' + rest)
}

// CPF writes a CPF punctuated, as 000.000.000-00.
//
// It takes the eleven digits bare or already punctuated, and returns anything
// else unchanged: it formats, and does not validate, so a page that prints a
// stored value shows what is stored rather than an empty cell. Validate with
// ParseCPF where the value is taken in.
func CPF(s string) string { return format(s, cpfMask) }
