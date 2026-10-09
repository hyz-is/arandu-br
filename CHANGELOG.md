# Changelog

Everything worth knowing about a release of Arandu BR is recorded here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and
the versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

A published module version is immutable: Go serves it from the proxy forever, so
a release is corrected by another release and never by moving a tag.

## [Unreleased]

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
