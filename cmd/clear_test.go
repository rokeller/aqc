package cmd

import (
	"testing"
)

func Test_runClearCmd(t *testing.T) {
	createTestQueue(t, "test-clear-empty")
	createTestQueue(t, "test-clear-some-msgs")
	addTestQueueMesssages(t, "test-clear-some-msgs", []string{"a", "b", "c"})

	tc := []testCase{
		{
			name:   "Empty queue",
			args:   []string{"clear", "--use-storage-emulator", "-q=test-clear-empty"},
			stdOut: "Deleted all messages.",
			verify: func(t *testing.T) {
				msgs := fetchTestQueueMesssages(t, "test-clear-empty")
				if len(msgs) != 0 {
					t.Errorf("The queue must not have messages, but has %d", len(msgs))
				}
			},
		},
		{
			name:   "Queue with messages",
			args:   []string{"clear", "--use-storage-emulator", "-q=test-clear-some-msgs"},
			stdOut: "Deleted all messages.",
			verify: func(t *testing.T) {
				msgs := fetchTestQueueMesssages(t, "test-clear-some-msgs")
				if len(msgs) != 0 {
					t.Errorf("The queue must not have messages, but has %d", len(msgs))
				}
			},
		},
	}

	executeTestCases(t, tc)
}
