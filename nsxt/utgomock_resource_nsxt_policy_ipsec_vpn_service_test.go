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

	tier0ipsecvpnsvcapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/ipsec_vpn_services"
	ipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/ipsec_vpn_services"
	t1localeservicesapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0Mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	t1Mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	t1LocaleServiceMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/locale_services"
)

var (
	ipsecSvcID          = "ipsec-svc-1"
	ipsecSvcDisplayName = "Test IPSec VPN Service"
	ipsecSvcDescription = "Test ipsec vpn service"
	ipsecSvcRevision    = int64(1)
	ipsecSvcGwPath      = "/infra/tier-1s/t1-gw-1"
	ipsecSvcGwID        = "t1-gw-1"
	ipsecSvcPath        = "/infra/tier-1s/t1-gw-1/ipsec-vpn-services/ipsec-svc-1"
)

func ipsecSvcAPIResponse() nsxModel.IPSecVpnService {
	enabled := true
	haSync := true
	logLevel := nsxModel.IPSecVpnService_IKE_LOG_LEVEL_INFO
	return nsxModel.IPSecVpnService{
		Id:          &ipsecSvcID,
		DisplayName: &ipsecSvcDisplayName,
		Description: &ipsecSvcDescription,
		Revision:    &ipsecSvcRevision,
		Path:        &ipsecSvcPath,
		Enabled:     &enabled,
		HaSync:      &haSync,
		IkeLogLevel: &logLevel,
	}
}

func minimalIPSecSvcData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  ipsecSvcDisplayName,
		"description":   ipsecSvcDescription,
		"nsx_id":        ipsecSvcID,
		"gateway_path":  ipsecSvcGwPath,
		"enabled":       true,
		"ha_sync":       true,
		"ike_log_level": nsxModel.IPSecVpnService_IKE_LOG_LEVEL_INFO,
	}
}

func setupIPSecSvcMock(t *testing.T, ctrl *gomock.Controller) (*t1Mocks.MockIpsecVpnServicesClient, func()) {
	mockSDK := t1Mocks.NewMockIpsecVpnServicesClient(ctrl)
	mockWrapper := &ipsecvpnapi.IPSecVpnServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnServicesClient
	cliTier1IpsecVpnServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ipsecvpnapi.IPSecVpnServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnServicesClient = original }
}

func setupIPSecSvcT0Mock(t *testing.T, ctrl *gomock.Controller) (*t0Mocks.MockIpsecVpnServicesClient, func()) {
	mockSDK := t0Mocks.NewMockIpsecVpnServicesClient(ctrl)
	mockWrapper := &tier0ipsecvpnsvcapi.IPSecVpnServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier0IpsecVpnServicesClient
	cliTier0IpsecVpnServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0ipsecvpnsvcapi.IPSecVpnServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0IpsecVpnServicesClient = original }
}

func setupIPSecSvcLocaleServiceMock(t *testing.T, ctrl *gomock.Controller) (*t1LocaleServiceMocks.MockIpsecVpnServicesClient, func()) {
	mockSDK := t1LocaleServiceMocks.NewMockIpsecVpnServicesClient(ctrl)
	mockWrapper := &t1localeservicesapi.IPSecVpnServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnLocaleServicesClient
	cliTier1IpsecVpnLocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1localeservicesapi.IPSecVpnServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnLocaleServicesClient = original }
}

