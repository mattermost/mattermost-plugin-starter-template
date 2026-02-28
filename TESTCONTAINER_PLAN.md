# Plan: Add Reusable Plugin Test Infrastructure to Starter Template

## Context

Mattermost plugins interact heavily with the server — KVStore, Post API, channels, users, hooks — but there's no easy way to integration-test these interactions. Plugin authors are left with either mocking everything via `plugintest.API` (fragile, doesn't test real behavior) or relying on E2E/Playwright tests (slow, hard to write, tests browser-to-server not plugin-to-server).

By adding a testcontainers-based test helper to the starter template, every new plugin gets integration testing infrastructure for free. The helper spins up a real Mattermost server + Postgres in Docker during `go test`, deploys the plugin under test, and provides an authenticated `model.Client4` for exercising the plugin via the REST API.

## Package: `server/testhelper/`

### `mmcontainer.go` — Container lifecycle

Starts Postgres + Mattermost as Docker containers via testcontainers-go:

- **Postgres**: `postgres:15-alpine`, credentials `mmuser`/`mostest`/`mattermost_test`, health check via `pg_isready`
- **Mattermost**: `mattermostdevelopment/mattermost-enterprise-edition:master` (overridable via `MM_TEST_IMAGE` env var), linked to Postgres via Docker network
- Server configured for testing:
  - `MM_PLUGINSETTINGS_ENABLEUPLOADS=true`
  - `MM_PLUGINSETTINGS_AUTOMATICPREPACKAGEDPLUGINS=false`
  - `MM_SERVICESETTINGS_ENABLETESTING=true`
  - `MM_TEAMSETTINGS_ENABLEOPENSERVER=true`
  - `MM_PASSWORDSETTINGS_MINIMUMLENGTH=5`
  - `MM_LOGSETTINGS_CONSOLELEVEL=ERROR`
- Waits for `/api/v4/system/ping` to return OK
- Creates admin user via `POST /api/v4/users`

Containers reuse across tests via `sync.Once` — started once per `go test` invocation, torn down when the binary exits.

### `helper.go` — TestHelper struct and convenience methods

```go
type TestHelper struct {
    ServerURL   string
    AdminClient *model.Client4
    AdminUser   *model.User
    Client      *model.Client4
    User        *model.User
    Team        *model.Team
    Channel     *model.Channel
}
```

Methods:
- **`Setup(t *testing.T) *TestHelper`** — Start containers (once), build and deploy the plugin (once via `sync.Once`), create fresh team/user/channel for test isolation, return helper. Registers cleanup via `t.Cleanup`.
  - Builds the plugin by running `make dist` from the repo root (detected automatically)
  - Finds the tarball in `dist/*.tar.gz`
  - Uploads and enables it via `AdminClient.UploadPlugin` / `EnablePlugin`
  - Polls `GetPluginStatuses` until the plugin is running
- **`CreateUser() *model.User`** — Create a user with random name, add to team.
- **`CreateChannel(channelType string) *model.Channel`** — Create a channel in the test team.
- **`PostAs(user *model.User, channelID, message string) *model.Post`** — Post a message as a specific user.

### `doc.go` — Package documentation

Documents usage, Docker requirement, and environment variables (`MM_TEST_IMAGE`).

## Usage by plugin authors

```go
package main

import (
    "testing"
    "github.com/mattermost/mattermost-plugin-starter-template/server/testhelper"
    "github.com/stretchr/testify/require"
)

func TestPluginActivation(t *testing.T) {
    // Setup starts containers, builds the plugin, deploys it — all automatically.
    th := testhelper.Setup(t)

    // Plugin is already deployed and running. Verify via API.
    statuses, _, err := th.AdminClient.GetPluginStatuses()
    require.NoError(t, err)

    found := false
    for _, s := range statuses {
        if s.PluginId == "com.mattermost.plugin-starter-template" {
            require.Equal(t, model.PluginStateRunning, s.State)
            found = true
        }
    }
    require.True(t, found)
}
```

When a plugin author forks the starter template, they get this test infrastructure automatically. They just change the plugin ID and tarball path.

## Build integration

No build tags or separate test targets. Any test that calls `testhelper.Setup(t)` will use the container — Docker is a prerequisite for running tests, same as any other test dependency.

The existing `make test` target runs all tests including those using the testcontainer. The `Setup` function uses `sync.Once` so the container starts once per `go test` invocation regardless of how many tests use it.

CI already runs tests in the `plugin-ci` reusable workflow which has Docker available on the runner.

## Dependencies to add to go.mod

```
github.com/testcontainers/testcontainers-go v0.35.0
github.com/testcontainers/testcontainers-go/modules/postgres v0.35.0
```

## Files to create/modify

1. **`server/testhelper/mmcontainer.go`** — Create: Postgres + Mattermost container setup
2. **`server/testhelper/helper.go`** — Create: TestHelper, Setup, data factories
3. **`server/testhelper/doc.go`** — Create: Package documentation
4. **`server/plugin_test.go`** — Modify: Replace mock-based test with real server test using testhelper
5. **`go.mod`** / **`go.sum`** — Modified by `go get` for testcontainers deps

## Verification

1. `go vet ./server/...`
2. `go test ./server/...` — All tests pass (Docker required). The testhelper automatically runs `make dist` and deploys the plugin on first use.
3. The smoke test in `plugin_test.go` verifies the plugin activates and responds to HTTP requests against a real Mattermost server.
