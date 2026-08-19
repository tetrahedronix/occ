# occ: Reproducible Container Sandboxes for AI Coding Agents

**occ** is a Go CLI that generates reproducible container sandboxes for AI coding agents based on a declarative `orchestrator.yml`. It produces a ready-to-use `Dockerfile` and a standalone shell launcher, with no runtime dependency on `occ` after generation.

Each generated file carries a provenance header with the `occ` version and a SHA-256 hash of the source `orchestrator.yml`, making stale or hand-edited output easy to detect. The project includes embedded default templates, per-file template overrides, and strict validation that fails fast on unsupported distros or runtimes instead of emitting a broken Dockerfile.

The complete development roadmap is available at the following URL:

📋 [occ - Development Roadmap](./TODO.org)

## Commands

| Command     | Description                                     | Flags                                      |
|-------------|-------------------------------------------------|--------------------------------------------|
| `occ init`  | Scaffold a starter `orchestrator.yml`                | `--config`, `--force`                      |
| `occ generate` | Render the `Dockerfile` and shell launcher   | `--config`, `--output-dir`, `--template-dir`, `--force`, `--diff` |
| `occ version` | Print the `occ` version                       |                                            |

The `generate` command refuses to overwrite an existing `Dockerfile` or launcher unless `--force` is given; `--diff` previews the would-be changes without writing anything.

## Generated Artifacts

| Artifact        | Description                                     | URL                                  |
|-----------------|-------------------------------------------------|--------------------------------------|
| `Dockerfile`    | Base image, distro-specific package manager invocation, and provenance header | [spec.md](./doc/spec.md) |
| `sandbox-occ.sh` | Launcher with volume mounts, SSH agent forwarding, and runtime-specific flags | [spec.md](./doc/spec.md) |

Podman-only flags (`--userns`, `--network`) are emitted only when `runtime: podman`, and shell-evaluated variables (`$SSH_AUTH_SOCK`, `$(pwd)`, `$HOME`) survive generation untouched.

## Documentation

The repository includes supporting documents covering the design philosophy, schema, and usage of occ:

| Type   | Doc            | Description                          | URL                        |
|--------|----------------|--------------------------------------|----------------------------|
| Roadmap | TODO.org       | Development roadmap and architecture | [TODO.org](./TODO.org)     |
| Spec   | spec.md        | `orchestrator.yml` schema and generated artifacts | [spec.md](./doc/spec.md) |
| Manual | manual.md      | Full user manual for `occ init` and `occ generate` | [manual.md](./doc/manual.md) |

## Workflow

```
occ init          # scaffold orchestrator.yml
occ generate      # write Dockerfile + sandbox-occ.sh
docker build -t occ-sandbox .    # or: podman build -t occ-sandbox .
./sandbox-occ.sh  # run the sandbox
```

## License

This repository is distributed under the terms of the **GNU General Public License v3.0**.