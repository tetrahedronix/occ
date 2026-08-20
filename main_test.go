package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newInitCommand is a helper that instantiates and configures the Cobra data
// structure for the init command. Test functions need this structure to pass
// it to runInit().
func newInitCommand() *cobra.Command {
	// Command name accepted by the `occ` program
	cmd := &cobra.Command{Use: "init"}
	// Registers the options (flags): --force (long form), -f (short form),
	// with a default boolean value of false
	cmd.Flags().BoolP("force", "f", false, "")
	return cmd
}

// captureLog takes a pointer to the test framework and returns a pointer to a
// buffer containing the captured log text.
func captureLog(t *testing.T) *bytes.Buffer {
	// Tells Go that captureLog is a test helper: if a method on t reports a
	// problem, Go will ignore this function and report in the test logs the
	// line of the file that called captureLog
	t.Helper()

	var buf bytes.Buffer
	// Retrieves the current stream where the log package is writing its output
	orig := log.Writer()

	// Redirects the output: from this point on, every call to a logging
	// method will write into the buf buffer
	log.SetOutput(&buf)

	// Registers a cleanup function: at the end of the test (whether it
	// succeeds or fails), the log destination is restored to the original one
	// saved in orig.
	t.Cleanup(func() { log.SetOutput(orig) })

	return &buf
}

// TestRunInitValidRuntimes covers the positive cases (happy path): it passes
// perfectly valid arguments (docker and podman) to runInit() and ensures that
// it does not return any error.
func TestRunInitValidRuntimes(t *testing.T) {
	// Log capture is currently enabled: the captureLog(t) call below is
	// commented out, so runInit() writes its logs to the default output.
	captureLog(t)
	// Valid arguments expected by runInit() with the init command
	validArgs := []string{"docker", "podman"}

	// For each valid argument:
	for _, rt := range validArgs {
		// Checks whether runInit() returns an error.
		if err := runInit(newInitCommand(), []string{rt}); err != nil {
			t.Errorf("runInit(%q): unexpected error: %v", rt, err)
		}
	}
}

func TestRunInitInvalidRuntime(t *testing.T) {
	buf := captureLog(t)
	err := runInit(newInitCommand(), []string{"foo"})

	// "foo" is an invalid argument: if err is nil the test fails because
	// runInit() accepted it instead.
	if err == nil {
		t.Fatal("runInit(\"foo\"): expected error, got nil")
	}

	// From this point on, the state is an error state.

	// Assertion on the error message: it must contain the substring expected
	// for an invalid runtime.
	if !strings.Contains(err.Error(), "non supportato") {
		t.Errorf("runInit(\"foo\") error = %q, want message mentioning 'non supportato'", err)
	}

	// runInit() only logs on success.
	// On error, runInit() must not produce any log output: if the captured
	// buffer is not empty, it means the function logged something before
	// returning the error, and the test fails.
	if buf.Len() != 0 {
		t.Errorf("runInit(\"foo\") logged output on failure: %q", buf.String())
	}
}

// TestRunInitForceFlagPropagated verifies that the value of the --force flag
// set to true is propagated down to runInit(): it calls runInit() with
// force=true and the podman argument and ensures that the logs contain
// force=true.
func TestRunInitForceFlagPropagated(t *testing.T) {
	buf := captureLog(t)
	cmd := newInitCommand()

	// Sets force=true; if it cannot do so, it exits with an error.
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatalf("setting force flag: %v", err)
	}

	// With force=true, checks the behaviour of runInit() with the podman
	// argument.
	if err := runInit(cmd, []string{"podman"}); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	// Temporary test, to be changed once newConfigFile is implemented.
	// Checks that runInit() behaved correctly after the call to
	// newConfigFile.
	if !strings.Contains(buf.String(), "force=true") {
		t.Errorf("log = %q, want it to contain force=true", buf.String())
	}
}

// TestRunInitForceFlagDefault verifies the default value of the --force flag:
// it calls runInit() without setting the flag and ensures that the logs
// contain force=false, i.e. that the flag keeps its default value.
func TestRunInitForceFlagDefault(t *testing.T) {
	buf := captureLog(t)

	if err := runInit(newInitCommand(), []string{"docker"}); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	if !strings.Contains(buf.String(), "force=false") {
		t.Errorf("log = %q, want it to contain force=false", buf.String())
	}
}
