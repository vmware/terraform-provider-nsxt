// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/orgs/projects/vpcs/GatewayPoliciesClient.go -package=mocks -source=<sdk>/services/nsxt/orgs/projects/vpcs/GatewayPoliciesClient.go GatewayPoliciesClient
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
	vpcGwPolicyID          = "gw-policy-id"
	vpcGwPolicyDisplayName = "test-gw-policy"
	vpcGwPolicyDescription = "Test VPC Gateway Policy"
	vpcGwPolicyRevision    = int64(1)
)

func vpcGwPolicyAPIResponse() nsxModel.GatewayPolicy {
	return nsxModel.GatewayPolicy{
		Id:          &vpcGwPolicyID,
		DisplayName: &vpcGwPolicyDisplayName,
		Description: &vpcGwPolicyDescription,
		Revision:    &vpcGwPolicyRevision,
	}
}

func minimalVpcGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    vpcGwPolicyDisplayName,
		"description":     vpcGwPolicyDescription,
		"nsx_id":          vpcGwPolicyID,
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

func setupVpcGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*vpcmocks.MockGatewayPoliciesClient, *nsxtmocks.MockOrgRootClient, func()) {
	mockSDK := vpcmocks.NewMockGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.GatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	originalGwPolicy := cliGatewayPoliciesClient
	cliGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.GatewayPolicyClientContext {
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
		cliGatewayPoliciesClient = originalGwPolicy
		cliOrgRootClient = originalOrgRoot
	}
}

func TestMockResourceNsxtVpcGatewayPolicyCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, mockOrgRoot, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(nsxModel.GatewayPolicy{}, notFoundErr),
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(vpcGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())

		err := resourceNsxtVPCGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcGwPolicyID, d.Id())
		assert.Equal(t, vpcGwPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails on NSX version check", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())

		err := resourceNsxtVPCGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(vpcGwPolicyAPIResponse(), nil)

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())

		err := resourceNsxtVPCGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtVpcGatewayPolicyRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(vpcGwPolicyAPIResponse(), nil)

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())
		d.SetId(vpcGwPolicyID)

		err := resourceNsxtVPCGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcGwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcGwPolicyDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(nsxModel.GatewayPolicy{}, notFoundErr)

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())
		d.SetId(vpcGwPolicyID)

		err := resourceNsxtVPCGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())

		err := resourceNsxtVPCGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcGatewayPolicyUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, mockOrgRoot, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(vpcGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())
		d.SetId(vpcGwPolicyID)

		err := resourceNsxtVPCGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())

		err := resourceNsxtVPCGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcGatewayPolicyDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(nil)

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())
		d.SetId(vpcGwPolicyID)

		err := resourceNsxtVPCGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())

		err := resourceNsxtVPCGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Delete API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, _, restore := setupVpcGwPolicyMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcGwPolicyID).Return(errors.New("delete failed"))

		res := resourceNsxtVPCGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGwPolicyData())
		d.SetId(vpcGwPolicyID)

		err := resourceNsxtVPCGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
