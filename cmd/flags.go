package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
)

const (
	FlagQueue           = "queue"
	FlagServiceURL      = "service-url"
	FlagQueueURL        = "queue-url"
	FlagStorageEmulator = "use-storage-emulator"
	FlagCount           = "count"
	FlagScript          = "script"
	FlagDecodeBase64    = "decode-base64"
	FlagDecodeJson      = "decode-json"
	FlagWhatIf          = "what-if"

	PrefixSource      = "src-"
	PrefixDestination = "dst-"
)

func addQueueConnectionFlags(cmd *cobra.Command) {
	addQueueConnectionFlagsWithPrefix(cmd, "", "")
}

func addQueueConnectionFlagsWithPrefix(cmd *cobra.Command, prefix, queueName string) {
	if queueName == "" {
		queueName = "queue"
	}

	cmd.Flags().Bool(prefix+FlagStorageEmulator, false,
		fmt.Sprintf("connect to the storage emulator for the %s", queueName))
	if prefix == "" {
		cmd.Flags().StringP(FlagQueue, "q", "",
			fmt.Sprintf("name of the %s", queueName))
	} else {
		cmd.Flags().String(prefix+FlagQueue, "",
			fmt.Sprintf("name of the %s", queueName))
	}
	cmd.Flags().String(prefix+FlagQueueURL, "",
		fmt.Sprintf("URL of the %s", queueName))
	cmd.Flags().String(prefix+FlagServiceURL, "",
		fmt.Sprintf("URL of service of the %s", queueName))

	cmd.MarkFlagsMutuallyExclusive(prefix+FlagQueue, prefix+FlagQueueURL)
	cmd.MarkFlagsMutuallyExclusive(prefix+FlagStorageEmulator, prefix+FlagServiceURL, prefix+FlagQueueURL)
	cmd.MarkFlagsOneRequired(prefix+FlagStorageEmulator, prefix+FlagServiceURL, prefix+FlagQueueURL)
}

func getQueueClientForCommand(cmd *cobra.Command) (*azqueue.QueueClient, error) {
	return getQueueClientForCommandWithPrefix(cmd, "")
}

func getQueueClientForCommandWithPrefix(cmd *cobra.Command, prefix string) (*azqueue.QueueClient, error) {
	useEmulator := getBoolFlagValue(cmd, prefix+FlagStorageEmulator)
	queueName := cmd.Flag(prefix + FlagQueue).Value.String()
	queueURL := cmd.Flag(prefix + FlagQueueURL).Value.String()
	serviceURL := cmd.Flag(prefix + FlagServiceURL).Value.String()

	if nil != useEmulator && *useEmulator {
		return getQueueClientForStorageEmulator(queueName)
	}

	if "" != queueURL {
		if hasSASSignature(queueURL) {
			return getQueueClientForQueueURLWithSASToken(queueURL)
		}

		return getQueueClientForQueueURLWithDefaultCredential(queueURL)
	}

	if "" != serviceURL {
		if "" == queueName {
			return nil, errors.New(
				"The queue parameter must be provided when the service-url parameter is passed")
		}

		return getQueueClientForServiceURLWithDefaultCredential(serviceURL, queueName)
	}

	return nil, errors.New("not enough flags for connection to queue")
}

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

func hasSASSignature(url string) bool {
	return strings.Contains(url, "&sig=")
}
