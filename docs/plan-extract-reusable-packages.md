# Plan: Extract Reusable Code from Plugin Starter Template

## Context

Plugins created from the Mattermost Plugin Starter Template get a snapshot of template code at creation time. Improvements to the template (build tools, test infrastructure, build config, webapp tooling) never reach existing plugins. The goal is to extract reusable code into shared, importable packages so plugins can receive updates via standard dependency management (`go get -u`, `npm update`).

Evidence of drift: comparing the starter template against `mattermost-plugin-github` reveals diverged manifest tools (different package names, different Go API usage), 42+ lines of Makefile diff, missing `pluginctl/logs.go` (~180 lines), and different linter versions.

## Strategy: Hybrid Extraction

Four focused packages, each using the update mechanism natural to its ecosystem. Implementation is ordered by value and ease of extraction.

---

## Package 1: Go Test Library (Priority: Highest)

**Repo**: `github.com/mattermost/mattermost-plugin-test` (new)

### What to Extract

All of `server/testhelper/` — container lifecycle, database reset, test data helpers.

### API Design

```go
package mmtest

// Setup starts containers (once per binary), resets DB, deploys plugin, creates test data.
func Setup(t *testing.T, opts ...Option) *TestHelper

// Options (functional options pattern):
func WithPluginBundle(path string) Option   // path to .tar.gz (required unless WithoutPlugin)
func WithPluginID(id string) Option         // plugin ID (required unless WithoutPlugin)
func WithMMImage(image string) Option       // override MM image (default: env MM_TEST_IMAGE)
func WithoutPlugin() Option                 // skip plugin deployment

// Utility functions for callers who want auto-discovery:
func FindPluginBundle(dir string) (string, error)  // glob dist/*.tar.gz
func ReadPluginID(manifestPath string) (string, error)  // parse plugin.json

// TestHelper has the same exported fields/methods as today:
// ServerURL, AdminClient, AdminUser, Client, User, Team, Channel
// CreateUser(), CreateChannel(), PostAs()
```

### Key Design Decisions

- **Plugin ID and bundle path are explicit parameters**, not auto-discovered from filesystem. The current `findRepoRoot()`/`getPluginID()` approach couples the library to a specific directory layout. Utility functions are provided for callers who want the convenience.
- **Container sharing via `sync.Once` preserved** — first `Setup()` call configures; subsequent calls reuse.
- **Hardcoded constants** (pgUser, admin creds) exported as `DefaultPassword`, `DefaultAdminUsername` etc. so callers can reference them.

### Changes to Starter Template

`server/testhelper/` becomes a thin wrapper (~20 lines) that calls `mmtest.Setup()` with `WithPluginID`/`WithPluginBundle` populated by the existing `findRepoRoot`/`getPluginID` logic. `go.mod` adds the new dependency; direct `testcontainers-go` deps become transitive.

### Migration for Existing Plugins

1. `go get github.com/mattermost/mattermost-plugin-test@latest`
2. Replace hand-rolled testhelper code with `mmtest.Setup(t, mmtest.WithPluginID("..."), mmtest.WithPluginBundle("dist/..."))`
3. No Makefile changes needed — `BUILD_FOR_TEST=1` and `make test` -> `make dist` chain stays the same

---

## Package 2: CLI Tool (Priority: High)

**Repo**: `github.com/mattermost/mattermost-plugin-build` (existing, currently near-empty)

### What to Extract

Merge `build/manifest/main.go` (222 lines) and `build/pluginctl/` (main.go + logs.go, ~370 lines) into a single CLI binary.

### CLI Commands

```
mmplugin manifest id [--manifest=plugin.json]
mmplugin manifest version [--manifest=plugin.json] [--hash-short=...] [--tag-latest=...] [--tag-current=...]
mmplugin manifest has-server [--manifest=plugin.json]
mmplugin manifest has-webapp [--manifest=plugin.json]
mmplugin manifest apply [--manifest=plugin.json] [--server-out=server/manifest.go] [--webapp-out=webapp/src/manifest.ts] [--go-package=main]
mmplugin manifest dist [--manifest=plugin.json] [--dist-dir=dist]
mmplugin manifest check [--manifest=plugin.json]

mmplugin deploy <plugin-id> <bundle-path>
mmplugin disable <plugin-id>
mmplugin enable <plugin-id>
mmplugin reset <plugin-id>
mmplugin logs <plugin-id>
mmplugin logs watch <plugin-id>
```

