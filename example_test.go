package fireworks_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/psyb0t/go-fireworks"
)

// This is the README usage example. A local test server stands in for
// https://api.fireworks.ai so the example runs without credentials.
func Example() {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(
				`{"models":[{"name":"accounts/my-account/models/m1"}]}`,
			))
		},
	))
	defer srv.Close()

	baseURL := srv.URL // https://api.fireworks.ai in real use

	auth := fireworks.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("FIREWORKS_API_KEY"),
			)

			return nil
		},
	)

	client, err := fireworks.NewClientWithResponses(baseURL, auth)
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

	// Output: accounts/my-account/models/m1
}
