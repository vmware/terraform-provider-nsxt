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

// newGoMockProviderClientCacheEnabled returns a provider client with config_scope caching
// turned on and a non-empty contextID (so getProviderManagedDefaultTags returns a tag to
// reconcile), for exercising a resource Read's isCacheEnabledForRead/CacheAwareResourceRead
// branch - in particular its patchFunc closure, which config_scope mode only invokes when
// the object read back is missing the provider-managed tag for this contextID.
//
// Pair this with a cache-miss: precompute the resource's own cache bucket key via
// getCacheQueryKey(resourceType, d, m) after building d, then before calling the resource's
// Read function do:
//
//	tc := gcache.getTypeCache(resourceType)
//	tc.data[query] = map[string]*data.StructValue{}   // present bucket, no entry: forces a miss without a live search
//	defer delete(gcache.byTyp, resourceType)
//
// (See TestUnitNsxt_cacheAwareResourceRead / TestUnitNsxt_invalidateCacheForResourceType in
// utgomock_cache_test.go for the same pattern used against the generic cache functions
// directly.)
func newGoMockProviderClientCacheEnabled() nsxtClients {
	m := newGoMockProviderClient()
	m.CommonConfig.CacheMode = "config_scope"
	m.CommonConfig.contextID = "ut-run-1"
	return m
}
