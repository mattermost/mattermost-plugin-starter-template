# AGENTS.md

Mattermost plugin starter template: a Go server plugin (`server/`) and a React/TypeScript webapp plugin (`webapp/`). Either half can be deleted; the Makefile skips whatever `plugin.json` doesn't declare.

The design, accessibility and testing rules below are adapted from the Mattermost webapp's `webapp/AGENTS.md` and `webapp/STYLE_GUIDE.md` so plugins match the product. When they change upstream, update them here.

## Commands

While iterating, run the narrow command; both `make test` and `make check-style` reinstall Go tools first, and `make test` also rebuilds the bundle and runs every test.

```bash
SKIP_DOCKER_TESTS=1 go test ./server/...       # server unit tests only
go test ./server/command -run TestHelloCommand # single test
cd webapp && npx jest src/manifest.test.tsx    # single webapp test
go vet ./... && golangci-lint run ./...        # server lint
cd webapp && npm run lint && npm run check-types
make test                                      # full check before pushing
make mock                                      # regenerate mocks after changing command.Command or kvstore.KVStore
go run golang.org/x/vuln/cmd/govulncheck@latest ./...  # Go vulnerabilities
cd webapp && npm audit --omit=dev              # vulnerabilities in the shipped bundle
```

`make help` lists the remaining targets (deploy, watch, release bumps, debugger attach).

## Architecture

- `plugin.json` is the source of truth for ID, version, and settings. `make apply` generates `server/manifest.go` and `webapp/src/manifest.ts` from it; both are gitignored, so never edit them. Never edit `dist/`, `server/dist/` or `webapp/dist/` either.
- Project-specific make targets go in `build/custom.mk`.
- Package layout follows the README's "Development guidance": stay in `main` unless there's a reason not to.
- `min_server_version` is 12.0.0. Don't use plugin API methods or host webapp exports newer than that without bumping it.

## Security

The server authenticates requests to `/plugins/{id}/...`, strips spoofed `Mattermost-User-ID` and `Mattermost-Plugin-ID` headers, and checks CSRF tokens on cookie-authenticated non-GET requests. Everything below is the plugin's job. `server/api.go` has a working example of each HTTP rule.

### Authentication and authorization

- Routes under `/api/v1` sit behind `MattermostAuthorizationRequired`. Take the acting user only from the `Mattermost-User-ID` header, never from the body, query or path.
- Authentication only proves *who* the user is. Before reading or changing anything named in the request, check access with `p.client.User.HasPermissionTo` (system), `HasPermissionToTeam` or `HasPermissionToChannel`. For a post, file or thread, check the channel it belongs to. Admin-only actions check `model.PermissionManageSystem`.
- KV records are shared by every user of the plugin. Key per-user data by user ID and check ownership on every read and write.
- Hiding a button in the webapp is cosmetic. Every permission is enforced on the server (see `handleResetGreetings`).
- Slash commands: the server only checks that the user can post in the channel. Admin subcommands check permissions themselves. Reply with ephemeral responses (the default) unless the whole channel should see the result.

### Requests that don't come from a logged-in user

- Routes outside `/api/v1` (e.g. `/webhook`) are unauthenticated. Verify the caller with a secret from a `"type": "generated"`, `"secret": true` setting, compared with `crypto/subtle.ConstantTimeCompare`, or with the provider's HMAC signature over the raw body. Reject everything while the secret is unset. Take secrets from headers, never the query string, since URLs end up in logs.
- Requests from other plugins arrive with an empty `Mattermost-User-ID` and the caller's ID in `Mattermost-Plugin-ID`. Routes meant for them check that ID against an allowlist.
- Files in `public/` are served to anyone at `/plugins/{id}/public/`. Never put secrets or private data there.
- Interactive message actions and dialog submissions are user-controlled input: re-check permissions, and if an action must come from the plugin's own post, look up the post and confirm its author is the plugin's bot.
- Don't trust `post.Props` in hooks like `MessageWillBePosted`; clients can set most props. Identify the plugin's own posts by `post.UserId`.

### Input, output and errors

- Validate IDs with `model.IsValidId`, check enums and lengths, and keep `limitRequestBody`. Paginate list endpoints and cap the page size.
- Anything that changes state uses POST, PUT, PATCH or DELETE. CSRF checks don't apply to GET.
- Return errors through `writeError`: log the cause, send a generic message. Never write `err.Error()` or upstream responses to the client.
- SQL through `p.client.Store` uses placeholders (`?`/squirrel), never string concatenation. Don't pass user input to `os/exec`. Before joining user input into a file path, check it with `filepath.IsLocal`.
- Outbound HTTP: use a client with a timeout. Before fetching a user-supplied URL, allow only `https` and an allowlist of hosts, so the plugin can't be used to reach internal services. Never forward Mattermost session tokens to third parties.

### Secrets and logging

