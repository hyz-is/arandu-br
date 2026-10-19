package unit_test

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

// The amounts applications read and wrote before this package had the
// variants, each written here as it was written there and kept as the oracle
// the package is measured against.

// legacyQuoteParse is how a quote form read the budget a visitor typed.
func legacyQuoteParse(s string) (int64, bool) {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "R$"))
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, false
	}
	whole, fraction, hasComma := strings.Cut(s, ",")
	whole = strings.ReplaceAll(whole, ".", "")
	if whole == "" || strings.Trim(whole, "0123456789") != "" || len(whole) > 12 {
		return 0, false
	}
	cents := int64(0)
	if hasComma {
		if len(fraction) == 0 || len(fraction) > 2 || strings.Trim(fraction, "0123456789") != "" {
			return 0, false
		}
		if len(fraction) == 1 {
			fraction += "0"
		}
		cents, _ = strconv.ParseInt(fraction, 10, 64)
	}
	reais, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, false
	}
	return reais*100 + cents, true
}

// legacyProposalParse is how a proposal form read the amount of an HTML number
// input, refusing zero as the form did.
func legacyProposalParse(raw string) (int64, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, fmt.Errorf("empty amount")
	}
	if strings.Contains(value, ",") && strings.Contains(value, ".") {
		value = strings.ReplaceAll(value, ".", "")
	}
	value = strings.ReplaceAll(value, ",", ".")
	parts := strings.Split(value, ".")
	if len(parts) > 2 || parts[0] == "" || strings.HasPrefix(parts[0], "-") || strings.HasPrefix(parts[0], "+") {
		return 0, fmt.Errorf("invalid amount")
	}
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return 0, fmt.Errorf("invalid amount")
	}
	fraction := int64(0)
	if len(parts) == 2 {
		if len(parts[1]) == 0 || len(parts[1]) > 2 {
			return 0, fmt.Errorf("invalid amount")
		}
		fracText := parts[1]
		if len(fracText) == 1 {
			fracText += "0"
		}
		fraction, err = strconv.ParseInt(fracText, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount")
		}
	}
	if whole > 92233720368547758 {
		return 0, fmt.Errorf("amount overflow")
	}
	cents := whole*100 + fraction
	if cents <= 0 {
		return 0, fmt.Errorf("amount must be positive")
	}
	return cents, nil
}

// legacyProposalFormat is how the proposal pages printed an amount.
func legacyProposalFormat(cents int64) string {
	return fmt.Sprintf("R$ %d,%02d", cents/100, cents%100)
}

// legacyPriceListFormat is how a pricing page printed a price.
func legacyPriceListFormat(centavos int64) string {
	reais, cents := centavos/100, centavos%100
	whole := fmt.Sprintf("%d", reais)
	for i := len(whole) - 3; i > 0; i -= 3 {
		whole = whole[:i] + "." + whole[i:]
	}
	if cents == 0 {
		return "R$ " + whole
	}
	return fmt.Sprintf("R$ %s,%02d", whole, cents)
}

// amounts is a spread of non-negative centavos: the small ones, the whole
// reais, the thousands boundaries, random ones, and the largest.
func amounts() []int64 {
	out := []int64{0, 1, 9, 10, 99, 100, 101, 4700, 19950, 99999, 100000, 123456, 1234500, 99999999, 100000000, math.MaxInt64}
	r := rand.New(rand.NewPCG(2026, 1015))
	for range 200 {
		out = append(out, r.Int64N(1_000_000_000_000), r.Int64N(1000)*100)
	}
	return out
}

