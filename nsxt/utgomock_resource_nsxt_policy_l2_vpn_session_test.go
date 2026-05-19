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

	l2vpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/l2vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	l2vpnMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/l2vpn_services"
)

var (
	l2SessionID          = "l2-session-1"
	l2SessionDisplayName = "Test L2 VPN Session"
	l2SessionDescription = "Test l2 vpn session"
	l2SessionRevision    = int64(1)
	l2SessionServicePath = "/infra/tier-1s/t1-gw-1/l2vpn-services/l2-svc-1"
	l2SessionGwID        = "t1-gw-1"
	l2SessionSvcID       = "l2-svc-1"
)

func l2SessionAPIResponse() nsxModel.L2VPNSession {
	return nsxModel.L2VPNSession{
		Id:          &l2SessionID,
		DisplayName: &l2SessionDisplayName,
		Description: &l2SessionDescription,
		Revision:    &l2SessionRevision,
	}
}

func minimalL2SessionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": l2SessionDisplayName,
		"description":  l2SessionDescription,
		"nsx_id":       l2SessionID,
		"service_path": l2SessionServicePath,
	}
}

func setupL2SessionMock(t *testing.T, ctrl *gomock.Controller) (*l2vpnMocks.MockSessionsClient, func()) {
	mockSDK := l2vpnMocks.NewMockSessionsClient(ctrl)
	mockWrapper := &l2vpnapi.L2VPNSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliT1L2vpnSessionsClient
	cliT1L2vpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *l2vpnapi.L2VPNSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliT1L2vpnSessionsClient = original }
}

func TestMockResourceNsxtPolicyL2VPNSessionCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nsxModel.L2VPNSession{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(l2SessionGwID, l2SessionSvcID, l2SessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2SessionID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VPNSessionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2SessionDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nsxModel.L2VPNSession{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VPNSessionUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(l2SessionGwID, l2SessionSvcID, l2SessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VPNSessionDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nil)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockNormalizeL2VpnTransportTunnelPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "flat T0 path returned unchanged",
			input:    "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "locale-services default segment stripped from T0 path",
			input:    "/infra/tier-0s/t0/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "non-default locale-service ID stripped",
			input:    "/infra/tier-0s/t0/locale-services/ls-custom/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-0s/t0/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "flat T1 path returned unchanged",
			input:    "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "locale-services default segment stripped from T1 path",
			input:    "/infra/tier-1s/t1/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1",
			expected: "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1",
		},
		{
			name:     "empty string returned unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "path ending at locale-services keyword without trailing slash unchanged",
			input:    "/infra/tier-0s/t0/locale-services",
			expected: "/infra/tier-0s/t0/locale-services",
		},
		{
			name:     "path with locale-service ID but no further segment unchanged",
			input:    "/infra/tier-0s/t0/locale-services/default",
			expected: "/infra/tier-0s/t0/locale-services/default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeL2VpnTransportTunnelPath(tt.input))
		})
	}
}

