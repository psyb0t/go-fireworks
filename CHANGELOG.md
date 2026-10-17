# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.0, 2026-10-16

- Added `fireworks`, a REST client generated with oapi-codegen from Fireworks AI's OpenAPI document (Gateway REST API, spec version 4.27.5, OpenAPI 3.1.0). It covers all 136 paths: accounts, API keys, users, models, datasets, deployments, fine-tuning and batch jobs, and the inference endpoints for chat completions, completions, messages and responses.
- Models are generated into `types.gen.go` and the client into `client.gen.go`.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` reads the current spec URL from the `.stats.yml` of Fireworks AI's Python SDK repo, downloads the spec from Stainless, and regenerates the client. The spec is not committed.
- Added the `go-fireworks` agent skill: `claude plugin install go-fireworks@psyb0t`, `codex plugin add go-fireworks@psyb0t`, `openclaw skills install @psyb0t/go-fireworks`.
