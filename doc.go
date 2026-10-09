// Package br validates and formats the Brazilian values an application takes
// from a form and prints on a page: the CPF, the CNPJ, the CEP, an amount in
// reais, a date, and the federative unit.
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
// The same checks are rules of hesape/validation, registered when this
// package is initialized and named by RuleCPF, RuleCNPJ, RuleCPFOrCNPJ and
// RuleCEP:
//
//	var StoreCustomer = validation.MustCompile(validation.Rules{
//		"document": "required|" + br.RuleCPFOrCNPJ,
//		"postcode": "required|" + br.RuleCEP,
//	})
//
// # Printing a value
//
// CPF, CNPJ and CEP punctuate a stored value; BRL writes centavos as reais,
// R$ 1.234,56; Date writes a day as 31/12/2026. Each is a function of its
// argument alone.
//
// UF is the closed set of the twenty-seven units, in the shape hesape/enum
// reads: UFValues, Label for the name a form shows, and Value and Scan for a
// column.
//
// The files are laid out by subject:
//
//	cpf.go, cnpj.go, cep.go  -> parsing, check digits and punctuation
//	shape.go                 -> the masks the three share
//	rules.go                 -> the hesape/validation rules
//	money.go, date.go        -> BRL and Date
//	uf.go                    -> the UF enum
package br
