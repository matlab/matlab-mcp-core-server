# Registry Publication

Files in this directory drive publication of this server to the
[official MCP registry](https://github.com/modelcontextprotocol/registry).

## Files

- `server.template.json`: source of truth for the registry entry. The
  `{{.Version}}`, `{{.FileSHA256}}`, and `{{.MCPBFilename}}` placeholders
  are expanded at publish time and at validation time with dummy values.
- `server.schema.json`: vendored copy of the registry schema referenced by
  the template's `$schema` field. The file lives at
  `internal/adaptors/registryvalidate/server.schema.json` and is embedded
  into the `registry-validate` binary. Refresh it periodically and bump the
  `$schema` URL in `server.template.json` to match.

The publish workflow generates a `server.json` on disk in CI, but
`make registry-validate` performs expansion in memory only.

## Local validation

```
make registry-validate
```

Runs offline schema validation. Does not contact the registry.

Registry publishing currently uses the manual
`.github/workflows/publish-registry.yml` workflow. Run it against the release
tag after the release is published. Automatic release publishing can be enabled
after the manual workflow is verified.

---

Copyright 2026 The MathWorks, Inc.

---
