// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/SharesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/SharesClient.go SharesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	sharemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	shareID          = "share-1"
	shareDisplayName = "share-fooname"
	shareDescription = "share mock"
	sharePath        = "/infra/shares/share-1"
	shareRevision    = int64(1)
	shareStrategy    = model.Share_SHARING_STRATEGY_NONE_DESCENDANTS
	shareSharedWith  = []string{"/orgs/default"}
)

func TestMockResourceNsxtPolicyShareCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShareSDK := sharemocks.NewMockSharesClient(ctrl)
	mockWrapper := &cliinfra.ShareClientContext{
		Client:     mockShareSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSharesClient
	defer func() { cliSharesClient = originalCli }()
	cliSharesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ShareClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockShareSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockShareSDK.EXPECT().Get(gomock.Any()).Return(model.Share{
			Id:              &shareID,
			DisplayName:     &shareDisplayName,
			Description:     &shareDescription,
			Path:            &sharePath,
			Revision:        &shareRevision,
			SharingStrategy: &shareStrategy,
			SharedWith:      shareSharedWith,
		}, nil)

		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalShareData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockShareSDK.EXPECT().Get("existing-id").Return(model.Share{Id: &shareID}, nil)

		res := resourceNsxtPolicyShare()
		data := minimalShareData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyShareRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShareSDK := sharemocks.NewMockSharesClient(ctrl)
	mockWrapper := &cliinfra.ShareClientContext{
		Client:     mockShareSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSharesClient
	defer func() { cliSharesClient = originalCli }()
	cliSharesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ShareClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockShareSDK.EXPECT().Get(shareID).Return(model.Share{
			Id:              &shareID,
			DisplayName:     &shareDisplayName,
			Description:     &shareDescription,
			Path:            &sharePath,
			Revision:        &shareRevision,
			SharingStrategy: &shareStrategy,
			SharedWith:      shareSharedWith,
		}, nil)

		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(shareID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, shareDisplayName, d.Get("display_name"))
		assert.Equal(t, shareDescription, d.Get("description"))
		assert.Equal(t, sharePath, d.Get("path"))
		assert.Equal(t, int(shareRevision), d.Get("revision"))
		assert.Equal(t, shareID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining Share ID")
	})
}

func TestMockResourceNsxtPolicyShareUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShareSDK := sharemocks.NewMockSharesClient(ctrl)
	mockWrapper := &cliinfra.ShareClientContext{
		Client:     mockShareSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSharesClient
	defer func() { cliSharesClient = originalCli }()
	cliSharesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ShareClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockShareSDK.EXPECT().Patch(shareID, gomock.Any()).Return(nil)
		mockShareSDK.EXPECT().Get(shareID).Return(model.Share{
			Id:              &shareID,
			DisplayName:     &shareDisplayName,
			Description:     &shareDescription,
			Path:            &sharePath,
			Revision:        &shareRevision,
			SharingStrategy: &shareStrategy,
			SharedWith:      shareSharedWith,
		}, nil)

		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalShareData())
		d.SetId(shareID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalShareData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining Share ID")
	})
}

func TestMockResourceNsxtPolicyShareDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShareSDK := sharemocks.NewMockSharesClient(ctrl)
	mockWrapper := &cliinfra.ShareClientContext{
		Client:     mockShareSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSharesClient
	defer func() { cliSharesClient = originalCli }()
	cliSharesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ShareClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockShareSDK.EXPECT().Delete(shareID).Return(nil)

		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(shareID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining Share ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockShareSDK.EXPECT().Delete(shareID).Return(errors.New("API error"))

		res := resourceNsxtPolicyShare()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(shareID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyShareDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalShareData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     shareDisplayName,
		"description":      shareDescription,
		"sharing_strategy": shareStrategy,
		"shared_with":      []interface{}{"/orgs/default"},
	}
}
