package br

import (
	"strconv"
	"strings"
)

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
func BRL(cents int64) string {
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
		if i > 0 && (len(reais)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteByte(reais[i])
	}
	b.WriteByte(',')
	b.WriteByte(byte('0' + centavos/10))
	b.WriteByte(byte('0' + centavos%10))
	return b.String()
}
