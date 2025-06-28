package cmd

import (
	"testing"
)

func Test_runMoveCmd(t *testing.T) {
	createTestQueue(t, "test-move-src-01")
	addTestQueueMesssages(t, "test-move-src-01", []string{"a", "b", "c"})

	createTestQueue(t, "test-move-src-02")
	addTestQueueMesssages(t, "test-move-src-02", []string{"a", "b", "c"})

	createTestQueue(t, "test-move-dst")

	tc := []testCase{
		{
			name: "Move messages works",
			args: []string{"move",
				"--src-use-storage-emulator", "--src-queue=test-move-src-01",
				"--dst-use-storage-emulator", "--dst-queue=test-move-dst"},
			stdOut: "Moved 3 message(s) of 3.",
			stdErr: "",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-move-src-01")
				if c != 0 {
					t.Errorf("got %d messages in source queue, want 0", c)
				}

				c = fetchTestQueueMessageCount(t, "test-move-dst")
				if c != 3 {
					t.Errorf("got %d messages in destination queue, want 3", c)
				}
			},
		},
		{
			name: "Destination queue does not exist",
			args: []string{"move",
				"--src-use-storage-emulator", "--src-queue=test-move-src-02",
				"--dst-use-storage-emulator", "--dst-queue=test-move-dst-does-not-exist"},
			stdOut:         "Moved 0 message(s) of 3.",
			stdErrContains: "RESPONSE 404: 404 The specified queue does not exist.",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-move-src-02")
				if c != 3 {
					t.Errorf("got %d messages in source queue, want 0", c)
				}

				c = fetchTestQueueMessageCount(t, "test-move-dst")
				if c != 3 {
					t.Errorf("got %d messages in destination queue, want 3", c)
				}
			},
		},
		{
			name: "Source queue does not exist",
			args: []string{"move",
				"--src-use-storage-emulator", "--src-queue=test-move-src",
				"--dst-use-storage-emulator", "--dst-queue=test-move-dst-does-not-exist"},
			stdOutContains: "Usage:",
			stdErrContains: "Failed to dequeue messages",
			errContains:    "QueueNotFound",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-move-src-02")
				if c != 3 {
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
