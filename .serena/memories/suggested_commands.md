# Suggested commands

## Build / install
- `make build`
- `make install` (also runs `build`)
- `go install ./...` (equivalent to installing the provider)

## Format / lint
- `make fmt` (runs `gofmt -s -w -e .`)
- `make lint` (runs `golangci-lint run`)

## Test
- `make test` (runs `go test -v -cover -timeout=120s -parallel=10 ./...`)
- `make testacc` (runs `TF_ACC=1 go test -v -cover -timeout 120m ./...`; creates real resources)

## Generate docs / housekeeping
- `make generate` (runs `go generate` in `tools/`)
  - generates headers (copywrite)
  - `terraform fmt -recursive examples/`
  - regenerates `docs/` via `tfplugindocs`

## Useful repo navigation
- `rg "pattern" .` (fast search)
- `go test ./...` / `go test ./internal/...`