- Credentials, tokens and keys in `plugin.json` settings get `"secret": true`. Never return plugin configuration to non-admins.
- Encrypt per-user OAuth tokens before storing them in the KV store (AES-GCM with a key from a generated secret setting), as the GitHub and Jira plugins do.
- Never log secrets, tokens, query strings or message content. `logRequests` omits the query string on purpose.
- Every state-changing route writes an audit record: `newAuditRecord` + `defer p.API.LogAuditRec`, then `Success()` when it completes. Denied attempts are recorded too. Put object IDs in records, never content, credentials or config values, and only after validating them.

### Cluster and availability

- Use `SetAtomicWithRetries` for KV read-modify-write and `cluster.Schedule` for background jobs; the plugin runs on every node in a cluster.
- Keep hooks that run inline (`MessageWillBePosted`, `MessageWillBeUpdated`, `FileWillBeUploaded`) fast and bounded; they block the user's request.
- `recoverPanic` keeps a handler panic from crashing the plugin. Goroutines you start need their own `recover`.

### Webapp

- Never use `dangerouslySetInnerHTML` (ESLint enforces it). To render Markdown, use the host's sanitizing `window.PostUtils.formatText` + `messageHtmlToComponent`.
- Before putting user-supplied URLs into `href` or `src`, check that the protocol is `http:` or `https:`. ESLint only catches literal `javascript:` URLs.
- Call the plugin server only through `src/client.ts`. `Client4.getOptions` adds the CSRF token; without it the server treats the request as logged out. Don't rely on the deprecated `X-Requested-With` fallback, which `ExperimentalStrictCSRFEnforcement` disables.
- Don't keep tokens or secrets in `localStorage` or Redux, and don't call third-party APIs with credentials from the browser; go through the plugin server.
- Check `event.origin` on every `message` event listener.

### Dependencies

- Dependabot opens monthly update PRs. Run `govulncheck` and `npm audit --omit=dev` before releasing.
- CI runs CodeQL (`.github/workflows/codeql-analysis.yml`) on Go and TypeScript for every PR, and OpenSSF Scorecard weekly. Fix CodeQL alerts rather than adding exclusions to `.github/codeql/codeql-config.yml`.
- Pin third-party GitHub Actions to a commit SHA with the version in a comment, and give workflows `permissions: {}` plus only what each job needs.
- Don't add a dependency for something the host, Compass or the standard library already provides.

## Webapp

### Host-provided packages

- React, ReactDOM, react-intl, Redux, react-redux, react-router-dom and prop-types are webpack `externals` provided by the host. Their `package.json` versions are dev-only and must match the host for `min_server_version`. Don't add another copy to the bundle or move them to `dependencies`.
- `@mattermost/shared/*` imports are loaded from the host at runtime via `window.loadSharedDependency` (see `webpack.config.js`). On 12.0 the host serves `components/button`, `components/emoji`, `components/shortcut_key`, `components/tooltip`, `utils/i18n` and `utils/user_agent`. Add `@mattermost/shared` as a devDependency for the types when you first use it.
- Plugin-facing host APIs live on `window.WebappUtils` (e.g. `modals.openModalById` for the published host modals) and `window.PostUtils`. Don't use `window.Components` or `window.ProductApi`; they're internal and change without notice.
- Don't add react-bootstrap, MUI, Emotion or Sass.

### UI components (Compass)

- Before building UI, check what `@mattermost/compass-ui` exports (package exports and component types). Use a Compass component when one exists, and don't recreate Compass building blocks or invent props.
- Import from subpaths: `import {Button} from '@mattermost/compass-ui/components/button';`. Never import the root barrel, and always use the full package name, never a relative path into `node_modules`.
- `index.tsx` imports `@mattermost/compass-ui/styles` once. Never import `@mattermost/compass-ui/styles/standalone` in a plugin.
- Don't use `@mattermost/compass-proto`; it's an unpublished prototyping package.
- Icons come from `@mattermost/compass-icons/components`. When passing one into a Compass icon slot (`leadingIcon`, `trailingIcon`, `icon`), wrap it in Compass `Icon` and omit `size`: `leadingIcon={<Icon glyph={<RefreshIcon/>}/>}`.
- Tooltips: `WithTooltip` from `@mattermost/shared/components/tooltip`. Emoji: `Emoji` from `@mattermost/shared/components/emoji`.
- Modals and menus: Compass `Modal` and menu components aren't adopted in the host yet. Open published host modals through `window.WebappUtils.modals`, and agree on anything else with design first.
- When migrating an existing plugin, switch to Compass only in the UI you're already changing; don't mass-migrate.

### Styling

