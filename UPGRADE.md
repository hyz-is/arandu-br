# Upgrade Guide

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
