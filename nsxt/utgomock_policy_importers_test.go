//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitNsxt_getPolicyPathResourceImporter(t *testing.T) {
	m := newGoMockProviderClient()
	imp := getPolicyPathResourceImporter("/infra/domains/default/groups/example")

	t.Run("valid policy path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("/infra/domains/default/groups/mygroup")

		rd, err := imp(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "mygroup", rd[0].Id())
	})

	t.Run("invalid path format", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("not-a-valid-uri")

		_, err := imp(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path")
	})
}

func TestUnitNsxt_getPolicyPathOrIDResourceImporter(t *testing.T) {
	m := newGoMockProviderClient()
	imp := getPolicyPathOrIDResourceImporter("/infra/foo/bar")

	t.Run("valid id without path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("plain-resource-id")

		rd, err := imp(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
	})

	t.Run("invalid id when not a path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("bad/id")

		_, err := imp(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path or id")
	})
}

func TestUnitNsxt_resourceNsxtPolicyVMTagsImporter(t *testing.T) {
	m := newGoMockProviderClient()
	contextSchema := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"project_id": {Type: schema.TypeString, Required: true, ForceNew: true},
				"vpc_id":     {Type: schema.TypeString, Optional: true, ForceNew: true},
			},
		},
	}

	t.Run("plain instance id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("503b13f0-22c7-4d72-a16e-aa7e6c8f2da6")

		rd, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "503b13f0-22c7-4d72-a16e-aa7e6c8f2da6", rd[0].Id())
	})

	t.Run("composite instance_id::project_id sets context", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("503b13f0-22c7-4d72-a16e-aa7e6c8f2da6::demoproj")

		rd, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "503b13f0-22c7-4d72-a16e-aa7e6c8f2da6", rd[0].Id())
		projectID, _, _ := getContextDataFromSchema(rd[0])
		assert.Equal(t, "demoproj", projectID)
	})

	t.Run("empty id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("")

		_, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.Error(t, err)
	})

	t.Run("empty instance_id in composite", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("::demoproj")

		_, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "non-empty")
	})

	t.Run("empty project_id in composite", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("503b13f0-22c7-4d72-a16e-aa7e6c8f2da6::")

		_, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "non-empty")
	})

	t.Run("slash in id without separator is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("/infra/realized-state/virtual-machines/vm1")

		_, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.Error(t, err)
	})

	t.Run("composite instance_id::project_id::vpc_id sets context", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("503b13f0-22c7-4d72-a16e-aa7e6c8f2da6::demoproj::demovpc")

		rd, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "503b13f0-22c7-4d72-a16e-aa7e6c8f2da6", rd[0].Id())
		projectID, vpcID, _ := getContextDataFromSchema(rd[0])
		assert.Equal(t, "demoproj", projectID)
		assert.Equal(t, "demovpc", vpcID)
	})

	t.Run("empty vpc_id in three-part composite", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("503b13f0-22c7-4d72-a16e-aa7e6c8f2da6::demoproj::")

		_, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "non-empty")
	})

	t.Run("empty project_id in three-part composite", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": contextSchema,
		}, map[string]interface{}{})
		d.SetId("503b13f0-22c7-4d72-a16e-aa7e6c8f2da6::::demovpc")

		_, err := resourceNsxtPolicyVMTagsImporter(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "non-empty")
	})
}

func TestUnitNsxt_getVpcPathResourceImporter(t *testing.T) {
	m := newGoMockProviderClient()
	imp := getVpcPathResourceImporter("/orgs/[org]/projects/[project]/vpcs/[vpc]/tier-1s/[tier1]")

	t.Run("valid vpc path sets context", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": getContextSchema(false, false, true),
		}, map[string]interface{}{})
		d.SetId("/orgs/myorg/projects/myproj/vpcs/myvpc/tier-1s/t1")

		rd, err := imp(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
	})

	t.Run("missing vpc in path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": getContextSchema(false, false, true),
		}, map[string]interface{}{})
		d.SetId("/infra/domains/default/groups/g1")

		_, err := imp(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path")
	})

	t.Run("invalid paths are rejected without changing the resource", func(t *testing.T) {
		base := "/orgs/myorg/projects/myproj/vpcs/myvpc"
		for name, path := range map[string]string{
			"wrong collection":         base + "/lb-pools/t1",
			"similar collection":       base + "/tier-1s-x/t1",
			"missing ID":               base + "/tier-1s/",
			"missing collection":       base + "/t1",
			"extra segment":            base + "/tier-1s/t1/extra",
			"missing VPC":              "/orgs/myorg/projects/myproj/tier-1s/t1",
			"empty project":            "/orgs/myorg/projects//vpcs/myvpc/tier-1s/t1",
			"wrong parent collections": "/orgs/myorg/projects/myproj/vpc-x/myvpc/tier-1s/t1",
			"plain ID":                 "t1",
		} {
			t.Run(name, func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
					"context": getContextSchema(false, false, true),
				}, map[string]interface{}{})
				d.SetId(path)

				_, err := imp(d, m)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "Invalid policy path")
				assert.Equal(t, path, d.Id())
				assert.Empty(t, d.Get("context"))
			})
		}
	})

	t.Run("empty ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": getContextSchema(false, false, true),
		}, map[string]interface{}{})
		d.SetId(" ")

		_, err := imp(d, m)
		require.ErrorIs(t, err, ErrEmptyImportID)
	})
}

// Every resource importing with getVpcPathResourceImporter should accept a path built from its own
// example, and reject the same path with a different collection
func TestUnitNsxt_getVpcPathResourceImporterExamples(t *testing.T) {
	m := newGoMockProviderClient()
	examples := map[string]string{
		"vpc endpoint":           vpcEndpointPathExample,
		"vpc gateway policy":     vpcGatewayPolicyPathExample,
		"vpc group":              vpcGroupPathExample,
		"vpc ip address alloc":   vpcIpAddressAllocationPathExample,
		"vpc lb app profile":     vpcLbAppProfilePathExample,
		"vpc lb client ssl":      vpcLbClientSslProfilePathExample,
		"vpc lb monitor profile": vpcLbMonitorProfilePathExample,
		"vpc lb persistence":     vpcLbPersistenceProfilePathExample,
		"vpc lb pool":            vpcLbPoolPathExample,
		"vpc lb server ssl":      vpcLbServerSslProfilePathExample,
		"vpc lb service":         vpcLbServicePathExample,
		"vpc lb virtual server":  vpcLbVirtualServerPathExample,
		"vpc security policy":    vpcSecurityPolicyPathExample,
		"vpc service endpoint":   vpcServiceEndpointPathExample,
		"vpc static routes":      vpcStaticRoutesPathExample,
		"vpc subnet":             vpcSubnetPathExample,
	}
	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			imp := getVpcPathResourceImporter(example)
			validPath := strings.NewReplacer("[org]", "default", "[project]", "proj", "[vpc]", "vpc1").Replace(example)
			segs := strings.Split(validPath, "/")
			segs[len(segs)-1] = "id1"
			validPath = strings.Join(segs, "/")

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"context": getContextSchema(false, false, true),
			}, map[string]interface{}{})
			d.SetId(validPath)
			rd, err := imp(d, m)
			require.NoError(t, err, validPath)
			require.Len(t, rd, 1)
			assert.Equal(t, "id1", rd[0].Id())

			// Replace the collection which precedes the object ID
			segs[len(segs)-2] = "wrong-collection"
			d = schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"context": getContextSchema(false, false, true),
			}, map[string]interface{}{})
			d.SetId(strings.Join(segs, "/"))
			_, err = imp(d, m)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "Invalid policy path")
		})
	}
}