- Co-locate styles: `my_component.css` next to `my_component.tsx`, imported by the component.
- Wrap each component in a root class matching its name in PascalCase (`.GreetingsSidebar`), name children with BEM suffixes (`.GreetingsSidebar__count`), and apply modifiers as a separate class.
- Colors come from the host theme variables (`--center-channel-color`, `--button-bg`, `--error-text`, ...); use `rgba(var(--center-channel-color-rgb), 0.72)` for transparency. Spacing, radius, elevation, type and motion come from Compass tokens (`--spacing-m`, `--font-size-100`, `--line-height-100`, `--font-weight-semibold`, ...). Don't hard-code those values. For pills use `var(--radius-pill)`; the webapp's `--radius-full` is `50%`.
- No `!important`, no inline `style={{...}}`, no element selectors.

### Components

- Functional components with hooks. Use `useSelector`/`useDispatch`, not `connect`.
- Split large components into child components or hooks.
- Wrap callbacks passed to children in `useCallback`, and objects/arrays passed to children in `useMemo`.

### Accessibility

- Use semantic elements (`button`, `input`, `h2`, `ul`) instead of `div`/`span` with roles.
- Every interactive element needs a short accessible name (1-3 words). Prefer visible text or a `label` over `aria-label`, and don't repeat the role ("button") in the name.
- Put extra context (help text, errors) in an accessible description with `aria-describedby`.
- Images and icons that carry meaning need alt text or `aria-label`; don't include "image" or "icon" in it. Decorative images get `alt=""`.
- Everything must work with the keyboard alone, following the [ARIA APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/). Focus must be visible when navigating by keyboard, move into modals when they open, and return to the trigger when they close.
- Announce content that updates asynchronously with `aria-live` (see the greeting count).

### Internationalization

- All UI text goes through react-intl. Add strings to `i18n/en.json` only; translators handle the rest.
- Prefer `<FormattedMessage>`; use `useIntl` only when a prop needs a string.
- For mixed formatting or combined strings, use react-intl rich text (`values={{b: (chunks) => <b>{chunks}</b>}}`) instead of concatenating strings or using Markdown.
- Outside React, pass around `MessageDescriptor`s (`{id, defaultMessage}`), not translated strings.

### Tests

- React Testing Library with a real Redux store and `IntlProvider`. Mock `client` at the module boundary (see `greetings_sidebar.test.tsx`).
- Interact with `userEvent` and always `await` it. Use `fireEvent` only where `userEvent` can't work (focus/blur, scroll, image load, fake timers).
- Query by role first: `getByRole` > `getByText`/`getByPlaceholderText` > `getByLabelText`/`getByAltText`/`getByTitle` > `getByTestId`.
- No snapshot tests. Assert what the user sees (`toBeVisible`) and ARIA state (`toHaveAttribute('aria-expanded', 'true')`).
- Only use `act` directly for updates that don't go through an RTL helper.

## Backend tests

- **Integration tests** (`server/integration_test.go`) are the default for anything that goes through the Mattermost server: plugin API calls, hooks, KV persistence, permissions, HTTP routing, activation. `testhelper.Setup(t)` runs Postgres and Mattermost via testcontainers and deploys the bundle from `dist/` (see `go doc github.com/mattermost/mattermost/server/public/pluginapi/testhelper` for fixtures and env vars). Each `Setup` resets the server (~10s), so group related assertions into one test.
- **Plain unit tests** (`server/plugin_test.go`) for pure logic that doesn't touch `p.API`/`p.client`, e.g. drive `ServeHTTP` with `httptest` after `initRouter()`.
- **Mocked-API unit tests** (`server/command/command_test.go`) only for error paths a real server can't easily produce: wrap a `plugintest.API` in `pluginapi.NewClient` and set `api.On(...)` expectations, or use the mockgen mocks in `server/command/mocks` and `server/store/kvstore/mocks` for this plugin's own interfaces. Avoid them otherwise; they encode assumptions about the server rather than testing it.
- Use testify `require`/`assert`, not `t.Fatal`/`t.Errorf`. Make tests table-driven when there's more than one case of the same shape. Don't add new test or mocking libraries.
- Before writing a test, ask whether its failure would mean a real bug. Don't assert implementation details like validation order.
- Security rules get tests: for each new route, cover the denied case (403/401), not just the happy path.

Integration test caveats specific to this repo:
- There's no build tag, so `go test ./...` includes them and fails without Docker. `SKIP_DOCKER_TESTS=1` skips them.
- They test whatever bundle is in `dist/`, which must be linux/amd64 (see the `test` target comment in the Makefile). After a server change, rebuild it before running one by hand:
  ```bash
  make apply server bundle MM_SERVICESETTINGS_ENABLEDEVELOPER=true DEFAULT_GOOS=linux DEFAULT_GOARCH=amd64
  go test ./server -run TestHelloEndpoint
  ```
  This skips the webapp build but requires an existing `webapp/dist`; use `make dist` with the same variables the first time.

## Pull requests and commits

- Commit subject: one succinct line.
- Don't add a `Co-Authored-By` line for the agent.
- In the PR description, call out any change to permissions, authentication, secrets or `min_server_version`.
