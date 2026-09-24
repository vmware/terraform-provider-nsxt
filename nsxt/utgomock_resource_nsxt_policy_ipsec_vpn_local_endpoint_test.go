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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	t0ipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/ipsec_vpn_services"
	t0nestedipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services/ipsec_vpn_services"
	ipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/ipsec_vpn_services"
	t1nestedipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/locale_services/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0IpsecMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/ipsec_vpn_services"
	t0nestedIpsecMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services/ipsec_vpn_services"
	t1IpsecMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/ipsec_vpn_services"
	t1nestedIpsecMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/locale_services/ipsec_vpn_services"
)

var (
	ipsecEPID          = "ep-1"
	ipsecEPDisplayName = "Test IPSec Local Endpoint"
	ipsecEPDescription = "Test ipsec local endpoint"
	ipsecEPRevision    = int64(1)
	ipsecEPServicePath = "/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1"
	ipsecEPGwID        = "t1-gw-1"
	ipsecEPSvcID       = "svc-1"
	ipsecEPLocalAddr   = "192.168.1.1"
)

func ipsecEPAPIResponse() nsxModel.IPSecVpnLocalEndpoint {
	return nsxModel.IPSecVpnLocalEndpoint{
		Id:           &ipsecEPID,
		DisplayName:  &ipsecEPDisplayName,
		Description:  &ipsecEPDescription,
		Revision:     &ipsecEPRevision,
		LocalAddress: &ipsecEPLocalAddr,
	}
}

func minimalIPSecEPData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  ipsecEPDisplayName,
		"description":   ipsecEPDescription,
		"nsx_id":        ipsecEPID,
		"service_path":  ipsecEPServicePath,
		"local_address": ipsecEPLocalAddr,
	}
}

