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

	groupsapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains/groups"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	groupsMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/groups"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	ffppbID          = "ffppb-1"
	ffppbName        = "ffppb-fooname"
	ffppbRevision    = int64(1)
	ffppbPath        = "/infra/domains/default/groups/group-1/firewall-flood-protection-profile-binding-maps/ffppb-1"
	ffppbProfilePath = "/infra/flood-protection-profiles/dfpp-1"
	ffppbGroupPath   = "/infra/domains/default/groups/group-1"
	ffppbDomainID    = "default"
	ffppbGroupID     = "group-1"
	ffppbSeqNum      = int64(10)
)

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileBindingCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBindingSDK := groupsMocks.NewMockFirewallFloodProtectionProfileBindingMapsClient(ctrl)
	bindingWrapper := &groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext{
		Client:     mockBindingSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFirewallFloodProtectionProfileBindingMapsClient
	defer func() { cliFirewallFloodProtectionProfileBindingMapsClient = originalCli }()
	cliFirewallFloodProtectionProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext {
		return bindingWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfileBinding()

	t.Run("Create_success", func(t *testing.T) {
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, gomock.Any()).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{}, vapiErrors.NotFound{})
		mockBindingSDK.EXPECT().Patch(ffppbDomainID, ffppbGroupID, gomock.Any(), gomock.Any()).Return(nil)
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, gomock.Any()).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{
			DisplayName:    &ffppbName,
			Path:           &ffppbPath,
			Revision:       &ffppbRevision,
			ProfilePath:    &ffppbProfilePath,
			SequenceNumber: &ffppbSeqNum,
			Description:    strPtr(""),
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":    ffppbName,
			"profile_path":    ffppbProfilePath,
			"group_path":      ffppbGroupPath,
			"sequence_number": int(ffppbSeqNum),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, ffppbName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_resource_already_exists", func(t *testing.T) {
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, gomock.Any()).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{
			DisplayName: &ffppbName,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"profile_path":    ffppbProfilePath,
			"group_path":      ffppbGroupPath,
			"sequence_number": int(ffppbSeqNum),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, gomock.Any()).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{}, vapiErrors.NotFound{})
		mockBindingSDK.EXPECT().Patch(ffppbDomainID, ffppbGroupID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"profile_path":    ffppbProfilePath,
			"group_path":      ffppbGroupPath,
			"sequence_number": int(ffppbSeqNum),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileBindingRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBindingSDK := groupsMocks.NewMockFirewallFloodProtectionProfileBindingMapsClient(ctrl)
	bindingWrapper := &groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext{
		Client:     mockBindingSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFirewallFloodProtectionProfileBindingMapsClient
	defer func() { cliFirewallFloodProtectionProfileBindingMapsClient = originalCli }()
	cliFirewallFloodProtectionProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext {
		return bindingWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfileBinding()

	t.Run("Read_success", func(t *testing.T) {
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, ffppbID).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{
			DisplayName:    &ffppbName,
			Path:           &ffppbPath,
			Revision:       &ffppbRevision,
			ProfilePath:    &ffppbProfilePath,
			SequenceNumber: &ffppbSeqNum,
			Description:    strPtr(""),
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		d.SetId(ffppbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, ffppbName, d.Get("display_name"))
		assert.Equal(t, ffppbProfilePath, d.Get("profile_path"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FloodProtectionProfile ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, ffppbID).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		d.SetId(ffppbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileBindingUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBindingSDK := groupsMocks.NewMockFirewallFloodProtectionProfileBindingMapsClient(ctrl)
	bindingWrapper := &groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext{
		Client:     mockBindingSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFirewallFloodProtectionProfileBindingMapsClient
	defer func() { cliFirewallFloodProtectionProfileBindingMapsClient = originalCli }()
	cliFirewallFloodProtectionProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext {
		return bindingWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfileBinding()

	t.Run("Update_success", func(t *testing.T) {
		mockBindingSDK.EXPECT().Patch(ffppbDomainID, ffppbGroupID, ffppbID, gomock.Any()).Return(nil)
		mockBindingSDK.EXPECT().Get(ffppbDomainID, ffppbGroupID, ffppbID).Return(model.PolicyFirewallFloodProtectionProfileBindingMap{
			DisplayName:    &ffppbName,
			Path:           &ffppbPath,
			Revision:       &ffppbRevision,
			ProfilePath:    &ffppbProfilePath,
			SequenceNumber: &ffppbSeqNum,
			Description:    strPtr(""),
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"profile_path":    ffppbProfilePath,
			"group_path":      ffppbGroupPath,
			"sequence_number": int(ffppbSeqNum),
		})
		d.SetId(ffppbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "DistributedFloodProtectionProfileBinding ID")
	})
}

func TestMockResourceNsxtPolicyDistributedFloodProtectionProfileBindingDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBindingSDK := groupsMocks.NewMockFirewallFloodProtectionProfileBindingMapsClient(ctrl)
	bindingWrapper := &groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext{
		Client:     mockBindingSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFirewallFloodProtectionProfileBindingMapsClient
	defer func() { cliFirewallFloodProtectionProfileBindingMapsClient = originalCli }()
	cliFirewallFloodProtectionProfileBindingMapsClient = func(sessionContext utl.SessionContext, connector client.Connector) *groupsapi.PolicyFirewallFloodProtectionProfileBindingMapClientContext {
		return bindingWrapper
	}

	res := resourceNsxtPolicyDistributedFloodProtectionProfileBinding()

	t.Run("Delete_success", func(t *testing.T) {
		mockBindingSDK.EXPECT().Delete(ffppbDomainID, ffppbGroupID, ffppbID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		d.SetId(ffppbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "DistributedFloodProtectionProfileBinding ID")
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockBindingSDK.EXPECT().Delete(ffppbDomainID, ffppbGroupID, ffppbID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"group_path": ffppbGroupPath,
		})
		d.SetId(ffppbID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedFloodProtectionProfileBindingDelete(d, m)
		require.Error(t, err)
	})
}
