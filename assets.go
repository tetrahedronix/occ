package main

import (
	"embed"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

//go:embed internal/templates/Dockerfile.tmpl
//go:embed internal/templates/orchestrator-occ.tmpl
var templateFS embed.FS

type ScriptOptions struct {
	SSH     bool `yaml:"ssh"`
	GitHub  bool `yaml:"github"`
	Reset   bool `yaml:"reset"`
	NoCache bool `yaml:"no-cache"`
}

type PodmanOptions struct {
	Userns  string `yaml:"userns,omitempty"`
	Network string `yaml:"network,omitempty"`
}

type Config struct {
	Distro        Distro         `yaml:"distro"`
	DistroVersion string         `yaml:"distro_version"`
	Runtime       Runtime        `yaml:"runtime"`
	ScriptName    string         `yaml:"script_name"`
	Packages      []string       `yaml:"packages"`
	NodeVersion   string         `yaml:"node_version"`
	Python        bool           `yaml:"python"`
	GoVersion     string         `yaml:"go_version"`
	Opencode      bool           `yaml:"opencode"`
	ScriptOptions ScriptOptions  `yaml:"script_options"`
	PodmanOptions *PodmanOptions `yaml:"podman_options,omitempty"`
}

func newConfig(runtime Runtime) (Config, error) {

	defaultConfig := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Runtime:       runtime,
		ScriptName:    string(runtime) + "-occ.sh",
		Packages:      []string{"golang", "curl", "ca-certificates", "git", "gh"},
		NodeVersion:   "22",
		Python:        false,
		GoVersion:     "1.26.6",
		Opencode:      true,
		ScriptOptions: ScriptOptions{
			SSH:     false,
			GitHub:  false,
			Reset:   false,
			NoCache: false,
		},
	}

	if runtime == RuntimePodman {
		defaultConfig.PodmanOptions = &PodmanOptions{
			Userns:  "keep-id",
			Network: "slirp4netns",
		}
	}

	return defaultConfig, nil
}

// Phase 3: "Implement YAML parser module"

func (c Config) Validate() error {

	if _, err := c.Distro.PackageManager(); err != nil {
		return err
	}

	if err := c.Runtime.Validate(); err != nil {
		return err
	}

	// if c.Runtime != RuntimePodman && c.Runtime != RuntimeDocker {
	// 	return fmt.Errorf("unsupported runtime: %q", c.Runtime)
	// }

	return nil
}

func newConfigFile(cfg Config, force bool) error {

	data, err := yaml.Marshal(&cfg)

	if err != nil {
		return fmt.Errorf("error during YAML serialization: %w", err)
	}

	filename := string(cfg.Runtime + "-occ.yml")

	// Scrivi direttamente
	if force {
		return os.WriteFile(filename, data, 0644)
	}

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("the file %s already exists. Use --force to overwrite", filename)
	}

	return os.WriteFile(filename, data, 0644)
}
