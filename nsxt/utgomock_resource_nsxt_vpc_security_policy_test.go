// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/orgs/projects/vpcs/SecurityPoliciesClient.go -package=mocks -source=<sdk>/services/nsxt/orgs/projects/vpcs/SecurityPoliciesClient.go SecurityPoliciesClient
// mockgen -destination=mocks/nsxt/OrgRootClient.go -package=mocks github.com/vmware/vsphere-automation-sdk-go/services/nsxt OrgRootClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apiRoot "github.com/vmware/terraform-provider-nsxt/api"
	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxtmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsxt"
	vpcmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcSecurityPolicyID          = "security-policy-id"
	vpcSecurityPolicyDisplayName = "test-security-policy"
	vpcSecurityPolicyDescription = "Test VPC Security Policy"
	vpcSecurityPolicyRevision    = int64(1)
)

func vpcSecurityPolicyAPIResponse() nsxModel.SecurityPolicy {
	path := "/orgs/default/projects/project1/vpcs/vpc1/security-policies/security-policy-id"
	return nsxModel.SecurityPolicy{
		Id:          &vpcSecurityPolicyID,
		DisplayName: &vpcSecurityPolicyDisplayName,
		Description: &vpcSecurityPolicyDescription,
		Revision:    &vpcSecurityPolicyRevision,
		Path:        &path,
	}
}

func minimalVpcSecurityPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    vpcSecurityPolicyDisplayName,
		"description":     vpcSecurityPolicyDescription,
		"nsx_id":          vpcSecurityPolicyID,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  "project1",
				"vpc_id":      "vpc1",
				"from_global": false,
			},
		},
	}
}

func setupVpcSecurityPolicyMock(t *testing.T, ctrl *gomock.Controller) (*vpcmocks.MockSecurityPoliciesClient, *nsxtmocks.MockOrgRootClient, func()) {
	mockSDK := vpcmocks.NewMockSecurityPoliciesClient(ctrl)
	mockWrapper := &apidomains.SecurityPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	originalSecPolicy := cliSecurityPoliciesClient
	cliSecurityPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.SecurityPolicyClientContext {
		return mockWrapper
	}

	mockOrgRoot := nsxtmocks.NewMockOrgRootClient(ctrl)
	orgRootWrapper := &apiRoot.OrgRootClientContext{
		Client:     mockOrgRoot,
		ClientType: utl.Local,
	}
	originalOrgRoot := cliOrgRootClient
	cliOrgRootClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiRoot.OrgRootClientContext {
		return orgRootWrapper
	}

	return mockSDK, mockOrgRoot, func() {
		cliSecurityPoliciesClient = originalSecPolicy
		cliOrgRootClient = originalOrgRoot
	}
}

func TestMockResourceNsxtVpcSecurityPolicyCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, mockOrgRoot, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(nsxModel.SecurityPolicy{}, notFoundErr),
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(vpcSecurityPolicyAPIResponse(), nil),
		)

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())

		err := resourceNsxtVPCSecurityPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcSecurityPolicyID, d.Id())
		assert.Equal(t, vpcSecurityPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails on NSX version check", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())

		err := resourceNsxtVPCSecurityPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(vpcSecurityPolicyAPIResponse(), nil)

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())

		err := resourceNsxtVPCSecurityPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtVpcSecurityPolicyRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(vpcSecurityPolicyAPIResponse(), nil)

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())
		d.SetId(vpcSecurityPolicyID)

		err := resourceNsxtVPCSecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcSecurityPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcSecurityPolicyDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(nsxModel.SecurityPolicy{}, notFoundErr)

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())
		d.SetId(vpcSecurityPolicyID)

		err := resourceNsxtVPCSecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())

		err := resourceNsxtVPCSecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcSecurityPolicyUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, mockOrgRoot, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(vpcSecurityPolicyAPIResponse(), nil),
		)

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())
		d.SetId(vpcSecurityPolicyID)

		err := resourceNsxtVPCSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcSecurityPolicyDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(nil)

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())
		d.SetId(vpcSecurityPolicyID)

		err := resourceNsxtVPCSecurityPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())

		err := resourceNsxtVPCSecurityPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Delete API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcSecurityPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcSecurityPolicyID).Return(errors.New("delete failed"))

		res := resourceNsxtVPCSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSecurityPolicyData())
		d.SetId(vpcSecurityPolicyID)

		err := resourceNsxtVPCSecurityPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
