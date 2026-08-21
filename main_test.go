package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "init"}
	cmd.Flags().BoolP("force", "f", false, "")
	return cmd
}

func runInitInTmpDir(t *testing.T, cmd *cobra.Command, args []string) error {
	t.Helper()
	origDir, err := os.Getwd()
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.Chdir(origDir) })
	tmpDir, err := os.MkdirTemp("", "occ-test")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })
	if err := os.Chdir(tmpDir); err != nil {
		return err
	}
	return runInit(cmd, args)
}

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	orig := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(orig) })
	return &buf
}

func TestRunInitValidRuntimes(t *testing.T) {
	for _, rt := range []string{"docker", "podman"} {
		t.Run(rt, func(t *testing.T) {
			if err := runInitInTmpDir(t, newInitCommand(), []string{rt}); err != nil {
				t.Errorf("runInit(%q): unexpected error: %v", rt, err)
			}
		})
	}
}

func TestRunInitInvalidRuntime(t *testing.T) {
	buf := captureLog(t)
	err := runInitInTmpDir(t, newInitCommand(), []string{"foo"})
	if err == nil {
		t.Fatal("runInit(\"foo\"): expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("runInit(\"foo\") error = %q, want message mentioning 'non supportato'", err)
	}
	if buf.Len() != 0 {
		t.Errorf("runInit(\"foo\") logged output on failure: %q", buf.String())
	}
}

func TestRunInitForceFlagWithForce(t *testing.T) {
	buf := captureLog(t)
	cmd := newInitCommand()
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatalf("setting force flag: %v", err)
	}
	if err := runInitInTmpDir(t, cmd, []string{"podman"}); err != nil {
		t.Fatalf("runInit: %v", err)
	}
	if !strings.Contains(buf.String(), "initialized successfully") {
		t.Errorf("log = %q, want success message", buf.String())
	}
}

func TestRunInitForceFlagDefault(t *testing.T) {
	buf := captureLog(t)
	if err := runInitInTmpDir(t, newInitCommand(), []string{"docker"}); err != nil {
		t.Fatalf("runInit: %v", err)
	}
	if !strings.Contains(buf.String(), "initialized successfully") {
		t.Errorf("log = %q, want success message", buf.String())
	}
}