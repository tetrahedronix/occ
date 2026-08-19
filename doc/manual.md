# occ - User Manual

`occ` builds a reproducible container sandbox for AI coding agents from a
`sandbox.yml` and renders a `Dockerfile` plus a standalone shell launcher.

## Commands

### occ init

Scaffold a starter `sandbox.yml` in the current directory.

```
occ init [--config <path>] [--force]
```

- `--config` writes to `<path>` instead of `sandbox.yml`.
- `--force` overwrites an existing file; otherwise `init` refuses.

### occ generate

Render the `Dockerfile` and the launcher from `sandbox.yml`.

```
occ generate [--config <path>] [--output-dir <dir>] [--template-dir <dir>]
             [--force] [--diff]
```

- `--config` reads `<path>` instead of `sandbox.yml`.
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
occ init          # scaffold sandbox.yml
occ generate      # write Dockerfile + sandbox-occ.sh
docker build -t occ-sandbox .    # or: podman build -t occ-sandbox .
./sandbox-occ.sh  # run the sandbox
```

The launcher forwards its first argument as the command to run inside the
container (default `opencode`).

## Template overrides

Place a custom `Dockerfile.tmpl` and/or `sandbox-occ.tmpl` in
`.occ/templates/`. Template values available:

- `{{.OccVersion}}`, `{{.ConfigFile}}`, `{{.ConfigHash}}`
- `{{.BaseImage}}`, `{{.PM}}` (`apt` or `apk`), `{{range .Packages}}`
- `{{.Runtime}}` (`podman` or `docker`)
- `{{.Userns}}`, `{{.Network}}` (Podman-only)
- `{{range .Volumes}}`, `{{.Image}}`, `{{.DefaultCmd}}`

Only generation-time values go through `text/template`; shell variables
(`$SSH_AUTH_SOCK`, `$(pwd)`, `$HOME`, ...) are emitted literally.
