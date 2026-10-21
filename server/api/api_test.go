package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest/mock"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mattermost/mattermost-plugin-starter-template/server/channelstatus"
	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore/mocks"
)

type logEntry struct {
	level  string
	msg    string
	fields map[string]any
}

// logRecorder captures log calls; plugintest.API can't match them because
// the request logger emits fields in random map order.
type logRecorder struct {
	*plugintest.API
	entries []logEntry
}

func (l *logRecorder) record(level, msg string, keyValuePairs []any) {
	fields := map[string]any{}
	for i := 0; i+1 < len(keyValuePairs); i += 2 {
		fields[keyValuePairs[i].(string)] = keyValuePairs[i+1]
	}
	l.entries = append(l.entries, logEntry{level: level, msg: msg, fields: fields})
}

func (l *logRecorder) LogDebug(msg string, keyValuePairs ...any) {
	l.record("debug", msg, keyValuePairs)
}

func (l *logRecorder) LogWarn(msg string, keyValuePairs ...any) {
	l.record("warn", msg, keyValuePairs)
}

func (l *logRecorder) LogError(msg string, keyValuePairs ...any) {
	l.record("error", msg, keyValuePairs)
}

func newTestHandler(t *testing.T) (*Handler, *logRecorder, *mocks.MockKVStore) {
	api := &logRecorder{API: &plugintest.API{}}
	t.Cleanup(func() { api.AssertExpectations(t) })
	store := mocks.NewMockKVStore(gomock.NewController(t))

	statuses := channelstatus.NewService(api, pluginapi.NewClient(api, &plugintest.Driver{}), store)
	return NewHandler(api, statuses, func() string { return "" }), api, store
}

func serve(h *Handler, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(nil, w, r)
	return w
}

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	h, _, _ := newTestHandler(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	h.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.NotNil(result)
	defer func() { _ = result.Body.Close() }()
	bodyBytes, err := io.ReadAll(result.Body)
	assert.Nil(err)
	bodyString := string(bodyBytes)

	assert.Equal("Hello, world!", bodyString)
}

