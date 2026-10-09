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

Pick the cheapest kind that can catch the bug.

- **Plain unit tests** (`server/plugin_test.go`): construct `Plugin{}`, call `initRouter()`, drive `ServeHTTP` with `httptest`. For handlers and logic that don't touch `p.API`/`p.client`.
- **Mocked-API unit tests** (`server/command/command_test.go`): wrap a `plugintest.API` in `pluginapi.NewClient` and set `api.On(...)` expectations; use the mockgen mocks in `server/command/mocks` for this plugin's own interfaces. For checking which API calls the code makes and exercising error paths that are hard to trigger for real. Brittle: they encode your assumptions about the server and need an expectation per call.
- **Integration tests** (`server/integration_test.go`): `testhelper.Setup(t)` runs Postgres and Mattermost via testcontainers and deploys the bundle from `dist/` (see `go doc github.com/mattermost/mattermost/server/public/pluginapi/testhelper` for fixtures and env vars). For behavior that depends on the real server: auth and routing into the plugin, hooks, KV persistence, permissions, activation. Each `Setup` resets the server (~10s), so keep these few, assert several things per test, and push edge cases down to unit tests.

Integration test caveats specific to this repo:
- There's no build tag, so `go test ./...` includes them and fails without Docker. `SKIP_DOCKER_TESTS=1` skips them.
- They test whatever bundle is in `dist/`, which must be linux/amd64 (see the `test` target comment in the Makefile). After a server change, rebuild it before running one by hand:
  ```bash
  make apply server bundle MM_SERVICESETTINGS_ENABLEDEVELOPER=true DEFAULT_GOOS=linux DEFAULT_GOARCH=amd64
  go test ./server -run TestHelloEndpoint
  ```
  This skips the webapp build but requires an existing `webapp/dist`; use `make dist` with the same variables the first time.