func TestMockSuppressL2VpnTransportTunnelsDiff(t *testing.T) {
	flatPath := "/infra/tier-1s/t1/ipsec-vpn-services/svc1/sessions/sess1"
	localeServicePath := "/infra/tier-1s/t1/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1"
	localeServicePathCustom := "/infra/tier-1s/t1/locale-services/custom/ipsec-vpn-services/svc1/sessions/sess1"
	differentPath := "/infra/tier-1s/t1/ipsec-vpn-services/svc2/sessions/sess2"

	tests := []struct {
		name     string
		k        string
		old      string
		new      string
		expected bool
	}{
		{
			name:     "count key is equal",
			k:        "transport_tunnels.#",
			old:      "1",
			new:      "1",
			expected: true,
		},
		{
			name:     "count key differs",
			k:        "transport_tunnels.#",
			old:      "1",
			new:      "2",
			expected: false,
		},
		{
			name:     "locale-service-scoped API path vs flat config path suppressed",
			k:        "transport_tunnels.0",
			old:      localeServicePath,
			new:      flatPath,
			expected: true,
		},
		{
			name:     "non-default locale-service path vs flat config path suppressed",
			k:        "transport_tunnels.0",
			old:      localeServicePathCustom,
			new:      flatPath,
			expected: true,
		},
		{
			name:     "flat path vs flat path is equal",
			k:        "transport_tunnels.0",
			old:      flatPath,
			new:      flatPath,
			expected: true,
		},
		{
			name:     "genuinely different paths are not suppressed",
			k:        "transport_tunnels.0",
			old:      flatPath,
			new:      differentPath,
			expected: false,
		},
		{
			name:     "empty old vs non-empty new is not suppressed",
			k:        "transport_tunnels.0",
			old:      "",
			new:      flatPath,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suppressL2VpnTransportTunnelsDiff(tt.k, tt.old, tt.new, nil)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMockResourceNsxtPolicyL2VPNSessionReadTransportTunnels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	flatPath := "/infra/tier-1s/" + l2SessionGwID + "/ipsec-vpn-services/svc1/sessions/sess1"
	localeServicePath := "/infra/tier-1s/" + l2SessionGwID + "/locale-services/default/ipsec-vpn-services/svc1/sessions/sess1"

	t.Run("Read stores locale-service-scoped API path in state unchanged", func(t *testing.T) {
		resp := l2SessionAPIResponse()
		resp.TransportTunnels = []string{localeServicePath}
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(resp, nil)

		res := resourceNsxtPolicyL2VPNSession()
		data := minimalL2SessionData()
		data["transport_tunnels"] = []interface{}{flatPath}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		tunnels := d.Get("transport_tunnels").([]interface{})
		require.Len(t, tunnels, 1)
		assert.Equal(t, localeServicePath, tunnels[0].(string), "Read must store the API-returned path verbatim; DiffSuppressFunc handles the comparison")
	})

	t.Run("Read stores flat API path in state unchanged", func(t *testing.T) {
		resp := l2SessionAPIResponse()
		resp.TransportTunnels = []string{flatPath}
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(resp, nil)

		res := resourceNsxtPolicyL2VPNSession()
		data := minimalL2SessionData()
		data["transport_tunnels"] = []interface{}{flatPath}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		tunnels := d.Get("transport_tunnels").([]interface{})
		require.Len(t, tunnels, 1)
		assert.Equal(t, flatPath, tunnels[0].(string))
	})

	t.Run("Read stores multiple API paths in state unchanged", func(t *testing.T) {
		localeServicePath2 := "/infra/tier-1s/" + l2SessionGwID + "/locale-services/default/ipsec-vpn-services/svc2/sessions/sess2"
		flatPath2 := "/infra/tier-1s/" + l2SessionGwID + "/ipsec-vpn-services/svc2/sessions/sess2"

		resp := l2SessionAPIResponse()
		resp.TransportTunnels = []string{localeServicePath, localeServicePath2}
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(resp, nil)

		res := resourceNsxtPolicyL2VPNSession()
		data := minimalL2SessionData()
		data["transport_tunnels"] = []interface{}{flatPath, flatPath2}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		tunnels := d.Get("transport_tunnels").([]interface{})
		require.Len(t, tunnels, 2)
		assert.Equal(t, localeServicePath, tunnels[0].(string))
		assert.Equal(t, localeServicePath2, tunnels[1].(string))
	})

	t.Run("Read always resets transport_tunnels when API returns empty list", func(t *testing.T) {
		resp := l2SessionAPIResponse()
		resp.TransportTunnels = []string{}
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(resp, nil)

		res := resourceNsxtPolicyL2VPNSession()
		data := minimalL2SessionData()
		data["transport_tunnels"] = []interface{}{flatPath}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		tunnels := d.Get("transport_tunnels").([]interface{})
		assert.Empty(t, tunnels, "stale transport_tunnels must be cleared when API returns empty list")
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nsxModel.L2VPNSession{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyL2VPNSession()
		data := minimalL2SessionData()
		data["transport_tunnels"] = []interface{}{flatPath}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
