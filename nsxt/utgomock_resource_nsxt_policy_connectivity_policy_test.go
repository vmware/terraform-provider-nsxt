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

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	connmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	connPolicyID         = "conn-policy-1"
	connPolicyParentPath = "/orgs/default/projects/project-1/transit-gateways/tgw-1"
	connPolicyOrgID      = "default"
	connPolicyProjectID  = "project-1"
	connPolicyTGWID      = "tgw-1"
	connPolicyGroupPath  = "/orgs/default/projects/project-1/infra/domains/domain1/groups/group1"
	connPolicyName       = "conn-policy-fooname"
	connPolicyRevision   = int64(1)
	connPolicyPath       = "/orgs/default/projects/project-1/transit-gateways/tgw-1/connectivity-policies/conn-policy-1"
	connPolicyScope      = model.ConnectivityPolicy_CONNECTIVITY_SCOPE_COMMUNITY
)

func TestMockResourceNsxtPolicyConnectivityPolicyCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnSDK := connmocks.NewMockConnectivityPoliciesClient(ctrl)
	connWrapper := &transitgateways.ConnectivityPolicyClientContext{
		Client:     mockConnSDK,
		ClientType: utl.Multitenancy,
	}

	originalCli := cliConnectivityPoliciesClient
	defer func() { cliConnectivityPoliciesClient = originalCli }()
	cliConnectivityPoliciesClient = func(sessionContext utl.SessionContext, connector client.Connector) *transitgateways.ConnectivityPolicyClientContext {
		return connWrapper
	}

	res := resourceNsxtPolicyConnectivityPolicy()

	t.Run("Create_success", func(t *testing.T) {
		mockConnSDK.EXPECT().Patch(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, gomock.Any(), gomock.Any()).Return(nil)
		mockConnSDK.EXPECT().Get(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, gomock.Any()).Return(model.ConnectivityPolicy{
			DisplayName:       &connPolicyName,
			Path:              &connPolicyPath,
			Revision:          &connPolicyRevision,
			ConnectivityScope: &connPolicyScope,
			Group:             &connPolicyGroupPath,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":       connPolicyName,
			"parent_path":        connPolicyParentPath,
			"group_path":         connPolicyGroupPath,
			"connectivity_scope": connPolicyScope,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, connPolicyName, d.Get("display_name"))
	})

	t.Run("Create_fails_below_version_9_1_0", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "9.1.0" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
			"group_path":  connPolicyGroupPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})

	t.Run("Create_fails_when_parent_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": "/orgs/default",
			"group_path":  connPolicyGroupPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyCreate(d, m)
		require.Error(t, err)
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockConnSDK.EXPECT().Patch(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
			"group_path":  connPolicyGroupPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyConnectivityPolicyRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnSDK := connmocks.NewMockConnectivityPoliciesClient(ctrl)
	connWrapper := &transitgateways.ConnectivityPolicyClientContext{
		Client:     mockConnSDK,
		ClientType: utl.Multitenancy,
	}

	originalCli := cliConnectivityPoliciesClient
	defer func() { cliConnectivityPoliciesClient = originalCli }()
	cliConnectivityPoliciesClient = func(sessionContext utl.SessionContext, connector client.Connector) *transitgateways.ConnectivityPolicyClientContext {
		return connWrapper
	}

	res := resourceNsxtPolicyConnectivityPolicy()

	t.Run("Read_success", func(t *testing.T) {
		mockConnSDK.EXPECT().Get(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID).Return(model.ConnectivityPolicy{
			DisplayName:       &connPolicyName,
			Path:              &connPolicyPath,
			Revision:          &connPolicyRevision,
			ConnectivityScope: &connPolicyScope,
			Group:             &connPolicyGroupPath,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		d.SetId(connPolicyID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, connPolicyName, d.Get("display_name"))
	})

	t.Run("Read_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ConnectivityPolicy ID")
	})

	t.Run("Read_fails_when_Get_returns_not_found", func(t *testing.T) {
		mockConnSDK.EXPECT().Get(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID).Return(model.ConnectivityPolicy{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		d.SetId(connPolicyID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyConnectivityPolicyUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnSDK := connmocks.NewMockConnectivityPoliciesClient(ctrl)
	connWrapper := &transitgateways.ConnectivityPolicyClientContext{
		Client:     mockConnSDK,
		ClientType: utl.Multitenancy,
	}

	originalCli := cliConnectivityPoliciesClient
	defer func() { cliConnectivityPoliciesClient = originalCli }()
	cliConnectivityPoliciesClient = func(sessionContext utl.SessionContext, connector client.Connector) *transitgateways.ConnectivityPolicyClientContext {
		return connWrapper
	}

	res := resourceNsxtPolicyConnectivityPolicy()

	t.Run("Update_success", func(t *testing.T) {
		mockConnSDK.EXPECT().Update(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID, gomock.Any()).Return(model.ConnectivityPolicy{}, nil)
		mockConnSDK.EXPECT().Get(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID).Return(model.ConnectivityPolicy{
			DisplayName:       &connPolicyName,
			Path:              &connPolicyPath,
			Revision:          &connPolicyRevision,
			ConnectivityScope: &connPolicyScope,
			Group:             &connPolicyGroupPath,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": connPolicyName,
			"parent_path":  connPolicyParentPath,
			"group_path":   connPolicyGroupPath,
		})
		d.SetId(connPolicyID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyUpdate(d, m)
		require.Error(t, err)
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockConnSDK.EXPECT().Update(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID, gomock.Any()).Return(model.ConnectivityPolicy{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
			"group_path":  connPolicyGroupPath,
		})
		d.SetId(connPolicyID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyConnectivityPolicyDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnSDK := connmocks.NewMockConnectivityPoliciesClient(ctrl)
	connWrapper := &transitgateways.ConnectivityPolicyClientContext{
		Client:     mockConnSDK,
		ClientType: utl.Multitenancy,
	}

	originalCli := cliConnectivityPoliciesClient
	defer func() { cliConnectivityPoliciesClient = originalCli }()
	cliConnectivityPoliciesClient = func(sessionContext utl.SessionContext, connector client.Connector) *transitgateways.ConnectivityPolicyClientContext {
		return connWrapper
	}

	res := resourceNsxtPolicyConnectivityPolicy()

	t.Run("Delete_success", func(t *testing.T) {
		mockConnSDK.EXPECT().Delete(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		d.SetId(connPolicyID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_ID_is_empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockConnSDK.EXPECT().Delete(connPolicyOrgID, connPolicyProjectID, connPolicyTGWID, connPolicyID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": connPolicyParentPath,
		})
		d.SetId(connPolicyID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyConnectivityPolicyDelete(d, m)
		require.Error(t, err)
	})
}
