package br

import (
	"strings"

	"github.com/arandu-io/hesape/validation"
)

// The names of the rules this package declares, for a rule string:
//
//	rules := validation.Rules{
//		"document": "required|" + br.RuleCPFOrCNPJ,
//		"postcode": "required|" + br.RuleCEPFormat,
//		"phone":    "required|" + br.RuleMobile,
//	}
//
// They are the names of the PHP packages geekcom/validator-docs and
// laravellegends/pt-br-validator, so a rule set written for one of them reads
// the same here. Rules lists them with what each accepts and what it says
// when it refuses.
//
// A rule exists in the catalogue of hesape/validation once Install has run,
// and a rule set compiled before that fails at boot with an unknown rule.
// Compile the sets after the application calls Install, from the function
// that wires the application. The seven rules the previous releases
// registered when this package was imported -- RuleCPF, RuleCNPJ,
// RuleCPFOrCNPJ, RuleCEP, RuleMobile, RuleBRL and the spelling cpf_or_cnpj --
// still are, so a rule set compiled in a package-level variable that names
// them by their constants keeps working. That registration is deprecated:
// call Install, which declares those rules too, and compile the sets after
// it.
//
// Each rule accepts only a string -- a document sent as a JSON number has
// already lost its leading zeros -- and a blank value is left to required or
// nullable. A rule that takes no argument refuses every value when it is
// written with one: cpf:strict is a typo nothing else would report.
const (
	// RuleCPF is a CPF, as ParseCPF accepts it.
	RuleCPF = "cpf"
	// RuleCNPJ is a CNPJ, numeric or alphanumeric, as ParseCNPJ accepts it.
	RuleCNPJ = "cnpj"
	// RuleCPFOrCNPJ is a value either ParseCPF or ParseCNPJ accepts: the
	// document field of a form a person or a company fills in. The rule was
	// registered as cpf_or_cnpj, which is still registered as another name for
	// it.
	RuleCPFOrCNPJ = "cpf_ou_cnpj"
	// RuleCEP is a CEP, as ParseCEP accepts it.
	RuleCEP = "cep"
	// RuleMobile is a Brazilian mobile number with its area code, as
	// ParseMobile accepts it. It was registered as celular, which now means
	// the number without the area code: see RuleLocalMobile.
	RuleMobile = "celular_com_ddd"
	// RuleBRL is an amount in reais, as ParseBRL accepts it: the price or the
	// budget field of a form.
	RuleBRL = "brl"

	// RuleLocalMobile is a mobile number without the area code, as
	// ParseLocalMobile accepts it.
	RuleLocalMobile = "celular"
	// RuleMobileWithCountry is a mobile number with the country code, as
	// ParseMobileWithCountry accepts it.
	RuleMobileWithCountry = "celular_com_codigo"
	// RuleMobileWithCountryBare is a mobile number with the country code and
	// no parentheses or hyphen: +5514999999999.
	RuleMobileWithCountryBare = "celular_com_codigo_sem_mascara"
	// RuleLocalLandline is a landline number without the area code, as
	// ParseLocalLandline accepts it.
	RuleLocalLandline = "telefone"
	// RuleLandline is a landline number with its area code, as ParseLandline
	// accepts it.
	RuleLandline = "telefone_com_ddd"
	// RuleLandlineWithCountry is a landline number with the country code, as
	// ParseLandlineWithCountry accepts it.
	RuleLandlineWithCountry = "telefone_com_codigo"
	// RuleAreaCode is an area code, as ParseAreaCode accepts it.
	RuleAreaCode = "ddd"

	// RuleCPFFormat is a CPF written 000.000.000-00, whatever its digits.
	RuleCPFFormat = "formato_cpf"
	// RuleCNPJFormat is a CNPJ written 00.000.000/0000-00, whatever its
	// characters.
	RuleCNPJFormat = "formato_cnpj"
	// RuleCPFOrCNPJFormat is a value written as either of the two.
	RuleCPFOrCNPJFormat = "formato_cpf_ou_cnpj"
	// RuleCEPFormat is a CEP written 00000-000 or 00.000-000.
	RuleCEPFormat = "formato_cep"

	// RulePlate is a vehicle plate, as ParsePlate accepts it.
	RulePlate = "placa"
	// RulePlateFormat is a plate written with three letters, an optional
	// hyphen, a digit, a digit or a letter, and two digits.
	RulePlateFormat = "formato_placa_de_veiculo"
	// RuleRenavam is a RENAVAM, as ParseRenavam accepts it.
	RuleRenavam = "renavam"
	// RuleCNH is a CNH register number, as ParseCNH accepts it.
	RuleCNH = "cnh"
	// RuleVoterTitle is a título de eleitor, as ParseVoterTitle accepts it.
	RuleVoterTitle = "titulo_eleitor"
	// RuleNIS is a NIS, as ParseNIS accepts it.
	RuleNIS = "nis"
	// RulePIS is a PIS: the same number as the NIS, under the other name.
	RulePIS = "pis"
	// RuleNISFormat is a NIS written 000.00000.00-0.
	RuleNISFormat = "formato_nis"
	// RulePISFormat is a PIS written 000.00000.00-0.
	RulePISFormat = "formato_pis"
	// RuleCNS is a Cartão Nacional de Saúde number, as ParseCNS accepts it.
	RuleCNS = "cns"
	// RuleCivilCertificate is the number of a birth, marriage or death
	// certificate, as ParseCivilCertificate accepts it.
	RuleCivilCertificate = "certidao"
	// RuleCivilCertificateFormat is a certificate number written in its nine
	// groups.
	RuleCivilCertificateFormat = "formato_certidao"
	// RuleUF is the code of a federative unit, as ParseUF accepts it.
	RuleUF = "uf"
	// RulePassport is a Brazilian passport number, as ParsePassport accepts
	// it.
	RulePassport = "passaporte"
	// RuleStateRegistration is an inscrição estadual of the unit named as the
	// argument, inscricao_estadual:SP, as ParseStateRegistration accepts it.
	RuleStateRegistration = "inscricao_estadual"
	// RuleDate is a day written 31/12/2026, as ParseDate accepts it.
	RuleDate = "data"
	// RuleTimeOfDay is a time written 15:04 or 15:04:05, as ParseTimeOfDay
	// accepts it.
	RuleTimeOfDay = "horario"
	// RuleMoney is an amount in reais under the name the PHP packages' users
	// look for: another name for RuleBRL.
	RuleMoney = "moeda"
)

