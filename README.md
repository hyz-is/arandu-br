<p align="center">
  <img src=".github/logo.svg" alt="Arandu" width="180">
</p>

# Arandu BR

The Brazilian values an Arandu application takes from a form and prints on a
page: the CPF, the CNPJ (numeric and alphanumeric), the CEP, the telephone
numbers, the documents with check digits (NIS, CNS, CNH, RENAVAM, título de
eleitor, certidão, the state registration), the vehicle plate, an amount in
reais, a date, and the federative unit — and a CPF or a CNPJ with part of it
hidden.

It is a library. It registers no route, owns no table and decides no
authorization, so there is nothing to construct and no `aru migrate` to run. An
application imports it, calls it, and hands its validation rules to the
Factory with one call, `factory.Use(br.Rules())`, where it is wired in
`bootstrap/app.go`.

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

The rules are a `validation.Provider`, `br.Rules()`, that the application hands
to the Factory it builds, once, where it is wired. Nothing is registered before
that call:

```go
validators := validation.NewFactory(validation.FromTranslation(tr))
validators.Use(br.Rules())

var StoreCustomer = validation.MustCompile(validation.Rules{
	"cpf":      []any{"required", br.RuleCPF},
	"cnpj":     "nullable|" + br.RuleCNPJ,
	"document": "required|" + br.RuleCPFOrCNPJ,
	"postcode": "required|" + br.RuleCEPFormat,
	"phone":    "required|" + br.RuleMobile,
	"state":    "required|" + br.RuleStateRegistration + ":SP",
}, validation.WithModuleRules())

err := validators.Check(StoreCustomer) // at boot: names any rule nothing declares
```

`Use` declares every rule and alias, each as a `validation.RuleDef` with its
pt-BR sentence, and a `validation.<rule>` line in the application's
translations is read before that sentence. `br.Defs()` returns the
declarations. A rule set compiled in a package-level variable runs before the
Factory exists, so it is compiled with `validation.WithModuleRules()` and
resolved by the Factory; a request type that declares `Rules()` needs neither.

The names are those of `geekcom/validator-docs` and
`laravellegends/pt-br-validator`, so a rule set written for either reads the
same, and the string form works: `"required|cpf"`.

| rule | constant | accepts |
| --- | --- | --- |
| `cpf` | `br.RuleCPF` | a CPF, `00000000000` or `000.000.000-00` |
| `cnpj` | `br.RuleCNPJ` | a CNPJ, numeric or alphanumeric, `00000000000000` or `00.000.000/0000-00` |
| `cpf_ou_cnpj` (`cpf_or_cnpj`, `cpf_cnpj`) | `br.RuleCPFOrCNPJ` | either of the two |
| `cep` | `br.RuleCEP` | a CEP, `00000000` or `00000-000`, not below `01000-000` |
| `formato_cpf`, `formato_cnpj`, `formato_cpf_ou_cnpj` (`formato_cpf_cnpj`), `formato_cep` | `br.RuleCPFFormat`, `br.RuleCNPJFormat`, `br.RuleCPFOrCNPJFormat`, `br.RuleCEPFormat` | the punctuated shape alone, whatever the digits: `000.000.000-00`, `00.000.000/0000-00`, `00000-000` or `00.000-000` |
| `celular` | `br.RuleLocalMobile` | a mobile number without the area code, `91234-5678` or `912345678` |
| `celular_com_ddd` | `br.RuleMobile` | `(11) 91234-5678`, `(11)91234-5678` or `11912345678`, with a real area code |
| `celular_com_codigo` | `br.RuleMobileWithCountry` | `+55(11)91234-5678`, `+55 (11) 91234-5678` or `+5511912345678` |
| `celular_com_codigo_sem_mascara` | `br.RuleMobileWithCountryBare` | `+5511912345678` or `+55 11 912345678` |
| `telefone`, `telefone_com_ddd`, `telefone_com_codigo` | `br.RuleLocalLandline`, `br.RuleLandline`, `br.RuleLandlineWithCountry` | a landline number, `3000-0000`, `(11) 3000-0000`, `+55(11)3000-0000` |
| `telefone_ou_celular` | `br.RuleAnyPhone` | a mobile or a landline number, with or without the area code: `91234-5678`, `3456-7890`, `(11) 91234-5678`, `(11) 3456-7890`, bare or punctuated |
| `ddd` | `br.RuleAreaCode` | one of the 67 area codes Anatel assigns |
| `placa` | `br.RulePlate` | any plate Brazil issued, the Mercosul one included, `ABC-1234`, `ABC1D23` |
| `formato_placa_de_veiculo` | `br.RulePlateFormat` | three letters, a hyphen or none, a digit, a digit or a letter, two digits |
| `renavam` | `br.RuleRenavam` | nine to eleven digits with the RENAVAM check digit |
| `cnh` | `br.RuleCNH` | the eleven digits of the register number, with the Denatran check digits |
| `titulo_eleitor` | `br.RuleVoterTitle` | the digits of a título de eleitor, five to thirteen, with its state and check digits |
| `nis`, `pis` | `br.RuleNIS`, `br.RulePIS` | `000.00000.00-0` or the eleven digits, with the check digit |
| `formato_nis`, `formato_pis` | `br.RuleNISFormat`, `br.RulePISFormat` | `000.00000.00-0`, whatever the digits |
| `cns` | `br.RuleCNS` | the fifteen digits of the Cartão Nacional de Saúde, definitive or provisional |
| `certidao`, `formato_certidao` | `br.RuleCivilCertificate`, `br.RuleCivilCertificateFormat` | the thirty-two digits of a birth, marriage or death certificate, and their nine groups |
| `uf` | `br.RuleUF` | one of the 27 codes, in capitals |
| `inscricao_estadual:UF` | `br.RuleStateRegistration` | the state registration number of the unit, `inscricao_estadual:SP` |
| `passaporte` | `br.RulePassport` | two letters and six digits |
| `data`, `horario` | `br.RuleDate`, `br.RuleTimeOfDay` | `31/12/2026`, a day the month has; `15:04` or `15:04:05` |
| `brl` (`moeda`) | `br.RuleBRL` | an amount in reais, as `ParseBRL` reads it |

