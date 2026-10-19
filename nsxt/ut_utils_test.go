//go:build unittest

// Shared helpers for unittest-tagged mock tests (see utgomock_*_test.go).

package nsxt

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"

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

func configurationChangeData(t *testing.T, res *schema.Resource, oldRaw, newRaw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	old := schema.TestResourceDataRaw(t, res.Schema, oldRaw)
	old.SetId("resource")
	diff, err := schema.InternalMap(res.Schema).Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(newRaw), nil, nil, false)
	require.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(old.State(), diff)
	require.NoError(t, err)
	return d
}

// Preserve a valid discriminator while corrupting a subtype field. This ensures
// errors from the second conversion are returned to Terraform rather than ignored.
func malformedAPIModel(discriminator, kind, field string) *data.StructValue {
	fields := map[string]data.DataValue{discriminator: data.NewStringValue(kind)}
	if field != "" {
		fields[field] = data.NewIntegerValue(1)
		if field == "ether_type" || field == "protocol_number" {
			fields[field] = data.NewStringValue("invalid")
		}
	}
	return data.NewStructValue("test", fields)
}

type unitPolicyTransport func(*http.Request) (*http.Response, error)

func (f unitPolicyTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// The existing mock client initializes the process-wide version on first use.
// Restore it so tests that construct clients do not change later tests' setup.
func newIsolatedMockProviderClient(t *testing.T) nsxtClients {
	t.Helper()
	originalVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = originalVersion })
	return newGoMockProviderClient()
}

func newIsolatedMockProviderClientCacheEnabled(t *testing.T) nsxtClients {
	t.Helper()
	originalVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = originalVersion })
	return newGoMockProviderClientCacheEnabled()
}

func testIdentityValidationBoundaries(t *testing.T, validate schema.SchemaValidateFunc, limit int, valid, invalid string) {
	t.Helper()

	for _, v := range []interface{}{valid, strings.Repeat("a", limit)} {
		_, errs := validate(v, t.Name())
		require.Empty(t, errs)
	}
	for _, v := range []interface{}{nil, 42, "", invalid, "has space", strings.Repeat("a", limit+1)} {
		_, errs := validate(v, t.Name())
		require.NotEmpty(t, errs, "value %v", v)
	}
}

func testVPNIPv6LocaleServiceDiff(t *testing.T, resource *schema.Resource, field string) {
	t.Helper()

	for _, tc := range []struct {
		name, address, path string
		wantError           bool
	}{
		{"ipv6_locale", "2001:db8::1", "/infra/tier-0s/t0/locale-services/ls/ipsec-vpn-services/vpn", true},
		{"ipv4_locale", "192.0.2.1", "/infra/tier-0s/t0/locale-services/ls/ipsec-vpn-services/vpn", false},
		{"ipv6_gateway", "2001:db8::1", "/infra/tier-0s/t0/ipsec-vpn-services/vpn", false},
		{"hostname_locale", "vpn.example.com", "/infra/tier-0s/t0/locale-services/ls/ipsec-vpn-services/vpn", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Run the real CustomizeDiff against Terraform's diff, so values have the same representation as a plan.
			res := &schema.Resource{Schema: resource.Schema, CustomizeDiff: resource.CustomizeDiff}
			_, err := res.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(map[string]interface{}{field: tc.address, "service_path": tc.path}), newIsolatedMockProviderClient(t))
			if tc.wantError {
				require.ErrorContains(t, err, "IPv6 addresses are not supported")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
