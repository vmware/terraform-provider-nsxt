package nsxt

import (
	api "github.com/vmware/go-vmware-nsxt"
)

func constructMockProviderClient() nsxtClients {
	commonConfig := commonProviderConfig{
		RemoteAuth:             false,
		ToleratePartialSuccess: false,
		MaxRetries:             2,
		MinRetryInterval:       0,
		MaxRetryInterval:       0,
		RetryStatusCodes:       []int{404, 400},
		Username:               "username",
		Password:               "password",
	}

	nsxtClient := nsxtClients{
		CommonConfig: commonConfig,
	}
	nsxtClient.NsxtClientConfig = &api.Configuration{
		BasePath:   "/api/v1",
		Scheme:     "https",
		UserAgent:  "terraform-provider-nsxt",
		UserName:   "username",
		Password:   "password",
		RemoteAuth: true,
		Insecure:   true,
	}

	return nsxtClient
}
