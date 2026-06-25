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
