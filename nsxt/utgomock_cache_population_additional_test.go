//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func TestUnitNsxt_CachePopulationPartialAndFailedSearches(t *testing.T) {
	for _, tc := range []struct {
		name, resourceType, mode                                            string
		parentPartial, childPartial, childFail, fallbackFail, invalidParent bool
		wantError                                                           bool
	}{
		{name: "search_unavailable", resourceType: resourceTypeGroup, mode: "global", wantError: true},
		{name: "partial_parents", resourceType: resourceTypeGroup, mode: "global", parentPartial: true},
		{name: "invalid_parent", resourceType: resourceTypeGroup, mode: "global", invalidParent: true, wantError: true},
		{name: "child_search_unavailable", resourceType: resourceTypeGatewayPolicy, mode: "global", childFail: true, wantError: true},
		{name: "partial_child_rules", resourceType: resourceTypeGatewayPolicy, mode: "global", childPartial: true},
		{name: "untagged_child_fallback", resourceType: resourceTypeGatewayPolicy, mode: "config_scope"},
		{name: "untagged_child_fallback_unavailable", resourceType: resourceTypeGatewayPolicy, mode: "config_scope", fallbackFail: true, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := newIsolatedMockProviderClient(t)
			m.Host = "https://nsx.invalid"
			m.CommonConfig.MaxRetries = 0
			m.CommonConfig.CacheMode = tc.mode
			m.CommonConfig.contextID = "run"
			d := schema.TestResourceDataRaw(t, cacheTestSchema(), nil)
			d.SetId("parent")
			calls := 0
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				query := r.URL.Query().Get("query")
				child := strings.Contains(query, "resource_type:rule")
				body := `{"results":[{"id":"parent","path":"/infra/domains/default/gateway-policies/parent","resource_type":"GatewayPolicy"}],"result_count":1}`
				status := http.StatusOK
				if tc.resourceType == resourceTypeGroup {
					body = `{"results":[{"id":"parent","resource_type":"Group"}],"result_count":1}`
				}
				if tc.invalidParent {
					body = `{"results":[{"id":42}],"result_count":1}`
				}
				if tc.name == "search_unavailable" {
					status = http.StatusServiceUnavailable
				}
				if tc.parentPartial {
					if r.URL.Query().Get("cursor") == "" {
						body = `{"results":[{"id":"parent"}],"result_count":2,"cursor":"next"}`
					} else {
						status = http.StatusServiceUnavailable
					}
				}
				if child {
					body = `{"results":[],"result_count":0}`
					if tc.childFail {
						status = http.StatusServiceUnavailable
					}
					if tc.childPartial {
						if r.URL.Query().Get("cursor") == "" {
							body = `{"results":[{"id":"rule","parent_path":"/infra/domains/default/gateway-policies/parent","resource_type":"Rule"}],"result_count":2,"cursor":"next"}`
						} else {
							status = http.StatusServiceUnavailable
						}
					}
					if tc.fallbackFail && !strings.Contains(query, "tags.scope") {
						status = http.StatusServiceUnavailable
					}
				}
				if status != http.StatusOK {
					body = `{"error_code":503,"error_message":"search unavailable","module_name":"test"}`
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}
			cache := &resourceTypeCache{data: make(map[string]map[string]*data.StructValue)}
			events, err := cache.getListOfPolicyResources("query", d, m, getPolicyConnector(m), utl.SessionContext{ClientType: utl.Local}, tc.resourceType, "run")
			if tc.wantError {
				require.Error(t, err)
				assert.NotContains(t, cache.data, "query")
			} else {
				require.NoError(t, err)
				require.Contains(t, cache.data["query"], "parent")
			}
			require.Positive(t, calls)
			if tc.parentPartial || tc.childPartial {
				found := false
				for _, event := range events {
					if event.level == "WARNING" {
						found = true
					}
				}
				assert.True(t, found, "partial results must emit a warning")
			}
			if tc.name == "untagged_child_fallback" {
				assert.Equal(t, 3, calls)
			}
		})
	}
}

func TestUnitNsxt_ManagedTagReflectionEdgeCases(t *testing.T) {
	m := newIsolatedMockProviderClientCacheEnabled(t)
	var nilObject *struct{ Tags []string }
	for _, obj := range []interface{}{nil, nilObject, 42, struct{ Other string }{}, struct{ Tags string }{}, struct{ Tags []string }{Tags: []string{"unrelated"}}} {
		// The stripping helper is intentionally tolerant of unrelated model types.
		stripProviderManagedTagsFromAny(obj)
	}
	_, err := ensureProviderManagedTagsWithPatchFunc(nilObject, m, func(*struct{ Tags []string }) error { t.Fatal("nil object patched"); return nil })
	require.NoError(t, err)
	_, err = ensureProviderManagedTagsWithPatchFunc(42, m, func(int) error { t.Fatal("scalar patched"); return nil })
	require.NoError(t, err)
	_, err = ensureProviderManagedTagsWithPatchFunc(struct{ Other string }{}, m, func(struct{ Other string }) error { t.Fatal("model without tags patched"); return nil })
	require.NoError(t, err)
}
