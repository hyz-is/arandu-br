package br

import (
	"errors"
	"strings"
)

// ErrInvalidStateRegistration is what ParseStateRegistration returns for a
// value that is not a state registration number of the unit it was asked
// about.
var ErrInvalidStateRegistration = errors.New("br: not a valid state registration number")

// ParseStateRegistration returns the digits of an inscrição estadual written
// for the unit uf, and ErrInvalidStateRegistration for anything else,
// including a uf that is not one of the twenty-seven.
//
// Every unit issues its own number, with its own length, its own prefix and
// its own check digits, so the unit is part of the question: the same digits
// can be valid in one and not in another. The number is read with its digits
// alone or with dots, hyphens and slashes between them, the way each unit
// prints it; any other character, or a separator that starts, ends or doubles,
// is refused. A repetition of one digit is refused.
//
// The value "ISENTO" a company without a registration writes in its place is
// not a number and is not accepted: a form that takes it names it beside the
// rule. It does not say whether the registration is active: only the unit's
// treasury can.
//
// The digits come back bare, which is the spelling to store.
func ParseStateRegistration(uf UF, s string) (string, error) {
	digits, ok := readStateRegistration(s)
	if !ok || repeated(digits) || !stateRegistrationChecks(uf, digits) {
		return "", ErrInvalidStateRegistration
	}
	return digits, nil
}

// ValidStateRegistration reports whether ParseStateRegistration accepts s for
// uf.
func ValidStateRegistration(uf UF, s string) bool {
	_, err := ParseStateRegistration(uf, s)
	return err == nil
}

// readStateRegistration returns the digits of s when s is digits with single
// dots, hyphens or slashes between them.
func readStateRegistration(s string) (string, bool) {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c >= '0' && c <= '9':
			out = append(out, c)
		case (c == '.' || c == '-' || c == '/') && i > 0 && i < len(s)-1 && s[i-1] >= '0' && s[i-1] <= '9':
		default:
			return "", false
		}
	}
	return string(out), len(out) > 0
}

// stateRegistrationChecks decides digits for one unit.
func stateRegistrationChecks(uf UF, d string) bool {
	switch uf {
	case UFAC:
		return len(d) == 13 && strings.HasPrefix(d, "01") && twoDigitsAscending(d)
	case UFAL:
		return len(d) == 9 && strings.HasPrefix(d, "24") && alagoasDigit(d)
	case UFAP:
		return len(d) == 9 && strings.HasPrefix(d, "03") && amapaDigit(d)
	case UFAM:
		return len(d) == 9 && amazonasDigit(d)
	case UFBA:
		return (len(d) == 8 || len(d) == 9) && bahiaDigits(d)
	case UFCE, UFES, UFPB, UFPI, UFSC, UFSE:
		return len(d) == 9 && cearaDigit(d)
	case UFMA:
		return len(d) == 9 && strings.HasPrefix(d, "12") && cearaDigit(d)
	case UFMS:
		return len(d) == 9 && strings.HasPrefix(d, "28") && cearaDigit(d)
	case UFPA:
		return len(d) == 9 && strings.HasPrefix(d, "15") && cearaDigit(d)
	case UFDF:
		return len(d) == 13 && twoDigitsAscending(d)
	case UFGO:
		return len(d) == 9 && (strings.HasPrefix(d, "10") || strings.HasPrefix(d, "11") || strings.HasPrefix(d, "15")) && goiasDigit(d)
	case UFMT:
		return len(d) == 11 && matoGrossoDigit(d)
	case UFMG:
		return len(d) == 13 && minasGeraisDigits(d)
	case UFPR:
		return len(d) == 10 && paranaDigits(d)
	case UFPE:
		return len(d) == 9 && pernambucoDigits(d)
	case UFRJ:
		return len(d) == 8 && rioDeJaneiroDigit(d)
	case UFRN:
		return (len(d) == 9 || len(d) == 10) && strings.HasPrefix(d, "20") && rioGrandeDoNorteDigit(d)
	case UFRS:
		return len(d) == 10 && rioGrandeDoSulDigit(d)
	case UFRO:
		return len(d) == 14 && rondoniaDigit(d)
	case UFRR:
		return len(d) == 9 && strings.HasPrefix(d, "24") && roraimaDigit(d)
	case UFSP:
		return len(d) == 12 && saoPauloDigits(d)
	case UFTO:
		return tocantinsDigit(d)
	}
	return false
}

// cyclicSum weighs the digits of d from start downwards, wrapping to reset
// when the weight reaches floor.
func cyclicSum(d string, start, floor, reset int) int {
	weight, sum := start, 0
	for i := 0; i < len(d); i++ {
		sum += int(d[i]-'0') * weight
		if weight--; weight == floor {
			weight = reset
		}
	}
	return sum
}

// descendingSum weighs the digits of d from start downwards by one.
func descendingSum(d string, start int) int {
	sum := 0
	for i := 0; i < len(d); i++ {
		sum += int(d[i]-'0') * (start - i)
	}
	return sum
}

// eleven is the check digit most units use: eleven less the remainder of the
// sum modulo eleven, and zero when that is ten or more.
func eleven(sum int) byte {
	if digit := 11 - sum%11; digit < 10 {
		return byte('0' + digit)
	}
	return '0'
}

// twoDigitsAscending is the pair of check digits of Acre and the Federal
// District: weights from the length less seven downwards, wrapping to nine.
func twoDigitsAscending(d string) bool {
	body := d[:len(d)-2]
	first := eleven(cyclicSum(body, len(body)-7, 1, 9))
	second := eleven(cyclicSum(body+string(first), len(body)-6, 1, 9))
	return d[len(d)-2] == first && d[len(d)-1] == second
}

