---
name: br-package
description: Validate and format Brazilian values in an Arandu application with the Arandu BR package (Go, github.com/hyz-is/arandu-br). Use when the request is to "validate a CPF", "validate a CNPJ", "accept CPF or CNPJ", "the alphanumeric CNPJ", "validate a CEP", "format a CPF", "mask the CNPJ", "format money in reais", "R$ 1.234,56", "R$ 1234,56 without the dot", "R$ 47 without the centavos", "parse an amount in reais", "read the budget the user typed", "mask the CPF", "hide the CPF", "***.456.789-**", "format the date as dd/mm/yyyy", "list the states", "a select of UFs", "go get github.com/hyz-is/arandu-br", "unknown rule cpf", "rule cpf is already in the catalogue", "the CPF with spaces is refused", "validate a placa", "validate a telefone", "formato_cep", "celular_com_ddd", "inscricao_estadual", "validate a NIS, CNS, CNH, RENAVAM or titulo de eleitor", "br.Rules", "factory.Use", "br.Install", or when the application has its own NormalizeCPF, ValidCPF, FormatCNPJ, maskCPF, MaskedCPF, maskCNPJ, ParseBRL or BRL formatter, or when a project's go.mod already requires github.com/hyz-is/arandu-br. Covers the validation rules of the PHP packages, `factory.Use(br.Rules())`, why they are named by constant and how a package-level rule set names them, the Parse functions and what they accept, ParseBRL, storing the bare spelling, the formatters, BRL, FormatBRL and Date, RedactCPF and RedactCNPJ, the UF enum, and the symptoms that mean something is wired wrong.
license: MIT
metadata:
  audience: app
---

# Using Arandu BR

A library of Brazilian values: the CPF, the CNPJ (numeric and alphanumeric), the
CEP, the telephone numbers, the other documents with check digits, the vehicle
plate, amounts in reais, dates, and the twenty-seven federative units, and a CPF
or a CNPJ with part of it hidden. It has no route, no table and no policy, so
the only line it adds to `bootstrap/app.go` is `validators.Use(br.Rules())`: there
is no `New`, no line in `k.Register` and no `aru migrate`. Import it, install its
rules, and call it.

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
| a hand check before `ctx.Validate`, a `MustCompile` of a CPF, CNPJ, CEP, phone or plate regexp | a rule in the rule set, named by its constant |
| a local `ValidatePlate`, `ValidateCNH`, `ValidatePIS`, `ValidateRenavam`, a phone regexp | `br.ParsePlate`, `br.ParseCNH`, `br.ParseNIS`, `br.ParseRenavam`, `br.ParseLocalMobile` and the other `Parse` functions, and the rule of the same document |
| a BRL formatter, or `number.WithLocale("pt-BR", ...)` around a currency call | `br.BRL(cents)` |
| `fmt.Sprintf("R$ %d,%02d", ...)`, a price printed as `R$ 47` | `br.FormatBRL(cents, br.BRLFormat{WithoutGrouping: true})`, `br.FormatBRL(cents, br.BRLFormat{WithoutZeroCents: true})` |
| a `ParseBRL`, `parseBRLCents` or any reader of an amount typed in reais | `br.ParseBRL(s)`, and `br.RuleBRL` in the rule set |
| `maskCPF`, `MaskCPF`, `MaskedCPF`, `maskCNPJ` | `br.RedactCPF(cpf, br.KeepCPFMiddle)` and the other named redactions, or a literal such as `"•••.###.###-••"` |
| a map or a slice of the states | `br.UFValues()` and the `br.UF` type |

Read what the local version accepted before switching: see *What is accepted*
below.

## Validate a form

The rules are the names of the PHP packages `geekcom/validator-docs` and
`laravellegends/pt-br-validator`: `cpf`, `cnpj`, `cpf_ou_cnpj`, `cep`,
`celular`, `celular_com_ddd`, `celular_com_codigo`, `telefone`,
`telefone_com_ddd`, `ddd`, `placa`, `renavam`, `cnh`, `titulo_eleitor`, `nis`,
`pis`, `cns`, `certidao`, `uf`, `passaporte`, `inscricao_estadual:UF`, `data`,
`horario`, `brl`, and the `formato_cpf`, `formato_cnpj`, `formato_cep` family
that checks the spelling alone. `br.Defs()` lists them, as `validation.RuleDef`
values, with what each accepts and its pt-BR sentence, and `br.Rules()` is the
`validation.Provider` that declares them. The application hands it to the
Factory once, where it is wired:

```go
validators := validation.NewFactory(validation.FromTranslation(tr))
validators.Use(br.Rules())
```

A rule set in a package-level variable is compiled before the Factory exists, so
it takes `validation.WithModuleRules()`, naming the rules by constant, and the
application checks it at boot:

```go
var StoreCustomer = validation.MustCompile(validation.Rules{
	"document": "required|" + br.RuleCPFOrCNPJ,
	"postcode": []any{"required", br.RuleCEPFormat},
	"phone":    "required|" + br.RuleMobile,
}, validation.WithModuleRules())

err := validators.Check(StoreCustomer)
```

The string form works the same, `"required|cpf"`, and so does mixing these rules
with the standard ones.

**Write the constant, never the word.** `br.RuleCPFOrCNPJ` is checked by the
compiler; `"cpf_ou_cnpj"` is a string nothing checks, and `cpf_or_cnpj`,
`cpf_cnpj` and the other aliases are the same rule under older names.

**Symptom to recognise:** `validators.Check` fails at boot naming an unknown
rule `formato_cep`, `placa` or any other from the list. The Factory never called
`Use(br.Rules())`, or the rule set is compiled on another Factory. A panic at
start naming the same rule in a package-level variable means the set was
compiled without `validation.WithModuleRules()`. The rules this package
registered at import in the earlier releases -- `cpf`, `cnpj`, `cpf_or_cnpj`,
`cpf_ou_cnpj`, `cep`, `celular`, `celular_com_ddd` and `brl` -- still are,
deprecated, so a set that names only those compiles without it.

**Symptom to recognise:** `br.Install` is called. It is the same as `Use`, and
deprecated: replace it with `validators.Use(br.Rules())`.

**Symptom to recognise:** `celular` refuses a number with its area code. In this
release `celular` is the number without the area code, as in the PHP packages;
the one with it is `celular_com_ddd`, which is what `br.RuleMobile` holds.

**Symptom to recognise:** `validation: rule cpf is already in the catalogue`.
The application registers its own `cpf` rule with `validation.Extend`. Delete
it: this package's is the one.

Each rule takes no argument -- `cpf:strict` does not compile; only
`inscricao_estadual` takes one, the unit -- and only a string, so a document sent
as a JSON number fails. A blank value is not checked; write `required` or
`nullable` beside it.

The sentences are Portuguese with the attribute, `O campo documento não é um
CPF válido.` Put the application's own under `validation.cpf`,
`validation.formato_cep` and so on in its translations; `br.Messages("pt-BR")`
returns every default by rule name.

A CPF and a CNPJ are **text**: keep them as strings from the form to a `text`
column, or they lose their leading zeros, and the alphanumeric CNPJ is not a
number at all.

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
The other documents follow the same rule with their own punctuation -- the NIS
`000.00000.00-0`, the CNS `000 0000 0000 0000`, a landline `(00) 0000-0000` -- and
the README lists each `Parse` function with what it returns. The PHP packages
forgive every character that is not a digit; these do not.

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
