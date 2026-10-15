package unit_test

import (
	"math/rand/v2"
	"strings"
	"testing"

	br "github.com/hyz-is/arandu-br"
)

// The redactions applications printed before this package had them, each
// written here as it was written there and kept as the oracle RedactCPF and
// RedactCNPJ are measured against. Each takes the bare document, which is what
// every one of them was called with.
var legacyRedactions = []struct {
	name     string
	document string // "cpf" or "cnpj"
	pattern  br.Redaction
	legacy   func(string) string
}{
	{
		// A back office's case file: the middle digits, behind bullets.
		"back office CPF", "cpf", "•••.###.###-••",
		func(cpf string) string { return "•••." + cpf[3:6] + "." + cpf[6:9] + "-••" },
	},
	{
		// The same case file's CNPJ, which may be alphanumeric.
		"back office CNPJ", "cnpj", "##.•••.•••/••••-##",
		func(cnpj string) string { return cnpj[:2] + ".•••.•••/••••-" + cnpj[12:] },
	},
	{
		// A claimed CPF written into a log or a stored input: the check digits.
		"document claim", "cpf", br.KeepCPFCheckDigits,
		func(d string) string { return "***.***.***-" + d[9:] },
	},
	{
		// A partner's CPF as the Receita Federal's open CNPJ data publishes it.
		"open CNPJ data partner", "cpf", br.KeepCPFMiddleBare,
		func(plain string) string { return "***" + plain[3:9] + "**" },
	},
	{
		// A public agent's CPF as the Portal da Transparência publishes it.
		"transparency portal", "cpf", br.KeepCPFMiddle,
		func(cpf string) string { return "***." + cpf[3:6] + "." + cpf[6:9] + "-**" },
	},
	{
		// A profile page reminding its owner which CPF it holds.
		"profile page", "cpf", br.KeepCPFLastFive,
		func(cpf string) string { return "***.***." + cpf[6:9] + "-" + cpf[9:] },
	},
	{
		// An API answer that keeps the check digits of either document, bare.
		"API identifier CPF", "cpf", "*********##",
		func(s string) string { return strings.Repeat("*", len(s)-2) + s[len(s)-2:] },
	},
	{
		"API identifier CNPJ", "cnpj", "************##",
		func(s string) string { return strings.Repeat("*", len(s)-2) + s[len(s)-2:] },
	},
}

// unvalidatedCPFs are eleven digits nobody checked: a redaction hides a stored
// value whether or not its check digits hold.
func unvalidatedCPFs(n int) []string {
	r := rand.New(rand.NewPCG(7, 11))
	out := make([]string, n)
	for i := range out {
		out[i] = randomFrom(r, digitChars, 11)
	}
	return out
}

// unvalidatedCNPJs are fourteen characters of the CNPJ's shape nobody checked.
func unvalidatedCNPJs(n int) []string {
	r := rand.New(rand.NewPCG(13, 17))
	out := make([]string, n)
	for i := range out {
		out[i] = randomFrom(r, alphanumericChars, 12) + randomFrom(r, digitChars, 2)
	}
	return out
}

func TestRedactionsReproduceWhatApplicationsPrinted(t *testing.T) {
	t.Parallel()

	cpfs := append(generatedCPFs(40), unvalidatedCPFs(40)...)
	cnpjs := append(generatedCNPJs(40), unvalidatedCNPJs(40)...)
	for _, tc := range legacyRedactions {
		documents, redact, punctuate := cpfs, br.RedactCPF, punctuatedCPF
		if tc.document == "cnpj" {
			documents, redact, punctuate = cnpjs, br.RedactCNPJ, punctuatedCNPJ
		}
		for _, d := range documents {
			want := tc.legacy(d)
			if got := redact(d, tc.pattern); got != want {
				t.Errorf("%s: redacting %q = %q, want %q", tc.name, d, got, want)
			}
			if got := redact(punctuate(d), tc.pattern); got != want {
				t.Errorf("%s: redacting %q = %q, want %q", tc.name, punctuate(d), got, want)
			}
		}
	}
}

