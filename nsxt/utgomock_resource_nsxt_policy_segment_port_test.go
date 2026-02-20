// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate mocks for this test, run mockgen for PortsClient and profile binding map clients
// in api/infra/segments and api/infra/segments/ports.

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
	segments "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	portprofiles "github.com/vmware/terraform-provider-nsxt/api/infra/segments/ports"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	portmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	profilemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments/ports"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	segPortSegmentPath = "/infra/segments/seg-1"
	segPortSegmentID   = "seg-1"
	segPortPortID      = "port-1"
	segPortDisplayName = "segment-port-fooname"
	segPortDescription = "segment port mock"
	segPortPath        = "/infra/segments/seg-1/ports/port-1"
	segPortRevision    = int64(1)
)

func TestMockResourceNsxtPolicySegmentPortRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPortsSDK := portmocks.NewMockPortsClient(ctrl)
	mockDiscoverySDK := profilemocks.NewMockPortDiscoveryProfileBindingMapsClient(ctrl)
	mockQosSDK := profilemocks.NewMockPortQosProfileBindingMapsClient(ctrl)
	mockSecuritySDK := profilemocks.NewMockPortSecurityProfileBindingMapsClient(ctrl)

	portsWrapper := &segments.SegmentPortClientContext{
		Client:     mockPortsSDK,
		ClientType: utl.Local,
	}
	discoveryWrapper := &portprofiles.PortDiscoveryProfileBindingMapClientContext{
		Client:     mockDiscoverySDK,
		ClientType: utl.Local,
	}
	qosWrapper := &portprofiles.PortQosProfileBindingMapClientContext{
		Client:     mockQosSDK,
		ClientType: utl.Local,
	}
	securityWrapper := &portprofiles.PortSecurityProfileBindingMapClientContext{
		Client:     mockSecuritySDK,
		ClientType: utl.Local,
	}

	originalPorts := cliPortsClient
	originalDiscovery := cliPortDiscoveryProfileBindingMapsClient
	originalQos := cliPortQosProfileBindingMapsClient
	originalSecurity := cliPortSecurityProfileBindingMapsClient
	defer func() {
		cliPortsClient = originalPorts
		cliPortDiscoveryProfileBindingMapsClient = originalDiscovery
		cliPortQosProfileBindingMapsClient = originalQos
		cliPortSecurityProfileBindingMapsClient = originalSecurity
	}()

	cliPortsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentPortClientContext {
		return portsWrapper
	}
	cliPortDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortDiscoveryProfileBindingMapClientContext {
		return discoveryWrapper
	}
	cliPortQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortQosProfileBindingMapClientContext {
		return qosWrapper
	}
	cliPortSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortSecurityProfileBindingMapClientContext {
		return securityWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockPortsSDK.EXPECT().Get(segPortSegmentID, segPortPortID).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{Results: []model.PortDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{Results: []model.PortQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{Results: []model.PortSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
		assert.Equal(t, segPortDescription, d.Get("description"))
		assert.Equal(t, segPortPath, d.Get("path"))
		assert.Equal(t, int(segPortRevision), d.Get("revision"))
	})
}

func TestMockResourceNsxtPolicySegmentPortCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPortsSDK := portmocks.NewMockPortsClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	portsWrapper := &segments.SegmentPortClientContext{
		Client:     mockPortsSDK,
		ClientType: utl.Local,
	}
	mockDiscoverySDK := profilemocks.NewMockPortDiscoveryProfileBindingMapsClient(ctrl)
	mockQosSDK := profilemocks.NewMockPortQosProfileBindingMapsClient(ctrl)
	mockSecuritySDK := profilemocks.NewMockPortSecurityProfileBindingMapsClient(ctrl)

	originalPorts := cliPortsClient
	originalInfra := cliInfraClient
	originalDiscovery := cliPortDiscoveryProfileBindingMapsClient
	originalQos := cliPortQosProfileBindingMapsClient
	originalSecurity := cliPortSecurityProfileBindingMapsClient
	defer func() {
		cliPortsClient = originalPorts
		cliInfraClient = originalInfra
		cliPortDiscoveryProfileBindingMapsClient = originalDiscovery
		cliPortQosProfileBindingMapsClient = originalQos
		cliPortSecurityProfileBindingMapsClient = originalSecurity
	}()

	cliPortsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentPortClientContext {
		return portsWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}
	cliPortDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortDiscoveryProfileBindingMapClientContext {
		return &portprofiles.PortDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
	}
	cliPortQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortQosProfileBindingMapClientContext {
		return &portprofiles.PortQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
	}
	cliPortSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortSecurityProfileBindingMapClientContext {
		return &portprofiles.PortSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
	}

	t.Run("Create success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPortsSDK.EXPECT().Get(segPortSegmentID, gomock.Any()).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{Results: []model.PortDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(segPortSegmentID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{Results: []model.PortQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(segPortSegmentID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{Results: []model.PortSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"display_name": segPortDisplayName,
			"description":  segPortDescription,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicySegmentPortUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPortsSDK := portmocks.NewMockPortsClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	portsWrapper := &segments.SegmentPortClientContext{
		Client:     mockPortsSDK,
		ClientType: utl.Local,
	}
	mockDiscoverySDK := profilemocks.NewMockPortDiscoveryProfileBindingMapsClient(ctrl)
	mockQosSDK := profilemocks.NewMockPortQosProfileBindingMapsClient(ctrl)
	mockSecuritySDK := profilemocks.NewMockPortSecurityProfileBindingMapsClient(ctrl)

	originalPorts := cliPortsClient
	originalInfra := cliInfraClient
	originalDiscovery := cliPortDiscoveryProfileBindingMapsClient
	originalQos := cliPortQosProfileBindingMapsClient
	originalSecurity := cliPortSecurityProfileBindingMapsClient
	defer func() {
		cliPortsClient = originalPorts
		cliInfraClient = originalInfra
		cliPortDiscoveryProfileBindingMapsClient = originalDiscovery
		cliPortQosProfileBindingMapsClient = originalQos
		cliPortSecurityProfileBindingMapsClient = originalSecurity
	}()

	cliPortsClient = func(sessionContext utl.SessionContext, connector client.Connector) *segments.SegmentPortClientContext {
		return portsWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}
	cliPortDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortDiscoveryProfileBindingMapClientContext {
		return &portprofiles.PortDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
	}
	cliPortQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortQosProfileBindingMapClientContext {
		return &portprofiles.PortQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
	}
	cliPortSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortSecurityProfileBindingMapClientContext {
		return &portprofiles.PortSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
	}

	t.Run("Update success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPortsSDK.EXPECT().Get(segPortSegmentID, segPortPortID).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{Results: []model.PortDiscoveryProfileBindingMap{}}, nil)
		mockQosSDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{Results: []model.PortQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{Results: []model.PortSecurityProfileBindingMap{}}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"display_name": segPortDisplayName,
			"description":  segPortDescription,
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"display_name": segPortDisplayName,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment ID")
	})
}

func TestMockResourceNsxtPolicySegmentPortDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	originalInfra := cliInfraClient
	defer func() { cliInfraClient = originalInfra }()
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Delete success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"display_name": segPortDisplayName,
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Segment Port ID")
	})
}
