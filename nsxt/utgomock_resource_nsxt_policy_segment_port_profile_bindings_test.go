//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
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
)

func minimalSegmentPortProfileBindingsData() map[string]interface{} {
	return map[string]interface{}{
		"segment_port_path": "invalid-port-path",
	}
}

var (
	segPortBindingsSegmentPath = "/infra/segments/seg-1"
	segPortBindingsSegmentID   = "seg-1"
	segPortBindingsPortID      = "port-1"
	segPortBindingsPortPath    = "/infra/segments/seg-1/ports/port-1"
)

type segmentPortProfileBindingsMocks struct {
	ports      *portmocks.MockPortsClient
	discovery  *profilemocks.MockPortDiscoveryProfileBindingMapsClient
	qos        *profilemocks.MockPortQosProfileBindingMapsClient
	security   *profilemocks.MockPortSecurityProfileBindingMapsClient
	infra      *inframocks.MockInfraClient
	restoreAll func()
}

// setupSegmentPortProfileBindingsMocks wires the mock seams used by getSegmentPort,
// nsxtPolicySegmentPortProfilesRead, and policyInfraPatch for the non-T1 (standalone
// segment) case exercised by segPortBindingsPortPath.
func setupSegmentPortProfileBindingsMocks(t *testing.T, ctrl *gomock.Controller) *segmentPortProfileBindingsMocks {
	mockPortsSDK := portmocks.NewMockPortsClient(ctrl)
	mockDiscoverySDK := profilemocks.NewMockPortDiscoveryProfileBindingMapsClient(ctrl)
	mockQosSDK := profilemocks.NewMockPortQosProfileBindingMapsClient(ctrl)
	mockSecuritySDK := profilemocks.NewMockPortSecurityProfileBindingMapsClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	portsWrapper := &segments.SegmentPortClientContext{Client: mockPortsSDK, ClientType: utl.Local}
	discoveryWrapper := &portprofiles.PortDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
	qosWrapper := &portprofiles.PortQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
	securityWrapper := &portprofiles.PortSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}

	originalPorts := cliPortsClient
	originalDiscovery := cliPortDiscoveryProfileBindingMapsClient
	originalQos := cliPortQosProfileBindingMapsClient
	originalSecurity := cliPortSecurityProfileBindingMapsClient
	originalInfra := cliInfraClient

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
	cliInfraClient = func(_ utl.SessionContext, _ client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	restore := func() {
		cliPortsClient = originalPorts
		cliPortDiscoveryProfileBindingMapsClient = originalDiscovery
		cliPortQosProfileBindingMapsClient = originalQos
		cliPortSecurityProfileBindingMapsClient = originalSecurity
		cliInfraClient = originalInfra
	}
	t.Cleanup(restore)

	return &segmentPortProfileBindingsMocks{
		ports:      mockPortsSDK,
		discovery:  mockDiscoverySDK,
		qos:        mockQosSDK,
		security:   mockSecuritySDK,
		infra:      mockInfraSDK,
		restoreAll: restore,
	}
}

