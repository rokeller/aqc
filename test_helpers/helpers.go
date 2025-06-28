package test_helpers

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ExecuteWithExit helps testing cases where os.Exit(...) is called. For more
// information see: https://go.dev/talks/2014/testing.slide#23
func ExecuteWithExit(
	t *testing.T,
	name string,
	fnWithExit func(*testing.T),
	expectedExitCode int) {
	t.Helper()

	if os.Getenv("ACTUALLY_EXECUTE") == "1" {
		fnWithExit(t)
		return
	}

	testArgs := append(os.Args[1:], "-test.run=^"+name+"$")
	cmd := exec.Command(os.Args[0], testArgs...)
	cmd.Env = append(os.Environ(), "ACTUALLY_EXECUTE=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, expectedExitCode, e.ExitCode(), "non-zero exit code matches")
		return
	}

	assert.Equal(t, expectedExitCode, 0, "exit code matches")
	assert.Nil(t, err,
		"Execute ran with error [%v], want exit status %d", err, expectedExitCode)
}
