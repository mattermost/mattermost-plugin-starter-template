package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (l *logRecorder) LogError(msg string, keyValuePairs ...any) {
	l.record("error", msg, keyValuePairs)
}

func newTestPlugin() (*Plugin, *logRecorder) {
	api := &logRecorder{API: &plugintest.API{}}
	plugin := &Plugin{}
	plugin.SetAPI(api)
	plugin.router = plugin.initRouter()
	return plugin, api
}

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin, _ := newTestPlugin()
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
		plugin, api := newTestPlugin()
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
		plugin, api := newTestPlugin()
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
	plugin, _ := newTestPlugin()
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
	plugin, api := newTestPlugin()
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
