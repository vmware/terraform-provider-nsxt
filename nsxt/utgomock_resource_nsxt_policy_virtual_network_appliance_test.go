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
	vapiData "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vnaID       = "vna-1"
	vnaName     = "vna-fooname"
	vnaRevision = int64(1)
	vnaPath     = "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1/virtual-network-appliances/vna-1"
	vnaHostname = "vna-host.example.com"
	// vnaClusterPath uses the vnaClusterPath declared in the cluster test file:
	// /infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1
)

func vnaStructValue(vna model.VirtualNetworkAppliance) *vapiData.StructValue {
	vna.ResourceType = model.VirtualNetworkAppliance__TYPE_IDENTIFIER
	sv, errs := vna.GetDataValue__()
	if errs != nil {
		panic(errs[0])
	}
	return sv.(*vapiData.StructValue)
}

func setupVNACRUDMock(ctrl *gomock.Controller) *epmocks.MockVirtualNetworkAppliancesInClusterClient {
	return epmocks.NewMockVirtualNetworkAppliancesInClusterClient(ctrl)
}

func setupVNACRUDClientOverride(mock *epmocks.MockVirtualNetworkAppliancesInClusterClient) func() {
	wrapper := &enforcementpoints.VirtualNetworkApplianceCRUDClientContext{
		Client:     mock,
		ClientType: utl.Local,
	}
	orig := cliVNACRUDClient
	cliVNACRUDClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.VirtualNetworkApplianceCRUDClientContext {
		return wrapper
	}
	return func() { cliVNACRUDClient = orig }
}

func TestUnitNsxt_getVNAClusterPathComponents(t *testing.T) {
	t.Run("parses a valid cluster path", func(t *testing.T) {
		site, ep, cluster, err := getVNAClusterPathComponents(vnaClusterPath)
		require.NoError(t, err)
		assert.Equal(t, vnaClusterSiteID, site)
		assert.Equal(t, vnaClusterEPID, ep)
		assert.Equal(t, vnaClusterID, cluster)
	})

	t.Run("errors when the site segment is missing", func(t *testing.T) {
		_, _, _, err := getVNAClusterPathComponents("/enforcement-points/default/virtual-network-appliance-clusters/c1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site ID")
	})

	t.Run("errors when the enforcement-point segment is missing", func(t *testing.T) {
		_, _, _, err := getVNAClusterPathComponents("/infra/sites/default/virtual-network-appliance-clusters/c1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "enforcement-point ID")
	})

	t.Run("errors when the cluster segment is missing", func(t *testing.T) {
		_, _, _, err := getVNAClusterPathComponents("/infra/sites/default/enforcement-points/default")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cluster ID")
	})
}

func TestUnitNsxt_getVNACredentialsFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, getVNACredentialsFromSchema(nil))
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		assert.Nil(t, getVNACredentialsFromSchema([]interface{}{}))
	})

	t.Run("parses all password fields", func(t *testing.T) {
		obj := getVNACredentialsFromSchema([]interface{}{
			map[string]interface{}{
				"cli_password":   "cli-pass",
				"root_password":  "root-pass",
				"audit_password": "audit-pass",
			},
		})
		require.NotNil(t, obj)
		assert.Equal(t, "cli-pass", *obj.CliPassword)
		assert.Equal(t, "root-pass", *obj.RootPassword)
		assert.Equal(t, "audit-pass", *obj.AuditPassword)
	})

	t.Run("empty password strings are omitted", func(t *testing.T) {
		obj := getVNACredentialsFromSchema([]interface{}{
			map[string]interface{}{
				"cli_password":   "",
				"root_password":  "",
				"audit_password": "",
			},
		})
		require.NotNil(t, obj)
		assert.Nil(t, obj.CliPassword)
		assert.Nil(t, obj.RootPassword)
		assert.Nil(t, obj.AuditPassword)
	})
}

