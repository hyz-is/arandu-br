package unit_test

import (
	"errors"
	"math/rand/v2"
	"regexp"
	"strings"
	"testing"

	br "github.com/hyz-is/arandu-br"
)

// Every valid document in this file is generated here, by the oracle below,
// from a fixed seed. None is copied from anywhere: a number that passes the
// check digits may belong to a person, and a test suite is published.
//
// The oracle is written apart from the package, in the textbook form of each
// algorithm -- explicit weight tables, eleven minus the remainder -- so that a
// mistake in the package's own arithmetic is a disagreement between two
// computations rather than one computation agreeing with itself.

// The shapes, as regular expressions. RE2's \d is ASCII only, which is the
// point: a full-width digit is not a digit of a document.
var (
	cpfShape  = regexp.MustCompile(`^(?:\d{11}|\d{3}\.\d{3}\.\d{3}-\d{2})$`)
	cnpjShape = regexp.MustCompile(`^(?:[0-9A-Za-z]{12}\d{2}|[0-9A-Za-z]{2}\.[0-9A-Za-z]{3}\.[0-9A-Za-z]{3}/[0-9A-Za-z]{4}-\d{2})$`)
	cepShape  = regexp.MustCompile(`^(?:\d{8}|\d{5}-\d{3})$`)
)

// bareOf drops the punctuation a shape allows.
func bareOf(s string) string {
	return strings.NewReplacer(".", "", "-", "", "/", "").Replace(s)
}

// oracleCPFDigit is one CPF check digit: weights from len+1 down to 2, eleven
// minus the remainder, and zero when that is ten or eleven.
func oracleCPFDigit(digits string) byte {
	sum := 0
	for i, weight := 0, len(digits)+1; i < len(digits); i, weight = i+1, weight-1 {
		sum += int(digits[i]-'0') * weight
	}
	d := 11 - sum%11
	if d >= 10 {
		d = 0
	}
	return byte('0' + d)
}

func oracleCPF(base string) string {
	first := oracleCPFDigit(base)
	return base + string(first) + string(oracleCPFDigit(base+string(first)))
}

// oracleCPFValid is the whole rule, decided without the package.
func oracleCPFValid(s string) bool {
	if !cpfShape.MatchString(s) {
		return false
	}
	d := bareOf(s)
	if strings.Count(d, d[:1]) == len(d) {
		return false
	}
	return oracleCPF(d[:9]) == d
}

var (
	cnpjFirstWeights  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondWeights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

// oracleCNPJDigit is one CNPJ check digit from the Receita Federal's weight
// table: each character worth its code minus 48, a remainder below two giving
// zero.
func oracleCNPJDigit(chars string, weights []int) byte {
	sum := 0
	for i := range chars {
		sum += (int(chars[i]) - 48) * weights[i]
	}
	if r := sum % 11; r >= 2 {
		return byte('0' + 11 - r)
	}
	return '0'
}

func oracleCNPJ(base string) string {
	base = strings.ToUpper(base)
	first := oracleCNPJDigit(base, cnpjFirstWeights)
	with := base + string(first)
	return with + string(oracleCNPJDigit(with, cnpjSecondWeights))
}

func oracleCNPJValid(s string) bool {
	if !cnpjShape.MatchString(s) {
		return false
	}
	c := strings.ToUpper(bareOf(s))
	if strings.Count(c, c[:1]) == len(c) {
		return false
	}
	return oracleCNPJ(c[:12]) == c
}

func oracleCEPValid(s string) bool {
	return cepShape.MatchString(s) && bareOf(s) >= "01000000"
}

// generator is the one source of the documents in this file, seeded so that a
// failure names the same document on every run.
func generator() *rand.Rand { return rand.New(rand.NewPCG(2026, 1009)) }

func randomFrom(r *rand.Rand, alphabet string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[r.IntN(len(alphabet))]
	}
	return string(b)
}