Which telephone rule: a field of a mobile number with the area code takes
`celular_com_ddd`; a field with a dynamic mask, mobile or landline, with or
without the area code, takes `telefone_ou_celular` (Laravel's rules add up, so
`celular_com_ddd|telefone_com_ddd` refuses everything); a field of a landline
number takes `telefone_com_ddd` or `telefone`.

Each rule checks what the matching `Parse` function checks, and only a string:
a CPF sent as a JSON number has already lost its leading zeros. A blank value is
left to `required` or `nullable`. A rule written with an argument it does not
take does not compile (`cpf:strict`; the eight names registered at import refuse
every value instead), and `inscricao_estadual` takes exactly one, the code of a
unit, in either case, and does not compile with another.

A failure reads as the PHP packages' sentences do, in Portuguese, with the
attribute: `O campo documento não é um CPF válido.` An application names its own
sentence under `validation.cpf`, `validation.cnpj`, `validation.formato_cep`
and so on in its translations, which `hesape/validation` reads before the
default; `br.Messages("pt-BR")` returns every default, by rule name, for an
application that loads them into its language files.

### Where it is narrower than the PHP packages

Both packages remove every character that is not a digit before checking. This
package takes a value bare or in its one conventional punctuation, so a form
never stores what the person did not type; trim or re-punctuate in the request
if the form should forgive it. Beyond that, a value the PHP rules accept and
these refuse:

| rule | refused here | why |
| --- | --- | --- |
| `cpf`, `cnpj`, `nis`, `cns`, `titulo_eleitor` and the rest | `094 050 986 59`, ` 094.050.986-59` | only the bare or the one punctuated spelling |
| `celular`, `celular_com_ddd` | `9800-1936`, `(31)9999-5555` | an eight-digit number is not a mobile number since 2016; it is a `telefone` |
| `celular_com_codigo`, `telefone_com_codigo` | `+1(31)99999-5555`, `+66(31)9999-5555` | the number is checked against Brazil's plan, so the code is `+55` |
| `cns` | a definitive card that begins with 0 or 3 to 6 | laravellegends computes it; a definitive card begins with 1 or 2 |
| `nis`, `pis` | a number shorter than eleven digits | geekcom pads it with zeros; a leading zero is a digit |
| `renavam` | `1970.73212`, `00000000000` | only digits; a repetition has no owner |
| `inscricao_estadual` | `8 2__6-7*9.3/41`, `00000000000` for MT | only dots, hyphens and slashes between digits; a repetition has no owner |
| `cnh` | — | laravellegends' arithmetic lacks the Denatran's correction and refuses 1 valid number in 40; the correct one is used |
| `placa` | — | `AAA12A3` passes `placa` (geekcom) and not `formato_placa_de_veiculo` (laravellegends) |

