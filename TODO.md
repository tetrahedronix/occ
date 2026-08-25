# TODO – occ (Container Sandbox Generator for AI Coding Agents)

## 1️⃣ Staged Plan (numbered)

### 1.1 Phase 1: Project Setup & Data Structures

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 1 | **Initialize Go module** | `go mod init`. | Go module initialized | ✅ |
| 2 | **Go structs for orchestrator.yml** | Define structs mirroring the schema: distro/base image, system packages list, runtime options & flags, `runtime` field (enum `podman`/`docker`), separate `podman_options` struct (`*PodmanOptions`, pointer, `omitempty`) populated only when `runtime: podman` (defaults `Userns="keep-id"`, `Network="slirp4netns"`), user-configurable output script name. | Complete, YAML-serializable `Config` struct | ✅ |
| 3 | **Distro → package-manager mapping** | Explicit list of supported distros (Debian/Ubuntu/Alpine), per-distro install syntax (`PMApt="apt"`, `PMApk="apk"`) already used in `Dockerfile.tmpl` conditional blocks, scope documented in `doc/spec.md`. Still missing: wiring `.PM` into the generator — `Config.Validate()` must be called **before** `template.Execute()` (Phase 3), otherwise `.PM=""` silently falls into the `apt` branch even for unsupported distros. | Mapping table + explicit-failure validation for unsupported distros | ⬜ |
| 4 | **yaml.v3 dependency** | Added `gopkg.in/yaml.v3`. | Dependency present in `go.mod` | ✅ |

