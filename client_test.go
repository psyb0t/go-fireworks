package fireworks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-fireworks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "test-key"

func newClient(
	t *testing.T,
	pathPrefix string,
	handler http.HandlerFunc,
) *fireworks.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := fireworks.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := fireworks.NewClientWithResponses(srv.URL+pathPrefix, auth)
	require.NoError(t, err)

	return client
}

func assertBearer(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
}

func TestGatewayGetAccount(t *testing.T) {
	t.Parallel()

	client := newClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/accounts/acct-x", r.URL.Path)
		assertBearer(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"name":"accounts/acct-x","displayName":"Test",` +
				`"email":"user@example.com"}`,
		))
	})

	res, err := client.GatewayGetAccountWithResponse(
		context.Background(),
		"acct-x",
		&fireworks.GatewayGetAccountParams{},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, "accounts/acct-x", *res.JSON200.Name)
	assert.Equal(t, "Test", *res.JSON200.DisplayName)
	assert.Equal(t, "user@example.com", res.JSON200.Email)
}

func TestGatewayListModels(t *testing.T) {
	t.Parallel()

	client := newClient(t, "", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/accounts/acct-x/models", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("pageSize"))
		assertBearer(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"models":[{"name":"accounts/acct-x/models/m1"}],` +
				`"nextPageToken":"next","totalSize":1}`,
		))
	})

	pageSize := int32(10)

	res, err := client.GatewayListModelsWithResponse(
		context.Background(),
		"acct-x",
		&fireworks.GatewayListModelsParams{PageSize: &pageSize},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.Models)
	require.Len(t, *res.JSON200.Models, 1)

	assert.Equal(t, "accounts/acct-x/models/m1", *(*res.JSON200.Models)[0].Name)
	assert.Equal(t, "next", *res.JSON200.NextPageToken)
}

func TestCreateChatCompletion(t *testing.T) {
	t.Parallel()

	client := newClient(t, "/inference",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/inference/v1/chat/completions", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assertBearer(t, r)

			var body fireworks.ChatCompletionRequest

			assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "model-x", body.Model)

			if assert.NotNil(t, body.Messages) {
				assert.Len(t, *body.Messages, 1)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(
				`{"id":"c1","created":1,"model":"model-x","choices":` +
					`[{"index":0,"message":{"role":"assistant",` +
					`"content":"hi"},"finish_reason":"stop"}]}`,
			))
		},
	)

	var content fireworks.ChatMessage_Content

	require.NoError(t, content.FromChatMessageContent0("hello"))

	res, err := client.CreateChatCompletionV1ChatCompletionsPostWithResponse(
		context.Background(),
		fireworks.ChatCompletionRequest{
			Model: "model-x",
			Messages: &[]fireworks.ChatMessage{
				{Role: "user", Content: &content},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Choices, 1)

	choice := res.JSON200.Choices[0]
	assert.Equal(t, "stop", *choice.FinishReason)

	text, err := choice.Message.Content.AsChatMessageContent0()
	require.NoError(t, err)
	assert.Equal(t, "hi", text)
}
