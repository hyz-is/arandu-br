# Upgrade Guide

## Unreleased

### Published views move out of `vendor/`

The view a package publishes lands in `resources/views/modules/<slug>/` and
compiles to `storage/framework/views/modules/<slug>`. It used to be `vendor/` in
both, and that address could not work: the go command reserves the name twice,
and a tree of published views hit both rules.

A file under a directory named `vendor` is left out of the module zip at any
depth. The file stays in the package's repository and is missing for everyone
who downloads it, so the `go:embed` that names its directory matches nothing and
the person building the project reads

```
pattern resources/views: no matching files found
```

— an error about the package, raised in their project. And a package whose
import path carries the element cannot be imported at all:

```
bootstrap/app.go:98:2: use of vendored package not allowed
```

which is exactly the import `(*Module).Boot` asks for. A published view is
compiled into a Go package the application has to import for its `init()` to
register anything, so the second rule refused the last step of the install.

Both were reproduced before this changed: `zip.CheckDir` reports the view as
`file is in vendor directory`, and a package under
`storage/framework/views/vendor/<slug>` is refused at import.

To move a package already released:

1. `git mv resources/views/vendor resources/views/modules`.
2. Rename the `vendorDir` constant in `views.go` to `moduleDir`, with the value
   `modules`.
3. Release the package, and tell the projects that installed it to publish
   again. The old files are theirs now, so `aru vendor:publish --apply` writes
   the new tree beside the old one and the old one is deleted by hand, along
   with its lines in `vendor-publish.lock` and its import in `bootstrap/app.go`.

Framework `v0.46.4` and Hesape `v0.37.0` refuse a publication that carries the
reserved name, so a package that has not moved fails its own tests with a
message naming both rules, rather than failing in the first project that
installs it.

### The entity is a concrete type over the non-generic model

Hesape `v0.47.0` removed the generic model layer, and the package now requires
Hesape `v0.50.1` and Framework `v0.51.0`; `arandu.mod.toml` says
`framework = ">= 0.51"`. `Br` embeds `model.Model`, its table is declared
once as `brTable` in `model.go`, and `aru model:build` generates
`BrQuery.go` beside it: `Brs`, `BrQuery` and
`BrCollection`. The file is committed, and `aru model:build --check` exits
1 when it is stale.

An application that uses the package through `New`, its routes and
`BrService` changes nothing for this section: every service method keeps
its signature. What breaks is code that reached the table through `Brs`:

| before | now |
|---|---|
| `func Brs(*data.DB) *model.Model[Br]` | `func Brs(model.DB) *BrQuery`; a `*data.DB` is a `model.DB`, so the call compiles unchanged |
| `Brs(db).NewQuery().Where(…)` | `Brs(db).Where(…)` |
| `Brs(db).NewInstance(nil, false)`, then `.Entity` | `Brs(db).New()`, which returns the `*Br` |
| `func(*model.Builder[Br])` in a grouped where | `func(*BrQuery)` |
| `record.Exists`, a field | `record.Exists()` |
| `record.Table`, a field | `record.Table()`, the `*model.Table` |
| the configuration the generic model carried as fields of `Br` -- `ConnectionName`, `CreatedAtColumn`, `DeletedAtColumn`, `Entity`, `Grammar`, `Incrementing`, `KeyType`, `NamedScopes`, `PerPage`, `PrimaryKey`, `Processor`, `RelationResolvers`, `SoftDeletes`, `TenantColumn`, `Timestamps`, `UpdatedAtColumn`, `WasRecentlyCreated` | gone; the table settings live in `brTable`, and a row keeps `Save`, `Delete`, `Fresh`, `Replicate`, `Exists`, `Table`, `WasRecentlyCreated()` and the attribute methods of `model.Model` |

Every one of those is `Br.<field>` in the API diff: `Br.ConnectionName`,
`Br.CreatedAtColumn`, `Br.DeletedAtColumn`, `Br.Entity`,
`Br.Exists`, `Br.Grammar`, `Br.Incrementing`,
`Br.KeyType`, `Br.NamedScopes`, `Br.PerPage`,
`Br.PrimaryKey`, `Br.Processor`, `Br.RelationResolvers`,
`Br.SoftDeletes`, `Br.Table`, `Br.TenantColumn`,
`Br.Timestamps`, `Br.UpdatedAtColumn` and
`Br.WasRecentlyCreated`.

**What changes without a compiler error.** A value copy of a `Br` cannot
be saved: the write is refused with `model.ErrUnwired`, where before it acted on
the row the copy was taken from. And `Brs(db)` is now one mutable query
rather than a model that opened a new one per chain, so two queries begun from
one value held in a variable share their clauses: start each at the constructor.

To move a package already configured from this template:

1. `go get github.com/arandu-io/hesape@v0.50.1 github.com/arandu-io/framework@v0.51.0`,
   then `go mod tidy`, and raise the `framework` floor in `arandu.mod.toml` to
   `>= 0.51`.
2. `go run github.com/arandu-io/aru/cmd/model-upgrade@v0.60.5 --dry-run ./...`,
   then without `--dry-run`. It refuses a constructor held in a variable and
   used twice; start each query at the constructor and run it again.
3. `go run github.com/arandu-io/aru@v0.60.5 model:build`, and commit the
   generated `BrQuery.go`.
