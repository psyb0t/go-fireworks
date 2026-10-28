# Go code

Refines `~/.claude/rule-details/working-with-generated-code.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `types.gen.go` and `client.gen.go` are generated. Never edit them by hand. Change the spec patch, `config.yaml` or `config-types.yaml`, then run `make generate`.
- **GO2** The spec is downloaded into the gitignored `api/` directory and never committed. A spec problem oapi-codegen cannot handle gets fixed by a patch step under `scripts/` (jq for JSON, a Go tool under `internal/cmd` for YAML) that `gen.go` runs before oapi-codegen.
- **GO3** `config.yaml` and `config-types.yaml` keep the same `type-mapping`, `compatibility` and `exclude-operation-ids` options, or the client and the models disagree on types.
- **GO4** Hand-written Go is limited to `gen.go` and the tests. The package never logs and never prints.
- **GO5** `make lint` passes with zero issues. Lines are at most 80 columns, the `.golangci.yml` limit; split calls one argument per line rather than adding `//nolint`.
