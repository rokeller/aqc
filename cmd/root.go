package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var rootCmd = &cobra.Command{
	Use:           "aqc",
	Short:         "aqc manages queues in Azure storage accounts",
	Long:          "aqc (Azure Queue Commands) helps managing messages in queues of Azure storage accounts",
	Args:          cobra.NoArgs,
	SilenceErrors: true,

	Version: version,
}

func Execute(errHandler func(error)) {
	if err := rootCmd.Execute(); nil != err {
		fmt.Fprintln(rootCmd.OutOrStdout(), err)
		errHandler(err)
	}
}
