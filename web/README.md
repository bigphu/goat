# `web`

> **Needs update:** ported from storeit. It now uses goat's `errs`, which doesn't marshal to JSON, so `problem.go` builds the response body itself. The 500 title changed to `Internal Server Error`. `ErrValidation` keeps the type `/errors/validation-failed`, which differs from `errs.Validation`'s `/errors/validation`. The Swagger page title in `docs.go` is still `StoreIt API`.

`web` is the shared HTTP layer: problem+json errors, OpenAPI request validation, and JSON helpers for hand-written handlers.

## Quick Start

### Hand-Written Handlers

A handler returns an error, and `Handle` turns it into the response. `Decode` reads the JSON body and validates it against `validate:` tags.

```go
r.Post("/things", web.Handle(func(w http.ResponseWriter, r *http.Request) error {
    var in CreateThing
    if err := web.Decode(r, &in); err != nil {
        return err
    }
    return web.JSON(w, http.StatusCreated, out)
}))
```

### Generated Handlers (oapi-codegen Strict Server)

Your `oapi.yaml` needs `strict-server: true`, `chi-server: true` and `embedded-spec: true`, and an `import-mapping` from your `$ref` to [`api/common.yaml`](../api/common.yaml) onto `goat/web/apicommon`.

```go
func (h *Handler) Mount(r chi.Router, tokens *jwt.Provider) error {
    spec, err := api.GetSwagger()
    if err != nil {
        return err
    }
    const base = "/api/v1" // must match `servers` in the spec, or ValidateRequests panics

    authMW := auth.Middleware(tokens,
        auth.Public(web.MustOperations(spec, base, "POST /api/v1/auth/login")...))
    upload := middleware.ForOperations(web.MustOperations(spec, base, "POST /api/v1/imports"),
        middleware.BodyLimit(50<<20), middleware.ReadTimeout(10*time.Minute))

    strict := api.NewStrictHandlerWithOptions(h, nil, api.StrictHTTPServerOptions{
        RequestErrorHandlerFunc:  web.RequestError,
        ResponseErrorHandlerFunc: web.WriteProblem,
    })
    api.HandlerWithOptions(strict, api.ChiServerOptions{
        BaseRouter: r,
        BaseURL:    base,
        // The last entry runs first: auth, then the overrides, then validation.
        Middlewares:      []api.MiddlewareFunc{web.ValidateRequests(spec, base), upload, authMW},
        ErrorHandlerFunc: web.RequestError,
    })
    return nil
}
```

The generated wrapper parses typed path and query params **before** these middlewares, and its errors go to `RequestError`.

## Error Responses

Any error becomes `application/problem+json`. An error with no `*errs.Error` in its chain becomes a generic 500, and its message is logged but never sent. Causes are never in the body.

```json
{
  "type": "/errors/validation-failed",
  "title": "Validation failed",
  "status": 422,
  "detail": "The request body has 1 invalid field.",
  "errors": [{ "field": "email", "detail": "is required" }]
}
```

Each entry in `errors` has `errs.ErrRequired` or `errs.ErrInvalid` as its cause, so `errors.Is` works on errors from `Validate`, `Decode` and `ValidateRequests`.

| Variable | Status | Type |
| --- | --- | --- |
| `ErrMalformedJSON` | 400 | `/errors/malformed-json` |
| `ErrInvalidRequest` | 400 | `/errors/invalid-request` |
| `ErrRouteNotFound` | 404 | `/errors/route-not-found` |
| `ErrMethodNotAllowed` | 405 | `/errors/method-not-allowed` |
| `ErrRequestTooLarge` | 413 | `/errors/request-too-large` |
| `ErrUnsupportedMediaType` | 415 | `/errors/unsupported-media-type` |
| `ErrValidation` | 422 | `/errors/validation-failed` |

## API Reference

### Validation

`ValidateRequests(spec, baseURL)` validates params, `Content-Type` and the body against the spec:

| Problem | Response |
| --- | --- |
| Invalid params | 400 `/errors/invalid-request`, with `errors[]` |
| Body fails the schema | 422 `/errors/validation-failed`, with `errors[]` (paths like `address.city`) |
| Malformed JSON | 400 |
| Wrong or missing `Content-Type` | 415, listing the allowed types |
| Body over the limit | 413 |
| Route not in the spec | 500, because it is a misconfiguration |

* Only JSON bodies that kin-openapi can decode are read and validated. Other types, such as uploads and custom `+json`, get a `Content-Type` check and a required-body check, and the body stays unread for streaming.
* Spec `security` is ignored, because the auth middleware handles it.
* Schema **defaults are not applied**, so apply them in Go. `page` and `page_size` arrive as nil.
* `Content-Type` is lowercased before the checks and passed on that way.
* `format: uuid` accepts any UUID version.

### Functions

* `MustOperations(spec, baseURL, ops...)`: Panics if an operation isn't in the spec. Operations are `"METHOD <chi route pattern>"`: uppercase method, the full path including the base URL, and param names exactly as in the spec.
* `RequestError(w, r, err)`: For oapi-codegen's `ErrorHandlerFunc` and `RequestErrorHandlerFunc`. Param errors become 400 with a field, body read errors map like `Decode`, and anything else returns a fixed message.
* `WriteProblem(w, r, err)`: Writes the problem and logs 5xx errors, plus anything after the response started, with the context logger.
* `RenderProblem(w, r, err)`: Writes without logging. Both write nothing once the response has started, drop headers such as `Content-Length` and `Content-Disposition`, and treat a nil error as 500.
* `Handle(h)`, `Decode(r, &dst)`, `Validate(v)`, `JSON(w, status, v)`, `NoContent(w)`, `Text(w, status, body)`: Helpers for hand-written handlers.
* `NotFoundHandler()`, `MethodNotAllowedHandler()`: Installed by `server.New`. The second sets `Allow`.
* `WrapWriter(w, r)`: Wraps `w` so `web` can tell whether the response has started.
* `MountDocs(r, docs...)`: Swagger UI at `GET /api/docs` and each spec at `GET /api/docs/<name>.json`. The page is public and loads Swagger UI from jsDelivr, so mount it in development only.
