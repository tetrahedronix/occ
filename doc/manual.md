# occ - User Manual

`occ` builds a reproducible container sandbox for AI coding agents from a
YAML file and renders a `Dockerfile` plus a standalone shell launcher.

## Commands

### occ init

Scaffold a starter `orchestrator-occ.yml` in the current directory.

```
occ init [--config <path>] [--force]
```

- `--config` writes to `<path>` instead of `orchestrator-occ.yml`.
- `--force` overwrites an existing file; otherwise `init` refuses.

### occ generate

Render the `Dockerfile` and the launcher from `orchestrator-occ.yml`.

```
occ generate [--config <path>] [--output-dir <dir>] [--template-dir <dir>]
             [--force] [--diff]
```

- `--config` reads `<path>` instead of `orchestrator-occ.yml`.
- `--output-dir` writes artifacts to `<dir>` (created if missing) instead of
  the current directory.
- `--template-dir` reads custom templates from `<dir>` (default
  `.occ/templates/`), falling back to embedded defaults per file.
- `--force` overwrites existing artifacts.
- `--diff` prints the diff of would-be changes and writes nothing.

### occ version

Print the `occ` version.

## Workflow

```
occ init          # scaffold orchestrator-occ.yml
occ generate      # write Dockerfile + orchestrator-occ.sh
docker build -t occ-sandbox .    # or: podman build -t occ-sandbox .
./orchestrator-occ.sh  # run the sandbox
```

The launcher forwards its first argument as the command to run inside the
container (default `opencode`).

## Template overrides

Place a custom `Dockerfile.tmpl` and/or `orchestrator-occ.tmpl` in
`.occ/templates/`. Template values available:

- `{{.OccVersion}}`, `{{.ConfigFile}}`, `{{.ConfigHash}}`
- `{{.BaseImage}}`, `{{.PM}}` (`apt` or `apk`), `{{range .Packages}}`
- `{{.Runtime}}` (`podman` or `docker`)
- `{{.PodmanOptions.Userns}}`, `{{.PodmanOptions.Network}}` (emitted only when `runtime: podman`)
- `{{range .Volumes}}`, `{{.Image}}`, `{{.DefaultCmd}}`

Only generation-time values go through `text/template`; shell variables
(`$SSH_AUTH_SOCK`, `$(pwd)`, `$HOME`, ...) are emitted literally.
