package cmd

import "github.com/spf13/cobra"

const (
	FlagQueue        = "queue"
	FlagCount        = "count"
	FlagScript       = "script"
	FlagDecodeBase64 = "decode-base64"
	FlagDecodeJson   = "decode-json"
	FlagWhatIf       = "what-if"
)

func getInt32FlagValue(c *cobra.Command, flagName string) *int32 {
	val, err := c.Flags().GetInt32(flagName)
	if nil != err {
		return nil
	}

	return &val
}

func getBoolFlagValue(c *cobra.Command, flagName string) *bool {
	val, err := c.Flags().GetBool(flagName)
	if nil != err {
		return nil
	}

	return &val
}
