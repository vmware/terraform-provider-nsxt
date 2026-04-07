// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/VniPoolsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/VniPoolsClient.go VniPoolsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vniPoolID          = "vni-pool-1"
	vniPoolDisplayName = "test-vni-pool"
	vniPoolDescription = "VNI pool description"
	vniPoolPath        = "/infra/vni-pools/vni-pool-1"
	vniPoolRevision    = int64(1)
	vniPoolStart       = int64(75001)
	vniPoolEnd         = int64(76000)
)

func vniPoolAPIResponse() model.VniPoolConfig {
	return model.VniPoolConfig{
		Id:          &vniPoolID,
		DisplayName: &vniPoolDisplayName,
		Description: &vniPoolDescription,
		Path:        &vniPoolPath,
		Revision:    &vniPoolRevision,
		Start:       &vniPoolStart,
		End:         &vniPoolEnd,
	}
}

func minimalVniPoolData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vniPoolDisplayName,
		"description":  vniPoolDescription,
		"start":        int(vniPoolStart),
		"end":          int(vniPoolEnd),
	}
}

func TestMockResourceNsxtPolicyVniPoolCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockVniPoolsClient(ctrl)
	mockWrapper := &cliinfra.VniPoolConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliVniPoolsClient
	defer func() { cliVniPoolsClient = originalCli }()
	cliVniPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.VniPoolConfigClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockSDK.EXPECT().Get(gomock.Any()).Return(vniPoolAPIResponse(), nil)

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVniPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get("existing-pool").Return(model.VniPoolConfig{Id: &vniPoolID}, nil)

		res := resourceNsxtPolicyVniPool()
		data := minimalVniPoolData()
		data["nsx_id"] = "existing-pool"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVniPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyVniPoolRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockVniPoolsClient(ctrl)
	mockWrapper := &cliinfra.VniPoolConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliVniPoolsClient
	defer func() { cliVniPoolsClient = originalCli }()
	cliVniPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.VniPoolConfigClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(vniPoolID).Return(vniPoolAPIResponse(), nil)

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vniPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vniPoolDisplayName, d.Get("display_name"))
		assert.Equal(t, vniPoolDescription, d.Get("description"))
		assert.Equal(t, vniPoolPath, d.Get("path"))
		assert.Equal(t, int(vniPoolRevision), d.Get("revision"))
		assert.Equal(t, int(vniPoolStart), d.Get("start"))
		assert.Equal(t, int(vniPoolEnd), d.Get("end"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining VNI Pool ID")
	})

	t.Run("Read clears ID when not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(vniPoolID).Return(model.VniPoolConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vniPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyVniPoolUpdate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockVniPoolsClient(ctrl)
	mockWrapper := &cliinfra.VniPoolConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliVniPoolsClient
	defer func() { cliVniPoolsClient = originalCli }()
	cliVniPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.VniPoolConfigClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(vniPoolID, gomock.Any()).Return(nil)
		mockSDK.EXPECT().Get(vniPoolID).Return(vniPoolAPIResponse(), nil)

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVniPoolData())
		d.SetId(vniPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, vniPoolDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVniPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining VNI Pool ID")
	})
}

func TestMockResourceNsxtPolicyVniPoolDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := inframocks.NewMockVniPoolsClient(ctrl)
	mockWrapper := &cliinfra.VniPoolConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliVniPoolsClient
	defer func() { cliVniPoolsClient = originalCli }()
	cliVniPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.VniPoolConfigClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(vniPoolID).Return(nil)

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vniPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining VNI Pool ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(vniPoolID).Return(errors.New("API error"))

		res := resourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vniPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVniPoolDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
