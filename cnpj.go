package br

import "errors"

// ErrInvalidCNPJ is what ParseCNPJ and CNPJCheckDigits return for a value that
// is not a CNPJ.
var ErrInvalidCNPJ = errors.New("br: not a valid CNPJ")

// ParseCNPJ returns the fourteen characters of a CNPJ, written bare
// (00000000000000) or punctuated (00.000.000/0000-00), and ErrInvalidCNPJ for
// anything else.
//
// The first twelve positions may hold letters as well as digits: that is the
// alphanumeric CNPJ the Receita Federal issues, and the numeric CNPJ is the case
// of it whose twelve positions are all digits, checked by the same computation.
// The last two positions are the check digits, and are always digits. A letter
// typed in lower case is read as the capital it stands for, and comes back as
// one.
//
// It refuses fourteen repetitions of one character, which no CNPJ is. It does
// not say whether the CNPJ was issued or is active: only the Receita Federal
// can.
//
// The characters come back bare, which is the spelling to store.
func ParseCNPJ(s string) (string, error) {
	chars, ok := read(s, cnpjMask)
	if !ok || repeated(chars) {
		return "", ErrInvalidCNPJ
	}
	check, err := CNPJCheckDigits(chars[:12])
	if err != nil || chars[12:] != check {
		return "", ErrInvalidCNPJ
	}
	return chars, nil
}

// ValidCNPJ reports whether ParseCNPJ accepts s.
func ValidCNPJ(s string) bool {
	_, err := ParseCNPJ(s)
	return err == nil
}

// CNPJCheckDigits returns the two check digits of a CNPJ whose first twelve
// characters are base, and ErrInvalidCNPJ when base is not twelve digits or
// letters. A lower-case letter is read as its capital.
//
// Each character is worth its ASCII code minus 48 -- a digit its own value, A
// seventeen, Z forty-two. A check digit is the sum of the characters before it,
// weighted from the right by 2 to 9 and again from 2, taken modulo eleven: a
// remainder below two gives zero, and any other gives eleven minus it.
func CNPJCheckDigits(base string) (string, error) {
	if len(base) != 12 {
		return "", ErrInvalidCNPJ
	}
	chars := make([]byte, 0, 13)
	for i := 0; i < len(base); i++ {
		c, ok := fill(base[i], 'A')
		if !ok {
			return "", ErrInvalidCNPJ
		}
		chars = append(chars, c)
	}
	first := cnpjDigit(chars)
	second := cnpjDigit(append(chars, first))
	return string([]byte{first, second}), nil
}

// cnpjDigit is the check digit that follows chars.
func cnpjDigit(chars []byte) byte {
	sum, weight := 0, 2
	for i := len(chars) - 1; i >= 0; i-- {
		sum += int(chars[i]-'0') * weight
		weight++
		if weight > 9 {
			weight = 2
		}
	}
	rest := sum % 11
	if rest < 2 {
		return '0'
	}
	return byte('0' + 11 - rest)
}

// CNPJ writes a CNPJ punctuated, as 00.000.000/0000-00, with its letters in
// capitals.
//
// It takes the fourteen characters bare or already punctuated, and returns
// anything else unchanged: it formats, and does not validate. Validate with
// ParseCNPJ where the value is taken in.
func CNPJ(s string) string { return format(s, cnpjMask) }
