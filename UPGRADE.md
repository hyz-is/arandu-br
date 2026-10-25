# Upgrade Guide

## Unreleased: v0.6.0

A minor release that breaks the registration API of v0.5.0 on purpose, to match
the one `hesape/validation` v0.88.0 provides for a module: the rules are a
`validation.Provider`, and the application hands it to the Factory it builds.
`apidiff` against `v0.5.0` reports `Rule` removed and `Rules` changed from
`func() []Rule` to `func() validation.Provider`, and nothing else incompatible.
The rules, their names, their aliases and what they accept are unchanged.

The package requires Hesape `v0.88.0` (was `v0.85.0`):

```bash
go get github.com/arandu-io/hesape@v0.88.0
```

`arandu.mod.toml` keeps `framework = ">= 0.66"`: no Framework release is built
on that Hesape yet, so there is no higher floor to name. An application whose
Framework does not hand a `validation.Factory` to the request still builds the
Factory itself, as below.

### Use the module on the Factory

```go
validators := validation.NewFactory(validation.FromTranslation(tr))
validators.Use(br.Rules())
```

`br.Install(validators)` is that call and is deprecated: it does nothing else,
and it ignores a nil factory, where v0.5.0 registered the rules in the
process-wide catalogue for one. Replace it:

```bash
grep -rn 'br\.Install' --include='*.go' .
```

The module declares its rules in the registry of the Factory it was handed to.
Nothing is registered before the application calls `Use`, in an `init` or
anywhere else, and a second Factory has its own registry.

### Rules and Defs

`br.Rules()` was the list of the rules; it is the Provider now, and the list is
`br.Defs()`, a `[]validation.RuleDef`:

| v0.5.0 | v0.6.0 |
| --- | --- |
| `br.Rules()` returned `[]br.Rule` | `br.Defs()` returns `[]validation.RuleDef` |
| `rule.Name`, `rule.Aliases` | the same fields |
| `rule.Check(value, parameters)` | `def.Check(validator, attribute, value, parameters)`, a `validation.ExtensionFunc` |
| `rule.Parameters` | `def.Params`, a `validation.Params{Min, Max}` |
| `rule.Messages["pt-BR"]`, `rule.Message(locale)` | `def.Messages["pt-BR"]`; `br.Messages(locale)` still lists them by name |
| `rule.Replace(message, parameters)` | `def.Replace`, a `validation.ReplacerFunc` |
| `rule.Extension()` | `def.Check` |
| `rule.Names()` | `append([]string{def.Name}, def.Aliases...)` |

Find the uses:

```bash
grep -rnE 'br\.(Rules|Rule)\b' --include='*.go' .
```

### Rule sets that name the rules

A rule set compiled in a package-level variable names a rule before any Factory
exists. Compile it with `validation.WithModuleRules()`, and check the sets at
boot, after `Use`, with the Factory:

```go
var StoreCustomer = validation.MustCompile(validation.Rules{
	"cpf":      []any{"required", "cpf"},
	"postcode": "required|formato_cep",
}, validation.WithModuleRules())

if err := validators.Check(StoreCustomer); err != nil {
	return err // names the field, the rule and the file of every rule nothing declares
}
```

A request type that declares `Rules()` is compiled by the Factory, after `Use`,
and needs neither. This replaces v0.5.0's instruction to compile the sets after
`Install`.

### What changed in how a rule is written

- **An argument is a boot failure.** A rule that takes none does not compile
  with one: `cpf:strict` is refused when the set is compiled, where v0.5.0
  compiled it and refused every value. The eight names registered when the
  package is imported (below) are the exception, and keep refusing every value.
- **`inscricao_estadual` is checked when it is compiled.** It takes one
  argument, the code of a unit: `inscricao_estadual:SP`. Without one, or with a
  code that is not one of the twenty-seven, the set does not compile.
- **The sentence is chosen by the locale of the request**, and a line the
  application defines under `validation.<rule>` is read before it. Every rule
  declares a pt-BR sentence only, as the PHP packages do.

### The registration at import

The eight names the earlier releases registered when the package is imported --
`cpf`, `cnpj`, `cpf_or_cnpj`, `cpf_ou_cnpj`, `cep`, `celular`,
`celular_com_ddd` and `brl` -- are still registered that way, deprecated, and
`Use` leaves those names to it. They are in the process-wide catalogue, so they
are not in the registry of the Factory, say their pt-BR sentence whatever the
locale, and are not counted among `validators.Registry().Names()`. A later
release removes the registration and declares them in the registry like the
others; nothing in an application changes when it does, as long as it calls
`Use`.

