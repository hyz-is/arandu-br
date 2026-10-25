package unit_test

import (
	"strconv"
	"strings"
	"testing"

	br "github.com/hyz-is/arandu-br"
)

// The numbers in this file are generated from a fixed seed, as the CPFs are,
// and judged by an oracle written apart from the package: a transliteration of
// the algorithm as the PHP packages state it, in the form they state it, so
// that the package's own arithmetic is checked against a second computation
// rather than against itself. None is the number of anybody.

// substitutions lists every string that differs from s in one character, the
// character replaced by each of the digits.
func substitutions(s string) []string {
	out := make([]string, 0, len(s)*9)
	for i := 0; i < len(s); i++ {
		for d := byte('0'); d <= '9'; d++ {
			if s[i] != d {
				out = append(out, s[:i]+string(d)+s[i+1:])
			}
		}
	}
	return out
}

// agree holds a Parse function to the oracle on a set of generated numbers and
// on every single-digit substitution of each: the same answer, and the digits
// back when it is yes.
func agree(t *testing.T, name string, parse func(string) (string, error), oracle func(string) bool, numbers []string) {
	t.Helper()
	accepted := 0
	for _, number := range numbers {
		for _, candidate := range append([]string{number}, substitutions(number)...) {
			got, err := parse(candidate)
			want := oracle(candidate)
			if (err == nil) != want {
				t.Fatalf("%s(%q) accepted = %v, oracle says %v", name, candidate, err == nil, want)
			}
			if err == nil {
				accepted++
				if got != candidate {
					t.Fatalf("%s(%q) = %q, want the digits back", name, candidate, got)
				}
			}
		}
	}
	if accepted < len(numbers) {
		t.Fatalf("%s accepted %d of %d generated numbers", name, accepted, len(numbers))
	}
}

func digitsOf(r interface{ IntN(int) int }, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + r.IntN(10))
	}
	return string(b)
}

// oracleNIS is the NIS as geekcom/validator-docs computes it: the digits from
// the tenth back to the first weighted 2 to 9 and over again.
func oracleNIS(s string) bool {
	if len(s) != 11 || strings.Count(s, s[:1]) == 11 {
		return false
	}
	d, p := 0, 2
	for c := 9; c >= 0; c-- {
		d += int(s[c]-'0') * p
		if p < 9 {
			p++
		} else {
			p = 2
		}
	}
	return int(s[10]-'0') == (10*d%11)%10
}

func TestParseNISAgreesWithTheOracle(t *testing.T) {
	t.Parallel()
	r := generator()
	var numbers []string
	for len(numbers) < 60 {
		base := digitsOf(r, 10)
		for d := byte('0'); d <= '9'; d++ {
			if oracleNIS(base + string(d)) {
				numbers = append(numbers, base+string(d))
			}
		}
	}
	agree(t, "ParseNIS", br.ParseNIS, oracleNIS, numbers)

	got, err := br.ParseNIS(br.NIS(numbers[0]))
	if err != nil || got != numbers[0] {
		t.Errorf("ParseNIS(NIS(%q)) = %q, %v", numbers[0], got, err)
	}
	if got := br.NIS(numbers[0]); len(got) != 14 || got[3] != '.' || got[9] != '.' || got[12] != '-' {
		t.Errorf("NIS(%q) = %q, want 000.00000.00-0", numbers[0], got)
	}
	if got := br.NIS("123"); got != "123" {
		t.Errorf("NIS(%q) = %q, want it unchanged", "123", got)
	}
}

// oracleCNS is the CNS as geekcom/validator-docs states it: a definitive card
// matches one pattern, a provisional card another, and both weigh to a
// multiple of eleven.
func oracleCNS(s string) bool {
	if len(s) != 15 {
		return false
	}
	definitive := (s[0] == '1' || s[0] == '2') && s[11] == '0' && s[12] == '0' && (s[13] == '0' || s[13] == '1')
	provisional := s[0] >= '7' && s[0] <= '9'
	if !definitive && !provisional {
		return false
	}
	sum := 0
	for i := 0; i < 15; i++ {
		sum += int(s[i]-'0') * (15 - i)
	}
	return sum%11 == 0
}

func TestParseCNSAgreesWithTheOracle(t *testing.T) {
	t.Parallel()
	r := generator()
	var numbers []string
	for len(numbers) < 60 {
		// A definitive card is the NIS it was issued from, 00, 0 or 1, and a
		// check digit; a provisional one is fifteen digits that add up.
		var base string
		if len(numbers)%2 == 0 {
			base = string("12"[r.IntN(2)]) + digitsOf(r, 10) + "00" + string("01"[r.IntN(2)])
		} else {
			base = string("789"[r.IntN(3)]) + digitsOf(r, 13)
		}
		for d := byte('0'); d <= '9'; d++ {
			if oracleCNS(base + string(d)) {
				numbers = append(numbers, base+string(d))
				break
			}
		}
	}
	agree(t, "ParseCNS", br.ParseCNS, oracleCNS, numbers)

	for _, in := range numbers[:4] {
		spaced := in[:3] + " " + in[3:7] + " " + in[7:11] + " " + in[11:]
		if got, err := br.ParseCNS(spaced); err != nil || got != in {
			t.Errorf("ParseCNS(%q) = %q, %v", spaced, got, err)
		}
		if got := br.CNS(in); got != spaced {
			t.Errorf("CNS(%q) = %q, want %q", in, got, spaced)
		}
	}
	for _, in := range []string{"", "0", "1111111111111111", "000000000000000", "111111111111111", "999999999999999", "116 3876 9194  0009"} {
		if _, err := br.ParseCNS(in); err == nil {
			t.Errorf("ParseCNS(%q) accepted", in)
		}
	}
}

