package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/pluginapi/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelloEndpoint verifies the plugin's HTTP API works end-to-end through a real
// Mattermost server. Setup() spins up Postgres + Mattermost containers, deploys the
// plugin built from this repository (via `make dist`), and enables it. This test then
// calls GET /plugins/<plugin-id>/api/v1/hello with a valid user auth token and expects
// a 200 response with "Hello, world!".
//
// This exercises the full request path: client -> Mattermost server (auth validation) ->
// plugin ServeHTTP -> MattermostAuthorizationRequired middleware -> HelloWorld handler.
// Unlike the unit test in plugin_test.go which uses httptest, this hits a real server.
func TestHelloEndpoint(t *testing.T) {
	th := testhelper.Setup(t)

	pluginURL := th.ServerURL + "/plugins/" + testhelper.PluginID() + "/api/v1/hello"

	req, err := http.NewRequest(http.MethodGet, pluginURL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+th.Client.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, world!", string(body))
}
