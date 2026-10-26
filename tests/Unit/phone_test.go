package unit_test

import (
	"errors"
	"testing"

	br "github.com/hyz-is/arandu-br"
)

// Every number here is invented: a valid area code and digits nobody chose.
// None is a real line.

type parseCase struct {
	in   string
	want string // the digits that come back; empty when the value is refused
}

// runParse holds one Parse function to a table: the value is accepted with
// those digits, or refused with err, and Valid agrees.
func runParse(t *testing.T, name string, parse func(string) (string, error), valid func(string) bool, err error, cases []parseCase) {
	t.Helper()
	for _, c := range cases {
		got, e := parse(c.in)
		switch {
		case c.want == "" && !errors.Is(e, err):
			t.Errorf("%s(%q) = %q, %v; want %v", name, c.in, got, e, err)
		case c.want != "" && (e != nil || got != c.want):
			t.Errorf("%s(%q) = %q, %v; want %q", name, c.in, got, e, c.want)
		}
		if valid(c.in) != (c.want != "") {
			t.Errorf("Valid for %s(%q) = %v, want %v", name, c.in, valid(c.in), c.want != "")
		}
	}
}

func TestParseAreaCode(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"11", "11"}, {"99", "99"}, {"85", "85"}, {"28", "28"}, {"79", "79"},
		{"10", ""}, {"20", ""}, {"23", ""}, {"36", ""}, {"50", ""}, {"52", ""}, {"72", ""}, {"76", ""}, {"78", ""}, {"90", ""},
		{"", ""}, {"1", ""}, {"011", ""}, {"(11)", ""}, {" 11", ""}, {"1a", ""}, {"٣١", ""},
	}
	runParse(t, "ParseAreaCode", br.ParseAreaCode, br.ValidAreaCode, br.ErrInvalidAreaCode, cases)

	// Sixty-seven codes between 11 and 99 are assigned.
	got := 0
	for a := 0; a < 100; a++ {
		if br.ValidAreaCode(string([]byte{byte('0' + a/10), byte('0' + a%10)})) {
			got++
		}
	}
	if got != 67 {
		t.Errorf("%d area codes are accepted, want the 67 Anatel assigns", got)
	}
}

func TestParseLocalMobile(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"98899-4444", "988994444"},
		{"988994444", "988994444"},
		{"99999-5555", "999995555"},
		{"91234-5678", "912345678"},
		{"90000-0000", ""}, // the nine and eight zeros a form fills
		{"900000000", ""},
		{"8899-4444", ""}, // eight digits: a landline's length
		{"9800-1936", ""}, // an eight-digit number that began with 9 before 2016
		{"900-1926", ""},
		{"88899-4444", ""}, // a mobile number starts with 9
		{"98899 4444", ""},
		{"98899-44444", ""},
		{" 98899-4444", ""},
		{"98899-4444 ", ""},
		{"9889944444", ""},
		{"(11) 98899-4444", ""}, // that has an area code
		{"", ""},
	}
	runParse(t, "ParseLocalMobile", br.ParseLocalMobile, br.ValidLocalMobile, br.ErrInvalidMobile, cases)
}

func TestParseMobileWithCountry(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"+55(99)98899-4444", "5599988994444"},
		{"+55 (99) 98899-4444", "5599988994444"},
		{"+55 (99)98899-4444", "5599988994444"},
		{"+55(99) 98899-4444", "5599988994444"},
		{"+55(31)99999-5555", "5531999995555"},
		{"+5514999999999", "5514999999999"},
		{"+55 14 999999999", "5514999999999"},
		{"+55 14999999999", "5514999999999"},
		{"+5514 999999999", "5514999999999"},
		{"+5599988994444", "5599988994444"},
		{"+1(31)99999-5555", ""}, // another country's code: its numbers do not follow this plan
		{"+5(99)800-1926", ""},
		{"+55(10)98899-4444", ""}, // no area code 10
		{"+55(99)88899-4444", ""}, // does not start with 9
		{"+55(99)98899-444", ""},
		{"+55(99)9889-94444", ""},
		{"+55(99)988994444", ""}, // the number is punctuated or the whole is not
		{"+55  (99)98899-4444", ""},
		{"+55 (99)  98899-4444", ""},
		{"+55 (99 98899-4444", ""},
		{"+55 99 98899-4444", ""}, // that is neither spelling
		{"5599988994444", ""},     // no plus
		{"55(99)98899-4444", ""},
		{"+55(99)90000-0000", ""},
		{"+55", ""},
		{"+55(", ""},
		{"", ""},
	}
	runParse(t, "ParseMobileWithCountry", br.ParseMobileWithCountry, br.ValidMobileWithCountry, br.ErrInvalidMobile, cases)
}