func TestServeHTTPRecoversFromPanic(t *testing.T) {
	t.Run("before the response is written", func(t *testing.T) {
		h, api, _ := newTestHandler(t)
		h.router.HandleFunc("/api/v1/panic", func(http.ResponseWriter, *http.Request) {
			panic("boom")
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/panic", nil)
		r.Header.Set("Mattermost-User-ID", "test-user-id")

		h.ServeHTTP(nil, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		require.Len(t, api.entries, 2)
		panicLog := api.entries[0]
		assert.Equal(t, "error", panicLog.level)
		assert.Equal(t, "Recovered from panic in HTTP handler", panicLog.msg)
		assert.Equal(t, "boom", panicLog.fields["panic"])
		assert.Equal(t, "test-user-id", panicLog.fields["user_id"])
		assert.NotEmpty(t, panicLog.fields["stack"])
		assert.Equal(t, http.StatusInternalServerError, api.entries[1].fields["status"])
		assert.Equal(t, "test-user-id", api.entries[1].fields["user_id"])
	})

	t.Run("after the response is written", func(t *testing.T) {
		h, api, _ := newTestHandler(t)
		h.router.HandleFunc("/api/v1/panic", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("partial"))
			panic("boom")
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/panic", nil)
		r.Header.Set("Mattermost-User-ID", "test-user-id")

		h.ServeHTTP(nil, w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "partial", w.Body.String())
		require.Len(t, api.entries, 2)
		assert.Equal(t, http.StatusOK, api.entries[1].fields["status"])
	})
}

func TestServeHTTPLimitsRequestBody(t *testing.T) {
	h, _, _ := newTestHandler(t)
	var readErr error
	h.router.HandleFunc("/api/v1/echo", func(_ http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/echo", strings.NewReader(strings.Repeat("a", maxRequestBodyBytes+1)))
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	h.ServeHTTP(nil, w, r)

	var maxBytesErr *http.MaxBytesError
	assert.ErrorAs(t, readErr, &maxBytesErr)
}

func TestServeHTTPLogsRequests(t *testing.T) {
	h, api, _ := newTestHandler(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/hello?token=secret", nil)

	h.ServeHTTP(nil, w, r)

	require.Len(t, api.entries, 1)
	entry := api.entries[0]
	assert.Equal(t, "debug", entry.level)
	assert.Equal(t, "HTTP request", entry.msg)
	assert.Equal(t, http.MethodGet, entry.fields["method"])
	assert.Equal(t, "/api/v1/hello", entry.fields["path"])
	assert.Equal(t, http.StatusUnauthorized, entry.fields["status"])
	assert.Equal(t, "", entry.fields["user_id"])
	assert.NotEmpty(t, entry.fields["duration"])
}

// The permission rules themselves are tested in the channelstatus package; these tests cover
// how the routes validate input and turn results into responses.
func TestChannelStatusAPI(t *testing.T) {
	channelID := model.NewId()
	request := func(method, channelID, body string) *http.Request {
		r := httptest.NewRequest(method, "/api/v1/channels/"+channelID+"/status", strings.NewReader(body))
		r.Header.Set("Mattermost-User-ID", "test-user-id")
		return r
	}
	allowEdit := func(api *logRecorder) {
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionReadChannel).Return(true)
		api.On("GetChannel", channelID).Return(&model.Channel{Id: channelID, Type: model.ChannelTypeOpen}, nil)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionManagePublicChannelProperties).Return(true)
		api.On("LogAuditRec", mock.Anything).Once()
	}

	t.Run("GET returns the channel's status", func(t *testing.T) {
		h, api, store := newTestHandler(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionReadChannel).Return(true)
		store.EXPECT().GetChannelStatus(channelID).Return(channelstatus.AtRisk, nil)

		w := serve(h, request(http.MethodGet, channelID, ""))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"status": "at-risk"}`, w.Body.String())
	})

	t.Run("GET rejects users who can't read the channel", func(t *testing.T) {
		h, api, _ := newTestHandler(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionReadChannel).Return(false)

		w := serve(h, request(http.MethodGet, channelID, ""))

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("rejects malformed channel IDs", func(t *testing.T) {
		h, _, _ := newTestHandler(t)

		w := serve(h, request(http.MethodPut, "not-a-channel-id", `{"status": "at-risk"}`))

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid channel ID"}`, w.Body.String())
	})

	t.Run("PUT sets the status", func(t *testing.T) {
		h, api, store := newTestHandler(t)
		allowEdit(api)
		store.EXPECT().SetChannelStatus(channelID, channelstatus.OnTrack).Return(nil)
		api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Once()

		w := serve(h, request(http.MethodPut, channelID, `{"status": "on-track"}`))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status": "on-track"}`, w.Body.String())
	})

	for name, body := range map[string]string{
		"invalid JSON":   `{`,
		"missing status": `{}`,
	} {
		t.Run("PUT rejects "+name, func(t *testing.T) {
			h, _, _ := newTestHandler(t)

			w := serve(h, request(http.MethodPut, channelID, body))

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, `{"error": "Request body must include a status"}`, w.Body.String())
		})
	}

	t.Run("PUT rejects unknown statuses", func(t *testing.T) {
		h, api, _ := newTestHandler(t)
		api.On("LogAuditRec", mock.Anything).Once()

		w := serve(h, request(http.MethodPut, channelID, `{"status": "shipped"}`))

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Unknown status"}`, w.Body.String())
	})

	t.Run("DELETE clears the status", func(t *testing.T) {
		h, api, store := newTestHandler(t)
		allowEdit(api)
		store.EXPECT().DeleteChannelStatus(channelID).Return(nil)
		api.On("PublishWebSocketEvent", mock.Anything, mock.Anything, mock.Anything).Once()

		w := serve(h, request(http.MethodDelete, channelID, ""))

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("logs store errors without returning them", func(t *testing.T) {
		h, api, store := newTestHandler(t)
		allowEdit(api)
		store.EXPECT().SetChannelStatus(channelID, channelstatus.OnTrack).Return(errors.New("pq: connection refused"))

		w := serve(h, request(http.MethodPut, channelID, `{"status": "on-track"}`))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Failed to update the channel status"}`, w.Body.String())
		require.NotEmpty(t, api.entries)
		assert.Equal(t, "error", api.entries[0].level)
		assert.Equal(t, "pq: connection refused", api.entries[0].fields["error"])
	})
}

func TestWebhook(t *testing.T) {
	for name, tc := range map[string]struct {
		configured string
		provided   string
		status     int
	}{
		"accepts the configured secret":            {configured: "s3cret", provided: "s3cret", status: http.StatusNoContent},
		"rejects a wrong secret":                   {configured: "s3cret", provided: "guess", status: http.StatusUnauthorized},
		"rejects a missing secret":                 {configured: "s3cret", provided: "", status: http.StatusUnauthorized},
		"rejects everything until it's configured": {configured: "", provided: "", status: http.StatusUnauthorized},
	} {
		t.Run(name, func(t *testing.T) {
			h, _, _ := newTestHandler(t)
			h.webhookSecret = func() string { return tc.configured }
			r := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(`{}`))
			if tc.provided != "" {
				r.Header.Set(webhookAuthHeader, tc.provided)
			}

			w := serve(h, r)

			assert.Equal(t, tc.status, w.Code)
		})
	}
}

// Exercises the real request path — server auth, then plugin ServeHTTP — unlike the
// httptest-based TestServeHTTP above.
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
