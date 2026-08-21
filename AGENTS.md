# AGENTS.md

## What this repo is

Two things live in one repo (`git@github.com:tetrahedronix/occ.git`, branch `main`):

1. **`occ` (planned)** — a Go CLI that renders `orchestrator-occ.yml` into a `Dockerfile` + shell launcher via embedded `text/template` templates. **No Go code exists yet.** The authoritative design is `TODO.org` (roadmap, tech stack, template conventions). Read it before writing code.
2. **This container workspace** — the `Dockerfile` + `podman-occ.sh` that build/run the `ai-occ` image, which runs OpenCode over the repo. These are NOT artifacts of the `occ` tool; do not confuse them with what `occ generate` will produce.

## Workflow

- Image name is `ai-occ`. Rebuild with `podman build -t ai-occ .`, then run `./podman-occ.sh [cmd]` (defaults to `opencode`).
- The launcher requires `podman` and an active SSH agent (`SSH_AUTH_SOCK`, `ssh-add`), mounts `$HOME/.gitconfig` read-only, and uses podman-only flags (`--userns=keep-id`, `--network=slirp4netns`, `-Z` SELinux labels).
- `Dockerfile` and `podman-occ.sh` have **uncommitted local changes** (removed `ENTRYPOINT`, tmpfs `HOME=/home/opencode`, gitconfig at `/tmp/.gitconfig`). Don't revert them; the committed version is stale.

## `.gitignore` is whitelist-based

It ignores `*` and unignores only listed extensions. Consequences for new work:

- `*.yml`, `*.yaml`, and `*.tmpl` **are** unignored, so `orchestrator-occ.yml` and `Dockerfile.tmpl`/`orchestrator-occ.tmpl` are tracked.
- `*.go`, `go.mod`, `go.sum`, `Dockerfile`, `*.sh`, `*.org`, `*.md`, `*.json` are tracked.

## Conventions to honor from `TODO.org`

- YAML parser: `gopkg.in/yaml.v3`; templates: stdlib `text/template` via `embed.FS`.
- Generated files get a provenance header (occ version + hash of `orchestrator-occ.yml`).
- Shell-evaluated variables (`$SSH_AUTH_SOCK`, `$(pwd)`, `$HOME`) must survive template rendering untouched — only generation-time values go through `text/template`.
- `occ generate` writes to `./Dockerfile` by default — an overwrite that collides with this repo's existing `Dockerfile`; honor the planned overwrite policy (refuse unless `--force`).
