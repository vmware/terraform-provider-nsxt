//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate mocks for this test, run mockgen for PortsClient and profile binding map clients
// in api/infra/segments and api/infra/segments/ports.

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

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
		assert.Equal(t, segPortPortID, d.Get("nsx_id"))
		assert.Equal(t, segPortPath, d.Get("path"))
		assert.Equal(t, int(segPortRevision), d.Get("revision"))
	})

	t.Run("Read with configured profiles reads only configured profile bindings", func(t *testing.T) {
		mockPortsSDK.EXPECT().Get(segPortSegmentID, segPortPortID).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		ipPath := "/infra/ip-discovery-profiles/p1"
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{
						IpDiscoveryProfilePath: &ipPath,
					},
				},
			}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path": ipPath,
				},
			},
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
		discovery := d.Get("discovery_profile").([]interface{})
		require.Len(t, discovery, 1)
		discoveryMap := discovery[0].(map[string]interface{})
		assert.Equal(t, ipPath, discoveryMap["ip_discovery_profile_path"])
	})

	t.Run("Read with only security profile reads only security profile binding", func(t *testing.T) {
		mockPortsSDK.EXPECT().Get(segPortSegmentID, segPortPortID).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		secPath := "/infra/segment-security-profiles/p1"
		mockSecuritySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{
				Results: []model.PortSecurityProfileBindingMap{
					{
						SegmentSecurityProfilePath: &secPath,
					},
				},
			}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"security_profile": []interface{}{
				map[string]interface{}{
					"security_profile_path": secPath,
				},
			},
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
		security := d.Get("security_profile").([]interface{})
		require.Len(t, security, 1)
		securityMap := security[0].(map[string]interface{})
		assert.Equal(t, secPath, securityMap["security_profile_path"])
	})

	t.Run("Read with multiple configured profiles reads only configured profile bindings", func(t *testing.T) {
		mockPortsSDK.EXPECT().Get(segPortSegmentID, segPortPortID).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		ipPath := "/infra/ip-discovery-profiles/p1"
		secPath := "/infra/segment-security-profiles/p1"
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{
						IpDiscoveryProfilePath: &ipPath,
					},
				},
			}, nil)
		mockSecuritySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{
				Results: []model.PortSecurityProfileBindingMap{
					{
						SegmentSecurityProfilePath: &secPath,
					},
				},
			}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path": ipPath,
				},
			},
			"security_profile": []interface{}{
				map[string]interface{}{
					"security_profile_path": secPath,
				},
			},
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
		discovery := d.Get("discovery_profile").([]interface{})
		require.Len(t, discovery, 1)
		security := d.Get("security_profile").([]interface{})
		require.Len(t, security, 1)
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

	t.Run("Create with profiles reads and populates profile bindings", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPortsSDK.EXPECT().Get(segPortSegmentID, gomock.Any()).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		ipPath := "/infra/ip-discovery-profiles/p1"
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{
						IpDiscoveryProfilePath: &ipPath,
					},
				},
			}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"display_name": segPortDisplayName,
			"description":  segPortDescription,
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path": ipPath,
				},
			},
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
		discovery := d.Get("discovery_profile").([]interface{})
		require.Len(t, discovery, 1)
	})

	t.Run("Create succeeds with short segment path", func(t *testing.T) {
		shortSegmentPath := "a/b/c"
		shortSegmentID := "c"
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPortsSDK.EXPECT().Get(shortSegmentID, gomock.Any()).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": shortSegmentPath,
			"display_name": segPortDisplayName,
		})

		m := newGoMockProviderClient()
		assert.NotPanics(t, func() {
			err := resourceNsxtPolicySegmentPortCreate(d, m)
			require.NoError(t, err)
		})
		assert.NotEmpty(t, d.Id())
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

	t.Run("Update with profiles reads profile bindings", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPortsSDK.EXPECT().Get(segPortSegmentID, segPortPortID).Return(model.SegmentPort{
			DisplayName: &segPortDisplayName,
			Description: &segPortDescription,
			Path:        &segPortPath,
			Revision:    &segPortRevision,
		}, nil)
		ipPath := "/infra/ip-discovery-profiles/p1"
		mockDiscoverySDK.EXPECT().List(segPortSegmentID, segPortPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{
						IpDiscoveryProfilePath: &ipPath,
					},
				},
			}, nil)

		res := resourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_path": segPortSegmentPath,
			"display_name": segPortDisplayName,
			"description":  segPortDescription,
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path": ipPath,
				},
			},
		})
		d.SetId(segPortPortID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentPortUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, segPortDisplayName, d.Get("display_name"))
		discovery := d.Get("discovery_profile").([]interface{})
		require.Len(t, discovery, 1)
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

func TestMockNsxtIsT1Segment(t *testing.T) {
	tests := []struct {
		name        string
		segmentPath string
		expected    bool
	}{
		{
			name:        "T1 segment path",
			segmentPath: "/infra/tier-1s/t1-1/segments/seg-1",
			expected:    true,
		},
		{
			name:        "non-T1 infra segment path",
			segmentPath: "/infra/segments/seg-1",
			expected:    false,
		},
		{
			name:        "short path with three tokens",
			segmentPath: "a/b/c",
			expected:    false,
		},
		{
			name:        "short path with two tokens",
			segmentPath: "a/b",
			expected:    false,
		},
		{
			name:        "empty string",
			segmentPath: "",
			expected:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isT1Segment(tt.segmentPath))
		})
	}
}

