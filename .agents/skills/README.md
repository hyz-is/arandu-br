# Skills

Procedures an assistant follows when working on this package.

They live in `.agents/skills/<name>/SKILL.md`, which is the path the coding
assistants read from — Cursor, Codex, Cline, Copilot, Gemini CLI, Amp, OpenCode,
Warp, Zed and the rest all look there. It is one directory rather than a file
per vendor, so a skill written once is read by whatever this package is being
written with.

Each file opens with frontmatter carrying a `name` and a `description`. The
`name` has to equal the directory name exactly, or the skill is not loaded. The
description is what a tool reads to decide whether the skill is relevant, so it
names the situation you are in rather than the topic it covers.

| skill | when it fires |
| --- | --- |
| `br-release` | the gates, the fuzz runs, the manifest, a dependency, a version, a tag |
| `br-vault-notes` | writing the note, the gap or the journal entry, when this checkout sits inside the Arandu Obsidian vault |
| `br-package` | validating or formatting a CPF, a CNPJ, a CEP, reais, a date or a UF **in an application** |

The template this package started from carried a `-policy` and a `-module`
skill. Both taught the HTTP module the template ships, and this package removed
it: there is no policy to open and no route, handler or migration to add.

`br-package` has a different audience from the other two, and that is on
purpose: it travels with the package so that an assistant working in somebody
else's project — the one running `go get` — has the rule constants, the accepted
shapes and the symptoms in front of it instead of writing a CPF check of its
own.

Its frontmatter says so, under `metadata`, with `audience: app`. That line is
what `aru skills:sync` reads to copy the skill into an application whose
`go.mod` requires this package, and no other skill here carries it: an
application that received the release procedure would follow it.

`br-vault-notes` fires on a condition rather than on a task: it applies
only when `MOC-arandu.md`, `plans/cmd/audit-vault/` and `45-modules/` are actually
beside this checkout. Outside the vault it is inert, and it says so first, so a
package cloned somewhere else never grows a folder tree imitating one.

## Why these exist

The common failure modes when changing this package are a check widened to
accept what a person did not type, a rule in `rules.go` that checks a value
itself instead of calling the Parse function, an `init` that does more than
register rules, a `Register` function an application is told to call from
`main` — after its rule sets were compiled — and a real person's document
pasted into a test. None belongs here.

The package is built to be checked rather than trusted. Every check is compared
with an oracle written apart from it, on generated documents, on every
single-character substitution and under fuzzing, and the structure test holds
the library shape. Running `go test -race ./...` exercises all of it.

## Adding your own

A skill in this directory is yours and travels with the repository. Keep it a
procedure rather than a description: a file that says "read the documentation"
never changes what anybody does.
