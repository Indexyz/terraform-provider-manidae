# Style and conventions

## Go
- Formatting: `gofmt -s` (repo uses `gofmt -s -w -e .` via `make fmt`).
- Linting: `golangci-lint run` (see `.golangci.yml` for enabled linters).
- Typical layout: provider implementation under `internal/`.

## Terraform/examples/docs
- Examples live under `examples/` and are used by `terraform-plugin-docs` for `docs/` generation.
- `make generate` runs `go generate` in `tools/`, which includes `terraform fmt -recursive ../examples/`.

## Notes
- The repo currently contains template/scaffolding identifiers (e.g., provider name `scaffolding`). When customizing, update `main.go` provider address, `internal/provider` type name, and the docs generation provider name in `tools/tools.go`.
