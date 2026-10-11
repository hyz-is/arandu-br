package unit_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

// Every number here is invented: a valid area code, a 9, and digits nobody
// chose. None is a real line.

func TestParseMobileAcceptsBothSpellingsAndReturnsTheDigits(t *testing.T) {
	for in, want := range map[string]string{
		"(11) 91234-5678": "11912345678",
		"11912345678":     "11912345678",
		"(85) 99000-1234": "85990001234",
	} {
		got, err := br.ParseMobile(in)
		if err != nil || got != want {
			t.Errorf("ParseMobile(%q) = %q, %v; want %q", in, got, err, want)
		}
		if !br.ValidMobile(in) {
			t.Errorf("ValidMobile(%q) = false", in)
		}
	}
}

func TestParseMobileRefuses(t *testing.T) {
	for _, in := range []string{
		"",
		"(10) 91234-5678", // no area code 10
		"(20) 91234-5678", // nor 20
		"(11) 81234-5678", // a mobile number starts with 9
		"(11) 90000-0000", // the nine and eight zeros a form fills
		"1191234567",      // ten digits: a landline's length
		"(11)91234-5678",  // neither spelling
		"+55 11 91234-5678",
		"(11) 9123A-5678",
	} {
		if _, err := br.ParseMobile(in); !errors.Is(err, br.ErrInvalidMobile) {
			t.Errorf("ParseMobile(%q) = %v, want ErrInvalidMobile", in, err)
		}
	}
}

func TestMobileFormatsAndLeavesAnythingElseAlone(t *testing.T) {
	for in, want := range map[string]string{
		"11912345678":     "(11) 91234-5678",
		"(11) 91234-5678": "(11) 91234-5678",
		"abc":             "abc",
	} {
		if got := br.Mobile(in); got != want {
			t.Errorf("Mobile(%q) = %q, want %q", in, got, want)
		}
	}
}

func FuzzParseMobile(f *testing.F) {
	for _, seed := range []string{"", "(11) 91234-5678", "11912345678", "(00) 00000-0000", "(11) 90000-0000", "１１"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, err := br.ParseMobile(s)
		if err != nil {
			return
		}
		if len(got) != 11 || strings.Trim(got, "0123456789") != "" || got[2] != '9' {
			t.Fatalf("ParseMobile(%q) = %q: not eleven digits with the nine", s, got)
		}
		if again, err := br.ParseMobile(got); err != nil || again != got {
			t.Fatalf("the digits ParseMobile returned for %q do not parse to themselves: %q, %v", s, again, err)
		}
		if again, err := br.ParseMobile(br.Mobile(got)); err != nil || again != got {
			t.Fatalf("the punctuated spelling of %q does not parse back: %q, %v", got, again, err)
		}
	})
}

var mobileRules = validation.MustCompile(validation.Rules{"mobile": "required|" + br.RuleMobile})

func TestTheMobileRulePlugsIntoAHesapeValidator(t *testing.T) {
	if errs := failures(t, mobileRules, validation.Data{"mobile": "(11) 91234-5678"}); errs.Any() {
		t.Fatalf("a valid mobile number was refused: %v", errs)
	}
	if errs := failures(t, mobileRules, validation.Data{"mobile": "(10) 91234-5678"}); !errs.Has("mobile") {
		t.Fatal("an invalid area code passed the rule")
	}
}
