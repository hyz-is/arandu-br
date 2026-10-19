# Upgrade Guide

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
