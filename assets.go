package main

import (
	"embed"
	"log"
)

//go:embed internal/templates/Dockerfile.tmpl
//go:embed internal/templates/orchestrator-occ.tmpl
var templateFS embed.FS

type Config struct {
	Distro        string   `yaml:"distro"`
	DistroVersion string   `yaml:"distro_version"`
	Runtime       string   `yaml:"runtime"`
	ScriptName    string   `yaml:"script_name"`
	Packages      []string `yaml:"packages"`
	NodeVersion   string   `yaml:"node_version"`
	Python        bool     `yaml:"python"`
	GoVersion     string   `yaml:"go_version"`
	Opencode      bool     `yaml:"opencode"`
	ScriptOptions []string `yaml:"script_options"`
}

func newConfig(runtime string) (Config, error) {

	defaultConfig := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Runtime:       runtime,
		ScriptName:    runtime + "-occ.sh",
		Packages:      []string{"golang", "curl", "ca-certificates", "git"},
		NodeVersion:   "22",
		Python:        false,
		GoVersion:     "1.26.6",
		Opencode:      true,
		ScriptOptions: []string{"ssh", "github", "reset", "no-cache"},
	}

	return defaultConfig, nil
}

func newConfigFile(cfg Config, force bool) error {

	// TODO: implementare la scrittura di orchestrator.yml
	// (creazione o sovrascrittura in base alla flag force)
	log.Printf("Stub: creazione di orchestrator.yml per runtime %s (force=%v)", cfg.Runtime, force)

	return nil
}