### 1.2 Phase 2: Template Design & Embedding

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 5 | **Draft Dockerfile.tmpl** | Base OS selection, package manager invocation branched per distro (apt, apk, etc.), environment & user setup, provenance header (version + config hash). | Default `Dockerfile.tmpl` ready | ⬜ |
| 6 | **Draft orchestrator-occ.tmpl (shell launcher)** | Dynamic handling of enabled flags, volume mounts and run options, per-project isolation of the OpenCode home (`$HOME/.occ/homes/<slug>/`, slug = `basename(project_dir) + "-" + sha256(absolute_path)[:8]`, computed at generation time in Go and embedded as a literal) replacing the shared `$HOME/.opencode-home` that leaks state across sandboxes, runtime-specific flags branched on `runtime` (Podman/Docker options are not interchangeable), shell-runtime variables (`$SSH_AUTH_SOCK`, `$(pwd)`, `$HOME`) left untouched by the template, provenance header. | Default `orchestrator-occ.tmpl` ready, with isolated per-project home | ⬜ |
| 7 | **Fix Git/SSH auth in generated sandbox (Issue #5)** | Pre-seed `known_hosts` in `Dockerfile.tmpl` (`ssh-keyscan -H github.com gitlab.com bitbucket.org`), replace the read-only `.gitconfig` mount with a writable copy in the persistent home, wire `ScriptOptions.SSH` to conditional `SSH_AUTH_SOCK` forwarding, wire `ScriptOptions.GitHub` to `GH_TOKEN`-based credential helper setup, document the required PAT scope (`contents: read/write`) in `doc/manual.md`. | Working Git/SSH authentication in the generated sandbox | ⬜ |
| 8 | **Embed default templates** | Use `embed.FS` to bundle templates into the binary. | Templates embedded, no external file dependency | ⬜ |
| 9 | **Template override convention** | Search order for user-provided `Dockerfile.tmpl` / `orchestrator-occ.tmpl` (project root vs `.occ/templates/`), per-file fallback to embedded default (a custom `Dockerfile.tmpl` shouldn't force a custom `orchestrator-occ.tmpl` too). | Lookup convention documented and implemented | ⬜ |

### 1.3 Phase 3: Core Logic & Generators

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 10 | **YAML parser** | Fallback handling if `orchestrator.yml` is missing or invalid, schema validation (required fields present, `distro`/`runtime` enum values within the supported set), fail fast with a clear error message before touching templates. | Robust YAML parser module | ⬜ |
| 11 | **Dockerfile generator** | Load user-provided template or fall back to the embedded default, render with `text/template`, guard against silent `<no value>` output on missing struct fields, write output to `Dockerfile`. | Complete Dockerfile generator | ⬜ |
| 12 | **Shell launcher generator** | Render `orchestrator-occ.tmpl` with parsed options, write output to the user-chosen `script_name.sh`, executable permissions (`chmod +x` / mode `0755`). | Complete launcher generator | ⬜ |
| 13 | **Overwrite policy** | Default: refuse to overwrite an existing `Dockerfile`/launcher without confirmation; `--force` flag for unconditional overwrite; optional `--diff` to preview changes. | Overwrite policy implemented | ⬜ |

### 1.4 Phase 4: CLI Interface & UX

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 14 | **CLI entrypoint (main.go)** | Subcommand `occ init` (scaffold `orchestrator.yml`), subcommand `occ generate` (read `orchestrator.yml`, produce `Dockerfile` + launcher), flags for custom YAML path, output directory and custom templates, `--force`/`--diff` flags. | Working CLI with `init`/`generate` subcommands | ⬜ |
| 15 | **Logging & user feedback** | Clear success/error messages. | Polished CLI UX | ⬜ |

### 1.5 Phase 5: Testing & Validation

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 16 | **Unit tests for YAML parsing & template rendering** | Valid config cases per supported distro/runtime combination, invalid config cases (missing required field, unsupported distro/runtime value). | Test coverage for parsing/rendering | ⬜ |
| 17 | **Unit tests for overwrite policy** | Tests for default refusal, `--force`, `--diff`. | Test coverage for overwrite policy | ⬜ |
| 18 | **Test generated Dockerfile** | Build with `docker build` and `podman build`. | Generated Dockerfile buildable on both runtimes | ⬜ |
| 19 | **Test shell launcher execution** | Verify runtime-specific flags appear only for the declared engine (Podman flags absent when `runtime: docker` and vice versa), verify shell-runtime variables (`$SSH_AUTH_SOCK`, etc.) survive generation unresolved. | Generated launcher tested end-to-end | ⬜ |

### 1.6 Phase 6: OpenCode Local LLM Configuration

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 20 | **OpenCode local LLM schema** | Add an `opencode` section to `orchestrator.yml` with `model`, `temperature`, `system_prompt` fields; validation: supported model list, temperature range, system_prompt length; fail fast on invalid config. | Schema extension defined | ⬜ |
| 21 | **Template changes for LLM-aware generation** | Option A: extend `Dockerfile.tmpl` with an `{{.Opencode}}` block for local model setup; Option B: modify `orchestrator-occ.tmpl` to inject LLM config into the shell environment; Option C: generate a separate artifact (e.g. `opencode-config.json`) referenced by the launcher. | Decision made and implemented | ⬜ |
| 22 | **Validation & acceptance criteria** | Supported compatible local models, valid temperature range (0.0–1.0), maximum character limit for system prompt, strict-mode guard for missing `opencode` fields. | Validation criteria implemented | ⬜ |
| 23 | **Issue template & documentation** | `feature-opencode-local-llm.md` template for tracking configuration requests, `.gitignore` update, documentation of the workflow (from manual `opencode.json` injection to declarative template-based generation). | Documentation and issue template ready | ⬜ |

## 2️⃣ Key Characteristics & Required Software

| Category | Detail | Software / Library |
|----------|--------|----------------------|
| Language | Core of the project | Go |
| YAML parser | Parsing/serialization of `orchestrator.yml` | `gopkg.in/yaml.v3` |
| Templating engine | Rendering of `Dockerfile` and shell launcher | `text/template` (stdlib) |
| Embedded templates | Default templates bundled into the binary | `embed` (stdlib) |
| Container runtime | Declared per-project (`podman` or `docker`), templates branch on this field | Podman (rootless, primary) / Docker |
| Package manager mapping | Distro-specific install syntax (`apt-get`, `apk`, ...), explicit scope, not implied by base-image name | In-house mapping table |
| Provenance | Header with `occ` version + SHA-256 hash of the source `orchestrator.yml`, to detect stale or hand-edited output | `crypto/sha256` (stdlib) |

---

*The tables above highlight all the non-obvious requirements agents should be aware of.*