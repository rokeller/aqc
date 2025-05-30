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
	client, err := getQueueClientForCommand(cmd)
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

	addQueueConnectionFlags(addCmd)

	// TODO: add flag for base64 encoding, flag for visibility timeout, flag for TTL, flag for number-of-times
}