func TestFormatBRLReproducesWhatApplicationsPrinted(t *testing.T) {
	t.Parallel()

	for _, cents := range amounts() {
		if got, want := br.FormatBRL(cents, br.BRLFormat{WithoutGrouping: true}), legacyProposalFormat(cents); got != want {
			t.Errorf("FormatBRL(%d, WithoutGrouping) = %q, the proposal pages printed %q", cents, got, want)
		}
		if got, want := br.FormatBRL(cents, br.BRLFormat{WithoutZeroCents: true}), legacyPriceListFormat(cents); got != want {
			t.Errorf("FormatBRL(%d, WithoutZeroCents) = %q, the pricing page printed %q", cents, got, want)
		}
	}

	for _, tc := range []struct {
		cents int64
		f     br.BRLFormat
		want  string
	}{
		{19950, br.BRLFormat{WithoutZeroCents: true}, "R$ 199,50"},
		{4700, br.BRLFormat{WithoutZeroCents: true}, "R$ 47"},
		{1234500, br.BRLFormat{WithoutZeroCents: true}, "R$ 12.345"},
		{0, br.BRLFormat{WithoutZeroCents: true}, "R$ 0"},
		{-4700, br.BRLFormat{WithoutZeroCents: true}, "-R$ 47"},
		{123456, br.BRLFormat{WithoutGrouping: true}, "R$ 1234,56"},
		{0, br.BRLFormat{WithoutGrouping: true}, "R$ 0,00"},
		{-123456, br.BRLFormat{WithoutGrouping: true}, "-R$ 1234,56"},
		{123400, br.BRLFormat{WithoutGrouping: true, WithoutZeroCents: true}, "R$ 1234"},
		{math.MinInt64, br.BRLFormat{WithoutGrouping: true}, "-R$ 92233720368547758,08"},
		{math.MaxInt64, br.BRLFormat{WithoutZeroCents: true}, "R$ 92.233.720.368.547.758,07"},
		{123456, br.BRLFormat{WithoutSymbol: true}, "1.234,56"},
		{-123456, br.BRLFormat{WithoutSymbol: true}, "-1.234,56"},
		{123456, br.BRLFormat{WithoutSymbol: true, WithoutGrouping: true}, "1234,56"},
		{5, br.BRLFormat{WithoutSymbol: true, WithoutGrouping: true}, "0,05"},
		{4700, br.BRLFormat{WithoutSymbol: true, WithoutZeroCents: true}, "47"},
	} {
		if got := br.FormatBRL(tc.cents, tc.f); got != tc.want {
			t.Errorf("FormatBRL(%d, %+v) = %q, want %q", tc.cents, tc.f, got, tc.want)
		}
	}
}

// TestParseBRLReadsWhatFormatBRLWritesWithoutTheSymbol holds the promise
// WithoutSymbol is for: a form field filled with the amount reads back as the
// amount.
func TestParseBRLReadsWhatFormatBRLWritesWithoutTheSymbol(t *testing.T) {
	t.Parallel()

	for _, cents := range amounts() {
		if cents < 0 {
			continue
		}
		for _, f := range []br.BRLFormat{{WithoutSymbol: true}, {WithoutSymbol: true, WithoutGrouping: true}} {
			written := br.FormatBRL(cents, f)
			if got, err := br.ParseBRL(written); err != nil || got != cents {
				t.Errorf("ParseBRL(%q) = %d, %v; FormatBRL wrote it from %d", written, got, err, cents)
			}
		}
	}
}

func TestBRLDecimal(t *testing.T) {
	t.Parallel()

	for in, want := range map[string]string{
		"1234.5678":                      "R$ 1.234,5678",
		"0.000001":                       "R$ 0,000001",
		"-0.000001":                      "-R$ 0,000001",
		"47":                             "R$ 47",
		"1234567.10":                     "R$ 1.234.567,10",
		"0.000000000000000001":           "R$ 0,000000000000000001",
		"123456789012345678901234567890": "R$ 123.456.789.012.345.678.901.234.567.890",
		"12.30":                          "R$ 12,30",
	} {
		if got := br.BRLDecimal(in); got != want {
			t.Errorf("BRLDecimal(%q) = %q, want %q", in, got, want)
		}
	}
	for _, in := range []string{"", "-", ".5", "1.", "1,5", "1.2.3", "R$ 1", " 1", "1e3", "--1", "+1", "１"} {
		if got := br.BRLDecimal(in); got != in {
			t.Errorf("BRLDecimal(%q) = %q, want it returned unchanged", in, got)
		}
	}
	for _, cents := range amounts() {
		decimal := fmt.Sprintf("%d.%02d", cents/100, cents%100)
		if cents < 0 {
			decimal = fmt.Sprintf("-%d.%02d", -(cents / 100), -(cents % 100))
		}
		if cents == math.MinInt64 {
			continue
		}
		if got, want := br.BRLDecimal(decimal), br.BRL(cents); got != want {
			t.Errorf("BRLDecimal(%q) = %q, BRL(%d) writes %q", decimal, got, cents, want)
		}
	}
}

func TestTheZeroBRLFormatIsBRL(t *testing.T) {
	t.Parallel()

	for _, cents := range append(amounts(), -1, -100, -123456, math.MinInt64) {
		if got, want := br.FormatBRL(cents, br.BRLFormat{}), br.BRL(cents); got != want {
			t.Errorf("FormatBRL(%d, BRLFormat{}) = %q, BRL writes %q", cents, got, want)
		}
	}
}

