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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	tier0s "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tier0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	evpnGatewayPath = "/infra/tier-0s/tier0-gw-1"
	evpnGatewayID   = "tier0-gw-1"
	evpnMode        = model.EvpnConfig_MODE_INLINE
	evpnRevision    = int64(1)
	evpnPath        = "/infra/tier-0s/tier0-gw-1/evpn-config"
	evpnName        = "evpn-config-fooname"
)

func TestMockResourceNsxtPolicyEvpnConfigCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := tier0mocks.NewMockEvpnClient(ctrl)
	evpnWrapper := &tier0s.EvpnConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalEvpn := cliEvpnClient
	defer func() { cliEvpnClient = originalEvpn }()
	cliEvpnClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0s.EvpnConfigClientContext {
		return evpnWrapper
	}

	res := resourceNsxtPolicyEvpnConfig()

	t.Run("Create_success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(evpnGatewayID, gomock.Any()).Return(nil)
		mockEvpnSDK.EXPECT().Get(evpnGatewayID).Return(model.EvpnConfig{
			DisplayName: &evpnName,
			Mode:        &evpnMode,
			Path:        &evpnPath,
			Revision:    &evpnRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": evpnName,
			"gateway_path": evpnGatewayPath,
			"mode":         evpnMode,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, evpnGatewayID, d.Id())
		assert.Equal(t, evpnName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_gateway_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": "/infra/tier-1s/tier1-gw-1",
			"mode":         evpnMode,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0 gateway path expected")
	})

	t.Run("Create_fails_when_gateway_path_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": "",
			"mode":         evpnMode,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gateway_path is not valid")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(evpnGatewayID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": evpnGatewayPath,
			"mode":         evpnMode,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEvpnConfigRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := tier0mocks.NewMockEvpnClient(ctrl)
	evpnWrapper := &tier0s.EvpnConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalEvpn := cliEvpnClient
	defer func() { cliEvpnClient = originalEvpn }()
	cliEvpnClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0s.EvpnConfigClientContext {
		return evpnWrapper
	}

	res := resourceNsxtPolicyEvpnConfig()

	t.Run("Read_success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Get(evpnGatewayID).Return(model.EvpnConfig{
			DisplayName: &evpnName,
			Mode:        &evpnMode,
			Path:        &evpnPath,
			Revision:    &evpnRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": evpnGatewayPath,
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, evpnMode, d.Get("mode"))
		assert.Equal(t, evpnName, d.Get("display_name"))
	})

	t.Run("Read_clears_id_when_not_found", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Get(evpnGatewayID).Return(model.EvpnConfig{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": evpnGatewayPath,
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read_fails_when_gateway_path_is_not_T0", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": "/infra/tier-1s/tier1-gw-1",
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0 gateway path expected")
	})
}

func TestMockResourceNsxtPolicyEvpnConfigUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := tier0mocks.NewMockEvpnClient(ctrl)
	evpnWrapper := &tier0s.EvpnConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalEvpn := cliEvpnClient
	defer func() { cliEvpnClient = originalEvpn }()
	cliEvpnClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0s.EvpnConfigClientContext {
		return evpnWrapper
	}

	res := resourceNsxtPolicyEvpnConfig()

	t.Run("Update_success", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(evpnGatewayID, gomock.Any()).Return(nil)
		mockEvpnSDK.EXPECT().Get(evpnGatewayID).Return(model.EvpnConfig{
			DisplayName: &evpnName,
			Mode:        &evpnMode,
			Path:        &evpnPath,
			Revision:    &evpnRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": evpnName,
			"gateway_path": evpnGatewayPath,
			"mode":         evpnMode,
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Patch_returns_error", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(evpnGatewayID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": evpnGatewayPath,
			"mode":         evpnMode,
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEvpnConfigDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvpnSDK := tier0mocks.NewMockEvpnClient(ctrl)
	evpnWrapper := &tier0s.EvpnConfigClientContext{
		Client:     mockEvpnSDK,
		ClientType: utl.Local,
	}

	originalEvpn := cliEvpnClient
	defer func() { cliEvpnClient = originalEvpn }()
	cliEvpnClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0s.EvpnConfigClientContext {
		return evpnWrapper
	}

	res := resourceNsxtPolicyEvpnConfig()

	t.Run("Delete_success_patches_mode_to_disable", func(t *testing.T) {
		// Delete patches EVPN config with mode=DISABLE instead of a real delete
		mockEvpnSDK.EXPECT().Patch(evpnGatewayID, gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": evpnGatewayPath,
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Patch_returns_error", func(t *testing.T) {
		mockEvpnSDK.EXPECT().Patch(evpnGatewayID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": evpnGatewayPath,
		})
		d.SetId(evpnGatewayID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEvpnConfigDelete(d, m)
		require.Error(t, err)
	})
}
