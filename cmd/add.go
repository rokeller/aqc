package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds messages to a queue",
	Args:  cobra.MinimumNArgs(1),

	RunE: runAddCmd,
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	queueName := cmd.Flag(FlagQueue).Value.String()
	client, err := getQueueClient(queueName)
	if nil != err {
		return nil
	}

	for _, msg := range args {
		_, err := client.EnqueueMessage(cmd.Context(), msg, nil)
		if nil != err {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to add message %q: %v\n", msg, err)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP(FlagQueue, "q", "", "name of the queue")
	addCmd.MarkFlagRequired(FlagQueue)

	// TODO: add flag for base64 encoding, flag for visibility timeout, flag for TTL, flag for number-of-times
}
