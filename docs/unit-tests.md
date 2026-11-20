# Unit tests

[Back to the contribution guide](../CONTRIBUTING.md)

Unit tests must be fast, repeatable, and independent of live Terraform
services. Put them next to the behavior they cover.

## What to test

For an MCP tool, cover:

- Tool name, description, and annotations
- Required inputs and schema constraints
- Output schema and empty-result shape
- Blank, invalid, and conflicting inputs
- Response mapping and pure helper behavior
- The regression demonstrated by a bug fix

Registration tests belong in `pkg/mcp-official/tools/tools_test.go` or the
nearest toolset test. Shared client, credential, session, and middleware
behavior belongs in the package that owns it, not in an individual tool test.

## Keep unit tests local

Unit tests must not require credentials or make external network calls. Use:

- `httptest.Server` for HTTP backend behavior
- MCP in-memory transports for server registration and calls
- `t.TempDir` and `t.Setenv` for isolated filesystem or environment behavior

Do not weaken an assertion or skip a test only to make a change pass.

## Run tests

Run the changed package while developing:

```shell
go test ./path/to/changed/package
```

Run the non-live application packages before review:

```shell
go test ./cmd/... ./pkg/... ./version/...
```

Run the CI unit-test target:

```shell
make test
```

> [!CAUTION]
> `make test` runs `go test -v ./...`. If HCPT credentials are exported, tests
> under `test/terraform/` may contact a real service. Unset those credentials
> unless you intend to run live tests against an isolated test organization.

For concurrency changes, also run the race detector on the affected package:

```shell
go test -race ./path/to/changed/package
```

Docker-backed public Registry and transport tests are documented in
[`e2e/README.md`](../e2e/README.md).
