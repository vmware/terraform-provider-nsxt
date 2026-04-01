//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run mockgen for SegmentsClient, InfraClient, and segment profile binding map clients.

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
	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	segments "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	segmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	segprofmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vlanSegDisplayName = "vlan-segment-fooname"
	vlanSegDescription = "vlan segment mock description"
	vlanSegPath        = "/infra/segments/vlan-seg-1"
	vlanSegRevision    = int64(1)
	vlanSegID          = "vlan-seg-1"
)

func vlanSegAPIResponse() model.Segment {
	return model.Segment{
		Id:          &vlanSegID,
		DisplayName: &vlanSegDisplayName,
		Description: &vlanSegDescription,
		Path:        &vlanSegPath,
		Revision:    &vlanSegRevision,
	}
}

func minimalVlanSegData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":        vlanSegDisplayName,
		"transport_zone_path": "/infra/sites/default/enforcement-points/default/transport-zones/zone-1",
		"vlan_ids":            []interface{}{},
		"subnet":              []interface{}{},
	}
}

func setupVlanSegReadProfileMocks(t *testing.T, ctrl *gomock.Controller) (
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
	cliSegmentDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentDiscoveryProfileBindingMapClientContext {
		return &segments.SegmentDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
	}
	cliSegmentQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentQosProfileBindingMapClientContext {
		return &segments.SegmentQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
	}
	cliSegmentSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentSecurityProfileBindingMapClientContext {
		return &segments.SegmentSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
	}
	return mockDiscoverySDK, mockQosSDK, mockSecuritySDK, func() {
		cliSegmentDiscoveryProfileBindingMapsClient = origDiscovery
		cliSegmentQosProfileBindingMapsClient = origQos
		cliSegmentSecurityProfileBindingMapsClient = origSecurity
	}
}

func TestMockResourceNsxtPolicyVlanSegmentRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSegmentsSDK := segmocks.NewMockSegmentsClient(ctrl)
	mockWrapper := &cliinfra.SegmentClientContext{
		Client:     mockSegmentsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsClient
	defer func() { cliSegmentsClient = originalCli }()
	cliSegmentsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK, cleanup := setupVlanSegReadProfileMocks(t, ctrl)
		defer cleanup()

		mockSegmentsSDK.EXPECT().Get(vlanSegID).Return(vlanSegAPIResponse(), nil)
		mockDiscoverySDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vlanSegID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vlanSegDisplayName, d.Get("display_name"))
		assert.Equal(t, vlanSegDescription, d.Get("description"))
		assert.Equal(t, vlanSegPath, d.Get("path"))
		assert.Equal(t, int(vlanSegRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})

	t.Run("Read clears ID when not found", func(t *testing.T) {
		mockSegmentsSDK.EXPECT().Get(vlanSegID).Return(model.Segment{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vlanSegID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyVlanSegmentCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSegmentsSDK := segmocks.NewMockSegmentsClient(ctrl)
	mockInfraSDK := segmocks.NewMockInfraClient(ctrl)
	mockWrapper := &cliinfra.SegmentClientContext{
		Client:     mockSegmentsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsClient
	originalInfra := cliInfraClient
	defer func() {
		cliSegmentsClient = originalCli
		cliInfraClient = originalInfra
	}()
	cliSegmentsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentClientContext {
		return mockWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Create success", func(t *testing.T) {
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK, cleanup := setupVlanSegReadProfileMocks(t, ctrl)
		defer cleanup()

		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockSegmentsSDK.EXPECT().Get(gomock.Any()).Return(vlanSegAPIResponse(), nil)
		mockDiscoverySDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVlanSegData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, vlanSegDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyVlanSegmentUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSegmentsSDK := segmocks.NewMockSegmentsClient(ctrl)
	mockInfraSDK := segmocks.NewMockInfraClient(ctrl)
	mockWrapper := &cliinfra.SegmentClientContext{
		Client:     mockSegmentsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentsClient
	originalInfra := cliInfraClient
	defer func() {
		cliSegmentsClient = originalCli
		cliInfraClient = originalInfra
	}()
	cliSegmentsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentClientContext {
		return mockWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Update success", func(t *testing.T) {
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK, cleanup := setupVlanSegReadProfileMocks(t, ctrl)
		defer cleanup()

		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockSegmentsSDK.EXPECT().Get(vlanSegID).Return(vlanSegAPIResponse(), nil)
		mockDiscoverySDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVlanSegData())
		d.SetId(vlanSegID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, vlanSegDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVlanSegData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})
}

func TestMockResourceNsxtPolicyVlanSegmentDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInfraSDK := segmocks.NewMockInfraClient(ctrl)
	mockPortsSDK := segprofmocks.NewMockPortsClient(ctrl)

	originalInfra := cliInfraClient
	originalPorts := cliPortsClient
	defer func() {
		cliInfraClient = originalInfra
		cliPortsClient = originalPorts
	}()
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}
	cliPortsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentPortClientContext {
		return &segments.SegmentPortClientContext{Client: mockPortsSDK, ClientType: utl.Local}
	}

	t.Run("Delete success", func(t *testing.T) {
		mockPortsSDK.EXPECT().List(vlanSegID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentPortListResult{Results: []model.SegmentPort{}}, nil)
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vlanSegID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVlanSegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVlanSegmentDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})
}