## v0.5.0

A minor release. It adds the rest of the Brazilian rules under the names the PHP
packages `geekcom/validator-docs` and `laravellegends/pt-br-validator` use, and
it breaks two things on purpose: the rule `celular` changes meaning, and the
sentences of every rule are in Portuguese now. Both are below, with what to
search for. `apidiff` against `v0.4.1` reports two constants whose value
changed, `RuleCPFOrCNPJ` and `RuleMobile`, and no other incompatible change.

### The rule celular is the number without the area code

`celular` accepted a mobile number with its area code, `(11) 91234-5678` or
`11912345678`. In the PHP packages it is the number without one,
`91234-5678`, and the one with the area code is `celular_com_ddd`. This package
follows the PHP packages, so the names mean what they mean there:

| the rule | accepts | it was |
| --- | --- | --- |
| `celular` | `91234-5678`, `912345678` | not a rule of this package |
| `celular_com_ddd` | `(11) 91234-5678`, `(11)91234-5678`, `11912345678` | what `celular` accepted |

Find the uses and rename them:

```bash
grep -rnE '"[^"]*\bcelular\b[^"]*"' --include='*.go' .
```

A rule string that holds `celular` and is checked against a number with an area
code (`required|celular`) becomes `required|celular_com_ddd`. A rule set that
names the constant, `br.RuleMobile`, changes nothing: the constant's value is
now `celular_com_ddd`, so it still checks the number with an area code. A new
constant, `RuleLocalMobile`, holds `celular`. A form that sends the number
without the area code and used to be refused now passes `celular`; one that
sends the area code and was passing `celular` now fails it, silently, until the
rule is renamed.

`RuleCPFOrCNPJ` holds `cpf_ou_cnpj` now, the PHP name, where it held
`cpf_or_cnpj`. The two names are the same rule and both stay registered, so
nothing needs to change.

### The sentences are Portuguese

Every rule says what the PHP packages say, in Portuguese, with the attribute:
`O campo :attribute não é um CPF válido.`, where it said `is not a valid CPF`
with no attribute. A line the application's translations hold under
`validation.cpf`, `validation.cnpj`, `validation.cpf_or_cnpj`, `validation.cep`,
`validation.celular` or `validation.brl` is read before the default, as it
always was, and is unchanged; one that depended on the English default, in a
test that compares a message, is the only thing to update. `br.Messages("pt-BR")`
returns the sentence of every rule by name, for an application that loads them
into its language files.

### Install the rules at boot

The rules are declared as plain values, `br.Rules()`, and installed by one
explicit call where the application is wired:

```go
validators := validation.NewFactory(translator)
br.Install(validators)
```

It registers every rule and alias below with the pt-BR sentence, and gives the
factory the filler of the `:uf` placeholder of `inscricao_estadual`. A rule set
that names one of them is compiled after that call, from the function that
wires the application, not in a package-level variable. Registering at import
is deprecated: when this package is imported it still registers `cpf`, `cnpj`,
`cpf_ou_cnpj`, `cpf_or_cnpj`, `cep`, `celular`, `celular_com_ddd` and `brl`, so
a rule set compiled before `main` that names them by their constants keeps
working in this release. Move the sets after `br.Install` and stop depending on
the import; a later release removes the registration at import, and
`br.Install` registers the same names, so nothing else changes.

### The rules added

Every row is a new rule unless it says otherwise. A value is accepted bare or in
its one conventional punctuation, as in the rest of the package; the PHP
packages forgive any character that is not a digit, and these do not.

