# AGENTS.md (server)

`server/api.go` has a working example of each HTTP rule below; copy its patterns.

## Security

The server authenticates `/plugins/{id}/...` requests, strips spoofed `Mattermost-User-ID` and `Mattermost-Plugin-ID` headers, and checks CSRF tokens on cookie-authenticated non-GET requests. The rest is the plugin's job.

- Take the acting user only from the `Mattermost-User-ID` header, never from the body, query or path.
- Before reading or changing anything named in a request, check `p.client.User.HasPermissionTo` (system), `HasPermissionToTeam` or `HasPermissionToChannel`. For a post, file or thread, check its channel. Check the permission for the resource being changed: channel admin actions use a channel permission (`handleResetGreetings` uses `model.PermissionManageChannelRoles`), team admin actions a team permission, and only system-wide actions `model.PermissionManageSystem`.
- KV records are shared by all users of the plugin: key per-user data by user ID and check ownership on every access.
- Slash commands: the server only checks that the user can post in the channel. Admin subcommands check permissions themselves; replies stay ephemeral unless the channel should see them.
- Routes outside `/api/v1` are unauthenticated. Verify callers with a `"type": "generated"`, `"secret": true` setting compared via `crypto/subtle.ConstantTimeCompare`, or the provider's HMAC over the raw body. Reject everything while the secret is unset; read it from a header, never the query string.
- Requests from other plugins have an empty `Mattermost-User-ID` and the caller in `Mattermost-Plugin-ID`; routes for them check an allowlist.
- `public/` is served to anyone. Never put secrets or private data there.
- Interactive actions, dialog submissions and `post.Props` in hooks are user-controlled. Re-check permissions, and identify the plugin's own posts by `post.UserId`, not by a prop.
- Validate IDs with `model.IsValidId`, check enums and lengths, keep `limitRequestBody`, and cap page sizes.
- Errors go through `writeError` (logs the cause, returns a generic message).
- SQL uses placeholders; no user input in `os/exec`; check user-supplied path parts with `filepath.IsLocal`.
- Outbound HTTP uses a client with a timeout. User-supplied URLs must be `https` and on an allowlist, so the plugin can't reach internal services. Never forward Mattermost tokens.
- Settings holding credentials get `"secret": true`. Encrypt per-user OAuth tokens before storing them (AES-GCM, key from a generated setting). Never return plugin config to non-admins.
- State-changing routes write an audit record: `newAuditRecord` + `defer p.API.LogAuditRec`, then `Success()`. Record denied attempts too, and only validated IDs, never content, credentials or config.
- KV read-modify-write uses `SetAtomicWithRetries`; background jobs use `cluster.Schedule`. The plugin runs on every cluster node.
- Inline hooks (`MessageWillBePosted`, `MessageWillBeUpdated`, `FileWillBeUploaded`) block the user's request: keep them fast. Goroutines you start need their own `recover`.
- Before releasing: `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`.

## Tests

- **Integration tests** (`server/integration_test.go`) are the default for anything that goes through the Mattermost server: plugin API calls, hooks, KV persistence, permissions, HTTP routing, activation. `testhelper.Setup(t)` runs Postgres and Mattermost via testcontainers and deploys the bundle from `dist/` (see `go doc github.com/mattermost/mattermost/server/public/pluginapi/testhelper`). Each `Setup` resets the server (~10s), so group related assertions into one test.
- **Plain unit tests** (`server/plugin_test.go`) for logic that doesn't need a real server, e.g. drive `ServeHTTP` with `httptest` after `initRouter()`.
- **Mocked-API unit tests** (`server/command/command_test.go`) only for error paths a real server can't easily produce: `plugintest.API` wrapped in `pluginapi.NewClient`, or the mockgen mocks in `server/command/mocks` and `server/store/kvstore/mocks`. They encode assumptions about the server, so avoid them otherwise.
- Every new route gets a test for the denied case (401/403), not just the happy path.
- testify `require`/`assert`, table-driven when there's more than one case of the same shape, no new test or mocking libraries. Don't assert implementation details like validation order.

Integration test caveats:
- There's no build tag, so `go test ./...` includes them and fails without Docker. `SKIP_DOCKER_TESTS=1` skips them. `MM_TEST_IMAGE` overrides the default `mattermost-enterprise-edition:latest` image.
- They test whatever bundle is in `dist/`, which must be linux/amd64. After a server change, rebuild before running one by hand:
  ```bash
  make apply server bundle MM_SERVICESETTINGS_ENABLEDEVELOPER=true DEFAULT_GOOS=linux DEFAULT_GOARCH=amd64
  go test ./server -run TestHelloEndpoint
  ```
  This needs an existing `webapp/dist`; use `make dist` with the same variables the first time.
