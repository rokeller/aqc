package cmd

import (
	"testing"
)

func Test_runMoveCmd(t *testing.T) {
	createTestQueue(t, "test-move-src")
	addTestQueueMesssages(t, "test-move-src", []string{"a", "b", "c"})

	createTestQueue(t, "test-move-dst")

	tc := []testCase{
		{
			name: "Move messages",
			args: []string{"move",
				"--src-use-storage-emulator", "--src-queue=test-move-src",
				"--dst-use-storage-emulator", "--dst-queue=test-move-dst"},
			stdOut: "Moved 3 message(s) of 3.",
			stdErr: "",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-move-src")
				if c != 0 {
					t.Errorf("got %d messages in source queue, want 0", c)
				}

				c = fetchTestQueueMessageCount(t, "test-move-dst")
				if c != 3 {
					t.Errorf("got %d messages in destination queue, want 3", c)
				}
			},
		},
	}

	executeTestCases(t, tc)
}
