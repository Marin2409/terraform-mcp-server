# Testing with MCP Inspector

[Back to the contribution guide](../CONTRIBUTING.md)

[MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is the
reference interactive client for checking MCP registration, schemas, calls,
and protocol messages. It supplements automated tests; it does not replace
them.

## Before you begin

Build the current server:

```shell
make build
```

Public Registry tools need no credentials. For an authenticated tool, export a
least-privileged test token from an approved secret store. Do not put tokens in
commands, Inspector screenshots, or pull request descriptions.

## Use the web interface

During the official SDK migration, launch Inspector against the official stdio
server with:

```shell
TF_X_OFFICIAL_SDK_ENABLED=true \
  npx @modelcontextprotocol/inspector \
  ./bin/terraform-mcp-server stdio
```

Inspector prints a local URL and normally opens it in a browser. In the web
interface:

1. Connect to the server.
2. Open **Tools** and refresh the tool list.
3. Select the changed tool.
4. Review its name, description, input schema, and annotations.
5. Enter representative arguments and select **Run Tool**.
6. Review the result and the JSON-RPC exchange.
7. Check the console for server errors or sensitive output.

Repeat the call with missing, invalid, and no-result inputs. For write tools,
use a disposable resource and verify the effect in the test organization.

## Use the CLI

Check registration and schema portability:

```shell
TF_X_OFFICIAL_SDK_ENABLED=true \
  npx @modelcontextprotocol/inspector --cli \
  ./bin/terraform-mcp-server stdio \
  --method tools/list \
  --strict
```

Call one tool:

```shell
TF_X_OFFICIAL_SDK_ENABLED=true \
  npx @modelcontextprotocol/inspector --cli \
  ./bin/terraform-mcp-server stdio \
  --method tools/call \
  --tool-name your_tool_name \
  --tool-args-json '{"example_argument":"example-value"}'
```

Use `--format json` when you need machine-readable output.

## Test Streamable HTTP

Start the server in one terminal:

```shell
TF_X_OFFICIAL_SDK_ENABLED=true \
  ./bin/terraform-mcp-server streamable-http \
  --transport-host 127.0.0.1 \
  --transport-port 8080
```

While the migration is active, connect Inspector to the official endpoint:

```shell
npx @modelcontextprotocol/inspector \
  --server-url http://127.0.0.1:8080/mcp/official \
  --transport http
```

Use `/mcp` when intentionally testing the legacy endpoint. The separate
`/mcp/official` endpoint and `TF_X_OFFICIAL_SDK_ENABLED` flag are temporary.

## What to verify

- The server initializes without errors.
- The tool appears once in `tools/list`.
- Its description and schema tell an agent how to call it.
- Required fields and annotations match actual behavior.
- Valid, invalid, and empty-result calls return stable shapes.
- API errors are useful and do not expose backend response data.
- Results and logs contain no credentials or unnecessary customer data.

Record the scenarios and sanitized results in the pull request. Do not include
real organization names, workspace names, tokens, or customer data.
