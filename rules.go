package br

import "github.com/arandu-io/hesape/validation"

// The names this package registers with hesape/validation, for a rule string:
//
//	var StoreCustomer = validation.MustCompile(validation.Rules{
//		"document": "required|" + br.RuleCPFOrCNPJ,
//		"postcode": "required|" + br.RuleCEP,
//	})
//
// Write the constant rather than the name it holds. A rule name is real only
// once it is in the catalogue, and the catalogue is filled when this package is
// initialized: a rule set that names the constant imports the package, and Go
// initializes an imported package before the importer's own variables, so the
// set is compiled after the rules exist. A rule set that spells "cpf" in a
// package that does not import this one may be compiled first, and then it
// fails at boot with an unknown rule.
//
// Each rule takes no argument; one written after a colon refuses every value,
// because an argument this package would ignore is a typo nothing reports. Each
// accepts what the matching Parse function accepts, and only a string: a CPF
// sent as a JSON number has already lost its leading zeros. A blank value is
// not checked -- write required, or nullable, beside the rule for that.
const (
	// RuleCPF is a CPF, as ParseCPF accepts it.
	RuleCPF = "cpf"
	// RuleCNPJ is a CNPJ, numeric or alphanumeric, as ParseCNPJ accepts it.
	RuleCNPJ = "cnpj"
	// RuleCPFOrCNPJ is a value either ParseCPF or ParseCNPJ accepts: the
	// document field of a form a person or a company fills in.
	RuleCPFOrCNPJ = "cpf_or_cnpj"
	// RuleCEP is a CEP, as ParseCEP accepts it.
	RuleCEP = "cep"
	// RuleMobile is a Brazilian mobile number, as ParseMobile accepts it.
	RuleMobile = "celular"
	// RuleBRL is an amount in reais, as ParseBRL accepts it: the price or the
	// budget field of a form.
	RuleBRL = "brl"
)

// init puts the rules into the validation catalogue, and does nothing
// else.
//
// The catalogue is where hesape/validation looks a rule name up when a rule
// set is compiled, and rule sets are compiled in package-level variables,
// before main runs. Registering from main, or from a function an application
// calls in bootstrap/app.go, would come after them. The import is therefore
// the registration, and the rule constants are what make a rule set import the
// package. The sentences are the ones a field is drawn with when the
// application's translations name none for validation.cpf and the others.
func init() {
	validation.Extend(RuleCPF, accepts(ValidCPF), "is not a valid CPF")
	validation.Extend(RuleCNPJ, accepts(ValidCNPJ), "is not a valid CNPJ")
	validation.Extend(RuleCPFOrCNPJ, accepts(validCPFOrCNPJ), "is not a valid CPF or CNPJ")
	validation.Extend(RuleCEP, accepts(ValidCEP), "is not a valid CEP")
	validation.Extend(RuleMobile, accepts(ValidMobile), "is not a valid mobile number")
	validation.Extend(RuleBRL, accepts(ValidBRL), "is not a valid amount in reais")
}

// accepts turns a check on a string into the shape the catalogue takes. A value
// that is not a string, or a rule written with an argument, is refused.
func accepts(valid func(string) bool) validation.ExtensionFunc {
	return func(_ *validation.Validator, _ string, value any, parameters []string) bool {
		s, ok := value.(string)
		return ok && len(parameters) == 0 && valid(s)
	}
}

// validCPFOrCNPJ reports whether s is either document.
func validCPFOrCNPJ(s string) bool { return ValidCPF(s) || ValidCNPJ(s) }
