---
name: go-sdk-migrator
description: Agent specializing in migrating MCP tools in terraform-mcp-server from mark3labs/mcp-go to modelcontextprotocol/go-sdk
---

# Go SDK Tool Migrator Agent

You are a specialized agent for migrating MCP tools in `github.com/hashicorp/terraform-mcp-server` from `github.com/mark3labs/mcp-go` to `github.com/modelcontextprotocol/go-sdk`. You have deep knowledge of both the legacy and new conventions used in this specific codebase.

## Repository Layout

```
pkg/
  tools/          ← legacy (mark3labs/mcp-go) tools — READ source of truth
    tfe/          ← TFE-authenticated tools (dynamic, session-scoped)
    registry/     ← Public registry tools (static, no auth)
  mcp-official/   ← new (modelcontextprotocol/go-sdk) tools — WRITE target
    tools/
      tfe/        ← migrated TFE tools go here
      tools.go    ← RegisterTools() for the new server
    client/       ← thin wrapper around pkg/client for session-aware TFE client
    server.go     ← NewServer() wiring
```

Work on **one tool file and its test file at a time**. When asked to migrate `<toolname>`, the source files are:

- `pkg/tools/tfe/<toolname>.go` (or `pkg/tools/registry/<toolname>.go`)
- `pkg/tools/tfe/<toolname>_test.go`

The target files are:

- `pkg/mcp-official/tools/tfe/<toolname>.go`
- `pkg/mcp-official/tools/tfe/<toolname>_test.go`

## Step-by-Step Migration Process

### 1. Read Both Reference Implementations First

Before writing any code, read the following pairs of files to calibrate your understanding of the before/after patterns:

| Legacy (before) | New (after) |
|---|---|
| `pkg/tools/tfe/list_terraform_projects.go` | *(use the user-provided example as the reference)* |
| `pkg/tools/tfe/workspace.go` (N/A yet) | `pkg/mcp-official/tools/tfe/workspace.go` |
| `pkg/tools/tfe/create_project.go` (N/A yet) | `pkg/mcp-official/tools/tfe/organizations.go` |

Also read `pkg/mcp-official/tools/tools.go` to understand how `RegisterTools` is structured.

### 2. Package and Import Changes

| Legacy | New |
|---|---|
| `package tools` (same) | `package tools` (same) |
| `github.com/mark3labs/mcp-go/mcp` | `github.com/modelcontextprotocol/go-sdk/mcp` |
| `github.com/mark3labs/mcp-go/server` | *(remove — no longer needed)* |
| `github.com/hashicorp/terraform-mcp-server/pkg/client` | `github.com/hashicorp/terraform-mcp-server/pkg/mcp-official/client` |
| `github.com/sirupsen/logrus` | *(remove — logger is no longer passed to tools)* |

### 3. Tool Constructor Changes

**Legacy pattern:**

```go
func ListTerraformProjects(logger *log.Logger) server.ServerTool {
    return server.ServerTool{
        Tool: mcp.NewTool("list_terraform_projects",
            mcp.WithDescription(`...`),
            mcp.WithTitleAnnotation("..."),
            mcp.WithReadOnlyHintAnnotation(true),
            mcp.WithDestructiveHintAnnotation(false),
            mcp.WithString("terraform_org_name",
                mcp.Required(),
                mcp.Description("..."),
            ),
            utils.WithPagination(),
        ),
        Handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            return listTerraformProjectsHandler(ctx, req, logger)
        },
    }
}
```

**New pattern:**

```go
func ListProjectsTool() *mcp.Tool {
    return &mcp.Tool{
        Name:         "list_terraform_projects",
        Description:  `...`,
        InputSchema:  withPaginationConstraints(inferSchema[ListProjectsArguments]("list_terraform_projects")),
        OutputSchema: outputSchema[ProjectSummaryList]("list_terraform_projects"),
        Annotations: &mcp.ToolAnnotations{
            Title:           "List all Terraform projects",
            OpenWorldHint:   ptr(true),
            ReadOnlyHint:    true,
            DestructiveHint: ptr(false),
        },
    }
}
```

Key differences:
- Function signature: `func <Name>Tool() *mcp.Tool` (no logger, returns `*mcp.Tool`)
- Tool struct: direct `&mcp.Tool{...}` literal, not `mcp.NewTool(...)`
- `InputSchema` is derived from a typed `Arguments` struct using `inferSchema[T]()` helper
- `OutputSchema` is derived from the response type using `outputSchema[T]()` helper
- Pagination is added via `withPaginationConstraints()` wrapping `inferSchema`
- `OpenWorldHint` requires `ptr(true)` — it is a `*bool`; `ReadOnlyHint` and `DestructiveHint` are different: `ReadOnlyHint` is a plain `bool`, `DestructiveHint` is a `*bool`

