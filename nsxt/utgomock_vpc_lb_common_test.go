//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

const (
	vpcLbPrfOrg     = utl.DefaultOrgID
	vpcLbPrfProject = "proj1"
	vpcLbPrfVpc     = "vpc1"
)

// vpcLbPrfContext returns the "context" block which makes getSessionContext yield a VPC session.
func vpcLbPrfContext() []interface{} {
	return []interface{}{map[string]interface{}{"project_id": vpcLbPrfProject, "vpc_id": vpcLbPrfVpc}}
}

func vpcLbPrfSetNsxVersion(t *testing.T, version string) {
	t.Helper()
	util.NsxVersion = version
	t.Cleanup(func() { util.NsxVersion = "" })
}

// runVpcLbImporterSuite validates the VPC path importer of a derived VPC LB resource.
func runVpcLbImporterSuite(t *testing.T, res *schema.Resource, apiSegment string) {
	t.Helper()
	path := "/orgs/default/projects/" + vpcLbPrfProject + "/vpcs/" + vpcLbPrfVpc + "/" + apiSegment + "/imported-id"

	t.Run("import by VPC path", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.2")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(path)
		out, err := res.Importer.State(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "imported-id", out[0].Id())
		projectID, vpcID, _ := getContextDataFromSchema(out[0])
		assert.Equal(t, vpcLbPrfProject, projectID)
		assert.Equal(t, vpcLbPrfVpc, vpcID)
	})

	t.Run("import by plain ID is rejected", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.2")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("imported-id")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), apiSegment)
	})

	t.Run("import by infra path is rejected", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.2")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/lb-persistence-profiles/imported-id")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("import by project path without VPC is rejected", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.2")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/" + vpcLbPrfProject + "/" + apiSegment + "/imported-id")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("import with old NSX version is rejected", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.1")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(path)
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.2")
	})
}