func expectEmptyProfileBindingLists(m *segmentPortProfileBindingsMocks) {
	m.discovery.EXPECT().
		List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.PortDiscoveryProfileBindingMapListResult{}, nil)
	m.qos.EXPECT().
		List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.PortQosProfileBindingMapListResult{}, nil)
	m.security.EXPECT().
		List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.PortSecurityProfileBindingMapListResult{}, nil)
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsCreate(t *testing.T) {
	t.Run("Create fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())

		err := resourceNsxtPolicySegmentPortProfileBindingsCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsRead(t *testing.T) {
	t.Run("Read fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())
		d.SetId("port-1")

		err := resourceNsxtPolicySegmentPortProfileBindingsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsUpdate(t *testing.T) {
	t.Run("Update fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())
		d.SetId("port-1")

		err := resourceNsxtPolicySegmentPortProfileBindingsUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsDelete(t *testing.T) {
	t.Run("Delete fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())
		d.SetId("port-1")

		err := resourceNsxtPolicySegmentPortProfileBindingsDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupSegmentPortProfileBindingsMocks(t, ctrl)

	t.Run("Create success", func(t *testing.T) {
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, nil)
		m.infra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, nil)
		expectEmptyProfileBindingLists(m)

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})

		err := resourceNsxtPolicySegmentPortProfileBindingsCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, segPortBindingsPortID, d.Id())
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsReadSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupSegmentPortProfileBindingsMocks(t, ctrl)

	t.Run("Read success with no bindings configured", func(t *testing.T) {
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, nil)
		expectEmptyProfileBindingLists(m)

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read success populates discovery, qos and security profile bindings", func(t *testing.T) {
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, nil)

		ipPath := "/infra/ip-discovery-profiles/ip-1"
		macPath := "/infra/mac-discovery-profiles/mac-1"
		discoveryBindingPath := "/infra/segments/seg-1/ports/port-1/port-discovery-profile-binding-maps/dbm-1"
		discoveryRevision := int64(1)
		m.discovery.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{
						IpDiscoveryProfilePath:  &ipPath,
						MacDiscoveryProfilePath: &macPath,
						Path:                    &discoveryBindingPath,
						Revision:                &discoveryRevision,
					},
				},
			}, nil)

		qosPath := "/infra/qos-profiles/qos-1"
		qosBindingPath := "/infra/segments/seg-1/ports/port-1/port-qos-profile-binding-maps/qbm-1"
		qosRevision := int64(2)
		m.qos.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{
				Results: []model.PortQosProfileBindingMap{
					{
						QosProfilePath: &qosPath,
						Path:           &qosBindingPath,
						Revision:       &qosRevision,
					},
				},
			}, nil)

		securityPath := "/infra/segment-security-profiles/sec-1"
		spoofguardPath := "/infra/spoofguard-profiles/sg-1"
		securityBindingPath := "/infra/segments/seg-1/ports/port-1/port-security-profile-binding-maps/sbm-1"
		securityRevision := int64(3)
		m.security.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{
				Results: []model.PortSecurityProfileBindingMap{
					{
						SegmentSecurityProfilePath: &securityPath,
						SpoofguardProfilePath:      &spoofguardPath,
						Path:                       &securityBindingPath,
						Revision:                   &securityRevision,
					},
				},
			}, nil)

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		discovery := d.Get("discovery_profile").([]interface{})
		require.Len(t, discovery, 1)
		discoveryMap := discovery[0].(map[string]interface{})
		assert.Equal(t, ipPath, discoveryMap["ip_discovery_profile_path"])
		assert.Equal(t, macPath, discoveryMap["mac_discovery_profile_path"])
		assert.Equal(t, discoveryBindingPath, discoveryMap["binding_map_path"])

		qos := d.Get("qos_profile").([]interface{})
		require.Len(t, qos, 1)
		qosMap := qos[0].(map[string]interface{})
		assert.Equal(t, qosPath, qosMap["qos_profile_path"])
		assert.Equal(t, qosBindingPath, qosMap["binding_map_path"])

		security := d.Get("security_profile").([]interface{})
		require.Len(t, security, 1)
		securityMap := security[0].(map[string]interface{})
		assert.Equal(t, securityPath, securityMap["security_profile_path"])
		assert.Equal(t, spoofguardPath, securityMap["spoofguard_profile_path"])
		assert.Equal(t, securityBindingPath, securityMap["binding_map_path"])
	})

	t.Run("Read clears ID when segment port not found", func(t *testing.T) {
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsUpdateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupSegmentPortProfileBindingsMocks(t, ctrl)

	t.Run("Update success", func(t *testing.T) {
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, nil)
		m.infra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		m.ports.EXPECT().Get(segPortBindingsSegmentID, segPortBindingsPortID).Return(model.SegmentPort{}, nil)
		expectEmptyProfileBindingLists(m)

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsDeleteSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupSegmentPortProfileBindingsMocks(t, ctrl)

	t.Run("Delete with no bindings configured skips all delete calls", func(t *testing.T) {
		expectEmptyProfileBindingLists(m)

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete removes all configured bindings", func(t *testing.T) {
		securityBindingPath := "/infra/segments/seg-1/ports/port-1/port-security-profile-binding-maps/sbm-1"
		discoveryBindingPath := "/infra/segments/seg-1/ports/port-1/port-discovery-profile-binding-maps/dbm-1"
		qosBindingPath := "/infra/segments/seg-1/ports/port-1/port-qos-profile-binding-maps/qbm-1"
		securityPath := "/infra/segment-security-profiles/sec-1"
		spoofguardPath := "/infra/spoofguard-profiles/sg-1"
		ipPath := "/infra/ip-discovery-profiles/ip-1"
		macPath := "/infra/mac-discovery-profiles/mac-1"
		qosPath := "/infra/qos-profiles/qos-1"

		m.discovery.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{
				Results: []model.PortDiscoveryProfileBindingMap{
					{IpDiscoveryProfilePath: &ipPath, MacDiscoveryProfilePath: &macPath, Path: &discoveryBindingPath},
				},
			}, nil)
		m.qos.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{
				Results: []model.PortQosProfileBindingMap{
					{QosProfilePath: &qosPath, Path: &qosBindingPath},
				},
			}, nil)
		m.security.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{
				Results: []model.PortSecurityProfileBindingMap{
					{SegmentSecurityProfilePath: &securityPath, SpoofguardProfilePath: &spoofguardPath, Path: &securityBindingPath},
				},
			}, nil)

		m.security.EXPECT().Delete(segPortBindingsSegmentID, segPortBindingsPortID, "sbm-1").Return(nil)
		m.discovery.EXPECT().Delete(segPortBindingsSegmentID, segPortBindingsPortID, "dbm-1").Return(nil)
		m.qos.EXPECT().Delete(segPortBindingsSegmentID, segPortBindingsPortID, "qbm-1").Return(nil)

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when security profile binding delete errors", func(t *testing.T) {
		securityBindingPath := "/infra/segments/seg-1/ports/port-1/port-security-profile-binding-maps/sbm-1"
		securityPath := "/infra/segment-security-profiles/sec-1"

		m.discovery.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortDiscoveryProfileBindingMapListResult{}, nil)
		m.qos.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortQosProfileBindingMapListResult{}, nil)
		m.security.EXPECT().
			List(segPortBindingsSegmentID, segPortBindingsPortID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PortSecurityProfileBindingMapListResult{
				Results: []model.PortSecurityProfileBindingMap{
					{SegmentSecurityProfilePath: &securityPath, Path: &securityBindingPath},
				},
			}, nil)
		m.security.EXPECT().Delete(segPortBindingsSegmentID, segPortBindingsPortID, "sbm-1").Return(errors.New("API error"))

		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"segment_port_path": segPortBindingsPortPath,
		})
		d.SetId(segPortBindingsPortID)

		err := resourceNsxtPolicySegmentPortProfileBindingsDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error deleting the security profile")
	})
}
