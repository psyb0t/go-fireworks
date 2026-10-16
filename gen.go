// Package fireworks is a Go client for the Fireworks AI REST API, generated
// with oapi-codegen from the OpenAPI document Fireworks AI publishes.
package fireworks

// The spec is large, so the models and the client are generated into
// separate files from the same spec.
//go:generate bash scripts/fetch-spec.sh
//go:generate go tool oapi-codegen --config=config-types.yaml api/openapi.yaml
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.yaml
