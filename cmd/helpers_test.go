package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name   string
	args   []string
	err    error
	stdOut string
	stdErr string
	verify func(t *testing.T)
}

func getStorageEmulatorQueueClient(t *testing.T, queueName string) *azqueue.QueueClient {
	svcClient, err := azqueue.NewServiceClientFromConnectionString(StorageEmulatorConnectionString, nil)
	if nil != err {
		t.Fatalf("failed to create service client: %v", err)
	}

	return svcClient.NewQueueClient(queueName)
}

func createTestQueue(t *testing.T, name string) {
	client := getStorageEmulatorQueueClient(t, name)

	client.Delete(t.Context(), nil)
	_, err := client.Create(t.Context(), nil)
	if nil != err {
		t.Fatalf("failed to create test queue %q: %v", name, err)
	}
}

func addTestQueueMesssages(t *testing.T, queueName string, msgs []string) {
	client := getStorageEmulatorQueueClient(t, queueName)

	for _, msg := range msgs {
		_, err := client.EnqueueMessage(t.Context(), msg, nil)
		if nil != err {
			t.Fatalf("failed to enqueue message %q to queue %q: %v", msg, queueName, err)
		}
	}
}

func fetchTestQueueMesssages(t *testing.T, queueName string) []*azqueue.DequeuedMessage {
	client := getStorageEmulatorQueueClient(t, queueName)

	resp, err := client.DequeueMessages(t.Context(), &azqueue.DequeueMessagesOptions{
		NumberOfMessages: to.Ptr(int32(32)),
	})
	if nil != err {
		t.Fatalf("failed to fetch messages from test queue %q: %v", queueName, err)
	}

	return resp.Messages
}

func executeTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, testCase := range testCases {
		t.Run(testCase.name, testCase.executeTestCase)
	}
}

func (testCase testCase) executeTestCase(t *testing.T) {
	t.Helper()

	stdOut, stdErr, err := execute(t, testCase.args...)

	assert.Equal(t, testCase.err, err, "expected error must match")

	if testCase.err == nil {
		assert.Equal(t, testCase.stdOut, stdOut, "stdout must match")
		assert.Equal(t, testCase.stdErr, stdErr, "stderr must match")
	}

	if testCase.verify != nil {
		testCase.verify(t)
	}
}

func execute(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	c := rootCmd

	// reset all flags to start clean
	cmdFn := func(c *cobra.Command) {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
			f.Changed = false
		})
	}
	visitCommands(c.Commands(), cmdFn)

	bufOut, bufErr := captureStdOutAndErr(c)
	c.SetArgs(args)
	err := c.Execute()

	return strings.TrimSpace(bufOut.String()),
		strings.TrimSpace(bufErr.String()),
		err
}

func captureStdOutAndErr(c *cobra.Command) (bufOut, bufErr *bytes.Buffer) {
	bufOut = new(bytes.Buffer)
	bufErr = new(bytes.Buffer)

	c.SetOut(bufOut)
	c.SetErr(bufErr)

	return
}

func visitCommands(cs []*cobra.Command, cmdFn func(c *cobra.Command)) {
	for _, c := range cs {
		cmdFn(c)
		visitCommands(c.Commands(), cmdFn)
	}
}