func TestUnitNsxt_getVNAManagementInterfaceFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		obj, err := getVNAManagementInterfaceFromSchema(nil)
		require.NoError(t, err)
		assert.Nil(t, obj)
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		obj, err := getVNAManagementInterfaceFromSchema([]interface{}{})
		require.NoError(t, err)
		assert.Nil(t, obj)
	})

	t.Run("parses network_id and ip_assignment", func(t *testing.T) {
		obj, err := getVNAManagementInterfaceFromSchema([]interface{}{
			map[string]interface{}{
				"network_id": "dvpg-1",
				"ip_assignment": []interface{}{
					map[string]interface{}{
						"dhcp_v4": true,
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, obj)
		assert.Equal(t, "dvpg-1", *obj.NetworkId)
		require.Len(t, obj.IpAssignmentSpecs, 1)
	})
}

func TestUnitNsxt_setVNAManagementInterfaceInSchema(t *testing.T) {
	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("nil object is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"cluster_path": vnaClusterPath})
		require.NoError(t, setVNAManagementInterfaceInSchema(d, nil))
		assert.Empty(t, d.Get("management_interface").([]interface{}))
	})

	t.Run("sets network_id and ip_assignment", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"cluster_path": vnaClusterPath})
		networkID := "dvpg-2"
		sv, err := vAPIConversion(model.Dhcpv4{IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_DHCPV4}, model.Dhcpv4BindingType())
		require.NoError(t, err)
		err = setVNAManagementInterfaceInSchema(d, &model.VirtualNetworkApplianceManagementInterface{
			NetworkId:         &networkID,
			IpAssignmentSpecs: []*vapiData.StructValue{sv},
		})
		require.NoError(t, err)
		list := d.Get("management_interface").([]interface{})
		require.Len(t, list, 1)
		assert.Equal(t, networkID, list[0].(map[string]interface{})["network_id"])
	})
}

func TestUnitNsxt_getVNADeploymentConfigFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, getVNADeploymentConfigFromSchema(nil))
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		assert.Nil(t, getVNADeploymentConfigFromSchema([]interface{}{}))
	})

	t.Run("parses fields plus reservation_info", func(t *testing.T) {
		obj := getVNADeploymentConfigFromSchema([]interface{}{
			map[string]interface{}{
				"compute_manager_id":          "cm-1",
				"cluster_or_resource_pool_id": "crp-1",
				"datastore_id":                "ds-1",
				"reservation_info": []interface{}{
					map[string]interface{}{
						"cpu_reservation_in_mhz":        1000,
						"cpu_reservation_in_shares":     model.CPUReservation_RESERVATION_IN_SHARES_HIGH_PRIORITY,
						"memory_reservation_percentage": 50,
					},
				},
			},
		})
		require.NotNil(t, obj)
		assert.Equal(t, "cm-1", *obj.ComputeManagerId)
		assert.Equal(t, "crp-1", *obj.ClusterOrResourcePoolId)
		assert.Equal(t, "ds-1", *obj.DatastoreId)
		require.NotNil(t, obj.ReservationInfo)
		assert.Equal(t, int64(1000), *obj.ReservationInfo.CpuReservation.ReservationInMhz)
		assert.Equal(t, int64(50), *obj.ReservationInfo.MemoryReservation.ReservationPercentage)
	})

	t.Run("omits fields left empty", func(t *testing.T) {
		obj := getVNADeploymentConfigFromSchema([]interface{}{map[string]interface{}{
			"compute_manager_id":          "",
			"cluster_or_resource_pool_id": "",
			"datastore_id":                "",
		}})
		require.NotNil(t, obj)
		assert.Nil(t, obj.ComputeManagerId)
		assert.Nil(t, obj.ClusterOrResourcePoolId)
		assert.Nil(t, obj.DatastoreId)
		assert.Nil(t, obj.ReservationInfo)
	})
}

func TestUnitNsxt_setVNADeploymentConfigInSchema(t *testing.T) {
	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("nil object is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"cluster_path": vnaClusterPath})
		require.NoError(t, setVNADeploymentConfigInSchema(d, nil))
		assert.Empty(t, d.Get("vm_deployment_config").([]interface{}))
	})

	t.Run("sets all fields plus reservation_info", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"cluster_path": vnaClusterPath})
		cm := "cm-2"
		crp := "crp-2"
		ds := "ds-2"
		mhz := int64(2000)
		shares := model.CPUReservation_RESERVATION_IN_SHARES_LOW_PRIORITY
		pct := int64(75)
		err := setVNADeploymentConfigInSchema(d, &model.VirtualNetworkApplianceDeploymentConfig{
			ComputeManagerId:        &cm,
			ClusterOrResourcePoolId: &crp,
			DatastoreId:             &ds,
			ReservationInfo: &model.ReservationInfo{
				CpuReservation: &model.CPUReservation{
					ReservationInMhz:    &mhz,
					ReservationInShares: &shares,
				},
				MemoryReservation: &model.MemoryReservation{
					ReservationPercentage: &pct,
				},
			},
		})
		require.NoError(t, err)
		list := d.Get("vm_deployment_config").([]interface{})
		require.Len(t, list, 1)
		elem := list[0].(map[string]interface{})
		assert.Equal(t, cm, elem["compute_manager_id"])
		ri := elem["reservation_info"].([]interface{})
		require.Len(t, ri, 1)
		assert.Equal(t, 2000, ri[0].(map[string]interface{})["cpu_reservation_in_mhz"])
	})
}

