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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier1sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	t1lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	t1GwID            = "t1-gw-1"
	t1DisplayName     = "Test Tier1 Gateway"
	t1Description     = "Test Tier1 gateway"
	t1Revision        = int64(1)
	t1Path            = "/infra/tier-1s/t1-gw-1"
	t1FailoverMode    = "PREEMPTIVE"
	t1HaMode          = "ACTIVE_STANDBY"
	t1DisableFirewall = false
)

func t1GatewayAPIResponse() nsxModel.Tier1 {
	return nsxModel.Tier1{
		Id:              &t1GwID,
		DisplayName:     &t1DisplayName,
		Description:     &t1Description,
		Revision:        &t1Revision,
		Path:            &t1Path,
		FailoverMode:    &t1FailoverMode,
		HaMode:          &t1HaMode,
		DisableFirewall: &t1DisableFirewall,
	}
}

func minimalT1GatewayData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    t1DisplayName,
		"description":     t1Description,
		"failover_mode":   t1FailoverMode,
		"ha_mode":         t1HaMode,
		"enable_firewall": true,
	}
}

func setupT1GatewayMocks(t *testing.T, ctrl *gomock.Controller) (
	*inframocks.MockTier1sClient,
	*inframocks.MockInfraClient,
	*t1lsmocks.MockLocaleServicesClient,
	func(),
) {
	mockT1SDK := inframocks.NewMockTier1sClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	mockLSSDK := t1lsmocks.NewMockLocaleServicesClient(ctrl)

	t1Wrapper := &cliinfra.Tier1ClientContext{
		Client:     mockT1SDK,
		ClientType: utl.Local,
	}
	lsWrapper := &tier1sapi.LocaleServicesClientContext{
		Client:     mockLSSDK,
		ClientType: utl.Local,
	}

	origT1 := cliTier1sClient
	origInfra := cliInfraClient
	origLS := cliTier1LocaleServicesClient

	cliTier1sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier1ClientContext {
		return t1Wrapper
	}
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}
	cliTier1LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier1sapi.LocaleServicesClientContext {
		return lsWrapper
	}

	return mockT1SDK, mockInfraSDK, mockLSSDK, func() {
		cliTier1sClient = origT1
		cliInfraClient = origInfra
		cliTier1LocaleServicesClient = origLS
	}
}

func TestMockResourceNsxtPolicyTier1GatewayCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockT1SDK, mockInfraSDK, mockLSSDK, restore := setupT1GatewayMocks(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		resultCount := int64(0)
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockT1SDK.EXPECT().Get(gomock.Any()).Return(t1GatewayAPIResponse(), nil),
			mockLSSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nsxModel.LocaleServicesListResult{Results: []nsxModel.LocaleServices{}, ResultCount: &resultCount}, nil),
		)

		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())

		err := resourceNsxtPolicyTier1GatewayCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, t1DisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyTier1GatewayRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockT1SDK, _, mockLSSDK, restore := setupT1GatewayMocks(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		resultCount := int64(0)
		mockT1SDK.EXPECT().Get(t1GwID).Return(t1GatewayAPIResponse(), nil)
		mockLSSDK.EXPECT().List(t1GwID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LocaleServicesListResult{Results: []nsxModel.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())
		d.SetId(t1GwID)

		err := resourceNsxtPolicyTier1GatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, t1DisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockT1SDK.EXPECT().Get(t1GwID).Return(nsxModel.Tier1{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())
		d.SetId(t1GwID)

		err := resourceNsxtPolicyTier1GatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())

		err := resourceNsxtPolicyTier1GatewayRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier1GatewayUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockT1SDK, mockInfraSDK, mockLSSDK, restore := setupT1GatewayMocks(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		resultCount := int64(0)
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockT1SDK.EXPECT().Get(t1GwID).Return(t1GatewayAPIResponse(), nil),
			mockLSSDK.EXPECT().List(t1GwID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nsxModel.LocaleServicesListResult{Results: []nsxModel.LocaleServices{}, ResultCount: &resultCount}, nil),
		)

		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())
		d.SetId(t1GwID)

		err := resourceNsxtPolicyTier1GatewayUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())

		err := resourceNsxtPolicyTier1GatewayUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier1GatewayDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, mockInfraSDK, _, restore := setupT1GatewayMocks(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())
		d.SetId(t1GwID)

		err := resourceNsxtPolicyTier1GatewayDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier1Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalT1GatewayData())

		err := resourceNsxtPolicyTier1GatewayDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
