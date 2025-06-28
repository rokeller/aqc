package cmd

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move messages from one queue to another",
	Args:  cobra.NoArgs,

	RunE: runMoveCmd,
}

func runMoveCmd(cmd *cobra.Command, args []string) error {
	start := time.Now().UTC()

	sourceClient, err := getQueueClientForCommandWithPrefix(cmd, PrefixSource)
	if nil != err {
		return nil
	}
	destinationClient, err := getQueueClientForCommandWithPrefix(cmd, PrefixDestination)
	if nil != err {
		return nil
	}

	numMoved := 0
	numSeen := 0

	for {
		opts := &azqueue.DequeueMessagesOptions{
			NumberOfMessages:  to.Ptr(int32(32)),
			VisibilityTimeout: to.Ptr(int32(30)),
		}
		var msgs []*azqueue.DequeuedMessage
		if resp, err := sourceClient.DequeueMessages(cmd.Context(), opts); nil != err {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to dequeue messages: %v\n", err)
			return err
		} else {
			msgs = resp.Messages
		}

		if len(msgs) <= 0 {
			break
		}

		for _, msg := range msgs {
			numSeen += 1
			if msg.InsertionTime.UTC().After(start) {
				continue
			}

			if _, err := destinationClient.EnqueueMessage(cmd.Context(), *msg.MessageText, nil); nil != err {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to enqueue target message: %v", err)
				continue
			}

			if _, err := sourceClient.DeleteMessage(cmd.Context(), *msg.MessageID, *msg.PopReceipt, nil); nil != err {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to delete source message: %v", err)
			}
			numMoved += 1
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Moved %d message(s) of %d.\n", numMoved, numSeen)

	return nil
}

func init() {
	rootCmd.AddCommand(moveCmd)

	addQueueConnectionFlagsWithPrefix(moveCmd, PrefixSource, "source queue")
	addQueueConnectionFlagsWithPrefix(moveCmd, PrefixDestination, "destination queue")
}