| rule | constant | accepts |
| --- | --- | --- |
| `cpf`, `cnpj`, `cep`, `brl` | `RuleCPF`, `RuleCNPJ`, `RuleCEP`, `RuleBRL` | unchanged |
| `cpf_ou_cnpj`, aliases `cpf_or_cnpj`, `cpf_cnpj` | `RuleCPFOrCNPJ` | a CPF or a CNPJ |
| `moeda`, alias of `brl` | `RuleMoney` | an amount in reais |
| `formato_cpf` | `RuleCPFFormat` | `000.000.000-00`, whatever the digits |
| `formato_cnpj` | `RuleCNPJFormat` | `00.000.000/0000-00`, digits or capitals in the first twelve |
| `formato_cpf_ou_cnpj`, alias `formato_cpf_cnpj` | `RuleCPFOrCNPJFormat` | either shape |
| `formato_cep` | `RuleCEPFormat` | `00000-000` or `00.000-000` |
| `celular` | `RuleLocalMobile` | `91234-5678` or the nine digits |
| `celular_com_ddd` | `RuleMobile` | `(11) 91234-5678`, `(11)91234-5678` or the eleven digits |
| `celular_com_codigo` | `RuleMobileWithCountry` | `+55(11)91234-5678`, `+55 (11) 91234-5678`, `+5511912345678` |
| `celular_com_codigo_sem_mascara` | `RuleMobileWithCountryBare` | `+5511912345678`, `+55 11 912345678` |
| `telefone` | `RuleLocalLandline` | `3000-0000` or the eight digits |
| `telefone_com_ddd` | `RuleLandline` | `(11) 3000-0000`, `(11)3000-0000` or the ten digits |
| `telefone_com_codigo` | `RuleLandlineWithCountry` | `+55(11)3000-0000`, `+5511 30000000` |
| `ddd` | `RuleAreaCode` | one of the 67 area codes Anatel assigns |
| `placa` | `RulePlate` | every plate Brazil issued, the Mercosul included, with or without the hyphen |
| `formato_placa_de_veiculo` | `RulePlateFormat` | three letters, a hyphen or none, a digit, a digit or a letter, two digits |
| `renavam` | `RuleRenavam` | nine, ten or eleven digits, with the RENAVAM check digit |
| `cnh` | `RuleCNH` | eleven digits, with the CNH check digits |
| `titulo_eleitor` | `RuleVoterTitle` | five to thirteen digits, or `0000 0000 0000`, with the state and the check digits |
| `nis`, `pis` | `RuleNIS`, `RulePIS` | `000.00000.00-0` or the eleven digits, with the check digit |
| `formato_nis`, `formato_pis` | `RuleNISFormat`, `RulePISFormat` | `000.00000.00-0`, whatever the digits |
| `cns` | `RuleCNS` | `000 0000 0000 0000` or the fifteen digits, definitive or provisional |
| `certidao` | `RuleCivilCertificate` | the thirty-two digits, bare or in nine groups, with the check digits |
| `formato_certidao` | `RuleCivilCertificateFormat` | the nine groups, with dots, hyphen or spaces |
| `uf` | `RuleUF` | one of the 27 codes, in capitals |
| `passaporte` | `RulePassport` | two letters and six digits |
| `inscricao_estadual:UF` | `RuleStateRegistration` | the state registration number of the unit in the argument, `inscricao_estadual:SP` |
| `data` | `RuleDate` | `31/12/2026`, a day the month has |
| `horario` | `RuleTimeOfDay` | `15:04` or `15:04:05` |

Each has a function that takes the value in, `ParseCNH`, `ParseNIS`,
`ParseCNS`, `ParseRenavam`, `ParseVoterTitle`, `ParseCivilCertificate`,
`ParsePlate`, `ParsePassport`, `ParseAreaCode`, `ParseLocalMobile`,
`ParseMobileWithCountry`, `ParseLandline`, `ParseLocalLandline`,
`ParseLandlineWithCountry`, `ParseStateRegistration`, `ParseDate` and
`ParseTimeOfDay`, with a `Valid` function beside it. `NormalizeCPF` and
`NormalizeCNPJ` drop the punctuation of text that already left the form, keeping
the letters of a CNPJ; `NIS`, `CNS` and `Landline` punctuate.

The CPF and the CNPJ are text from the form to the column: a document read as a
number has lost its leading zeros, and the alphanumeric CNPJ cannot be a number
at all. Every rule refuses a value that is not a string.

### Where the two PHP packages disagree

Where `geekcom/validator-docs` and `laravellegends/pt-br-validator` accept
different values, the rule takes one reading and the reason is here. The
README lists what each rule refuses that a PHP package accepts.

- **Punctuation.** Both packages remove every character that is not a digit
  before checking, so `094 050 986 59` passes. This package accepts the value
  bare or in its one conventional punctuation, which is what `ParseCPF` always
  did, so a value the person did not type is never stored.
