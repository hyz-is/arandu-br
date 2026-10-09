# Security policy

## Supported versions

The latest minor release of Arandu BR receives security fixes. Older minor
releases do not: a published Go module version is immutable, so a fix is a new
release rather than a moved tag.

## Reporting a vulnerability

Report it privately, through GitHub's advisory form:

<https://github.com/hyz-is/arandu-br/security/advisories/new>

Do not open a public issue and do not describe the problem in a pull request.
A report that arrives in public is a report every reader of this repository can
act on before there is a release to upgrade to.

Include what you need to reproduce it: the version and the input that triggers
it. Never a real person's document number: a generated one shows the same
defect.

Expect an acknowledgement within a few days. If the report is confirmed, the
fix, the release and the advisory are published together, and you are credited
unless you ask not to be.

## What is in scope

Anything that makes a value pass that is not what the package says it accepts,
or makes a call fail in a way its signature does not allow. In particular:

- a value `ParseCPF`, `ParseCNPJ` or `ParseCEP` accepts that is not a CPF, a
  CNPJ or a CEP in the shape the documentation names -- a character outside the
  mask, a check digit that does not compute, a repetition the package refuses;
- a rule registered with `hesape/validation` that accepts what the matching
  `Parse` function refuses;
- a panic on any input, from any exported function;
- `Value` writing a `UF` outside the set to a column.

## What is not

- A document that passes and was never issued, or a CEP that passes and
  addresses nothing. The package checks shape and check digits, and says so;
  existence is a question for the Receita Federal and the Correios.
- A vulnerability in Hesape itself. Report that to Hesape.
