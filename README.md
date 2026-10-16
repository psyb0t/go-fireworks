# go-fireworks

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-fireworks.svg)](https://pkg.go.dev/github.com/psyb0t/go-fireworks)
[![CI](https://github.com/psyb0t/go-fireworks/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-fireworks/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-fireworks/badges/version.svg)](https://github.com/psyb0t/go-fireworks/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-fireworks/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-fireworks/badges/importers.svg)](https://github.com/psyb0t/go-fireworks/blob/badges/importers.md)

Unofficial Go client for the [Fireworks AI API](https://docs.fireworks.ai), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Fireworks AI publishes. I am not affiliated with Fireworks AI.

The package is `fireworks`, in the repository root. Everything in `types.gen.go` and `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-fireworks
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-fireworks"
)

func main() {
	auth := fireworks.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("FIREWORKS_API_KEY"))

		return nil
	})

	client, err := fireworks.NewClientWithResponses("https://api.fireworks.ai", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.GatewayListModelsWithResponse(
		context.Background(),
		"my-account",
		&fireworks.GatewayListModelsParams{},
	)
	if err != nil {
		panic(err)
	}

	for _, m := range *res.JSON200.Models {
		fmt.Println(*m.Name)
	}
}
```

`example_test.go` runs the same code against a local test server, so `make test` checks it.

Method names come from the spec's operation IDs. The control plane ones carry a `Gateway` prefix (`GatewayListModels`, `GatewayCreateDeployment`). The inference ones are long, because the spec's IDs are (`CreateChatCompletionV1ChatCompletionsPost`).

## Base URLs

The spec puts the control plane and inference on different base URLs, and one client has one base URL. Most paths live under `https://api.fireworks.ai`. These five live under `https://api.fireworks.ai/inference`:

| Path | Methods |
|---|---|
| `/v1/chat/completions` | `CreateChatCompletionV1ChatCompletionsPost` |
| `/v1/completions` | `CreateCompletionV1CompletionsPost` |
| `/v1/messages` | `MessagesPost` |
| `/v1/responses` | `CreateResponseV1ResponsesPost`, `ListResponsesV1ResponsesGet` |
| `/v1/responses/{response_id}` | `GetResponseV1ResponsesResponseIdGet`, `DeleteResponseV1ResponsesResponseIdDelete` |

For those, make a second client with the inference base URL:

```go
inference, err := fireworks.NewClientWithResponses("https://api.fireworks.ai/inference", auth)
```

## Auth

Fireworks AI uses a Bearer token: send the API key as `Authorization: Bearer <key>`. Create a key in the Fireworks AI dashboard or with `firectl`. The request editor in the usage example sets the header on every request.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` reads `openapi_spec_url` from the `.stats.yml` of [fw-ai-external/python-sdk](https://github.com/fw-ai-external/python-sdk) and downloads that URL to `api/openapi.yaml`. Stainless hosts the spec under a URL with a content hash that changes on every update, so the script never hardcodes it.
- oapi-codegen runs twice: `config-types.yaml` writes the models to `types.gen.go`, `config.yaml` writes the client to `client.gen.go`. One file would be about 60,000 lines. Both configs map `number` to `float64`, because oapi-codegen would otherwise use `float32`.

## Agent integrations

The [skill](.agents/skills/go-fireworks) teaches a coding agent how to use go-fireworks when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-fireworks@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-fireworks@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-fireworks
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Fireworks AI's published API specification. The specification itself is not redistributed here.
