package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
)

var peekCmd = &cobra.Command{
	Use:   "peek",
	Short: "Peek messages on a queue",
	Long:  "Peek the first few messages from a queue",
	Args:  cobra.NoArgs,

	RunE: runPeekCmd,
}

func runPeekCmd(cmd *cobra.Command, args []string) error {
	client, err := getQueueClientForCommand(cmd)
	if nil != err {
		return err
	}

	count := getInt32FlagValue(cmd, FlagCount)
	if nil == count {
		count = to.Ptr(int32(32))
	} else if *count <= 0 {
		*count = 1
	} else if *count > 32 {
		*count = 32
	}

	opts := &azqueue.PeekMessagesOptions{
		NumberOfMessages: count,
	}
	resp, err := client.PeekMessages(cmd.Context(), opts)
	if nil != err {
		return err
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Peeking %d message(s).\n", len(resp.Messages))
	for _, msg := range resp.Messages {
		json, err := json.Marshal(msg)
		if nil != err {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), string(json))
		fmt.Fprintln(cmd.OutOrStdout())
	}

	return nil
}

func init() {
	rootCmd.AddCommand(peekCmd)

	addQueueConnectionFlags(peekCmd)

	peekCmd.Flags().Int32P(FlagCount, "c", 32, "number of messages to peek")
}
