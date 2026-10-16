package main

import (
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

func newTestPlugin(t *testing.T) (*Plugin, *logRecorder, *mocks.MockKVStore) {
	api := &logRecorder{API: &plugintest.API{}}
	t.Cleanup(func() { api.AssertExpectations(t) })
	store := mocks.NewMockKVStore(gomock.NewController(t))

	plugin := &Plugin{kvstore: store}
	plugin.SetAPI(api)
	plugin.client = pluginapi.NewClient(api, &plugintest.Driver{})
	plugin.router = plugin.initRouter()
	return plugin, api, store
}

func serve(plugin *Plugin, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	plugin.ServeHTTP(nil, w, r)
	return w
}

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin, _, _ := newTestPlugin(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	plugin.ServeHTTP(nil, w, r)

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
		plugin, api, _ := newTestPlugin(t)
		plugin.router.HandleFunc("/api/v1/panic", func(http.ResponseWriter, *http.Request) {
			panic("boom")
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/panic", nil)
		r.Header.Set("Mattermost-User-ID", "test-user-id")

		plugin.ServeHTTP(nil, w, r)

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
		plugin, api, _ := newTestPlugin(t)
		plugin.router.HandleFunc("/api/v1/panic", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("partial"))
			panic("boom")
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/panic", nil)
		r.Header.Set("Mattermost-User-ID", "test-user-id")

		plugin.ServeHTTP(nil, w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "partial", w.Body.String())
		require.Len(t, api.entries, 2)
		assert.Equal(t, http.StatusOK, api.entries[1].fields["status"])
	})
}

func TestServeHTTPLimitsRequestBody(t *testing.T) {
	plugin, _, _ := newTestPlugin(t)
	var readErr error
	plugin.router.HandleFunc("/api/v1/echo", func(_ http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/echo", strings.NewReader(strings.Repeat("a", maxRequestBodyBytes+1)))
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	plugin.ServeHTTP(nil, w, r)

	var maxBytesErr *http.MaxBytesError
	assert.ErrorAs(t, readErr, &maxBytesErr)
}

func TestServeHTTPLogsRequests(t *testing.T) {
	plugin, api, _ := newTestPlugin(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/hello?token=secret", nil)

	plugin.ServeHTTP(nil, w, r)

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

func TestGetGreetings(t *testing.T) {
	channelID := model.NewId()
	request := func(channelID string) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/channels/"+channelID+"/greetings", nil)
		r.Header.Set("Mattermost-User-ID", "test-user-id")
		return r
	}

	t.Run("returns the channel's greeting count", func(t *testing.T) {
		plugin, api, store := newTestPlugin(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionReadChannel).Return(true)
		store.EXPECT().GetGreetingCount(channelID).Return(2, nil)

		w := serve(plugin, request(channelID))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"count": 2}`, w.Body.String())
	})

	t.Run("rejects malformed channel IDs", func(t *testing.T) {
		plugin, _, _ := newTestPlugin(t)

		w := serve(plugin, request("not-a-channel-id"))

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "Invalid channel ID"}`, w.Body.String())
	})

	t.Run("rejects users who can't read the channel", func(t *testing.T) {
		plugin, api, _ := newTestPlugin(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionReadChannel).Return(false)

		w := serve(plugin, request(channelID))

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("logs store errors without returning them", func(t *testing.T) {
		plugin, api, store := newTestPlugin(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionReadChannel).Return(true)
		store.EXPECT().GetGreetingCount(channelID).Return(0, errors.New("pq: connection refused"))

		w := serve(plugin, request(channelID))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Failed to load greetings"}`, w.Body.String())
		require.NotEmpty(t, api.entries)
		assert.Equal(t, "error", api.entries[0].level)
		assert.Equal(t, "pq: connection refused", api.entries[0].fields["error"])
	})
}

func TestResetGreetings(t *testing.T) {
	channelID := model.NewId()
	request := func(channelID string) *http.Request {
		r := httptest.NewRequest(http.MethodDelete, "/api/v1/channels/"+channelID+"/greetings", nil)
		r.Header.Set("Mattermost-User-ID", "test-user-id")
		return r
	}
	captureAudit := func(api *logRecorder) *model.AuditRecord {
		rec := &model.AuditRecord{}
		api.On("LogAuditRec", mock.Anything).Run(func(args mock.Arguments) {
			*rec = *args.Get(0).(*model.AuditRecord)
		}).Once()
		return rec
	}

	t.Run("lets channel admins reset the count", func(t *testing.T) {
		plugin, api, store := newTestPlugin(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionManageChannelRoles).Return(true)
		store.EXPECT().ResetGreetingCount(channelID).Return(nil)
		audit := captureAudit(api)

		w := serve(plugin, request(channelID))

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "resetGreetingCount", audit.EventName)
		assert.Equal(t, model.AuditStatusSuccess, audit.Status)
		assert.Equal(t, "test-user-id", audit.Actor.UserId)
		assert.Equal(t, channelID, audit.EventData.Parameters["channel_id"])
	})

	t.Run("audits denied attempts", func(t *testing.T) {
		plugin, api, _ := newTestPlugin(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionManageChannelRoles).Return(false)
		audit := captureAudit(api)

		w := serve(plugin, request(channelID))

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, model.AuditStatusFail, audit.Status)
		assert.Equal(t, channelID, audit.EventData.Parameters["channel_id"])
	})

	t.Run("keeps unvalidated input out of the audit record", func(t *testing.T) {
		plugin, api, _ := newTestPlugin(t)
		audit := captureAudit(api)

		w := serve(plugin, request("not-a-channel-id"))

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, model.AuditStatusFail, audit.Status)
		assert.NotContains(t, audit.EventData.Parameters, "channel_id")
	})

	t.Run("logs store errors without returning them", func(t *testing.T) {
		plugin, api, store := newTestPlugin(t)
		api.On("HasPermissionToChannel", "test-user-id", channelID, model.PermissionManageChannelRoles).Return(true)
		store.EXPECT().ResetGreetingCount(channelID).Return(errors.New("pq: connection refused"))
		audit := captureAudit(api)

		w := serve(plugin, request(channelID))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Failed to reset greetings"}`, w.Body.String())
		assert.Equal(t, model.AuditStatusFail, audit.Status)
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
			plugin, _, _ := newTestPlugin(t)
			plugin.setConfiguration(&configuration{WebhookSecret: tc.configured})
			r := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(`{}`))
			if tc.provided != "" {
				r.Header.Set(webhookAuthHeader, tc.provided)
			}

			w := serve(plugin, r)

			assert.Equal(t, tc.status, w.Code)
		})
	}
}
