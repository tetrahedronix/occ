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

type Config struct {
	Distro        string        `yaml:"distro"`
	DistroVersion string        `yaml:"distro_version"`
	Runtime       string        `yaml:"runtime"`
	ScriptName    string        `yaml:"script_name"`
	Packages      []string      `yaml:"packages"`
	NodeVersion   string        `yaml:"node_version"`
	Python        bool          `yaml:"python"`
	GoVersion     string        `yaml:"go_version"`
	Opencode      bool          `yaml:"opencode"`
	ScriptOptions ScriptOptions `yaml:"script_options"`
}

func newConfig(runtime string) (Config, error) {

	defaultConfig := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Runtime:       runtime,
		ScriptName:    runtime + "-occ.sh",
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

	return defaultConfig, nil
}

func newConfigFile(cfg Config, force bool) error {

	data, err := yaml.Marshal(&cfg)

	if err != nil {
		return fmt.Errorf("error during YAML serialization: %w", err)
	}

	filename := cfg.Runtime + "-occ.yml"

	// Scrivi direttamente
	if force {
		return os.WriteFile(filename, data, 0644)
	}

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("the file %s already exists. Use --force to overwrite", filename)
	}

	return os.WriteFile(filename, data, 0644)
}
