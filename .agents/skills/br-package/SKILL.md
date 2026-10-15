---
name: br-package
description: Validate and format Brazilian values in an Arandu application with the Arandu BR package (Go, github.com/hyz-is/arandu-br). Use when the request is to "validate a CPF", "validate a CNPJ", "accept CPF or CNPJ", "the alphanumeric CNPJ", "validate a CEP", "format a CPF", "mask the CNPJ", "format money in reais", "R$ 1.234,56", "R$ 1234,56 without the dot", "R$ 47 without the centavos", "parse an amount in reais", "read the budget the user typed", "mask the CPF", "hide the CPF", "***.456.789-**", "format the date as dd/mm/yyyy", "list the states", "a select of UFs", "go get github.com/hyz-is/arandu-br", "unknown rule cpf", "rule cpf is already in the catalogue", "the CPF with spaces is refused", or when the application has its own NormalizeCPF, ValidCPF, FormatCNPJ, maskCPF, MaskedCPF, maskCNPJ, ParseBRL or BRL formatter, or when a project's go.mod already requires github.com/hyz-is/arandu-br. Covers the six validation rules and why they are named by constant, the Parse functions and what they accept, ParseBRL, storing the bare spelling, the formatters, BRL, FormatBRL and Date, RedactCPF and RedactCNPJ, the UF enum, and the symptoms that mean something is wired wrong.
license: MIT
metadata:
  audience: app
---

# Using Arandu BR

A library of Brazilian values: the CPF, the CNPJ (numeric and alphanumeric), the
CEP, the mobile number, amounts in reais, dates, and the twenty-seven federative
units, and a CPF or a CNPJ with part of it hidden. It has no
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
| `fmt.Sprintf("R$ %d,%02d", ...)`, a price printed as `R$ 47` | `br.FormatBRL(cents, br.BRLFormat{WithoutGrouping: true})`, `br.FormatBRL(cents, br.BRLFormat{WithoutZeroCents: true})` |
| a `ParseBRL`, `parseBRLCents` or any reader of an amount typed in reais | `br.ParseBRL(s)`, and `br.RuleBRL` in the rule set |
| `maskCPF`, `MaskCPF`, `MaskedCPF`, `maskCNPJ` | `br.RedactCPF(cpf, br.KeepCPFMiddle)` and the other named redactions, or a literal such as `"•••.###.###-••"` |
| a map or a slice of the states | `br.UFValues()` and the `br.UF` type |

Read what the local version accepted before switching: see *What is accepted*
below.

## Validate a form

The package registers `cpf`, `cnpj`, `cpf_or_cnpj`, `cep`, `celular` and `brl`
with `hesape/validation` when it is initialized (`br.RuleCPF`, `br.RuleCNPJ`,
`br.RuleCPFOrCNPJ`, `br.RuleCEP`, `br.RuleMobile`, `br.RuleBRL`). A rule set names them by constant,
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
`cpf`, `cnpj`, `cpf_or_cnpj`, `cep`, `celular` or `brl`. The rule set spells the word; spell the
constant.

**Symptom to recognise:** `validation: rule cpf is already in the catalogue`.
The application registers its own `cpf` rule with `validation.Extend`. Delete
it: this package's is the one.

Each rule takes no argument — `cpf:strict` refuses every value — and only a
string, so a document sent as a JSON number fails. A blank value is not
checked; write `required` or `nullable` beside it.

The default sentences are English: `is not a valid CPF`, `is not a valid
CNPJ`, `is not a valid CPF or CNPJ`, `is not a valid CEP`, `is not a valid
mobile number`, `is not a valid amount in reais`. Put the application's own
under `validation.cpf`, `validation.cnpj`, `validation.cpf_or_cnpj`,
`validation.cep`, `validation.celular` and `validation.brl` in its
translations.

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

## Take an amount in reais

```go
var StoreQuote = validation.MustCompile(validation.Rules{
	"budget": "required|" + br.RuleBRL,
})

cents, err := br.ParseBRL(in.String("budget")) // int64 centavos, or br.ErrInvalidBRL
```

Store the centavos as an `int64`, never a float. `ParseBRL` reads:

| typed | centavos |
| --- | --- |
| `R$ 1.234,56`, `R$1.234,56`, `1.234,56`, `1234,56`, with spaces around | `123456` |
| `1.500` — a dot before three digits groups thousands | `150000` |
| `1234.56`, `1.5` — a dot before one or two digits at the end is the decimal point, as an HTML number input submits | `123456`, `150` |
| `0` | `0` — a form that needs a positive amount checks `cents > 0` itself |

It refuses a sign, another currency, `r$`, a space or a dot anywhere else among
the digits (`1 234`, `1.2345`, `.0,01`), a third decimal (`1,234`) and an amount
past the largest `int64`.

**Symptom to recognise:** an `<input type="number">` without `step="0.01"`
sends `1.500` for one and a half, which `ParseBRL` reads as fifteen hundred.
Give a number input that takes centavos `step="0.01"`.

## Print a value

| call | prints |
| --- | --- |
| `br.CPF(stored)` | `000.000.000-00` |
| `br.CNPJ(stored)` | `00.000.000/0000-00` |
| `br.CEP(stored)` | `00000-000` |
| `br.BRL(cents)` | `R$ 1.234,56`, and `-R$ 1.234,56` for a negative amount |
| `br.FormatBRL(cents, br.BRLFormat{WithoutGrouping: true})` | `R$ 1234,56` |
| `br.FormatBRL(cents, br.BRLFormat{WithoutZeroCents: true})` | `R$ 47`, and `R$ 199,50` when there are centavos |
| `br.RedactCPF(stored, br.KeepCPFMiddle)` | `***.456.789-**` |
| `br.RedactCPF(stored, br.KeepCPFMiddleBare)` | `***456789**` |
| `br.RedactCPF(stored, br.KeepCPFCheckDigits)` | `***.***.***-00` |
| `br.RedactCPF(stored, br.KeepCPFLastFive)` | `***.***.789-00` |
| `br.RedactCNPJ(stored, br.KeepCNPJEnds)` | `12.***.***/****-35` |
| `br.Date(t.In(location))` | `31/12/2026` |

The three document formatters return what they cannot punctuate unchanged; they
do not validate. `BRL` and `FormatBRL` take the amount in centavos as an
`int64`, never a float, and read no process-wide locale. `Date` prints the day
in the location the time carries, so convert a UTC time to the reader's
location first, or a late evening prints as the next day.

A `Redaction` is spelled like the document: `#` shows the character in that
position and anything else is written as it stands, so
`br.RedactCPF(stored, "•••.###.###-••")` prints `•••.456.789-••`. It has the
punctuated or the bare length, and the output takes that spelling.
`RedactCPF` and `RedactCNPJ` read the value bare or punctuated, do not check
the check digits, and return `""` — never the value — for anything they cannot
read: write the application's own answer for that case beside the call. Do not
trim or re-punctuate a value into the redaction to make it readable; store the
bare spelling, and redact that.

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
