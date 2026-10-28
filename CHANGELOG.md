# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.0, 2026-10-28

- The client is generated from the spec Fireworks publishes at `https://docs.fireworks.ai/merged.openapi.yaml`. The Stainless-hosted copy of the spec it used before is no longer publicly readable, so scheduled regeneration had stopped.
- Breaking: the published spec does not include batch jobs, environments, node pools, MCP servers, leaderboards, snapshots, generic deployments, AWS IAM role bindings or `ExecuteRlorTrainStep`, so 56 operations are removed. Get API Key and Update API Key are left out because the spec lists them under a path without the key ID.
- Breaking: request body types are renamed from `GatewayXBody` to `GatewayGatewayXBody`, and `ChatCompletionRequest.Messages` is a `*[]ChatMessage`.
- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.1.3, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

## v0.1.2, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.1, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.0, 2026-10-16

- Added `fireworks`, a REST client generated with oapi-codegen from Fireworks AI's OpenAPI document (Gateway REST API, spec version 4.27.5, OpenAPI 3.1.0). It covers all 136 paths: accounts, API keys, users, models, datasets, deployments, fine-tuning and batch jobs, and the inference endpoints for chat completions, completions, messages and responses.
- Models are generated into `types.gen.go` and the client into `client.gen.go`.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` reads the current spec URL from the `.stats.yml` of Fireworks AI's Python SDK repo, downloads the spec from Stainless, and regenerates the client. The spec is not committed.
- Added the `go-fireworks` agent skill: `claude plugin install go-fireworks@psyb0t`, `codex plugin add go-fireworks@psyb0t`, `openclaw skills install @psyb0t/go-fireworks`.
