# `mail`

> **Needs update:** ported from an earlier application and due to be re-implemented by hand. The ported tests were removed, so this package currently has no tests.

`mail` sends email through a transport chosen explicitly with `MAIL_TRANSPORT`: `smtp`, `resend` or `log`. It only handles delivery. The content (templates, links) belongs to the module sending it.

## Quick Start

Build the sender once in the worker. It makes no network connection until `Send`. Send from a background job, never inside a request.

```go
sender, err := mail.New(cfg.Mail, log)

err = sender.Send(ctx, mail.Message{
    To:             mail.Address{Name: a.Name, Email: a.Email},
    Subject:        "...",
    Text:           text,
    HTML:           html,
    IdempotencyKey: "accounts/invite/" + tokenID.String(),
})
if mail.IsPermanent(err) {
    return river.JobCancel(err) // retrying won't help
}
```

From and Reply-To come from config.

## API Reference

* `New(cfg, log) (Sender, error)`: Builds the sender for the configured transport.
* `IsPermanent(err) bool`: True for Resend 4xx except 408 and 429, SMTP 5xx, and invalid addresses. Everything else retries.
* `IsReservedTestDomain(email) bool`: True for `.test`, `.example`, `.invalid` and `.localhost`.
* `StatusError`: An HTTP error from Resend, with `Status` and `Reason`. `Reason` holds only the `name` and `message` fields, truncated, with the API key redacted.

## Transports

| Transport | Notes |
| --- | --- |
| `smtp` | Opens one connection per send, with a 30s timeout. `MAIL_SMTP_TLS=none` is for Mailpit only |
| `resend` | Drops a repeat with the same `IdempotencyKey` for 24h, so a job retried after a successful send doesn't mail twice. `smtp` and `log` ignore the key |
| `log` | Prints the message, including the text body and any link, at WARN. For development and tests only |

Resend and SMTP with TLS drop recipients in reserved test domains and log a WARN with only the domain. `log` and Mailpit still show them.

## Configuration

| Env var | Default |
| --- | --- |
| `MAIL_TRANSPORT` | required: smtp, resend or log |
| `MAIL_FROM` | required: `Name <addr>` |
| `MAIL_REPLY_TO` | |
| `MAIL_SMTP_HOST` | |
| `MAIL_SMTP_PORT` | 587 |
| `MAIL_SMTP_USERNAME` | |
| `MAIL_SMTP_PASSWORD_FILE` | |
| `MAIL_SMTP_TLS` | starttls, tls or none |
| `MAIL_RESEND_API_KEY_FILE` | |
