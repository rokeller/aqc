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
	client, err := getQueueClientForCommand(cmd)
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

	addQueueConnectionFlags(clearCmd)
}
