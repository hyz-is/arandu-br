---
name: br-release
description: Get this Arandu package green, declared and published. Use when the request is to "run the gates", "check it builds", "the tests fail", "run the fuzzers", "go build says the directory prefix does not contain modules listed in go.work", "aru doctor says this is not an Arandu project", "update arandu.mod.toml", "add a dependency", "bump hesape", "tag a release", "publish it", "go get cannot find the module", "write the changelog", or before opening a pull request. Covers the four gates and the fuzz runs, what the manifest declares and who checks it, why the module requires Hesape and nothing else, the immutability of a published module version, and what breaks somebody else's build without touching a signature.
license: MIT
---

# Getting it green and out

## The four gates

Nothing is finished until all four exit zero.

```sh
export GOWORK=off
gofmt -l $(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go')
go build ./...
go vet ./...
go test -race ./...
```

`GOWORK=off` first, and here it is not a habit. This checkout may sit beside a
Go workspace that lists the framework repositories and does not list this one.
When it does, `go build ./...` exits 1 before compiling anything:

```
pattern ./...: directory prefix . does not contain modules listed in go.work
or their selected dependencies
```

That is the workspace refusing to answer for a module it was not told about, and
it says nothing about the code. With the workspace off, the module resolves the
versions in `go.mod` — which is what CI compiles against, and what somebody's
`go get` will get.

`aru doctor` is not a gate. It exits 1 here with *"this is not an Arandu
project: no go.mod, main.go and arandu.toml together"* — it reads applications,
and this is a library.

## The fuzz runs

`go test` runs the seeds of the four fuzz targets every time. Before a release,
and after any change to a parser, run each one for a while; `go test -fuzz`
takes one target at a time:

```sh
go test -run='^$' -fuzz='^FuzzParseCPF$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseCNPJ$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzParseCEP$' -fuzztime=30s ./tests/Unit
go test -run='^$' -fuzz='^FuzzBRLAgreesWithHesapeUnderPtBR$' -fuzztime=30s ./tests/Unit
```

A failure writes its input under `tests/Unit/testdata/fuzz/`. Commit that file
with the fix: it becomes a seed, and the input runs on every `go test` after.
Check first that it is not a real person's document — a fuzzer that reaches a
valid CPF reached it by arithmetic, but the file is published.

## What the manifest declares, and who reads it

`arandu.mod.toml` is a promise, and the check that gives it weight runs here.
Nowhere else does: `aru doctor` audits the tree of the application it is run
inside and never loads a dependency, so the manifest of an installed package is
opened by nothing but that package's own suite.

```toml
name = "hyz-is/br"
framework = ">= 0.51"
profiles = ["conventional", "performance"]

[permissions]
network = false
filesystem = false
exec = false
migrations = false
```

`TestTheDeclaredCapabilitiesAreWhatTheCodeDoes`, in `tests/Unit/audit_test.go`,
compares the four permissions against what the code *calls*, by AST. All four
are `false`, and they stay that way: the package turns strings into strings. A
lookup at the Receita Federal or the Correios would make it `network = true`,
and a validator that fails when the network does — that is a different package,
not a change to this one.

The `framework` line is the first Framework release built on the Hesape this
module requires, because the framework itself is not in `go.mod`. When Hesape
is bumped, find the Framework release that requires the new version and move
the floor with it, in the same commit. An application on an older Framework
would otherwise have its Hesape raised past the one it was built with.

## Adding a dependency

Do not. This module goes into other people's builds, and it requires
`github.com/arandu-io/hesape` and nothing else — not even the framework, which
is the web stack and is not needed to format a string.
`TestTheModuleRequiresNothingButHesape` refuses any other direct requirement,
and `TestThePackageIsALibrary` refuses an import beyond the standard library,
`hesape/validation` and `hesape/enum`.

```sh
export GOWORK=off
go list -m -f '{{if and (not .Indirect) (not .Main)}}{{.Path}} {{.Version}}{{end}}' all
```

## Publishing

A published Go module version is **immutable**. The proxy serves it forever, and
a mistake is corrected by another release, never by moving a tag.

Before tagging:

- the four gates pass with `GOWORK=off`, and each fuzz target has run;
- `arandu.mod.toml` matches the code;
- `CHANGELOG.md` has the entry under the version, in the Keep a Changelog
  sections — `Added`, `Changed`, `Fixed`, `Removed` — and `UPGRADE.md` has a
  `## vX.Y.Z` heading for the same version;
- `currentVersion` in `tests/Unit/release_test.go` names that version;
- every rule constant in `rules.go` is named in a released changelog entry.

While the version starts with `v0.`, breaking is allowed. Breaking *quietly* is
not. Run `apidiff` against the release baseline; the CI step named `api diff
against the last release` fails unless `UPGRADE.md` names each removal.

Four changes break an installer without touching a signature, and each one is
worth a line in the changelog and in `UPGRADE.md`:

- **A rule name.** A rule set that spells the old name refuses to boot.
- **What a parser accepts.** Narrowing it refuses values an application
  already stored; widening it lets a form accept what it used to refuse. Either
  way, every installer's data changes meaning.
- **The bare spelling a parser returns.** It is what applications store and
  search for, so changing it splits their columns in two.
- **A UF code or a constant name.** Codes are in columns; renaming a constant
  must never rewrite a row.

## Reporting a vulnerability

Not in an issue and not in a pull request. `SECURITY.md` has the private
advisory address. A value accepted that the documentation says is refused is in
scope, and so is a panic on any input.
