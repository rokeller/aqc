package cmd

import (
	"testing"
)

func Test_runPeekCmd(t *testing.T) {
	createTestQueue(t, "test-peek")
	addTestQueueMesssages(t, "test-peek", []string{"1", "2", "3", "4"})

	tc := []testCase{
		{
			name:           "Count lower bound enforced",
			args:           []string{"peek", "--use-storage-emulator", "-q=test-peek", "-c0"},
			stdOutContains: "\"MessageText\":\"1\"",
			stdErr:         "Peeking 1 message(s).",
		},
		{
			name:           "Count upper bound enforced",
			args:           []string{"peek", "--use-storage-emulator", "-q=test-peek", "-c33"},
			stdOutContains: "\"MessageText\":\"4\"",
			stdErr:         "Peeking 4 message(s).",
		},
		{
			name:           "Count upper bound enforced",
			args:           []string{"peek", "--use-storage-emulator", "-q=test-peek", "-c2"},
			stdOutContains: "\"MessageText\":\"2\"",
			stdErr:         "Peeking 2 message(s).",
		},
	}

	executeTestCases(t, tc)
}
