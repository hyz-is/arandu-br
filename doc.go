// Package br validates and formats the Brazilian values an application takes
// from a form and prints on a page: the CPF, the CNPJ, the CEP, the telephone
// numbers, the documents with check digits, an amount in reais, a date, and the
// federative unit -- and a CPF or a CNPJ with part of it hidden.
//
// It is a library. It registers no route, owns no table and decides no
// authorization, so nothing about it goes in bootstrap/app.go but one call to
// Install: a caller imports it and calls it.
//
// # Taking a value in
//
// ParseCPF, ParseCNPJ and ParseCEP read a value written bare or in its one
// conventional punctuation, and return the bare spelling to store. ValidCPF,
// ValidCNPJ and ValidCEP are the same question answered with a bool. The CPF
// and the CNPJ are checked against the check digits the Receita Federal
// computes; the CNPJ may be alphanumeric. A CPF and a CNPJ are text from the
// form to the column: a document read as a number has lost its leading zeros,
// and the alphanumeric CNPJ cannot be a number at all.
//
// ParseMobile reads a mobile number the same way, and ParseLocalMobile,
// ParseMobileWithCountry, ParseLandline, ParseLocalLandline and
// ParseLandlineWithCountry the other spellings of a telephone number.
// ParseNIS, ParseCNS, ParseCNH, ParseRenavam, ParseVoterTitle,
// ParseCivilCertificate, ParsePlate, ParsePassport and ParseStateRegistration
// read the other documents, and ParseDate and ParseTimeOfDay read 31/12/2026
// and 15:04. NormalizeCPF and NormalizeCNPJ drop the punctuation of text that
// already left the form. ParseBRL reads an amount in reais as a person types
// it, R$ 1.234,56, or as an HTML number input submits it, 1234.56, and returns
// centavos.
//
// # The rules
//
// The same checks are rules of hesape/validation, under the names of the PHP
// packages geekcom/validator-docs and laravellegends/pt-br-validator: cpf,
// cnpj, cpf_ou_cnpj, cep, celular, celular_com_ddd, telefone, placa, renavam,
// cnh, titulo_eleitor, nis, cns, certidao, uf, inscricao_estadual:SP, data,
// horario, brl and the formato_ rules, which check a spelling alone. Rules
// lists them as plain values -- the name, the aliases, the check, the arguments
// and the pt-BR sentence with :attribute -- and Install puts them in the
// catalogue, once, where the application is wired:
//
//	br.Install(factory)
//
// The rule sets that name them are compiled after that call. Each rule has a
// constant:
//
//	rules := validation.Rules{
//		"document": "required|" + br.RuleCPFOrCNPJ,
//		"postcode": "required|" + br.RuleCEPFormat,
//		"state":    "required|" + br.RuleStateRegistration + ":SP",
//	}
//
// Importing this package still registers cpf, cnpj, cpf_or_cnpj, cpf_ou_cnpj,
// cep, celular, celular_com_ddd and brl, so that a rule set compiled in a
// package-level variable keeps finding them. That registration is deprecated
// and goes in a later release; Install registers the same names.
//
// # Printing a value
//
// CPF, CNPJ, CEP, Mobile, Landline, NIS and CNS punctuate a stored value; BRL
// writes centavos as reais, R$ 1.234,56, and FormatBRL writes them without the
// grouping, without zero centavos or without the symbol; BRLDecimal writes an
// exact decimal amount finer than a centavo, every digit kept; Date writes a day
// as 31/12/2026 and DateTime a day and a time as 31/12/2026 15:04. RedactCPF
// and RedactCNPJ write a document with part of it hidden, through a Redaction
// such as KeepCPFMiddle, ***.456.789-**. Each is a function of its arguments
// alone.
//
// UF is the closed set of the twenty-seven units, in the shape hesape/enum
// reads: UFValues, Label for the name a form shows, and Value and Scan for a
// column.
//
// The files are laid out by subject:
//
//	cpf.go, cnpj.go, cep.go  -> parsing, check digits and punctuation
//	mobile.go, phone.go      -> the telephone numbers and the area code
//	nis.go, cns.go, cnh.go,
//	renavam.go, voter.go,
//	certificate.go           -> the other documents with check digits
//	plate.go                 -> the vehicle plate and the passport
//	stateregistration.go     -> the inscricao estadual of the 27 units
//	normalize.go             -> NormalizeCPF and NormalizeCNPJ
//	shape.go, masked.go      -> the masks they share, and the formato_ checks
//	redact.go                -> RedactCPF, RedactCNPJ and Redaction
//	rules.go, catalog.go     -> Rule, Rules, Install and the rule names
//	money.go, date.go        -> BRL, FormatBRL, BRLDecimal, ParseBRL, Date, DateTime,
//	                            ParseDate and ParseTimeOfDay
//	uf.go                    -> the UF enum
package br
