# terraform-provider-manidae (overview)

## Purpose
A Terraform provider implemented in Go using the HashiCorp Terraform Plugin Framework. The repository currently matches the upstream “scaffolding” template (e.g., provider type name `scaffolding` and docs/examples under `docs/`/`examples/`).

## Tech stack
- Go (module: `github.com/indexyz/terraform-provider-manidae`, `go 1.24.6`)
- Terraform Plugin Framework (`github.com/hashicorp/terraform-plugin-framework`)
- Testing: `go test`, `terraform-plugin-testing`, `testify`
- Linting: `golangci-lint`
- Docs generation: `terraform-plugin-docs` via `go generate` in `tools/`

## Entrypoint
- `main.go`: runs `providerserver.Serve(...)` for the provider binary.
