# Contributing

Thank you for considering a contribution.

## Before you open a pull request

Run the four gates. They are what CI runs, and a pull request that fails one of
them fails for everybody who pulls after it.

```bash
export GOWORK=off
gofmt -l $(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go')
go build ./...
go vet ./...
go test -race ./...
```

The two filters on `gofmt` are not preference. `gofmt` is the only tool in the
chain that ignores build tags, so it reads files the compiler deliberately
excludes; `testdata/` holds fixtures that are invalid on purpose.

A change to a parser gets a fuzz run too, of the target that covers it:

```bash
go test -run='^$' -fuzz='^FuzzParseCNPJ$' -fuzztime=30s ./tests/Unit
```

## What a change has to keep true

1. **The oracle decides.** The tests carry a second implementation of every
   check, written from the Receita Federal's weight tables rather than from
   this package. A change the oracle disagrees with is not merged by changing
   the oracle.
2. **A rule and its Parse function are one check.** The rules in `rules.go`
   call `ValidCPF`, `ValidCNPJ` and `ValidCEP`; they check nothing of their own.
3. **What is accepted stays narrow.** The value bare, or in its one
   conventional punctuation. Widening it to another separator, to whitespace or
   to another spelling is a change to what every installer stores, and it needs
   an argument in the pull request.
4. **This stays a library.** No route, no table, no policy, no `Module`, no
   import beyond the standard library, `hesape/validation` and `hesape/enum`,
   and an `init` that only registers rules. `TestThePackageIsALibrary` holds
   all of it.
5. **`arandu.mod.toml` matches the code.** Adding an outbound call, a file
   access or a process means declaring it there in the same commit.

No real person's document goes into a test, a comment or a document. Generate
one with the oracle in `tests/Unit/documents_test.go`.

## Style

Everything in the source is in English: identifiers, doc comments, internal
comments, error messages, and the names and messages of tests. `pkg.go.dev`
publishes the doc comments, and its readers are users of this package.

Every exported symbol carries a doc comment saying what it does. The comment
documents the symbol and nothing else: no dates, no issue numbers, no
references to other repositories. Why a signature is what it is belongs there
when it is a fact about the code; anything else belongs in the pull request.

## Commits

One commit per logical change. Formatting is not mixed with behaviour, because
a diff that does both is a diff nobody reviews.

## Reporting a vulnerability

Do not open an issue. See [SECURITY.md](SECURITY.md).

## Licence

By contributing you agree that your contribution is licensed under the MIT
licence, the same as the rest of github.com/hyz-is/arandu-br.