Install: `go install github.com/mattermost/mattermost-plugin-build/cmd/mmplugin@latest`

### Key Design Decisions

- **All paths are CLI flags with sensible defaults** matching current behavior (e.g., `--manifest` defaults to `plugin.json` in cwd).
- **`--go-package` flag** on `manifest apply` — existing plugins use `package plugin` while the starter template uses `package main`. This flag resolves the divergence.
- **Auto-detects git state** for version info, removing the need for ldflags injection from `setup.mk`. Still supports ldflags for backward compatibility.
- **Uses cobra** for command structure (standard in the Mattermost ecosystem).

### Changes to Starter Template

`build/setup.mk` simplifies — instead of compiling tools inline, it installs/invokes `mmplugin`:

```make
MMPLUGIN ?= $(shell command -v mmplugin 2>/dev/null)
ifeq ($(MMPLUGIN),)
    $(shell $(GO) install github.com/mattermost/mattermost-plugin-build/cmd/mmplugin@v0.1.0)
    MMPLUGIN = $(GOBIN)/mmplugin
endif

PLUGIN_ID ?= $(shell $(MMPLUGIN) manifest id)
PLUGIN_VERSION ?= $(shell $(MMPLUGIN) manifest version)
```

`build/manifest/` and `build/pluginctl/` directories are deleted.

### Migration for Existing Plugins

1. `go install github.com/mattermost/mattermost-plugin-build/cmd/mmplugin@latest`
2. Copy updated `build/setup.mk` from starter template
3. Update Makefile references from `./build/bin/manifest` to `$(MMPLUGIN) manifest` and `./build/bin/pluginctl` to `$(MMPLUGIN)`
4. Delete `build/manifest/` and `build/pluginctl/`
5. Plugins using `package plugin` in manifest.go: add `--go-package=plugin` to the apply call

---

## Package 3: Shared Makefile (Priority: High)

**Repo**: Same as Package 2 (`github.com/mattermost/mattermost-plugin-build`)

### What to Extract

~300 lines of generic targets from the current 452-line Makefile into `plugin.mk`:

- Build: `server`, `webapp`, `bundle`, `dist`, `apply`, `manifest-check`
- Test: `test`, `test-ci`, `coverage`
- Lint: `check-style`, `install-go-tools`
- Deploy: `deploy`, `deploy-from-watch`, `watch`
- Control: `enable`, `disable`, `reset`, `kill`, `logs`, `logs-watch`
- Debug: `attach`, `attach-headless`, `detach`
- Release: `patch`, `minor`, `major` and `-rc` variants
- Utility: `clean`, `mock`, `i18n-extract`

### Distribution Mechanism

Use Go module cache — every plugin already has Go installed:

```make
PLUGIN_BUILD_VERSION ?= v0.1.0
PLUGIN_BUILD_PKG := github.com/mattermost/mattermost-plugin-build@$(PLUGIN_BUILD_VERSION)
PLUGIN_BUILD_DIR := $(shell $(GO) env GOPATH)/pkg/mod/$(PLUGIN_BUILD_PKG)
ifeq ($(wildcard $(PLUGIN_BUILD_DIR)/plugin.mk),)
    $(shell $(GO) mod download $(PLUGIN_BUILD_PKG))
endif
include $(PLUGIN_BUILD_DIR)/plugin.mk
```

### Key Design Decisions

- **All variables use `?=`** (conditional assignment) so plugins can override anything before the include.
- **`build/custom.mk`** is preserved as the local customization point, loaded after `plugin.mk`.
- **`setup.mk` logic absorbed into `plugin.mk`** — tool installation, plugin metadata extraction.

### Changes to Starter Template

`Makefile` shrinks from 452 lines to ~30 lines (variable definitions + include). `build/setup.mk` is deleted (absorbed into `plugin.mk`).

