# Plan: Add Reusable Plugin Test Infrastructure to Starter Template

## Context

Mattermost plugins interact heavily with the server — KVStore, Post API, channels, users, hooks — but there's no easy way to integration-test these interactions. Plugin authors are left with either mocking everything via `plugintest.API` (fragile, doesn't test real behavior) or relying on E2E/Playwright tests (slow, hard to write, tests browser-to-server not plugin-to-server).

By adding a testcontainers-based test helper to the starter template, every new plugin gets integration testing infrastructure for free. The helper spins up a real Mattermost server + Postgres in Docker during `go test`, deploys the plugin under test, and provides an authenticated `model.Client4` for exercising the plugin via the REST API.

## Package: `server/testhelper/`

### `mmcontainer.go` — Container lifecycle

Starts Postgres + Mattermost as Docker containers via testcontainers-go:

- **Postgres**: `postgres:15-alpine`, credentials `mmuser`/`mostest`/`mattermost_test`, health check via `pg_isready`
- **Mattermost**: `mattermost/mattermost-enterprise-edition:latest` by default (currently v11.4.x), linked to Postgres via Docker network
  - **Image override**: The `MM_TEST_IMAGE` env var overrides the full image reference (e.g., `MM_TEST_IMAGE=mattermostdevelopment/mattermost-enterprise-edition:master` for bleeding-edge, or `mattermost/mattermost-enterprise-edition:11.3` for an older release, or a specific commit SHA from the `mattermostdevelopment` repo)
- Server configured for testing:
  - `MM_PLUGINSETTINGS_ENABLEUPLOADS=true`
  - `MM_PLUGINSETTINGS_AUTOMATICPREPACKAGEDPLUGINS=false`
  - `MM_SERVICESETTINGS_ENABLETESTING=true`
  - `MM_TEAMSETTINGS_ENABLEOPENSERVER=true`
  - `MM_PASSWORDSETTINGS_MINIMUMLENGTH=5`
  - `MM_LOGSETTINGS_CONSOLELEVEL=ERROR`

#### Readiness check

Uses testcontainers-go `wait.ForHTTP("/api/v4/system/ping")` strategy with:
- **Timeout**: 120 seconds (Mattermost can take 15-30s to start, plus migration time on first boot)
- **Poll interval**: 1 second
- **Expected status**: 200

#### Admin user creation

After the container is ready, create the admin user by executing `mmctl` **inside** the Mattermost container via `container.Exec()`:

```
mmctl --local user create \
  --email admin@example.com \
  --username admin \
  --password Password1! \
  --system-admin \
  --email-verified
```

The `--local` flag connects via the Unix socket inside the container, bypassing authentication — which is necessary because no users exist yet. The `mmctl` binary is bundled in the official Mattermost Docker images.

#### Container reuse

Containers are started once per `go test` invocation via `sync.Once`.

#### Teardown

- **Ryuk reaper** (testcontainers-go default): A sidecar container that automatically removes test containers when the parent process exits, even on crashes or `kill -9`. This is the primary safety net.
- **Explicit cleanup**: `testcontainers.CleanupContainer(t, container)` is registered via `t.Cleanup()` on the *first* test that triggers `Setup()`. This provides orderly shutdown when tests complete normally.
- **Documentation**: A note in `doc.go` will mention that Ryuk must not be disabled (`TESTCONTAINERS_RYUK_DISABLED=false` is the default) to prevent container leaks in CI.

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

