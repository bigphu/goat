# `apicommon`

> **Needs update:** ported from storeit. The source spec moved to `api/common.yaml` at the module root, and regenerating from there gives the same output. The comments in `common.yaml` still describe storeit's paths and `make generate`.

`apicommon` holds the Go types generated from [`api/common.yaml`](../../api/common.yaml): errors, paging and IDs. Modules' generated code shares them through `import-mapping`. `apicommon.gen.go` is generated, so don't edit it by hand.

## Types

| Type | Purpose |
| --- | --- |
| `Problem` | The problem+json error body |
| `FieldError` | One entry of `Problem.Errors` |
| `ID` | A UUID |
| `Page`, `PageSize` | Paging params |
| `Paged` | The `total` count, combined with each list's `items` through `allOf` |

`page` is 1 to 100000 (default 1) and `page_size` is 1 to 200 (default 50), so `(page-1)*page_size` always fits an `int32` OFFSET. The validator doesn't apply these defaults, so handle them in Go.

`problem_test.go` checks that the body `web.RenderProblem` writes fits `Problem`, so the spec and the code can't drift apart.

## Regenerating

```bash
go generate ./web/apicommon
```