func TestUnitNsxt_resourceNsxtPolicyVirtualNetworkApplianceImporter(t *testing.T) {
	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("valid import ID sets cluster_path and seeds credentials", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vnaPath)

		out, err := resourceNsxtPolicyVirtualNetworkApplianceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, vnaClusterPath, out[0].Get("cluster_path"))
		creds := out[0].Get("credentials").([]interface{})
		require.Len(t, creds, 1)
	})

	t.Run("malformed ID errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("not-a-valid-path")

		_, err := resourceNsxtPolicyVirtualNetworkApplianceImporter(d, nil)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceCreate(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNA := setupVNACRUDMock(ctrl)
	restore := setupVNACRUDClientOverride(mockVNA)
	defer restore()

	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("Create_success", func(t *testing.T) {
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
			Hostname:    &vnaHostname,
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any()).Return(nil, vapiErrors.NotFound{})
		mockVNA.EXPECT().Patch(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any(), gomock.Any()).Return(nil)
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any()).Return(returnSV, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": vnaName,
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, vnaName, d.Get("display_name"))
		assert.Equal(t, vnaHostname, d.Get("hostname"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		existingSV := vnaStructValue(model.VirtualNetworkAppliance{Id: &vnaID})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(existingSV, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":       vnaID,
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_cluster_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": "/infra/sites/default/enforcement-points/default/no-clusters/x",
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cluster ID")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any()).Return(nil, vapiErrors.NotFound{})
		mockVNA.EXPECT().Patch(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceRead(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNA := setupVNACRUDMock(ctrl)
	restore := setupVNACRUDClientOverride(mockVNA)
	defer restore()

	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("Read_success", func(t *testing.T) {
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vnaName, d.Get("display_name"))
		assert.Equal(t, vnaPath, d.Get("path"))
	})

	t.Run("Read_not_found_clears_id", func(t *testing.T) {
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(nil, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	// Verify credentials are preserved after Read when the API returns nil
	// credentials (passwords are write-only and never included in GET).
	t.Run("Read_preserves_credentials_when_api_returns_nil", func(t *testing.T) {
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
			// Credentials intentionally absent: API never returns passwords.
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		cliPass := "TestCli@Secret99"
		rootPass := "TestRoot@Secret99"
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
			"credentials": []interface{}{
				map[string]interface{}{
					"cli_password":   cliPass,
					"root_password":  rootPass,
					"audit_password": "",
					"cli_username":   "",
					"audit_username": "",
				},
			},
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)

		creds := d.Get("credentials").([]interface{})
		require.Len(t, creds, 1, "credentials block must be preserved in state")
		credsMap := creds[0].(map[string]interface{})
		assert.Equal(t, cliPass, credsMap["cli_password"], "cli_password must be preserved")
		assert.Equal(t, rootPass, credsMap["root_password"], "root_password must be preserved")
	})

	// Verify credentials (including passwords) are preserved and computed
	// usernames are updated when the API returns a Credentials object.
	t.Run("Read_preserves_passwords_and_updates_usernames", func(t *testing.T) {
		cliUsername := "admin"
		auditUsername := "audit"
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
			Credentials: &model.VirtualNetworkApplianceCredential{
				CliUsername:   &cliUsername,
				AuditUsername: &auditUsername,
			},
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		cliPass := "TestCli@Secret99"
		rootPass := "TestRoot@Secret99"
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
			"credentials": []interface{}{
				map[string]interface{}{
					"cli_password":   cliPass,
					"root_password":  rootPass,
					"audit_password": "",
					"cli_username":   "",
					"audit_username": "",
				},
			},
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)

		creds := d.Get("credentials").([]interface{})
		require.Len(t, creds, 1, "credentials block must be preserved in state")
		credsMap := creds[0].(map[string]interface{})
		assert.Equal(t, cliPass, credsMap["cli_password"], "cli_password must be preserved")
		assert.Equal(t, rootPass, credsMap["root_password"], "root_password must be preserved")
		assert.Equal(t, cliUsername, credsMap["cli_username"], "cli_username must be set from API response")
		assert.Equal(t, auditUsername, credsMap["audit_username"], "audit_username must be set from API response")
	})

	// Reproduce the acceptance-test panic (importBasic): the Terraform Plugin
	// SDK v2 normalises an empty TypeList element to nil when read back via
	// d.Get, so c[0].(map[string]interface{}) panics. This test forces that
	// nil-element path by calling d.Set("credentials", []interface{}{nil})
	// after creating the ResourceData, then verifies that Read neither panics
	// nor returns an error and that the API-returned usernames are written.
	t.Run("Read_with_nil_credentials_element_does_not_panic", func(t *testing.T) {
		cliUser := "admin"
		auditUser := "audit"
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
			Credentials: &model.VirtualNetworkApplianceCredential{
				CliUsername:   &cliUser,
				AuditUsername: &auditUser,
			},
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		// Simulate the SDK normalisation: empty map → nil element.
		require.NoError(t, d.Set("credentials", []interface{}{nil}))

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)

		credList := d.Get("credentials").([]interface{})
		require.Len(t, credList, 1, "credentials block must be written when list had a nil element")
		creds := credList[0].(map[string]interface{})
		assert.Equal(t, cliUser, creds["cli_username"], "cli_username must reflect NSX response")
		assert.Equal(t, auditUser, creds["audit_username"], "audit_username must reflect NSX response")
	})

	// Verify that no credentials block is written to state when the API
	// returns a Credentials object but no credentials are configured locally.
	t.Run("Read_no_credentials_block_when_not_configured", func(t *testing.T) {
		cliUsername := "admin"
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
			Credentials: &model.VirtualNetworkApplianceCredential{
				CliUsername: &cliUsername,
			},
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)

		creds := d.Get("credentials").([]interface{})
		assert.Empty(t, creds, "credentials block must not appear when not configured")
	})

	// Simulate the import path: the Importer seeds an empty credentials block
	// before Read is called. Read must populate the block with NSX-returned
	// usernames so that suppressIfEmptyPriorState can suppress the subsequent
	// password diff and the plan shows zero drift (bug 3715433).
	t.Run("Read_with_importer_seeded_block_writes_usernames_to_state", func(t *testing.T) {
		cliUser := "admin"
		auditUser := "audit"
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
			Credentials: &model.VirtualNetworkApplianceCredential{
				CliUsername:   &cliUser,
				AuditUsername: &auditUser,
			},
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		// Seed an empty credentials block — exactly what the Importer does.
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
			"credentials":  []interface{}{map[string]interface{}{}},
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)

		credList := d.Get("credentials").([]interface{})
		require.Len(t, credList, 1, "credentials block must be written when importer-seeded block is present")
		creds := credList[0].(map[string]interface{})
		assert.Equal(t, cliUser, creds["cli_username"], "cli_username must reflect NSX response")
		assert.Equal(t, auditUser, creds["audit_username"], "audit_username must reflect NSX response")
		assert.Equal(t, "", creds["cli_password"], "cli_password must remain empty (write-only)")
		assert.Equal(t, "", creds["root_password"], "root_password must remain empty (write-only)")
	})
}

func TestSuppressIfEmptyPriorState(t *testing.T) {
	res := resourceNsxtPolicyVirtualNetworkAppliance()

	// Build a ResourceData that simulates an existing (imported) resource:
	// no credentials in state yet, resource ID is set.
	existing := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_path": vnaClusterPath,
	})
	existing.SetId(vnaID)

	// Build a ResourceData that simulates a new resource (no ID yet).
	fresh := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cluster_path": vnaClusterPath,
	})
	// fresh has no ID set.

	// Existing resource, old state is empty: suppress (import scenario).
	assert.True(t, suppressIfEmptyPriorState("cli_password", "", "VMware123!", existing),
		"must suppress diff when old (state) is empty and resource exists (import)")
	assert.True(t, suppressIfEmptyPriorState("root_password", "", "VMware123!", existing),
		"must suppress diff when old (state) is empty and resource exists (import)")

	// New resource (no ID), old state is empty: do not suppress so passwords
	// are included in the Create diff.
	assert.False(t, suppressIfEmptyPriorState("cli_password", "", "VMware123!", fresh),
		"must not suppress diff for a new resource (no ID)")

	// Non-empty old value: never suppress so password changes are applied.
	assert.False(t, suppressIfEmptyPriorState("cli_password", "OldPass!", "NewPass!", existing),
		"must not suppress diff when old (state) is non-empty")
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceUpdate(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNA := setupVNACRUDMock(ctrl)
	restore := setupVNACRUDClientOverride(mockVNA)
	defer restore()

	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("Update_success", func(t *testing.T) {
		updatedName := vnaName + "-updated"
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &updatedName,
			Path:        &vnaPath,
			Revision:    &vnaRevision,
		})
		mockVNA.EXPECT().Update(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID, gomock.Any()).Return(returnSV, nil)
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": updatedName,
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, updatedName, d.Get("display_name"))
	})

	t.Run("Update_fails_on_API_error", func(t *testing.T) {
		mockVNA.EXPECT().Update(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID, gomock.Any()).Return(nil, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceDelete(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNA := setupVNACRUDMock(ctrl)
	restore := setupVNACRUDClientOverride(mockVNA)
	defer restore()

	res := resourceNsxtPolicyVirtualNetworkAppliance()

	t.Run("Delete_success", func(t *testing.T) {
		mockVNA.EXPECT().Delete(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID, (*bool)(nil)).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_on_API_error", func(t *testing.T) {
		mockVNA.EXPECT().Delete(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID, (*bool)(nil)).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		d.SetId(vnaID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceDelete(d, m)
		require.Error(t, err)
	})
}
