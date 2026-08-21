---
name: Feature Request - Local LLM Configuration for OpenCode
about: Suggest support for automatic OpenCode local LLM configuration in generated sandboxes
title: 'feat: Add support for OpenCode local LLM configuration via templates'
labels: enhancement
assignees: ''
---

## Problem Statement

Currently, when using **OpenCode** with a local language model (LLM), users must manually configure the `opencode.json` file with the schema, provider, model, endpoint address, and other parameters, then inject it into the sandbox container. This manual process is error-prone and not reproducible across different environments or team members.

The `occ` tool generates reproducible container sandboxes from a declarative `orchestrator.yml`, but it currently lacks built-in support for automatically configuring OpenCode's local LLM settings.

## Current Workflow (Manual)

1. Generate sandbox with `occ generate`
2. Build the Docker image
3. Manually create/configure `opencode.json` with local LLM settings:
   ```json
   {
     "schema": "...",
     "provider": "...",
     "model": "...",
     "endpoint": "...",
     "other_params": {}
   }
   ```
4. Mount or copy `opencode.json` into the running container
5. Ensure the configuration persists across rebuilds

This approach has several issues:
- **Not reproducible**: Each developer must manually configure their environment
- **Error-prone**: Manual injection can lead to misconfigurations
- **Not version-controlled**: Configuration may drift between environments
- **Violates infrastructure-as-code principles**: Configuration exists outside the declarative spec

## Proposed Solution

Extend `occ` to support automatic OpenCode local LLM configuration through template-based generation. This would allow users to declare their LLM provider settings directly in `orchestrator.yml`, which `occ` would then use to generate both the `opencode.json` configuration and ensure it's properly mounted in the sandbox.

### Implementation Options

#### Option A: Extend `Dockerfile.tmpl` (Recommended)
Add logic to the `Dockerfile.tmpl` to generate and embed `opencode.json` during the image build process. This ensures the configuration is baked into the image itself.

**Template changes:**
```dockerfile
{{- if .OpenCode}}
RUN mkdir -p /etc/opencode && \
    cat > /etc/opencode/opencode.json <<EOF
{
  "schema": "{{.OpenCode.Schema}}",
  "provider": "{{.OpenCode.Provider}}",
  "model": "{{.OpenCode.Model}}",
  "endpoint": "{{.OpenCode.Endpoint}}"{{if .OpenCode.ExtraParams}},
  {{range $key, $value := .OpenCode.ExtraParams}}  "{{$key}}": "{{$value}}"{{end}}{{end}}
}
EOF
ENV OPENCODE_CONFIG=/etc/opencode/opencode.json
{{- end}}
```

**Schema changes (`orchestrator.yml`):**
```yaml
distro: debian:12
runtime: docker
packages:
  - curl
  - git
opencode:
  enabled: true
  schema: v1
  provider: ollama
  model: llama3.2
  endpoint: http://host.docker.internal:11434
  extra_params:
    timeout: 30
    temperature: 0.7
```

#### Option B: Extend `orchestrator-occ.tmpl` (Launcher Script)
Modify the shell launcher to mount a host-side `opencode.json` or generate it at runtime before starting the container.

**Template changes:**
```bash
{{- if .OpenCode}}
# Generate opencode.json at runtime
OPENCODE_CONFIG_DIR="$HOME/.occ/opencode"
mkdir -p "$OPENCODE_CONFIG_DIR"
cat > "$OPENCODE_CONFIG_DIR/opencode.json" <<EOF
{
  "schema": "{{.OpenCode.Schema}}",
  "provider": "{{.OpenCode.Provider}}",
  ...
}
EOF

docker run ... \
  -v "$OPENCODE_CONFIG_DIR/opencode.json:/etc/opencode/opencode.json:ro" \
  ...
{{- end}}
```

#### Option C: Generate Separate `opencode.json` Artifact
Similar to how `occ` generates `Dockerfile` and `sandbox-occ.sh`, add a third output file `opencode.json` that can be optionally mounted.

**Pros:**
- Clean separation of concerns
- Users can version-control the generated config independently
- No template complexity in Dockerfile or launcher

**Cons:**
- Requires additional file management
- Another artifact to track and potentially overwrite

## Recommended Approach

**Option A (Dockerfile.tmpl extension)** is the most idiomatic solution because:

1. **Reproducibility**: The configuration becomes part of the immutable image
2. **Simplicity**: No runtime generation or volume mount complexity
3. **Consistency**: Follows the existing pattern of embedding configuration in generated artifacts
4. **Security**: Configuration is read-only within the container
5. **Team-friendly**: All developers get identical configurations from the same `orchestrator.yml`

However, **Option B** may be preferable if:
- The LLM endpoint varies per developer (e.g., different local ports)
- Users want to override settings without rebuilding the image
- The configuration contains secrets that shouldn't be baked into images

A hybrid approach could support both: default to Option A, but allow runtime overrides via Option B when specific flags are set.

## Required Changes

### 1. Schema Extension (`doc/spec.md`)
Add new optional `opencode` section to `orchestrator.yml` schema:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `opencode` | object | no | OpenCode local LLM configuration |
| `opencode.enabled` | boolean | no | Enable OpenCode config generation (default: false) |
| `opencode.schema` | string | yes* | API schema version |
| `opencode.provider` | string | yes* | LLM provider (ollama, lmstudio, vllm, etc.) |
| `opencode.model` | string | yes* | Model name/identifier |
| `opencode.endpoint` | string | yes* | LLM API endpoint URL |
| `opencode.extra_params` | map | no | Additional provider-specific parameters |

*Required only if `opencode.enabled: true`

### 2. Template Updates
- Modify `internal/templates/Dockerfile.tmpl` to conditionally generate `opencode.json`
- Optionally modify `internal/templates/orchestrator-occ.tmpl` for runtime mounting support

### 3. Go Struct Changes
Update the configuration struct in `main.go` or relevant data structure file to include the `OpenCode` configuration object.

### 4. Validation Logic
Add validation to ensure:
- If `opencode.enabled` is true, all required fields are present
- Endpoint URLs are valid HTTP/HTTPS URLs
- Provider values are from a supported list (or allow arbitrary strings with a warning)

### 5. Documentation
Update:
- `README.md` - Mention OpenCode support in the overview
- `doc/spec.md` - Full schema documentation
- `doc/manual.md` - Usage examples and best practices

## Acceptance Criteria

- [ ] Users can declare OpenCode LLM settings in `orchestrator.yml`
- [ ] `occ generate` produces a valid `opencode.json` (either embedded in Dockerfile or as separate artifact)
- [ ] Generated configuration is correctly mounted/available in the running container
- [ ] Invalid configurations fail validation with clear error messages
- [ ] Existing sandboxes without OpenCode config continue to work unchanged
- [ ] Documentation includes example `orchestrator.yml` with OpenCode configuration
- [ ] Unit tests cover OpenCode config generation and validation

## Alternative Considered

**Environment Variables Instead of JSON Config**: Some LLM providers support configuration via environment variables (e.g., `OLLAMA_HOST`, `OPENAI_API_KEY`). This could be an alternative or complementary approach.

**Pros:**
- Simpler implementation (just add ENV directives to Dockerfile)
- More flexible for secret management via Docker secrets

**Cons:**
- Less standardized (each provider uses different env vars)
- Harder to validate and document
- May not support all OpenCode features

## Additional Context

- OpenCode documentation on local LLM configuration: [link to docs if available]
- Related issue: None yet
- This feature aligns with occ's goal of reproducible, declarative sandbox configuration