func setupIPSecEPMock(t *testing.T, ctrl *gomock.Controller) (*t1IpsecMocks.MockLocalEndpointsClient, func()) {
	mockSDK := t1IpsecMocks.NewMockLocalEndpointsClient(ctrl)
	mockWrapper := &ipsecvpnapi.IPSecVpnLocalEndpointClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnLocalEndpointsClient
	cliTier1IpsecVpnLocalEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ipsecvpnapi.IPSecVpnLocalEndpointClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnLocalEndpointsClient = original }
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(nsxModel.IPSecVpnLocalEndpoint{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipsecEPGwID, ipsecEPSvcID, ipsecEPID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(ipsecEPAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecEPID, d.Id())
	})
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(ipsecEPAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecEPDisplayName, d.Get("display_name"))
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(nsxModel.IPSecVpnLocalEndpoint{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ipsecEPGwID, ipsecEPSvcID, ipsecEPID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(ipsecEPAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_localEndpointClientAllDispatchPaths(t *testing.T) {
	t.Run("tier0 flat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK := t0IpsecMocks.NewMockLocalEndpointsClient(ctrl)
		wrapper := &t0ipsecvpnapi.IPSecVpnLocalEndpointClientContext{Client: mockSDK, ClientType: utl.Local}
		original := cliTier0IpsecVpnLocalEndpointsClient
		defer func() { cliTier0IpsecVpnLocalEndpointsClient = original }()
		cliTier0IpsecVpnLocalEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0ipsecvpnapi.IPSecVpnLocalEndpointClientContext {
			return wrapper
		}
		c := &localEndpointClient{isT0: true, gwID: "t0-1", serviceID: "svc-1", sessionContext: utl.SessionContext{ClientType: utl.Local}}

		mockSDK.EXPECT().Get("t0-1", "svc-1", "ep-1").Return(ipsecEPAPIResponse(), nil)
		_, err := c.Get(nil, "ep-1")
		require.NoError(t, err)

		mockSDK.EXPECT().Patch("t0-1", "svc-1", "ep-1", gomock.Any()).Return(nil)
		require.NoError(t, c.Patch(nil, "ep-1", ipsecEPAPIResponse()))

		mockSDK.EXPECT().Delete("t0-1", "svc-1", "ep-1").Return(nil)
		require.NoError(t, c.Delete(nil, "ep-1"))
	})

	t.Run("tier0 nested locale service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK := t0nestedIpsecMocks.NewMockLocalEndpointsClient(ctrl)
		wrapper := &t0nestedipsecvpnapi.IPSecVpnLocalEndpointClientContext{Client: mockSDK, ClientType: utl.Local}
		original := cliTier0LocaleServiceIpsecVpnLocalEndpointsClient
		defer func() { cliTier0LocaleServiceIpsecVpnLocalEndpointsClient = original }()
		cliTier0LocaleServiceIpsecVpnLocalEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0nestedipsecvpnapi.IPSecVpnLocalEndpointClientContext {
			return wrapper
		}
		c := &localEndpointClient{isT0: true, gwID: "t0-1", localeServiceID: "ls-1", serviceID: "svc-1", sessionContext: utl.SessionContext{ClientType: utl.Local}}

		mockSDK.EXPECT().Get("t0-1", "ls-1", "svc-1", "ep-1").Return(ipsecEPAPIResponse(), nil)
		_, err := c.Get(nil, "ep-1")
		require.NoError(t, err)

		mockSDK.EXPECT().Patch("t0-1", "ls-1", "svc-1", "ep-1", gomock.Any()).Return(nil)
		require.NoError(t, c.Patch(nil, "ep-1", ipsecEPAPIResponse()))

		mockSDK.EXPECT().Delete("t0-1", "ls-1", "svc-1", "ep-1").Return(nil)
		require.NoError(t, c.Delete(nil, "ep-1"))
	})

	t.Run("tier1 nested locale service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK := t1nestedIpsecMocks.NewMockLocalEndpointsClient(ctrl)
		wrapper := &t1nestedipsecvpnapi.IPSecVpnLocalEndpointClientContext{Client: mockSDK, ClientType: utl.Local}
		original := cliTier1LocaleServiceIpsecVpnLocalEndpointsClient
		defer func() { cliTier1LocaleServiceIpsecVpnLocalEndpointsClient = original }()
		cliTier1LocaleServiceIpsecVpnLocalEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1nestedipsecvpnapi.IPSecVpnLocalEndpointClientContext {
			return wrapper
		}
		c := &localEndpointClient{isT0: false, gwID: "t1-1", localeServiceID: "ls-1", serviceID: "svc-1", sessionContext: utl.SessionContext{ClientType: utl.Local}}

		mockSDK.EXPECT().Get("t1-1", "ls-1", "svc-1", "ep-1").Return(ipsecEPAPIResponse(), nil)
		_, err := c.Get(nil, "ep-1")
		require.NoError(t, err)

		mockSDK.EXPECT().Patch("t1-1", "ls-1", "svc-1", "ep-1", gomock.Any()).Return(nil)
		require.NoError(t, c.Patch(nil, "ep-1", ipsecEPAPIResponse()))

		mockSDK.EXPECT().Delete("t1-1", "ls-1", "svc-1", "ep-1").Return(nil)
		require.NoError(t, c.Delete(nil, "ep-1"))
	})

	t.Run("tier1 nested locale service with project context is rejected", func(t *testing.T) {
		c := &localEndpointClient{
			isT0:            false,
			gwID:            "t1-1",
			localeServiceID: "ls-1",
			serviceID:       "svc-1",
			sessionContext:  utl.SessionContext{ClientType: utl.Multitenancy, ProjectID: "proj-1"},
		}

		_, err := c.Get(nil, "ep-1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project context")

		err = c.Patch(nil, "ep-1", ipsecEPAPIResponse())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project context")

		err = c.Delete(nil, "ep-1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project context")
	})
}

func TestUnitNsxt_getIPSecVpnLocalEndpointSessionContext(t *testing.T) {
	res := resourceNsxtPolicyIPSecVpnLocalEndpoint()

	t.Run("non-project path leaves ClientType unchanged", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		ctx := getIPSecVpnLocalEndpointSessionContext(d, newGoMockProviderClient(), ipsecEPServicePath)
		assert.EqualValues(t, utl.Local, ctx.ClientType)
	})

	t.Run("project-scoped path sets Multitenancy client type", func(t *testing.T) {
		data := minimalIPSecEPData()
		data["service_path"] = "/orgs/default/projects/proj-1/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		ctx := getIPSecVpnLocalEndpointSessionContext(d, newGoMockProviderClient(), "/orgs/default/projects/proj-1/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1")
		assert.EqualValues(t, utl.Multitenancy, ctx.ClientType)
		assert.Equal(t, "proj-1", ctx.ProjectID)
	})
}
