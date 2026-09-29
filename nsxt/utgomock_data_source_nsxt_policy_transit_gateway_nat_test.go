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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

const tgwNatPath = "/orgs/default/projects/proj-1/transit-gateways/tgw-1"

func setupTransitGatewayNatMock(t *testing.T, ctrl *gomock.Controller) (*tgwmocks.MockNatClient, func()) {
	mockSDK := tgwmocks.NewMockNatClient(ctrl)
	mockWrapper := &transitgateways.TransitGatewayNatClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "proj-1",
	}
	original := cliTransitGatewayNatClient
	cliTransitGatewayNatClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.TransitGatewayNatClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayNatClient = original }
}

func TestMockDataSourceNsxtPolicyTransitGatewayNatSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTransitGatewayNatMock(t, ctrl)
	defer restore()

	t.Run("Read succeeds for the default USER nat type", func(t *testing.T) {
		id := "USER"
		displayName := "default"
		description := "USER NAT section"
		path := tgwNatPath + "/nat/USER"
		mockSDK.EXPECT().Get("default", "proj-1", "tgw-1", "USER").Return(nsxModel.TransitGatewayNat{
			Id:          &id,
			DisplayName: &displayName,
			Description: &description,
			Path:        &path,
		}, nil)

		ds := dataSourceNsxtPolicyTransitGatewayNat()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"transit_gateway_path": tgwNatPath,
			"nat_type":             "USER",
		})

		err := dataSourceNsxtPolicyTransitGatewayNatRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "USER", d.Id())
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, path, d.Get("path"))
	})

	t.Run("Read fails when the API returns an error", func(t *testing.T) {
		mockSDK.EXPECT().Get("default", "proj-1", "tgw-1", "DEFAULT").Return(nsxModel.TransitGatewayNat{}, assert.AnError)

		ds := dataSourceNsxtPolicyTransitGatewayNat()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"transit_gateway_path": tgwNatPath,
			"nat_type":             "DEFAULT",
		})

		err := dataSourceNsxtPolicyTransitGatewayNatRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}

func TestMockDataSourceNsxtPolicyTransitGatewayNatInvalidPath(t *testing.T) {
	t.Run("Invalid transit_gateway_path returns error immediately", func(t *testing.T) {
		ds := dataSourceNsxtPolicyTransitGatewayNat()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"transit_gateway_path": "/invalid/path",
		})

		err := dataSourceNsxtPolicyTransitGatewayNatRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transit_gateway_path")
	})

	t.Run("Too-short transit_gateway_path returns error immediately", func(t *testing.T) {
		ds := dataSourceNsxtPolicyTransitGatewayNat()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"transit_gateway_path": "/orgs/default",
		})

		err := dataSourceNsxtPolicyTransitGatewayNatRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transit_gateway_path")
	})
}
