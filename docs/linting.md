# Linting checks

[Back to the contribution guide](../CONTRIBUTING.md)

Run these checks before marking a Go pull request ready for review.

## Format changed files

Run `gofmt` on each Go file you changed:

```shell
gofmt -w path/to/changed.go path/to/changed_test.go
```

Do not include repository-wide formatting in an unrelated change.

## Run static and diff checks

```shell
go vet ./...
git diff --check
```

The repository does not currently have a `golangci-lint` configuration or a
`make lint` target. Do not add or require a new linter in an unrelated pull
request.

## Check source headers

Copy the copyright and SPDX header from an adjacent source file. The repository
license-header configuration is in [`.copywrite.hcl`](../.copywrite.hcl).

## Check dependency changes

If the pull request intentionally changes Go dependencies, run:

```shell
go mod tidy
```

Review `go.mod` and `go.sum`, and remove unrelated changes. Do not add a
dependency for behavior already covered by the standard library or an existing
dependency.

For Markdown-only changes, run `git diff --check` and preview the rendered
headings, tables, code blocks, and links.
