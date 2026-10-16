package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/bot/logger"
)

// initRouter initializes the HTTP router for the plugin.
func (p *Plugin) initRouter() *mux.Router {
	router := mux.NewRouter()

	// The Mattermost server forwards plugin requests without the panic recovery,
	// body size limit, or security headers it applies to its own REST API.
	// logRequests comes first so it records the 500 written by recoverPanic.
	router.Use(p.logRequests, p.recoverPanic, limitRequestBody, setSecurityHeaders)

	// External services have no Mattermost session, so webhooks sit outside the
	// authenticated subrouter and prove themselves with a shared secret instead.
	router.HandleFunc("/webhook", p.handleWebhook).Methods(http.MethodPost)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(p.MattermostAuthorizationRequired)

	// The server only checks CSRF tokens on non-GET requests, so anything that changes
	// state must use POST, PUT, PATCH or DELETE.
	apiRouter.HandleFunc("/hello", p.HelloWorld).Methods(http.MethodGet)
	apiRouter.HandleFunc("/channels/{channel_id}/greetings", p.handleGetGreetings).Methods(http.MethodGet)
	apiRouter.HandleFunc("/channels/{channel_id}/greetings", p.handleResetGreetings).Methods(http.MethodDelete)

	return router
}

type pluginContextKey struct{}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
// The root URL is currently <siteUrl>/plugins/com.mattermost.plugin-starter-template/api/v1/. Replace com.mattermost.plugin-starter-template with the plugin ID.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pluginContextKey{}, c)))
}

// MattermostAuthorizationRequired can trust Mattermost-User-ID because the server sets it only after authenticating the request.
// Requests from other plugins arrive without a user ID and are rejected here; routes meant for
// them should check Mattermost-Plugin-ID against an allowlist instead.
func (p *Plugin) MattermostAuthorizationRequired(next http.Handler) http.Handler {
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
func (p *Plugin) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := logger.New(p.API).With(logger.LogContext{
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
func (p *Plugin) recoverPanic(next http.Handler) http.Handler {
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

func (p *Plugin) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello, world!")); err != nil {
		requestLogger(r).WithError(err).Warnf("Failed to write response")
	}
}

type greetingsResponse struct {
	Count int `json:"count"`
}

func (p *Plugin) handleGetGreetings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	channelID := mux.Vars(r)["channel_id"]

	if !model.IsValidId(channelID) {
		writeError(w, r, http.StatusBadRequest, "Invalid channel ID", nil)
		return
	}

	// The server only authenticates the user. Every endpoint that takes an ID must check
	// the user may access that object, or any user can read any channel's data.
	if !p.client.User.HasPermissionToChannel(userID, channelID, model.PermissionReadChannel) {
		writeError(w, r, http.StatusForbidden, "You don't have access to this channel", nil)
		return
	}

	count, err := p.kvstore.GetGreetingCount(channelID)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "Failed to load greetings", err)
		return
	}

	writeJSON(w, r, http.StatusOK, greetingsResponse{Count: count})
}

func (p *Plugin) handleResetGreetings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	channelID := mux.Vars(r)["channel_id"]

	// Denied attempts are recorded too; they're what an investigation looks for.
	auditRec := p.newAuditRecord(r, "resetGreetingCount")
	defer p.API.LogAuditRec(auditRec)

	if !model.IsValidId(channelID) {
		writeError(w, r, http.StatusBadRequest, "Invalid channel ID", nil)
		return
	}
	model.AddEventParameterToAuditRec(auditRec, "channel_id", channelID)

	// Hiding the button in the webapp is cosmetic; this check is what enforces it.
	if !p.client.User.HasPermissionToChannel(userID, channelID, model.PermissionManageChannelRoles) {
		writeError(w, r, http.StatusForbidden, "Only channel admins can reset greetings", nil)
		return
	}

	if err := p.kvstore.ResetGreetingCount(channelID); err != nil {
		writeError(w, r, http.StatusInternalServerError, "Failed to reset greetings", err)
		return
	}

	auditRec.Success()
	w.WriteHeader(http.StatusNoContent)
}

// newAuditRecord starts a failed record for a state-changing request; call Success once it
// completes. Record IDs and outcomes only, never message content, secrets or config values.
func (p *Plugin) newAuditRecord(r *http.Request, event string) *model.AuditRecord {
	pluginCtx, _ := r.Context().Value(pluginContextKey{}).(*plugin.Context)
	if pluginCtx == nil {
		pluginCtx = &plugin.Context{}
	}
	return plugin.MakeAuditRecordWithContext(event, model.AuditStatusFail, pluginCtx, r.Header.Get("Mattermost-User-ID"), r.URL.Path)
}

// webhookAuthHeader carries the secret generated in the plugin settings. Prefer a header
// over a query parameter: URLs end up in proxy and browser logs.
const webhookAuthHeader = "X-Webhook-Secret"

func (p *Plugin) handleWebhook(w http.ResponseWriter, r *http.Request) {
	secret := p.getConfiguration().WebhookSecret
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
