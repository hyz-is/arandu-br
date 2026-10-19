<p align="center">
  <img src=".github/logo.svg" alt="Arandu" width="180">
</p>

# Arandu BR

The Brazilian values an Arandu application takes from a form and prints on a
page: the CPF, the CNPJ (numeric and alphanumeric), the CEP, the mobile number, an amount in reais,
a date, and the federative unit — and a CPF or a CNPJ with part of it hidden.

It is a library. It registers no route, owns no table and decides no
authorization, so there is nothing to wire in `bootstrap/app.go`, nothing to
construct, and no `aru migrate` to run. An application imports it and calls it.

## Install

```bash
go get github.com/hyz-is/arandu-br
```

```go
import br "github.com/hyz-is/arandu-br"
```

The module requires Hesape and nothing else. `arandu.mod.toml` declares no
capability: no network, no filesystem, no process, no migration.

## Validate a form

The package registers six rules with `hesape/validation` when it is
initialized. Name them by their constants:

```go
var StoreCustomer = validation.MustCompile(validation.Rules{
	"cpf":      "required|" + br.RuleCPF,
	"cnpj":     "nullable|" + br.RuleCNPJ,
	"document": "required|" + br.RuleCPFOrCNPJ,
	"postcode": "required|" + br.RuleCEP,
})
```

| constant | rule | accepts |
| --- | --- | --- |
| `br.RuleCPF` | `cpf` | a CPF, `00000000000` or `000.000.000-00` |
| `br.RuleCNPJ` | `cnpj` | a CNPJ, `00000000000000` or `00.000.000/0000-00`, letters allowed in the first twelve positions |
| `br.RuleCPFOrCNPJ` | `cpf_or_cnpj` | either of the two |
| `br.RuleCEP` | `cep` | a CEP, `00000000` or `00000-000` |
| `br.RuleMobile` | `celular` | a mobile number, `00000000000` or `(00) 00000-0000`, with a real area code and the leading 9 |
| `br.RuleBRL` | `brl` | an amount in reais, as `ParseBRL` reads it |

Write the constant, not the word it holds. A rule name exists once it is in the
catalogue, and a rule set is compiled in a package-level variable, before
`main`. Naming `br.RuleCPF` imports this package into the file that declares
the set, and Go initializes an imported package before the variables of the one
importing it, so the rules are registered before the set is compiled. The word
`"cpf"` in a package that does not import this one may be compiled first, and
then the application refuses to boot with an unknown rule.

Each rule checks what the matching `Parse` function checks, takes no argument
(`cpf:anything` refuses every value) and only a string: a CPF sent as a JSON
number has already lost its leading zeros. A blank value is left to `required`
or `nullable`.

A failure reads `is not a valid CPF` (`CNPJ`, `CPF or CNPJ`, `CEP`, `mobile
number`, `amount in reais`). An application that speaks Portuguese names its
own sentence under `validation.cpf`, `validation.cnpj`,
`validation.cpf_or_cnpj`, `validation.cep`, `validation.celular` and
`validation.brl` in its translations, which `hesape/validation` reads before
the default.

## Take a value in

```go
cpf, err := br.ParseCPF(in.String("cpf"))    // eleven digits, bare
cnpj, err := br.ParseCNPJ(in.String("cnpj")) // fourteen characters, bare, letters in capitals
cep, err := br.ParseCEP(in.String("postcode"))
```

Each `Parse` function returns the bare spelling, which is the one to store, or
`ErrInvalidCPF`, `ErrInvalidCNPJ` or `ErrInvalidCEP`. `ValidCPF`, `ValidCNPJ`
and `ValidCEP` answer the same question with a `bool`.

What is accepted is deliberately narrow: the value bare, or in its one
conventional punctuation, and nothing else. No other separator, no separator
out of place, no whitespace around it — a form that trimmed or re-punctuated
input would accept a value the person did not type.

