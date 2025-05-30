package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all messages from a queue",
	Args:  cobra.NoArgs,

	RunE: runClearCmd,
}

func runClearCmd(cmd *cobra.Command, args []string) error {
	queueName := cmd.Flag(FlagQueue).Value.String()
	client, err := getQueueClient(queueName)
	if nil != err {
		return nil
	}

	_, err = client.ClearMessages(cmd.Context(), nil)
	if nil != err {
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Deleted all messages.")

	return nil
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.Flags().StringP(FlagQueue, "q", "", "name of the queue")
	clearCmd.MarkFlagRequired(FlagQueue)
}