func TestParseLocalLandline(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"3598-4550", "35984550"},
		{"35984550", "35984550"},
		{"9999-5555", "99995555"}, // an eight-digit number that began with 9 reads as a telephone
		{"2000-0000", "20000000"},
		{"1598-4550", ""},
		{"0598-4550", ""},
		{"99800-1926", ""},
		{"3598 4550", ""},
		{"3598-45500", ""},
		{"598-4550", ""},
		{"(11) 3598-4550", ""},
		{"", ""},
	}
	runParse(t, "ParseLocalLandline", br.ParseLocalLandline, br.ValidLocalLandline, br.ErrInvalidLandline, cases)
}

func TestParseLandline(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"(99)3500-4444", "9935004444"},
		{"(99) 3500-4444", "9935004444"},
		{"(31)9999-5555", "3199995555"},
		{"(77)9999-3333", "7799993333"},
		{"9935004444", "9935004444"},
		{"(99)9-1926", ""},
		{"(10)3500-4444", ""}, // no area code 10
		{"(99)1500-4444", ""}, // the first digit is not 0 or 1
		{"(99)99800-1926", ""},
		{"(99)  3500-4444", ""},
		{"99 3500-4444", ""},
		{"+55(99)3500-4444", ""},
		{"3500-4444", ""},
		{"99350044444", ""},
		{"", ""},
	}
	runParse(t, "ParseLandline", br.ParseLandline, br.ValidLandline, br.ErrInvalidLandline, cases)
}

func TestParseLandlineWithCountry(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		{"+55(99)3500-4444", "559935004444"},
		{"+55 (99) 3500-4444", "559935004444"},
		{"+55 99 35004444", "559935004444"},
		{"+559935004444", "559935004444"},
		{"+66(31)9999-5555", ""},
		{"+5(99)9-1926", ""},
		{"+55(99)99800-1926", ""},
		{"+55(10)3500-4444", ""},
		{"(99)3500-4444", ""},
		{"", ""},
	}
	runParse(t, "ParseLandlineWithCountry", br.ParseLandlineWithCountry, br.ValidLandlineWithCountry, br.ErrInvalidLandline, cases)
}

func TestLandlineFormatsAndLeavesAnythingElseAlone(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"1130000000":     "(11) 3000-0000",
		"(11) 3000-0000": "(11) 3000-0000",
		"abc":            "abc",
		"11300000":       "11300000",
	} {
		if got := br.Landline(in); got != want {
			t.Errorf("Landline(%q) = %q, want %q", in, got, want)
		}
	}
}

func FuzzParsePhones(f *testing.F) {
	for _, seed := range []string{"", "(11) 91234-5678", "+55 (11) 3000-0000", "+5511912345678", "(11)3000-0000", "98899-4444", "\x00"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if d, err := br.ParseLocalMobile(s); err == nil && (len(d) != 9 || d[0] != '9') {
			t.Fatalf("ParseLocalMobile(%q) = %q", s, d)
		}
		if d, err := br.ParseMobileWithCountry(s); err == nil && (len(d) != 13 || d[:2] != "55") {
			t.Fatalf("ParseMobileWithCountry(%q) = %q", s, d)
		}
		if d, err := br.ParseLandline(s); err == nil && len(d) != 10 {
			t.Fatalf("ParseLandline(%q) = %q", s, d)
		}
		if d, err := br.ParseLandlineWithCountry(s); err == nil && (len(d) != 12 || d[:2] != "55") {
			t.Fatalf("ParseLandlineWithCountry(%q) = %q", s, d)
		}
		if d, err := br.ParseLocalLandline(s); err == nil && len(d) != 8 {
			t.Fatalf("ParseLocalLandline(%q) = %q", s, d)
		}
	})
}