- **CPF** — the two check digits the Receita Federal computes. The eleven
  repetitions of one digit compute too and belong to nobody, so they are
  refused by name.
- **CNPJ** — the alphanumeric scheme: twelve digits or capital letters, each
  worth its ASCII code minus 48, then two numeric check digits. The numeric CNPJ
  is the all-digit case and is checked by the same computation. A lower-case
  letter is read as its capital. Fourteen repetitions of one character are
  refused.
- **CEP** — eight digits, and not below `01000-000`, where the Correios ranges
  begin. There is no check digit, so a CEP that passes may still address
  nothing.

None of them says whether a document was issued, is regular or exists: that is
a question for the Receita Federal or the Correios, and this package makes no
call.

`CPFCheckDigits` and `CNPJCheckDigits` compute the check digits of a base, which
is what a seeder or a factory needs to write a valid document without borrowing
somebody's.

### An amount in reais

```go
cents, err := br.ParseBRL(in.String("budget")) // int64 centavos, or ErrInvalidBRL
```

`ParseBRL` reads an amount as a person types it, and as an HTML number input
submits it:

| typed | centavos |
| --- | --- |
| `R$ 1.234,56`, `R$1.234,56`, `1.234,56`, `1234,56` | `123456` |
| `1.500` | `150000` — a dot followed by three digits groups thousands |
| `1234.56`, `1,5`, `1.5` | `123456`, `150`, `150` — a dot followed by one or two digits is the decimal point |
| `R$ 47`, `0` | `4700`, `0` |

Whitespace around the amount and after `R$` is ignored. A sign, another
currency, a space or a dot elsewhere among the digits, a group of other than
three digits, a third decimal and an amount past the largest `int64` are
refused. Zero is an amount; a form that needs a positive one says so itself. A
number input that takes centavos carries `step="0.01"`, which keeps a third
decimal — `1.500` meaning one and a half — out of what it submits.

## Print a value

```go
br.CPF("11122233344")      // "111.222.333-44", check digits or not
br.CNPJ("12abc34501de35")  // "12.ABC.345/01DE-35"
br.CEP("70040010")         // "70040-010"
br.BRL(123456)             // "R$ 1.234,56"
br.BRL(-99)                // "-R$ 0,99"
br.Date(t.In(saoPaulo))    // "31/12/2026"
br.DateTime(t.In(saoPaulo)) // "31/12/2026 15:04"

br.FormatBRL(123456, br.BRLFormat{WithoutGrouping: true}) // "R$ 1234,56"
br.FormatBRL(4700, br.BRLFormat{WithoutZeroCents: true})  // "R$ 47"
br.FormatBRL(19950, br.BRLFormat{WithoutZeroCents: true}) // "R$ 199,50"
br.FormatBRL(123456, br.BRLFormat{WithoutSymbol: true, WithoutGrouping: true}) // "1234,56", for a form field

br.BRLDecimal("1234.5678") // "R$ 1.234,5678", every digit kept
```

`CPF`, `CNPJ` and `CEP` punctuate a value of the right shape and return
anything else unchanged: they format and do not validate, so a page that prints
a stored value shows what is stored.

`BRL` takes centavos as an `int64` — an amount that has to add up is held as an
integer — and writes every one, the most negative included. It writes the same
string as `number.CurrencyFromCents(cents, "BRL")` with the process-wide locale
set to `pt-BR`, and depends on no process-wide setting.

`FormatBRL` is `BRL` departing from it in what its `BRLFormat` names, and the
zero `BRLFormat` is `BRL` itself. Every amount it prints, `ParseBRL` reads back
-- `WithoutSymbol` is how a form field is filled with an amount it hands back.

`BRLDecimal` writes an amount finer than a centavo -- a price per unit of
usage, the exact total a charge was rounded from -- held as decimal text, the
way `math/big` and a decimal column write it: the digits of the reais, then a
dot and every digit of the fraction. Nothing is rounded, and a value that is
not a decimal is returned unchanged.

