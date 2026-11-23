# Contributing to Terraform MCP Server

> [!IMPORTANT]
> **Pre-refactor guidance:** This guide describes the repository before the official MCP Go SDK migration and planned 2.0 refactor are complete. Paths, commands, and environment variables may change. Review the [refactoring plan](https://github.com/hashicorp/proj-terraform-mcp-server/pull/3) and open pull requests before starting.
>
> The refactoring plan may only be accessible to HashiCorp contributors. In
> summary, the refactor will move application code from `pkg/` to `internal/`,
> restructure the CLI, centralize configuration, and rename several environment
> variables.

Thank you for contributing. Contributions are released under the project's
[MPL-2.0 license](LICENSE). By participating, you agree to follow the
[HashiCorp Community Guidelines](https://www.hashicorp.com/en/community-guidelines).

## Index

| Section | Purpose |
| --- | --- |
| [Welcome and scope](#welcome-and-scope) | What we accept and when to open an issue |
| [Quickstart](#quickstart) | Build and run the project locally |
| [Submit changes](#submit-changes) | Branches, commits, pull requests, and size limits |
| [Code ownership](#code-ownership) | How ownership works and how teams request it |
| [Code and style standards](#code-and-style-standards) | Required Go and MCP conventions |
| [Testing and linting](#testing-and-linting) | Checks to run before review |
| [AI-assisted contributions](#ai-assisted-contributions) | Expectations for AI-generated work |
| [Get help](#get-help) | Where to ask questions and report problems |

## Welcome and scope

We welcome:

- Bug fixes
- Documentation improvements
- Tests
- Performance and security improvements
- New features and MCP tools with a clear user need

Open an **issue** before implementing:

- A new tool, resource, or public capability
- A breaking change
- A new dependency or package
- An architectural, authentication, or security change
- A pull request expected to exceed the size targets below

For new tools, explain the user workflow, why existing tools are insufficient,
expected inputs and outputs, security implications, and testing plan. An API
endpoint alone is not a sufficient reason to add a tool.

The following are out of scope:

- New tools using the legacy Mark3Labs SDK
- Thin API wrappers without a clear agent workflow
- Unrelated refactoring mixed into a feature or bug fix
- HCP Terraform or Terraform Enterprise product support requests
- Terraform language or general Terraform workflow questions

Do not report vulnerabilities in a public issue. Follow the
[HashiCorp vulnerability reporting process](https://www.hashicorp.com/en/trust/security/vulnerability-management).

## Quickstart

### Requirements

| Dependency | Version | Required for |
| --- | --- | --- |
| Go | `1.26.10` | Build and tests |
| Git | Current supported release | Source control |
| Make | Current supported release | Build and test commands |
| Docker | Current supported release | Images and end-to-end tests |
| Node.js | `24.22.1` or newer 24.x | MCP Inspector only |

The Go version in [`.go-version`](.go-version) is the source of truth.

### Build from source

```shell
git clone https://github.com/YOUR-USERNAME/terraform-mcp-server.git
cd terraform-mcp-server

git remote add upstream https://github.com/hashicorp/terraform-mcp-server.git
make deps
make build
```

Run the Streamable HTTP server:

```shell
./bin/terraform-mcp-server streamable-http \
  --transport-host 127.0.0.1 \
  --transport-port 8080
```

From another terminal, check its health:

```shell
curl --fail http://127.0.0.1:8080/health
```

Registry tools work without Terraform credentials. HCP Terraform and Terraform
Enterprise tools require a valid token.

## Submit changes

### Git workflow

Create a branch from the latest upstream `main`:

```shell
git fetch upstream
git switch --create TYPE/short-description upstream/main
```

Use one of these branch types: `fix`, `feat`, `docs`, `test`, `refactor`, or
`chore`.

Write short, imperative commit messages using the same type:

```text
fix: reject blank workspace names
feat: add organization details tool
docs: clarify Inspector setup
```

### Keep pull requests small and focused

We strongly prefer smaller, focused pull requests over large, monolithic chunks of code. Smaller PRs are significantly easier to review, test, debug, and revert safely if needed.

- **One coherent change**: Focus each PR on a single tool, feature, or bug fix.
- **Separate concerns**: Do not combine behavior changes with unrelated formatting, dependency upgrades, file moves, or broad refactors.
- **Split large work**: If a feature is large, break it down into a series of smaller, self-contained PRs. If a larger change cannot be split safely, open an issue first to discuss the approach with maintainers.

### Pull request checklist

Before marking a pull request ready for review:

- [Link PR to issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword).
- Explain the user problem and the chosen solution.
- Include tests for new behavior and bug fixes.
- Run all applicable checks in [Testing and linting](#testing-and-linting).
- Update documentation for user-visible changes.
- Add an entry under `# Unreleased` in [`CHANGELOG.md`](CHANGELOG.md), or explain
  why one is unnecessary.
- Describe security, compatibility, and rollback considerations.
- Remove credentials and customer data from logs, output, and screenshots.
- Complete the repository pull request template.
- Sign the [HashiCorp Contributor License Agreement](https://www.hashicorp.com/en/cla)
  when prompted.

All required CI checks must pass before merge.

## Code ownership

[`CODEOWNERS`](CODEOWNERS) tells GitHub which teams maintain specific parts of
the repository. When a pull request changes a matching path, GitHub requests
review from that path's owners. Branch protection may require their approval.

Code ownership assigns review and maintenance responsibility. It does not grant
repository access or permission to merge.

### How it works

The repository-wide default owner is:

```text
* @hashicorp/team-proj-mcp-servers
```

More specific entries assign a domain to another team:

```text
/pkg/mcp-official/tools/search/** @hashicorp/team-tf-search
```

CODEOWNERS patterns are evaluated from top to bottom, and the last matching
pattern wins. Owners from earlier matching patterns do not carry over. List
every team that must review the path on the final matching line:

```text
/path/to/domain/** @hashicorp/team-proj-mcp-servers @hashicorp/team-name
```

Use narrow patterns that cover the implementation, tests, registration,
configuration, and documentation the team will maintain. See
[GitHub's CODEOWNERS documentation](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
for pattern syntax.

### Request ownership for your team

Before editing `CODEOWNERS`, open a
[GitHub issue](https://github.com/hashicorp/terraform-mcp-server/issues/new/choose)
and tag `@hashicorp/team-proj-mcp-servers`. Include:

- The GitHub team handle
- The paths the team wants to own
- The product or technical domain covered by those paths
- Why the team is the appropriate long-term owner
- Confirmation that the team will review changes and maintain the code
- Any existing owners that should remain co-owners

The GitHub team must be visible and have explicit write access to the
repository. Project maintainers will confirm the final patterns and whether
shared ownership is required.

After approval, submit a focused pull request that updates `CODEOWNERS` and
links to the issue. Do not bury an ownership change in an unrelated pull
request or request ownership for a one-time contribution.

## Code and style standards

Follow [Effective Go](https://go.dev/doc/effective_go),
[Go Code Review Comments](https://go.dev/wiki/CodeReviewComments), and the
surrounding repository code.

- Format changed Go files with `gofmt`.
- Use `github.com/modelcontextprotocol/go-sdk/mcp` for new MCP functionality.
- Do not add new functionality under `pkg/tools/` or `pkg/mcp-mark3labs/`.
- Reuse the clients, sessions, middleware, and pagination already in the
  repository.
- Validate untrusted input before calling external services.
- Propagate `context.Context` and wrap operational errors with context and `%w`.
- Use typed tool inputs and focused response types.
- Keep list responses bounded and paginated.
- Set accurate MCP read-only, destructive, and open-world annotations.
- Preserve public tool names, JSON fields, and behavior unless an approved issue
  changes them.
- Copy the copyright and SPDX header from an adjacent source file.
- Never log or return tokens, signed URLs, state data, sensitive variables, or
  unnecessary customer data.

Do not add an abstraction, interface, or dependency without a current need. Do
not return complete `go-tfe` models when a smaller response is sufficient.

### Adding an MCP tool

Follow [Adding an MCP tool](docs/adding-mcp-tool.md) for the implementation,
registration, testing, and documentation checklist.

## Testing and linting

Pull requests must pass the applicable checks before they are marked ready for
review.

| Guide | Use it for |
| --- | --- |
| [Unit tests](docs/unit-tests.md) | Package, schema, validation, and registration tests |
| [HCPT tests](docs/hcpt-tests.md) | Tests against real HCP Terraform or Terraform Enterprise APIs |
| [Linting checks](docs/linting.md) | Formatting, static analysis, headers, and diff checks |
| [MCP Inspector](docs/mcp-inspector.md) | Manual tool and schema verification |

### When relevant

| Change | Command |
| --- | --- |
| Public Registry or transport behavior | `make test-e2e` |
| Official SDK public Registry behavior | `make test-e2e-official-public-registry` |
| HCP Terraform tool behavior | `make test-hcpt RUN=TestName` |
| Docker image | `make docker-build` |
| Helm chart | `helm lint helm/terraform-mcp-server` |

Docker-backed tests require Docker. HCP Terraform tests require a running
server, test credentials, and an isolated test organization. If you cannot run
a required integration test, say so in the pull request and provide the exact
test maintainers should run.

### MCP Inspector

Follow [Testing with MCP Inspector](docs/mcp-inspector.md) for web, CLI, stdio,
and Streamable HTTP instructions.

## AI-assisted contributions

AI tools are permitted, but AI assistance does not lower the review standard.

If you use AI:

- Review, understand, and test every submitted line.
- Verify APIs and behavior against current code and official documentation.
- Remove duplicated code, unnecessary abstractions, fabricated APIs, and
  irrelevant tests.
- Never send credentials, customer data, private source code, or unreleased information to an unapproved service.

Do not submit generated output that you cannot explain or maintain. Maintainers
may close AI-generated pull requests that are unreviewed, oversized, incorrect,
or repeatedly ignore repository conventions. The contributor remains
responsible for the entire change, regardless of how it was produced.

## Get help

| Topic | Channel |
| --- | --- |
| Contribution design or project questions | [GitHub Discussions](https://github.com/hashicorp/terraform-mcp-server/discussions) |
| Bugs and feature requests | [GitHub Issues](https://github.com/hashicorp/terraform-mcp-server/issues) |
| HCP Terraform or Terraform Enterprise support | [HashiCorp Support](https://support.hashicorp.com/hc/en-us/requests/new) |
| Terraform language and workflow questions | [HashiCorp Discuss](https://discuss.hashicorp.com/) |
| Security vulnerabilities | `security@hashicorp.com` |
| Community behavior | [HashiCorp Community Guidelines](https://www.hashicorp.com/en/community-guidelines) |

Ask early when you are unsure. A short discussion is easier than rewriting a
completed contribution.
