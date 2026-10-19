# AGENTS.md

Mattermost plugin starter template: a Go server plugin (`server/`) and a React/TypeScript webapp plugin (`webapp/`). Either half can be deleted; the Makefile skips whatever `plugin.json` doesn't declare.

**Before editing `server/`, read `server/AGENTS.md`. Before editing `webapp/`, read `webapp/AGENTS.md`.** Step-by-step guides for adding a slash command, HTTP endpoint, webhook, setting or webapp UI are in `README.md` under "Modifying the server boilerplate" and "Modifying the webapp boilerplate".

## Commands

While iterating, run the narrow command; `make test` and `make check-style` reinstall Go tools first, and `make test` also rebuilds the bundle.

```bash
SKIP_DOCKER_TESTS=1 go test ./server/...       # server unit tests only
go test ./server/command -run TestStatusCommand # single test
cd webapp && npx jest src/manifest.test.tsx    # single webapp test
go vet ./... && golangci-lint run ./...        # server lint
cd webapp && npm run lint && npm run check-types
make test                                      # full check before pushing
make mock                                      # regenerate mocks after changing command.Command or kvstore.KVStore
```

`make help` lists the rest (deploy, watch, release bumps, debugger attach).

## Architecture

- `plugin.json` is the source of truth for ID, version, and settings. `make apply` generates `server/manifest.go` and `webapp/src/manifest.ts` from it; both are gitignored, so never edit them. A fresh checkout doesn't compile until `make apply` runs. Never edit `dist/`, `server/dist/` or `webapp/dist/`.
- Project-specific make targets go in `build/custom.mk`.
- `min_server_version` is 12.0.0. Don't use plugin API methods or host webapp exports newer than that without bumping it.

## Rules that apply everywhere

- The server only authenticates plugin requests. Every route and command that touches data checks the user's permission on that object, on the server; hiding UI is not a check.
- Anything that changes state uses POST, PUT, PATCH or DELETE. The server's CSRF check doesn't cover GET.
- Never send internal errors to clients or write secrets, tokens or message content to logs.
- Don't add dependencies the host, Compass or the standard library already provide. Fix CodeQL alerts instead of excluding them in `.github/codeql/codeql-config.yml`.
- In GitHub workflows, pin third-party actions to a commit SHA with the version in a comment, and start from `permissions: {}`.

## Pull requests and commits

- Commit subject: one succinct line. Don't add a `Co-Authored-By` line for the agent.
- In the PR description, call out any change to permissions, authentication, secrets or `min_server_version`.
