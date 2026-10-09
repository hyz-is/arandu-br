package unit_test

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/arandu-io/hesape/enum"
	"github.com/arandu-io/hesape/number"
	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

func TestBRL(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		cents int64
		want  string
	}{
		{0, "R$ 0,00"},
		{1, "R$ 0,01"},
		{10, "R$ 0,10"},
		{99, "R$ 0,99"},
		{100, "R$ 1,00"},
		{123456, "R$ 1.234,56"},
		{100000, "R$ 1.000,00"},
		{99999999, "R$ 999.999,99"},
		{100000000, "R$ 1.000.000,00"},
		{-1, "-R$ 0,01"},
		{-123456, "-R$ 1.234,56"},
		{math.MaxInt64, "R$ 92.233.720.368.547.758,07"},
		{math.MinInt64, "-R$ 92.233.720.368.547.758,08"},
	} {
		if got := br.BRL(tc.cents); got != tc.want {
			t.Errorf("BRL(%d) = %q, want %q", tc.cents, got, tc.want)
		}
	}
}

// hesapeBRL is what hesape/number writes for the amount with the process-wide
// locale set to pt-BR. It swaps that locale, so the tests calling it do not run
// in parallel: nothing else in this package formats through hesape/number, and
// the swap is undone before it returns.
func hesapeBRL(cents int64) string {
	return number.WithLocale("pt-BR", func() string { return number.CurrencyFromCents(cents, "BRL") })
}

// TestBRLAgreesWithHesapeUnderPtBR holds the two ways of writing reais in this
// collection to one answer. BRL exists because a library cannot set the
// process-wide locale; it must not exist as a second spelling of the amount.
func TestBRLAgreesWithHesapeUnderPtBR(t *testing.T) {
	for _, cents := range []int64{0, 1, -1, 5, 99, 100, 123456, -123456, 100000000, math.MaxInt64, math.MinInt64} {
		if got, want := br.BRL(cents), hesapeBRL(cents); got != want {
			t.Errorf("BRL(%d) = %q, hesape/number under pt-BR writes %q", cents, got, want)
		}
	}
}

func FuzzBRLAgreesWithHesapeUnderPtBR(f *testing.F) {
	for _, seed := range []int64{0, -1, 123456, math.MaxInt64, math.MinInt64} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, cents int64) {
		if got, want := br.BRL(cents), hesapeBRL(cents); got != want {
			t.Fatalf("BRL(%d) = %q, hesape/number under pt-BR writes %q", cents, got, want)
		}
	})
}