- **CNPJ.** Both packages compute the alphanumeric check digits the same way,
  with each character worth its ASCII code minus 48. geekcom refuses only the
  zeroed CNPJ and laravellegends every repetition of one character; every
  repetition is refused. The `formato_cnpj` rule takes capitals only, as in the
  laravellegends rule.
- **CNH.** The packages compute different second digits. geekcom's is the
  Denatran's, with the correction that adds eleven back to a negative result;
  laravellegends' lacks it, refuses about one valid number in forty, such as
  `04463004100`, and accepts a few that are not valid. The Denatran's is used.
- **Mobile and landline numbers.** The packages accept eight and nine digits
  for both. A mobile number has nine digits that start with 9 since 2016, so
  `celular` refuses `9800-1936`; the eight-digit forms are `telefone`, whose
  first digit is not 0 or 1. The country code is `+55`: the number is checked
  against Brazil's plan, so `+1(31)99999-5555` is refused where laravellegends
  accepts it.
- **CNS.** laravellegends reads a definitive card that begins with any digit
  below 7; the number begins with 1 or 2, as geekcom's pattern has it, and a
  provisional one with 7, 8 or 9.
- **NIS.** geekcom pads a number shorter than eleven digits with zeros; a
  leading zero is a digit of the number, so a shorter one is refused. `nis` and
  `pis` check the same number, as both packages do.
- **RENAVAM.** The nine and ten digit numbers are read with the zeros they
  lost, as geekcom does, and come back with eleven. Eleven zeros are refused.
- **Plate.** `placa` is the six shapes geekcom lists, with laravellegends' hyphen
  after the three letters; `formato_placa_de_veiculo` is laravellegends' shape
  alone, so `AAA12A3` passes `placa` and not `formato_placa_de_veiculo`.
- **Voter title.** geekcom's rule, from five to thirteen digits.
- **State registration.** The twenty-seven algorithms of the library geekcom
  uses. Eleven zeros, which that library accepts for Mato Grosso, are refused.
- **Repetitions.** A number of one digit repeated is refused by every document
  with check digits, as `ParseCPF` always did.

### Not a rule yet

`data` and `horario` are checked here; neither PHP package registers a rule
by those names, and the Arandu mission lists them. Laravel's own `date_format`
rule reads other shapes and is unchanged.

### Go 1.26.10 and Hesape

`go.mod` declares `go 1.26.10`. Every vulnerability an Arandu application
reaches is in the Go standard library (crypto/x509, the HTTP/2 and CONNECT
handling in net/http, html/template, net/textproto and crypto/tls) and is
fixed in Go 1.26.10, so a module that declared less let a build run on a
toolchain that still carries them. Build with Go 1.26.10 or later.

The package requires Hesape `v0.85.0` (was `v0.54.0`), and `arandu.mod.toml`
declares `framework = ">= 0.66"`, the first Framework release built on that
Hesape:

```bash
go get github.com/arandu-io/framework@v0.66.0
go get github.com/arandu-io/hesape@v0.85.0
```

An application below those follows the Framework and Hesape upgrade guides
between the versions it had and these first.

## v0.4.1

Nothing to change. `ParseMobile` now also accepts `(00)00000-0000`, the
punctuated number without the space after the area code, which it refused
before; every value it accepted still parses to the same digits.

## v0.4.0

Nothing changes for code already using the package: every addition is new,
and `BRL`, `FormatBRL` and `Date` write what they wrote before. An application
that wrote these by hand replaces its own with the package's:

- a hand-written `fmt.Sprintf("%d,%02d", cents/100, cents%100)` that fills a
  form field becomes `br.FormatBRL(cents, br.BRLFormat{WithoutSymbol: true,
  WithoutGrouping: true})`, which writes a negative amount with its sign
  (`-12,34`) where the hand-written one wrote `-12,-34`;
- a local formatter of an exact decimal amount in reais -- a usage price, a
  residual of a rounding -- becomes `br.BRLDecimal(decimal)`, given the text
  `math/big` or a decimal column writes (`1234.5678`), not a float;
- `t.Format("02/01/2006 15:04")` becomes `br.DateTime(t)`.

## v0.3.0

Nothing changes for code already using the package: every addition is new,
and `BRL` writes what it wrote before. An application that hid part of a
document, read an amount in reais or printed one in a variant of `BRL`
replaces its own function with the package's:

