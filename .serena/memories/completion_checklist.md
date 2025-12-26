# What to do when a task is completed

- Run `make fmt` and `make lint`.
- Run `make test` (and `make testacc` only when acceptance coverage is needed and credentials/env are configured).
- If the task affects schema/docs/examples, run `make generate` and review `docs/` changes.
- Ensure `go mod tidy` is run when dependencies change (also run by GoReleaser pre-hooks).
