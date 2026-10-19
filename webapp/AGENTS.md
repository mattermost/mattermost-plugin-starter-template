# AGENTS.md (webapp)

Adapted from the Mattermost webapp's `webapp/AGENTS.md` and `STYLE_GUIDE.md` so plugins match the product. ESLint and `tsc` enforce formatting, imports, `dangerouslySetInnerHTML` and `javascript:` URLs; the rules below are the ones they can't.

## Host-provided packages

- React, ReactDOM, react-intl, Redux, react-redux, react-router-dom and prop-types are webpack `externals` from the host. Their versions in `package.json` are dev-only and must match the host for `min_server_version`. Never bundle a second copy or move them to `dependencies`.
- `@mattermost/shared/*` loads from the host via `window.loadSharedDependency` (see `webpack.config.js`). On 12.0 the host serves `components/button`, `components/emoji`, `components/shortcut_key`, `components/tooltip`, `utils/i18n`, `utils/user_agent`. Add `@mattermost/shared` as a devDependency for types on first use.
- Host APIs for plugins live on `window.WebappUtils` (e.g. `modals.openModalById` for published host modals) and `window.PostUtils`. Don't use `window.Components` or `window.ProductApi`; they change without notice.
- Don't add react-bootstrap, MUI, Emotion or Sass.

## Compass

- Check what `@mattermost/compass-ui` exports before building UI. Use the Compass component when one exists; don't recreate Compass building blocks or invent props.
- Import from subpaths (`@mattermost/compass-ui/components/button`), never the root barrel or a relative path into `node_modules`.
- `src/index.tsx` imports `@mattermost/compass-ui/styles` once. Never import `@mattermost/compass-ui/styles/standalone`: it overrides the user's theme.
- Never use `@mattermost/compass-proto` (unpublished prototypes).
- Icons come from `@mattermost/compass-icons/components`. In a Compass slot (`leadingIcon`, `trailingIcon`, `icon`), wrap them in Compass `Icon` with no `size`: `leadingIcon={<Icon glyph={<RefreshIcon/>}/>}`.
- Tooltips: `WithTooltip` from `@mattermost/shared/components/tooltip`. Emoji: `Emoji` from `@mattermost/shared/components/emoji`. Use these instead of the compass-ui versions until the host adopts those.
- Compass modals and menus aren't adopted in the host yet. Use `window.WebappUtils.modals` for published host modals and agree on anything else with design first.
- When migrating existing UI, switch to Compass only in code you're already changing.

## Styling

- Co-locate `my_component.css` with `my_component.tsx`. Root class in PascalCase matching the component (`.ChannelStatusSidebar`), BEM children (`.ChannelStatusSidebar__heading`), modifiers as a separate class. No `!important`, inline `style`, or element selectors.
- Colors from host theme variables (`--center-channel-color`, `--button-bg`, `--error-text`; `rgba(var(--center-channel-color-rgb), 0.72)` for transparency). Spacing, radius, elevation, type and motion from Compass tokens (`--spacing-m`, `--font-size-100`, `--line-height-100`, `--font-weight-semibold`). Never hard-code those. Pills use `var(--radius-pill)`: the webapp's `--radius-full` is `50%`.

## Accessibility and i18n

- Accessible names come from visible text or a `label` before `aria-label`, and never repeat the role ("button"). Help and error text go in `aria-describedby`. Content that updates asynchronously gets `aria-live`.
- Everything works by keyboard; focus moves into modals and back to the trigger on close.
- All UI text goes through react-intl, with strings added to `i18n/en.json` only. Prefer `<FormattedMessage>`; use `useIntl` only for string props. Use rich text (`values={{b: (chunks) => <b>{chunks}</b>}}`) instead of concatenation or Markdown.

## Security

- Call the plugin server only through `src/client.ts`; `Client4.getOptions` adds the CSRF token. Don't rely on the server's deprecated `X-Requested-With` fallback, which `ExperimentalStrictCSRFEnforcement` disables.
- Render Markdown with the host's sanitizing `window.PostUtils.formatText` + `messageHtmlToComponent`.
- Only `http:`/`https:` URLs from user data go into `href` or `src`.
- No tokens or secrets in `localStorage` or Redux, and no third-party API calls with credentials from the browser; route them through the plugin server.
- Check `event.origin` in every `message` listener.
- Before releasing: `npm audit --omit=dev` (0 advisories is the bar for the shipped bundle).

## Tests

- React Testing Library with a real Redux store and `IntlProvider`; mock `client` at the module boundary (see `channel_status_sidebar.test.tsx` and `channel_status/test_utils.tsx`).
- `await userEvent...`; `fireEvent` only for focus/blur, scroll, image load or fake timers.
- Query by role first (`getByRole` > text/placeholder > label/alt/title > `getByTestId`). No snapshots: assert `toBeVisible` and ARIA attributes.
