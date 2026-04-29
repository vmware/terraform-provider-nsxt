//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
)

func setupTransportZoneDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockTransportZonesClient, func()) {
	t.Helper()
	mockTZSDK := epmocks.NewMockTransportZonesClient(ctrl)
	mockTZWrapper := &enforcementpoints.PolicyTransportZoneClientContext{
		Client:     mockTZSDK,
		ClientType: utl.Local,
	}
	orig := cliTransportZonesClient
	cliTransportZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyTransportZoneClientContext {
		return mockTZWrapper
	}
	return mockTZSDK, func() { cliTransportZonesClient = orig }
}

func TestMockDataSourceNsxtPolicyTransportZoneRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTransportZoneDataSourceMock(t, ctrl)
	defer restore()

	includeFalse := false
	m := newGoMockProviderClient()

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, getPolicyEnforcementPoint(m), tzID).Return(tzAPIResponse(), nil)

		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": tzID,
		})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, tzID, d.Id())
		assert.Equal(t, tzDisplayName, d.Get("display_name"))
		assert.Equal(t, tzType, d.Get("transport_type"))
	})

	t.Run("by id read error", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, getPolicyEnforcementPoint(m), tzID).Return(nsxModel.PolicyTransportZone{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": tzID,
		})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.Error(t, err)
	})

	t.Run("by display_name single match from list", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, getPolicyEnforcementPoint(m), nil, &includeFalse, nil, nil, &includeFalse, nil).Return(nsxModel.PolicyTransportZoneListResult{
			Results: []nsxModel.PolicyTransportZone{tzAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": tzDisplayName,
		})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, tzID, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, getPolicyEnforcementPoint(m), nil, &includeFalse, nil, nil, &includeFalse, nil).Return(nsxModel.PolicyTransportZoneListResult{}, vapiErrors.ServiceUnavailable{})

		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": tzDisplayName,
		})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.Error(t, err)
	})

	t.Run("is_default and transport_type selects default zone", func(t *testing.T) {
		def := true
		tzDef := tzAPIResponse()
		tzDef.IsDefault = &def

		mockSDK.EXPECT().List(defaultSite, getPolicyEnforcementPoint(m), nil, &includeFalse, nil, nil, &includeFalse, nil).Return(nsxModel.PolicyTransportZoneListResult{
			Results: []nsxModel.PolicyTransportZone{tzDef},
		}, nil)

		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"is_default":     true,
			"transport_type": tzType,
		})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, tzID, d.Id())
		assert.True(t, d.Get("is_default").(bool))
	})

	t.Run("missing id name and default selection criteria", func(t *testing.T) {
		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "please specify")
	})

	t.Run("multiple exact display_name matches", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, getPolicyEnforcementPoint(m), nil, &includeFalse, nil, nil, &includeFalse, nil).Return(nsxModel.PolicyTransportZoneListResult{
			Results: []nsxModel.PolicyTransportZone{tzAPIResponse(), tzAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicyTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": tzDisplayName,
		})

		err := dataSourceNsxtPolicyTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})
}
