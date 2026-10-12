package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/pluginapi/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Exercises the real request path — server auth, then plugin ServeHTTP — unlike the
// httptest-based unit test in plugin_test.go.
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
