package unit_test

import (
	"testing"
	"time"

	br "github.com/hyz-is/arandu-br"
)

func TestParsePlate(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		// One of each shape the plates have had.
		{"P15186", "P15186"},
		{"1234567", "1234567"},
		{"AA1234", "AA1234"},
		{"ABC1234", "ABC1234"},
		{"BEE4R22", "BEE4R22"},
		{"ABC12D3", "ABC12D3"},
		// The hyphen the plates of 1990 to 2018 are written with, and the
		// letters in either case.
		{"ABC-1234", "ABC1234"},
		{"abc-1234", "ABC1234"},
		{"aBc1234", "ABC1234"},
		{"FUM-0B05", "FUM0B05"},
		{"FUM-5L58", "FUM5L58"},
		{"bee4r22", "BEE4R22"},
		// Refused.
		{"X1234", ""},
		{"a2c-1234", ""},
		{"abc-12ed", ""},
		{"abc 1234", ""},
		{"ãBC1234", ""},
		{"ABCD1234", ""},
		{"ABC12345", ""},
		{"ab1234a", ""},
		{"ab123a4", ""},
		{"abc+1234", ""},
		{"AB-1234", ""}, // the hyphen belongs to the three-letter plate
		{"ABC--1234", ""},
		{"ABC-123", ""},
		{"12345A7", ""},
		{"", ""},
	}
	runParse(t, "ParsePlate", br.ParsePlate, br.ValidPlate, br.ErrInvalidPlate, cases)
}

func TestParsePassport(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"GB139485", "GB139485"},
		{"gb139485", "GB139485"},
		{"AB000000", "AB000000"},
		{"A3948", ""},
		{"GB13948", ""},
		{"GB1394850", ""},
		{"G1139485", ""},
		{"GBA39485", ""},
		{"GB 39485", ""},
		{"", ""},
	}
	runParse(t, "ParsePassport", br.ParsePassport, br.ValidPassport, br.ErrInvalidPassport, cases)
}

func TestParseDate(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]time.Time{
		"31/12/2026": time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		"01/01/0001": time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"29/02/2028": time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC),
		"29/02/2000": time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
		"30/04/2026": time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
	} {
		got, err := br.ParseDate(in)
		if err != nil || !got.Equal(want) {
			t.Errorf("ParseDate(%q) = %v, %v; want %v", in, got, err, want)
		}
		if !br.ValidDate(in) {
			t.Errorf("ValidDate(%q) = false", in)
		}
	}
	for _, in := range []string{
		"", "31122026", "2026-12-31", "1/1/2026", "01/01/26", "31/04/2026", "29/02/2027", "29/02/1900", "00/01/2026",
		"01/00/2026", "01/13/2026", "32/01/2026", "01/01/0000", " 31/12/2026", "31/12/2026 ", "31-12-2026", "31/12/20261",
		"3a/12/2026", "٣١/12/2026",
	} {
		if _, err := br.ParseDate(in); err != br.ErrInvalidDate {
			t.Errorf("ParseDate(%q) = %v, want ErrInvalidDate", in, err)
		}
		if br.ValidDate(in) {
			t.Errorf("ValidDate(%q) = true", in)
		}
	}

	// Date and ParseDate are one pair.
	day := time.Date(2026, 12, 31, 18, 30, 0, 0, time.UTC)
	if back, err := br.ParseDate(br.Date(day)); err != nil || back.Format("2006-01-02") != "2026-12-31" {
		t.Errorf("ParseDate(Date(%v)) = %v, %v", day, back, err)
	}
}

func TestParseTimeOfDay(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]time.Duration{
		"00:00":    0,
		"15:04":    15*time.Hour + 4*time.Minute,
		"23:59":    23*time.Hour + 59*time.Minute,
		"09:30:15": 9*time.Hour + 30*time.Minute + 15*time.Second,
		"23:59:59": 23*time.Hour + 59*time.Minute + 59*time.Second,
	} {
		got, err := br.ParseTimeOfDay(in)
		if err != nil || got != want {
			t.Errorf("ParseTimeOfDay(%q) = %v, %v; want %v", in, got, err, want)
		}
		if !br.ValidTimeOfDay(in) {
			t.Errorf("ValidTimeOfDay(%q) = false", in)
		}
	}
	for _, in := range []string{
		"", "24:00", "12:60", "12:30:60", "9:30", "09:3", "0930", "093015", "09.30", "09:30:1", "09:30:15:00", "12:30 PM", " 12:30",
		"12:30 ", "ab:cd",
	} {
		if _, err := br.ParseTimeOfDay(in); err != br.ErrInvalidTimeOfDay {
			t.Errorf("ParseTimeOfDay(%q) = %v, want ErrInvalidTimeOfDay", in, err)
		}
	}
}

func TestNormalizeKeepsTextAndLetters(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"111.222.333-44":     "11122233344",
		" 111 222 333 44 ":   "11122233344",
		"000.000.001-91":     "00000000191",
		"cpf: 111.222.333-4": "1112223334",
		"":                   "",
		"٣٣٣":                "",
	} {
		if got := br.NormalizeCPF(in); got != want {
			t.Errorf("NormalizeCPF(%q) = %q, want %q", in, got, want)
		}
	}
	for in, want := range map[string]string{
		"12.ABC.345/01DE-35":   "12ABC34501DE35",
		"12.abc.345/01de-35":   "12ABC34501DE35",
		" 00.000.000/0001-91 ": "00000000000191",
		"12abc34501de35":       "12ABC34501DE35",
		"é.ç/-":                "",
	} {
		if got := br.NormalizeCNPJ(in); got != want {
			t.Errorf("NormalizeCNPJ(%q) = %q, want %q", in, got, want)
		}
	}
}