// Rule declares one rule of hesape/validation: its name, the check behind it,
// the sentence it says in each language, and what it takes as an argument.
//
// It is plain data and plain functions, so an application installs it through
// Install and a test calls Check without a validator. Rules returns them all.
type Rule struct {
	// Name is the name a rule string uses, such as "cpf" or "formato_cep".
	Name string
	// Aliases are other names the same rule answers to.
	Aliases []string
	// Check decides one string. parameters are the arguments written after
	// the colon, comma separated, in the order Parameters names them.
	Check func(value string, parameters []string) bool
	// Implicit marks a rule that runs on a blank value, the way required does.
	// None of these is: a blank value is left to required and nullable.
	Implicit bool
	// Parameters names the arguments the rule takes, in order. A rule with
	// none refuses every value when it is written with an argument, and a rule
	// with them refuses every value written with fewer or more.
	Parameters []string
	// Messages is the sentence the rule says when it refuses a value, by
	// locale. "pt-BR" is always there. A sentence holds :attribute and one
	// placeholder per name in Parameters, which Replace fills.
	Messages map[string]string
	// Replace fills the placeholders of a sentence that come from the
	// arguments. It is nil for a rule whose sentence has none.
	Replace func(message string, parameters []string) string
}

// Names lists the name of the rule and then its aliases.
func (r Rule) Names() []string { return append([]string{r.Name}, r.Aliases...) }

// Message is the sentence for locale, and the pt-BR one for a locale the rule
// has no sentence in.
func (r Rule) Message(locale string) string {
	if message, ok := r.Messages[locale]; ok {
		return message
	}
	return r.Messages["pt-BR"]
}

// Extension is the rule as hesape/validation calls it: a value that is not a
// string is refused, and so is a rule written with other arguments than it
// takes.
func (r Rule) Extension() validation.ExtensionFunc {
	return func(_ *validation.Validator, _ string, value any, parameters []string) bool {
		s, ok := value.(string)
		return ok && len(parameters) == len(r.Parameters) && r.Check(s, parameters)
	}
}

// Messages lists the sentence of every rule in locale, by the name a rule
// string writes -- aliases included -- which is the key under validation.
// An application that keeps its sentences in language files can load them
// from it, and the lines it defines itself take precedence over these.
func Messages(locale string) map[string]string {
	out := map[string]string{}
	for _, rule := range Rules() {
		for _, name := range rule.Names() {
			out[name] = rule.Message(locale)
		}
	}
	return out
}

