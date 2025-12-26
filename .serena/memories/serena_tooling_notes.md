# Serena tooling notes

- `list_dir`/`find_file` work as expected.
- Go LSP-backed symbol tools (`get_symbols_overview`) returned `file not found (-32098)` for `main.go` and `internal/provider/provider.go` in this environment; fallback to `rg`/`sed` reads if this persists.
