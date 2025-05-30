package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var rootCmd = &cobra.Command{
	Use:           "aqc",
	Short:         "aqc manages queues in Azure storage accounts",
	Long:          "aqc (Azure Queue Commands) helps managing queues and messages in Azure storage accounts",
	Args:          cobra.NoArgs,
	SilenceErrors: true,

	Version: version,
}

func Execute(errHandler func(error)) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.OutOrStdout(), err)
		errHandler(err)
	}
}
