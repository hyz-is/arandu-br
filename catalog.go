package br

import (
	"fmt"
	"strings"

	"github.com/arandu-io/hesape/validation"
)

// Defs lists the rules this package declares, each as the RuleDef that
// hesape/validation takes: its name, its aliases, its check, the arguments it
// takes and its pt-BR sentence. It returns a new slice every time; the
// application changes nothing by editing it. Rules is the module that declares
// them.
//
// Where the two PHP packages the names come from disagree on what a rule
// accepts, the rule here accepts what the Parse function of its document
// accepts, which is the narrow reading: the value bare or in its one
// conventional punctuation, and nothing else.
func Defs() []validation.RuleDef {
	return []validation.RuleDef{
		def(RuleCPF, nil, ValidCPF, "O campo :attribute não é um CPF válido."),
		def(RuleCNPJ, nil, ValidCNPJ, "O campo :attribute não é um CNPJ válido."),
		def(RuleCPFOrCNPJ, []string{"cpf_or_cnpj", "cpf_cnpj"}, validCPFOrCNPJ, "O campo :attribute não é um CPF ou CNPJ válido."),
		def(RuleCPFFormat, nil, maskedCPF, "O campo :attribute não possui o formato válido de CPF."),
		def(RuleCNPJFormat, nil, maskedCNPJ, "O campo :attribute não possui o formato válido de CNPJ."),
		def(RuleCPFOrCNPJFormat, []string{"formato_cpf_cnpj"}, maskedCPFOrCNPJ, "O campo :attribute não possui o formato válido de CPF ou CNPJ."),
		def(RuleCEP, nil, ValidCEP, "O campo :attribute não é um CEP válido."),
		def(RuleCEPFormat, nil, maskedCEP, "O campo :attribute não possui um formato válido de CEP."),

		def(RuleLocalMobile, nil, ValidLocalMobile, "O campo :attribute não é um celular válido."),
		def(RuleMobile, nil, ValidMobile, "O campo :attribute não é um celular com DDD válido."),
		def(RuleMobileWithCountry, nil, ValidMobileWithCountry, "O campo :attribute não é um celular com código do país e DDD válido."),
		def(RuleMobileWithCountryBare, nil, validUnpunctuatedMobileWithCountry, "O campo :attribute não é um celular válido. Exemplo de celular válido: +5514999999999."),
		def(RuleLocalLandline, nil, ValidLocalLandline, "O campo :attribute não é um telefone válido."),
		def(RuleLandline, nil, ValidLandline, "O campo :attribute não é um telefone com DDD válido."),
		def(RuleLandlineWithCountry, nil, ValidLandlineWithCountry, "O campo :attribute não é um telefone com código do país e DDD válido."),
		def(RuleAreaCode, nil, ValidAreaCode, "O campo :attribute não é um DDD válido."),

		def(RulePlate, nil, ValidPlate, "O campo :attribute não é uma placa de veículo válida."),
		def(RulePlateFormat, nil, platePatternShape, "O campo :attribute não possui um formato válido de placa de veículo."),
		def(RuleRenavam, nil, ValidRenavam, "O campo :attribute não é um RENAVAM válido."),
		def(RuleCNH, nil, ValidCNH, "O campo :attribute não é uma carteira nacional de habilitação válida."),
		def(RuleVoterTitle, nil, ValidVoterTitle, "O campo :attribute não é um título de eleitor válido."),
		def(RuleNIS, nil, ValidNIS, "O campo :attribute não é um NIS válido."),
		def(RulePIS, nil, ValidNIS, "O campo :attribute não é um PIS válido."),
		def(RuleNISFormat, nil, maskedNIS, "O campo :attribute não possui o formato válido de NIS."),
		def(RulePISFormat, nil, maskedNIS, "O campo :attribute não é um PIS com formato válido."),
		def(RuleCNS, nil, ValidCNS, "O campo :attribute não é um CNS válido."),
		def(RuleCivilCertificate, nil, ValidCivilCertificate, "O campo :attribute não é um número de certidão válido."),
		def(RuleCivilCertificateFormat, nil, maskedCivilCertificate, "O campo :attribute não possui o formato válido de certidão."),
		def(RuleUF, nil, func(s string) bool { return UF(s).Valid() }, "O campo :attribute não é uma UF válida."),
		def(RulePassport, nil, ValidPassport, "O campo :attribute não é um passaporte válido."),
		stateRegistrationDef(),

		def(RuleDate, nil, ValidDate, "O campo :attribute não é uma data válida. Use o formato dd/mm/aaaa."),
		def(RuleTimeOfDay, nil, ValidTimeOfDay, "O campo :attribute não é um horário válido. Use o formato hh:mm ou hh:mm:ss."),
		def(RuleBRL, []string{RuleMoney}, ValidBRL, "O campo :attribute não é um valor em reais válido."),
	}
}

// def declares a rule that takes no argument and checks one string. The
// registry refuses an argument when a rule set is compiled; the check refuses
// it too, for the rules the process-wide catalogue holds, which accept any.
func def(name string, aliases []string, valid func(string) bool, sentence string) validation.RuleDef {
	return validation.RuleDef{
		Name:    name,
		Aliases: aliases,
		Check: func(_ *validation.Validator, _ string, value any, parameters []string) bool {
			s, ok := value.(string)
			return ok && len(parameters) == 0 && valid(s)
		},
		Messages: map[string]string{"pt-BR": sentence},
	}
}

// stateRegistrationDef declares inscricao_estadual, whose one argument is the
// unit. A unit that is not one of the twenty-seven is refused when the rule set
// is compiled, and the sentence names it.
func stateRegistrationDef() validation.RuleDef {
	return validation.RuleDef{
		Name: RuleStateRegistration,
		Check: func(_ *validation.Validator, _ string, value any, parameters []string) bool {
			s, ok := value.(string)
			return ok && len(parameters) == 1 && ValidStateRegistration(UF(upperUnit(parameters)), s)
		},
		Params: validation.Params{Min: 1, Max: 1},
		Args: func(parameters []string) error {
			if !UF(upperUnit(parameters)).Valid() {
				return fmt.Errorf("%s takes the code of a federative unit, such as SP, and not %q", RuleStateRegistration, parameters[0])
			}
			return nil
		},
		Messages: map[string]string{"pt-BR": "O campo :attribute não é uma inscrição estadual válida para :uf."},
		Replace: func(message, _, _ string, parameters []string, _ *validation.Validator) string {
			return strings.ReplaceAll(message, ":uf", upperUnit(parameters))
		},
	}
}

// validCPFOrCNPJ reports whether s is either document.
func validCPFOrCNPJ(s string) bool { return ValidCPF(s) || ValidCNPJ(s) }
