package br

import "strings"

// Rules lists the rules this package declares, each with its name, its
// aliases, its check and its pt-BR sentence. It returns a new slice every
// time; the application changes nothing by editing it.
//
// Where the two PHP packages the names come from disagree on what a rule
// accepts, the rule here accepts what the Parse function of its document
// accepts, which is the narrow reading: the value bare or in its one
// conventional punctuation, and nothing else.
func Rules() []Rule {
	return []Rule{
		{Name: RuleCPF, Check: plain(ValidCPF),
			Messages: ptBR("O campo :attribute não é um CPF válido.")},
		{Name: RuleCNPJ, Check: plain(ValidCNPJ),
			Messages: ptBR("O campo :attribute não é um CNPJ válido.")},
		{Name: RuleCPFOrCNPJ, Aliases: []string{"cpf_or_cnpj", "cpf_cnpj"}, Check: plain(validCPFOrCNPJ),
			Messages: ptBR("O campo :attribute não é um CPF ou CNPJ válido.")},
		{Name: RuleCPFFormat, Check: plain(maskedCPF),
			Messages: ptBR("O campo :attribute não possui o formato válido de CPF.")},
		{Name: RuleCNPJFormat, Check: plain(maskedCNPJ),
			Messages: ptBR("O campo :attribute não possui o formato válido de CNPJ.")},
		{Name: RuleCPFOrCNPJFormat, Aliases: []string{"formato_cpf_cnpj"}, Check: plain(maskedCPFOrCNPJ),
			Messages: ptBR("O campo :attribute não possui o formato válido de CPF ou CNPJ.")},
		{Name: RuleCEP, Check: plain(ValidCEP),
			Messages: ptBR("O campo :attribute não é um CEP válido.")},
		{Name: RuleCEPFormat, Check: plain(maskedCEP),
			Messages: ptBR("O campo :attribute não possui um formato válido de CEP.")},

		{Name: RuleLocalMobile, Check: plain(ValidLocalMobile),
			Messages: ptBR("O campo :attribute não é um celular válido.")},
		{Name: RuleMobile, Check: plain(ValidMobile),
			Messages: ptBR("O campo :attribute não é um celular com DDD válido.")},
		{Name: RuleMobileWithCountry, Check: plain(ValidMobileWithCountry),
			Messages: ptBR("O campo :attribute não é um celular com código do país e DDD válido.")},
		{Name: RuleMobileWithCountryBare, Check: plain(validUnpunctuatedMobileWithCountry),
			Messages: ptBR("O campo :attribute não é um celular válido. Exemplo de celular válido: +5514999999999.")},
		{Name: RuleLocalLandline, Check: plain(ValidLocalLandline),
			Messages: ptBR("O campo :attribute não é um telefone válido.")},
		{Name: RuleLandline, Check: plain(ValidLandline),
			Messages: ptBR("O campo :attribute não é um telefone com DDD válido.")},
		{Name: RuleLandlineWithCountry, Check: plain(ValidLandlineWithCountry),
			Messages: ptBR("O campo :attribute não é um telefone com código do país e DDD válido.")},
		{Name: RuleAreaCode, Check: plain(ValidAreaCode),
			Messages: ptBR("O campo :attribute não é um DDD válido.")},

		{Name: RulePlate, Check: plain(ValidPlate),
			Messages: ptBR("O campo :attribute não é uma placa de veículo válida.")},
		{Name: RulePlateFormat, Check: plain(platePatternShape),
			Messages: ptBR("O campo :attribute não possui um formato válido de placa de veículo.")},
		{Name: RuleRenavam, Check: plain(ValidRenavam),
			Messages: ptBR("O campo :attribute não é um RENAVAM válido.")},
		{Name: RuleCNH, Check: plain(ValidCNH),
			Messages: ptBR("O campo :attribute não é uma carteira nacional de habilitação válida.")},
		{Name: RuleVoterTitle, Check: plain(ValidVoterTitle),
			Messages: ptBR("O campo :attribute não é um título de eleitor válido.")},
		{Name: RuleNIS, Check: plain(ValidNIS),
			Messages: ptBR("O campo :attribute não é um NIS válido.")},
		{Name: RulePIS, Check: plain(ValidNIS),
			Messages: ptBR("O campo :attribute não é um PIS válido.")},
		{Name: RuleNISFormat, Check: plain(maskedNIS),
			Messages: ptBR("O campo :attribute não possui o formato válido de NIS.")},
		{Name: RulePISFormat, Check: plain(maskedNIS),
			Messages: ptBR("O campo :attribute não é um PIS com formato válido.")},
		{Name: RuleCNS, Check: plain(ValidCNS),
			Messages: ptBR("O campo :attribute não é um CNS válido.")},
		{Name: RuleCivilCertificate, Check: plain(ValidCivilCertificate),
			Messages: ptBR("O campo :attribute não é um número de certidão válido.")},
		{Name: RuleCivilCertificateFormat, Check: plain(maskedCivilCertificate),
			Messages: ptBR("O campo :attribute não possui o formato válido de certidão.")},
		{Name: RuleUF, Check: plain(func(s string) bool { return UF(s).Valid() }),
			Messages: ptBR("O campo :attribute não é uma UF válida.")},
		{Name: RulePassport, Check: plain(ValidPassport),
			Messages: ptBR("O campo :attribute não é um passaporte válido.")},
		{Name: RuleStateRegistration, Parameters: []string{"uf"},
			Check: func(s string, parameters []string) bool {
				return ValidStateRegistration(UF(upperUnit(parameters)), s)
			},
			Messages: ptBR("O campo :attribute não é uma inscrição estadual válida para :uf."),
			Replace: func(message string, parameters []string) string {
				if unit := upperUnit(parameters); unit != "" {
					return strings.ReplaceAll(message, ":uf", unit)
				}
				return strings.ReplaceAll(message, ":uf", "a UF informada")
			}},

		{Name: RuleDate, Check: plain(ValidDate),
			Messages: ptBR("O campo :attribute não é uma data válida. Use o formato dd/mm/aaaa.")},
		{Name: RuleTimeOfDay, Check: plain(ValidTimeOfDay),
			Messages: ptBR("O campo :attribute não é um horário válido. Use o formato hh:mm ou hh:mm:ss.")},
		{Name: RuleBRL, Aliases: []string{RuleMoney}, Check: plain(ValidBRL),
			Messages: ptBR("O campo :attribute não é um valor em reais válido.")},
	}
}

// validCPFOrCNPJ reports whether s is either document.
func validCPFOrCNPJ(s string) bool { return ValidCPF(s) || ValidCNPJ(s) }