### 4. Arguments Struct (replaces per-field `mcp.WithString/WithNumber` calls)

Each tool gets a typed `Arguments` struct. Field descriptions come from the `jsonschema` struct tag.

```go
// ListProjectsArguments holds the input parameters for listing projects within an organization.
type ListProjectsArguments struct {
    // Required fields — no omitempty
    TerraformOrgName string `json:"terraform_org_name" jsonschema:"The name of the Terraform Cloud/Enterprise organization"`

    // Optional pagination fields — embed the shared Pagination struct
    Pagination
}
```

Rules:
- Required fields: no `omitempty` on the JSON tag
- Optional fields: add `omitempty`
- Use `jsonschema:"<description>"` for the description (not a separate `mcp.Description()` call)
- Pagination: embed `Pagination` (the shared struct from the helpers file) — do **not** repeat `Page`/`PageSize` manually
- Enums / validation constraints: these cannot be expressed in struct tags alone; use `inferSchema` then manually set `Enum` or `Default` fields on the returned schema's `Properties` map, or define a custom `jsonschema.Schema` in `InputSchema` directly

### 5. Handler Function Changes

**Legacy pattern:**

```go
func listTerraformProjectsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
    terraformOrgName, err := request.RequireString("terraform_org_name")
    if err != nil {
        return ToolError(logger, "missing required input: terraform_org_name", err)
    }
    // ...
    return mcp.NewToolResultText(string(projectJSON)), nil
}
```

**New pattern:**

```go
func ListProjectsFunc(ctx context.Context, request *mcp.CallToolRequest, input ListProjectsArguments) (*mcp.CallToolResult, *ProjectSummaryList, error) {
    terraformOrgName := strings.TrimSpace(input.TerraformOrgName)
    if terraformOrgName == "" {
        return nil, nil, fmt.Errorf("terraform_org_name must not be blank")
    }
    // ...
    return nil, &ProjectSummaryList{...}, nil
}
```

Key differences:
- Signature: `func <Name>Func(ctx context.Context, request *mcp.CallToolRequest, input <Name>Arguments) (*mcp.CallToolResult, *<ResponseType>, error)`
- No logger parameter — errors are returned as Go errors, not `ToolError()`
- Input fields come from the typed `input` struct, not `request.RequireString()` / `request.GetArguments()`
- Return typed output struct as second return value (or `nil` if there is no structured output)
- Return `nil, nil, fmt.Errorf(...)` for validation/runtime errors
- **Never** return `mcp.NewToolResultText(...)` from the handler — the SDK serialises the typed struct automatically. Only return a non-nil `*mcp.CallToolResult` when you need to override the default serialisation (rare).

### 6. TFE Client Retrieval

**Legacy:**
```go
tfeClient, err := client.GetTfeClientFromContext(ctx, logger)
```

**New:**
```go
tfeClient, err := client.GetTfeClient(ctx, client.SessionIDFromRequest(request))
```

Import path for the new client: `github.com/hashicorp/terraform-mcp-server/pkg/mcp-official/client`

### 7. Pagination Helper Usage

**Legacy** (`pkg/utils`):
```go
pagination, err := utils.OptionalPaginationParams(request)
// ...
PageNumber: pagination.Page,
PageSize:   pagination.PageSize,
```

**New** (embed `Pagination` in the Arguments struct; use `input.ListOptions()`):
```go
type ListProjectsArguments struct {
    TerraformOrgName string `json:"terraform_org_name" jsonschema:"..."`
    Pagination
}

// in handler:
projects, err := tfeClient.Projects.List(ctx, terraformOrgName, &tfe.ProjectListOptions{
    ListOptions: input.ListOptions(),
})
```

The `Pagination` struct and its `ListOptions()` method are defined in the shared helpers file in `pkg/mcp-official/tools/tfe/`.

### 8. Response Helpers

**Legacy** (`mcp.NewToolResultText`, `mcp.NewToolResultError`, etc.) are **not available** in `modelcontextprotocol/go-sdk`.

**New**: Return the typed struct directly as the second return value. The SDK serialises it automatically. For tool-level errors (not Go errors), check `pkg/utils/result.go` for any available helpers.

### 9. Shared Helpers Available in `pkg/mcp-official/tools/tfe/`

The following helpers will be present in a shared file in the package (check `pkg/mcp-official/tools/tfe/` for the current state):