func TestParseAnyPhone(t *testing.T) {
	t.Parallel()
	cases := []parseCase{
		// A mobile number, with or without the area code.
		{"(11) 91234-5678", "11912345678"},
		{"(11)91234-5678", "11912345678"},
		{"11912345678", "11912345678"},
		{"91234-5678", "912345678"},
		{"912345678", "912345678"},
		// A landline, with or without it.
		{"1134567890", "1134567890"},
		{"(11) 3456-7890", "1134567890"},
		{"(11)3456-7890", "1134567890"},
		{"3456-7890", "34567890"},
		{"34567890", "34567890"},
		{"(85) 2000-0001", "8520000001"},
		// An area code that Anatel does not assign.
		{"(10) 91234-5678", ""},
		{"(20) 3456-7890", ""},
		{"1091234567", ""},
		{"(00) 3456-7890", ""},
		// A length that is not 8, 9, 10 or 11.
		{"", ""},
		{"3456789", ""},
		{"1234567890123", ""},
		{"119123456789", ""},
		{"(11) 91234-56789", ""},
		// Ten digits whose third is 9, and eleven with no 9 after the area code.
		{"1191234567", ""},
		{"(11) 9123-4567", ""},
		{"11812345678", ""},
		{"(11) 81234-5678", ""},
		// A landline starts with 2 to 9 and a mobile number with 9.
		{"1456-7890", ""},
		{"0456-7890", ""},
		{"812345678", ""},
		// Repeated digits.
		{"22222222", ""},
		{"99999999", ""},
		{"999999999", ""},
		{"2222222222", ""},
		{"99999999999", ""},
		{"(11) 11111-1111", ""},
		// Letters, a plus sign, the country code and punctuation out of place.
		{"(11) 9123A-5678", ""},
		{"abcdefgh", ""},
		{"+55 (11) 91234-5678", ""},
		{"+5511912345678", ""},
		{"5511912345678", ""},
		{"11 91234-5678", ""},
		{"(11)  91234-5678", ""},
		{"(11) 9123-45678", ""},
		{"91234 5678", ""},
		{" 91234-5678", ""},
		{"91234-5678 ", ""},
		{"٩١٢٣٤٥٦٧٨", ""},
	}
	runParse(t, "ParseAnyPhone", br.ParseAnyPhone, br.ValidAnyPhone, br.ErrInvalidPhone, cases)
}

func TestAnyPhoneAcceptsWhatTheFourKindsAccept(t *testing.T) {
	t.Parallel()
	// Every value the dedicated parsers accept is a telephone number, except a
	// repetition and the ten digits whose third is 9, which the field refuses.
	for _, in := range []string{"(11) 91234-5678", "11912345678", "91234-5678", "(11) 3000-0000", "1130000000", "3000-0000", "9800-1936", "(31)9999-5555"} {
		dedicated := br.ValidMobile(in) || br.ValidLandline(in) || br.ValidLocalMobile(in) || br.ValidLocalLandline(in)
		if !dedicated {
			t.Fatalf("%q is not a number of any dedicated kind", in)
		}
	}
	// The dynamic field takes 9800-1936, an eight-digit number, as telefone does,
	// and refuses (31)9999-5555, a ten-digit one whose third digit is 9.
	if !br.ValidAnyPhone("9800-1936") || br.ValidAnyPhone("(31)9999-5555") {
		t.Error("the eight-digit and the ten-digit rules for a leading 9 are not as documented")
	}
}

func TestNormalizePhoneAndFormatPhone(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"(11) 91234-5678":    "11912345678",
		"+55 (11) 3456-7890": "551134567890",
		" 3456-7890 ":        "34567890",
		"":                   "",
		"abc":                "",
		"١٢٣":                "",
	} {
		if got := br.NormalizePhone(in); got != want {
			t.Errorf("NormalizePhone(%q) = %q, want %q", in, got, want)
		}
	}
	for in, want := range map[string]string{
		"34567890":       "3456-7890",
		"912345678":      "91234-5678",
		"1134567890":     "(11) 3456-7890",
		"11912345678":    "(11) 91234-5678",
		"":               "",
		"3456789":        "3456789",
		"551134567890":   "551134567890",
		"(11) 3456-7890": "(11) 3456-7890",
		"3456-7890":      "3456-7890",
		"1134A67890":     "1134A67890",
	} {
		if got := br.FormatPhone(in); got != want {
			t.Errorf("FormatPhone(%q) = %q, want %q", in, got, want)
		}
	}
	// What the field stores is what it shows: format, parse, the same digits.
	for _, in := range []string{"34567890", "912345678", "1134567890", "11912345678"} {
		if got, err := br.ParseAnyPhone(br.FormatPhone(in)); err != nil || got != in {
			t.Errorf("ParseAnyPhone(FormatPhone(%q)) = %q, %v", in, got, err)
		}
	}
}

func FuzzParseAnyPhone(f *testing.F) {
	for _, seed := range []string{"", "(11) 91234-5678", "1134567890", "34567890", "912345678", "(11)3456-7890", "1191234567", "\x00"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		digits, err := br.ParseAnyPhone(s)
		if err != nil {
			return
		}
		if n := len(digits); n < 8 || n > 11 {
			t.Fatalf("ParseAnyPhone(%q) = %q", s, digits)
		}
		if br.NormalizePhone(digits) != digits {
			t.Fatalf("ParseAnyPhone(%q) = %q, not digits", s, digits)
		}
		if again, err := br.ParseAnyPhone(br.FormatPhone(digits)); err != nil || again != digits {
			t.Fatalf("the punctuated %q of %q does not parse back: %q, %v", br.FormatPhone(digits), digits, again, err)
		}
	})
}
