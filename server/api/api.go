// Package api serves the plugin's HTTP routes under /plugins/{plugin_id}.
package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/bot/logger"

	"github.com/mattermost/mattermost-plugin-starter-template/server/channelstatus"
)

// Handler routes the requests the server forwards to the plugin's ServeHTTP hook.
type Handler struct {
	pluginAPI     plugin.API
	statuses      *channelstatus.Service
	webhookSecret func() string
	router        *mux.Router
}

// NewHandler reads webhookSecret on every request, so a changed setting applies without reactivating the plugin.
func NewHandler(pluginAPI plugin.API, statuses *channelstatus.Service, webhookSecret func() string) *Handler {
	h := &Handler{pluginAPI: pluginAPI, statuses: statuses, webhookSecret: webhookSecret}
	h.router = h.initRouter()
	return h
}

func (h *Handler) initRouter() *mux.Router {
	router := mux.NewRouter()

	// The Mattermost server forwards plugin requests without the panic recovery,
	// body size limit, or security headers it applies to its own REST API.
	// logRequests comes first so it records the 500 written by recoverPanic.
	router.Use(h.logRequests, h.recoverPanic, limitRequestBody, setSecurityHeaders)

	// External services have no Mattermost session, so webhooks sit outside the
	// authenticated subrouter and prove themselves with a shared secret instead.
	router.HandleFunc("/webhook", h.handleWebhook).Methods(http.MethodPost)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(h.MattermostAuthorizationRequired)

	// The server only checks CSRF tokens on non-GET requests, so anything that changes
	// state must use POST, PUT, PATCH or DELETE.
	apiRouter.HandleFunc("/hello", h.HelloWorld).Methods(http.MethodGet)
	apiRouter.HandleFunc("/channels/{channel_id}/status", h.handleGetStatus).Methods(http.MethodGet)
	apiRouter.HandleFunc("/channels/{channel_id}/status", h.handleSetStatus).Methods(http.MethodPut)
	apiRouter.HandleFunc("/channels/{channel_id}/status", h.handleClearStatus).Methods(http.MethodDelete)

	return router
}

type pluginContextKey struct{}

// ServeHTTP has the plugin hook's signature; the plugin's ServeHTTP passes its arguments straight through.
// The root URL is currently <siteUrl>/plugins/com.mattermost.plugin-starter-template/api/v1/. Replace com.mattermost.plugin-starter-template with the plugin ID.
func (h *Handler) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pluginContextKey{}, c)))
}

// MattermostAuthorizationRequired can trust Mattermost-User-ID because the server sets it only after authenticating the request.
// Requests from other plugins arrive without a user ID and are rejected here; routes meant for
// them should check Mattermost-Plugin-ID against an allowlist instead.
func (h *Handler) MattermostAuthorizationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// maxRequestBodyBytes bounds how much of a request body a handler will read.
// Raise it for endpoints that legitimately accept larger payloads.
const maxRequestBodyBytes = 1 << 20 // 1 MiB

type loggerContextKey struct{}

// requestLogger returns the logger attached by logRequests, which tags every line with the requesting user.
func requestLogger(r *http.Request) logger.Logger {
	return r.Context().Value(loggerContextKey{}).(logger.Logger)
}

// logRequests attaches a request-scoped logger and logs every request at debug level,
// since the server doesn't log plugin requests.
func (h *Handler) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := logger.New(h.pluginAPI).With(logger.LogContext{
			"user_id": r.Header.Get("Mattermost-User-ID"),
		})
		r = r.WithContext(context.WithValue(r.Context(), loggerContextKey{}, log))
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		// The query string is omitted because it may carry tokens or other secrets.
		log.With(logger.LogContext{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   rec.status,
			"duration": time.Since(start).String(),
		}).Debugf("HTTP request")
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// WriteHeader records only the first status, since later ones never reach the client.
func (r *statusRecorder) WriteHeader(status int) {
	if !r.wroteHeader {
		r.status = status
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	r.wroteHeader = true
	return r.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

// recoverPanic stops a panicking handler from crashing the plugin process.
// The server deactivates a plugin that crashes repeatedly.
func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			if panicVal := recover(); panicVal != nil {
				requestLogger(r).With(logger.LogContext{
					"path":  r.URL.Path,
					"panic": fmt.Sprint(panicVal),
					"stack": string(debug.Stack()),
				}).Errorf("Recovered from panic in HTTP handler")
				// Once the response has started, an error body would be appended to it.
				if !rec.wroteHeader {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(rec, r)
	})
}

func limitRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		next.ServeHTTP(w, r)
	})
}

func setSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello, world!")); err != nil {
		requestLogger(r).WithError(err).Warnf("Failed to write response")
	}
}

type statusBody struct {
	Status string `json:"status"`
}

// channelIDFromPath rejects malformed IDs before they reach a permission check, the store or a log.
func channelIDFromPath(w http.ResponseWriter, r *http.Request) (string, bool) {
	channelID := mux.Vars(r)["channel_id"]
	if !model.IsValidId(channelID) {
		writeError(w, r, http.StatusBadRequest, "Invalid channel ID", nil)
		return "", false
	}
	return channelID, true
}

func (h *Handler) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	channelID, ok := channelIDFromPath(w, r)
	if !ok {
		return
	}

	status, err := h.statuses.Get(r.Header.Get("Mattermost-User-ID"), channelID)
	if err != nil {
		writeStatusError(w, r, err, "Failed to load the channel status")
		return
	}

	writeJSON(w, r, http.StatusOK, statusBody{Status: status})
}

func (h *Handler) handleSetStatus(w http.ResponseWriter, r *http.Request) {
	channelID, ok := channelIDFromPath(w, r)
	if !ok {
		return
	}

	var body statusBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, r, http.StatusBadRequest, "Request body must include a status", nil)
		return
	}

	if err := h.statuses.Set(pluginContext(r), r.URL.Path, r.Header.Get("Mattermost-User-ID"), channelID, body.Status); err != nil {
		writeStatusError(w, r, err, "Failed to update the channel status")
		return
	}

	writeJSON(w, r, http.StatusOK, body)
}

func (h *Handler) handleClearStatus(w http.ResponseWriter, r *http.Request) {
	channelID, ok := channelIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.statuses.Set(pluginContext(r), r.URL.Path, r.Header.Get("Mattermost-User-ID"), channelID, ""); err != nil {
		writeStatusError(w, r, err, "Failed to clear the channel status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeStatusError turns channelstatus errors into responses. Unexpected errors are logged
// and replaced with failure, since they can hold internal details.
func writeStatusError(w http.ResponseWriter, r *http.Request, err error, failure string) {
	switch {
	case errors.Is(err, channelstatus.ErrForbidden):
		writeError(w, r, http.StatusForbidden, "You don't have permission to do that in this channel", nil)
	case errors.Is(err, channelstatus.ErrInvalidStatus):
		writeError(w, r, http.StatusBadRequest, "Unknown status", nil)
	default:
		writeError(w, r, http.StatusInternalServerError, failure, err)
	}
}

// pluginContext returns the session details ServeHTTP received, for audit records.
func pluginContext(r *http.Request) *plugin.Context {
	c, _ := r.Context().Value(pluginContextKey{}).(*plugin.Context)
	return c
}

// webhookAuthHeader carries the secret generated in the plugin settings. Prefer a header
// over a query parameter: URLs end up in proxy and browser logs.
const webhookAuthHeader = "X-Webhook-Secret"

func (h *Handler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	secret := h.webhookSecret()
	provided := r.Header.Get(webhookAuthHeader)

	// ConstantTimeCompare keeps response timing from revealing how much of the secret matched.
	if secret == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(secret)) != 1 {
		writeError(w, r, http.StatusUnauthorized, "Invalid webhook secret", nil)
		return
	}

	// Handle the payload here. limitRequestBody already caps how much of r.Body can be read.
	w.WriteHeader(http.StatusNoContent)
}

type errorResponse struct {
	Error string `json:"error"`
}

// writeError sends a message that is safe to show users and logs err, which may hold
// internal details such as SQL errors or upstream responses, only on the server.
func writeError(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	if err != nil {
		requestLogger(r).WithError(err).Errorf("%s", message)
	}
	writeJSON(w, r, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		requestLogger(r).WithError(err).Warnf("Failed to write response")
	}
}