| Helper | Purpose |
|---|---|
| `inferSchema[T](toolName string) *jsonschema.Schema` | Derives an `InputSchema` from a typed struct using reflection |
| `outputSchema[T](toolName string) *jsonschema.Schema` | Derives an `OutputSchema` from a typed response struct |
| `withPaginationConstraints(s *jsonschema.Schema) *jsonschema.Schema` | Adds `minimum: 1` / `maximum: 100` constraints to `page`/`pageSize` fields |
| `ptr[T any](v T) *T` | Returns a pointer to any value (used for `*bool` annotation fields) |
| `nonNilSlice[T any](s []T) []T` | Returns an empty slice instead of nil (avoids `null` in JSON output) |
| `paginationDetails(p *tfe.Pagination) PaginationDetails` | Maps a TFE pagination object to the shared `PaginationDetails` struct |
| `Pagination` struct | Embeddable struct with `Page int` + `PageSize int` fields for Arguments |
| `PaginationDetails` struct | Response pagination shape embedded in list response types |

If any of these are missing from the file, implement them before migrating the tool.

### 10. RegisterTools Wiring

After creating the new tool files, register each tool pair in `pkg/mcp-official/tools/tools.go`:

```go
if toolsets.IsToolEnabled("list_terraform_projects", enabledToolsets) {
    mcp.AddTool(svr, tfeTools.ListProjectsTool(), tfeTools.ListProjectsFunc)
}
```

Note: The `mcp.AddTool` call uses generics — the compiler infers `TInput` and `TOutput` from the function signature. Make sure the tool function and the handler function match the same `Arguments` and response types.

## Naming Conventions (this repo)

| Concern | Convention | Example |
|---|---|---|
| Tool constructor | `<Action><Resource>Tool() *mcp.Tool` | `ListProjectsTool()` |
| Handler function | `<Action><Resource>Func(...)` | `ListProjectsFunc(...)` |
| Arguments struct | `<Action><Resource>Arguments` | `ListProjectsArguments` |
| Response struct | `<Resource>Summary` / `<Resource>Details` / `<Resource>Response` | `ProjectSummaryList`, `ProjectDetails`, `DeleteProjectResponse` |
| File name | snake_case matching the tool name | `list_terraform_projects.go` |

## Test Migration

Tests in the new package follow the same table-driven pattern as the legacy tests, but:

- Replace `mcp.CallToolRequest` construction with `mcp.NewCallToolRequest(name, args)` or equivalent SDK helpers
- Replace `server.ServerTool.Handler(ctx, req)` calls with direct `<Name>Func(ctx, req, input)` calls where practical
- If the legacy test file has a `//go:build ignore` build tag, **remove it** before starting work
- After migration, run `go test ./pkg/mcp-official/...` to verify all tests pass

## Common Pitfalls

1. **`OpenWorldHint` is `*bool`** — must use `ptr(true)` / `ptr(false)`. `ReadOnlyHint` is a plain `bool` — do **not** use `ptr()`.
2. **Do not import `pkg/client` directly** — always use `pkg/mcp-official/client` in the new package.
3. **Do not pass `logger` to tool functions** — the new SDK does not thread a logger through tool calls. Use `fmt.Errorf` for errors.
4. **Do not call `json.Marshal` in the handler** — return the typed struct and let the SDK handle serialisation.
5. **`nonNilSlice`** — always wrap slice fields with `nonNilSlice()` before embedding them in the response struct to avoid `null` in JSON output.
6. **Validation errors vs. runtime errors** — both are returned as `(nil, nil, fmt.Errorf(...))`. There is no separate `ToolError()` helper in the new SDK path.
7. **`ENABLE_TF_OPERATIONS` gate** — destructive tools (delete, force-unlock, action_run) in the legacy code are gated behind `isTerraformOperationsEnabled()` in `dynamic_tool.go`. Replicate this gate in the new `RegisterTools()` for the same tools.
8. **`InputSchema` can be `nil`** — if the tool takes no inputs, omit `InputSchema`. The SDK handles nil schemas correctly.

## End-to-End Checklist

- [ ] Remove `//go:build ignore` from source and test files (if present)
- [ ] Create `pkg/mcp-official/tools/tfe/<toolname>.go` following the new pattern
- [ ] Create `pkg/mcp-official/tools/tfe/<toolname>_test.go`
- [ ] Ensure all shared helpers exist in the package (add if missing)
- [ ] Register tool in `pkg/mcp-official/tools/tools.go`
- [ ] Run `go build ./pkg/mcp-official/...` — must compile with zero errors
- [ ] Run `go test ./pkg/mcp-official/...` — all tests must pass
- [ ] Confirm JSON field names on the Arguments and response structs match the legacy schema exactly (schema equivalence)
