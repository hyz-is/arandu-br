// Package br validates and formats the Brazilian values an application takes
// from a form and prints on a page: the CPF, the CNPJ, the CEP, an amount in
// reais, a date, and the federative unit -- and a CPF or a CNPJ with part of
// it hidden.
//
// It is a library. It registers no route, owns no table and decides no
// authorization, so nothing about it goes in bootstrap/app.go: a caller
// imports it and calls it.
//
// # Taking a value in
//
// ParseCPF, ParseCNPJ and ParseCEP read a value written bare or in its one
// conventional punctuation, and return the bare spelling to store. ValidCPF,
// ValidCNPJ and ValidCEP are the same question answered with a bool. The CPF
// and the CNPJ are checked against the check digits the Receita Federal
// computes; the CNPJ may be alphanumeric.
//
// ParseMobile reads a mobile number the same way. ParseBRL reads an amount in
// reais as a person types it, R$ 1.234,56, or as an HTML number input submits
// it, 1234.56, and returns centavos.
//
// The same checks are rules of hesape/validation, registered when this
// package is initialized and named by RuleCPF, RuleCNPJ, RuleCPFOrCNPJ,
// RuleCEP, RuleMobile and RuleBRL:
//
//	var StoreCustomer = validation.MustCompile(validation.Rules{
//		"document": "required|" + br.RuleCPFOrCNPJ,
//		"postcode": "required|" + br.RuleCEP,
//	})
//
// # Printing a value
//
// CPF, CNPJ, CEP and Mobile punctuate a stored value; BRL writes centavos as
// reais, R$ 1.234,56, and FormatBRL writes them without the grouping, without
// zero centavos or without the symbol; BRLDecimal writes an exact decimal
// amount finer than a centavo, every digit kept; Date writes a day as
// 31/12/2026 and DateTime a day and a time as 31/12/2026 15:04. RedactCPF and
// RedactCNPJ write a document with part of it hidden, through a Redaction
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
//	mobile.go                -> the mobile number
//	shape.go                 -> the masks they share
//	redact.go                -> RedactCPF, RedactCNPJ and Redaction
//	rules.go                 -> the hesape/validation rules
//	money.go, date.go        -> BRL, FormatBRL, BRLDecimal, ParseBRL, Date and DateTime
//	uf.go                    -> the UF enum
package br
