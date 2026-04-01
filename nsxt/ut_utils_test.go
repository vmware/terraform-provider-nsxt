//go:build unittest

// Shared helpers for unittest-tagged mock tests (see utgomock_*_test.go).

package nsxt

import (
	api "github.com/vmware/go-vmware-nsxt"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func constructMockProviderClient() nsxtClients {
	// Prevent auto-initialization of NsxVersion via real API calls when empty.
	// Tests that need a specific version set it explicitly before calling this.
	if util.NsxVersion == "" {
		util.NsxVersion = "3.0.0"
	}
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