func alagoasDigit(d string) bool {
	digit := descendingSum(d[:8], 9) * 10 % 11
	if digit >= 10 {
		digit = 0
	}
	return int(d[8]-'0') == digit
}

func amapaDigit(d string) bool {
	body := d[:8]
	// Two ranges of numbers weigh the sum differently and give a different
	// digit when the remainder leaves eleven.
	extra, fallback := 0, 0
	switch {
	case "03000001" <= body && body <= "03017000":
		extra, fallback = 5, 0
	case "03017001" <= body && body <= "03019022":
		extra, fallback = 9, 1
	}
	digit := 11 - (extra+descendingSum(body, 9))%11
	switch digit {
	case 10:
		digit = 0
	case 11:
		digit = fallback
	}
	return int(d[8]-'0') == digit
}

func amazonasDigit(d string) bool {
	sum := descendingSum(d[:8], 9)
	digit := 11 - sum
	if sum >= 11 {
		if digit = 11 - sum%11; digit >= 10 {
			digit = 0
		}
	}
	return int(d[8]-'0') == digit
}

func bahiaDigits(d string) bool {
	// The first digit of an eight-digit number, or the second of a nine-digit
	// one, tells which modulus is used.
	pivot := d[0]
	if len(d) == 9 {
		pivot = d[1]
	}
	modulus := 11
	if strings.IndexByte("0123458", pivot) >= 0 {
		modulus = 10
	}
	digit := func(body string) byte {
		rest := modulus - descendingSum(body, len(body)+1)%modulus
		if rest >= 10 {
			rest = 0
		}
		return byte('0' + rest)
	}
	body := d[:len(d)-2]
	second := digit(body)
	first := digit(body + string(second))
	return d[len(d)-2] == first && d[len(d)-1] == second
}

func cearaDigit(d string) bool {
	return d[8] == eleven(descendingSum(d[:8], 9))
}

func goiasDigit(d string) bool {
	body := d[:8]
	digit := 11 - descendingSum(body, 9)%11
	if digit >= 10 {
		if digit == 11 && "10103105" <= body && body <= "10119997" {
			digit = 1
		} else {
			digit = 0
		}
	}
	return int(d[8]-'0') == digit
}

func matoGrossoDigit(d string) bool {
	return d[10] == eleven(cyclicSum(d[:10], 3, 1, 9))
}

func minasGeraisDigits(d string) bool {
	// The first digit is the sum of the digits of the products of the body, with
	// a zero put in front of the fourth digit, weighted 1 and 2 in turn.
	body := d[:3] + "0" + d[3:11]
	sum := 0
	for i := 0; i < len(body); i++ {
		weight := 1
		if i%2 == 1 {
			weight = 2
		}
		for product := int(body[i]-'0') * weight; product > 0; product /= 10 {
			sum += product % 10
		}
	}
	first := byte('0')
	if last := sum % 10; last != 0 {
		first = byte('0' + 10 - last)
	}
	second := eleven(cyclicSum(d[:11]+string(first), 3, 1, 11))
	return d[11] == first && d[12] == second
}

func paranaDigits(d string) bool {
	body := d[:8]
	first := eleven(cyclicSum(body, len(body)-5, 1, 7))
	second := eleven(cyclicSum(body+string(first), len(body)-4, 1, 7))
	return d[8] == first && d[9] == second
}

func pernambucoDigits(d string) bool {
	body := d[:7]
	first := eleven(descendingSum(body, 8))
	second := eleven(descendingSum(body+string(first), 9))
	return d[7] == first && d[8] == second
}

func rioDeJaneiroDigit(d string) bool {
	return d[7] == eleven(cyclicSum(d[:7], 2, 1, 7))
}

func rioGrandeDoNorteDigit(d string) bool {
	last := len(d) - 1
	digit := descendingSum(d[:last], len(d)) * 10 % 11
	if digit == 10 {
		digit = 0
	}
	return int(d[last]-'0') == digit
}

func rioGrandeDoSulDigit(d string) bool {
	return d[9] == eleven(cyclicSum(d[:9], 2, 1, 9))
}

func rondoniaDigit(d string) bool {
	digit := 11 - cyclicSum(d[:13], 6, 1, 9)%11
	if digit >= 10 {
		digit -= 10
	}
	return int(d[13]-'0') == digit
}

func roraimaDigit(d string) bool {
	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(d[i]-'0') * (i + 1)
	}
	return int(d[8]-'0') == sum%9
}

func saoPauloDigits(d string) bool {
	weights := [8]int{1, 3, 4, 5, 6, 7, 8, 10}
	sum := 0
	for i, weight := range weights {
		sum += int(d[i]-'0') * weight
	}
	// The digit is the last one of the remainder, so ten is written as zero.
	first := byte('0' + sum%11%10)
	second := byte('0' + cyclicSum(d[:11], 3, 1, 10)%11%10)
	return d[8] == first && d[11] == second
}

func tocantinsDigit(d string) bool {
	switch len(d) {
	case 11:
		// The older number carries a category in the third and fourth digits,
		// left out of the digits that are weighed.
		switch d[2:4] {
		case "01", "02", "03", "99":
			body := d[:2] + d[4:]
			return cearaDigit(body)
		}
	case 9:
		return cearaDigit(d)
	}
	return false
}
