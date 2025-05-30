package cmd

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
)

const (
	StorageEmulatorConnectionString = "AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;DefaultEndpointsProtocol=http;QueueEndpoint=http://127.0.0.1:10001/devstoreaccount1"
)

func getQueueClientForStorageEmulator(name string) (*azqueue.QueueClient, error) {
	sc, err := azqueue.NewServiceClientFromConnectionString(StorageEmulatorConnectionString, nil)
	if nil != err {
		return nil, err
	}

	return sc.NewQueueClient(name), nil
}

func getQueueClientForQueueURLWithSASToken(queueURL string) (*azqueue.QueueClient, error) {
	return azqueue.NewQueueClientWithNoCredential(queueURL, nil)
}

func getQueueClientForQueueURLWithDefaultCredential(queueURL string) (*azqueue.QueueClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if nil != err {
		return nil, err
	}

	return azqueue.NewQueueClient(queueURL, cred, nil)
}

func getQueueClientForServiceURLWithDefaultCredential(serviceURL string, name string) (*azqueue.QueueClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if nil != err {
		return nil, err
	}

	sc, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	if nil != err {
		return nil, err
	}

	return sc.NewQueueClient(name), nil
}