func TestParseBRLAccepts(t *testing.T) {
	t.Parallel()

	for in, want := range map[string]int64{
		"1.234,56":                    123456,
		"1234,56":                     123456,
		"R$ 1.234,56":                 123456,
		"R$1.234,56":                  123456,
		"R$\u00a01.234,56":            123456, // the no-break space a browser formats with
		"  R$  1.234,56  ":            123456,
		"\t1234\n":                    123400,
		"1234":                        123400,
		"R$ 80.000,00":                8000000,
		"R$ 50.000,00":                5000000,
		"50000":                       5000000,
		"1.500":                       150000,
		"1.234.567,89":                123456789,
		"1,5":                         150,
		"1,05":                        105,
		"0,01":                        1,
		"0":                           0,
		"0,00":                        0,
		"007":                         700,
		"1234.56":                     123456, // an HTML number input
		"1234.5":                      123450,
		"0.01":                        1,
		"1.5":                         150,
		"12.34":                       1234,
		"R$ 47":                       4700,
		"R$ 1234,56":                  123456,
		"92.233.720.368.547.758,07":   math.MaxInt64,
		"92233720368547758.07":        math.MaxInt64,
		"0000000000000000000000001,5": 150,
	} {
		got, err := br.ParseBRL(in)
		if err != nil || got != want {
			t.Errorf("ParseBRL(%q) = %d, %v; want %d", in, got, err, want)
		}
		if !br.ValidBRL(in) {
			t.Errorf("ValidBRL(%q) = false", in)
		}
	}
}

func TestParseBRLRefuses(t *testing.T) {
	t.Parallel()

	for _, in := range []string{
		"",
		"   ",
		"R$",
		"R$ ",
		"r$ 10",
		"BRL 10",
		"US$ 10",
		"R $10",
		"-10",
		"-R$ 10",
		"R$ -10",
		"+10",
		"1,",
		",50",
		".50",
		"1.",
		"1,234",       // a third digit of centavos
		"1,2,3",       // two commas
		"1,234.56",    // the comma grouping, the dot the decimals
		"1.234.56",    // grouping and a decimal dot together
		"1234.567",    // three digits after a lone dot of no grouping
		"12.34.567",   // a group of two
		"1234.567,00", // a first group of four
		"1.2345",      // a group of four
		"1..000",
		"1 234",
		"1 234,56",
		"1.234 ,56",
		"1e3",
		"0x10",
		"１２３", // full-width digits
		"12,3a",
		"1.-5",
		"1.+5",
		"92.233.720.368.547.758,08", // past the largest int64
		"99999999999999999999",
	} {
		got, err := br.ParseBRL(in)
		if !errors.Is(err, br.ErrInvalidBRL) || got != 0 {
			t.Errorf("ParseBRL(%q) = %d, %v; want ErrInvalidBRL", in, got, err)
		}
		if br.ValidBRL(in) {
			t.Errorf("ValidBRL(%q) = true", in)
		}
	}
}

// TestParseBRLWhereTheApplicationsDisagree names each spelling the two forms
// read differently, or one read and the other refused, and what the package
// reads it as. Every case here is a line of the upgrade guide.
func TestParseBRLWhereTheApplicationsDisagree(t *testing.T) {
	t.Parallel()

	const refused = -1
	for _, tc := range []struct {
		in                     string
		quote, proposal, wants int64
	}{
		// A dot followed by one or two digits: the quote form dropped the dot,
		// the proposal form read a decimal point, and so does the package.
		{"1.50", 15000, 150, 150},
		{"12.34", 123400, 1234, 1234},
		{"1.5", 1500, 150, 150},
		// A dot followed by three digits groups thousands, which the proposal
		// form refused as a third decimal.
		{"1.234", 123400, refused, 123400},
		{"1.000.000", 100000000, refused, 100000000},
		// The symbol, which the proposal form refused.
		{"R$ 10", 1000, refused, 1000},
		// Zero, which the proposal form refused for being no price.
		{"0", 0, refused, 0},
		// A dot anywhere, and a space inside the digits, which the quote form
		// dropped.
		{"1.2345", 1234500, refused, refused},
		{"1..000", 100000, refused, refused},
		{"1 234", 123400, refused, refused},
		// A sign after a decimal dot, which the proposal form read as part of
		// the centavos.
		{"1.-5", refused, 95, refused},
		// More than twelve digits of reais, which the quote form refused.
		{"1234567890123", refused, 123456789012300, 123456789012300},
	} {
		quote, ok := legacyQuoteParse(tc.in)
		if !ok {
			quote = refused
		}
		proposal, err := legacyProposalParse(tc.in)
		if err != nil {
			proposal = refused
		}
		if quote != tc.quote || proposal != tc.proposal {
			t.Errorf("%q: the forms read %d and %d, the table says %d and %d", tc.in, quote, proposal, tc.quote, tc.proposal)
		}
		got, err := br.ParseBRL(tc.in)
		if err != nil {
			got = refused
		}
		if got != tc.wants {
			t.Errorf("ParseBRL(%q) = %d, want %d", tc.in, got, tc.wants)
		}
	}
}

