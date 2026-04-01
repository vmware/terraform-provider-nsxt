//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/EvpnTenantConfigsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/EvpnTenantConfigsClient.go EvpnTenantConfigsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	evpnmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	evpnTenantID           = "evpn-tenant-1"
	evpnTenantDisplayName  = "evpn-tenant-fooname"
	evpnTenantDescription  = "evpn tenant mock"
	evpnTenantPath         = "/infra/evpn-tenant-configs/evpn-tenant-1"
	evpnTenantRevision     = int64(1)
	evpnTenantTzPath       = "/infra/sites/default/enforcement-points/default/transport-zones/tz1"
	evpnTenantVniPoolPath  = "/infra/vni-pools/vni1"
	evpnTenantMappingVlans = "1-100"
	evpnTenantMappingVnis  = "10000-10100"
)

func TestMockResourceNsxtPolicyEvpnTenantCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := evpnmocks.NewMockEvpnTenantConfigsClient(ctrl)
	mockWrapper := &cliinfra.EvpnTenantConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalCli := cliEvpnTenantConfigsClient
	defer func() { cliEvpnTenantConfigsClient = originalCli }()
	cliEvpnTenantConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.EvpnTenantConfigClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockEvpnSDK.EXPECT().Get(gomock.Any()).Return(model.EvpnTenantConfig{
			Id:                &evpnTenantID,
			DisplayName:       &evpnTenantDisplayName,
			Description:       &evpnTenantDescription,
			Path:              &evpnTenantPath,
			Revision:          &evpnTenantRevision,
			TransportZonePath: &evpnTenantTzPath,
			VniPoolPath:       &evpnTenantVniPoolPath,
			Mappings:          []model.VlanVniRangePair{{Vlans: &evpnTenantMappingVlans, Vnis: &evpnTenantMappingVnis}},
		}, nil)

		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEvpnTenantData(t))
		setEvpnTenantMapping(d, t)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Get("existing-id").Return(model.EvpnTenantConfig{Id: &evpnTenantID}, nil)

		res := resourceNsxtPolicyEvpnTenant()
		data := minimalEvpnTenantData(t)
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		setEvpnTenantMapping(d, t)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyEvpnTenantRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := evpnmocks.NewMockEvpnTenantConfigsClient(ctrl)
	mockWrapper := &cliinfra.EvpnTenantConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalCli := cliEvpnTenantConfigsClient
	defer func() { cliEvpnTenantConfigsClient = originalCli }()
	cliEvpnTenantConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.EvpnTenantConfigClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Get(evpnTenantID).Return(model.EvpnTenantConfig{
			Id:                &evpnTenantID,
			DisplayName:       &evpnTenantDisplayName,
			Description:       &evpnTenantDescription,
			Path:              &evpnTenantPath,
			Revision:          &evpnTenantRevision,
			TransportZonePath: &evpnTenantTzPath,
			VniPoolPath:       &evpnTenantVniPoolPath,
			Mappings:          []model.VlanVniRangePair{{Vlans: &evpnTenantMappingVlans, Vnis: &evpnTenantMappingVnis}},
		}, nil)

		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(evpnTenantID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, evpnTenantDisplayName, d.Get("display_name"))
		assert.Equal(t, evpnTenantDescription, d.Get("description"))
		assert.Equal(t, evpnTenantPath, d.Get("path"))
		assert.Equal(t, int(evpnTenantRevision), d.Get("revision"))
		assert.Equal(t, evpnTenantID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Evpn Tenant ID")
	})
}

func TestMockResourceNsxtPolicyEvpnTenantUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := evpnmocks.NewMockEvpnTenantConfigsClient(ctrl)
	mockWrapper := &cliinfra.EvpnTenantConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalCli := cliEvpnTenantConfigsClient
	defer func() { cliEvpnTenantConfigsClient = originalCli }()
	cliEvpnTenantConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.EvpnTenantConfigClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(evpnTenantID, gomock.Any()).Return(nil)
		mockEvpnSDK.EXPECT().Get(evpnTenantID).Return(model.EvpnTenantConfig{
			Id:                &evpnTenantID,
			DisplayName:       &evpnTenantDisplayName,
			Description:       &evpnTenantDescription,
			Path:              &evpnTenantPath,
			Revision:          &evpnTenantRevision,
			TransportZonePath: &evpnTenantTzPath,
			VniPoolPath:       &evpnTenantVniPoolPath,
			Mappings:          []model.VlanVniRangePair{{Vlans: &evpnTenantMappingVlans, Vnis: &evpnTenantMappingVnis}},
		}, nil)

		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEvpnTenantData(t))
		setEvpnTenantMapping(d, t)
		d.SetId(evpnTenantID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEvpnTenantData(t))
		setEvpnTenantMapping(d, t)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Evpn Tenant ID")
	})
}

func TestMockResourceNsxtPolicyEvpnTenantDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := evpnmocks.NewMockEvpnTenantConfigsClient(ctrl)
	mockWrapper := &cliinfra.EvpnTenantConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalCli := cliEvpnTenantConfigsClient
	defer func() { cliEvpnTenantConfigsClient = originalCli }()
	cliEvpnTenantConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.EvpnTenantConfigClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Delete(evpnTenantID).Return(nil)

		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(evpnTenantID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Evpn Tenant ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Delete(evpnTenantID).Return(errors.New("API error"))

		res := resourceNsxtPolicyEvpnTenant()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(evpnTenantID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnTenantDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalEvpnTenantData(t *testing.T) map[string]interface{} {
	return map[string]interface{}{
		"display_name":        evpnTenantDisplayName,
		"description":         evpnTenantDescription,
		"transport_zone_path": evpnTenantTzPath,
		"vni_pool_path":       evpnTenantVniPoolPath,
	}
}

// setEvpnTenantMapping sets the mapping TypeSet on d (must be called after TestResourceDataRaw for Create/Update tests).
func setEvpnTenantMapping(d *schema.ResourceData, t *testing.T) {
	res := resourceNsxtPolicyEvpnTenant()
	mappingElem := res.Schema["mapping"].Elem.(*schema.Resource)
	mappingSet := schema.NewSet(schema.HashResource(mappingElem), []interface{}{
		map[string]interface{}{"vlans": evpnTenantMappingVlans, "vnis": evpnTenantMappingVnis},
	})
	if err := d.Set("mapping", mappingSet); err != nil {
		t.Fatal(err)
	}
}