### Migration for Existing Plugins

1. Set `PLUGIN_BUILD_VERSION` in Makefile
2. Replace Makefile with thin version from updated starter template
3. Move any custom targets into `build/custom.mk`
4. Delete `build/setup.mk`, `build/manifest/`, `build/pluginctl/`
5. Verify with `make check-style test dist`

**Risk**: This is the highest-risk migration because many plugins have heavily customized Makefiles. Mitigation: `?=` variables mean any existing overrides take precedence, and adoption is opt-in.

---

## Package 4: Webapp Build Package (Priority: Medium)

**Package**: `@mattermost/plugin-webpack-config` on npm (new)

### What to Extract

- `webpack.config.js` (99 lines) -> `createWebpackConfig(pluginId, options)` function
- `babel.config.js` (42 lines) -> exported base config
- `tsconfig.json` (34 lines) -> `tsconfig.base.json` to extend from
- Jest config from `package.json` -> Jest preset
- `webapp/src/types/mattermost-webapp/index.d.ts` (1133 lines) -> `types/plugin-registry.d.ts`

### Key Design Decisions

- **`createWebpackConfig` returns a plain object** — plugins can spread/override specific fields for customization.
- **Externals list is fixed** (React, ReactDOM, Redux, etc.) — these are provided by the Mattermost webapp runtime and cannot be changed without breaking.
- **PluginRegistry types** are the most maintenance-heavy artifact. They define the Mattermost webapp plugin API and are currently copy-pasted into every plugin. Extracting them into this package (or eventually `@mattermost/types`) eliminates 1133 lines of drift-prone code per plugin.
- **Watch-mode deploy hook is opt-in** via `options.watchDeploy`.

### Changes to Starter Template

- `webpack.config.js`: 99 -> ~8 lines (calls `createWebpackConfig`)
- `babel.config.js`: 42 -> ~3 lines (re-exports base)
- `tsconfig.json`: extends `@mattermost/plugin-webpack-config/tsconfig.base.json`
- `index.d.ts`: deleted (comes from package)
- Several devDependencies become transitive and can be removed

### Migration for Existing Plugins

1. `npm install --save-dev @mattermost/plugin-webpack-config@latest`
2. Replace config files with thin versions
3. Delete local `types/mattermost-webapp/index.d.ts`
4. Run `npm run build && npm run test` to verify

**Risk**: Webapp configs tend to accumulate plugin-specific tweaks. Plugins should diff their configs against the starter template before adopting.

---

## Implementation Sequence

| Phase | Package | Estimated Effort | Dependencies |
|-------|---------|-----------------|--------------|
| 1 | Go Test Library | 2-3 weeks | None |
| 2 | CLI Tool | 2-3 weeks | None (parallel with Phase 1) |
| 3 | Shared Makefile | 2 weeks | Phase 2 (uses mmplugin) |
| 4 | Webapp Build Package | 3-4 weeks | None (parallel with Phase 3) |

Phases 1 and 2 can proceed in parallel. Phase 3 depends on Phase 2. Phase 4 can proceed in parallel with Phase 3.

## Risks and Mitigations

1. **Breaking existing plugins**: All packages are opt-in. No forced migration.
2. **Version coordination**: Test library, CLI, and Makefile all depend on `mattermost/server/public/model`. Use dependabot; release together when server bumps.
3. **Two-repo coordination**: Changes spanning shared package + starter template require two PRs. Mitigated by semver pinning.
4. **Makefile fragility**: Make's `include` is whitespace-sensitive. Mitigated by using Go module cache (reliable) and clear error messages.
5. **Webapp config divergence**: Some plugins have heavily customized webpack configs. `createWebpackConfig` returns a mutable object; plugins with extreme customization can decline adoption.

## Verification

For each package, before tagging a release:
1. Update the starter template to consume the new package
2. Run `make check-style test dist` (full build pipeline)
3. Run integration tests (`make test` with Docker)
4. Test `make deploy` against a local Mattermost server
5. Test migration on at least one real plugin (e.g., `mattermost-plugin-github`)
