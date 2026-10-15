package br

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidBRL is what ParseBRL returns for a value that is not an amount in
// reais.
var ErrInvalidBRL = errors.New("br: not a valid amount in reais")

// BRL writes an amount of reais held as a count of centavos, the way a
// Brazilian page prints money: R$ 1.234,56.
//
// A dot separates each three digits of the reais, a comma stands before the
// two digits of the centavos, and the symbol leads with a space after it. A
// negative amount leads with the sign, before the symbol -- -R$ 1.234,56 --
// because a sign written after the symbol reads as part of the amount. Zero is
// R$ 0,00, and every int64 is written, its most negative value included.
//
// The amount is an integer because an amount that has to add up is held as
// one; no float64 appears on the way from it to the page. The result depends
// on nothing but the argument: hesape/number writes the same string from
// CurrencyFromCents(cents, "BRL") only while the process-wide locale is pt-BR,
// and a library cannot set that locale without setting it for everything else
// the application formats.
//
// BRL is FormatBRL with the zero BRLFormat.
func BRL(cents int64) string { return FormatBRL(cents, BRLFormat{}) }

// BRLFormat is how FormatBRL departs from BRL. Its zero value departs in
// nothing, so each field names what it leaves out.
type BRLFormat struct {
	// WithoutGrouping writes the reais with no dot between their thousands:
	// R$ 1234,56.
	WithoutGrouping bool
	// WithoutZeroCents leaves the centavos off an amount of whole reais, the
	// way a price list writes one: R$ 47, and still R$ 199,50.
	WithoutZeroCents bool
}

// FormatBRL writes an amount of reais held as a count of centavos, as BRL
// does, departing from it in what f says: R$ 1234,56 without the grouping,
// R$ 47 without the zero centavos. The sign, the symbol and the comma are
// BRL's, and every int64 is written.
func FormatBRL(cents int64, f BRLFormat) string {
	// The magnitude is taken in uint64 because the most negative int64 has no
	// positive counterpart to negate it into.
	negative := cents < 0
	magnitude := uint64(cents)
	if negative {
		magnitude = uint64(-(cents + 1)) + 1
	}
	reais := strconv.FormatUint(magnitude/100, 10)
	centavos := magnitude % 100

	var b strings.Builder
	if negative {
		b.WriteByte('-')
	}
	b.WriteString("R$ ")
	for i := 0; i < len(reais); i++ {
		if !f.WithoutGrouping && i > 0 && (len(reais)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteByte(reais[i])
	}
	if f.WithoutZeroCents && centavos == 0 {
		return b.String()
	}
	b.WriteByte(',')
	b.WriteByte(byte('0' + centavos/10))
	b.WriteByte(byte('0' + centavos%10))
	return b.String()
}

// ParseBRL reads an amount in reais as a person types it into a form, and
// returns it as a count of centavos, or ErrInvalidBRL.
//
// The amount may be led by the symbol, R$, with or without a space after it,
// and whitespace around the whole is ignored. The reais are digits, bare
// (1234) or grouped by dots in threes (1.234); a comma then gives one or two
// digits of centavos (1.234,56, and 1.234,5 for fifty centavos). An amount
// with no comma may instead give its centavos after a single dot, when that
// dot is followed by one or two digits and nothing else (1234.56): that is
// how an HTML number input submits a value, and no grouping of thousands
// leaves fewer than three digits after a dot. A dot followed by three digits
// always groups thousands, so 1.500 is fifteen hundred reais -- a number
// input that should take centavos says step="0.01", which keeps a third
// decimal out of it.
//
// Everything else is refused: a sign, another currency or a lower-case r$, a
// space or a dot inside the digits other than as above, a group of other than
// three digits, a third digit of centavos, a comma with nothing after it, a
// digit outside ASCII, and an amount past the largest int64. Zero is an
// amount; a form that needs a positive one says so itself.
func ParseBRL(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if rest, ok := strings.CutPrefix(s, "R$"); ok {
		s = strings.TrimLeftFunc(rest, unicode.IsSpace)
	}

	whole, fraction := s, ""
	if before, after, ok := strings.Cut(s, ","); ok {
		whole, fraction = before, after
		if fraction == "" || !reaisDigits(whole) {
			return 0, ErrInvalidBRL
		}
	} else if !reaisDigits(whole) {
		dot := strings.LastIndexByte(s, '.')
		if dot < 0 || len(s)-dot-1 == 0 || len(s)-dot-1 > 2 {
			return 0, ErrInvalidBRL
		}
		whole, fraction = s[:dot], s[dot+1:]
		if whole == "" || !allDigits(whole) {
			return 0, ErrInvalidBRL
		}
	}
	if len(fraction) > 2 || !allDigits(fraction) {
		return 0, ErrInvalidBRL
	}

	var cents int64
	for _, digits := range []string{strings.ReplaceAll(whole, ".", ""), (fraction + "00")[:2]} {
		for i := 0; i < len(digits); i++ {
			d := int64(digits[i] - '0')
			if cents > (math.MaxInt64-d)/10 {
				return 0, ErrInvalidBRL
			}
			cents = cents*10 + d
		}
	}
	return cents, nil
}

// ValidBRL reports whether ParseBRL accepts s.
func ValidBRL(s string) bool {
	_, err := ParseBRL(s)
	return err == nil
}

// reaisDigits reports whether s is the reais of an amount: one digit or more,
// bare, or grouped by dots with one to three digits before the first dot and
// three after each.
func reaisDigits(s string) bool {
	groups := strings.Split(s, ".")
	if len(groups) == 1 {
		return s != "" && allDigits(s)
	}
	if len(groups[0]) == 0 || len(groups[0]) > 3 || !allDigits(groups[0]) {
		return false
	}
	for _, g := range groups[1:] {
		if len(g) != 3 || !allDigits(g) {
			return false
		}
	}
	return true
}

// allDigits reports whether every byte of s is an ASCII digit. The empty
// string has none that is not, and the callers that need a digit say so.
func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
