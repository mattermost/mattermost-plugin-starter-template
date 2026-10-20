package main

import (
	"context"
	"io"
	"net/http"
	"strings"
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

// Drives the slash command and the API against the server's real permissions: members
// read the status, non-members can't, and the role permission for editing a channel
// header decides who can change it.
func TestChannelStatusEndpoint(t *testing.T) {
	th := testhelper.Setup(t)
	ctx := context.Background()

	_, _, err := th.AdminClient.AddChannelMember(ctx, th.Channel.Id, th.User.Id)
	require.NoError(t, err)

	command := func(text string) string {
		resp, _, cmdErr := th.Client.ExecuteCommand(ctx, th.Channel.Id, text)
		require.NoError(t, cmdErr)
		return resp.Text
	}
	assert.Equal(t, "This channel has no status. Set one with `/channelstatus on-track`, `at-risk` or `off-track`.", command("/channelstatus"))
	assert.Equal(t, "Unknown status. Use `on-track`, `at-risk`, `off-track` or `clear`.", command("/channelstatus shipped"))
	assert.Equal(t, "Channel status set to **At risk**.", command("/channelstatus at-risk"))
	assert.Equal(t, "This channel is **At risk**.", command("/channelstatus"))

	channelStatus := func(client *model.Client4, method, channelID, body string) (int, string) {
		url := th.ServerURL + "/plugins/" + testhelper.PluginID() + "/api/v1/channels/" + channelID + "/status"
		req, reqErr := http.NewRequest(method, url, strings.NewReader(body))
		require.NoError(t, reqErr)
		req.Header.Set("Authorization", "Bearer "+client.AuthToken)

		resp, doErr := http.DefaultClient.Do(req)
		require.NoError(t, doErr)
		defer func() { _ = resp.Body.Close() }()
		respBody, readErr := io.ReadAll(resp.Body)
		require.NoError(t, readErr)
		return resp.StatusCode, string(respBody)
	}

	status, body := channelStatus(th.Client, http.MethodGet, th.Channel.Id, "")
	assert.Equal(t, http.StatusOK, status)
	assert.JSONEq(t, `{"status": "at-risk"}`, body)

	privateChannel := th.CreateChannel(model.ChannelTypePrivate)
	status, _ = channelStatus(th.Client, http.MethodGet, privateChannel.Id, "")
	assert.Equal(t, http.StatusForbidden, status)
	status, _ = channelStatus(th.Client, http.MethodPut, privateChannel.Id, `{"status": "on-track"}`)
	assert.Equal(t, http.StatusForbidden, status)

	assert.Equal(t, "Channel status cleared.", command("/channelstatus clear"))
	_, body = channelStatus(th.Client, http.MethodGet, th.Channel.Id, "")
	assert.JSONEq(t, `{"status": ""}`, body)

	status, _ = channelStatus(th.Client, http.MethodPut, th.Channel.Id, `{"status": "on-track"}`)
	assert.Equal(t, http.StatusOK, status)

	role, _, err := th.AdminClient.GetRoleByName(ctx, model.ChannelUserRoleId)
	require.NoError(t, err)
	original := role.Permissions
	t.Cleanup(func() {
		_, _, _ = th.AdminClient.PatchRole(context.Background(), role.Id, &model.RolePatch{Permissions: &original})
	})
	restricted := []string{}
	for _, permission := range original {
		if permission != model.PermissionManagePublicChannelProperties.Id {
			restricted = append(restricted, permission)
		}
	}
	_, _, err = th.AdminClient.PatchRole(ctx, role.Id, &model.RolePatch{Permissions: &restricted})
	require.NoError(t, err)

	assert.Equal(t, "You don't have permission to change this channel's status.", command("/channelstatus clear"))
	status, _ = channelStatus(th.Client, http.MethodDelete, th.Channel.Id, "")
	assert.Equal(t, http.StatusForbidden, status)
	_, body = channelStatus(th.Client, http.MethodGet, th.Channel.Id, "")
	assert.JSONEq(t, `{"status": "on-track"}`, body)
}
