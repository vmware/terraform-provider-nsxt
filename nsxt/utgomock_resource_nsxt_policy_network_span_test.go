//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/NetworkSpansClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/NetworkSpansClient.go NetworkSpansClient

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
	spanmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	networkSpanID          = "network-span-1"
	networkSpanDisplayName = "network-span-fooname"
	networkSpanDescription = "network span mock"
	networkSpanPath        = "/infra/network-spans/network-span-1"
	networkSpanRevision    = int64(1)
)

func TestMockResourceNsxtPolicyNetworkSpanCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpanSDK := spanmocks.NewMockNetworkSpansClient(ctrl)
	mockWrapper := &cliinfra.NetworkSpanClientContext{
		Client:     mockSpanSDK,
		ClientType: utl.Local,
	}

	originalCli := cliNetworkSpansClient
	defer func() { cliNetworkSpansClient = originalCli }()
	cliNetworkSpansClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.NetworkSpanClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockSpanSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockSpanSDK.EXPECT().Get(gomock.Any()).Return(model.NetworkSpan{
			Id:          &networkSpanID,
			DisplayName: &networkSpanDisplayName,
			Description: &networkSpanDescription,
			Path:        &networkSpanPath,
			Revision:    &networkSpanRevision,
		}, nil)

		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNetworkSpanData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSpanSDK.EXPECT().Get("existing-id").Return(model.NetworkSpan{Id: &networkSpanID}, nil)

		res := resourceNsxtPolicyNetworkSpan()
		data := minimalNetworkSpanData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when NSX version below 9.1.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "9.1.0" }()

		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNetworkSpanData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})
}

func TestMockResourceNsxtPolicyNetworkSpanRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpanSDK := spanmocks.NewMockNetworkSpansClient(ctrl)
	mockWrapper := &cliinfra.NetworkSpanClientContext{
		Client:     mockSpanSDK,
		ClientType: utl.Local,
	}

	originalCli := cliNetworkSpansClient
	defer func() { cliNetworkSpansClient = originalCli }()
	cliNetworkSpansClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.NetworkSpanClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockSpanSDK.EXPECT().Get(networkSpanID).Return(model.NetworkSpan{
			Id:          &networkSpanID,
			DisplayName: &networkSpanDisplayName,
			Description: &networkSpanDescription,
			Path:        &networkSpanPath,
			Revision:    &networkSpanRevision,
		}, nil)

		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(networkSpanID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, networkSpanDisplayName, d.Get("display_name"))
		assert.Equal(t, networkSpanDescription, d.Get("description"))
		assert.Equal(t, networkSpanPath, d.Get("path"))
		assert.Equal(t, int(networkSpanRevision), d.Get("revision"))
		assert.Equal(t, networkSpanID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining NetworkSpan ID")
	})
}

func TestMockResourceNsxtPolicyNetworkSpanUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpanSDK := spanmocks.NewMockNetworkSpansClient(ctrl)
	mockWrapper := &cliinfra.NetworkSpanClientContext{
		Client:     mockSpanSDK,
		ClientType: utl.Local,
	}

	originalCli := cliNetworkSpansClient
	defer func() { cliNetworkSpansClient = originalCli }()
	cliNetworkSpansClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.NetworkSpanClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockSpanSDK.EXPECT().Update(networkSpanID, gomock.Any()).Return(model.NetworkSpan{
			Id:          &networkSpanID,
			DisplayName: &networkSpanDisplayName,
			Description: &networkSpanDescription,
			Path:        &networkSpanPath,
			Revision:    &networkSpanRevision,
		}, nil)
		mockSpanSDK.EXPECT().Get(networkSpanID).Return(model.NetworkSpan{
			Id:          &networkSpanID,
			DisplayName: &networkSpanDisplayName,
			Description: &networkSpanDescription,
			Path:        &networkSpanPath,
			Revision:    &networkSpanRevision,
		}, nil)

		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNetworkSpanData())
		d.Set("revision", 1)
		d.SetId(networkSpanID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNetworkSpanData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining NetworkSpan ID")
	})
}

func TestMockResourceNsxtPolicyNetworkSpanDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpanSDK := spanmocks.NewMockNetworkSpansClient(ctrl)
	mockWrapper := &cliinfra.NetworkSpanClientContext{
		Client:     mockSpanSDK,
		ClientType: utl.Local,
	}

	originalCli := cliNetworkSpansClient
	defer func() { cliNetworkSpansClient = originalCli }()
	cliNetworkSpansClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.NetworkSpanClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockSpanSDK.EXPECT().Delete(networkSpanID).Return(nil)

		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(networkSpanID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining NetworkSpan ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSpanSDK.EXPECT().Delete(networkSpanID).Return(errors.New("API error"))

		res := resourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(networkSpanID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyNetworkSpanDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalNetworkSpanData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": networkSpanDisplayName,
		"description":  networkSpanDescription,
	}
}
