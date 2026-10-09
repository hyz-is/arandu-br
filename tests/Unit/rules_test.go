package unit_test

import (
	"strings"
	"testing"

	"github.com/arandu-io/hesape/validation"
	br "github.com/hyz-is/arandu-br"
)

// documentRules is a rule set as an application writes one: in a package-level
// variable, compiled before any test runs, naming the rules by their constants.
// That it compiled at all is the first assertion here -- an unknown rule name
// panics in MustCompile, which would stop this package's tests before main.
var documentRules = validation.MustCompile(validation.Rules{
	"cpf":      br.RuleCPF,
	"cnpj":     br.RuleCNPJ,
	"document": "required|" + br.RuleCPFOrCNPJ,
	"postcode": "nullable|" + br.RuleCEP,
})

// failures runs the set over the data and returns the messages by field.
func failures(t *testing.T, set *validation.Set, data validation.Data) validation.Errors {
	t.Helper()
	v := validation.Make(data, set)
	v.Passes()
	return v.Errors()
}

func TestTheRulesPlugIntoAHesapeValidator(t *testing.T) {
	t.Parallel()

	cpf := generatedCPFs(1)[0]
	cnpj := generatedCNPJs(2)[1]

	passing := validation.Data{
		"cpf":      punctuatedCPF(cpf),
		"cnpj":     strings.ToLower(cnpj),
		"document": cnpj,
		"postcode": "70040-010",
	}
	if errs := failures(t, documentRules, passing); errs.Any() {
		t.Fatalf("valid documents were refused: %v", errs)
	}

	asCPF := validation.Data{"document": cpf}
	if errs := failures(t, documentRules, asCPF); errs.Any() {
		t.Fatalf("cpf_or_cnpj refused a CPF: %v", errs)
	}

	refused := validation.Data{
		"cpf":      cnpj,
		"cnpj":     cpf,
		"document": cpf[:10] + string('0'+(cpf[10]-'0'+1)%10),
		"postcode": "00000-000",
	}
	errs := failures(t, documentRules, refused)
	for field, want := range map[string]string{
		"cpf":      "is not a valid CPF",
		"cnpj":     "is not a valid CNPJ",
		"document": "is not a valid CPF or CNPJ",
		"postcode": "is not a valid CEP",
	} {
		if got := errs.Get(field); len(got) != 1 || !strings.Contains(got[0], want) {
			t.Errorf("%s: messages = %q, want one saying %q", field, got, want)
		}
	}
}

func TestTheRulesLeaveABlankValueToRequired(t *testing.T) {
	t.Parallel()

	errs := failures(t, documentRules, validation.Data{"cpf": "", "cnpj": "", "document": "", "postcode": nil})
	if got := errs.Keys(); len(got) != 1 || got[0] != "document" {
		t.Fatalf("fields refused = %v, want only document, which is required", got)
	}
}

func TestTheRulesRefuseAValueThatIsNotAString(t *testing.T) {
	t.Parallel()

	// A CPF that starts with zero cannot survive being a number; one that
	// does not would, and is refused all the same, so that whether a form
	// passes does not depend on which CPF somebody happened to type.
	cpf := generatedCPFs(40)
	var number int64
	for _, c := range cpf {
		if c[0] != '0' {
			for _, d := range c {
				number = number*10 + int64(d-'0')
			}
			break
		}
	}
	errs := failures(t, documentRules, validation.Data{"cpf": number, "document": cpf[0]})
	if got := errs.Get("cpf"); len(got) != 1 {
		t.Fatalf("a CPF sent as the number %d was not refused: %v", number, errs)
	}
}

func TestARuleWrittenWithAnArgumentRefusesEverything(t *testing.T) {
	t.Parallel()

	set, err := validation.Compile(validation.Rules{"cpf": br.RuleCPF + ":strict"})
	if err != nil {
		t.Fatalf("compiling: %v", err)
	}
	if errs := failures(t, set, validation.Data{"cpf": generatedCPFs(1)[0]}); !errs.Has("cpf") {
		t.Fatal("cpf:strict accepted a valid CPF; an argument the rule ignores is a typo nothing reports")
	}
}

func TestTheRuleNamesAreTheOnesHesapeKnows(t *testing.T) {
	t.Parallel()

	// The names as an application would spell them by hand, after this
	// package has been initialized: the constants are the spelling, and these
	// are what they hold.
	for _, name := range []string{"cpf", "cnpj", "cpf_or_cnpj", "cep"} {
		if _, err := validation.Compile(validation.Rules{"field": name}); err != nil {
			t.Errorf("rule %q is not in the catalogue: %v", name, err)
		}
	}
	if _, err := validation.Compile(validation.Rules{"field": "cpf_cnpj"}); err == nil {
		t.Error("an unregistered rule name compiled, so the check above proves nothing")
	}
}
