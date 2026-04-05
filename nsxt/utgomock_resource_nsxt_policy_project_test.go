//go:build unittest

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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	sdkprojects "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	orgsapi "github.com/vmware/terraform-provider-nsxt/api/orgs"
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	orgsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	projectID          = "test-project-id"
	projectDisplayName = "test-project"
	projectDescription = "Test Project"
	projectRevision    = int64(1)
	projectPath        = "/orgs/default/projects/test-project-id"
	ipv6BlockPath      = "/orgs/default/projects/default/ip-blocks/test-ipv6-block"
)

func projectAPIResponse() nsxModel.Project {
	return nsxModel.Project{
		Id:          &projectID,
		DisplayName: &projectDisplayName,
		Description: &projectDescription,
		Revision:    &projectRevision,
		Path:        &projectPath,
	}
}

func projectAPIResponseWithIPv6Blocks() nsxModel.Project {
	p := projectAPIResponse()
	p.Ipv6Blocks = []string{ipv6BlockPath}
	return p
}

func minimalProjectData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": projectDisplayName,
		"description":  projectDescription,
		"nsx_id":       projectID,
	}
}

func setupProjectMock(t *testing.T, ctrl *gomock.Controller) (*orgsmocks.MockProjectsClient, func()) {
	mockSDK := orgsmocks.NewMockProjectsClient(ctrl)
	mockWrapper := &orgsapi.ProjectClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliProjectsClient
	cliProjectsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *orgsapi.ProjectClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliProjectsClient = original }
}

// vpcSecurityProfilesClientStub implements sdkprojects.VpcSecurityProfilesClient for unit tests.
// Get returns NotFound so setVpcSecurityProfileInSchema exits without further API calls.
type vpcSecurityProfilesClientStub struct{}

func (vpcSecurityProfilesClientStub) Get(string, string, string) (nsxModel.VpcSecurityProfile, error) {
	return nsxModel.VpcSecurityProfile{}, vapiErrors.NotFound{}
}

func (vpcSecurityProfilesClientStub) List(string, string, *string, *bool, *string, *int64, *bool, *string) (nsxModel.VpcSecurityProfileListResult, error) {
	return nsxModel.VpcSecurityProfileListResult{}, nil
}

func (vpcSecurityProfilesClientStub) Patch(string, string, string, nsxModel.VpcSecurityProfile) error {
	return nil
}

func (vpcSecurityProfilesClientStub) Update(string, string, string, nsxModel.VpcSecurityProfile) (nsxModel.VpcSecurityProfile, error) {
	return nsxModel.VpcSecurityProfile{}, nil
}

var _ sdkprojects.VpcSecurityProfilesClient = vpcSecurityProfilesClientStub{}

func setupVpcSecurityProfilesStub(t *testing.T) func() {
	t.Helper()
	stub := vpcSecurityProfilesClientStub{}
	mockWrapper := &projects.VpcSecurityProfileClientContext{
		Client:     stub,
		ClientType: utl.Multitenancy,
	}
	original := cliVpcSecurityProfilesClient
	cliVpcSecurityProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *projects.VpcSecurityProfileClientContext {
		return mockWrapper
	}
	return func() { cliVpcSecurityProfilesClient = original }
}

func TestMockResourceNsxtPolicyProjectCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupProjectMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "4.1.0"
		defer func() { util.NsxVersion = "" }()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(nsxModel.Project{}, notFoundErr),
			mockSDK.EXPECT().Patch(utl.DefaultOrgID, projectID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(projectAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())

		err := resourceNsxtPolicyProjectCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectID, d.Id())
		assert.Equal(t, projectDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		util.NsxVersion = "4.1.0"
		defer func() { util.NsxVersion = "" }()

		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(projectAPIResponse(), nil)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())

		err := resourceNsxtPolicyProjectCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyProjectRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupProjectMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(projectAPIResponse(), nil)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())
		d.SetId(projectID)

		err := resourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectDisplayName, d.Get("display_name"))
		assert.Equal(t, projectDescription, d.Get("description"))
	})

	t.Run("Read sets ipv6_blocks when NSX 9.2.0+", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		restoreSec := setupVpcSecurityProfilesStub(t)
		defer restoreSec()

		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(projectAPIResponseWithIPv6Blocks(), nil)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())
		d.SetId(projectID)

		err := resourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		blocks := d.Get("ipv6_blocks").([]interface{})
		require.Len(t, blocks, 1)
		assert.Equal(t, ipv6BlockPath, blocks[0])
	})

	t.Run("Read does not populate ipv6_blocks when NSX below 9.2", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		restoreSec := setupVpcSecurityProfilesStub(t)
		defer restoreSec()

		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(projectAPIResponseWithIPv6Blocks(), nil)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())
		d.SetId(projectID)

		err := resourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		_, ok := d.GetOk("ipv6_blocks")
		assert.False(t, ok)
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(nsxModel.Project{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())
		d.SetId(projectID)

		err := resourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())

		err := resourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyProjectUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupProjectMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		util.NsxVersion = "4.1.0"
		defer func() { util.NsxVersion = "" }()

		gomock.InOrder(
			mockSDK.EXPECT().Patch(utl.DefaultOrgID, projectID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(projectAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())
		d.SetId(projectID)

		err := resourceNsxtPolicyProjectUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())

		err := resourceNsxtPolicyProjectUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyProjectDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupProjectMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(utl.DefaultOrgID, projectID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())
		d.SetId(projectID)

		err := resourceNsxtPolicyProjectDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectData())

		err := resourceNsxtPolicyProjectDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
