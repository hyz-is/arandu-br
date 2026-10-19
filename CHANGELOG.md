# Changelog

Everything worth knowing about a release of Arandu BR is recorded here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and
the versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

A published module version is immutable: Go serves it from the proxy forever, so
a release is corrected by another release and never by moving a tag.

## [Unreleased]

## [0.4.0] - 2026-10-19

### Added

- `BRLFormat.WithoutSymbol`: `FormatBRL` without the `R$` and the space after
  it, `1.234,56`, the way a form field is filled with an amount; with
  `WithoutGrouping`, `1234,56`. `ParseBRL` reads both back.
- `BRLDecimal`: an exact amount of reais held as decimal text, `1234.5678`,
  written with the symbol, the grouping and the comma and every digit of the
  fraction kept, `R$ 1.234,5678`. A value that is not a decimal is returned
  unchanged.
- `DateTime`: a day and a time of day, `31/12/2026 15:04`, in the location the
  time holds.

## [0.3.0] - 2026-10-15

### Added

- `RedactCPF` and `RedactCNPJ`: a CPF or a CNPJ written with part of it
  hidden, through a `Redaction` spelled like the document -- `#` shows a
  position, every other character is written as it stands -- in the
  punctuated or the bare length. A value of another shape, or a redaction of
  neither length, gives `""`, never the value itself; the check digits are
  not checked.
- The `Redaction` type and five named ones: `KeepCPFMiddle`
  (`***.456.789-**`), `KeepCPFMiddleBare` (`***456789**`),
  `KeepCPFCheckDigits` (`***.***.***-00`), `KeepCPFLastFive`
  (`***.***.789-00`) and `KeepCNPJEnds` (`12.***.***/****-35`).
- `ParseBRL` and `ValidBRL`: an amount in reais as a person types it --
  `R$ 1.234,56`, `1234,56`, `1.500`, `R$ 47`, whitespace around it -- or as an
  HTML number input submits it, `1234.56`; returned as centavos in an `int64`.
  A sign, a misplaced dot or space, a third decimal and an overflow are
  refused.
- `ErrInvalidBRL`.
- The `hesape/validation` rule `brl`, named by `RuleBRL`, registered with the
  other five when the package is initialized.
- `FormatBRL` and `BRLFormat`: `BRL` without the grouping (`R$ 1234,56`,
  `WithoutGrouping`) or without the centavos of a whole amount (`R$ 47`,
  `WithoutZeroCents`). `BRL` is `FormatBRL` with the zero `BRLFormat`, and
  writes what it wrote before.

## [0.2.0] - 2026-10-11

### Added

- `ParseMobile`, `ValidMobile` and `Mobile`: a Brazilian mobile number written
  bare or as `(00) 00000-0000`, with an area code Anatel assigns, the leading
  `9`, and the nine followed by eight zeros refused; returned as its eleven
  digits, and punctuated for a page.
- `ErrInvalidMobile`.
- The `hesape/validation` rule `celular`, named by `RuleMobile`, registered
  with the other four when the package is initialized.

## [0.1.2] - 2026-10-09

### Changed

- Requires Hesape v0.54.0. The manifest's framework floor is `>= 0.56`:
  Framework v0.56.0 is the first release built on that Hesape. No symbol,
  rule name, accepted value or returned spelling changes; `apidiff` against
  v0.1.1 reports nothing, and `hesape/validation` and `hesape/enum`, the two
  packages this one imports, did not change between Hesape v0.52.0 and
  v0.54.0.
- An application below those releases upgrades them first, following their
  upgrade guides from the versions it required before. Framework v0.56.0
  removed `config.Config.SessionTTL` and stopped reading `SESSION_TTL` in
  `config.Load`; an application that boots with `bootstrap.LoadConfiguration`
  compiles and boots unchanged.

## [0.1.1] - 2026-10-09

### Changed

- Requires Hesape v0.52.0. The manifest's framework floor is `>= 0.55`:
  Framework v0.55.1 is the first release built on that Hesape. No symbol,
  rule name, accepted value or returned spelling changes; `apidiff` against
  v0.1.0 reports nothing.
- An application below those releases upgrades them first, following their
  upgrade guides from the versions it required before. Hesape v0.52.0 removed
  the names it had deprecated; from Framework v0.55.0 an unread `SESSION_*`
  setting stops the boot, and from v0.54.0 a boolean setting that does not
  read as one does.

## [0.1.0] - 2026-10-09

### Added

- `ParseCPF`, `ValidCPF`, `CPFCheckDigits` and `CPF`: a CPF written bare or as
  `000.000.000-00`, checked against its two check digits, the repetitions of
  one digit refused, returned bare; and punctuated for a page.
- `ParseCNPJ`, `ValidCNPJ`, `CNPJCheckDigits` and `CNPJ`: the alphanumeric
  CNPJ, of which the numeric one is the all-digit case, written bare or as
  `00.000.000/0000-00`, a lower-case letter read as its capital.
- `ParseCEP`, `ValidCEP` and `CEP`: eight digits, bare or as `00000-000`, and
  not below `01000-000`.
- `ErrInvalidCPF`, `ErrInvalidCNPJ`, `ErrInvalidCEP` and `ErrInvalidUF`.
- The `hesape/validation` rules `cpf`, `cnpj`, `cpf_or_cnpj` and `cep`,
  registered when the package is initialized and named by `RuleCPF`,
  `RuleCNPJ`, `RuleCPFOrCNPJ` and `RuleCEP`. Each takes no argument and only a
  string, and leaves a blank value to `required` or `nullable`.
- `BRL`: centavos as `R$ 1.234,56`, the sign before the symbol, for every
  `int64`, reading no process-wide setting.
- `Date`: the day a `time.Time` holds in its own location, as `dd/mm/yyyy`.
- `UF`, the twenty-seven federative units: `UFValues`, `Valid`, `String`,
  `Label`, `ParseUF`, and `Value` and `Scan` for a column. It satisfies
  `enum.Enum`.
- `arandu.mod.toml` as `hyz-is/br`, declaring no capability and
  `framework = ">= 0.51"`; the module requires Hesape `v0.50.1` and nothing
  else.
- The `br-package` skill, marked `audience: app`, for an application that
  installs the package.
