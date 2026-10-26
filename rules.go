package br

import (
	"strings"

	"github.com/arandu-io/hesape/validation"
)

// The names of the rules this package declares, for a rule string:
//
//	rules := validation.Rules{
//		"document": "required|" + br.RuleCPFOrCNPJ,
//		"postcode": []any{"required", br.RuleCEPFormat},
//		"phone":    "required|" + br.RuleMobile,
//	}
//
// They are the names of the PHP packages geekcom/validator-docs and
// laravellegends/pt-br-validator, so a rule set written for one of them reads
// the same here. Defs lists them with what each accepts and what it says when
// it refuses, and Rules is the module that declares them to a Factory:
//
//	factory := validation.NewFactory(validation.FromTranslation(tr)).Use(br.Rules())
//
// A rule set compiled in a package-level variable may name them if it is
// compiled with validation.WithModuleRules: the Factory resolves it when it is
// used, and Factory.Check fails at boot naming any rule no module declared. The names the earlier releases registered when this
// package was imported -- cpf, cnpj, cpf_or_cnpj, cep, celular and brl, with
// cpf_ou_cnpj and celular_com_ddd beside them -- are still registered that way,
// which is deprecated, and Rules leaves them to that registration.
//
// Each rule accepts only a string -- a document sent as a JSON number has
// already lost its leading zeros -- and a blank value is left to required or
// nullable. A rule that takes no argument does not compile when it is written
// with one: cpf:strict is a typo, and a boot failure reports it.
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
	// RuleAnyPhone is a telephone number, mobile or landline, with or without
	// the area code, as ParseAnyPhone accepts it: the field whose mask changes
	// as the person types.
	RuleAnyPhone = "telefone_ou_celular"
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

// Rules is the module of the rules in Defs, for validation.Factory.Use:
//
//	factory := validation.NewFactory(validation.FromTranslation(tr))
//	factory.Use(br.Rules())
//
// It declares every rule and alias of Defs in the factory's registry, with the
// pt-BR sentence and the filler of the :uf placeholder, and nothing else:
// nothing is registered before the application calls Use. A name already in
// the process-wide catalogue, which is what the deprecated registration at
// import leaves, or already declared in the registry, is skipped, so Use may be
// called on a factory that Install already served.
func Rules() validation.Provider { return provider{} }

type provider struct{}

// RegisterValidation declares the rules in the registry.
func (provider) RegisterValidation(registry *validation.Registry) {
	for _, def := range Defs() {
		var names []string
		for _, name := range append([]string{def.Name}, def.Aliases...) {
			if !registered(name) && !registry.Has(name) {
				names = append(names, name)
			}
		}
		if len(names) == 0 {
			continue
		}
		def.Name, def.Aliases = names[0], names[1:]
		registry.Add(def)
	}
}

// Install declares the rules of Rules on a Factory, and does nothing for a nil
// one.
//
// Deprecated: use factory.Use(br.Rules()), which is the same declaration through
// the call hesape/validation provides for a module.
func Install(f *validation.Factory) {
	if f != nil {
		f.Use(Rules())
	}
}

// Messages lists the sentence of every rule in locale, by the name a rule
// string writes -- aliases included -- which is the key under validation.
// An application that keeps its sentences in language files can load them
// from it, and the lines it defines itself take precedence over these.
func Messages(locale string) map[string]string {
	out := map[string]string{}
	for _, def := range Defs() {
		sentence, ok := def.Messages[locale]
		if !ok {
			sentence = def.Messages["pt-BR"]
		}
		for _, name := range append([]string{def.Name}, def.Aliases...) {
			out[name] = sentence
		}
	}
	return out
}

// registered reports whether the process-wide catalogue of hesape/validation
// already knows a rule name.
func registered(name string) bool {
	_, err := validation.Compile(validation.Rules{"field": name})
	return err == nil
}

// legacy returns the declaration of the rule of that name from Defs.
func legacy(name string) validation.RuleDef {
	for _, def := range Defs() {
		if def.Name == name {
			return def
		}
	}
	panic("br: no rule named " + name)
}

// legacyAs is the rule under another of its names.
func legacyAs(name, alias string) validation.RuleDef {
	def := legacy(name)
	def.Name = alias
	return def
}

// init registers, when the package is imported, the rules the earlier
// releases registered -- cpf, cnpj, cpf_or_cnpj, cep, celular and brl -- and
// the two names the constants RuleCPFOrCNPJ and RuleMobile now hold, and does
// nothing else. It exists so that a rule set compiled in a package-level
// variable of an application that does not call Factory.Use keeps finding
// them; it is deprecated, and factory.Use(br.Rules()) is what an application
// calls instead.
//
// celular is the rule without the area code now, as the PHP packages name it;
// the mobile number with an area code, which celular was, is celular_com_ddd.
func init() {
	validation.Extend(RuleCPF, legacy(RuleCPF).Check, legacy(RuleCPF).Messages["pt-BR"])
	validation.Extend(RuleCNPJ, legacy(RuleCNPJ).Check, legacy(RuleCNPJ).Messages["pt-BR"])
	validation.Extend(RuleCPFOrCNPJ, legacy(RuleCPFOrCNPJ).Check, legacy(RuleCPFOrCNPJ).Messages["pt-BR"])
	validation.Extend("cpf_or_cnpj", legacyAs(RuleCPFOrCNPJ, "cpf_or_cnpj").Check, legacy(RuleCPFOrCNPJ).Messages["pt-BR"])
	validation.Extend(RuleCEP, legacy(RuleCEP).Check, legacy(RuleCEP).Messages["pt-BR"])
	validation.Extend(RuleLocalMobile, legacy(RuleLocalMobile).Check, legacy(RuleLocalMobile).Messages["pt-BR"])
	validation.Extend(RuleMobile, legacy(RuleMobile).Check, legacy(RuleMobile).Messages["pt-BR"])
	validation.Extend(RuleBRL, legacy(RuleBRL).Check, legacy(RuleBRL).Messages["pt-BR"])
}

// upperUnit is the argument of inscricao_estadual in capitals, and empty when
// the rule was written without one.
func upperUnit(parameters []string) string {
	if len(parameters) == 0 {
		return ""
	}
	return strings.ToUpper(parameters[0])
}