- a local `maskCPF`, `MaskCPF`, `MaskedCPF` or `maskCNPJ` becomes
  `br.RedactCPF(cpf, r)` or `br.RedactCNPJ(cnpj, r)`, with a named
  `Redaction` (`br.KeepCPFMiddle`, `br.KeepCPFMiddleBare`,
  `br.KeepCPFCheckDigits`, `br.KeepCPFLastFive`, `br.KeepCNPJEnds`) or its own
  literal, such as `"•••.###.###-••"`. The redaction answers `""` for a value
  it cannot read, where a local version may have hidden all of it, kept it, or
  trimmed and re-punctuated it first; an application that wants another
  answer for that case writes it beside the call;
- a local BRL parser becomes `br.ParseBRL`, and a rule set names
  `br.RuleBRL`;
- a hand-written `fmt.Sprintf("R$ %d,%02d", ...)` becomes
  `br.FormatBRL(cents, br.BRLFormat{WithoutGrouping: true})`, and a price
  printed without zero centavos becomes
  `br.FormatBRL(cents, br.BRLFormat{WithoutZeroCents: true})`.

`ParseBRL` reads one set of spellings, and two that applications read
before differ from it. Check what a form received before switching:

- a dot followed by one or two digits is the decimal point (`1.50` is one
  real and fifty centavos, `12.34` twelve reais and thirty-four), and a dot
  followed by three digits groups thousands (`1.234` is one thousand two
  hundred and thirty-four reais). A parser that dropped every dot read the
  first as a hundred and fifty reais; one that read every dot as the decimal
  point refused the second;
- a dot anywhere else (`1.2345`, `1..000`, `.0,01`), a space between the
  digits (`1 234`) and a sign (`1.-5`, `-10`) are refused;
- zero is an amount: a form that needs a positive one keeps checking
  `cents > 0` itself;
- the symbol `R$` and whitespace around the amount are accepted, so a form
  that refused them now takes them.

## v0.2.0

Nothing changes for code already using the package. The release adds the
mobile number: an application that validated or formatted one itself (a local
`NormalizeMobile`, a phone regexp) replaces it with `br.ParseMobile`,
`br.ValidMobile` or `br.Mobile`, and a rule set names `br.RuleMobile`. The
parser accepts the eleven digits bare or as `(00) 00000-0000` and nothing
else: a `+55` prefix or a space-separated spelling is refused, so a form uses
the mask.

## v0.1.2

Nothing in the package changes. The release requires Hesape v0.54.0 and
declares `framework = ">= 0.56"`, so `go get` raises an application's Hesape
to v0.54.0. An application on an older Framework moves to v0.56.0 first,
reading the Framework and Hesape upgrade guides from the versions it required
before: Framework v0.56.0 removed `config.Config.SessionTTL`, which an
application that boots with `bootstrap.LoadConfiguration` never read.

## v0.1.1

Nothing in the package changes. The release requires Hesape v0.52.0 and
declares `framework = ">= 0.55"`, so `go get` raises an application's Hesape
to v0.52.0. An application on an older Framework moves to v0.55.1 first,
reading the Framework and Hesape upgrade guides from the versions it required
before: Hesape v0.52.0 removed the names it had deprecated, and Framework
v0.54.0 and v0.55.0 refuse at boot a malformed boolean setting and an unread
`SESSION_*` setting.

## v0.1.0

The first release. There is nothing to upgrade from.

An application that validated or formatted these values itself moves to this
package by replacing its own functions, not by wrapping them:

- a local `NormalizeCPF`, `ValidCPF` or `FormatCPF` becomes `br.ParseCPF`,
  `br.ValidCPF` or `br.CPF`, and the same for the CNPJ and the CEP;
- a rule set that called one of them by hand names `br.RuleCPF`,
  `br.RuleCNPJ`, `br.RuleCPFOrCNPJ` or `br.RuleCEP` instead;
- a BRL formatter becomes `br.BRL`, with the amount in centavos as an `int64`;
- a list of the twenty-seven units becomes `br.UFValues()`.

Two of those can change what an application accepts, and both are worth a
look before the switch. The parsers accept the value bare or in its one
conventional punctuation, and refuse surrounding whitespace, where a local
version may have trimmed it. And a CNPJ may now carry letters in its first
twelve positions, which a column or a check that expected digits only has to
make room for.