`Date` writes the day `t` holds in its own location, and `DateTime` the day
and the time of day, `31/12/2026 15:04`. A time read from the
database is usually in UTC, where a late evening in Brasília is already the next
day, so convert it with `t.In(location)` first.

## Hide part of a document

```go
br.RedactCPF("12345678900", br.KeepCPFMiddle)      // "***.456.789-**"
br.RedactCPF("123.456.789-00", br.KeepCPFLastFive) // "***.***.789-00"
br.RedactCPF("12345678900", "•••.###.###-••")      // "•••.456.789-••"
br.RedactCNPJ("12ABC34501DE35", br.KeepCNPJEnds)   // "12.***.***/****-35"
```

A `Redaction` is spelled like the document: `#` shows the character in that
position and every other character is written as it stands, the separators
included. It has the length of the punctuated spelling or of the bare one, and
the output takes that spelling.

| redaction | pattern | prints |
| --- | --- | --- |
| `br.KeepCPFMiddle` | `***.###.###-**` | as the Portal da Transparência publishes a CPF |
| `br.KeepCPFMiddleBare` | `***######**` | as the Receita Federal's open CNPJ data publishes a partner's CPF |
| `br.KeepCPFCheckDigits` | `***.***.***-##` | only the two check digits, which identify nobody |
| `br.KeepCPFLastFive` | `***.***.###-##` | what a profile page shows its owner |
| `br.KeepCNPJEnds` | `##.***.***/****-##` | the first two and the last two characters |

`RedactCPF` and `RedactCNPJ` read the value as `ParseCPF` and `ParseCNPJ` do,
bare or punctuated, and do not check the check digits. Anything else — and a
redaction of neither length — gives `""`, never the value: a redaction that
fell back to its input would print what it exists to hide.

## The federative units

`br.UF` is the closed set of the twenty-seven units, shaped like an enum `aru
make:enum` writes:

```go
br.UFSP.String()  // "SP"
br.UFSP.Label()   // "São Paulo"
uf, err := br.ParseUF(in.String("uf")) // exact code, in capitals, or ErrInvalidUF
```

`br.UFValues()` lists them in the alphabetical order of their names. The type
satisfies `enum.Enum`, so the rest of the collection reads it without a second
list:

```go
enum.Options(br.UFValues()...)        // the options of a select
validation.EnumOf(br.UFValues()...)   // the rule
```

As a column, `Value` and `Scan` refuse a code outside the set.

## What this package does not do

| reached for | why not |
| --- | --- |
| `New`, a `Module`, a line in `k.Register` | there is nothing to register: no route, no table, no policy |
| a lookup at the Receita Federal or the Correios | a validator that calls out is a validator that fails when the network does, and the manifest declares `network = false` |
| accepting `123.456.789 09`, ` 01000-000` or full-width digits | the caller would store a value the person did not type |
| a `float64` amount | `BRL` takes centavos; a float has already rounded |
| reading the process locale | the same call would print differently depending on what another part of the application set |

## Layout

```
cpf.go, cnpj.go, cep.go   parsing, check digits and punctuation
mobile.go                 the mobile number
shape.go                  the masks they share
redact.go                 RedactCPF, RedactCNPJ and the Redaction type
rules.go                  the hesape/validation rules
money.go, date.go         BRL, FormatBRL, BRLDecimal, ParseBRL, Date and DateTime
uf.go                     the UF enum
```

## Tests

```bash
GOWORK=off go test -race ./...
```

Every valid document in the suite is generated from a fixed seed by an oracle
written apart from the package, from the Receita Federal's weight tables; none
is copied from anywhere. The parsers are fuzzed against that oracle:

```bash
GOWORK=off go test -run='^$' -fuzz='^FuzzParseCNPJ$' -fuzztime=30s ./tests/Unit
```

## Licence

MIT. See [LICENSE.md](LICENSE.md). Copyright Paulo R. Lima.
