---
name: go-fireworks
description: Unofficial Go client for the Fireworks AI REST API, generated with oapi-codegen. `fireworks.NewClientWithResponses(baseURL, fireworks.WithRequestEditorFn(...))` from `github.com/psyb0t/go-fireworks` gives one typed method per spec operation, covering all 136 paths of the Fireworks AI Gateway REST API, the control plane (accounts, API keys, users, models, datasets, deployments, fine-tuning and batch jobs) under `https://api.fireworks.ai` and the inference endpoints (chat completions, completions, messages, responses) under `https://api.fireworks.ai/inference`. Responses decode into typed `JSON200` fields. Use when Go code needs to call Fireworks AI, or already imports go-fireworks.
homepage: https://github.com/psyb0t/go-fireworks
user-invocable: true
metadata:
  openclaw:
    emoji: "🎆"
    requires:
      bins:
        - go
      env:
        - FIREWORKS_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the client sends HTTPS requests only to the base URL you give it, normally https://api.fireworks.ai or https://api.fireworks.ai/inference."
---

# go-fireworks

Typed Go client for the Fireworks AI REST API. Every operation in the spec is a method on `ClientWithResponses`, named after its operation ID.

For install, base URLs, method naming and error handling in detail, see [references/setup.md](references/setup.md).

## Security & safety

- Read the API key from an environment variable such as `FIREWORKS_API_KEY`. Never hardcode it, commit it or log it, and never print request headers.
- Requests go only to the base URL passed to `NewClientWithResponses`. Prompts, datasets and model files sent through the client reach Fireworks AI.
- Calls that spend money or delete things (create a deployment, start a fine-tuning or batch job, delete a model, dataset, deployment or API key) run only when the user asked for that action, and only on resources the task owns.

## When to use

- Go code calls Fireworks AI: chat or text completions, listing models, managing deployments, datasets or fine-tuning jobs.
- The project already imports `github.com/psyb0t/go-fireworks`.

## When NOT to use

- You need streaming. The generated client returns the whole body; for `stream: true` read `HTTPResponse.Body` from the plain `Client` methods yourself.
- The code is not Go. Use Fireworks AI's own SDKs.

## Usage

```go
import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/psyb0t/go-fireworks"
)

auth := fireworks.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FIREWORKS_API_KEY"))

	return nil
})

client, err := fireworks.NewClientWithResponses("https://api.fireworks.ai", auth)
if err != nil {
	return err
}

res, err := client.GatewayListModelsWithResponse(ctx, "my-account", &fireworks.GatewayListModelsParams{})
if err != nil {
	return err // transport error
}

if res.JSON200 == nil {
	return fmt.Errorf("list models: status %d: %s", res.StatusCode(), res.Body)
}

for _, m := range *res.JSON200.Models {
	fmt.Println(*m.Name)
}
```

Chat completions live on the inference base URL, so they need a second client:

```go
inference, err := fireworks.NewClientWithResponses("https://api.fireworks.ai/inference", auth)
if err != nil {
	return err
}

var content fireworks.ChatMessage_Content
if err := content.FromChatMessageContent0("Say hi"); err != nil {
	return err
}

res, err := inference.CreateChatCompletionV1ChatCompletionsPostWithResponse(ctx, fireworks.ChatCompletionRequest{
	Model:    "accounts/fireworks/models/<model-id>",
	Messages: []fireworks.ChatMessage{{Role: "user", Content: &content}},
})
if err != nil {
	return err
}

if res.JSON200 == nil {
	return fmt.Errorf("chat: status %d: %s", res.StatusCode(), res.Body)
}

text, err := res.JSON200.Choices[0].Message.Content.AsChatMessageContent0()
```

Check `JSON200` for nil before reading it. A non-2xx status is not a Go error; the status is in `StatusCode()` and the raw body in `Body`.
