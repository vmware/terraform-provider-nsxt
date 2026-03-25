// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/LbPoolsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/LbPoolsClient.go LbPoolsClient

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
	lbPoolMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	lbPoolID          = "lb-pool-1"
	lbPoolDisplayName = "lb-pool-fooname"
	lbPoolDescription = "lb pool mock"
	lbPoolPath        = "/infra/lb-pools/lb-pool-1"
	lbPoolRevision    = int64(1)
	lbPoolAlgorithm   = model.LBPool_ALGORITHM_ROUND_ROBIN
	lbPoolMinActive   = int64(1)
)

func TestMockResourceNsxtPolicyLBPoolCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockPoolSDK.EXPECT().Get(gomock.Any()).Return(minimalLBPoolModel(), nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get("existing-id").Return(model.LBPool{Id: &lbPoolID}, nil)

		res := resourceNsxtPolicyLBPool()
		data := minimalLBPoolData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBPoolRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(minimalLBPoolModel(), nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, lbPoolDisplayName, d.Get("display_name"))
		assert.Equal(t, lbPoolDescription, d.Get("description"))
		assert.Equal(t, lbPoolPath, d.Get("path"))
		assert.Equal(t, int(lbPoolRevision), d.Get("revision"))
		assert.Equal(t, lbPoolID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBPool ID")
	})
}

func TestMockResourceNsxtPolicyLBPoolUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Update(lbPoolID, gomock.Any()).Return(minimalLBPoolModel(), nil)
		mockPoolSDK.EXPECT().Get(lbPoolID).Return(minimalLBPoolModel(), nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())
		d.Set("revision", 1)
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPoolData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBPool ID")
	})
}

func TestMockResourceNsxtPolicyLBPoolDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPoolSDK := lbPoolMocks.NewMockLbPoolsClient(ctrl)
	mockWrapper := &cliinfra.LBPoolClientContext{
		Client:     mockPoolSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbPoolsClient
	defer func() { cliLbPoolsClient = originalCli }()
	cliLbPoolsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBPoolClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockPoolSDK.EXPECT().Delete(lbPoolID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBPool ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockPoolSDK.EXPECT().Delete(lbPoolID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbPoolID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBPoolDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalLBPoolModel() model.LBPool {
	return model.LBPool{
		Id:               &lbPoolID,
		DisplayName:      &lbPoolDisplayName,
		Description:      &lbPoolDescription,
		Path:             &lbPoolPath,
		Revision:         &lbPoolRevision,
		Algorithm:        &lbPoolAlgorithm,
		MinActiveMembers: &lbPoolMinActive,
	}
}

func minimalLBPoolData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbPoolDisplayName,
		"description":  lbPoolDescription,
	}
}
