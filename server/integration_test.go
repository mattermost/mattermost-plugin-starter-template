package main

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
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
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
}

// Runs /hello for real, then checks the count is readable by channel members only
// and resettable by channel admins only.
func TestGreetingsEndpoint(t *testing.T) {
	th := testhelper.Setup(t)
	ctx := context.Background()

	_, _, err := th.AdminClient.AddChannelMember(ctx, th.Channel.Id, th.User.Id)
	require.NoError(t, err)
	_, _, err = th.Client.ExecuteCommand(ctx, th.Channel.Id, "/hello @someone")
	require.NoError(t, err)

	greetings := func(client *model.Client4, method, channelID string) (int, string) {
		url := th.ServerURL + "/plugins/" + testhelper.PluginID() + "/api/v1/channels/" + channelID + "/greetings"
		req, reqErr := http.NewRequest(method, url, nil)
		require.NoError(t, reqErr)
		req.Header.Set("Authorization", "Bearer "+client.AuthToken)

		resp, doErr := http.DefaultClient.Do(req)
		require.NoError(t, doErr)
		defer func() { _ = resp.Body.Close() }()
		body, readErr := io.ReadAll(resp.Body)
		require.NoError(t, readErr)
		return resp.StatusCode, string(body)
	}

	status, body := greetings(th.Client, http.MethodGet, th.Channel.Id)
	assert.Equal(t, http.StatusOK, status)
	assert.JSONEq(t, `{"count": 1}`, body)

	privateChannel := th.CreateChannel(model.ChannelTypePrivate)
	status, _ = greetings(th.Client, http.MethodGet, privateChannel.Id)
	assert.Equal(t, http.StatusForbidden, status)

	status, _ = greetings(th.Client, http.MethodDelete, th.Channel.Id)
	assert.Equal(t, http.StatusForbidden, status)

	status, _ = greetings(th.AdminClient, http.MethodDelete, th.Channel.Id)
	assert.Equal(t, http.StatusNoContent, status)
	_, body = greetings(th.Client, http.MethodGet, th.Channel.Id)
	assert.JSONEq(t, `{"count": 0}`, body)
}