// oracleCNH is the CNH as geekcom/validator-docs computes it: the first check
// digit from the nine digits weighted 2 to 10, the second from the ten that
// result when it is put in front, weighted 2 to 11.
func oracleCNH(s string) bool {
	if len(s) != 11 {
		return false
	}
	if n, _ := strconv.Atoi(s); n == 0 {
		return false
	}
	check := func(partial string) int {
		sum := 0
		for i, w := 0, 2; i < len(partial); i, w = i+1, w+1 {
			sum += int(partial[i]-'0') * w
		}
		if rest := sum % 11; rest > 1 {
			return 11 - rest
		}
		return 0
	}
	first := check(s[:9])
	second := check(strconv.Itoa(first) + s[:9])
	return strconv.Itoa(first)+strconv.Itoa(second) == s[9:] && strings.Count(s, s[:1]) != 11
}

func TestParseCNHAgreesWithTheOracle(t *testing.T) {
	t.Parallel()
	r := generator()
	var numbers []string
	for len(numbers) < 80 {
		base := digitsOf(r, 9)
		for a := byte('0'); a <= '9'; a++ {
			for b := byte('0'); b <= '9'; b++ {
				if oracleCNH(base + string(a) + string(b)) {
					numbers = append(numbers, base+string(a)+string(b))
				}
			}
		}
	}
	agree(t, "ParseCNH", br.ParseCNH, oracleCNH, numbers)

	// The register number whose second digit needs the correction that adds
	// eleven back: the first digit came from a ten, and the second sum left a
	// remainder below the discount.
	if !br.ValidCNH("04463004100") {
		t.Error("ParseCNH refused 04463004100, which only an implementation without the correction refuses")
	}
	for _, in := range []string{"", "00000000000", "11111111111", "1234567890", "123456789012", "96784547 43", "967.845.479-43"} {
		if br.ValidCNH(in) {
			t.Errorf("ParseCNH(%q) accepted", in)
		}
	}
}

// oracleRenavam is the RENAVAM as geekcom/validator-docs computes it: the ten
// digits reversed, weighted 2 to 9 and then 2 and 3.
func oracleRenavam(s string) bool {
	if len(s) < 9 || len(s) > 11 {
		return false
	}
	s = strings.Repeat("0", 11-len(s)) + s
	if strings.Count(s, s[:1]) == 11 {
		return false
	}
	reversed := []byte(s[:10])
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(reversed[i]-'0') * (i + 2)
	}
	sum += int(reversed[8]-'0')*2 + int(reversed[9]-'0')*3
	check := 11 - sum%11
	if check >= 10 {
		check = 0
	}
	return int(s[10]-'0') == check
}

func TestParseRenavamAgreesWithTheOracle(t *testing.T) {
	t.Parallel()
	r := generator()
	var numbers []string
	for len(numbers) < 60 {
		base := digitsOf(r, 10)
		for d := byte('0'); d <= '9'; d++ {
			if oracleRenavam(base + string(d)) {
				numbers = append(numbers, base+string(d))
			}
		}
	}
	agree(t, "ParseRenavam", br.ParseRenavam, oracleRenavam, numbers)

	// An older number has nine or ten digits and is read with the zeros it
	// lost; it comes back with them.
	for _, in := range []string{"00197073212", "0197073212", "197073212"} {
		if got, err := br.ParseRenavam(in); err != nil || got != "00197073212" {
			t.Errorf("ParseRenavam(%q) = %q, %v; want 00197073212", in, got, err)
		}
	}
	for _, in := range []string{"", "19707321", "1970.73212", "1234555582", "00000000000", "11111111111", "000000000", "197073212123", "1970732 2"} {
		if br.ValidRenavam(in) {
			t.Errorf("ParseRenavam(%q) accepted", in)
		}
	}
}

// oracleVoterTitle is the título de eleitor as geekcom/validator-docs
// computes it, statement by statement.
func oracleVoterTitle(input string) bool {
	if len(input) < 5 || len(input) > 13 {
		return false
	}
	uf := input[len(input)-4 : len(input)-2]
	ufNumber, _ := strconv.Atoi(uf)
	if strings.Count(input, input[:1]) == len(input) || ufNumber < 1 || ufNumber > 28 {
		return false
	}
	dv := input[len(input)-2:]
	sequence := input[:len(input)-4]
	for i := 0; i < 2; i++ {
		factor, sum := 9, 0
		for j := len(sequence) - 1; j > -1; j-- {
			sum += int(sequence[j]-'0') * factor
			if factor == 2 {
				factor = 10
			}
			factor--
		}
		digit := sum % 11
		if digit == 0 && ufNumber < 3 {
			digit = 1
		} else if digit == 10 {
			digit = 0
		}
		if int(dv[i]-'0') != digit {
			return false
		}
		if i == 0 {
			sequence = uf + strconv.Itoa(digit)
		}
	}
	return true
}