func TestDate(t *testing.T) {
	t.Parallel()

	brasilia := time.FixedZone("BRT", -3*60*60)
	lateEvening := time.Date(2026, time.December, 31, 22, 30, 0, 0, brasilia)
	for _, tc := range []struct {
		name string
		at   time.Time
		want string
	}{
		{"a day with one digit in a month with one digit", time.Date(2026, time.March, 5, 12, 0, 0, 0, time.UTC), "05/03/2026"},
		{"the last day of the year", time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC), "31/12/2026"},
		{"a leap day", time.Date(2028, time.February, 29, 0, 0, 0, 0, time.UTC), "29/02/2028"},
		{"in the location it holds", lateEvening, "31/12/2026"},
		{"the same instant in UTC is the next day", lateEvening.UTC(), "01/01/2027"},
		{"converted back", lateEvening.UTC().In(brasilia), "31/12/2026"},
	} {
		if got := br.Date(tc.at); got != tc.want {
			t.Errorf("%s: Date = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// fold lowers a name and drops the accents Portuguese names carry, which is
// how a person alphabetizes them.
func fold(s string) string {
	return strings.NewReplacer("á", "a", "â", "a", "ã", "a", "é", "e", "ê", "e", "í", "i",
		"ó", "o", "ô", "o", "õ", "o", "ú", "u", "ç", "c").Replace(strings.ToLower(s))
}

func TestUFValuesAreTheTwentySevenUnitsInTheOrderOfTheirNames(t *testing.T) {
	t.Parallel()

	values := br.UFValues()
	if len(values) != 27 {
		t.Fatalf("UFValues has %d units, want 27", len(values))
	}
	codes, labels := map[string]bool{}, map[string]bool{}
	for i, v := range values {
		if !v.Valid() || len(v.String()) != 2 || strings.ToUpper(v.String()) != v.String() {
			t.Errorf("%q is not a valid two-letter code", v)
		}
		if codes[v.String()] || labels[v.Label()] {
			t.Errorf("%q or its name %q appears twice", v, v.Label())
		}
		codes[v.String()], labels[v.Label()] = true, true
		if i > 0 && fold(values[i-1].Label()) >= fold(v.Label()) {
			t.Errorf("%s (%s) comes after %s (%s)", v, v.Label(), values[i-1], values[i-1].Label())
		}
	}
	for code, name := range map[br.UF]string{br.UFSP: "São Paulo", br.UFDF: "Distrito Federal", br.UFAP: "Amapá"} {
		if code.Label() != name {
			t.Errorf("%s.Label() = %q, want %q", code, code.Label(), name)
		}
	}
}

func TestParseUF(t *testing.T) {
	t.Parallel()

	for _, v := range br.UFValues() {
		if got, err := br.ParseUF(v.String()); err != nil || got != v {
			t.Errorf("ParseUF(%q) = %q, %v", v, got, err)
		}
	}
	for _, s := range []string{"", "sp", " SP", "SP ", "XX", "BR", "São Paulo"} {
		got, err := br.ParseUF(s)
		if !errors.Is(err, br.ErrInvalidUF) || got != "" || got.Valid() {
			t.Errorf("ParseUF(%q) = %q, %v; want ErrInvalidUF", s, got, err)
		}
	}
}

func TestUFValueAndScanRefuseACodeOutsideTheSet(t *testing.T) {
	t.Parallel()

	if got, err := br.UFRJ.Value(); err != nil || got != "RJ" {
		t.Errorf("UFRJ.Value() = %v, %v", got, err)
	}
	if _, err := br.UF("").Value(); !errors.Is(err, br.ErrInvalidUF) {
		t.Errorf("an empty UF was written: %v", err)
	}

	for _, src := range []any{"PE", []byte("PE")} {
		var v br.UF
		if err := v.Scan(src); err != nil || v != br.UFPE {
			t.Errorf("Scan(%#v) = %q, %v", src, v, err)
		}
	}
	for _, src := range []any{"pe", "XX", []byte(""), int64(26), nil} {
		v := br.UFPE
		if err := v.Scan(src); !errors.Is(err, br.ErrInvalidUF) || v != br.UFPE {
			t.Errorf("Scan(%#v) = %q, %v; want ErrInvalidUF and the value untouched", src, v, err)
		}
	}
}

func TestUFIsReadByHesapeEnumAndValidation(t *testing.T) {
	t.Parallel()

	options := enum.Options(br.UFValues()...)
	if len(options) != 27 || options[0] != (enum.Option{Value: "AC", Label: "Acre"}) {
		t.Fatalf("enum.Options = %v", options[:1])
	}
	if !enum.Holds(br.UFSP, "SC") || enum.Holds(br.UFSP, "XX") {
		t.Error("enum.Holds does not read the set off the type")
	}

	rule := validation.EnumOf(br.UFValues()...)
	for value, want := range map[any]bool{"SP": true, br.UFBA: true, "sp": false, br.UF("XX"): false, "": false} {
		if got := rule.Passes("uf", value); got != want {
			t.Errorf("EnumOf(UFValues).Passes(%#v) = %v, want %v", value, got, want)
		}
	}

	set := validation.MustCompile(validation.Rules{"uf": "required|enum:" + strings.Join(enum.Names(br.UFValues()...), ",")})
	if errs := failures(t, set, validation.Data{"uf": "MG"}); errs.Any() {
		t.Errorf("enum over UFValues refused MG: %v", errs)
	}
	if errs := failures(t, set, validation.Data{"uf": "XX"}); !errs.Has("uf") {
		t.Error("enum over UFValues accepted XX")
	}
}
