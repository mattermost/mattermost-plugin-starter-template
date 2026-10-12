# AGENTS.md

Mattermost plugin starter template: a Go server plugin (`server/`) and a React/TypeScript webapp plugin (`webapp/`). Either half can be deleted; the Makefile skips whatever `plugin.json` doesn't declare.

## Commands

While iterating, run the narrow command; both `make test` and `make check-style` reinstall Go tools first, and `make test` also rebuilds the bundle and runs every test.

```bash
SKIP_DOCKER_TESTS=1 go test ./server/...       # server unit tests only
go test ./server/command -run TestHelloCommand # single test
cd webapp && npx jest src/manifest.test.tsx    # single webapp test
go vet ./... && golangci-lint run ./...        # server lint
cd webapp && npm run lint && npm run check-types
make test                                      # full check before pushing
make mock                                      # regenerate server/command/mocks after changing command.Command
```

`make help` lists the remaining targets (deploy, watch, release bumps, debugger attach).

## Architecture

- `plugin.json` is the source of truth for ID, version, and settings. `make apply` generates `server/manifest.go` and `webapp/src/manifest.ts` from it; both are gitignored, so never edit them.
- Project-specific make targets go in `build/custom.mk`.
- Package layout follows the README's "Development guidance": stay in `main` unless there's a reason not to.

## Backend tests

- **Integration tests** (`server/integration_test.go`) are the default for anything that goes through the Mattermost server: plugin API calls, hooks, KV persistence, permissions, HTTP routing, activation. `testhelper.Setup(t)` runs Postgres and Mattermost via testcontainers and deploys the bundle from `dist/` (see `go doc github.com/mattermost/mattermost/server/public/pluginapi/testhelper` for fixtures and env vars). Each `Setup` resets the server (~10s), so group related assertions into one test.
- **Plain unit tests** (`server/plugin_test.go`) for pure logic that doesn't touch `p.API`/`p.client`, e.g. drive `ServeHTTP` with `httptest` after `initRouter()`.
- **Mocked-API unit tests** (`server/command/command_test.go`) only for error paths a real server can't easily produce: wrap a `plugintest.API` in `pluginapi.NewClient` and set `api.On(...)` expectations, or use the mockgen mocks in `server/command/mocks` for this plugin's own interfaces. Avoid them otherwise; they encode assumptions about the server rather than testing it.

Integration test caveats specific to this repo:
- There's no build tag, so `go test ./...` includes them and fails without Docker. `SKIP_DOCKER_TESTS=1` skips them.
- They test whatever bundle is in `dist/`, which must be linux/amd64 (see the `test` target comment in the Makefile). After a server change, rebuild it before running one by hand:
  ```bash
  make apply server bundle MM_SERVICESETTINGS_ENABLEDEVELOPER=true DEFAULT_GOOS=linux DEFAULT_GOARCH=amd64
  go test ./server -run TestHelloEndpoint
  ```
  This skips the webapp build but requires an existing `webapp/dist`; use `make dist` with the same variables the first time.