const (
	digitChars        = "0123456789"
	alphanumericChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// generatedCPFs are valid CPFs, bare.
func generatedCPFs(n int) []string {
	r := generator()
	out := make([]string, 0, n)
	for len(out) < n {
		cpf := oracleCPF(randomFrom(r, digitChars, 9))
		if oracleCPFValid(cpf) {
			out = append(out, cpf)
		}
	}
	return out
}

// generatedCNPJs are valid CNPJs, bare: the first half numeric, the second
// alphanumeric.
func generatedCNPJs(n int) []string {
	r := generator()
	out := make([]string, 0, n)
	for len(out) < n {
		alphabet := digitChars
		if len(out) >= n/2 {
			alphabet = alphanumericChars
		}
		cnpj := oracleCNPJ(randomFrom(r, alphabet, 12))
		if oracleCNPJValid(cnpj) {
			out = append(out, cnpj)
		}
	}
	return out
}

func punctuatedCPF(d string) string {
	return d[0:3] + "." + d[3:6] + "." + d[6:9] + "-" + d[9:11]
}

func punctuatedCNPJ(c string) string {
	return c[0:2] + "." + c[2:5] + "." + c[5:8] + "/" + c[8:12] + "-" + c[12:14]
}

func TestTheOracleAgreesWithThePublishedAlphanumericExample(t *testing.T) {
	t.Parallel()

	// The example the Receita Federal publishes with the alphanumeric CNPJ. It
	// is a specimen, not an establishment, and it is here because it is the
	// one value that checks the oracle against the source rather than against
	// the package.
	const specimen = "12.ABC.345/01DE-35"
	if !oracleCNPJValid(specimen) {
		t.Fatalf("the oracle refuses the published specimen %s", specimen)
	}
	if got, err := br.ParseCNPJ(specimen); err != nil || got != "12ABC34501DE35" {
		t.Fatalf("ParseCNPJ(%q) = %q, %v; want 12ABC34501DE35", specimen, got, err)
	}
}

func TestParseCPFAcceptsEveryGeneratedCPFInBothSpellings(t *testing.T) {
	t.Parallel()

	for _, cpf := range generatedCPFs(2000) {
		for _, written := range []string{cpf, punctuatedCPF(cpf)} {
			got, err := br.ParseCPF(written)
			if err != nil || got != cpf {
				t.Fatalf("ParseCPF(%q) = %q, %v; want %q", written, got, err, cpf)
			}
			if !br.ValidCPF(written) {
				t.Fatalf("ValidCPF(%q) = false", written)
			}
		}
	}
}

// TestParseCPFAgreesWithTheOracleOnEverySubstitution changes every position of
// every generated CPF to every other digit. Most of those are refused, and not
// all: a check digit that computes as ten is written as zero, so a substitution
// that moves a remainder between ten and eleven passes. The oracle decides
// which, and the package has to agree with it each time.
func TestParseCPFAgreesWithTheOracleOnEverySubstitution(t *testing.T) {
	t.Parallel()

	refused := 0
	for _, cpf := range generatedCPFs(300) {
		for i := range cpf {
			for d := byte('0'); d <= '9'; d++ {
				if cpf[i] == d {
					continue
				}
				changed := cpf[:i] + string(d) + cpf[i+1:]
				if got, want := br.ValidCPF(changed), oracleCPFValid(changed); got != want {
					t.Fatalf("ValidCPF(%q) = %v, the oracle says %v", changed, got, want)
				}
				if !oracleCPFValid(changed) {
					refused++
				}
			}
		}
	}
	if refused == 0 {
		t.Fatal("no substitution was refused, so the comparison above proved nothing")
	}
}

func TestParseCPFRefuses(t *testing.T) {
	t.Parallel()

	cpf := generatedCPFs(1)[0]
	wrongCheck := cpf[:10] + string('0'+(cpf[10]-'0'+1)%10)
	for _, tc := range []struct{ name, value string }{
		{"empty", ""},
		{"wrong check digit", wrongCheck},
		{"wrong check digit, punctuated", punctuatedCPF(wrongCheck)},
		{"ten digits", cpf[:10]},
		{"twelve digits", cpf + "0"},
		{"leading space", " " + cpf},
		{"trailing space", punctuatedCPF(cpf) + " "},
		{"a slash for the dash", cpf[0:3] + "." + cpf[3:6] + "." + cpf[6:9] + "/" + cpf[9:]},
		{"dots in the wrong place", cpf[0:4] + "." + cpf[4:7] + "." + cpf[7:9] + "-" + cpf[9:]},
		{"only the dash", cpf[:9] + "-" + cpf[9:]},
		{"a letter", "A" + cpf[1:]},
		{"full-width digits", strings.Map(func(r rune) rune { return r - '0' + '０' }, cpf)},
		{"all zeros", "00000000000"},
		{"all ones, punctuated", "111.111.111-11"},
		{"all nines", "99999999999"},
	} {
		if got, err := br.ParseCPF(tc.value); !errors.Is(err, br.ErrInvalidCPF) || got != "" {
			t.Errorf("%s: ParseCPF(%q) = %q, %v; want ErrInvalidCPF", tc.name, tc.value, got, err)
		}
	}
}

func TestRepeatedDigitsComputeTheirOwnCheckDigits(t *testing.T) {
	t.Parallel()

	// The reason the repetitions are refused by name: their check digits
	// compute, so the arithmetic alone would accept them.
	for d := byte('0'); d <= '9'; d++ {
		base := strings.Repeat(string(d), 9)
		if check, err := br.CPFCheckDigits(base); err != nil || check != strings.Repeat(string(d), 2) {
			t.Errorf("CPFCheckDigits(%q) = %q, %v; want %q", base, check, err, strings.Repeat(string(d), 2))
		}
	}
}

func TestCPFCheckDigitsAgreesWithTheOracle(t *testing.T) {
	t.Parallel()

	for _, cpf := range generatedCPFs(500) {
		if got, err := br.CPFCheckDigits(cpf[:9]); err != nil || got != cpf[9:] {
			t.Fatalf("CPFCheckDigits(%q) = %q, %v; want %q", cpf[:9], got, err, cpf[9:])
		}
	}
	for _, base := range []string{"", "12345678", "1234567890", "12345678A", "123.456.789"} {
		if _, err := br.CPFCheckDigits(base); !errors.Is(err, br.ErrInvalidCPF) {
			t.Errorf("CPFCheckDigits(%q) error = %v; want ErrInvalidCPF", base, err)
		}
	}
}

func TestParseCNPJAcceptsEveryGeneratedCNPJ(t *testing.T) {
	t.Parallel()

	for _, cnpj := range generatedCNPJs(2000) {
		for _, written := range []string{
			cnpj, punctuatedCNPJ(cnpj),
			strings.ToLower(cnpj), strings.ToLower(punctuatedCNPJ(cnpj)),
		} {
			got, err := br.ParseCNPJ(written)
			if err != nil || got != cnpj {
				t.Fatalf("ParseCNPJ(%q) = %q, %v; want %q", written, got, err, cnpj)
			}
			if !br.ValidCNPJ(written) {
				t.Fatalf("ValidCNPJ(%q) = false", written)
			}
		}
	}
}

func TestParseCNPJAgreesWithTheOracleOnEverySubstitution(t *testing.T) {
	t.Parallel()

	refused := 0
	for _, cnpj := range generatedCNPJs(100) {
		for i := range cnpj {
			for _, c := range []byte(alphanumericChars) {
				if cnpj[i] == c {
					continue
				}
				changed := cnpj[:i] + string(c) + cnpj[i+1:]
				if got, want := br.ValidCNPJ(changed), oracleCNPJValid(changed); got != want {
					t.Fatalf("ValidCNPJ(%q) = %v, the oracle says %v", changed, got, want)
				}
				if !oracleCNPJValid(changed) {
					refused++
				}
			}
		}
	}
	if refused == 0 {
		t.Fatal("no substitution was refused, so the comparison above proved nothing")
	}
}

func TestParseCNPJRefuses(t *testing.T) {
	t.Parallel()

	cnpj := generatedCNPJs(2)[1]
	wrongCheck := cnpj[:13] + string('0'+(cnpj[13]-'0'+1)%10)
	for _, tc := range []struct{ name, value string }{
		{"empty", ""},
		{"wrong check digit", wrongCheck},
		{"a letter where a check digit goes", cnpj[:13] + "A"},
		{"thirteen characters", cnpj[:13]},
		{"fifteen characters", cnpj + "0"},
		{"trailing space", punctuatedCNPJ(cnpj) + " "},
		{"a dash for the slash", cnpj[0:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "-" + cnpj[8:12] + "-" + cnpj[12:]},
		{"an accented letter", "Á" + cnpj[2:]},
		{"a symbol", "*" + cnpj[1:]},
		{"all zeros", "00000000000000"},
		{"all zeros, punctuated", "00.000.000/0000-00"},
	} {
		if got, err := br.ParseCNPJ(tc.value); !errors.Is(err, br.ErrInvalidCNPJ) || got != "" {
			t.Errorf("%s: ParseCNPJ(%q) = %q, %v; want ErrInvalidCNPJ", tc.name, tc.value, got, err)
		}
	}
}

func TestCNPJCheckDigitsAgreesWithTheOracle(t *testing.T) {
	t.Parallel()

	for _, cnpj := range generatedCNPJs(500) {
		for _, base := range []string{cnpj[:12], strings.ToLower(cnpj[:12])} {
			if got, err := br.CNPJCheckDigits(base); err != nil || got != cnpj[12:] {
				t.Fatalf("CNPJCheckDigits(%q) = %q, %v; want %q", base, got, err, cnpj[12:])
			}
		}
	}
	for _, base := range []string{"", "12345678901", "1234567890123", "12.345.678/0001", "12345678901_"} {
		if _, err := br.CNPJCheckDigits(base); !errors.Is(err, br.ErrInvalidCNPJ) {
			t.Errorf("CNPJCheckDigits(%q) error = %v; want ErrInvalidCNPJ", base, err)
		}
	}
}

func TestParseCEP(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		value string
		want  string
	}{
		{"01000-000", "01000000"},
		{"01000000", "01000000"},
		{"99999-999", "99999999"},
		{"70040-010", "70040010"},
	} {
		if got, err := br.ParseCEP(tc.value); err != nil || got != tc.want {
			t.Errorf("ParseCEP(%q) = %q, %v; want %q", tc.value, got, err, tc.want)
		}
		if !br.ValidCEP(tc.value) {
			t.Errorf("ValidCEP(%q) = false", tc.value)
		}
	}
	for _, value := range []string{
		"", "00000-000", "00000000", "00999-999", "0100000", "010000000",
		"01000 000", "01.000-000", "1000-0000", "01000-00A", " 01000-000",
	} {
		if got, err := br.ParseCEP(value); !errors.Is(err, br.ErrInvalidCEP) || got != "" {
			t.Errorf("ParseCEP(%q) = %q, %v; want ErrInvalidCEP", value, got, err)
		}
	}
}

func TestTheFormattersPunctuateWithoutValidating(t *testing.T) {
	t.Parallel()

	cpf := generatedCPFs(1)[0]
	cnpj := generatedCNPJs(2)[1]
	for _, tc := range []struct {
		name, got, want string
	}{
		{"CPF bare", br.CPF(cpf), punctuatedCPF(cpf)},
		{"CPF punctuated", br.CPF(punctuatedCPF(cpf)), punctuatedCPF(cpf)},
		{"CPF with a wrong check digit", br.CPF("12345678900"), "123.456.789-00"},
		{"CPF too short", br.CPF("1234"), "1234"},
		{"CPF empty", br.CPF(""), ""},
		{"CPF with letters", br.CPF("not a cpf"), "not a cpf"},
		{"CNPJ bare", br.CNPJ(cnpj), punctuatedCNPJ(cnpj)},
		{"CNPJ lower case", br.CNPJ(strings.ToLower(cnpj)), punctuatedCNPJ(cnpj)},
		{"CNPJ punctuated", br.CNPJ(punctuatedCNPJ(cnpj)), punctuatedCNPJ(cnpj)},
		{"CNPJ too long", br.CNPJ(cnpj + "1"), cnpj + "1"},
		{"CEP bare", br.CEP("70040010"), "70040-010"},
		{"CEP punctuated", br.CEP("70040-010"), "70040-010"},
		{"CEP below every range", br.CEP("00000000"), "00000-000"},
		{"CEP with a space", br.CEP("70040 010"), "70040 010"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

// The fuzz targets hold each parser against the oracle on input nobody chose:
// they agree on every value, what is accepted comes back bare and survives a
// round trip through its formatter, and a formatter applied twice changes
// nothing the second time. The seeds run on every go test.

func FuzzParseCPF(f *testing.F) {
	for _, cpf := range generatedCPFs(8) {
		f.Add(cpf)
		f.Add(punctuatedCPF(cpf))
	}
	first := generatedCPFs(1)[0]
	for _, seed := range []string{
		"", "111.111.111-11", " " + first, first[:9] + "/" + first[9:], "１２３",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, err := br.ParseCPF(s)
		if (err == nil) != oracleCPFValid(s) {
			t.Fatalf("ParseCPF(%q) error = %v, the oracle says valid = %v", s, err, oracleCPFValid(s))
		}
		formatted := br.CPF(s)
		if br.CPF(formatted) != formatted {
			t.Fatalf("CPF is not idempotent on %q: %q then %q", s, formatted, br.CPF(formatted))
		}
		if err != nil {
			if got != "" || !errors.Is(err, br.ErrInvalidCPF) {
				t.Fatalf("ParseCPF(%q) = %q, %v on refusal", s, got, err)
			}
			return
		}
		if got != bareOf(s) || len(got) != 11 {
			t.Fatalf("ParseCPF(%q) = %q, want the bare digits", s, got)
		}
		if again, err := br.ParseCPF(br.CPF(got)); err != nil || again != got {
			t.Fatalf("round trip of %q through CPF gave %q, %v", got, again, err)
		}
	})
}

func FuzzParseCNPJ(f *testing.F) {
	for _, cnpj := range generatedCNPJs(8) {
		f.Add(cnpj)
		f.Add(strings.ToLower(punctuatedCNPJ(cnpj)))
	}
	for _, seed := range []string{"", "00.000.000/0000-00", "12.ABC.345/01DE-35", "12ABC34501DE3A", "Á"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, err := br.ParseCNPJ(s)
		if (err == nil) != oracleCNPJValid(s) {
			t.Fatalf("ParseCNPJ(%q) error = %v, the oracle says valid = %v", s, err, oracleCNPJValid(s))
		}
		formatted := br.CNPJ(s)
		if br.CNPJ(formatted) != formatted {
			t.Fatalf("CNPJ is not idempotent on %q: %q then %q", s, formatted, br.CNPJ(formatted))
		}
		if err != nil {
			if got != "" || !errors.Is(err, br.ErrInvalidCNPJ) {
				t.Fatalf("ParseCNPJ(%q) = %q, %v on refusal", s, got, err)
			}
			return
		}
		if got != strings.ToUpper(bareOf(s)) || len(got) != 14 {
			t.Fatalf("ParseCNPJ(%q) = %q, want the bare characters in capitals", s, got)
		}
		if again, err := br.ParseCNPJ(br.CNPJ(got)); err != nil || again != got {
			t.Fatalf("round trip of %q through CNPJ gave %q, %v", got, again, err)
		}
	})
}

func FuzzParseCEP(f *testing.F) {
	for _, seed := range []string{"", "01000-000", "00999999", "70040010", "70040 010", "7004-0010"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, err := br.ParseCEP(s)
		if (err == nil) != oracleCEPValid(s) {
			t.Fatalf("ParseCEP(%q) error = %v, the oracle says valid = %v", s, err, oracleCEPValid(s))
		}
		formatted := br.CEP(s)
		if br.CEP(formatted) != formatted {
			t.Fatalf("CEP is not idempotent on %q: %q then %q", s, formatted, br.CEP(formatted))
		}
		if err != nil {
			if got != "" || !errors.Is(err, br.ErrInvalidCEP) {
				t.Fatalf("ParseCEP(%q) = %q, %v on refusal", s, got, err)
			}
			return
		}
		if got != bareOf(s) || len(got) != 8 {
			t.Fatalf("ParseCEP(%q) = %q, want the bare digits", s, got)
		}
		if again, err := br.ParseCEP(br.CEP(got)); err != nil || again != got {
			t.Fatalf("round trip of %q through CEP gave %q, %v", got, again, err)
		}
	})
}