var amountRules = validation.MustCompile(validation.Rules{
	"budget": "required|" + br.RuleBRL,
	"price":  "nullable|" + br.RuleBRL,
})

func TestTheBRLRulePlugsIntoAHesapeValidator(t *testing.T) {
	t.Parallel()

	if errs := failures(t, amountRules, validation.Data{"budget": "R$ 1.234,56", "price": "1234.56"}); errs.Any() {
		t.Fatalf("valid amounts were refused: %v", errs)
	}
	if errs := failures(t, amountRules, validation.Data{"budget": "R$ 10", "price": ""}); errs.Any() {
		t.Fatalf("a blank amount beside nullable was refused: %v", errs)
	}
	errs := failures(t, amountRules, validation.Data{"budget": "1,234.56", "price": 1234})
	for _, field := range []string{"budget", "price"} {
		if got := errs.Get(field); len(got) != 1 || !strings.Contains(got[0], "is not a valid amount in reais") {
			t.Errorf("%s: messages = %q, want one saying it is not a valid amount in reais", field, got)
		}
	}
	set, err := validation.Compile(validation.Rules{"price": br.RuleBRL + ":2"})
	if err != nil {
		t.Fatalf("compiling: %v", err)
	}
	if errs := failures(t, set, validation.Data{"price": "10"}); !errs.Has("price") {
		t.Error("brl:2 accepted an amount; an argument the rule ignores is a typo nothing reports")
	}
}

// FuzzParseBRL holds the parser to three things on input nobody chose: what it
// accepts is not negative and comes back from every BRLFormat it can be
// printed in; and where the two forms it replaces both read an amount, and
// read the same one, it reads that amount or refuses the spelling -- it never
// reads another amount. A refusal there is a dot out of place, such as
// ".0,01", which both forms dropped; TestParseBRLAccepts holds the spellings
// both read and the package reads too.
func FuzzParseBRL(f *testing.F) {
	for _, seed := range []string{"", "R$ 1.234,56", "1234.56", "1.500", "1,5", "0", "1.-5", "92.233.720.368.547.758,07", "１２３"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, err := br.ParseBRL(s)
		if err == nil {
			if got < 0 {
				t.Fatalf("ParseBRL(%q) = %d, a negative amount", s, got)
			}
			for _, format := range []br.BRLFormat{{}, {WithoutGrouping: true}, {WithoutZeroCents: true}, {WithoutGrouping: true, WithoutZeroCents: true}} {
				printed := br.FormatBRL(got, format)
				if again, err := br.ParseBRL(printed); err != nil || again != got {
					t.Fatalf("ParseBRL(%q) = %d, printed as %q, reads back as %d, %v", s, got, printed, again, err)
				}
			}
		}
		quote, quoteOK := legacyQuoteParse(s)
		proposal, proposalErr := legacyProposalParse(s)
		if quoteOK && proposalErr == nil && quote == proposal && err == nil && got != quote {
			t.Fatalf("both forms read %q as %d; ParseBRL reads %d", s, quote, got)
		}
	})
}

func TestParseBRLReadsBackEveryAmountItPrints(t *testing.T) {
	t.Parallel()

	for _, cents := range amounts() {
		for _, format := range []br.BRLFormat{{}, {WithoutGrouping: true}, {WithoutZeroCents: true}, {WithoutGrouping: true, WithoutZeroCents: true}} {
			printed := br.FormatBRL(cents, format)
			if got, err := br.ParseBRL(printed); err != nil || got != cents {
				t.Errorf("ParseBRL(%q) = %d, %v; want %d", printed, got, err, cents)
			}
		}
	}
}
