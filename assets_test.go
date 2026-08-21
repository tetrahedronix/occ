package main

import (
	"os"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

// TestNewConfig verifies that newConfig() returns the expected default Config
// for each supported runtime. It builds the base default configuration, then
// for every runtime in the cases table it checks that the returned Config
// matches the base one with only Runtime and ScriptName overridden.
func TestNewConfig(t *testing.T) {
	base := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Packages:      []string{"golang", "curl", "ca-certificates", "git", "gh"},
		NodeVersion:   "22",
		GoVersion:     "1.26.6",
		Opencode:      true,
		ScriptOptions: []string{"ssh", "github", "reset", "no-cache"},
	}

	cases := []struct {
		runtime string
		want    Config
	}{
		{"docker", Config{Runtime: "docker", ScriptName: "docker-occ.sh"}},
		{"podman", Config{Runtime: "podman", ScriptName: "podman-occ.sh"}},
	}

	// For each test case, calls newConfig() with the runtime under test and
	// compares the result against the expected base configuration.
	for _, tc := range cases {
		// Calls newConfig() with the current runtime.
		got, err := newConfig(tc.runtime)
		// If newConfig() returns an error, the test fails immediately.
		if err != nil {
			t.Fatalf("newConfig(%q): unexpected error: %v", tc.runtime, err)
		}

		// Overrides Runtime and ScriptName in the base configuration with the
		// expected values for this runtime, so it can be compared with got.
		base.Runtime = tc.runtime
		base.ScriptName = tc.want.ScriptName

		// Compares the returned Config with the expected one field by field.
		if !reflect.DeepEqual(got, base) {
			t.Errorf("newConfig(%q) = %+v, want %+v", tc.runtime, got, base)
		}
	}
}

// TestNewConfigFile tests the newConfigFile function which creates a YAML
// configuration file from a Config struct. It verifies:
// - File creation when force is true
// - File creation when force is false and file doesn't exist
// - Error when force is false and file already exists
// - Correct YAML content in the created file
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
		ScriptOptions: []string{"ssh", "github"},
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
	})
}