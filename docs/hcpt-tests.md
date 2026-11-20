# HCPT tests

[Back to the contribution guide](../CONTRIBUTING.md)

Tests under `test/terraform/` call tools through a running MCP server and verify
behavior against real HCP Terraform or Terraform Enterprise APIs. They are not
unit tests.

> [!WARNING]
> Use only an approved test organization, a least-privileged test token, and
> disposable resources. Never point these tests at production.

## When to add HCPT coverage

Add or update an HCPT test when correctness depends on:

- Real API request and response behavior
- Token permissions or product entitlements
- Resource creation, updates, or deletion
- Behavior shared by the legacy and official SDK endpoints

A migration that preserves an existing tool contract can normally reuse that
tool's existing HCPT coverage.

## Required configuration

The test harness reads:

| Variable | Purpose |
| --- | --- |
| `TF_MCP_ENDPOINT` | Running MCP endpoint |
| `TFE_TOKEN` | Test account token |
| `TFE_ORG_NAME` | Isolated test organization |
| `TFE_USERNAME` | Expected test username |
| `TFE_USER_EMAIL` | Expected test email |
| `TFE_ADDRESS` | Terraform Enterprise address; omit for HCP Terraform |
| `ENABLE_TF_OPERATIONS` | Enables tests and tools for gated operations |

Keep tokens out of source files, command arguments, logs, screenshots, and pull
request descriptions. Load them from your approved secret store.

## Compile without running

Every contributor can verify that the suite compiles:

```shell
go test ./test/terraform/... -run '^$'
```

## Run against the official SDK endpoint

The following setup is temporary while `/mcp/official` and
`TF_X_OFFICIAL_SDK_ENABLED` exist.

Build and start the server in one terminal:

```shell
make build

TF_X_OFFICIAL_SDK_ENABLED=true \
ENABLE_TF_OPERATIONS=true \
  ./bin/terraform-mcp-server streamable-http \
  --transport-host 127.0.0.1 \
  --transport-port 8080
```

In another terminal, load test credentials from your secret store and set:

```shell
export TF_MCP_ENDPOINT=http://127.0.0.1:8080/mcp/official
export TFE_ORG_NAME=your-test-organization
export TFE_USERNAME=your-test-username
export TFE_USER_EMAIL=your-test-email
export ENABLE_TF_OPERATIONS=true
```

Then run the narrowest relevant test:

```shell
make test-hcpt RUN='^TestWorkspaceHappyPath$'
```

Use `/mcp` and omit `TF_X_OFFICIAL_SDK_ENABLED` when intentionally testing the
legacy endpoint. CI runs the suite against both endpoints; see
[`.github/workflows/hcpt_test.yml`](../.github/workflows/hcpt_test.yml).

Only set `ENABLE_TF_OPERATIONS=true` when testing operations that need the
gate. The server and test process must use the same setting.

## Test safely

- Generate unique resource names so concurrent runs do not collide.
- Register cleanup immediately after creating a resource.
- Check entitlements and skip with a clear reason when unavailable.
- Verify mutations through the API when the MCP result alone is insufficient.
- Keep polling bounded by a timeout.
- Do not hide a missing tool on the legacy endpoint; that is a regression.

If you cannot run an applicable HCPT test, say so in the pull request and name
the exact command maintainers should run.
