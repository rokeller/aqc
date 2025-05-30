package cmd

import (
	"errors"
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
)

func addQueueConnectionFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagStorageEmulator, false, "connect to the storage emulator")
	cmd.Flags().StringP(FlagQueue, "q", "", "name of the queue")
	cmd.Flags().String(FlagQueueURL, "", "URL of the queue")
	cmd.Flags().String(FlagServiceURL, "", "URL of the queue service")

	cmd.MarkFlagsMutuallyExclusive(FlagQueue, FlagQueueURL)
	cmd.MarkFlagsMutuallyExclusive(FlagStorageEmulator, FlagServiceURL, FlagQueueURL)
	cmd.MarkFlagsOneRequired(FlagStorageEmulator, FlagServiceURL, FlagQueueURL)
}

func getQueueClientForCommand(cmd *cobra.Command) (*azqueue.QueueClient, error) {
	useEmulator := getBoolFlagValue(cmd, FlagStorageEmulator)
	queueName := cmd.Flag(FlagQueue).Value.String()
	queueURL := cmd.Flag(FlagQueueURL).Value.String()
	serviceURL := cmd.Flag(FlagServiceURL).Value.String()

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
			return nil, errors.New("The queue parameter must be provided when the service-url parameter is passed")
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
