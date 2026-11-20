# Adding an MCP tool

[Back to the contribution guide](../CONTRIBUTING.md)

> [!IMPORTANT]
> These paths describe the repository before the official SDK migration and
> planned 2.0 refactor are complete. Check current code and open pull requests
> before starting.

## Before writing code

Open an issue and describe:

- The user workflow the tool supports
- Why existing tools are insufficient
- Example inputs and expected output
- Whether the tool reads, creates, updates, or deletes data
- Authentication, authorization, and data-exposure risks
- Expected response size and pagination behavior

An API endpoint alone is not a reason to add a tool. Wait for maintainer
agreement before implementing it.

## Choose the package

| Tool domain | Current location |
| --- | --- |
| Public Terraform Registry | `pkg/mcp-official/tools/registry/` |
| HCP Terraform or Terraform Enterprise | `pkg/mcp-official/tools/tfe/` |
| Terraform Search | `pkg/mcp-official/tools/search/` |

Use `github.com/modelcontextprotocol/go-sdk/mcp` for new tools. Do not add a new
Mark3Labs implementation unless maintainers request migration compatibility
work.

Read a nearby tool and its test before adding files. Follow its naming,
response, error, and registration patterns.

## Implement the contract

A tool normally has:

- A typed arguments structure with JSON field names
- A purpose-specific response structure
- A `*mcp.Tool` descriptor
- A typed handler
- Input validation before API access
- Accurate read-only, destructive, and open-world annotations

Keep the contract small and stable:

- Return only fields needed by the workflow, not a complete `go-tfe` model.
- Paginate or bound list responses.
- Encode empty lists as `[]`, not `null`.
- Do not expose tokens, signed URLs, state data, sensitive variables, or
  unnecessary customer data.
- Preserve existing names and JSON fields unless the approved issue changes
  them.

Reuse the shared clients, sessions, middleware, pagination, and API error
handling. Do not build a second client or authentication path inside a tool.

## Register the tool

Update each current registration point:

1. Add the implementation and adjacent unit test under
   `pkg/mcp-official/tools/`.
2. Register the typed handler in `pkg/mcp-official/tools/tools.go`.
3. Add its toolset and runtime gates to `pkg/toolsets/registry.go`.
4. Update registration or toolset tests if applicable.

Discuss write and destructive tools with maintainers. Use the existing
`RequiresTFE` and `RequiresTFOps` policies where applicable; do not classify a
tool from its HTTP method alone.

## Test and document it

- Follow the [unit test guide](unit-tests.md).
- Add [HCPT coverage](hcpt-tests.md) when correctness depends on real API
  behavior, permissions, entitlements, or mutations.
- Complete the [linting checks](linting.md).
- Exercise the contract with [MCP Inspector](mcp-inspector.md).
- Add a concise entry under `# Unreleased` in `CHANGELOG.md`.
- Update user documentation for new or changed public behavior.

## Definition of done

- [ ] The issue and user workflow are linked from the pull request.
- [ ] The tool uses the official MCP Go SDK.
- [ ] Input, output, annotations, and runtime gates match its behavior.
- [ ] Unit and registration tests pass without credentials or network access.
- [ ] Applicable HCPT and Inspector scenarios are recorded in the pull request.
- [ ] Logs and results contain no sensitive or unnecessary data.
- [ ] The changelog and user documentation are updated.