func TestMockResourceNsxtPolicyIPSecVpnServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(nsxModel.IPSecVpnService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipsecSvcGwID, ipsecSvcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcDisplayName, d.Get("display_name"))
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(nsxModel.IPSecVpnService{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(ipsecSvcGwID, ipsecSvcID, gomock.Any()).Return(ipsecSvcAPIResponse(), nil),
			mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecSvcGwID, ipsecSvcID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceT0GatewayPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcT0Mock(t, ctrl)
	defer restore()

	data := minimalIPSecSvcData()
	data["gateway_path"] = "/infra/tier-0s/t0-gw-1"

	t.Run("Read success against a T0 gateway", func(t *testing.T) {
		mockSDK.EXPECT().Get("t0-gw-1", ipsecSvcID).Return(ipsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcDisplayName, d.Get("display_name"))
	})

	t.Run("Delete success against a T0 gateway", func(t *testing.T) {
		mockSDK.EXPECT().Delete("t0-gw-1", ipsecSvcID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceLocaleServicePath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcLocaleServiceMock(t, ctrl)
	defer restore()

	data := minimalIPSecSvcData()
	delete(data, "gateway_path")
	data["locale_service_path"] = "/infra/tier-1s/t1-gw-1/locale-services/default"

	t.Run("Read success against a locale-service-scoped T1 gateway", func(t *testing.T) {
		mockSDK.EXPECT().Get("t1-gw-1", "default", ipsecSvcID).Return(ipsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcDisplayName, d.Get("display_name"))
	})

	t.Run("Delete success against a locale-service-scoped T1 gateway", func(t *testing.T) {
		mockSDK.EXPECT().Delete("t1-gw-1", "default", ipsecSvcID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestUnitNsxt_ipsecVpnServiceMultitenancyLocaleServiceError(t *testing.T) {
	t.Run("locale-service-scoped VPN under a project path is rejected", func(t *testing.T) {
		data := minimalIPSecSvcData()
		delete(data, "gateway_path")
		data["locale_service_path"] = "/orgs/default/projects/proj-1/infra/tier-1s/t1-gw-1/locale-services/default"

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project context is not supported")
	})
}

func TestUnitNsxt_getLocaleServiceAndGatewayPath(t *testing.T) {
	t.Run("fails when neither path is set", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		_, _, err := getLocaleServiceAndGatewayPath(d)
		require.Error(t, err)
	})

	t.Run("succeeds when gateway_path is set", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"gateway_path": ipsecSvcGwPath})

		gw, ls, err := getLocaleServiceAndGatewayPath(d)
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcGwPath, gw)
		assert.Equal(t, "", ls)
	})
}

func TestUnitNsxt_extractProjectIDFromPolicyPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"project path", "/orgs/default/projects/proj-1/infra/tier-1s/t1-gw-1", "proj-1"},
		{"non-project path", "/infra/tier-1s/t1-gw-1", ""},
		{"projects at end with no id", "/orgs/default/projects", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractProjectIDFromPolicyPath(tt.path))
		})
	}
}

func TestUnitNsxt_resourceNsxtPolicyIPSecVpnServiceImport(t *testing.T) {
	res := resourceNsxtPolicyIPSecVpnService()

	t.Run("gateway-scoped path sets gateway_path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/tier-0s/t0-gw-1/ipsec-vpn-services/svc-1")

		out, err := resourceNsxtPolicyIPSecVpnServiceImport(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "svc-1", d.Id())
		assert.Equal(t, "/infra/tier-0s/t0-gw-1", d.Get("gateway_path"))
	})

	t.Run("locale-service-scoped path sets locale_service_path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/tier-0s/t0-gw-1/locale-services/default/ipsec-vpn-services/svc-1")

		out, err := resourceNsxtPolicyIPSecVpnServiceImport(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "svc-1", d.Id())
		assert.Equal(t, "/infra/tier-0s/t0-gw-1/locale-services/default", d.Get("locale_service_path"))
	})

	t.Run("project-scoped path sets context", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/proj-1/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1")

		out, err := resourceNsxtPolicyIPSecVpnServiceImport(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		ctxList := d.Get("context").([]interface{})
		require.Len(t, ctxList, 1)
		assert.Equal(t, "proj-1", ctxList[0].(map[string]interface{})["project_id"])
	})

	t.Run("path without ipsec-vpn-services segment fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/tier-0s/t0-gw-1")

		_, err := resourceNsxtPolicyIPSecVpnServiceImport(d, nil)
		require.Error(t, err)
	})

	t.Run("empty trailing segment fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/tier-0s/t0-gw-1/ipsec-vpn-services/")

		_, err := resourceNsxtPolicyIPSecVpnServiceImport(d, nil)
		require.Error(t, err)
	})
}
