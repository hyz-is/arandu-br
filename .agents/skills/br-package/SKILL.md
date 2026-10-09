---
name: br-package
description: Validate and format Brazilian values in an Arandu application with the Arandu BR package (Go, github.com/hyz-is/arandu-br). Use when the request is to "validate a CPF", "validate a CNPJ", "accept CPF or CNPJ", "the alphanumeric CNPJ", "validate a CEP", "format a CPF", "mask the CNPJ", "format money in reais", "R$ 1.234,56", "format the date as dd/mm/yyyy", "list the states", "a select of UFs", "go get github.com/hyz-is/arandu-br", "unknown rule cpf", "rule cpf is already in the catalogue", "the CPF with spaces is refused", or when the application has its own NormalizeCPF, ValidCPF, FormatCNPJ or BRL formatter, or when a project's go.mod already requires github.com/hyz-is/arandu-br. Covers the four validation rules and why they are named by constant, the Parse functions and what they accept, storing the bare spelling, the formatters, BRL and Date, the UF enum, and the symptoms that mean something is wired wrong.
license: MIT
metadata:
  audience: app
---

# Using Arandu BR

A library of Brazilian values: the CPF, the CNPJ (numeric and alphanumeric), the
CEP, amounts in reais, dates, and the twenty-seven federative units. It has no
route, no table and no policy, so nothing about it goes in `bootstrap/app.go`:
there is no `New`, no line in `k.Register` and no `aru migrate`. Import it and
call it.

## Install

```bash
go get github.com/hyz-is/arandu-br
```

```go
import br "github.com/hyz-is/arandu-br"
```

## If the application already does this itself

Replace it; do not wrap it. A second CPF check in the application is the one
that drifts, and `aru doctor` reports a reimplemented helper.

| the application has | it becomes |
| --- | --- |
| `NormalizeCPF`, `ValidCPF`, `FormatCPF` | `br.ParseCPF`, `br.ValidCPF`, `br.CPF` |
| the same for the CNPJ and the CEP | `br.ParseCNPJ`, `br.ValidCNPJ`, `br.CNPJ`; `br.ParseCEP`, `br.ValidCEP`, `br.CEP` |
| a hand check before `ctx.Validate` | a rule in the rule set, named by its constant |
| a BRL formatter, or `number.WithLocale("pt-BR", ...)` around a currency call | `br.BRL(cents)` |
| a map or a slice of the states | `br.UFValues()` and the `br.UF` type |

Read what the local version accepted before switching: see *What is accepted*
below.

## Validate a form

The package registers `cpf`, `cnpj`, `cpf_or_cnpj` and `cep` with
`hesape/validation` when it is initialized. A rule set names them by constant,
in the package-level variable where every rule set lives:

```go
var StoreCustomer = validation.MustCompile(validation.Rules{
	"document": "required|" + br.RuleCPFOrCNPJ,
	"postcode": "required|" + br.RuleCEP,
})
```

**Write the constant, never the word.** A rule set is compiled before `main`,
and a rule exists only once this package's `init` has run. Naming
`br.RuleCPFOrCNPJ` imports the package into the file that declares the set, and
Go runs an imported package's `init` before the importer's variables. The bare
word `"cpf_or_cnpj"`, in a package that does not import this one, may be
compiled first.

**Symptom to recognise:** the application panics at boot naming an unknown rule
`cpf`, `cnpj`, `cpf_or_cnpj` or `cep`. The rule set spells the word; spell the
constant.

**Symptom to recognise:** `validation: rule cpf is already in the catalogue`.
The application registers its own `cpf` rule with `validation.Extend`. Delete
it: this package's is the one.

Each rule takes no argument — `cpf:strict` refuses every value — and only a
string, so a document sent as a JSON number fails. A blank value is not
checked; write `required` or `nullable` beside it.

The default sentences are English: `is not a valid CPF`, `is not a valid
CNPJ`, `is not a valid CPF or CNPJ`, `is not a valid CEP`. Put the
application's own under `validation.cpf`, `validation.cnpj`,
`validation.cpf_or_cnpj` and `validation.cep` in its translations.

## Take the value in, and store it bare

```go
in, err := ctx.Validate(requests.StoreCustomer)
if err != nil {
	return err
}
document, err := br.ParseCPF(in.String("document"))
```

`ParseCPF`, `ParseCNPJ` and `ParseCEP` return the bare spelling — eleven
digits, fourteen characters with letters in capitals, eight digits — or
`br.ErrInvalidCPF`, `br.ErrInvalidCNPJ`, `br.ErrInvalidCEP`. Store that
spelling. A column holding some values punctuated and some bare cannot be
searched for either, and a unique index on it does not stop the same document
twice.

A CNPJ column has to hold letters: the alphanumeric CNPJ carries them in its
first twelve positions. A column or a check that assumed digits only refuses
CNPJs the Receita Federal issues.

## What is accepted

The value bare, or in its one conventional punctuation, and nothing else:

| document | bare | punctuated |
| --- | --- | --- |
| CPF | `00000000000` | `000.000.000-00` |
| CNPJ | `00000000000000` | `00.000.000/0000-00`, letters in the first twelve |
| CEP | `00000000` | `00000-000` |

Not whitespace around it, not another separator, not a separator out of place.

**Symptom to recognise:** a document pasted with a trailing space is refused.
That is deliberate. If the form should forgive it, trim in the request before
validating, where the application decides it, rather than widening the check
for everybody.

The CPF and the CNPJ are checked against their check digits; a repetition of
one character is refused. A CEP has no check digit, so only a code below
`01000-000` is refused for its value. None of the three says whether the
document was issued or the address exists — the package makes no call.

A seeder or a factory that needs a valid document generates one:
`br.CPFCheckDigits(base)` and `br.CNPJCheckDigits(base)` return the check
digits for a nine- or twelve-character base. Never seed a real person's number.

## Print a value

| call | prints |
| --- | --- |
| `br.CPF(stored)` | `000.000.000-00` |
| `br.CNPJ(stored)` | `00.000.000/0000-00` |
| `br.CEP(stored)` | `00000-000` |
| `br.BRL(cents)` | `R$ 1.234,56`, and `-R$ 1.234,56` for a negative amount |
| `br.Date(t.In(location))` | `31/12/2026` |

The three document formatters return what they cannot punctuate unchanged; they
do not validate. `BRL` takes the amount in centavos as an `int64`, never a
float, and reads no process-wide locale. `Date` prints the day in the location
the time carries, so convert a UTC time to the reader's location first, or a
late evening prints as the next day.

## The federative units

```go
type Address struct {
	State br.UF // a column: Value and Scan refuse a code outside the set
}

uf, err := br.ParseUF(in.String("state"))  // "SP", exactly, or br.ErrInvalidUF
enum.Options(br.UFValues()...)             // the options of a select: SP / São Paulo
validation.EnumOf(br.UFValues()...)        // the rule object
```

In a rule string, the cases come off the type rather than being written again:

```go
"state": "required|enum:" + strings.Join(enum.Names(br.UFValues()...), ","),
```

## Reporting a problem

Issues and pull requests go to `github.com/hyz-is/arandu-br`. A vulnerability —
a value accepted that is not what the documentation says is accepted — goes to
the private advisory form named in that repository's `SECURITY.md`, never into
an issue, and never with a real person's document in it.

MIT licensed. Copyright Paulo R. Lima.