4. Let the compiler name what is left; tests that read the model's settings
   observe what the table compiles instead.

<details>
<summary>Every method the generic model promoted on <code>*Br</code> that the non-generic one does not, as the API diff names them</summary>

The five whose signature changed rather than disappeared:

- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Is`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).MakeHidden`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).MakeVisible`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetRelation`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SyncOriginal`

The rest are gone from the method set of `*Br`:

- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).AddGlobalScope, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).All, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Append, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).AttributesToArray, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).CallNamedScope, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Create, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Destroy, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).DiscardChanges, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Except, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Find, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).FindMany, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).FindOrFail, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).FindOrNew, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).First, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).FirstOrCreate, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).FirstOrNew, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ForceCreate, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ForceDeleteQuietly, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ForceDeleted, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ForceDeleting, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ForceDestroy, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).FreshTimestamp, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetAppends, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetConnectionName, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetCreatedAtColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetDeletedAtColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetForeignKey, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetGlobalScopes, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetHidden, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetIncrementing, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetKeyName, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetKeyType, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetMorphClass, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetPerPage, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetPrevious, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQualifiedCreatedAtColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQualifiedDeletedAtColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQualifiedKeyName, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQualifiedUpdatedAtColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQueueableConnection, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQueueableID, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetQueueableRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetRawOriginal, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetRelation, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetRouteKey, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetRouteKeyName, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetTable, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetTouchedRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetUpdatedAtColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).GetVisible, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).HasAppended, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).HasGlobalScope, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).HasNamedScope, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).IsForceDeleting, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).IsIgnoringTouch, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).IsNot, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).IsRelation, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).IsSoftDeletable, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadAggregate, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorph, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorphAggregate, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorphAvg, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorphCount, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorphMax, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorphMin, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).LoadMorphSum, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewBaseQueryBuilder, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewCollection, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewFromBuilder, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewInstance, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewModelQuery, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewQuery, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewQueryForRestoration, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewQueryWithoutRelationships, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewQueryWithoutScope, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewQueryWithoutScopes, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).NewTypedBuilder, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).On, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).OnWriteConnection, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Only, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).OnlyTrashed, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).OriginalIsEquivalent, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).PushQuietly, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).QualifyColumn, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).QualifyColumns, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Query, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Ref, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).RegisterGlobalScopes, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).RegisterModelEvent, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ReplicateQuietly, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ResolveRouteBinding, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ResolveRouteBindingQuery, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ResolveSoftDeletableRouteBinding, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).RestoreQuietly, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Restored, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Restoring, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetAppends, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetConnection, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetHidden, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetIncrementing, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetKeyName, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetKeyType, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetPerPage, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetTable, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetTouchedRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SetVisible, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SoftDeleted, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SyncChanges, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SyncOriginalAttribute, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).SyncOriginalAttributes, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).ToPrettyJSON, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Touches, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UnsetAttribute, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UnsetRelation, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UnsetRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UpdateOrCreate, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UpdateOrFail, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UpdateQuietly, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UpdateTimestamps, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).UsesTimestamps, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).Where, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).WhereKey, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).With, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).WithTrashed, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).WithoutRelations, method set of *Br`
- `github.com/arandu-io/hesape/database/model.(*Model[github.com/hyz-is/arandu-br.Br]).WithoutTimestamps, method set of *Br`

</details>

### Routes need a session, and the router answers the errors

Every route is registered behind `RequireAuth`, and each handler reads who is
asking with `ctx.User()` and returns what the service returned. The module no
longer loads the session itself, reads no guest, and maps no error to a status.

`Config.Tenant` was removed. It was the tenant a visitor with no session was
read as, and with every route behind the guard there is no such visitor:

```go
// Before.
br.New(br.Config{Tenant: cfg.Auth.Tenant}, db, sessions)

// After.
br.New(br.Config{}, db, sessions)
```

What a client sees changes with it, and none of it is a compiler error:

| request | before | now |
|---|---|---|
| no session | 403, after the policy refused a guest | 303 to `/auth/login`, before any handler runs |
| input rejected | 422, with the field names in the body | 303 back to the referring page, with the messages and what was typed in the flash |
| policy refused | 403 `forbidden` | 403 `Forbidden`, the standard sentence |
| no such record | 404 `not found` | 404 `Not Found` |

`ErrNotFound` now wraps `model.ErrModelNotFound`, which is what makes the router
answer it with 404. `errors.Is(err, ErrNotFound)` still holds; its text is
longer.

To move a package already configured from this template: replace each
handler's `m.subject(ctx.Request)` with `who, _ := ctx.User()` and each
`return m.answer(ctx, err)` with `return err`; delete `subject` and `answer`;
open `Routes` with `r = r.Group("", middleware.RequireAuth(m.sessions))`; make
`ErrNotFound` wrap `model.ErrModelNotFound`; and remove `Tenant` from `Config`
and from the wiring in `bootstrap/app.go`.

The table is what a page sees. A client that asks for JSON -- an `Accept`
naming `application/json`, or an XHR that is not htmx -- is not redirected after
a rejected input: Framework `v0.51.0` answers it 422 `application/problem+json`,
with the messages in `errors` keyed by field, and answers the 403 and the 404 as
problem documents too. Only the sign-in redirect is the same for every client.

