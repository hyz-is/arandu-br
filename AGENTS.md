# Working on Arandu BR

This is an Arandu package that validates and formats Brazilian values: the CPF,
the CNPJ (numeric and alphanumeric), the CEP, the telephone numbers, the other
documents with check digits, the vehicle plate, amounts in reais, dates, and the
federative unit, and writes a CPF or a CNPJ with part of it hidden. It is a Go
module somebody `go get`s and imports.

It began as the package skeleton, configured as `hyz-is/br`, and the HTTP
module the template ships was removed on purpose: there is no entity, no table,
no migration, no policy, no route, no view and no `Module` to register. What a
skeleton package wires in `bootstrap/app.go`, this one does not need, and
nothing here should grow back into it.

Read `.agents/skills/` before writing code. Each skill is a procedure, and the
one you need is named by the situation you are in.

## The gates

Nothing is finished until all four exit zero.

```sh
export GOWORK=off
gofmt -l $(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go')
go build ./...
go vet ./...
go test -race ./...
```

`GOWORK=off` is not a preference. This checkout may sit beside a Go workspace
that lists the framework repositories and does not list this one; when it does,
every command above fails before it compiles anything:

```
pattern ./...: directory prefix . does not contain modules listed in go.work
or their selected dependencies
```

With the workspace off, the module resolves the Hesape version in `go.mod` —
which is what CI compiles against, and what somebody's `go get` will get.

Before a release, give each parser a short fuzz run as well. The seeds run on
every `go test`; the run is what looks past them:

```sh
go test -run='^$' -fuzz='^FuzzParseCPF$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseCNPJ$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseCEP$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseMobile$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParsePhones$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseAnyPhone$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseStateRegistration$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzBRLAgreesWithHesapeUnderPtBR$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseBRL$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzRedactCPF$' -fuzztime=30s ./tests/Unit
```

`aru doctor` is not one of the gates. It reads applications, and exits 1 here.
It does not read this package after it is installed either, which is why
`tests/Unit/audit_test.go` compares `arandu.mod.toml` with what the code calls.

## What this repository holds

| | measured with |
| --- | --- |
| 23 Go files in one package at the root | `grep -l '^package br' *.go` |
| 15 test files, 10 fuzz targets | `find tests -name '*_test.go'` · `grep -c '^func Fuzz' tests/Unit/*.go` |
| 0 routes, 0 tables, 0 capabilities | `arandu.mod.toml`, `TestThePackageIsALibrary` |
| 1 direct dependency, `github.com/arandu-io/hesape` | `go list -m -f '{{if and (not .Indirect) (not .Main)}}{{.Path}} {{.Version}}{{end}}' all` |

```
doc.go                    the package comment
cpf.go, cnpj.go, cep.go   parsing, check digits and punctuation
mobile.go, phone.go       the telephone numbers and the area code
nis.go, cns.go, cnh.go    the other documents with check digits
renavam.go, voter.go
certificate.go
plate.go                  the vehicle plate and the passport
stateregistration.go      the state registration number of the 27 units
normalize.go              NormalizeCPF and NormalizeCNPJ
shape.go, masked.go       the masks they share, and the checks of the formato_ rules
redact.go                 RedactCPF, RedactCNPJ and the Redaction type
rules.go, catalog.go      Rules (the Provider), Defs, Install, the rule names and the init of the deprecated registration
money.go, date.go         BRL, FormatBRL, ParseBRL, Date, DateTime, ParseDate and ParseTimeOfDay
uf.go                     the UF enum
```

## What does not exist here, on purpose

| A model reaches for | What is here instead |
| --- | --- |
| `New`, a `Module`, a `Config`, a line in `k.Register` | nothing. A caller imports the package and calls a function |
| a second registration function beside `br.Rules()`, or a registry the package fills as it is imported | `br.Defs()`, plain `validation.RuleDef` values, and `br.Rules()`, the `validation.Provider` the application passes to `factory.Use` where it is wired. `br.Install` is a deprecated call to `Use`. The `init` in `rules.go` registers only the eight names the earlier releases did, is deprecated, and `TestThePackageIsALibrary` refuses an `init` that does anything but `validation.Extend` |
| the rule name `"cpf"` written in a rule set | `br.RuleCPF`. A rule set in a package-level variable that names a rule `Use` declares is compiled with `validation.WithModuleRules()` and checked with `Factory.Check` at boot |
| a rule of the package that accepts what only one PHP package accepts, and a rule that forgives punctuation | one reading per rule, in `UPGRADE.md` and the README, and the value bare or in its one punctuation |
| trimming, re-punctuating, or accepting any separator | the one conventional mask or the bare value, nothing else. A value the person did not type is a value nobody can explain |
| a lookup at the Receita Federal or the Correios | nothing. `network = false`; a validator that calls out fails when the network does |
| `number.CurrencyFromCents` behind `number.WithLocale` | `BRL`, which reads no process-wide setting. The tests hold the two to the same string under `pt-BR`, so this is the same answer and not a second one |
| a real person's CPF or CNPJ in a test, a comment or a README | a document generated by the oracle in `tests/Unit/documents_test.go`, or the Receita Federal's published specimen `12.ABC.345/01DE-35`. The one place that holds other numbers is `tests/Unit/rule_examples_test.go`, which rewrites the specimens the two PHP packages publish in their own test suites, each written by a generator site and shared by every project that tests these rules; add to it only what a PHP package publishes |
| an import of the framework, `hesape/database` or `hesape/http` | the standard library, `hesape/validation` and `hesape/enum`. `TestThePackageIsALibrary` holds the list |

## The four properties

A change that breaks one of them is not merged, whatever else it improves.

1. **The oracle decides.** Every check is compared with an implementation
   written apart from the package, from the weight tables, on generated input,
   on every single-character substitution and under fuzzing. A change to the
   arithmetic that the oracle disagrees with is a bug in the change.
2. **A rule and its Parse function are one check.** `RuleCPF` accepts what
   `ParseCPF` accepts, through `ValidCPF`; nothing in `rules.go` or `catalog.go`
   checks a value itself. The `formato_` rules, which ask about the spelling
   alone, are the exception, and their checks are in `masked.go`.
3. **Formatting never validates and never fails.** `CPF`, `CNPJ`, `CEP` and
   `Mobile` return what they cannot punctuate unchanged; `BRL`, `FormatBRL`
   and `Date` are total. `RedactCPF` and `RedactCNPJ` return `""` for what
   they cannot read, and never their input: a redaction that fell back to the
   value would print what it exists to hide.
4. **`arandu.mod.toml` matches the code.** All four capabilities are `false`,
   and `tests/Unit/audit_test.go` compares that with what the code calls.

## Writing code

Everything in the source is in English: identifiers, doc comments, internal
comments, error messages, and the names and messages of tests. The names of the
federative units in `UF.Label` are data, and are the names as they are written.

Every exported symbol carries a doc comment, and the comment documents the
symbol and nothing else. Why a signature is what it is belongs there when it is
a fact about the code. A date, an issue number, a version in progress or the
name of another repository does not.

Tests go under `tests/Unit`, declaring `package unit_test`, and import the
package by its module path, which is what makes them see exactly what a caller
sees.
