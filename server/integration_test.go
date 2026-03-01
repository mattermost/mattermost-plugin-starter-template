package main

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-starter-template/server/testhelper"
)

// TestPluginActivation verifies that the plugin (built from this repository via `make dist`)
// starts successfully on a real Mattermost server. Setup() spins up Postgres + Mattermost
// containers, resets the database, uploads the plugin bundle, and enables it. This test
// then confirms the plugin reaches the Running state by querying the plugin status API.
//
// This is the most basic integration test — if this fails, the plugin cannot activate at all
// (e.g., OnActivate returns an error, missing dependencies, manifest issues).
func TestPluginActivation(t *testing.T) {
	th := testhelper.Setup(t)

	ctx := context.Background()
	statuses, _, err := th.AdminClient.GetPluginStatuses(ctx)
	require.NoError(t, err)

	found := false
	for _, s := range statuses {
		if s.PluginId == testhelper.PluginID() && s.State == model.PluginStateRunning {
			found = true
			break
		}
	}
	require.True(t, found, "plugin %s should be running", testhelper.PluginID())
}

// TestHelloEndpoint verifies the plugin's HTTP API works end-to-end through the Mattermost
// server. It calls GET /plugins/<plugin-id>/api/v1/hello with a valid user auth token and
// expects a 200 response with "Hello, world!".
//
// This exercises the full request path: client → Mattermost server (auth validation) →
// plugin ServeHTTP → MattermostAuthorizationRequired middleware → HelloWorld handler.
// Unlike the unit test in plugin_test.go which uses httptest, this hits the real server.
func TestHelloEndpoint(t *testing.T) {
	th := testhelper.Setup(t)

	// Build the full URL to the plugin's API endpoint. Mattermost serves plugin
	// routes at /plugins/<plugin-id>/, forwarding them to the plugin's ServeHTTP.
	pluginURL := th.ServerURL + "/plugins/" + testhelper.PluginID() + "/api/v1/hello"

	req, err := http.NewRequest(http.MethodGet, pluginURL, nil)
	require.NoError(t, err)

	// Authenticate as the test user. The Mattermost server validates this token,
	// then sets the Mattermost-User-ID header before forwarding to the plugin.
	req.Header.Set("Authorization", "Bearer "+th.Client.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, world!", string(body))
}