// Install puts every rule of Rules into the catalogue of hesape/validation,
// under its name and its aliases, with the pt-BR sentence as the sentence the
// rule says, and gives the factory the fillers of the placeholders that come
// from arguments.
//
// It is called once, explicitly, where the application is wired -- from
// bootstrap/app.go, with the factory the application builds -- and before the
// rule sets that name these rules are compiled. A name already in the
// catalogue is left as it is, so calling Install a second time, or after the
// import-time registration of the seven earlier rules, registers nothing
// twice. A nil factory registers the rules in the catalogue without the
// fillers, so a rule that takes an argument says its placeholder as written.
//
// The sentences are the default: a line the application's translations hold
// under validation.<rule> is read before them.
func Install(f *validation.Factory) {
	for _, rule := range Rules() {
		for _, name := range rule.Names() {
			if registered(name) {
				continue
			}
			switch {
			case f != nil && rule.Implicit:
				f.ExtendImplicit(name, rule.Extension(), rule.Message("pt-BR"))
			case f != nil:
				f.Extend(name, rule.Extension(), rule.Message("pt-BR"))
			case rule.Implicit:
				validation.ExtendImplicit(name, rule.Extension(), rule.Message("pt-BR"))
			default:
				validation.Extend(name, rule.Extension(), rule.Message("pt-BR"))
			}
		}
		if f != nil && rule.Replace != nil {
			replace := rule.Replace
			filler := func(message, _, _ string, parameters []string, _ *validation.Validator) string {
				return replace(message, parameters)
			}
			for _, name := range rule.Names() {
				f.Replacer(name, filler)
			}
		}
	}
}

// registered reports whether the catalogue of hesape/validation already knows
// a rule name.
func registered(name string) bool {
	_, err := validation.Compile(validation.Rules{"field": name})
	return err == nil
}

// legacy returns the rule of that name from the table.
func legacy(name string) Rule {
	for _, rule := range Rules() {
		if rule.Name == name {
			return rule
		}
	}
	panic("br: no rule named " + name)
}

// legacyAs is the rule under another of its names.
func legacyAs(name, alias string) Rule {
	rule := legacy(name)
	rule.Name = alias
	return rule
}

// init registers, when the package is imported, the rules the earlier
// releases registered -- cpf, cnpj, cpf_or_cnpj, cep, celular and brl -- and
// the two names the constants RuleCPFOrCNPJ and RuleMobile now hold, and does
// nothing else. It exists so that a rule set compiled in a package-level
// variable, before main, keeps finding them; it is deprecated, and Install is
// what an application calls instead.
//
// celular is the rule without the area code now, as the PHP packages name it;
// the mobile number with an area code, which celular was, is celular_com_ddd.
func init() {
	validation.Extend(RuleCPF, legacy(RuleCPF).Extension(), legacy(RuleCPF).Message("pt-BR"))
	validation.Extend(RuleCNPJ, legacy(RuleCNPJ).Extension(), legacy(RuleCNPJ).Message("pt-BR"))
	validation.Extend(RuleCPFOrCNPJ, legacy(RuleCPFOrCNPJ).Extension(), legacy(RuleCPFOrCNPJ).Message("pt-BR"))
	validation.Extend("cpf_or_cnpj", legacyAs(RuleCPFOrCNPJ, "cpf_or_cnpj").Extension(), legacy(RuleCPFOrCNPJ).Message("pt-BR"))
	validation.Extend(RuleCEP, legacy(RuleCEP).Extension(), legacy(RuleCEP).Message("pt-BR"))
	validation.Extend(RuleLocalMobile, legacy(RuleLocalMobile).Extension(), legacy(RuleLocalMobile).Message("pt-BR"))
	validation.Extend(RuleMobile, legacy(RuleMobile).Extension(), legacy(RuleMobile).Message("pt-BR"))
	validation.Extend(RuleBRL, legacy(RuleBRL).Extension(), legacy(RuleBRL).Message("pt-BR"))
}

// plain turns a check on a string that takes no argument into the shape Rule
// holds.
func plain(valid func(string) bool) func(string, []string) bool {
	return func(s string, _ []string) bool { return valid(s) }
}

// ptBR is a Messages map with a pt-BR sentence only.
func ptBR(sentence string) map[string]string { return map[string]string{"pt-BR": sentence} }

// upperUnit is the argument of inscricao_estadual in capitals, and empty when
// the rule was written without one.
func upperUnit(parameters []string) string {
	if len(parameters) == 0 {
		return ""
	}
	return strings.ToUpper(parameters[0])
}
