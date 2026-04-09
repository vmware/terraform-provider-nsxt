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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	"github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	tier1s "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	segprofmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	t1segmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	fixedSegDisplayName = "fixed-segment-test"
	fixedSegDescription = "fixed segment mock description"
	fixedSegPath        = "/infra/tier-1s/t1-gw-for-seg/segments/fixed-seg-1"
	fixedSegRevision    = int64(1)
	fixedSegID          = "fixed-seg-1"
	fixedSegT1Path      = "/infra/tier-1s/t1-gw-for-seg"
	fixedSegT1ID        = "t1-gw-for-seg"
)

func fixedSegAPIResponse() model.Segment {
	return model.Segment{
		Id:               &fixedSegID,
		DisplayName:      &fixedSegDisplayName,
		Description:      &fixedSegDescription,
		Path:             &fixedSegPath,
		Revision:         &fixedSegRevision,
		ConnectivityPath: &fixedSegT1Path,
	}
}

func minimalFixedSegData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":      fixedSegDisplayName,
		"connectivity_path": fixedSegT1Path,
		"subnet":            []interface{}{},
	}
}

func setupFixedSegT1Mocks(t *testing.T, ctrl *gomock.Controller) (*t1segmocks.MockSegmentsClient, func()) {
	mockSDK := t1segmocks.NewMockSegmentsClient(ctrl)
	mockWrapper := &tier1s.SegmentClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTier1SegmentsClient
	cliTier1SegmentsClient = func(_ utl.SessionContext, _ client.Connector) *tier1s.SegmentClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1SegmentsClient = original }
}

func setupFixedSegReadProfileMocks(t *testing.T, ctrl *gomock.Controller) (
	*segprofmocks.MockSegmentDiscoveryProfileBindingMapsClient,
	*segprofmocks.MockSegmentQosProfileBindingMapsClient,
	*segprofmocks.MockSegmentSecurityProfileBindingMapsClient,
	func(),
) {
	mockDiscoverySDK := segprofmocks.NewMockSegmentDiscoveryProfileBindingMapsClient(ctrl)
	mockQosSDK := segprofmocks.NewMockSegmentQosProfileBindingMapsClient(ctrl)
	mockSecuritySDK := segprofmocks.NewMockSegmentSecurityProfileBindingMapsClient(ctrl)

	origDiscovery := cliSegmentDiscoveryProfileBindingMapsClient
	origQos := cliSegmentQosProfileBindingMapsClient
	origSecurity := cliSegmentSecurityProfileBindingMapsClient
	cliSegmentDiscoveryProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *segments.SegmentDiscoveryProfileBindingMapClientContext {
		return &segments.SegmentDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
	}
	cliSegmentQosProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *segments.SegmentQosProfileBindingMapClientContext {
		return &segments.SegmentQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
	}
	cliSegmentSecurityProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *segments.SegmentSecurityProfileBindingMapClientContext {
		return &segments.SegmentSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
	}

	mockDiscoverySDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.SegmentDiscoveryProfileBindingMapListResult{}, nil).AnyTimes()
	mockQosSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.SegmentQosProfileBindingMapListResult{}, nil).AnyTimes()
	mockSecuritySDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.SegmentSecurityProfileBindingMapListResult{}, nil).AnyTimes()

	return mockDiscoverySDK, mockQosSDK, mockSecuritySDK, func() {
		cliSegmentDiscoveryProfileBindingMapsClient = origDiscovery
		cliSegmentQosProfileBindingMapsClient = origQos
		cliSegmentSecurityProfileBindingMapsClient = origSecurity
	}
}

func TestMockResourceNsxtPolicyFixedSegmentRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockT1SegSDK, restoreT1Seg := setupFixedSegT1Mocks(t, ctrl)
	defer restoreT1Seg()

	_, _, _, restoreProfiles := setupFixedSegReadProfileMocks(t, ctrl)
	defer restoreProfiles()

	t.Run("Read success", func(t *testing.T) {
		mockT1SegSDK.EXPECT().Get(fixedSegT1ID, fixedSegID).Return(fixedSegAPIResponse(), nil)

		res := resourceNsxtPolicyFixedSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFixedSegData())
		d.SetId(fixedSegID)

		err := resourceNsxtPolicyFixedSegmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fixedSegDisplayName, d.Get("display_name"))
		assert.Equal(t, fixedSegT1Path, d.Get("connectivity_path"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockT1SegSDK.EXPECT().Get(fixedSegT1ID, fixedSegID).Return(model.Segment{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyFixedSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFixedSegData())
		d.SetId(fixedSegID)

		err := resourceNsxtPolicyFixedSegmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtPolicyFixedSegmentCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockT1SegSDK, restoreT1Seg := setupFixedSegT1Mocks(t, ctrl)
	defer restoreT1Seg()

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	defer func() { cliInfraClient = originalInfra }()
	cliInfraClient = func(_ utl.SessionContext, _ client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	_, _, _, restoreProfiles := setupFixedSegReadProfileMocks(t, ctrl)
	defer restoreProfiles()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockT1SegSDK.EXPECT().Get(fixedSegT1ID, gomock.Any()).Return(fixedSegAPIResponse(), nil),
		)

		res := resourceNsxtPolicyFixedSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFixedSegData())

		err := resourceNsxtPolicyFixedSegmentCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyFixedSegmentUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockT1SegSDK, restoreT1Seg := setupFixedSegT1Mocks(t, ctrl)
	defer restoreT1Seg()

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	defer func() { cliInfraClient = originalInfra }()
	cliInfraClient = func(_ utl.SessionContext, _ client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	_, _, _, restoreProfiles := setupFixedSegReadProfileMocks(t, ctrl)
	defer restoreProfiles()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockT1SegSDK.EXPECT().Get(fixedSegT1ID, fixedSegID).Return(fixedSegAPIResponse(), nil),
		)

		res := resourceNsxtPolicyFixedSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFixedSegData())
		d.SetId(fixedSegID)

		err := resourceNsxtPolicyFixedSegmentUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyFixedSegmentDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	originalInfra := cliInfraClient
	defer func() { cliInfraClient = originalInfra }()
	cliInfraClient = func(_ utl.SessionContext, _ client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Delete success", func(t *testing.T) {
		// Fixed segment deletion uses H-API (policyInfraPatch with marked-for-delete flag)
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyFixedSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"connectivity_path": fixedSegT1Path,
		})
		d.SetId(fixedSegID)

		err := resourceNsxtPolicyFixedSegmentDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