- **`Setup(t *testing.T) *TestHelper`** — Start containers (once), deploy the plugin (once via `sync.Once`), **reset server state**, create fresh team/user/channel for test isolation, return helper.
  - **Plugin deployment** (once): Finds the pre-built tarball in `dist/*.tar.gz` (built by `make dist` which runs as a Makefile dependency before tests — see Build Integration below). Uploads and enables it via `AdminClient.UploadPlugin` / `EnablePlugin`. Polls `GetPluginStatuses` with a 30-second timeout until the plugin reports `PluginStateRunning`.
  - **State reset** (every test): Executes `mattermost db reset --confirm` inside the Mattermost container via `container.Exec()`. This truncates all tables except `db_migrations` (using the same `DropAllTables()` logic the Mattermost server uses internally — a single dynamic `TRUNCATE TABLE ... CASCADE` across all public-schema tables). After reset, re-creates the admin user via `mmctl --local` and re-deploys the plugin. This guarantees each test starts with a clean database — no KVStore leakage, no leftover users/channels/posts from prior tests.
  - Creates a fresh team, regular user, and channel for the test.
  - Registers `t.Cleanup` to handle per-test teardown (no DB-level cleanup needed since next test's `Setup` resets anyway).

- **`CreateUser() *model.User`** — Create a user with random name, add to team.
- **`CreateChannel(channelType string) *model.Channel`** — Create a channel in the test team.
- **`PostAs(user *model.User, channelID, message string) *model.Post`** — Post a message as a specific user.

#### Parallel test safety

Tests using `TestHelper` **can** call `t.Parallel()` only if they do not share state beyond what `Setup` creates. Since each call to `Setup` resets the DB, parallel tests within the same `go test` binary would conflict. If parallel execution is needed in the future, the approach should shift to one container per parallel group (similar to Mattermost server's `TestPool` pattern). For now, tests run sequentially which is the safe default for plugin integration tests.

### `testmain_test.go` — TestMain for graceful Docker detection

```go
func TestMain(m *testing.M) {
    if os.Getenv("SKIP_DOCKER_TESTS") != "" {
        fmt.Println("Skipping integration tests (SKIP_DOCKER_TESTS is set)")
        os.Exit(0)
    }

    // Check Docker availability
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    provider, err := testcontainers.ProviderDocker.GetProvider()
    if err != nil {
        fmt.Printf("Skipping integration tests: Docker not available: %v\n", err)
        os.Exit(0)
    }
    defer provider.Close()

    if _, err := provider.Health(ctx); err != nil {
        fmt.Printf("Skipping integration tests: Docker not healthy: %v\n", err)
        os.Exit(0)
    }

    os.Exit(m.Run())
}
```

This ensures:
- Tests skip gracefully (exit 0, not failure) when Docker is unavailable
- Developers without Docker can still run `make test` — integration tests simply don't execute
- The `SKIP_DOCKER_TESTS` env var provides an explicit opt-out

### `doc.go` — Package documentation

Documents usage, Docker requirement, Ryuk reaper dependency, and environment variables:
- `MM_TEST_IMAGE` — Override Mattermost Docker image (default: `mattermost/mattermost-enterprise-edition:latest`)
- `SKIP_DOCKER_TESTS` — Set to any value to skip integration tests

## Usage by plugin authors

```go
package main

import (
    "testing"
    "github.com/mattermost/mattermost-plugin-starter-template/server/testhelper"
    "github.com/stretchr/testify/require"
)

func TestPluginActivation(t *testing.T) {
    // Setup starts containers (once), resets DB, deploys plugin, creates test data.
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

When a plugin author forks the starter template, they get this test infrastructure automatically. They just change the plugin ID in `plugin.json`.

## Build integration

### Makefile change: make `dist` a prerequisite of `test`

The current `test` target:
```makefile
test: apply webapp/node_modules install-go-tools
```

Updated to:
```makefile
test: dist install-go-tools
```

Since `dist` already depends on `apply`, `server`, `webapp`, and `bundle`, this ensures the plugin tarball (`dist/*.tar.gz`) is built **before** any Go tests run. The `testhelper.Setup()` function simply locates the pre-built tarball — it never shells out to `make`.

The `webapp/node_modules` dependency is dropped from `test` because `dist` → `webapp` already handles it.

**Trade-off**: `make test` now builds the full plugin before running tests, adding ~10-30s to the test cycle. This is acceptable because:
1. Integration tests require the built bundle anyway
2. The `dist` target is a no-op if sources haven't changed (Make dependency tracking)
3. It eliminates the fragility of shelling out to `make` from inside `go test`

### No build tags

No build tags or separate test targets. Any test that calls `testhelper.Setup(t)` will use the container. The `TestMain` Docker check ensures graceful degradation when Docker is unavailable.

CI already runs tests in the `plugin-ci` reusable workflow which has Docker available on the runner.

## Dependencies to add to go.mod

```
github.com/testcontainers/testcontainers-go v0.35.0
github.com/testcontainers/testcontainers-go/modules/postgres v0.35.0
```

(Verify latest version at implementation time — the library releases frequently.)

## Files to create/modify

1. **`server/testhelper/mmcontainer.go`** — Create: Postgres + Mattermost container setup, readiness check, admin creation via mmctl
2. **`server/testhelper/helper.go`** — Create: TestHelper, Setup (with DB reset via `mattermost db reset --confirm`), data factories
3. **`server/testhelper/testmain_test.go`** — Create: TestMain with Docker availability check
4. **`server/testhelper/doc.go`** — Create: Package documentation
5. **`server/plugin_test.go`** — Modify: Replace mock-based test with real server test using testhelper
6. **`Makefile`** — Modify: Make `dist` a prerequisite of `test`
7. **`go.mod`** / **`go.sum`** — Modified by `go get` for testcontainers deps

## Verification

1. `go vet ./server/...`
2. `make test` — All tests pass (Docker required). The `dist` target builds the plugin bundle, then `go test` deploys it into the container.
3. `SKIP_DOCKER_TESTS=1 make test` — Integration tests are skipped gracefully, other tests still pass.
4. The smoke test in `plugin_test.go` verifies the plugin activates and responds to HTTP requests against a real Mattermost server.
