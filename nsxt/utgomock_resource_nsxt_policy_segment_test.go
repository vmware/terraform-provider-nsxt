// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run mockgen for SegmentsClient and segment profile binding map clients in api/infra/segments.

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	segDisplayName = "segment-fooname"
	segDescription = "segment mock description"
	segPath        = "/infra/segments/seg-1"
	segRevision    = int64(1)
	segID          = "seg-1"
)

func TestMockResourceNsxtPolicySegmentRead(t *testing.T) {
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
		mockDiscoverySDK := segprofmocks.NewMockSegmentDiscoveryProfileBindingMapsClient(ctrl)
		mockQosSDK := segprofmocks.NewMockSegmentQosProfileBindingMapsClient(ctrl)
		mockSecuritySDK := segprofmocks.NewMockSegmentSecurityProfileBindingMapsClient(ctrl)

		mockSegmentsSDK.EXPECT().Get(segID).Return(model.Segment{
			DisplayName: &segDisplayName,
			Description: &segDescription,
			Path:        &segPath,
			Revision:    &segRevision,
		}, nil)
		mockDiscoverySDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{}}, nil)

		originalDiscovery := cliSegmentDiscoveryProfileBindingMapsClient
		originalQos := cliSegmentQosProfileBindingMapsClient
		originalSecurity := cliSegmentSecurityProfileBindingMapsClient
		cliSegmentDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentDiscoveryProfileBindingMapClientContext {
			return &segments.SegmentDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
		}
		cliSegmentQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentQosProfileBindingMapClientContext {
			return &segments.SegmentQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
		}
		cliSegmentSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentSecurityProfileBindingMapClientContext {
			return &segments.SegmentSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
		}
		defer func() {
			cliSegmentDiscoveryProfileBindingMapsClient = originalDiscovery
			cliSegmentQosProfileBindingMapsClient = originalQos
			cliSegmentSecurityProfileBindingMapsClient = originalSecurity
		}()

		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"connectivity_path": "",
		})
		d.SetId(segID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, segDisplayName, d.Get("display_name"))
		assert.Equal(t, segDescription, d.Get("description"))
		assert.Equal(t, segPath, d.Get("path"))
		assert.Equal(t, int(segRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"connectivity_path": "",
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})
}

func TestMockResourceNsxtPolicySegmentCreate(t *testing.T) {
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
		mockDiscoverySDK := segprofmocks.NewMockSegmentDiscoveryProfileBindingMapsClient(ctrl)
		mockQosSDK := segprofmocks.NewMockSegmentQosProfileBindingMapsClient(ctrl)
		mockSecuritySDK := segprofmocks.NewMockSegmentSecurityProfileBindingMapsClient(ctrl)
		originalDiscovery := cliSegmentDiscoveryProfileBindingMapsClient
		originalQos := cliSegmentQosProfileBindingMapsClient
		originalSecurity := cliSegmentSecurityProfileBindingMapsClient
		cliSegmentDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentDiscoveryProfileBindingMapClientContext {
			return &segments.SegmentDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
		}
		cliSegmentQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentQosProfileBindingMapClientContext {
			return &segments.SegmentQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
		}
		cliSegmentSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentSecurityProfileBindingMapClientContext {
			return &segments.SegmentSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
		}
		defer func() {
			cliSegmentDiscoveryProfileBindingMapsClient = originalDiscovery
			cliSegmentQosProfileBindingMapsClient = originalQos
			cliSegmentSecurityProfileBindingMapsClient = originalSecurity
		}()

		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockSegmentsSDK.EXPECT().Get(gomock.Any()).Return(model.Segment{
			DisplayName: &segDisplayName,
			Description: &segDescription,
			Path:        &segPath,
			Revision:    &segRevision,
		}, nil)
		mockDiscoverySDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":         segDisplayName,
			"transport_zone_path":  "/infra/sites/default/enforcement-points/default/transport-zones/zone-1",
			"connectivity_path":    "",
			"vlan_ids":             []interface{}{},
			"subnet":               []interface{}{},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, segDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicySegmentUpdate(t *testing.T) {
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
		mockDiscoverySDK := segprofmocks.NewMockSegmentDiscoveryProfileBindingMapsClient(ctrl)
		mockQosSDK := segprofmocks.NewMockSegmentQosProfileBindingMapsClient(ctrl)
		mockSecuritySDK := segprofmocks.NewMockSegmentSecurityProfileBindingMapsClient(ctrl)
		originalDiscovery := cliSegmentDiscoveryProfileBindingMapsClient
		originalQos := cliSegmentQosProfileBindingMapsClient
		originalSecurity := cliSegmentSecurityProfileBindingMapsClient
		cliSegmentDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentDiscoveryProfileBindingMapClientContext {
			return &segments.SegmentDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
		}
		cliSegmentQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentQosProfileBindingMapClientContext {
			return &segments.SegmentQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
		}
		cliSegmentSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentSecurityProfileBindingMapClientContext {
			return &segments.SegmentSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
		}
		defer func() {
			cliSegmentDiscoveryProfileBindingMapsClient = originalDiscovery
			cliSegmentQosProfileBindingMapsClient = originalQos
			cliSegmentSecurityProfileBindingMapsClient = originalSecurity
		}()

		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockSegmentsSDK.EXPECT().Get(segID).Return(model.Segment{
			DisplayName: &segDisplayName,
			Description: &segDescription,
			Path:        &segPath,
			Revision:    &segRevision,
		}, nil)
		mockDiscoverySDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":         segDisplayName,
			"connectivity_path":    "",
			"transport_zone_path":  "/infra/sites/default/enforcement-points/default/transport-zones/zone-1",
			"vlan_ids":             []interface{}{},
			"subnet":               []interface{}{},
		})
		d.SetId(segID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, segDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":        segDisplayName,
			"connectivity_path":   "",
			"transport_zone_path": "/infra/sites/default/enforcement-points/default/transport-zones/zone-1",
			"vlan_ids":            []interface{}{},
			"subnet":              []interface{}{},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})
}

func TestMockResourceNsxtPolicySegmentDelete(t *testing.T) {
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
		mockPortsSDK.EXPECT().List(segID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.SegmentPortListResult{Results: []model.SegmentPort{}}, nil)
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(segID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegment()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})
}