// TestRedactionsSpellThePublishedExamples holds the literal strings the
// applications' own tests and comments print. None of the numbers is a CPF:
// each fails its check digits, which the test asserts, so that no person's
// document is in a published suite.
func TestRedactionsSpellThePublishedExamples(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in      string
		pattern br.Redaction
		want    string
	}{
		{"12398224725", "•••.###.###-••", "•••.982.247-••"},
		{"12398224725", br.KeepCPFMiddleBare, "***982247**"},
		{"12398224725", br.KeepCPFLastFive, "***.***.247-25"},
		{"12345678900", br.KeepCPFMiddle, "***.456.789-**"},
		{"123.456.789-00", br.KeepCPFCheckDigits, "***.***.***-00"},
		{"12345678900", "*********##", "*********00"},
	} {
		if oracleCPFValid(tc.in) {
			t.Fatalf("%q is a valid CPF; a published example must not be one", tc.in)
		}
		if got := br.RedactCPF(tc.in, tc.pattern); got != tc.want {
			t.Errorf("RedactCPF(%q, %q) = %q, want %q", tc.in, tc.pattern, got, tc.want)
		}
	}

	// The Receita Federal's published specimen of the alphanumeric CNPJ.
	for in, want := range map[string]string{
		"12ABC34501DE35":     "12.•••.•••/••••-35",
		"12.abc.345/01de-35": "12.•••.•••/••••-35",
	} {
		if got := br.RedactCNPJ(in, "##.•••.•••/••••-##"); got != want {
			t.Errorf("RedactCNPJ(%q) = %q, want %q", in, got, want)
		}
	}
	if got, want := br.RedactCNPJ("12abc34501de35", br.KeepCNPJEnds), "12.***.***/****-35"; got != want {
		t.Errorf("RedactCNPJ(KeepCNPJEnds) = %q, want %q", got, want)
	}
	if got, want := br.RedactCNPJ("12abc34501de35", "##.###.###/####-##"), "12.ABC.345/01DE-35"; got != want {
		t.Errorf("a redaction that shows everything = %q, want the CNPJ as CNPJ writes it, %q", got, want)
	}
}

func TestARedactionChoosesThePunctuationToo(t *testing.T) {
	t.Parallel()

	for pattern, want := range map[br.Redaction]string{
		"***.###.###-**": "***.456.789-**",
		"***-###-###/**": "***-456-789/**",
		"***#########**": "***.456.789-**", // '#' over a separator writes the separator
		"xxx######xx":    "xxx456789xx",
		"###########":    "12345678900",
	} {
		if got := br.RedactCPF("12345678900", pattern); got != want {
			t.Errorf("RedactCPF(%q) = %q, want %q", pattern, got, want)
		}
	}
}

// TestARedactionNeverFallsBackToItsInput is the property the type exists for:
// whatever it cannot read, it answers with nothing rather than with the value
// it was asked to hide.
func TestARedactionNeverFallsBackToItsInput(t *testing.T) {
	t.Parallel()

	cpf := generatedCPFs(1)[0]
	cnpj := generatedCNPJs(1)[0]
	for _, tc := range []struct {
		name string
		got  string
	}{
		{"an empty CPF", br.RedactCPF("", br.KeepCPFMiddle)},
		{"a CPF too short", br.RedactCPF(cpf[:10], br.KeepCPFMiddle)},
		{"a CPF too long", br.RedactCPF(cpf+"0", br.KeepCPFMiddle)},
		{"a CPF with spaces around it", br.RedactCPF(" "+cpf+" ", br.KeepCPFMiddle)},
		{"a CPF with another separator", br.RedactCPF(cpf[:3]+"."+cpf[3:6]+"."+cpf[6:9]+" "+cpf[9:], br.KeepCPFMiddle)},
		{"a CPF with a letter", br.RedactCPF("1234567890A", br.KeepCPFMiddle)},
		{"a CPF already redacted", br.RedactCPF("***456789**", br.KeepCPFMiddleBare)},
		{"a CPF under a redaction of neither length", br.RedactCPF(cpf, "***.###")},
		{"a CPF under the CNPJ's redaction", br.RedactCPF(cpf, br.KeepCNPJEnds)},
		{"a CPF under an empty redaction", br.RedactCPF(cpf, "")},
		{"an empty CNPJ", br.RedactCNPJ("", br.KeepCNPJEnds)},
		{"a CNPJ with letters in its check digits", br.RedactCNPJ(cnpj[:12]+"AB", br.KeepCNPJEnds)},
		{"a CNPJ under the bare CPF redaction", br.RedactCNPJ(cnpj, br.KeepCPFMiddleBare)},
		{"a CPF read as a CNPJ", br.RedactCNPJ(cpf, br.KeepCNPJEnds)},
	} {
		if tc.got != "" {
			t.Errorf("%s: got %q, want \"\"", tc.name, tc.got)
		}
	}
}

func FuzzRedactCPF(f *testing.F) {
	for _, seed := range append(generatedCPFs(4), "", "12345678900", "123.456.789-00", "***456789**", "１２３") {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got := br.RedactCPF(s, br.KeepCPFMiddle)
		if !cpfShape.MatchString(s) {
			if got != "" {
				t.Fatalf("RedactCPF(%q) = %q for a value of another shape, want \"\"", s, got)
			}
			return
		}
		d := bareOf(s)
		if want := "***." + d[3:6] + "." + d[6:9] + "-**"; got != want {
			t.Fatalf("RedactCPF(%q) = %q, want %q", s, got, want)
		}
	})
}
