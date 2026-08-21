package main

import (
	"os"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestNewConfig(t *testing.T) {
	base := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Packages:      []string{"golang", "curl", "ca-certificates", "git", "gh"},
		NodeVersion:   "22",
		Python:        false,
		GoVersion:     "1.26.6",
		Opencode:      true,
	}

	cases := []struct {
		runtime string
		want    Config
	}{
		{"docker", Config{Runtime: "docker", ScriptName: "docker-occ.sh"}},
		{"podman", Config{Runtime: "podman", ScriptName: "podman-occ.sh"}},
	}

	for _, tc := range cases {
		got, err := newConfig(tc.runtime)
		if err != nil {
			t.Fatalf("newConfig(%q): unexpected error: %v", tc.runtime, err)
		}

		// Initialize want from base, overriding only Runtime and ScriptName
		want := base
		want.Runtime = tc.want.Runtime
		want.ScriptName = tc.want.ScriptName

		if !reflect.DeepEqual(got, want) {
			t.Errorf("newConfig(%q) = %+v, want %+v", tc.runtime, got, want)
		}
	}
}

func TestNewConfigFile(t *testing.T) {
	testDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	cfg := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Runtime:       "docker",
		ScriptName:    "docker-occ.sh",
		Packages:      []string{"golang", "curl"},
		NodeVersion:   "22",
		Python:        false,
		GoVersion:     "1.26.6",
		Opencode:      true,
		ScriptOptions: ScriptOptions{
			SSH:     true,
			GitHub:  false,
			Reset:   true,
			NoCache: false,
		},
	}

	expectedFilename := "docker-occ.yml"

	t.Run("CreateNewFile", func(t *testing.T) {
		err := newConfigFile(cfg, false)
		if err != nil {
			t.Fatalf("newConfigFile() unexpected error: %v", err)
		}

		data, err := os.ReadFile(expectedFilename)
		if err != nil {
			t.Fatalf("Failed to read created file: %v", err)
		}

		var readCfg Config
		err = yaml.Unmarshal(data, &readCfg)
		if err != nil {
			t.Fatalf("Failed to unmarshal YAML: %v", err)
		}

		if !reflect.DeepEqual(readCfg, cfg) {
			t.Errorf("Config mismatch: got %+v, want %+v", readCfg, cfg)
		}
	})

	t.Run("FileExistsNoForce", func(t *testing.T) {
		_ = newConfigFile(cfg, false)

		err := newConfigFile(cfg, false)
		if err == nil {
			t.Fatal("Expected error when file exists and force=false, got nil")
		}
	})

	t.Run("FileExistsWithForce", func(t *testing.T) {
		modifiedCfg := cfg
		modifiedCfg.Distro = "ubuntu"
		modifiedCfg.ScriptName = "ubuntu-occ.sh"

		err := newConfigFile(modifiedCfg, true)
		if err != nil {
			t.Fatalf("newConfigFile() with force=true unexpected error: %v", err)
		}

		data, err := os.ReadFile(expectedFilename)
		if err != nil {
			t.Fatalf("Failed to read modified file: %v", err)
		}

		var readCfg Config
		err = yaml.Unmarshal(data, &readCfg)
		if err != nil {
			t.Fatalf("Failed to unmarshal YAML: %v", err)
		}

		if readCfg.Distro != "ubuntu" {
			t.Errorf("Expected distro to be 'ubuntu', got '%s'", readCfg.Distro)
		}
		if readCfg.ScriptName != "ubuntu-occ.sh" {
			t.Errorf("Expected script_name to be 'ubuntu-occ.sh', got '%s'", readCfg.ScriptName)
		}
	})
}