package main

import (
	"reflect"
	"testing"
)

// TestNewConfig verifies that newConfig() returns the expected default Config
// for each supported runtime. It builds the base default configuration, then
// for every runtime in the cases table it checks that the returned Config
// matches the base one with only Runtime and ScriptName overridden.
func TestNewConfig(t *testing.T) {
	base := Config{
		Distro:        "debian",
		DistroVersion: "13.6",
		Packages:      []string{"golang", "curl", "ca-certificates", "git"},
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