# AGENTS.md

Guidance for coding agents working in this repository.

This is the Mattermost plugin starter template. It has a Go server plugin (`server/`), a React/TypeScript webapp plugin (`webapp/`), and shared build tooling (`build/`, `Makefile`). Either half can be deleted. The Makefile detects which halves exist through `HAS_SERVER`/`HAS_WEBAPP` (derived from `plugin.json`) and skips the missing parts.

## Commands

```bash
make                  # check-style + test + dist
make dist             # build server binaries + webapp, bundle to dist/<id>-<version>.tar.gz
make check-style      # eslint + tsc (webapp), go vet + golangci-lint (server)
make test             # builds a linux/amd64 bundle, then runs all Go tests (gotestsum) + jest
make deploy           # build and install into a running server (local mode socket or MM_* creds)
make watch            # rebuild/redeploy the webapp on change
make mock             # regenerate server/command/mocks with mockgen
make help             # list all targets
```

Running a single test:

```bash
go test ./server/command -run TestHelloCommand          # unit test, no setup needed
make dist MM_SERVICESETTINGS_ENABLEDEVELOPER=true DEFAULT_GOOS=linux DEFAULT_GOARCH=amd64
go test ./server -run TestHelloEndpoint                 # integration test, needs Docker + bundle
SKIP_DOCKER_TESTS=1 go test ./...                       # run everything except integration tests
cd webapp && npx jest src/manifest.test.tsx             # single webapp test
```

`MM_DEBUG=1` produces unminified webapp and debug-friendly server builds. `MM_SERVICESETTINGS_ENABLEDEVELOPER=true` builds the server for a single `DEFAULT_GOOS`/`DEFAULT_GOARCH` pair instead of every platform.

## Architecture

- `plugin.json` is the source of truth for the plugin's ID, version, and settings schema. `make apply` (run by `build/manifest`) generates `server/manifest.go` and `webapp/src/manifest.ts` from it. Both files are gitignored, so edit `plugin.json` and never the generated files. If no version is set, it is derived from git tags at build time.
- `server/plugin.go` holds the `Plugin` struct. `OnActivate` wires up the `pluginapi.Client`, the KV store wrapper (`server/store/kvstore`), the slash command handler (`server/command`), the HTTP router, and a cluster-wide background job (`server/job.go`, scheduled via `pluginapi/cluster`).
- `server/api.go`: `ServeHTTP` delegates to a gorilla/mux router mounted under `/plugins/<id>/`. The `MattermostAuthorizationRequired` middleware trusts the `Mattermost-User-ID` header, which the Mattermost server sets only after it authenticates the request.
- `server/configuration.go` uses the standard copy-on-write configuration pattern behind `configurationLock`.
- `build/pluginctl` is the tool behind `deploy`, `disable`, `logs`, and similar targets. Project-specific make targets go in `build/custom.mk`.
- Following the README: keep code in the `main` package unless it introduces a new interface or wraps an upstream integration.
- CI (`.github/workflows/ci.yml`) delegates to the shared `mattermost/actions-workflows` plugin-ci workflow, which runs `make test-ci`.

## Backend tests

The server has three kinds of tests. Use the cheapest one that can catch the bug you care about.

### 1. Plain unit tests (`httptest`, no mocks)

Example: `server/plugin_test.go` (`TestServeHTTP`). These build a `Plugin{}` directly, call `initRouter()`, and drive `ServeHTTP` with `httptest`. They fake the auth header that the server would normally set.

Use them for pure logic and HTTP handlers that don't touch `p.API` or `p.client`, such as routing, request parsing, response formatting, and middleware. They are fast and need nothing installed.

### 2. Unit tests with mocked plugin API (`plugintest` / mockgen)

Example: `server/command/command_test.go`. These create a `plugintest.API` (testify mock) and wrap it in `pluginapi.NewClient`, then set `api.On(...)` expectations for each server call the code makes. For mocking this plugin's own interfaces (e.g. `command.Command`), use the mockgen mocks in `server/command/mocks`, and regenerate them with `make mock` after changing the interface.

Use them when code calls the Mattermost API and you want to check what it calls and how it handles specific return values, especially error paths that are hard to trigger on a real server. Keep in mind that these tests only verify your assumptions about the server's behavior. Every API call needs an expectation, which makes them brittle when the implementation changes.

### 3. Integration tests (`pluginapi/testhelper`)

Example: `server/integration_test.go` (`TestHelloEndpoint`). `testhelper.Setup(t)` uses testcontainers to start Postgres and a real Mattermost server, deploys the bundle from `dist/`, and provides `th.AdminClient`, `th.Client`, `th.Team`, `th.User`, and `th.Channel`. Helpers include `CreateUser`, `CreateChannel`, and `PostAs`.

Use them for behavior that depends on the real server: authentication and request routing into the plugin, hooks firing, KV store persistence, permissions, and plugin activation and configuration. Also use them for anything where a mock could silently disagree with the real API.

Requirements and caveats:
- Docker must be running. Without it, the test fails instead of skipping. Set `SKIP_DOCKER_TESTS=1` to skip.
- Exactly one bundle must exist in `dist/*.tar.gz`, built for linux/amd64, because the Mattermost image only ships for that platform. `make test` handles this. When running `go test` by hand, rebuild the bundle after any server change, or you will be testing stale code.
- These tests have no build tag, so a plain `go test ./...` includes them.
- Every `Setup` call resets the database and restarts the server, which costs about 10 seconds per test. Prefer a few integration tests with several assertions each over many small ones.
- `MM_TEST_IMAGE` overrides the server image (default `mattermost/mattermost-enterprise-edition:latest`).

### Rule of thumb

Start with a plain unit test. Add `plugintest` mocks when the code under test calls the plugin API and you need to control its responses. Write an integration test when correctness depends on how the real server behaves, or when the test should cover a full user-visible flow end to end. Keep integration tests few and focused, and push edge cases down to unit tests.
