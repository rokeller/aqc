package cmd

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"

func getQueueClient(name string) (*azqueue.QueueClient, error) {
	sc, err := azqueue.NewServiceClientFromConnectionString("AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;DefaultEndpointsProtocol=http;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;QueueEndpoint=http://127.0.0.1:10001/devstoreaccount1;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	if nil != err {
		return nil, err
	}

	return sc.NewQueueClient(name), nil
}