func TestParseVoterTitleAgreesWithTheOracle(t *testing.T) {
	t.Parallel()
	r := generator()
	var numbers []string
	for len(numbers) < 80 {
		length := []int{7, 8, 12, 12, 12, 13}[r.IntN(6)]
		state := []string{"01", "02", "03", "12", "28", "10"}[r.IntN(6)]
		base := digitsOf(r, length-4) + state
		for a := byte('0'); a <= '9'; a++ {
			for b := byte('0'); b <= '9'; b++ {
				if oracleVoterTitle(base + string(a) + string(b)) {
					numbers = append(numbers, base+string(a)+string(b))
				}
			}
		}
	}
	agree(t, "ParseVoterTitle", br.ParseVoterTitle, oracleVoterTitle, numbers)

	// A title of twelve digits is also read in groups of four.
	for _, in := range numbers {
		if len(in) != 12 {
			continue
		}
		spaced := in[:4] + " " + in[4:8] + " " + in[8:]
		if got, err := br.ParseVoterTitle(spaced); err != nil || got != in {
			t.Errorf("ParseVoterTitle(%q) = %q, %v", spaced, got, err)
		}
		break
	}
	for _, in := range []string{"", "1234", "00000000000000", "0000000000", "111111111111", "123456789012345", "30212 60", "3021260 "} {
		if br.ValidVoterTitle(in) {
			t.Errorf("ParseVoterTitle(%q) accepted", in)
		}
	}
}

// oracleCertificate is the certificate number as geekcom/validator-docs
// computes it: the weights start at 32 less the length and climb to 10, then
// fall to 0 and climb again.
func oracleCertificate(s string) bool {
	if len(s) != 32 || strings.Count(s, s[:1]) == 32 {
		return false
	}
	weighted := func(value string) int {
		sum, multiplier := 0, 32-len(value)
		for i := 0; i < len(value); i++ {
			sum += int(value[i]-'0') * multiplier
			multiplier++
			if multiplier > 10 {
				multiplier = 0
			}
		}
		return sum
	}
	number, dv := s[:30], s[30:]
	dv1 := weighted(number) % 11
	if dv1 > 9 {
		dv1 = 1
	}
	dv2 := weighted(number+strconv.Itoa(dv1)) % 11
	if dv2 > 9 {
		dv2 = 1
	}
	return dv == strconv.Itoa(dv1)+strconv.Itoa(dv2)
}

func TestParseCivilCertificateAgreesWithTheOracle(t *testing.T) {
	t.Parallel()
	r := generator()
	var numbers []string
	for len(numbers) < 40 {
		base := digitsOf(r, 30)
		for a := byte('0'); a <= '9'; a++ {
			for b := byte('0'); b <= '9'; b++ {
				if oracleCertificate(base + string(a) + string(b)) {
					numbers = append(numbers, base+string(a)+string(b))
				}
			}
		}
	}
	agree(t, "ParseCivilCertificate", br.ParseCivilCertificate, oracleCertificate, numbers[:40])

	in := numbers[0]
	dotted := in[:6] + "." + in[6:8] + "." + in[8:10] + "." + in[10:14] + "." + in[14:15] + "." + in[15:20] + "." + in[20:23] + "." + in[23:30] + "-" + in[30:]
	spaced := strings.NewReplacer(".", " ", "-", " ").Replace(dotted)
	mixed := in[:6] + " " + in[6:8] + " " + in[8:10] + " " + in[10:14] + " " + in[14:15] + " " + in[15:20] + " " + in[20:23] + " " + in[23:30] + "-" + in[30:]
	for _, spelling := range []string{dotted, spaced, mixed, in} {
		if got, err := br.ParseCivilCertificate(spelling); err != nil || got != in {
			t.Errorf("ParseCivilCertificate(%q) = %q, %v; want %q", spelling, got, err, in)
		}
	}
	for _, bad := range []string{
		"", in[:31], in + "0",
		strings.Replace(dotted, "-", ".", 1),           // the last group follows a hyphen or a space
		strings.Replace(dotted, ".", "-", 1),           // the others a dot or a space
		strings.Replace(dotted, ".", "/", 1),           // nothing else
		strings.Replace(dotted, ".", "  ", 1),          // not two
		dotted + " ",                                   // nothing around it
		" " + dotted,                                   //
		dotted[:len(dotted)-3],                         // the last group is missing
		strings.Replace(dotted, in[:6], "00000A", 1),   // a letter
		strings.Replace(spaced, in[:6]+" ", in[:6], 1), // a separator missing
		strings.Repeat("0", 32),                        // its check digits compute, and it belongs to nobody
	} {
		if br.ValidCivilCertificate(bad) {
			t.Errorf("ParseCivilCertificate(%q) accepted", bad)
		}
	}
}