Examples the PHP packages test with are rewritten as data in
`tests/Unit/rule_examples_test.go`, and the ones refused here carry the reason.

### The import-time registration is deprecated

Until v0.5.0 the package registered `cpf`, `cnpj`, `cpf_or_cnpj`, `cep`,
`celular` and `brl` when it was imported. It still does, now with `cpf_ou_cnpj`
and `celular_com_ddd` beside them, so a rule set that names them keeps compiling
without a Factory. `Use` leaves those eight names to it. That registration goes
in a later release; call `factory.Use(br.Rules())`. `br.Install(factory)` is
that call under its old name, and is deprecated too.

`celular` means the number without the area code now, as in the PHP packages;
what it was is `celular_com_ddd`. `UPGRADE.md` says what to search for.

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

A CPF and a CNPJ are **text**, from the form to the column. A document read as
a number has lost its leading zeros (`00000000191` becomes `191`), and the
alphanumeric CNPJ cannot be a number at all. Declare the column `text` (or
`varchar`), never an integer, and take the value with `in.String`. For text that
already left the form, `NormalizeCPF` keeps the digits of a string and
`NormalizeCNPJ` keeps the digits and the letters, in capitals, so a stored
`12.abc.345/01de-35` compares with `12ABC34501DE35`; neither validates, and
`ParseCPF` and `ParseCNPJ` are what take a value in.

The other values follow the same shape, a `Parse` function that returns the
bare spelling or an `Err…` and a `Valid` function beside it:

| value | function | returns | an error |
| --- | --- | --- | --- |
| mobile number with area code | `ParseMobile` | 11 digits | `ErrInvalidMobile` |
| mobile number without it, or with `+55` | `ParseLocalMobile`, `ParseMobileWithCountry` | 9 digits, `55` and 11 digits | `ErrInvalidMobile` |
| landline | `ParseLandline`, `ParseLocalLandline`, `ParseLandlineWithCountry` | 10, 8, `55` and 10 digits | `ErrInvalidLandline` |
| any telephone number | `ParseAnyPhone` | 8, 9, 10 or 11 digits | `ErrInvalidPhone` |
| area code | `ParseAreaCode` | 2 digits | `ErrInvalidAreaCode` |
| NIS, PIS | `ParseNIS` | 11 digits | `ErrInvalidNIS` |
| Cartão Nacional de Saúde | `ParseCNS` | 15 digits | `ErrInvalidCNS` |
| CNH | `ParseCNH` | 11 digits | `ErrInvalidCNH` |
| RENAVAM | `ParseRenavam` | 11 digits, the zeros of an older number put back | `ErrInvalidRenavam` |
| título de eleitor | `ParseVoterTitle` | its digits | `ErrInvalidVoterTitle` |
| certidão | `ParseCivilCertificate` | 32 digits | `ErrInvalidCivilCertificate` |
| vehicle plate | `ParsePlate` | in capitals, no hyphen | `ErrInvalidPlate` |
| passport | `ParsePassport` | in capitals | `ErrInvalidPassport` |
| inscrição estadual | `ParseStateRegistration(uf, s)` | its digits | `ErrInvalidStateRegistration` |
| day, time of day | `ParseDate`, `ParseTimeOfDay` | `time.Time` at midnight UTC, `time.Duration` since midnight | `ErrInvalidDate`, `ErrInvalidTimeOfDay` |

`NIS`, `CNS` and `Landline` punctuate a stored value like `CPF`, and `FormatPhone` picks the mask of a telephone number by its digits (8, 9, 10 or 11). `NormalizePhone` keeps the digits of text that already left the form.

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
rules.go, catalog.go      Rules, Defs, Install and the rule names
phone.go                  the landline, the area code and the other spellings of a number
nis.go, cns.go, cnh.go    the documents with check digits
renavam.go, voter.go
certificate.go
plate.go                  the vehicle plate and the passport
stateregistration.go      the state registration number of the 27 units
normalize.go              NormalizeCPF and NormalizeCNPJ
masked.go                 the checks of the formato_ rules
money.go, date.go         BRL, FormatBRL, BRLDecimal, ParseBRL, Date, DateTime, ParseDate, ParseTimeOfDay
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
