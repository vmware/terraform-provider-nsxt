//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/realized_state/RealizedEntitiesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/RealizedEntitiesClient.go RealizedEntitiesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	realizedstate "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	realizedmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state"
)

var giRealizationGatewayInterfacePath = "/infra/tier-0s/t0-1/locale-services/default/interfaces/iface-1"

func setupRealizedEntitiesMock(t *testing.T, ctrl *gomock.Controller) (*realizedmocks.MockRealizedEntitiesClient, func()) {
	mockSDK := realizedmocks.NewMockRealizedEntitiesClient(ctrl)
	mockWrapper := &realizedstate.RealizedEntityClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliRealizedEntitiesClient
	cliRealizedEntitiesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *realizedstate.RealizedEntityClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliRealizedEntitiesClient = original }
}

func TestMockDataSourceNsxtPolicyGatewayInterfaceRealizationRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRealizedEntitiesMock(t, ctrl)
	defer restore()

	ds := dataSourceNsxtPolicyGatewayInterfaceRealization()

	t.Run("Read success by id", func(t *testing.T) {
		state := "REALIZED"
		giID := "gi-1"
		ipKey := "IpAddresses"
		macKey := "MacAddress"
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{
					Id:    &giID,
					State: &state,
					ExtendedAttributes: []nsxModel.AttributeVal{
						{Key: &ipKey, Values: []string{"10.0.0.1"}},
						{Key: &macKey, Values: []string{"00:11:22:33:44:55"}},
					},
				},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           giID,
			"gateway_path": giRealizationGatewayInterfacePath,
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, giID, d.Id())
		assert.Equal(t, state, d.Get("state"))
		assert.Equal(t, "00:11:22:33:44:55", d.Get("mac_address"))
	})

	t.Run("Read site_path set on local manager fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"site_path":    "/infra/sites/site-1",
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})

	t.Run("Read list error", func(t *testing.T) {
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{}, errors.New("list failed"))

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read site_path missing on global manager fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
		})
		c := newGoMockProviderClient()
		c.PolicyGlobalManager = true

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site_path")
	})

	t.Run("Read by display_name perfect match", func(t *testing.T) {
		state := "REALIZED"
		giID := "gi-2"
		name := "my-iface"
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{Id: &giID, State: &state, DisplayName: &name},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"display_name": name,
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, giID, d.Id())
	})

	t.Run("Read by display_name multiple perfect matches fails", func(t *testing.T) {
		state := "REALIZED"
		name := "dup-iface"
		id1, id2 := "gi-3", "gi-4"
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{Id: &id1, State: &state, DisplayName: &name},
				{Id: &id2, State: &state, DisplayName: &name},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"display_name": name,
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple gateway interfaces with name")
	})

	t.Run("Read by display_name contains match", func(t *testing.T) {
		state := "REALIZED"
		giID := "gi-5"
		name := "prefix-my-iface-suffix"
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{Id: &giID, State: &state, DisplayName: &name},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"display_name": "my-iface",
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, giID, d.Id())
	})

	t.Run("Read by display_name multiple contains matches fails", func(t *testing.T) {
		state := "REALIZED"
		id1, id2 := "gi-6", "gi-7"
		name1, name2 := "aa-my-iface", "bb-my-iface"
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{Id: &id1, State: &state, DisplayName: &name1},
				{Id: &id2, State: &state, DisplayName: &name2},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"display_name": "my-iface",
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple gateway interfaces whose name contains")
	})

	t.Run("Read with no id or display_name returns first result", func(t *testing.T) {
		state := "REALIZED"
		giID := "gi-8"
		mockSDK.EXPECT().List(giRealizationGatewayInterfacePath, (*string)(nil)).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{Id: &giID, State: &state},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"gateway_path": giRealizationGatewayInterfacePath,
			"delay":        0,
			"timeout":      5,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, giID, d.Id())
	})
}