func TestMockNsxtGetT1IdFromSegPath(t *testing.T) {
	tests := []struct {
		name        string
		segmentPath string
		expectedID  string
	}{
		{
			name:        "T1 segment path",
			segmentPath: "/infra/tier-1s/t1-1/segments/seg-1",
			expectedID:  "t1-1",
		},
		{
			name:        "non-T1 segment path",
			segmentPath: "/infra/segments/seg-1",
			expectedID:  "",
		},
		{
			name:        "short path with three tokens",
			segmentPath: "a/b/c",
			expectedID:  "",
		},
		{
			name:        "empty string",
			segmentPath: "",
			expectedID:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, getT1IdFromSegPath(tt.segmentPath))
		})
	}
}

func TestMockNsxtSegmentPortImporter(t *testing.T) {
	res := resourceNsxtPolicySegmentPort()

	t.Run("Import with invalid ID (not policy path) returns error without panic", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("port-123")
		rd, err := getSegmentPortPathOrIDResourceImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path port-123")
		assert.Nil(t, rd)
	})

	t.Run("Import with empty ID returns error without panic", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("")
		rd, err := getSegmentPortPathOrIDResourceImporter(d, nil)
		require.Error(t, err)
		assert.Nil(t, rd)
	})

	t.Run("Import with segment path (missing /ports/) returns error without panic", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/segments/seg-1")
		rd, err := getSegmentPortPathOrIDResourceImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path /infra/segments/seg-1")
		assert.Nil(t, rd)
	})

	t.Run("Import with valid infra segment port path succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/segments/seg-1/ports/port-123")
		rd, err := getSegmentPortPathOrIDResourceImporter(d, nil)
		require.NoError(t, err)
		require.NotNil(t, rd)
		assert.Equal(t, "port-123", d.Id())
		assert.Equal(t, "/infra/segments/seg-1", d.Get("segment_path"))
	})

	t.Run("Import with valid tier-1 segment port path succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/tier-1s/t1-1/segments/seg-1/ports/port-123")
		rd, err := getSegmentPortPathOrIDResourceImporter(d, nil)
		require.NoError(t, err)
		require.NotNil(t, rd)
		assert.Equal(t, "port-123", d.Id())
		assert.Equal(t, "/infra/tier-1s/t1-1/segments/seg-1", d.Get("segment_path"))
	})

	t.Run("Import with valid multitenancy segment port path succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/proj-1/infra/segments/seg-1/ports/port-123")
		rd, err := getSegmentPortPathOrIDResourceImporter(d, nil)
		require.NoError(t, err)
		require.NotNil(t, rd)
		assert.Equal(t, "port-123", d.Id())
		assert.Equal(t, "/orgs/default/projects/proj-1/infra/segments/seg-1", d.Get("segment_path"))
	})

	t.Run("Import with provider client populates attached profile bindings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDiscoverySDK := profilemocks.NewMockPortDiscoveryProfileBindingMapsClient(ctrl)
		mockQosSDK := profilemocks.NewMockPortQosProfileBindingMapsClient(ctrl)
		mockSecuritySDK := profilemocks.NewMockPortSecurityProfileBindingMapsClient(ctrl)

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

		originalDiscovery := cliPortDiscoveryProfileBindingMapsClient
		originalQos := cliPortQosProfileBindingMapsClient
		originalSecurity := cliPortSecurityProfileBindingMapsClient
		defer func() {
			cliPortDiscoveryProfileBindingMapsClient = originalDiscovery
			cliPortQosProfileBindingMapsClient = originalQos
			cliPortSecurityProfileBindingMapsClient = originalSecurity
		}()

		cliPortDiscoveryProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortDiscoveryProfileBindingMapClientContext {
			return discoveryWrapper
		}
		cliPortQosProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortQosProfileBindingMapClientContext {
			return qosWrapper
		}
		cliPortSecurityProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *portprofiles.PortSecurityProfileBindingMapClientContext {
			return securityWrapper
		}

		ipPath := "/infra/ip-discovery-profiles/p1"
		mockDiscoverySDK.EXPECT().List("seg-1", "port-123", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{
						IpDiscoveryProfilePath: &ipPath,
					},
				},
			}, nil)
		mockQosSDK.EXPECT().List("seg-1", "port-123", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{Results: []model.PortQosProfileBindingMap{}}, nil)
		mockSecuritySDK.EXPECT().List("seg-1", "port-123", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{Results: []model.PortSecurityProfileBindingMap{}}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/segments/seg-1/ports/port-123")
		m := newGoMockProviderClient()
		rd, err := getSegmentPortPathOrIDResourceImporter(d, m)
		require.NoError(t, err)
		require.NotNil(t, rd)
		assert.Equal(t, "port-123", d.Id())
		assert.Equal(t, "/infra/segments/seg-1", d.Get("segment_path"))

		discovery := d.Get("discovery_profile").([]interface{})
		require.Len(t, discovery, 1)
		discoveryMap := discovery[0].(map[string]interface{})
		assert.Equal(t, ipPath, discoveryMap["ip_discovery_profile_path"])
	})
}
