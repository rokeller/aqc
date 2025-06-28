package cmd

import (
	"errors"
	"testing"
)

func Test_runAddCmd(t *testing.T) {
	createTestQueue(t, "test-add")

	tc := []testCase{
		{
			name: "Unsupported flag",
			args: []string{"add", "--foo"},
			err:  errors.New("unknown flag: --foo"),
		},
		{
			name:   "Add single message",
			args:   []string{"add", "--use-storage-emulator", "-q=test-add", "message-one"},
			stdOut: "Added 1 message(s).",
			stdErr: "",
		},
		{
			name:   "Add multiple messages",
			args:   []string{"add", "--use-storage-emulator", "-q=test-add", "message-a", "message-b", "message-c"},
			stdOut: "Added 3 message(s).",
			stdErr: "",
		},
	}

	executeTestCases(t, tc)
}
